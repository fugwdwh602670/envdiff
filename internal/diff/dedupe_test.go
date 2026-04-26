package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makePairs(kv ...string) [][2]string {
	if len(kv)%2 != 0 {
		panic("makePairs requires even number of arguments")
	}
	pairs := make([][2]string, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		pairs = append(pairs, [2]string{kv[i], kv[i+1]})
	}
	return pairs
}

func TestDedupeEnv_NoDuplicates(t *testing.T) {
	pairs := makePairs("FOO", "bar", "BAZ", "qux")
	res := DedupeEnv(pairs, DefaultDedupeOptions())
	if len(res.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %d", len(res.Duplicates))
	}
	if res.Env["FOO"] != "bar" || res.Env["BAZ"] != "qux" {
		t.Fatalf("unexpected env: %v", res.Env)
	}
}

func TestDedupeEnv_PreferFirst(t *testing.T) {
	pairs := makePairs("KEY", "first", "OTHER", "x", "KEY", "second")
	opts := DefaultDedupeOptions() // PreferLast = false
	res := DedupeEnv(pairs, opts)
	if len(res.Duplicates) != 1 {
		t.Fatalf("expected 1 duplicate, got %d", len(res.Duplicates))
	}
	if res.Env["KEY"] != "first" {
		t.Errorf("expected 'first', got %q", res.Env["KEY"])
	}
	if res.Duplicates[0].Kept != "first" {
		t.Errorf("Kept should be 'first', got %q", res.Duplicates[0].Kept)
	}
}

func TestDedupeEnv_PreferLast(t *testing.T) {
	pairs := makePairs("KEY", "first", "KEY", "second", "KEY", "third")
	opts := DedupeOptions{PreferLast: true}
	res := DedupeEnv(pairs, opts)
	if res.Env["KEY"] != "third" {
		t.Errorf("expected 'third', got %q", res.Env["KEY"])
	}
	if len(res.Duplicates[0].Values) != 3 {
		t.Errorf("expected 3 values recorded, got %d", len(res.Duplicates[0].Values))
	}
}

func TestDedupeEnv_ReportOnly(t *testing.T) {
	pairs := makePairs("A", "1", "A", "2", "B", "3")
	opts := DedupeOptions{ReportOnly: true}
	res := DedupeEnv(pairs, opts)
	if len(res.Env) != 0 {
		t.Errorf("expected empty env in report-only mode, got %v", res.Env)
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0].Key != "A" {
		t.Errorf("expected duplicate for A, got %v", res.Duplicates)
	}
}

func TestWriteDedupeReport_NoDuplicates(t *testing.T) {
	res := DedupeResult{}
	var buf bytes.Buffer
	WriteDedupeReport(&buf, res)
	if !strings.Contains(buf.String(), "No duplicate") {
		t.Errorf("expected no-duplicate message, got: %s", buf.String())
	}
}

func TestWriteDedupeReport_WithDuplicates(t *testing.T) {
	pairs := makePairs("SECRET", "alpha", "SECRET", "beta")
	res := DedupeEnv(pairs, DefaultDedupeOptions())
	var buf bytes.Buffer
	WriteDedupeReport(&buf, res)
	out := buf.String()
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected SECRET in report, got: %s", out)
	}
	if !strings.Contains(out, "2 occurrences") {
		t.Errorf("expected occurrence count in report, got: %s", out)
	}
}
