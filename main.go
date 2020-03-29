package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var envName, fileName string
var verbose bool
var tasks []Task

func main() {
	flag.StringVar(&envName, "env", "SCULLION_TASKS", "load configuration from environment `variable`")
	flag.StringVar(&fileName, "file", "", "read configuration from given `file`")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose output")
	flag.Parse()

	config := readConfiguration(envName, fileName)
	fmt.Println(string(config))
	tasks = make([]Task, 0)
	err := json.Unmarshal(config, &tasks)
	if err != nil {
		fmt.Printf("Unable to unmarshal configuration: %v\n", err)
		os.Exit(3)
	}
	fmt.Printf("%v\n", tasks)
}

func readConfiguration(envName string, fileName string) []byte {
	var config []byte
	if fileName != "" {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Printf("Unable to read configuration: %v\n", err)
			os.Exit(2)
		}
		config = data
	} else if envName != "" {
		config = []byte(os.Getenv(envName))
	} else {
		fmt.Println("Please supply a file or an environment variable for configuration")
		flag.Usage()
		os.Exit(1)
	}
	return config
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
