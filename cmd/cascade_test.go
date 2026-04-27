package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runCascadeCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"cascade"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func writeCascadeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func resetCascadeCmd(t *testing.T) {
	t.Helper()
	cascadeOverwrite = true
	cascadeSkipEmpty = false
	cascadeOutput = ""
	rootCmd.ResetFlags()
	for _, c := range rootCmd.Commands() {
		if c.Use == "cascade <base.env> [layer.env...]" {
			c.ResetFlags()
		}
	}
	_ = cobra.EnableCommandSorting
}

func TestCascadeCmd_BasicOutput(t *testing.T) {
	dir := t.TempDir()
	base := writeCascadeEnv(t, dir, "base.env", "KEY=base\nSHARED=from_base\n")
	top := writeCascadeEnv(t, dir, "top.env", "SHARED=from_top\nNEW=added\n")

	out, err := runCascadeCmd(base, top)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "SHARED") {
		t.Errorf("expected SHARED in output, got: %s", out)
	}
	if !strings.Contains(out, "from_top") {
		t.Errorf("expected from_top value in output, got: %s", out)
	}
}

func TestCascadeCmd_MissingArgs(t *testing.T) {
	_, err := runCascadeCmd("only_one.env")
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestCascadeCmd_InvalidFile(t *testing.T) {
	_, err := runCascadeCmd("/nonexistent/a.env", "/nonexistent/b.env")
	if err == nil {
		t.Error("expected error for invalid file")
	}
}
