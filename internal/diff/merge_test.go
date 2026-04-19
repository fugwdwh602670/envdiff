package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMerge_PreferB(t *testing.T) {
	a := map[string]string{"FOO": "a", "ONLY_A": "1"}
	b := map[string]string{"FOO": "b", "ONLY_B": "2"}

	out := Merge(a, b, DefaultMergeOptions())

	if out["FOO"] != "b" {
		t.Errorf("expected FOO=b, got %s", out["FOO"])
	}
	if out["ONLY_A"] != "1" {
		t.Errorf("expected ONLY_A=1, got %s", out["ONLY_A"])
	}
	if out["ONLY_B"] != "2" {
		t.Errorf("expected ONLY_B=2, got %s", out["ONLY_B"])
	}
}

func TestMerge_PreferA(t *testing.T) {
	a := map[string]string{"FOO": "a"}
	b := map[string]string{"FOO": "b"}

	opts := DefaultMergeOptions()
	opts.PreferA = true
	out := Merge(a, b, opts)

	if out["FOO"] != "a" {
		t.Errorf("expected FOO=a, got %s", out["FOO"])
	}
}

func TestMerge_SkipMissing(t *testing.T) {
	a := map[string]string{"SHARED": "a", "ONLY_A": "1"}
	b := map[string]string{"SHARED": "b", "ONLY_B": "2"}

	opts := DefaultMergeOptions()
	opts.SkipMissing = true
	out := Merge(a, b, opts)

	if _, ok := out["ONLY_A"]; ok {
		t.Error("expected ONLY_A to be excluded")
	}
	if _, ok := out["ONLY_B"]; ok {
		t.Error("expected ONLY_B to be excluded")
	}
	if out["SHARED"] == "" {
		t.Error("expected SHARED to be present")
	}
}

func TestWriteMerged(t *testing.T) {
	merged := map[string]string{"BAR": "2", "FOO": "1"}
	dir := t.TempDir()
	path := filepath.Join(dir, "merged.env")

	if err := WriteMerged(merged, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}

	content := string(data)
	if content != "BAR=2\nFOO=1\n" {
		t.Errorf("unexpected content:\n%s", content)
	}
}
