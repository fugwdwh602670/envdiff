package diff

import (
	"strings"
	"testing"
)

func TestValidateSchema_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
	}
	opts := SchemaOptions{
		Rules: []SchemaRule{
			{Key: "APP_ENV", Required: true},
			{Key: "APP_PORT", Required: true},
		},
	}
	issues := ValidateSchema(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := map[string]string{}
	opts := SchemaOptions{
		Rules: []SchemaRule{
			{Key: "DATABASE_URL", Required: true},
		},
	}
	issues := ValidateSchema(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "DATABASE_URL" {
		t.Errorf("unexpected key: %s", issues[0].Key)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "sqlite://local",
	}
	opts := SchemaOptions{
		Rules: []SchemaRule{
			{Key: "DATABASE_URL", Required: true, Pattern: "postgres://"},
		},
	}
	issues := ValidateSchema(env, opts)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if !strings.Contains(issues[0].Message, "postgres://") {
		t.Errorf("expected pattern in message, got: %s", issues[0].Message)
	}
}

func TestValidateSchema_PatternMatch(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
	}
	opts := SchemaOptions{
		Rules: []SchemaRule{
			{Key: "DATABASE_URL", Required: true, Pattern: "postgres://"},
		},
	}
	issues := ValidateSchema(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestWriteSchemaReport_WithIssues(t *testing.T) {
	issues := []SchemaIssue{
		{Key: "SECRET_KEY", Message: "required key is missing"},
	}
	var sb strings.Builder
	WriteSchemaReport(issues, &sb)
	if !strings.Contains(sb.String(), "SECRET_KEY") {
		t.Errorf("expected key in output, got: %s", sb.String())
	}
}

func TestWriteSchemaReport_NoIssues(t *testing.T) {
	var sb strings.Builder
	WriteSchemaReport(nil, &sb)
	if !strings.Contains(sb.String(), "no issues") {
		t.Errorf("expected no-issues message, got: %s", sb.String())
	}
}
