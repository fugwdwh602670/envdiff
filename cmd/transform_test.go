package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runTransformCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"transform"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func writeTransformEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func resetTransformCmd() {
	transformUpper = false
	transformLower = false
	transformTrimmed = false
	transformPrefix = ""
	transformOutput = ""
}

func TestTransformCmd_NoOp(t *testing.T) {
	defer resetTransformCmd()
	p := writeTransformEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	out, err := runTransformCmd(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "0 key(s) changed") {
		t.Errorf("expected 0 changes, got: %s", out)
	}
}

func TestTransformCmd_Upper(t *testing.T) {
	defer resetTransformCmd()
	p := writeTransformEnv(t, "APP_NAME=myapp\nDB_HOST=localhost\n")
	out, err := runTransformCmd("--upper", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2 key(s) changed") {
		t.Errorf("expected 2 changes, got: %s", out)
	}
}

func TestTransformCmd_PrefixTransform(t *testing.T) {
	defer resetTransformCmd()
	p := writeTransformEnv(t, "DB_HOST=localhost\nDB_PASS=secret\nAPP_NAME=myapp\n")
	out, err := runTransformCmd("--prefix-transform", "DB_=upper", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "2 key(s) changed") {
		t.Errorf("expected 2 DB_ keys changed, got: %s", out)
	}
}

func TestTransformCmd_MissingArg(t *testing.T) {
	defer resetTransformCmd()
	// Suppress cobra error output
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"transform"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestTransformCmd_InvalidPrefixFormat(t *testing.T) {
	defer resetTransformCmd()
	p := writeTransformEnv(t, "DB_HOST=localhost\n")
	_ = p
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetOut(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"transform", "--prefix-transform", "BADFORMAT", p})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for bad prefix-transform format")
	}
}

func TestTransformCmd_OutputFile(t *testing.T) {
	defer resetTransformCmd()
	p := writeTransformEnv(t, "APP_NAME=myapp\n")
	outFile := filepath.Join(t.TempDir(), "report.txt")
	_, err := runTransformCmd("--upper", "--output", outFile, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	if !strings.Contains(string(data), "key(s) changed") {
		t.Errorf("expected report in output file, got: %s", string(data))
	}
}

var _ = func() *cobra.Command { return rootCmd }()
