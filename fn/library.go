package fn

import (
	"scullion/ctx"

	"github.com/lxc/lxd/shared/logger"
)

func NewLibraryRegistrar(library map[string]ctx.Subprogram) Registrar {
	return func(state *ctx.State) {
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

func (l *Library) Call(name string) {
	p, ok := l.library[name]
	if !ok {
		logger.Errorf("unable to locate subprogram '%s' in library", name)
		return
	}
	newPgm := p.Dup()
	l.state.EmitSub(newPgm)
}
