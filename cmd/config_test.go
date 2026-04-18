package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/yourorg/envdiff/cmd"
	"github.com/yourorg/envdiff/internal/config"
)

func runConfigCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestConfigShow_Defaults(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".envdiff.json")
	t.Setenv("ENVDIFF_CONFIG", cfgPath)

	out, err := runConfigCmd("config", "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected output")
	}
}

func TestConfigSetFormat_Valid(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".envdiff.json")
	t.Setenv("ENVDIFF_CONFIG", cfgPath)

	_, err := runConfigCmd("config", "set-format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg, _ := config.Load(cfgPath)
	if cfg.DefaultFormat != "json" {
		t.Errorf("expected json, got %s", cfg.DefaultFormat)
	}
}

func TestConfigSetFormat_Invalid(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVDIFF_CONFIG", filepath.Join(dir, ".envdiff.json"))

	_, err := runConfigCmd("config", "set-format", "yaml")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestConfigShow_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".envdiff.json")
	cfg := config.Config{
		DefaultFormat: "csv",
		IgnoreKeys:    []string{"DEBUG"},
		ShowMissing:   true,
		ShowMismatch:  false,
	}
	_ = config.Save(cfgPath, cfg)
	_ = os.Setenv("ENVDIFF_CONFIG", cfgPath)
	defer os.Unsetenv("ENVDIFF_CONFIG")

	out, err := runConfigCmd("config", "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
	_ = cobra.Command{} // ensure import used
}
