// Package templateprompt provides functionalities for defining and rendering prompts.
package templateprompt

import (
	"bytes"
	"text/template"
)

// TemplateRenderer is a struct that holds a parsed template.
type TemplateRenderer struct {
	tmpl *template.Template
}

// NewTemplateRenderer is a constructor function that takes a template string,
// parses it, and returns a new TemplateRenderer with the parsed template.
func NewTemplateRenderer(tmpl string) (*TemplateRenderer, error) {
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{tmpl: t}, nil
}

// Render takes a map of values and uses the template to generate the final output.
func (tr *TemplateRenderer) Render(params map[string]string) (string, error) {
	var tmplBytes bytes.Buffer
	if err := tr.tmpl.Execute(&tmplBytes, params); err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
