package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeTransformEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
	}
}

func TestTransformEnv_ApplyToAll(t *testing.T) {
	env := makeTransformEnv()
	opts := DefaultTransformOptions()
	opts.ApplyToAll = strings.ToUpper

	out, results := TransformEnv(env, opts)

	for k, v := range out {
		if v != strings.ToUpper(env[k]) {
			t.Errorf("key %s: expected %q, got %q", k, strings.ToUpper(env[k]), v)
		}
	}
	if len(results) != len(env) {
		t.Errorf("expected %d results, got %d", len(env), len(results))
	}
}

func TestTransformEnv_PatternTransform(t *testing.T) {
	env := makeTransformEnv()
	opts := DefaultTransformOptions()
	opts.Transforms["DB_"] = func(v string) string { return "***" }

	out, results := TransformEnv(env, opts)

	if out["DB_HOST"] != "***" {
		t.Errorf("expected DB_HOST to be redacted, got %q", out["DB_HOST"])
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}

	changedCount := 0
	for _, r := range results {
		if r.Changed {
			changedCount++
		}
	}
	if changedCount != 2 {
		t.Errorf("expected 2 changed results, got %d", changedCount)
	}
}

func TestTransformEnv_NoChanges(t *testing.T) {
	env := makeTransformEnv()
	opts := DefaultTransformOptions()

	out, results := TransformEnv(env, opts)

	for k, v := range env {
		if out[k] != v {
			t.Errorf("key %s changed unexpectedly", k)
		}
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no changes, but %s changed", r.Key)
		}
	}
}

func TestTransformEnv_DoesNotMutateOriginal(t *testing.T) {
	env := makeTransformEnv()
	opts := DefaultTransformOptions()
	opts.ApplyToAll = strings.ToUpper

	origCopy := make(map[string]string, len(env))
	for k, v := range env {
		origCopy[k] = v
	}

	TransformEnv(env, opts)

	for k, v := range origCopy {
		if env[k] != v {
			t.Errorf("original mutated at key %s", k)
		}
	}
}

func TestWriteTransformReport_WithChanges(t *testing.T) {
	results := []TransformResult{
		{Key: "DB_HOST", Original: "localhost", Result: "***", Changed: true},
		{Key: "APP_NAME", Original: "myapp", Result: "myapp", Changed: false},
	}
	var buf bytes.Buffer
	WriteTransformReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "1 key(s) changed") {
		t.Errorf("expected change count in report, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in report")
	}
}
