package diff

import (
	"bytes"
	"testing"
)

func makeAnnotateResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMismatch},
		{Key: "API_KEY", Status: StatusMissingInB},
		{Key: "PORT", Status: StatusMatch},
	}
}

func TestAnnotateResults_WithMatches(t *testing.T) {
	opts := DefaultAnnotateOptions()
	opts.Annotations["DB_HOST"] = "Database hostname"
	opts.Annotations["API_KEY"] = "Third-party API key"

	anns := AnnotateResults(makeAnnotateResults(), opts)
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	if anns[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY first (sorted), got %s", anns[0].Key)
	}
}

func TestAnnotateResults_NoMatches(t *testing.T) {
	opts := DefaultAnnotateOptions()
	anns := AnnotateResults(makeAnnotateResults(), opts)
	if len(anns) != 0 {
		t.Fatalf("expected 0 annotations, got %d", len(anns))
	}
}

func TestWriteAnnotationReport_WithAnnotations(t *testing.T) {
	anns := []Annotation{
		{Key: "DB_HOST", Comment: "Database hostname"},
	}
	var buf bytes.Buffer
	WriteAnnotationReport(&buf, anns)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty report")
	}
	if !bytes.Contains([]byte(out), []byte("DB_HOST")) {
		t.Error("expected DB_HOST in report")
	}
}

func TestWriteAnnotationReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteAnnotationReport(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("No annotations")) {
		t.Error("expected 'No annotations' message")
	}
}
