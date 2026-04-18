package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/config"
)

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected text, got %s", cfg.DefaultFormat)
	}
	if !cfg.ShowMissing || !cfg.ShowMismatch {
		t.Error("expected show flags to be true")
	}
}

func TestLoad_NotExist(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/.envdiff.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DefaultFormat != "text" {
		t.Errorf("expected default format, got %s", cfg.DefaultFormat)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiff.json")

	orig := config.Config{
		DefaultFormat: "json",
		IgnoreKeys:    []string{"SECRET", "TOKEN"},
		ShowMissing:   false,
		ShowMismatch:  true,
	}

	if err := config.Save(path, orig); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if loaded.DefaultFormat != orig.DefaultFormat {
		t.Errorf("format mismatch: %s vs %s", loaded.DefaultFormat, orig.DefaultFormat)
	}
	if loaded.ShowMissing != orig.ShowMissing {
		t.Errorf("show_missing mismatch")
	}
	if len(loaded.IgnoreKeys) != 2 || loaded.IgnoreKeys[0] != "SECRET" {
		t.Errorf("ignore_keys mismatch: %v", loaded.IgnoreKeys)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiff.json")
	_ = os.WriteFile(path, []byte("not json{"), 0644)
	_, err := config.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
