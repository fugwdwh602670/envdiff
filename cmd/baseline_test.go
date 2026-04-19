package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runBaselineCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestBaselineSaveAndCheck_NoDiff(t *testing.T) {
	dir := t.TempDir()
	envA := filepath.Join(dir, ".env")
	envB := filepath.Join(dir, ".env.prod")
	baseline := filepath.Join(dir, "baseline.json")

	os.WriteFile(envA, []byte("KEY=value\n"), 0644)
	os.WriteFile(envB, []byte("KEY=value\n"), 0644)

	out, err := runBaselineCmd("baseline", "save", envA, envB, "--baseline", baseline)
	if err != nil {
		t.Fatalf("save error: %v\n%s", err, out)
	}
	if !strings.Contains(out, "Baseline saved") {
		t.Errorf("expected saved message, got: %s", out)
	}

	out, _ = runBaselineCmd("baseline", "check", envA, envB, "--baseline", baseline)
	if !strings.Contains(out, "New issues: 0") {
		t.Errorf("expected 0 new issues, got: %s", out)
	}
	if !strings.Contains(out, "Resolved: 0") {
		t.Errorf("expected 0 resolved, got: %s", out)
	}
}

func TestBaselineCheck_MissingBaseline(t *testing.T) {
	dir := t.TempDir()
	envA := filepath.Join(dir, ".env")
	envB := filepath.Join(dir, ".env.prod")
	os.WriteFile(envA, []byte("KEY=val\n"), 0644)
	os.WriteFile(envB, []byte("KEY=val\n"), 0644)

	_, err := runBaselineCmd("baseline", "check", envA, envB, "--baseline", "/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing baseline")
	}
}

func TestBaselineSave_MissingFile(t *testing.T) {
	dir := t.TempDir()
	baseline := filepath.Join(dir, "baseline.json")
	_, err := runBaselineCmd("baseline", "save", "/no/such/file", "/no/such/file2", "--baseline", baseline)
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

var _ = func() *cobra.Command { return rootCmd }()
