package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/cloudfoundry-community/go-cfclient"
)

func Validate(tasks []Task) bool {
	fails := 0

	var org cfclient.Org
	err := json.Unmarshal([]byte(orgJSON), &org)
	if err != nil {
		fmt.Printf("Unable to parse org payload for validation: %s\n", err)
		fails += 1
	}

	var space cfclient.Space
	err = json.Unmarshal([]byte(spaceJSON), &space)
	if err != nil {
		fmt.Printf("Unable to parse space payload for validation: %s\n", err)
		fails += 1
	}

	var app cfclient.App
	err = json.Unmarshal([]byte(applicationJSON), &app)
	if err != nil {
		fmt.Printf("Unable to parse app payload for validation: %s\n", err)
		fails += 1
	}

	orgVar := TaskVariables{
		Org: org,
	}
	spaceVar := TaskVariables{
		Org:   org,
		Space: space,
	}
	appVar := TaskVariables{
		Org:   org,
		Space: space,
		App:   app,
	}

	for _, task := range tasks {
		m, err := createMetadata(task, nil)
		if err != nil {
			fmt.Printf("Unable to compile expressions for task '%s': %s\n", task.Name, err)
			fails += 1
			continue
		}

		fails += validateExpression(task.Name, "org", m.OrgExpr, orgVar)
		fails += validateExpression(task.Name, "space", m.SpaceExpr, spaceVar)
		fails += validateExpression(task.Name, "app", m.AppExpr, appVar)

		if _, err := time.ParseDuration(task.Schedule.Frequency); err != nil {
			fmt.Printf("Unable to evaluate task '%s' frequency: %s\n", task.Name, err)
			fails += 1
			continue
		}
	}
	return fails == 0
}

func validateExpression(taskName, testName string, pgm *vm.Program, vars TaskVariables) int {
	result, err := expr.Run(pgm, vars)
	if err != nil {
		fmt.Printf("Unable to evaluate task '%s' %s expression: %s\n", taskName, testName, err)
		return 1
	}
	_, err = IsTrue(result)
	if err != nil {
		fmt.Printf("Unable to evaluate task '%s' %s expression: %s\n", taskName, testName, err)
		return 1
	}
	return 0
}
