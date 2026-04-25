package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func makeSnapshotEnv() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/dev",
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	env := makeSnapshotEnv()

	if err := SaveSnapshot(path, env, "test-label"); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if snap.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", snap.Label)
	}
	if snap.Env["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", snap.Env["APP_HOST"])
	}
	if snap.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
}

func TestLoadSnapshot_NotExist(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiffSnapshot_NoDiff(t *testing.T) {
	env := makeSnapshotEnv()
	snap := DefaultSnapshotOptions("baseline")
	for k, v := range env {
		snap.Env[k] = v
	}

	results := DiffSnapshot(snap, env)
	for _, r := range results {
		if r.Status != StatusMatch {
			t.Errorf("expected all matches, got %s for key %s", r.Status, r.Key)
		}
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	env := makeSnapshotEnv()
	snap := DefaultSnapshotOptions("v1")
	for k, v := range env {
		snap.Env[k] = v
	}

	live := map[string]string{
		"APP_HOST": "prod.example.com", // changed
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/dev",
		"NEW_KEY":  "new-value", // added
	}

	results := DiffSnapshot(snap, live)
	statuses := map[string]Status{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}

	if statuses["APP_HOST"] != StatusMismatch {
		t.Errorf("expected APP_HOST mismatch, got %s", statuses["APP_HOST"])
	}
	if statuses["NEW_KEY"] != StatusMissingInA {
		t.Errorf("expected NEW_KEY missing-in-a, got %s", statuses["NEW_KEY"])
	}
}
