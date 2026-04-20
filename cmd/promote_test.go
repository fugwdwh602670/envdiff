package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runPromoteCmd(args ...string) (string, error) {
	return runCmd(append([]string{"promote"}, args...)...)
}

func TestPromoteCmd_MissingArgs(t *testing.T) {
	_, err := runPromoteCmd()
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestPromoteCmd_DryRun_AddsKey(t *testing.T) {
	src := writeTempEnv(t, "A=1\nB=2\n")
	tgt := writeTempEnv(t, "A=99\n")

	out, err := runPromoteCmd("--dry-run", src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[added]") {
		t.Errorf("expected [added] in dry-run output, got: %s", out)
	}
	if strings.Contains(out, "Written to") {
		t.Error("dry-run should not write output")
	}

	// target file should be unchanged
	data, _ := os.ReadFile(tgt)
	if strings.Contains(string(data), "B=") {
		t.Error("target should not be modified during dry run")
	}
}

func TestPromoteCmd_WritesOutput(t *testing.T) {
	src := writeTempEnv(t, "A=1\nB=2\n")
	tgt := writeTempEnv(t, "A=99\n")

	outFile := filepath.Join(t.TempDir(), "result.env")
	_, err := runPromoteCmd("--output", outFile, src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	contents := string(data)
	if !strings.Contains(contents, "B=2") {
		t.Errorf("expected B=2 in output, got: %s", contents)
	}
	if !strings.Contains(contents, "A=99") {
		t.Errorf("expected original A=99 preserved, got: %s", contents)
	}
}

func TestPromoteCmd_Overwrite(t *testing.T) {
	src := writeTempEnv(t, "A=new\n")
	tgt := writeTempEnv(t, "A=old\n")

	outFile := filepath.Join(t.TempDir(), "result.env")
	out, err := runPromoteCmd("--overwrite", "--output", outFile, src, tgt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "[overwritten]") {
		t.Errorf("expected [overwritten] in output, got: %s", out)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "A=new") {
		t.Errorf("expected A=new in output file, got: %s", string(data))
	}
}
