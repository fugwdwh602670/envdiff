package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func makeBaselineResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMismatch, ValueA: "localhost", ValueB: "prod-db"},
		{Key: "API_KEY", Status: StatusMissingInB},
	}
}

func TestSaveAndLoadBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := makeBaselineResults()

	if err := SaveBaseline(path, ".env", ".env.prod", results); err != nil {
		t.Fatalf("SaveBaseline error: %v", err)
	}

	b, err := LoadBaseline(path)
	if err != nil {
		t.Fatalf("LoadBaseline error: %v", err)
	}
	if b.FileA != ".env" || b.FileB != ".env.prod" {
		t.Errorf("unexpected file names: %s %s", b.FileA, b.FileB)
	}
	if len(b.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(b.Results))
	}
}

func TestLoadBaseline_NotExist(t *testing.T) {
	_, err := LoadBaseline("/nonexistent/baseline.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadBaseline_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadBaseline(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffBaseline_NewAndResolved(t *testing.T) {
	baseline := []Result{
		{Key: "OLD_KEY", Status: StatusMissingInB},
		{Key: "SHARED", Status: StatusMismatch},
	}
	current := []Result{
		{Key: "NEW_KEY", Status: StatusMissingInB},
		{Key: "SHARED", Status: StatusMismatch},
	}
	newIssues, resolved := DiffBaseline(baseline, current)
	if len(newIssues) != 1 || newIssues[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY as new issue, got %+v", newIssues)
	}
	if len(resolved) != 1 || resolved[0].Key != "OLD_KEY" {
		t.Errorf("expected OLD_KEY as resolved, got %+v", resolved)
	}
}

func TestDiffBaseline_NoDiff(t *testing.T) {
	results := makeBaselineResults()
	newIssues, resolved := DiffBaseline(results, results)
	if len(newIssues) != 0 || len(resolved) != 0 {
		t.Errorf("expected no diff, got new=%v resolved=%v", newIssues, resolved)
	}
}
