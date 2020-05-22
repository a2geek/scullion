package ctx

import (
	"scullion/log"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
)

// State stores the current execution context for this rule
type State struct {
	vars map[string]interface{}
	code []Subprogram
	ch   chan *State
}

func NewState(subpgm Subprogram, stateChan chan *State) *State {
	return &State{
		vars: make(map[string]interface{}),
		code: []Subprogram{subpgm},
		ch:   stateChan,
	}
}

func (s *State) AddFunc(n string, fn interface{}) {
	s.vars[n] = fn
}

func (ctx *State) Dup() *State {
	// Clone existing code
	nextCode := make([]Subprogram, len(ctx.code))
	copy(nextCode, ctx.code)
	// Clone existing variables
	nextVars := make(map[string]interface{})
	for k, v := range ctx.vars {
		nextVars[k] = v
	}
	return &State{
		vars: nextVars,
		code: nextCode,
		ch:   ctx.ch,
	}
}

func (ctx *State) Execute() {
	sub := ctx.code[len(ctx.code)-1]
	if len(sub.Pipeline) == 0 {
		// Reached end of Pipeline, need to execute _all_ Actions
		actions := sub.Actions
		sub.Actions = nil
		for _, action := range actions {
			nextCtx := ctx.Dup()
			nextCtx.code[len(nextCtx.code)-1].Pipeline = []*vm.Program{action}
			ctx.ch <- nextCtx
		}
		// If there are prior Pipelines to run, kick it off
		if len(ctx.code) > 0 {
			nextCtx := ctx.Dup()
			nextCtx.code = nextCtx.code[1:]
			ctx.ch <- nextCtx
		}
	} else {
		// In the Pipeline, each step determines if the pipeline continues via calling Emit
		pgm := sub.Pipeline[0]
		sub.Pipeline = sub.Pipeline[1:]
		x, err := expr.Compile(pgm.Source.Content(),
			expr.Env(ctx.vars),

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
		)
		if err != nil {
			ctx.Errorf("compile: %v", err)
			return
		}

		out, err := expr.Run(x, ctx.vars)
		if err != nil {
			ctx.Errorf("run: %v", err)
			return
		}
		if e, ok := out.(error); ok && e != nil {
			ctx.Errorf("run: %v", e)
			return
		}
	}
}

func (ctx *State) EmitVar(newVars map[string]interface{}) {
	nextCtx := ctx.Dup()
	for k, v := range newVars {
		nextCtx.vars[k] = v
	}
	ctx.ch <- nextCtx
}

func (ctx *State) EmitSub(pgm Subprogram) {
	nextCtx := ctx.Dup()
	nextCtx.code = append(nextCtx.code, pgm)
	ctx.ch <- nextCtx
}

func (ctx *State) ReEmit() {
	ctx.ch <- ctx
}

func (s *State) LoggerWrapper() log.Logger              { return s }
func (s *State) Debug(v ...interface{})                 { s.activeLogger().Debug(v...) }
func (s *State) Debugf(format string, v ...interface{}) { s.activeLogger().Debugf(format, v...) }
func (s *State) Info(v ...interface{})                  { s.activeLogger().Info(v...) }
func (s *State) Infof(format string, v ...interface{})  { s.activeLogger().Infof(format, v...) }
func (s *State) Warn(v ...interface{})                  { s.activeLogger().Warn(v...) }
func (s *State) Warnf(format string, v ...interface{})  { s.activeLogger().Warnf(format, v...) }
func (s *State) Error(v ...interface{})                 { s.activeLogger().Error(v...) }
func (s *State) Errorf(format string, v ...interface{}) { s.activeLogger().Errorf(format, v...) }

func (s *State) activeLogger() log.Logger {
	return s.activeSubprogram().Logger
}

func (s *State) activeSubprogram() Subprogram {
	return s.code[len(s.code)-1]
}
