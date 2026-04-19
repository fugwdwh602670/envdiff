package diff

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

// TemplateOptions controls template rendering.
type TemplateOptions struct {
	TemplatePath string
	TemplateStr  string
}

// DefaultTemplateOptions returns sensible defaults.
func DefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{}
}

const defaultTemplate = `{{range .}}[{{.Status}}] {{.Key}}{{if .ValueA}} (a={{.ValueA}}){{end}}{{if .ValueB}} (b={{.ValueB}}){{end}}
{{end}}`

// RenderTemplate renders diff results using a Go text/template.
// If opts.TemplateStr is set it is used directly; if opts.TemplatePath is set
// the file is read; otherwise a built-in template is used.
func RenderTemplate(results []Result, opts TemplateOptions, w io.Writer) error {
	tmplStr := defaultTemplate
	if opts.TemplateStr != "" {
		tmplStr = opts.TemplateStr
	} else if opts.TemplatePath != "" {
		data, err := readFile(opts.TemplatePath)
		if err != nil {
			return fmt.Errorf("template: read %s: %w", opts.TemplatePath, err)
		}
		tmplStr = string(data)
	}

	tmpl, err := template.New("envdiff").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template: parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, results); err != nil {
		return fmt.Errorf("template: execute: %w", err)
	}
	_, err = w.Write(buf.Bytes())
	return err
}
