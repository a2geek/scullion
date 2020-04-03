package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	TaskEnvVar          func(string) `short:"e" long:"env" default:"SCULLION_TASKS" description:"load configuration from environment variable"`
	TaskFileName        func(string) `short:"f" long:"file" description:"read configuration from given file"`
	Verbose             []bool       `short:"v" long:"verbose" description:"enable verbose output"`
	CfAPI               string       `short:"a" long:"api" env:"CF_API" description:"Cloud Foundry API URL" required:"yes"`
	CfUsername          string       `short:"u" long:"username" env:"CF_USERNAME" description:"Cloud Foundry username" required:"yes"`
	CfPassword          string       `short:"p" long:"password" env:"CF_PASSWORD" description:"Cloud Foundry password" required:"yes"`
	CfSkipSslValidation bool         `short:"k" long:"skip-ssl-validation" env:"CF_SKIP_SSL_VALIDATION" description:"Skip SSL validation of Cloud Foundry endpoint. Not recommended."`
}

var tasks []Task

func main() {
	var options Options
	options.TaskEnvVar = readConfigurationFromString
	options.TaskFileName = readConfigurationFromFile
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(0)
	}

	fmt.Printf("%v\n", tasks)

	// fmt.Println("Please supply a file or an environment variable for configuration")
	// flag.Usage()
	// os.Exit(1)

	c := &cfclient.Config{
		ApiAddress:        options.CfAPI,
		Username:          options.CfUsername,
		Password:          options.CfPassword,
		SkipSslValidation: options.CfSkipSslValidation,
	}
	client, err := cfclient.NewClient(c)
	if err != nil {
		panic(err)
	}
	q := url.Values{}
	apps, err := client.ListAppsByQuery(q)
	if err != nil {
		panic(err)
	}
	for _, app := range apps {
		fmt.Printf("Name: %s (%d instances)\n", app.Name, app.Instances)
	}
}

func loadTasks(config []byte) {
	tasks = make([]Task, 0)
	err := json.Unmarshal(config, &tasks)
	if err != nil {
		fmt.Printf("Unable to unmarshal configuration: %v\n", err)
		os.Exit(3)
	}
}

func readConfigurationFromString(envName string) {
	if envName != "" {
		loadTasks([]byte(os.Getenv(envName)))
	}
}

func readConfigurationFromFile(fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to read configuration: %v\n", err)
		os.Exit(2)
	}
	loadTasks(data)
}

type Task struct {
	Schedule Schedule `json:"schedule"`
	Filters  Filter   `json:"filters"`
}
type Schedule struct {
	Frequency string `json:"frequency"`
}
type Filter struct {
	Organization string `json:"organization"`
	Space        string `json:"space"`
	Application  string `json:"application"`
	Action       string `json:"action"`
}
