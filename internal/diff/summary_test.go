package diff

import (
	"testing"
)

func makeSummaryResults() []Result {
	return []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMismatch},
		{Key: "C", Status: StatusMissingInA},
		{Key: "D", Status: StatusMissingInB},
		{Key: "E", Status: StatusMissingInB},
	}
}

func TestSummarize_Counts(t *testing.T) {
	s := Summarize(makeSummaryResults())
	if s.Total != 5 {
		t.Errorf("Total: got %d, want 5", s.Total)
	}
	if s.Match != 1 {
		t.Errorf("Match: got %d, want 1", s.Match)
	}
	if s.Mismatch != 1 {
		t.Errorf("Mismatch: got %d, want 1", s.Mismatch)
	}
	if s.MissingInA != 1 {
		t.Errorf("MissingInA: got %d, want 1", s.MissingInA)
	}
	if s.MissingInB != 2 {
		t.Errorf("MissingInB: got %d, want 2", s.MissingInB)
	}
}

func TestSummarize_Empty(t *testing.T) {
	s := Summarize([]Result{})
	if s.Total != 0 || s.HasDiff() {
		t.Errorf("expected empty summary with no diff")
	}
}

func TestSummary_HasDiff_True(t *testing.T) {
	s := Summarize(makeSummaryResults())
	if !s.HasDiff() {
		t.Error("expected HasDiff to be true")
	}
}

func TestSummary_HasDiff_False(t *testing.T) {
	results := []Result{
		{Key: "X", Status: StatusMatch},
		{Key: "Y", Status: StatusMatch},
	}
	s := Summarize(results)
	if s.HasDiff() {
		t.Error("expected HasDiff to be false")
	}
}
