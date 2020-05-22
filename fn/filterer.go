package fn

import "scullion/ctx"

func NewFiltererRegistrar() Registrar {
	return func(state *ctx.State) {
		f := Filterer{
			state: state,
		}
		state.AddFunc("Filter", f.Filter)
	}
}

type Filterer struct {
	state *ctx.State
}

func (f *Filterer) Filter(flag bool) error {
	if flag {
		f.state.ReEmit()
	}
	return nil
}
