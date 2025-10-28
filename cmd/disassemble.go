package cmd

import (
	"fmt"
	"scullion/action"
	"scullion/option"
	"scullion/task"
)

type Disassemble struct {
	option.TaskOptions `group:"Task Options"`
}

func (cmd *Disassemble) Execute(args []string) error {
	taskDefs, err := cmd.ReadConfiguration()
	if err != nil {
		return err
	}

	runOpts := option.RunOptions{
		DryRun: false,
		Level:  "INFO",
		NoDate: false,
	}
	for _, taskDef := range taskDefs {
		m, err := task.NewMetadata(taskDef, nil, action.Log, runOpts)
		if err != nil {
			fmt.Printf("Unable to compile expressions for task '%s': %s\n", taskDef.Name, err)
			continue
		}

		fmt.Printf("['%s' : org (src={%s})]\n%s\n\n", taskDef.Name, taskDef.Filters.Organization, m.OrgExpr.Disassemble())
		fmt.Printf("['%s' : space (src={%s})]\n%s\n\n", taskDef.Name, taskDef.Filters.Space, m.SpaceExpr.Disassemble())
		fmt.Printf("['%s' : app (src={%s})]\n%s\n\n", taskDef.Name, taskDef.Filters.Application, m.AppExpr.Disassemble())
	}
	return nil
}
