package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestPromote_OnlyMissing(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	tgt := map[string]string{"A": "99", "D": "4"}

	out, results := Promote(src, tgt, DefaultPromoteOptions())

	if out["A"] != "99" {
		t.Errorf("expected A to remain 99, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %s", out["C"])
	}

	actions := actionsMap(results)
	if actions["B"] != "added" {
		t.Errorf("expected B added, got %s", actions["B"])
	}
	if actions["A"] != "skipped" {
		t.Errorf("expected A skipped, got %s", actions["A"])
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := map[string]string{"A": "new", "B": "2"}
	tgt := map[string]string{"A": "old"}

	opts := DefaultPromoteOptions()
	opts.OverwriteExisting = true
	opts.OnlyMissing = false

	out, results := Promote(src, tgt, opts)

	if out["A"] != "new" {
		t.Errorf("expected A=new, got %s", out["A"])
	}
	actions := actionsMap(results)
	if actions["A"] != "overwritten" {
		t.Errorf("expected A overwritten, got %s", actions["A"])
	}
}

func TestPromote_DryRun(t *testing.T) {
	src := map[string]string{"X": "10"}
	tgt := map[string]string{}

	opts := DefaultPromoteOptions()
	opts.DryRun = true

	out, results := Promote(src, tgt, opts)

	if _, ok := out["X"]; ok {
		t.Error("dry run should not modify target")
	}
	if len(results) != 1 || results[0].Action != "added" {
		t.Error("expected one 'added' result in dry run")
	}
}

func TestPromote_IdenticalValues(t *testing.T) {
	src := map[string]string{"K": "same"}
	tgt := map[string]string{"K": "same"}

	_, results := Promote(src, tgt, DefaultPromoteOptions())
	if len(results) != 0 {
		t.Errorf("expected no results for identical values, got %d", len(results))
	}
}

func TestWritePromoteReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WritePromoteReport(&buf, nil)
	if !strings.Contains(buf.String(), "No keys promoted") {
		t.Error("expected 'No keys promoted' message")
	}
}

func TestWritePromoteReport_Actions(t *testing.T) {
	results := []PromoteResult{
		{Key: "A", OldValue: "", NewValue: "1", Action: "added"},
		{Key: "B", OldValue: "old", NewValue: "new", Action: "overwritten"},
		{Key: "C", OldValue: "x", NewValue: "y", Action: "skipped"},
	}
	var buf bytes.Buffer
	WritePromoteReport(&buf, results)
	out := buf.String()
	if !strings.Contains(out, "[added]") {
		t.Error("expected [added] in report")
	}
	if !strings.Contains(out, "[overwritten]") {
		t.Error("expected [overwritten] in report")
	}
	if !strings.Contains(out, "[skipped]") {
		t.Error("expected [skipped] in report")
	}
}

func actionsMap(results []PromoteResult) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Action
	}
	return m
}
