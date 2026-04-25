package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeInterpolateEnv() map[string]string {
	return map[string]string{
		"HOST":     "localhost",
		"PORT":     "5432",
		"DB_URL":   "postgres://${HOST}:${PORT}/mydb",
		"APP_NAME": "envdiff",
		"GREETING": "Hello, ${APP_NAME}!",
	}
}

func TestInterpolateEnv_ResolvesRefs(t *testing.T) {
	env := makeInterpolateEnv()
	results, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		switch r.Key {
		case "DB_URL":
			if r.Resolved != "postgres://localhost:5432/mydb" {
				t.Errorf("DB_URL resolved = %q, want postgres://localhost:5432/mydb", r.Resolved)
			}
		case "GREETING":
			if r.Resolved != "Hello, envdiff!" {
				t.Errorf("GREETING resolved = %q, want 'Hello, envdiff!'", r.Resolved)
			}
		}
	}
}

func TestInterpolateEnv_MissingRef_Error(t *testing.T) {
	env := map[string]string{
		"URL": "http://${UNDEFINED_HOST}/path",
	}
	_, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err == nil {
		t.Fatal("expected error for missing ref, got nil")
	}
	if !strings.Contains(err.Error(), "UNDEFINED_HOST") {
		t.Errorf("error should mention missing key, got: %v", err)
	}
}

func TestInterpolateEnv_AllowMissing(t *testing.T) {
	env := map[string]string{
		"URL": "http://${MISSING}/path",
	}
	opts := DefaultInterpolateOptions()
	opts.AllowMissing = true
	results, err := InterpolateEnv(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if len(results[0].Missing) == 0 || results[0].Missing[0] != "MISSING" {
		t.Errorf("expected missing ref MISSING, got %v", results[0].Missing)
	}
}

func TestInterpolateEnv_NoRefs(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	results, err := InterpolateEnv(env, DefaultInterpolateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for env with no refs, got %d", len(results))
	}
}

func TestWriteInterpolateReport_WithResults(t *testing.T) {
	results := []InterpolateResult{
		{Key: "DB_URL", Original: "postgres://${HOST}/db", Resolved: "postgres://localhost/db", Refs: []string{"HOST"}},
	}
	var buf bytes.Buffer
	WriteInterpolateReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "postgres://localhost/db") {
		t.Errorf("expected resolved value in output, got: %s", out)
	}
}

func TestWriteInterpolateReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteInterpolateReport(&buf, nil)
	if !strings.Contains(buf.String(), "No interpolated") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
