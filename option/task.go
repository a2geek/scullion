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
}

func (o *TaskOptions) loadRuleDefs(data []byte) (config.Config, error) {
	cfg := config.Config{}
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("unable to unmarshal configuration: %w", err)
	}
	if len(cfg.Rules) == 0 {
		return cfg, errors.New("please specify a configuration with tasks; nothing to do")
	}
	fmt.Printf("Loaded %d rules, %d templates, and library of %d subprograms\n", len(cfg.Rules), len(cfg.Templates), len(cfg.Library))
	return cfg, nil
}

func (o *TaskOptions) ReadConfiguration() (config.Config, error) {
	if o.EnvVar != "" {
		envValue := os.Getenv(o.EnvVar)
		if envValue != "" {
			if o.FileName != "" {
				return config.Config{}, errors.New("both configuration options specified; please choose one")
			}
			fmt.Printf("Using environment variable %s\n", o.EnvVar)
			return o.loadRuleDefs([]byte(envValue))
		}
	}
	if o.FileName != "" {
		data, err := ioutil.ReadFile(o.FileName)
		if err != nil {
			return config.Config{}, fmt.Errorf("unable to read configuration file: %w", err)
		}
		return o.loadRuleDefs(data)
	}
	return config.Config{}, fmt.Errorf("no configuration specified and variable '%s' is unset", o.EnvVar)
}
