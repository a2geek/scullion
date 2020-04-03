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
	TaskOptions         `group:"Task Options" required:"yes"`
	Verbose             []bool `short:"v" long:"verbose" description:"enable verbose output"`
	CloudFoundryOptions `group:"Cloud Foundry Configuration"`
}
type TaskOptions struct {
	EnvVar   func(string) `short:"e" long:"env" default:"SCULLION_TASKS" description:"load configuration from environment variable"`
	FileName func(string) `short:"f" long:"file" description:"read configuration from given file"`
}
type CloudFoundryOptions struct {
	API               string `short:"a" long:"api" env:"CF_API" description:"API URL" required:"yes"`
	Username          string `short:"u" long:"username" env:"CF_USERNAME" description:"Username" required:"yes"`
	Password          string `short:"p" long:"password" env:"CF_PASSWORD" description:"Password" required:"yes"`
	SkipSslValidation bool   `short:"k" long:"skip-ssl-validation" env:"CF_SKIP_SSL_VALIDATION" description:"Skip SSL validation of Cloud Foundry endpoint. Not recommended."`
}

var tasks []Task

func main() {
	var options Options
	options.TaskOptions.EnvVar = readConfigurationFromString
	options.TaskOptions.FileName = readConfigurationFromFile
	parser := flags.NewParser(&options, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("Please supply a file or an environment variable for configuration")
		os.Exit(1)
	}

	c := &cfclient.Config{
		ApiAddress:        options.API,
		Username:          options.Username,
		Password:          options.Password,
		SkipSslValidation: options.SkipSslValidation,
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
		envValue := os.Getenv(envName)
		if envValue != "" {
			fmt.Printf("Using environment variable %s", envName)
			loadTasks([]byte(envValue))
		}
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
