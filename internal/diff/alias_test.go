package diff

import (
	"bytes"
	"testing"
)

func makeAliasEnv() map[string]string {
	return map[string]string{
		"DB_HOST":       "localhost",
		"DATABASE_HOST": "remotehost",
		"API_KEY":       "secret",
		"UNRELATED":     "value",
	}
}

func TestResolveAliases_RenamesAlias(t *testing.T) {
	env := map[string]string{
		"DB_HOST":   "localhost",
		"UNRELATED": "value",
	}
	aliases := AliasMap{
		"DATABASE_HOST": []string{"DB_HOST"},
	}
	result := ResolveAliases(env, aliases)
	if _, ok := result["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if result["DATABASE_HOST"] != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", result["DATABASE_HOST"])
	}
	if result["UNRELATED"] != "value" {
		t.Error("expected UNRELATED to be preserved")
	}
}

func TestResolveAliases_DoesNotOverwriteCanonical(t *testing.T) {
	env := map[string]string{
		"DB_HOST":       "alias-value",
		"DATABASE_HOST": "canonical-value",
	}
	aliases := AliasMap{
		"DATABASE_HOST": []string{"DB_HOST"},
	}
	result := ResolveAliases(env, aliases)
	// canonical already exists; alias value should not overwrite it
	if result["DATABASE_HOST"] != "canonical-value" {
		t.Errorf("expected canonical-value, got %q", result["DATABASE_HOST"])
	}
}

func TestResolveAliases_EmptyAliasMap(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	result := ResolveAliases(env, AliasMap{})
	if result["FOO"] != "bar" {
		t.Error("expected FOO to be preserved with empty alias map")
	}
}

func TestAuditAliases_DetectsMatches(t *testing.T) {
	env := makeAliasEnv()
	aliases := AliasMap{
		"DATABASE_HOST": []string{"DB_HOST"},
		"SECRET_KEY":    []string{"API_KEY"},
	}
	reports := AuditAliases(env, aliases)
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
	if reports[0].Alias != "API_KEY" || reports[0].Canonical != "SECRET_KEY" {
		t.Errorf("unexpected first report: %+v", reports[0])
	}
	if reports[1].Alias != "DB_HOST" || reports[1].Canonical != "DATABASE_HOST" {
		t.Errorf("unexpected second report: %+v", reports[1])
	}
}

func TestAuditAliases_NoMatches(t *testing.T) {
	env := map[string]string{"UNRELATED": "val"}
	aliases := AliasMap{"DATABASE_HOST": []string{"DB_HOST"}}
	reports := AuditAliases(env, aliases)
	if len(reports) != 0 {
		t.Errorf("expected 0 reports, got %d", len(reports))
	}
}

func TestWriteAliasReport_WithResults(t *testing.T) {
	var buf bytes.Buffer
	reports := []AliasReport{
		{Alias: "DB_HOST", Canonical: "DATABASE_HOST", Value: "localhost"},
	}
	WriteAliasReport(&buf, reports)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("DB_HOST")) {
		t.Error("expected DB_HOST in output")
	}
}

func TestWriteAliasReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	WriteAliasReport(&buf, nil)
	if !bytes.Contains(buf.Bytes(), []byte("No alias")) {
		t.Error("expected 'No alias' message for empty report")
	}
}
