package cmd

import (
	"scullion/option"
)

type Disassemble struct {
	option.TaskOptions `group:"Task Options"`
}

func (cmd *Disassemble) Execute(args []string) error {
	// cfg, err := cmd.ReadConfiguration()
	// if err != nil {
	// 	return err
	// }

	// runOpts := option.RunOptions{
	// 	DryRun: false,
	// 	Level:  "INFO",
	// 	NoDate: false,
	// }
	// for _, ruleDef := range cfg.Rules {
	// 	m, err := task.NewMetadata(ruleDef, nil, action.Log, runOpts)
	// 	if err != nil {
	// 		fmt.Printf("Unable to compile expressions for task '%s': %s\n", ruleDef.Name, err)
	// 		continue
	// 	}

	// 	fmt.Printf("['%s' : org]\n%s\n\n", ruleDef.Name, m.OrgExpr.Disassemble())
	// 	fmt.Printf("['%s' : space]\n%s\n\n", ruleDef.Name, m.SpaceExpr.Disassemble())
	// 	fmt.Printf("['%s' : app]\n%s\n\n", ruleDef.Name, m.AppExpr.Disassemble())
	// }
	return nil
}
