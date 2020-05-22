package ctx

import (
	"scullion/log"
)

type Subprogram struct {
	Name     string
	Pipeline []string
	Actions  []string
	Logger   log.Logger
}

func NewSubprogram(name string, pipeline, actions []string, logger log.Logger) (Subprogram, error) {
	sub := Subprogram{
		Name:     name,
		Pipeline: pipeline,
		Actions:  actions,
		Logger:   logger,
	}
	// FIXME: We used to do stuff. Leaving error just in case we do stuff again!
	return sub, nil
}

func (sub *Subprogram) Dup() Subprogram {
	newP := make([]string, len(sub.Pipeline))
	copy(newP, sub.Pipeline)
	newA := make([]string, len(sub.Actions))
	copy(newA, sub.Actions)
	return Subprogram{
		Name:     sub.Name,
		Pipeline: newP,
		Actions:  newA,
		Logger:   sub.Logger,
	}
}
