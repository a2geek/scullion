package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"scullion/action"
	"scullion/config"
	"scullion/option"
	"scullion/payload"
	"scullion/task"
	"time"

	"github.com/cloudfoundry-community/go-cfclient"
)

type Validate struct {
	option.TaskOptions `group:"Task Options"`
}

func (cmd *Validate) Execute(args []string) error {
	taskDefs, err := cmd.ReadConfiguration()
	if err != nil {
		return err
	}

	if !cmd.validate(taskDefs) {
		fmt.Println("Some tasks and schedules DID NOT pass validation.")
		os.Exit(1)
	} else {
		fmt.Println("All tasks and schedules passed validation.")
	}

	return nil
}

func (cmd *Validate) validate(taskDefs []config.TaskDef) bool {
	totalFails := 0

	var org cfclient.Org
	err := json.Unmarshal([]byte(payload.OrgJSON), &org)
	if err != nil {
		fmt.Printf("Unable to parse org payload for validation: %s\n", err)
		totalFails++
	}

	var space cfclient.Space
	err = json.Unmarshal([]byte(payload.SpaceJSON), &space)
	if err != nil {
		fmt.Printf("Unable to parse space payload for validation: %s\n", err)
		totalFails++
	}

	var app cfclient.App
	err = json.Unmarshal([]byte(payload.ApplicationJSON), &app)
	if err != nil {
		fmt.Printf("Unable to parse app payload for validation: %s\n", err)
		totalFails++
	}

	orgVar := task.Variables{
		Org: org,
	}
	spaceVar := task.Variables{
		Org:   org,
		Space: space,
	}
	appVar := task.Variables{
		Org:   org,
		Space: space,
		App:   app,
	}

	for _, taskDef := range taskDefs {
		fails := 0
		m, err := task.NewMetadata(taskDef, nil, action.Log)
		if err != nil {
			fmt.Printf("Unable to compile expressions for task '%s': %s\n", taskDef.Name, err)
			fails++
		} else {
			fails += validateExpression(taskDef.Name, "org", m.IsOrgMatch, orgVar)
			fails += validateExpression(taskDef.Name, "space", m.IsSpaceMatch, spaceVar)
			fails += validateExpression(taskDef.Name, "app", m.IsAppMatch, appVar)
		}

		if _, err := time.ParseDuration(taskDef.Schedule.Frequency); err != nil {
			fmt.Printf("Unable to evaluate task '%s' frequency: %s\n", taskDef.Name, err)
			fails++
		}

		totalFails += fails
		if fails == 0 {
			fmt.Printf("Task '%s' looks ok.\n", taskDef.Name)
		} else {
			fmt.Printf("Task '%s' had errors.\n", taskDef.Name)
		}
	}
	return totalFails == 0
}

func validateExpression(taskName, testName string, pgm func(task.Variables) (bool, error), vars task.Variables) int {
	_, err := pgm(vars)
	if err != nil {
		fmt.Printf("Unable to evaluate task '%s' %s expression: %s\n", taskName, testName, err)
		return 1
	}
	return 0
}
