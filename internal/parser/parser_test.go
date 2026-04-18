package parser

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY1=value1\nKEY2=value2\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "value1" || env["KEY2"] != "value2" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParseFile_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=val\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 || env["KEY"] != "val" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParseFile_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `KEY1="hello world"` + "\n" + `KEY2='single'` + "\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY1"] != "hello world" {
		t.Errorf("KEY1: got %q", env["KEY1"])
	}
	if env["KEY2"] != "single" {
		t.Errorf("KEY2: got %q", env["KEY2"])
	}
}

func TestParseFile_ExportPrefix(t *testing.T) {
	path := writeTempEnv(t, "export MY_VAR=123\n")
	env, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["MY_VAR"] != "123" {
		t.Errorf("MY_VAR: got %q", env["MY_VAR"])
	}
}

func TestParseFile_MalformedLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := ParseFile(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
