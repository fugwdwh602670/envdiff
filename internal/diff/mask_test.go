package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeMaskEnv() map[string]string {
	return map[string]string{
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123xyz",
		"APP_NAME":     "envdiff",
		"SMTP_TOKEN":   "tok_live_abc",
		"DEBUG":        "true",
	}
}

func findMaskResult(results []MaskResult, key string) (MaskResult, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return MaskResult{}, false
}

func TestMaskEnv_MasksMatchingKeys(t *testing.T) {
	env := makeMaskEnv()
	opts := DefaultMaskOptions()
	results := MaskEnv(env, opts)

	r, ok := findMaskResult(results, "DB_PASSWORD")
	if !ok {
		t.Fatal("expected DB_PASSWORD in results")
	}
	if !r.WasMasked {
		t.Error("expected DB_PASSWORD to be masked")
	}
	if r.Masked != strings.Repeat("*", len(r.Original)) {
		t.Errorf("unexpected masked value: %s", r.Masked)
	}
}

func TestMaskEnv_DoesNotMaskPlainKeys(t *testing.T) {
	env := makeMaskEnv()
	opts := DefaultMaskOptions()
	results := MaskEnv(env, opts)

	r, ok := findMaskResult(results, "APP_NAME")
	if !ok {
		t.Fatal("expected APP_NAME in results")
	}
	if r.WasMasked {
		t.Error("APP_NAME should not be masked")
	}
	if r.Masked != "envdiff" {
		t.Errorf("expected original value, got %s", r.Masked)
	}
}

func TestMaskEnv_VisibleChars(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "abcdefgh"}
	opts := DefaultMaskOptions()
	opts.VisibleChars = 3
	results := MaskEnv(env, opts)

	r, ok := findMaskResult(results, "SECRET_KEY")
	if !ok {
		t.Fatal("expected SECRET_KEY")
	}
	if r.Masked != "*****fgh" {
		t.Errorf("expected *****fgh, got %s", r.Masked)
	}
}

func TestMaskEnv_CustomPattern(t *testing.T) {
	env := map[string]string{"MY_CUSTOM_FIELD": "value123", "NORMAL": "ok"}
	opts := MaskOptions{
		Patterns:     []string{`(?i)custom`},
		MaskChar:     "#",
		VisibleChars: 0,
	}
	results := MaskEnv(env, opts)

	r, ok := findMaskResult(results, "MY_CUSTOM_FIELD")
	if !ok {
		t.Fatal("expected MY_CUSTOM_FIELD")
	}
	if r.Masked != "########" {
		t.Errorf("expected ########, got %s", r.Masked)
	}
}

func TestWriteMaskReport_Output(t *testing.T) {
	env := map[string]string{"API_KEY": "secret123", "HOST": "localhost"}
	opts := DefaultMaskOptions()
	results := MaskEnv(env, opts)

	var buf bytes.Buffer
	WriteMaskReport(&buf, results)
	out := buf.String()

	if !strings.Contains(out, "[MASKED]") {
		t.Error("expected [MASKED] in output")
	}
	if !strings.Contains(out, "[plain]") {
		t.Error("expected [plain] in output")
	}
	if strings.Contains(out, "secret123") {
		t.Error("original secret value should not appear in report")
	}
}
