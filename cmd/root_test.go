package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return p
}

func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestRootCmd_NoDiff(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	b := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	out, err := runCmd(t, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestRootCmd_WithDiff(t *testing.T) {
	a := writeTempEnv(t, "FOO=bar\nONLY_A=1\n")
	b := writeTempEnv(t, "FOO=changed\nONLY_B=2\n")
	out, err := runCmd(t, a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "ONLY_A") {
		t.Errorf("expected ONLY_A in output, got: %s", out)
	}
	if !strings.Contains(out, "ONLY_B") {
		t.Errorf("expected ONLY_B in output, got: %s", out)
	}
}

func TestRootCmd_MissingArgs(t *testing.T) {
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing args")
	}
}
