package diff

import (
	"strings"
	"testing"
)

func TestFormatResults_NoDiff_Text(t *testing.T) {
	out, err := FormatResults(nil, FormatText, "a.env", "b.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %q", out)
	}
}

func TestFormatResults_NoDiff_JSON(t *testing.T) {
	out, err := FormatResults(nil, FormatJSON, "a.env", "b.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("expected empty JSON array, got: %q", out)
	}
}

func TestFormatResults_Text(t *testing.T) {
	results := []Result{
		{Key: "FOO", Kind: MissingInB, ValueA: "bar", ValueB: ""},
		{Key: "BAZ", Kind: Mismatch, ValueA: "old", ValueB: "new"},
	}
	out, err := FormatResults(results, FormatText, "dev.env", "prod.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING in output, got: %q", out)
	}
	if !strings.Contains(out, "MISMATCH") {
		t.Errorf("expected MISMATCH in output, got: %q", out)
	}
	if !strings.Contains(out, "prod.env") {
		t.Errorf("expected prod.env in output, got: %q", out)
	}
}

func TestFormatResults_JSON(t *testing.T) {
	results := []Result{
		{Key: "KEY", Kind: MissingInA, ValueA: "", ValueB: "val"},
	}
	out, err := FormatResults(results, FormatJSON, "a.env", "b.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"KEY\"") {
		t.Errorf("expected key in JSON, got: %q", out)
	}
	if !strings.Contains(out, "missing_in_a") {
		t.Errorf("expected kind in JSON, got: %q", out)
	}
}

func TestFormatResults_CSV(t *testing.T) {
	results := []Result{
		{Key: "X", Kind: Mismatch, ValueA: "1", ValueB: "2"},
	}
	out, err := FormatResults(results, FormatCSV, "a.env", "b.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "key,kind,value_a,value_b") {
		t.Errorf("expected CSV header, got: %q", out)
	}
	if !strings.Contains(out, "X,mismatch,1,2") {
		t.Errorf("expected CSV row, got: %q", out)
	}
}

func TestFormatResults_UnknownFormat(t *testing.T) {
	_, err := FormatResults(nil, OutputFormat("xml"), "a.env", "b.env")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
