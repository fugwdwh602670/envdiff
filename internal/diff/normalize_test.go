package diff

import (
	"strings"
	"testing"
)

func makeNormalizeEnv() map[string]string {
	return map[string]string{
		"API_KEY":   "  secret123  ",
		"DB_PASS":   "'mypassword'",
		"APP_NAME":  `"My App"`,
		"LOG_LEVEL": "INFO",
	}
}

func TestNormalizeEnv_TrimSpace(t *testing.T) {
	env := makeNormalizeEnv()
	opts := DefaultNormalizeOptions()
	opts.StripQuotes = false
	out := NormalizeEnv(env, opts)
	if out["API_KEY"] != "secret123" {
		t.Errorf("expected trimmed value, got %q", out["API_KEY"])
	}
	if out["LOG_LEVEL"] != "INFO" {
		t.Errorf("expected unchanged value, got %q", out["LOG_LEVEL"])
	}
}

func TestNormalizeEnv_StripQuotes(t *testing.T) {
	env := makeNormalizeEnv()
	opts := DefaultNormalizeOptions()
	opts.TrimSpace = false
	out := NormalizeEnv(env, opts)
	if out["DB_PASS"] != "mypassword" {
		t.Errorf("expected unquoted value, got %q", out["DB_PASS"])
	}
	if out["APP_NAME"] != "My App" {
		t.Errorf("expected unquoted value, got %q", out["APP_NAME"])
	}
}

func TestNormalizeEnv_LowerKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultNormalizeOptions()
	opts.LowerKeys = true
	out := NormalizeEnv(env, opts)
	if _, ok := out["foo"]; !ok {
		t.Error("expected lowercased key 'foo'")
	}
	if _, ok := out["baz"]; !ok {
		t.Error("expected lowercased key 'baz'")
	}
}

func TestNormalizeEnv_LowerValues(t *testing.T) {
	env := map[string]string{"KEY": "VALUE"}
	opts := DefaultNormalizeOptions()
	opts.LowerValues = true
	out := NormalizeEnv(env, opts)
	if out["KEY"] != "value" {
		t.Errorf("expected lowercase value, got %q", out["KEY"])
	}
}

func TestNormalizeEnv_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"X": "  hello  "}
	opts := DefaultNormalizeOptions()
	NormalizeEnv(env, opts)
	if env["X"] != "  hello  " {
		t.Error("original map was mutated")
	}
}

func TestWriteNormalizeReport_WithChanges(t *testing.T) {
	orig := map[string]string{"KEY": "  val  "}
	norm := map[string]string{"KEY": "val"}
	var sb strings.Builder
	WriteNormalizeReport(orig, norm, &sb)
	if !strings.Contains(sb.String(), "KEY") {
		t.Error("expected KEY in report output")
	}
}

func TestWriteNormalizeReport_NoChanges(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	var sb strings.Builder
	WriteNormalizeReport(env, env, &sb)
	if !strings.Contains(sb.String(), "No normalization changes") {
		t.Error("expected no-changes message")
	}
}
