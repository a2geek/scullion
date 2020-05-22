package fn

import (
	"fmt"
	"scullion/config"
	"scullion/ctx"
)

func NewLibraryRegistrar(libraryDefs []config.LibraryDef) Registrar {
	return func(state *ctx.State) {
		library := make(map[string]ctx.Subprogram)
		for _, libraryDef := range libraryDefs {
			pgm, err := ctx.NewSubprogram(libraryDef.Name, libraryDef.Pipeline, libraryDef.Actions, state.LoggerWrapper())
			if err != nil {
				state.Errorf("Error with subprogram '%s': %v", libraryDef.Name, err)
			}
			library[libraryDef.Name] = pgm
		}

		lib := Library{
			state:   state,
			library: library,
		}
		state.AddFunc("Call", lib.Call)
	}
}

type Library struct {
	state   *ctx.State
	library map[string]ctx.Subprogram
}

func (l *Library) Call(name string) error {
	p, ok := l.library[name]
	if !ok {
		return fmt.Errorf("unable to locate subprogram '%s' in library", name)
	}
	newPgm := p.Dup()
	l.state.EmitSub(newPgm)
	return nil
}
