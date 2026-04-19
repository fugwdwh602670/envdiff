package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runAnnotateCmd(args ...string) (string, error) {
	return runCmd(append([]string{"annotate"}, args...)...)
}

func TestAnnotateCmd_NoAnnotations(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.env")
	b := filepath.Join(dir, "b.env")
	os.WriteFile(a, []byte("KEY=val\n"), 0644)
	os.WriteFile(b, []byte("KEY=val\n"), 0644)

	out, err := runAnnotateCmd(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No annotations") {
		t.Errorf("expected 'No annotations', got: %s", out)
	}
}

func TestAnnotateCmd_WithAnnotation(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.env")
	b := filepath.Join(dir, "b.env")
	os.WriteFile(a, []byte("DB_HOST=localhost\n"), 0644)
	os.WriteFile(b, []byte("DB_HOST=prod.db\n"), 0644)

	out, err := runAnnotateCmd(a, b, "--annotate", "DB_HOST=Database hostname")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "Database hostname") {
		t.Errorf("expected comment in output, got: %s", out)
	}
}

func TestAnnotateCmd_MissingArgs(t *testing.T) {
	_, err := runAnnotateCmd()
	if err == nil {
		t.Error("expected error for missing args")
	}
}
