package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func NewOptions(tasks *[]Task) Options {
	options = Options{
		tasks: tasks,
	}
	options.TaskOptions.EnvVar = options.readConfigurationFromString
	options.TaskOptions.FileName = options.readConfigurationFromFile
	return options
}

var options Options

type Options struct {
	tasks *[]Task

	TaskOptions         `group:"Task Options"`
	WorkerPools         `group:"Worker Pools" namespace:"worker" env-namespace:"WORKER"`
	Verbose             []bool `short:"v" long:"verbose" description:"enable verbose output"`
	CloudFoundryOptions `group:"Cloud Foundry Configuration" namespace:"cf" env-namespace:"CF" reqired:"yes"`
}
type TaskOptions struct {
	EnvVar   func(string) `short:"e" long:"env" default:"SCULLION_TASKS" description:"load configuration from environment variable"`
	FileName func(string) `short:"f" long:"file" description:"read configuration from given file"`
}
type WorkerPools struct {
	OrgPool    int `long:"org-pool" env:"ORG_POOL" default:"1" description:"set the number of organization workers in the pool"`
	SpacePool  int `long:"space-pool" env:"SPACE_POOL" default:"1" description:"set the number of space workers in the pool"`
	AppPool    int `long:"app-pool" env:"APP_POOL" default:"1" description:"set the number of application workers in the pool"`
	ActionPool int `long:"action-pool" env:"ACTION_POOL" default:"1" description:"set the number of action (stop/start) workers in the pool"`
}
type CloudFoundryOptions struct {
	API               string `short:"a" long:"api" env:"API" description:"API URL"`
	Username          string `short:"u" long:"username" env:"USERNAME" description:"Username"`
	Password          string `short:"p" long:"password" env:"PASSWORD" description:"Password"`
	SkipSslValidation bool   `short:"k" long:"skip-ssl-validation" env:"SKIP_SSL_VALIDATION" description:"Skip SSL validation of Cloud Foundry endpoint. Not recommended."`
}

func (o *Options) loadTasks(config []byte) {
	if len(*options.tasks) != 0 {
		log.Fatal("please do not load tasks twice")
	}
	err := json.Unmarshal(config, &options.tasks)
	if err != nil {
		fmt.Printf("Unable to unmarshal configuration: %v\n", err)
		os.Exit(3)
	}
	fmt.Printf("Loaded %d tasks\n", len(*options.tasks))
}

func (o *Options) readConfigurationFromString(envName string) {
	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			fmt.Printf("[%s]\n", envValue)
			fmt.Printf("Using environment variable %s\n", envName)
			o.loadTasks([]byte(envValue))
		}
	}
}

func (o *Options) readConfigurationFromFile(fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to read configuration: %v\n", err)
		os.Exit(2)
	}
	o.loadTasks(data)
}
