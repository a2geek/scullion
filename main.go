package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	envName := flag.String("env", "SCULLION_TASKS", "load configuration from environment `variable`")
	fileName := flag.String("file", "", "read configuration from given `file`")
	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	readConfigFunc := evaluateFlags(envName, fileName, verbose)

	config := readConfigFunc()
	fmt.Println(string(config))
	tasks := loadTasks(config)
	fmt.Printf("%v\n", tasks)
}

func evaluateFlags(envName, fileName *string, verbose *bool) func() []byte {
	var readConfigFunc func() []byte
	if *fileName != "" {
		readConfigFunc = func() []byte { return readConfigurationFromFile(*fileName) }
	} else if *envName != "" {
		readConfigFunc = func() []byte { return []byte(os.Getenv(*envName)) }
	} else {
		fmt.Println("Please supply a file or an environment variable for configuration")
		flag.Usage()
		os.Exit(1)
	}
	return readConfigFunc
}

func loadTasks(config []byte) []Task {
	tasks := make([]Task, 0)
	err := json.Unmarshal(config, &tasks)
	if err != nil {
		fmt.Printf("Unable to unmarshal configuration: %v\n", err)
		os.Exit(3)
	}
	return tasks
}

func readConfigurationFromFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to read configuration: %v\n", err)
		os.Exit(2)
	}
	return data
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
