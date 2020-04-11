package option

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"scullion/config"
)

type TaskOptions struct {
	EnvVar   string `short:"e" long:"env" default:"SCULLION_TASKS" description:"Load configuration from environment variable"`
	FileName string `short:"f" long:"file" description:"Read configuration from given file"`
	DryRun   bool   `long:"dry-run" description:"Perform a dry run and log actions that would be taken"`
}

func (o *TaskOptions) loadTaskDefs(data []byte) ([]config.TaskDef, error) {
	taskDefs := make([]config.TaskDef, 0)
	err := json.Unmarshal(data, &taskDefs)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal configuration: %w", err)
	}
	if len(taskDefs) == 0 {
		return nil, errors.New("please specify a configuration with tasks; nothing to do")
	}
	fmt.Printf("Loaded %d tasks\n", len(taskDefs))
	return taskDefs, nil
}

func (o *TaskOptions) ReadConfiguration() ([]config.TaskDef, error) {
	if o.EnvVar != "" {
		envValue := os.Getenv(o.EnvVar)
		if envValue != "" {
			if o.FileName != "" {
				return nil, errors.New("both configuration options specified; please choose one")
			}
			fmt.Printf("Using environment variable %s\n", o.EnvVar)
			return o.loadTaskDefs([]byte(envValue))
		}
	}
	if o.FileName != "" {
		data, err := ioutil.ReadFile(o.FileName)
		if err != nil {
			return nil, fmt.Errorf("unable to read configuration file: %w", err)
		}
		return o.loadTaskDefs(data)
	}
	return nil, fmt.Errorf("no configuration specified and variable '%s' is unset", o.EnvVar)
}
