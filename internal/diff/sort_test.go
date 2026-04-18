package diff

import (
	"testing"
)

func makeSortResults() []Result {
	return []Result{
		{Key: "ZEBRA", Status: StatusMatch},
		{Key: "ALPHA", Status: StatusMismatch},
		{Key: "MONGO", Status: StatusMissingInA},
		{Key: "BETA", Status: StatusMissingInB},
		{Key: "DELTA", Status: StatusMatch},
	}
}

func TestSortResults_ByKey(t *testing.T) {
	results := makeSortResults()
	sorted := SortResults(results, SortByKey)

	expectedKeys := []string{"ALPHA", "BETA", "DELTA", "MONGO", "ZEBRA"}
	for i, key := range expectedKeys {
		if sorted[i].Key != key {
			t.Errorf("index %d: got %q, want %q", i, sorted[i].Key, key)
		}
	}
}

func TestSortResults_ByStatus(t *testing.T) {
	results := makeSortResults()
	sorted := SortResults(results, SortByStatus)

	// First should be MissingInB, then MissingInA, then Mismatch, then Match
	if sorted[0].Status != StatusMissingInB {
		t.Errorf("expected first status MissingInB, got %q", sorted[0].Status)
	}
	if sorted[1].Status != StatusMissingInA {
		t.Errorf("expected second status MissingInA, got %q", sorted[1].Status)
	}
	if sorted[2].Status != StatusMismatch {
		t.Errorf("expected third status Mismatch, got %q", sorted[2].Status)
	}
}

func TestSortResults_DoesNotMutateOriginal(t *testing.T) {
	results := makeSortResults()
	originalFirst := results[0].Key
	SortResults(results, SortByKey)
	if results[0].Key != originalFirst {
		t.Errorf("original slice was mutated")
	}
}

func TestSortResults_DefaultIsKey(t *testing.T) {
	results := makeSortResults()
	byDefault := SortResults(results, "")
	byKey := SortResults(results, SortByKey)
	for i := range byDefault {
		if byDefault[i].Key != byKey[i].Key {
			t.Errorf("default sort differs from key sort at index %d", i)
		}
	}
}
