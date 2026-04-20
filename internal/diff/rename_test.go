package diff

import (
	"bytes"
	"testing"
)

func TestApplyRenames_Basic(t *testing.T) {
	env := map[string]string{
		"OLD_KEY": "value1",
		"KEEP_KEY": "value2",
	}
	opts := DefaultRenameOptions()
	opts.Renames = RenameMap{"OLD_KEY": "NEW_KEY"}

	result := ApplyRenames(env, opts)

	if _, ok := result["OLD_KEY"]; ok {
		t.Error("expected OLD_KEY to be removed")
	}
	if v, ok := result["NEW_KEY"]; !ok || v != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", v)
	}
	if v, ok := result["KEEP_KEY"]; !ok || v != "value2" {
		t.Errorf("expected KEEP_KEY=value2, got %q", v)
	}
}

func TestApplyRenames_IgnoreUnmatched_False(t *testing.T) {
	env := map[string]string{
		"OLD_KEY":  "value1",
		"KEEP_KEY": "value2",
	}
	opts := RenameOptions{
		Renames:         RenameMap{"OLD_KEY": "NEW_KEY"},
		IgnoreUnmatched: false,
	}

	result := ApplyRenames(env, opts)

	if _, ok := result["KEEP_KEY"]; ok {
		t.Error("expected KEEP_KEY to be dropped when IgnoreUnmatched=false")
	}
	if v, ok := result["NEW_KEY"]; !ok || v != "value1" {
		t.Errorf("expected NEW_KEY=value1, got %q", v)
	}
}

func TestApplyRenames_EmptyMap(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	opts := DefaultRenameOptions()

	result := ApplyRenames(env, opts)

	if len(result) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestWriteRenameReport_Applied(t *testing.T) {
	original := map[string]string{"OLD_KEY": "secret"}
	renamed := map[string]string{"NEW_KEY": "secret"}
	renames := RenameMap{"OLD_KEY": "NEW_KEY"}

	var buf bytes.Buffer
	WriteRenameReport(&buf, original, renamed, renames)

	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty report")
	}
	if !bytes.Contains(buf.Bytes(), []byte("OLD_KEY -> NEW_KEY")) {
		t.Errorf("expected rename entry in report, got: %s", out)
	}
}

func TestWriteRenameReport_NoRenames(t *testing.T) {
	original := map[string]string{"FOO": "bar"}
	renamed := map[string]string{"FOO": "bar"}
	renames := RenameMap{"MISSING_KEY": "OTHER_KEY"}

	var buf bytes.Buffer
	WriteRenameReport(&buf, original, renamed, renames)

	if !bytes.Contains(buf.Bytes(), []byte("No renames applied.")) {
		t.Errorf("expected no-renames message, got: %s", buf.String())
	}
}
