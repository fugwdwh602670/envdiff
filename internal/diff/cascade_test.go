package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeLayers() []map[string]string {
	return []map[string]string{
		{"BASE": "base_val", "SHARED": "from_base", "ONLY_BASE": "yes"},
		{"SHARED": "from_mid", "MID_ONLY": "mid"},
		{"SHARED": "from_top", "TOP_ONLY": "top"},
	}
}

func TestCascadeEnv_Overwrite(t *testing.T) {
	layers := makeLayers()
	opts := DefaultCascadeOptions()
	merged, entries := CascadeEnv(layers, opts)

	if merged["SHARED"] != "from_top" {
		t.Errorf("expected SHARED=from_top, got %q", merged["SHARED"])
	}
	if merged["BASE"] != "base_val" {
		t.Errorf("expected BASE=base_val, got %q", merged["BASE"])
	}
	if len(entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(entries))
	}
}

func TestCascadeEnv_NoOverwrite(t *testing.T) {
	layers := makeLayers()
	opts := DefaultCascadeOptions()
	opts.Overwrite = false
	merged, _ := CascadeEnv(layers, opts)

	if merged["SHARED"] != "from_base" {
		t.Errorf("expected SHARED=from_base (no overwrite), got %q", merged["SHARED"])
	}
}

func TestCascadeEnv_SkipEmpty(t *testing.T) {
	layers := []map[string]string{
		{"KEY": "base"},
		{"KEY": ""},
	}
	opts := DefaultCascadeOptions()
	opts.SkipEmpty = true
	merged, _ := CascadeEnv(layers, opts)

	if merged["KEY"] != "base" {
		t.Errorf("expected KEY=base (skip empty), got %q", merged["KEY"])
	}
}

func TestCascadeEnv_EmptyLayers(t *testing.T) {
	merged, entries := CascadeEnv(nil, DefaultCascadeOptions())
	if len(merged) != 0 {
		t.Errorf("expected empty merged map")
	}
	if len(entries) != 0 {
		t.Errorf("expected empty entries")
	}
}

func TestWriteCascadeReport_WithNames(t *testing.T) {
	layers := makeLayers()
	_, entries := CascadeEnv(layers, DefaultCascadeOptions())

	var buf bytes.Buffer
	WriteCascadeReport(&buf, entries, []string{"base.env", "mid.env", "prod.env"})
	out := buf.String()

	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected layer name prod.env in output")
	}
	if !strings.Contains(out, "SHARED") {
		t.Errorf("expected SHARED key in output")
	}
}

func TestWriteCascadeReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteCascadeReport(&buf, nil, nil)
	if !strings.Contains(buf.String(), "No keys resolved") {
		t.Errorf("expected empty message")
	}
}
