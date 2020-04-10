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

	for _, taskDef := range taskDefs {
		m, err := task.NewMetadata(taskDef, nil, action.Log)
		if err != nil {
			fmt.Printf("Unable to compile expressions for task '%s': %s\n", taskDef.Name, err)
			continue
		}

		fmt.Printf("['%s' : org]\n%s\n\n", taskDef.Name, m.OrgExpr.Disassemble())
		fmt.Printf("['%s' : space]\n%s\n\n", taskDef.Name, m.SpaceExpr.Disassemble())
		fmt.Printf("['%s' : app]\n%s\n\n", taskDef.Name, m.AppExpr.Disassemble())
	}
	return nil
}
