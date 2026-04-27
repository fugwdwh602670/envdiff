package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runMaskCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	maskCmd.SetOut(&buf)
	maskCmd.SetErr(&buf)
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs(append([]string{"mask"}, args...))
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func resetMaskCmd() {
	maskPatterns = nil
	maskChar = ""
	maskVisibleChars = 0
	maskShowAll = false
	maskCmd.ResetFlags()
	maskCmd.Flags().StringArrayVar(&maskPatterns, "pattern", nil, "")
	maskCmd.Flags().StringVar(&maskChar, "char", "", "")
	maskCmd.Flags().IntVar(&maskVisibleChars, "visible", 0, "")
	maskCmd.Flags().BoolVar(&maskShowAll, "report", false, "")
}

func writeMaskEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestMaskCmd_MasksSecretKeys(t *testing.T) {
	defer resetMaskCmd()
	f := writeMaskEnv(t, "API_KEY=abc123\nAPP_NAME=envdiff\n")
	out, err := runMaskCmd(t, f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "abc123") {
		t.Error("secret value should be masked in output")
	}
	if !strings.Contains(out, "APP_NAME=envdiff") {
		t.Error("non-secret key should appear unmasked")
	}
}

func TestMaskCmd_ReportFlag(t *testing.T) {
	defer resetMaskCmd()
	f := writeMaskEnv(t, "DB_PASSWORD=hunter2\nHOST=localhost\n")
	out, err := runMaskCmd(t, "--report", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[MASKED]") {
		t.Error("expected [MASKED] in report output")
	}
	if strings.Contains(out, "hunter2") {
		t.Error("original password should not appear in report")
	}
}

func TestMaskCmd_MissingArg(t *testing.T) {
	defer resetMaskCmd()
	// Temporarily suppress cobra error output
	maskCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"mask"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestMaskCmd_CustomPattern(t *testing.T) {
	defer resetMaskCmd()
	f := writeMaskEnv(t, "MY_CUSTOM=value\nNORMAL=ok\n")
	out, err := runMaskCmd(t, "--pattern", "(?i)custom", f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "value") {
		t.Error("custom-matched key value should be masked")
	}
	if !strings.Contains(out, "NORMAL=ok") {
		t.Error("non-matching key should be unmasked")
	}
	_ = cobra.Command{}
}
