package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestStatEnv_Basic(t *testing.T) {
	env := map[string]string{
		"FOO": "bar",
		"BAZ": "",
		"QUX": "value",
	}
	stats := StatEnv(env)
	if stats.TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", stats.TotalKeys)
	}
	if stats.EmptyValues != 1 {
		t.Errorf("expected 1 empty value, got %d", stats.EmptyValues)
	}
	if stats.UniqueKeys != 3 {
		t.Errorf("expected 3 unique keys, got %d", stats.UniqueKeys)
	}
}

func TestStatEnv_Empty(t *testing.T) {
	stats := StatEnv(map[string]string{})
	if stats.TotalKeys != 0 {
		t.Errorf("expected 0 keys")
	}
	if stats.EmptyValues != 0 {
		t.Errorf("expected 0 empty values")
	}
}

func TestStatEnv_AllEmpty(t *testing.T) {
	env := map[string]string{
		"A": "",
		"B": "",
	}
	stats := StatEnv(env)
	if stats.EmptyValues != 2 {
		t.Errorf("expected 2 empty values, got %d", stats.EmptyValues)
	}
}

func TestWriteStatsReport(t *testing.T) {
	env := map[string]string{
		"KEY1": "val",
		"KEY2": "",
	}
	stats := StatEnv(env)
	var buf bytes.Buffer
	WriteStatsReport(&buf, "test.env", stats)
	out := buf.String()
	if !strings.Contains(out, "test.env") {
		t.Error("expected label in output")
	}
	if !strings.Contains(out, "Total keys") {
		t.Error("expected Total keys in output")
	}
	if !strings.Contains(out, "Empty values") {
		t.Error("expected Empty values in output")
	}
}
