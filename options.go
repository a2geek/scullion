package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func NewOptions() Options {
	options = Options{}
	options.TaskOptions.EnvVar = readConfigurationFromString
	options.TaskOptions.FileName = readConfigurationFromFile
	return options
}

var options Options

type Options struct {
	Tasks []Task `no-flag:"yes"`

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

func loadTasks(config []byte) {
	options.Tasks = make([]Task, 0)
	err := json.Unmarshal(config, &options.Tasks)
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
