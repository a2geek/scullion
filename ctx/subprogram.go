package ctx

import (
	"fmt"
	"scullion/log"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

type Subprogram struct {
	Name     string
	Pipeline []*vm.Program
	Actions  []*vm.Program
	Logger   log.Logger
}

func NewSubprogram(name string, pipeline, actions []string, logger log.Logger) (Subprogram, error) {
	sub := Subprogram{
		Name:   name,
		Logger: logger,
	}
	for i, code := range pipeline {
		pgm, err := expr.Compile(code, expr.AllowUndefinedVariables())
		if err != nil {
			return sub, fmt.Errorf("entry #%d: %w", i, err)
		}
		sub.Pipeline = append(sub.Pipeline, pgm)
	}
	for i, code := range actions {
		pgm, err := expr.Compile(code, expr.AllowUndefinedVariables())
		if err != nil {
			return sub, fmt.Errorf("entry #%d: %w", i, err)
		}
		sub.Actions = append(sub.Actions, pgm)
	}
	return sub, nil
}

func (sub *Subprogram) Dup() Subprogram {
	newP := make([]*vm.Program, len(sub.Pipeline))
	copy(newP, sub.Pipeline)
	newA := make([]*vm.Program, len(sub.Actions))
	copy(newA, sub.Actions)
	return Subprogram{
		Name:     sub.Name,
		Pipeline: newP,
		Actions:  newA,
		Logger:   sub.Logger,
	}
}
