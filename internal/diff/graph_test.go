package diff

import (
	"bytes"
	"strings"
	"testing"
)

func makeGraphEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"APP_NAME":    "envdiff",
		"APP_VERSION": "1.0",
		"STANDALONE":  "yes",
	}
}

func TestBuildGraph_Groups(t *testing.T) {
	env := makeGraphEnv()
	opts := DefaultGraphOptions()
	root := BuildGraph(env, opts)

	if root.Label != "(root)" {
		t.Errorf("expected root label (root), got %s", root.Label)
	}
	if len(root.Children) != 2 {
		t.Errorf("expected 2 group children, got %d", len(root.Children))
	}
	if len(root.Keys) != 1 || root.Keys[0] != "STANDALONE" {
		t.Errorf("expected STANDALONE in root keys, got %v", root.Keys)
	}
}

func TestBuildGraph_ChildKeys(t *testing.T) {
	env := makeGraphEnv()
	opts := DefaultGraphOptions()
	root := BuildGraph(env, opts)

	var dbNode *GraphNode
	for _, c := range root.Children {
		if c.Label == "DB" {
			dbNode = c
		}
	}
	if dbNode == nil {
		t.Fatal("expected DB group node")
	}
	if len(dbNode.Keys) != 2 {
		t.Errorf("expected 2 keys in DB group, got %d", len(dbNode.Keys))
	}
}

func TestBuildGraph_EmptyEnv(t *testing.T) {
	root := BuildGraph(map[string]string{}, DefaultGraphOptions())
	if len(root.Children) != 0 || len(root.Keys) != 0 {
		t.Error("expected empty graph for empty env")
	}
}

func TestWriteGraphReport_NoValues(t *testing.T) {
	env := makeGraphEnv()
	opts := DefaultGraphOptions()
	root := BuildGraph(env, opts)

	var buf bytes.Buffer
	WriteGraphReport(&buf, root, env, opts)
	out := buf.String()

	if !strings.Contains(out, "DB/") {
		t.Error("expected DB/ group in output")
	}
	if !strings.Contains(out, "- STANDALONE") {
		t.Error("expected STANDALONE in output")
	}
	if strings.Contains(out, "localhost") {
		t.Error("values should not appear when ShowValues=false")
	}
}

func TestWriteGraphReport_WithValues(t *testing.T) {
	env := makeGraphEnv()
	opts := DefaultGraphOptions()
	opts.ShowValues = true
	root := BuildGraph(env, opts)

	var buf bytes.Buffer
	WriteGraphReport(&buf, root, env, opts)
	out := buf.String()

	if !strings.Contains(out, "localhost") {
		t.Error("expected value 'localhost' in output when ShowValues=true")
	}
}
