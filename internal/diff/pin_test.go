package diff

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func makePinEnv() map[string]string {
	return map[string]string{
		"APP_SECRET": "abc123",
		"DB_PASS":    "hunter2",
		"API_KEY":    "key-xyz",
	}
}

func TestPinEnv_BasicPin(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	added, err := PinEnv(env, []string{"APP_SECRET", "DB_PASS"}, path, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(added) != 2 {
		t.Errorf("expected 2 added, got %d", len(added))
	}

	pf, err := LoadPinFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(pf.Pins) != 2 {
		t.Errorf("expected 2 pins in file, got %d", len(pf.Pins))
	}
}

func TestPinEnv_NoOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	_, _ = PinEnv(env, []string{"APP_SECRET"}, path, DefaultPinOptions())

	env2 := map[string]string{"APP_SECRET": "changed"}
	added, err := PinEnv(env2, []string{"APP_SECRET"}, path, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("expected no additions (no overwrite), got %d", len(added))
	}

	pf, _ := LoadPinFile(path)
	if pf.Pins[0].Value != "abc123" {
		t.Errorf("expected original value preserved, got %s", pf.Pins[0].Value)
	}
}

func TestPinEnv_Overwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	_, _ = PinEnv(env, []string{"APP_SECRET"}, path, DefaultPinOptions())

	env2 := map[string]string{"APP_SECRET": "newval"}
	opts := DefaultPinOptions()
	opts.Overwrite = true
	added, err := PinEnv(env2, []string{"APP_SECRET"}, path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(added) != 1 {
		t.Errorf("expected 1 overwrite, got %d", len(added))
	}
	pf, _ := LoadPinFile(path)
	if pf.Pins[0].Value != "newval" {
		t.Errorf("expected updated value, got %s", pf.Pins[0].Value)
	}
}

func TestCheckPins_NoViolations(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	_, _ = PinEnv(env, []string{"APP_SECRET", "DB_PASS"}, path, DefaultPinOptions())

	violations, err := CheckPins(env, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheckPins_Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	_, _ = PinEnv(env, []string{"APP_SECRET"}, path, DefaultPinOptions())

	env["APP_SECRET"] = "tampered"
	violations, err := CheckPins(env, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Expected != "abc123" || violations[0].Actual != "tampered" {
		t.Errorf("unexpected violation values: %+v", violations[0])
	}
}

func TestCheckPins_MissingKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	env := makePinEnv()

	_, _ = PinEnv(env, []string{"API_KEY"}, path, DefaultPinOptions())

	delete(env, "API_KEY")
	violations, err := CheckPins(env, path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 1 || !violations[0].Missing {
		t.Errorf("expected missing violation, got %+v", violations)
	}
}

func TestLoadPinFile_NotExist(t *testing.T) {
	pf, err := LoadPinFile("/nonexistent/path/pins.json")
	if err != nil {
		t.Fatalf("expected empty pin file, got error: %v", err)
	}
	if len(pf.Pins) != 0 {
		t.Errorf("expected empty pins")
	}
}

func TestLoadPinFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadPinFile(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestWritePinReport_Output(t *testing.T) {
	var buf bytes.Buffer
	added := []PinEntry{{Key: "FOO", Value: "bar", Comment: "locked"}}
	violations := []PinViolation{{Key: "BAZ", Expected: "x", Actual: "y"}}
	WritePinReport(&buf, added, violations)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty report")
	}
}
