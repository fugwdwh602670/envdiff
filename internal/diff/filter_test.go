package diff

import (
	"testing"
)

func makeResults() []Result {
	return []Result{
		{Key: "ONLY_A", Status: MissingInB},
		{Key: "ONLY_B", Status: MissingInA},
		{Key: "DIFF_VAL", ValueA: "foo", ValueB: "bar", Status: Mismatch},
	}
}

func TestFilter_AllEnabled(t *testing.T) {
	results := makeResults()
	opts := DefaultFilterOptions()
	out := Filter(results, opts)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestFilter_OnlyMissingInB(t *testing.T) {
	results := makeResults()
	opts := FilterOptions{ShowMissingInB: true}
	out := Filter(results, opts)
	if len(out) != 1 || out[0].Key != "ONLY_A" {
		t.Fatalf("expected only ONLY_A, got %+v", out)
	}
}

func TestFilter_OnlyMissingInA(t *testing.T) {
	results := makeResults()
	opts := FilterOptions{ShowMissingInA: true}
	out := Filter(results, opts)
	if len(out) != 1 || out[0].Key != "ONLY_B" {
		t.Fatalf("expected only ONLY_B, got %+v", out)
	}
}

func TestFilter_OnlyMismatch(t *testing.T) {
	results := makeResults()
	opts := FilterOptions{ShowMismatch: true}
	out := Filter(results, opts)
	if len(out) != 1 || out[0].Key != "DIFF_VAL" {
		t.Fatalf("expected only DIFF_VAL, got %+v", out)
	}
}

func TestFilter_NoneEnabled(t *testing.T) {
	results := makeResults()
	opts := FilterOptions{}
	out := Filter(results, opts)
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}
