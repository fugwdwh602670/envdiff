package cmd

import (
	"os"
	"strings"
	"testing"
)

func runLintCmd(t *testing.T, args ...string) (string, int) {
	t.Helper()
	out, err := runCmd(append([]string{"lint"}, args...)...)
	code := 0
	if err != nil {
		code = 1
	}
	return out, code
}

func TestLintCmd_NoIssues(t *testing.T) {
	f := writeTempEnv(t, "APP_NAME=myapp\nPORT=8080\n")
	out, code := runLintCmd(t, f)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d. output: %s", code, out)
	}
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestLintCmd_EmptyValue(t *testing.T) {
	f := writeTempEnv(t, "DB_PASSWORD=\nAPP=ok\n")
	out, _ := runLintCmd(t, f)
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in output, got: %s", out)
	}
	if !strings.Contains(out, "empty") {
		t.Errorf("expected 'empty' in output, got: %s", out)
	}
}

func TestLintCmd_BadNaming(t *testing.T) {
	f := writeTempEnv(t, "myKey=value\n")
	out, _ := runLintCmd(t, f)
	if !strings.Contains(out, "myKey") {
		t.Errorf("expected myKey in output, got: %s", out)
	}
}

func TestLintCmd_MissingArg(t *testing.T) {
	_, err := runCmd("lint")
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestLintCmd_FileNotFound(t *testing.T) {
	_, err := runCmd("lint", "/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func writeTempEnvLint(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}
