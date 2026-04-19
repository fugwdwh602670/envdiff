package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplate_Default(t *testing.T) {
	results := []Result{
		{Key: "FOO", Status: StatusMissingInB},
		{Key: "BAR", Status: StatusMismatch, ValueA: "x", ValueB: "y"},
	}
	var buf bytes.Buffer
	if err := RenderTemplate(results, DefaultTemplateOptions(), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "b=y") {
		t.Errorf("expected b=y in output, got: %s", out)
	}
}

func TestRenderTemplate_CustomStr(t *testing.T) {
	results := []Result{
		{Key: "HELLO", Status: StatusMatch},
	}
	opts := DefaultTemplateOptions()
	opts.TemplateStr = `{{range .}}KEY={{.Key}}{{end}}`
	var buf bytes.Buffer
	if err := RenderTemplate(results, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "KEY=HELLO" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderTemplate_FromFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tmpl.txt")
	_ = os.WriteFile(p, []byte(`{{range .}}{{.Key}};{{end}}`), 0644)

	results := []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMissingInA},
	}
	opts := DefaultTemplateOptions()
	opts.TemplatePath = p
	var buf bytes.Buffer
	if err := RenderTemplate(results, opts, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "A;") || !strings.Contains(buf.String(), "B;") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestRenderTemplate_BadTemplate(t *testing.T) {
	opts := DefaultTemplateOptions()
	opts.TemplateStr = `{{.Unclosed`
	var buf bytes.Buffer
	if err := RenderTemplate(nil, opts, &buf); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestRenderTemplate_MissingFile(t *testing.T) {
	opts := DefaultTemplateOptions()
	opts.TemplatePath = "/nonexistent/path/tmpl.txt"
	var buf bytes.Buffer
	if err := RenderTemplate(nil, opts, &buf); err == nil {
		t.Fatal("expected error for missing file")
	}
}
