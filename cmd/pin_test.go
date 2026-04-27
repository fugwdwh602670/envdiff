package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runPinCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func resetPinCmd(pinFile string) {
	pinFilePath = pinFile
	pinOverwrite = false
	pinReportOnly = false
	pinComment = ""
}

func writePinEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPinCmd_BasicPin(t *testing.T) {
	envFile := writePinEnv(t, "APP_SECRET=abc123\nDB_PASS=hunter2\n")
	pinFile := filepath.Join(t.TempDir(), "pins.json")
	resetPinCmd(pinFile)

	out, err := runPinCmd(t, "pin", "--pin-file", pinFile, envFile, "APP_SECRET")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "APP_SECRET") {
		t.Errorf("expected APP_SECRET in output, got: %s", out)
	}
}

func TestPinCmd_CheckNoViolations(t *testing.T) {
	envFile := writePinEnv(t, "APP_SECRET=abc123\n")
	pinFile := filepath.Join(t.TempDir(), "pins.json")
	resetPinCmd(pinFile)

	_, _ = runPinCmd(t, "pin", "--pin-file", pinFile, envFile, "APP_SECRET")

	out, err := runPinCmd(t, "pin", "check", "--pin-file", pinFile, envFile)
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(out, "No pin violations") {
		t.Errorf("expected no violations message, got: %s", out)
	}
}

func TestPinCmd_CheckViolation(t *testing.T) {
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env")
	pinFile := filepath.Join(dir, "pins.json")

	_ = os.WriteFile(envFile, []byte("APP_SECRET=original\n"), 0o644)
	resetPinCmd(pinFile)
	_, _ = runPinCmd(t, "pin", "--pin-file", pinFile, envFile, "APP_SECRET")

	_ = os.WriteFile(envFile, []byte("APP_SECRET=tampered\n"), 0o644)
	out, err := runPinCmd(t, "pin", "check", "--pin-file", pinFile, envFile)
	if err == nil {
		t.Errorf("expected error for violation, got none. output: %s", out)
	}
}

func TestPinCmd_MissingArgs(t *testing.T) {
	_ = cobra.Command{}
	_, err := runPinCmd(t, "pin")
	if err == nil {
		t.Error("expected error for missing args")
	}
}
