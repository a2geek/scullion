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
		pgm, err := compile(code)
		if err != nil {
			return sub, fmt.Errorf("entry #%d: %w", i, err)
		}
		sub.Pipeline = append(sub.Pipeline, pgm)
	}
	for i, code := range actions {
		pgm, err := compile(code)
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

func compile(code string) (*vm.Program, error) {
	options := []expr.Option{
		//		expr.Env(vars),

		// Operators override for date comprising.
		expr.Operator("==", "Equal"),
		expr.Operator("<", "Before"),
		expr.Operator("<=", "BeforeOrEqual"),
		expr.Operator(">", "After"),
		expr.Operator(">=", "AfterOrEqual"),

		// Time and duration manipulation.
		expr.Operator("+", "Add"),
		expr.Operator("-", "Sub"),

		// Operators override for duration comprising.
		expr.Operator("==", "EqualDuration"),
		expr.Operator("<", "BeforeDuration"),
		expr.Operator("<=", "BeforeOrEqualDuration"),
		expr.Operator(">", "AfterDuration"),
		expr.Operator(">=", "AfterOrEqualDuration"),
	}

	pgm, err := expr.Compile(code, options...)
	if err != nil {
		return nil, err
	}
	return pgm, nil
}
