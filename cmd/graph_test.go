package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runGraphCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"graph"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func writeGraphEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func resetGraphCmd() {
	rootCmd.ResetFlags()
	rootCmd.ResetCommands()
	// Re-register by resetting cobra state for test isolation
}

func TestGraphCmd_BasicOutput(t *testing.T) {
	file := writeGraphEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=test\nSTANDALONE=yes\n")

	cmd := &cobra.Command{Use: "graph"}
	_ = cmd

	out, err := runGraphCmd(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB/") {
		t.Errorf("expected DB/ group in output, got:\n%s", out)
	}
	if !strings.Contains(out, "STANDALONE") {
		t.Errorf("expected STANDALONE in output, got:\n%s", out)
	}
}

func TestGraphCmd_WithValues(t *testing.T) {
	file := writeGraphEnv(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	out, err := runGraphCmd("--values", file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected value 'localhost' in output, got:\n%s", out)
	}
}

func TestGraphCmd_MissingArg(t *testing.T) {
	_, err := runGraphCmd()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}
