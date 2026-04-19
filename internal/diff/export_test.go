package diff

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeExportResults() []Result {
	return []Result{
		{Key: "APP_NAME", Status: StatusMatch, ValueA: "myapp", ValueB: "myapp"},
		{Key: "DB_PASS", Status: StatusMissingInB, ValueA: "secret", ValueB: ""},
		{Key: "API_KEY", Status: StatusMismatch, ValueA: "abc", ValueB: "xyz"},
	}
}

func TestExportResults_EnvFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.env")
	err := ExportResults(makeExportResults(), out, DefaultExportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME in env output")
	}
}

func TestExportResults_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.json")
	opts := DefaultExportOptions()
	opts.Format = ExportFormatJSON
	err := ExportResults(makeExportResults(), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var results []Result
	if err := json.Unmarshal(data, &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestExportResults_MarkdownFormat(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.md")
	opts := DefaultExportOptions()
	opts.Format = ExportFormatMarkdown
	err := ExportResults(makeExportResults(), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if !strings.Contains(string(data), "| Key |") {
		t.Errorf("expected markdown table header")
	}
	if !strings.Contains(string(data), "DB_PASS") {
		t.Errorf("expected DB_PASS in markdown output")
	}
}

func TestExportResults_OnlyMissing(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "out.env")
	opts := DefaultExportOptions()
	opts.OnlyMissing = true
	err := ExportResults(makeExportResults(), out, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "APP_NAME") {
		t.Errorf("APP_NAME should be excluded when OnlyMissing=true")
	}
	if !strings.Contains(string(data), "DB_PASS") {
		t.Errorf("DB_PASS should be included when OnlyMissing=true")
	}
}
