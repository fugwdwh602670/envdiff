package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoDiff(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res := diff.Compare(a, b)
	if res.HasDiff() {
		t.Errorf("expected no diff, got: %s", res)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}
	res := diff.Compare(a, b)
	if len(res.MissingInB) != 1 || res.MissingInB[0] != "ONLY_A" {
		t.Errorf("expected ONLY_A missing in B, got %v", res.MissingInB)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}
	res := diff.Compare(a, b)
	if len(res.MissingInA) != 1 || res.MissingInA[0] != "ONLY_B" {
		t.Errorf("expected ONLY_B missing in A, got %v", res.MissingInA)
	}
}

func TestCompare_Mismatch(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "baz"}
	res := diff.Compare(a, b)
	if len(res.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(res.Mismatched))
	}
	vals := res.Mismatched["FOO"]
	if vals[0] != "bar" || vals[1] != "baz" {
		t.Errorf("unexpected mismatch values: %v", vals)
	}
}

func TestCompare_Mixed(t *testing.T) {
	a := map[string]string{"SHARED": "same", "CHANGED": "old", "GONE": "x"}
	b := map[string]string{"SHARED": "same", "CHANGED": "new", "NEW": "y"}
	res := diff.Compare(a, b)
	if !res.HasDiff() {
		t.Fatal("expected diff")
	}
	if len(res.MissingInB) != 1 || res.MissingInB[0] != "GONE" {
		t.Errorf("unexpected MissingInB: %v", res.MissingInB)
	}
	if len(res.MissingInA) != 1 || res.MissingInA[0] != "NEW" {
		t.Errorf("unexpected MissingInA: %v", res.MissingInA)
	}
	if _, ok := res.Mismatched["CHANGED"]; !ok {
		t.Error("expected CHANGED in Mismatched")
	}
}
