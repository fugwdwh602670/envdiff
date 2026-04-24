package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeAuditResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMatch, ValueA: "localhost", ValueB: "localhost"},
		{Key: "DB_PASS", Status: StatusMismatch, ValueA: "secret", ValueB: "newsecret"},
		{Key: "API_KEY", Status: StatusMissingInB, ValueA: "abc123", ValueB: ""},
		{Key: "NEW_VAR", Status: StatusMissingInA, ValueA: "", ValueB: "hello"},
	}
}

func TestAuditResults_ExcludesUnchanged(t *testing.T) {
	results := makeAuditResults()
	entries := AuditResults(results, DefaultAuditOptions())
	for _, e := range entries {
		if e.Event == "unchanged" {
			t.Errorf("expected unchanged entries to be excluded, got key %s", e.Key)
		}
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestAuditResults_IncludesUnchanged(t *testing.T) {
	results := makeAuditResults()
	opts := DefaultAuditOptions()
	opts.IncludeUnchanged = true
	entries := AuditResults(results, opts)
	if len(entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(entries))
	}
}

func TestAuditResults_EventMapping(t *testing.T) {
	results := makeAuditResults()
	opts := DefaultAuditOptions()
	opts.IncludeUnchanged = true
	entries := AuditResults(results, opts)

	events := map[string]string{}
	for _, e := range entries {
		events[e.Key] = e.Event
	}

	expected := map[string]string{
		"DB_HOST": "unchanged",
		"DB_PASS": "changed",
		"API_KEY": "removed",
		"NEW_VAR": "added",
	}
	for k, want := range expected {
		if got := events[k]; got != want {
			t.Errorf("key %s: expected event %q, got %q", k, want, got)
		}
	}
}

func TestAuditResults_RedactValues(t *testing.T) {
	results := makeAuditResults()
	opts := DefaultAuditOptions()
	opts.RedactValues = true
	entries := AuditResults(results, opts)
	for _, e := range entries {
		if e.OldValue != "" || e.NewValue != "" {
			t.Errorf("key %s: expected redacted values, got old=%q new=%q", e.Key, e.OldValue, e.NewValue)
		}
	}
}

func TestAuditResults_SortedByKey(t *testing.T) {
	results := makeAuditResults()
	entries := AuditResults(results, DefaultAuditOptions())
	for i := 1; i < len(entries); i++ {
		if entries[i].Key < entries[i-1].Key {
			t.Errorf("entries not sorted: %s before %s", entries[i-1].Key, entries[i].Key)
		}
	}
}

func TestWriteAuditReport_WithEntries(t *testing.T) {
	results := makeAuditResults()
	entries := AuditResults(results, DefaultAuditOptions())
	var buf bytes.Buffer
	WriteAuditReport(&buf, entries)
	out := buf.String()
	if !strings.Contains(out, "Audit Log") {
		t.Error("expected header in audit report")
	}
	if !strings.Contains(out, "changed") {
		t.Error("expected 'changed' event in report")
	}
}

func TestWriteAuditReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteAuditReport(&buf, []AuditEntry{})
	if !strings.Contains(buf.String(), "No audit events") {
		t.Error("expected empty message")
	}
}
