package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runValidateCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return runCmd(t, append([]string{"validate"}, args...)...)
}

func writeValidateEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestValidateCmd_NoIssues(t *testing.T) {
	p := writeValidateEnv(t, "HOST=localhost\nPORT=8080\n")
	out, err := runValidateCmd(t, p)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "no validation issues") {
		t.Errorf("expected no-issue message, got: %s", out)
	}
}

func TestValidateCmd_EmptyValue(t *testing.T) {
	p := writeValidateEnv(t, "API_KEY=\n")
	out, _ := runValidateCmd(t, p)
	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty value warning, got: %s", out)
	}
}

func TestValidateCmd_LowercaseKey(t *testing.T) {
	p := writeValidateEnv(t, "secret=abc\n")
	out, _ := runValidateCmd(t, p)
	if !strings.Contains(out, "uppercase") {
		t.Errorf("expected uppercase warning, got: %s", out)
	}
}

func TestValidateCmd_MissingArg(t *testing.T) {
	_, err := runValidateCmd(t)
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestValidateCmd_DisableUppercase(t *testing.T) {
	p := writeValidateEnv(t, "lowercase_key=value\n")
	out, err := runValidateCmd(t, "--require-uppercase=false", "--forbid-empty=false", p)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "no validation issues") {
		t.Errorf("expected no issues with uppercase disabled, got: %s", out)
	}
}
