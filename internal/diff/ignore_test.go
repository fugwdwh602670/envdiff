package diff

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "ignore-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestIgnoreList_AddAndContains(t *testing.T) {
	il := NewIgnoreList()
	il.Add("SECRET_KEY")
	il.Add("  DB_PASS  ")

	if !il.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if !il.Contains("DB_PASS") {
		t.Error("expected DB_PASS to be ignored (trimmed)")
	}
	if il.Contains("OTHER") {
		t.Error("OTHER should not be in ignore list")
	}
}

func TestLoadIgnoreFile(t *testing.T) {
	path := writeTempIgnore(t, "# comment\nSECRET_KEY\n\nDB_PASS\n# another comment\nAPI_TOKEN\n")
	il, err := LoadIgnoreFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if il.Len() != 3 {
		t.Errorf("expected 3 keys, got %d", il.Len())
	}
	for _, key := range []string{"SECRET_KEY", "DB_PASS", "API_TOKEN"} {
		if !il.Contains(key) {
			t.Errorf("expected %s in ignore list", key)
		}
	}
}

func TestLoadIgnoreFile_NotExist(t *testing.T) {
	_, err := LoadIgnoreFile("/nonexistent/path.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestApplyIgnoreList(t *testing.T) {
	results := []Result{
		{Key: "SECRET_KEY", Status: MissingInB},
		{Key: "APP_ENV", Status: Mismatch},
		{Key: "DB_PASS", Status: MissingInA},
	}
	il := NewIgnoreList()
	il.Add("SECRET_KEY")
	il.Add("DB_PASS")

	out := ApplyIgnoreList(results, il)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", out[0].Key)
	}
}

func TestApplyIgnoreList_Nil(t *testing.T) {
	results := []Result{{Key: "A", Status: Mismatch}}
	out := ApplyIgnoreList(results, nil)
	if len(out) != 1 {
		t.Error("nil ignore list should return all results")
	}
}
