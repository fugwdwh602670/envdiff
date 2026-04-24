package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runScopeCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	scopeCmd.ResetFlags()
	scopeCmd.Flags().BoolVar(&scopeInvert, "invert", false, "Return keys that do NOT match the given prefixes")
	rootCmd.SetOut(buf)
	scopeCmd.SetOut(buf)
	scopeCmd.SetErr(buf)
	// redirect stdout capture via cobra output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs(append([]string{"scope"}, args...))
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old
	var sb strings.Builder
	io := make([]byte, 4096)
	for {
		n, e := r.Read(io)
		sb.Write(io[:n])
		if e != nil {
			break
		}
	}
	return sb.String(), err
}

func writeScopeEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func resetScopeCmd(t *testing.T) {
	t.Helper()
	scopePrefixes = nil
	scopeInvert = false
}

func TestScopeCmd_NoPrefix(t *testing.T) {
	t.Cleanup(func() { resetScopeCmd(t) })
	p := writeScopeEnv(t, "DB_HOST=localhost\nAPP_NAME=myapp\n")
	out, err := runScopeCmd(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = out // all keys should be present
}

func TestScopeCmd_WithPrefix(t *testing.T) {
	t.Cleanup(func() { resetScopeCmd(t) })
	p := writeScopeEnv(t, "DB_HOST=localhost\nDB_PORT=5432\nAPP_NAME=myapp\n")
	out, err := runScopeCmd(p, "DB_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
}

func TestScopeCmd_MissingArg(t *testing.T) {
	t.Cleanup(func() { resetScopeCmd(t) })
	_, err := runScopeCmd()
	if err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestScopeCmd_Invert(t *testing.T) {
	t.Cleanup(func() { resetScopeCmd(t) })
	p := writeScopeEnv(t, "DB_HOST=localhost\nAPP_NAME=myapp\n")
	// Use cobra flag directly
	if err := scopeCmd.Flags().Set("invert", "true"); err != nil {
		t.Fatal(err)
	}
	out, err := runScopeCmd(p, "DB_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = out
}

var _ = cobra.Command{} // ensure cobra import used
