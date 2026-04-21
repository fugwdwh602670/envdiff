package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeGroupResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMatch, ValueA: "localhost", ValueB: "localhost"},
		{Key: "DB_PORT", Status: StatusMismatch, ValueA: "5432", ValueB: "5433"},
		{Key: "APP_ENV", Status: StatusMissingInB, ValueA: "production", ValueB: ""},
		{Key: "APP_DEBUG", Status: StatusMatch, ValueA: "false", ValueB: "false"},
		{Key: "CACHE_TTL", Status: StatusMissingInA, ValueA: "", ValueB: "300"},
		{Key: "NOPREFIXKEY", Status: StatusMatch, ValueA: "x", ValueB: "x"},
	}
}

func TestGroupResults_ByPrefix(t *testing.T) {
	opts := DefaultGroupOptions()
	groups := GroupResults(makeGroupResults(), opts)

	if _, ok := groups["DB"]; !ok {
		t.Error("expected group 'DB'")
	}
	if _, ok := groups["APP"]; !ok {
		t.Error("expected group 'APP'")
	}
	if _, ok := groups["CACHE"]; !ok {
		t.Error("expected group 'CACHE'")
	}
	if _, ok := groups["(ungrouped)"]; !ok {
		t.Error("expected group '(ungrouped)' for keys without separator")
	}
	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(groups["DB"]))
	}
}

func TestGroupResults_ByStatus(t *testing.T) {
	opts := GroupOptions{GroupBy: "status", Separator: "_"}
	groups := GroupResults(makeGroupResults(), opts)

	if len(groups[string(StatusMatch)]) != 3 {
		t.Errorf("expected 3 match results, got %d", len(groups[string(StatusMatch)]))
	}
	if len(groups[string(StatusMismatch)]) != 1 {
		t.Errorf("expected 1 mismatch result, got %d", len(groups[string(StatusMismatch)]))
	}
	if len(groups[string(StatusMissingInA)]) != 1 {
		t.Errorf("expected 1 missing-in-A result, got %d", len(groups[string(StatusMissingInA)]))
	}
}

func TestGroupResults_Empty(t *testing.T) {
	groups := GroupResults([]Result{}, DefaultGroupOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestWriteGroupReport_Output(t *testing.T) {
	opts := DefaultGroupOptions()
	groups := GroupResults(makeGroupResults(), opts)

	var buf bytes.Buffer
	WriteGroupReport(&buf, groups)
	out := buf.String()

	if !strings.Contains(out, "[APP]") {
		t.Error("expected '[APP]' header in output")
	}
	if !strings.Contains(out, "[DB]") {
		t.Error("expected '[DB]' header in output")
	}
	if !strings.Contains(out, "missing in B") {
		t.Error("expected 'missing in B' annotation")
	}
	if !strings.Contains(out, "!=") {
		t.Error("expected mismatch annotation with '!='")
	}
}
