package fn

import (
	"fmt"
	"scullion/ctx"
	"scullion/log"
)

func NewTemplateRegistrar(templates map[string]string) Registrar {
	return func(state *ctx.State) {
		tmpl := Template{
			templates: templates,
			logger:    state.LoggerWrapper(),
		}
		state.AddFunc("Template", tmpl.Template)
	}
}

type Template struct {
	templates map[string]string
	logger    log.Logger
}

func (t Template) Template(name string, parameters ...string) string {
	if t, ok := t.templates[name]; ok {
		return fmt.Sprintf(t, parameters)
	}
	t.logger.Errorf("template '%s' not found", name)
	return "template not found"
}
