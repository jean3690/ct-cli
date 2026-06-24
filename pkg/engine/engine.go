package engine

import "text/template"

// TemplateEngine defines the interface for rendering templates.
type TemplateEngine interface {
	Render(tmplContent string, vars map[string]string) (string, error)
	RenderFile(tmplPath, destPath string, vars map[string]string) error
	RenderString(tmplStr string, vars map[string]string) (string, error)
	AddFuncs(funcs template.FuncMap)
}

// New creates a new TemplateEngine with default custom functions.
func New() TemplateEngine {
	e := &textEngine{}
	e.registerDefaultFuncs()
	return e
}
