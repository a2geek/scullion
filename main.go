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
	WorkerPools         `group:"Worker Pools" namespace:"worker" env-namespace:"WORKER"`
	Verbose             []bool `short:"v" long:"verbose" description:"enable verbose output"`
	CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}
type TaskOptions struct {
	EnvVar   func(string) `short:"e" long:"env" default:"SCULLION_TASKS" description:"load configuration from environment variable"`
	FileName func(string) `short:"f" long:"file" description:"read configuration from given file"`
}
type WorkerPools struct {
	OrgPool   int `long:"org-pool" env:"ORG_POOL" default:"1" description:"set the number of organization workers in the pool"`
	SpacePool int `long:"space-pool" env:"SPACE_POOL" default:"1" description:"set the number of space workers in the pool"`
	AppPool   int `long:"app-pool" env:"APP_POOL" default:"1" description:"set the number of application workers in the pool"`
}
type CloudFoundryOptions struct {
	API               string `short:"a" long:"api" env:"API" description:"API URL"`
	Username          string `short:"u" long:"username" env:"USERNAME" description:"Username"`
	Password          string `short:"p" long:"password" env:"PASSWORD" description:"Password"`
	SkipSslValidation bool   `short:"k" long:"skip-ssl-validation" env:"SKIP_SSL_VALIDATION" description:"Skip SSL validation of Cloud Foundry endpoint. Not recommended."`
}

var tasks []Task

func main() {
	var options Options
	options.TaskOptions.EnvVar = readConfigurationFromString
	options.TaskOptions.FileName = readConfigurationFromFile
	parser := flags.NewParser(&options, flags.Default)
	parser.NamespaceDelimiter = "-"
	_, err := parser.Parse()
	if err != nil {
		if !flags.WroteHelp(err) {
			fmt.Println(err)
		}
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

	// 3 worker pools, configured by env/flags for size:
	//   <-(Filter, Org)
	//   <-(Filter, Space)
	//   <-(Filter, App)
	// Per task worker:
	//   <-(Tick) and delivers to proper pool based on starting filter
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
			fmt.Printf("Using environment variable %s\n", envName)
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
