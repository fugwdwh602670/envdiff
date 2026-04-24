package diff

import (
	"bytes"
	"testing"
)

func makeScopeEnv() map[string]string {
	return map[string]string{
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"REDIS_URL":   "redis://localhost",
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0",
		"SECRET_KEY":  "abc123",
	}
}

func TestScopeEnv_NoFilter(t *testing.T) {
	env := makeScopeEnv()
	opts := DefaultScopeOptions()
	result := ScopeEnv(env, opts)
	if len(result) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result))
	}
}

func TestScopeEnv_SinglePrefix(t *testing.T) {
	env := makeScopeEnv()
	opts := ScopeOptions{Prefixes: []string{"DB_"}}
	result := ScopeEnv(env, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := result["APP_NAME"]; ok {
		t.Error("did not expect APP_NAME in result")
	}
}

func TestScopeEnv_MultiPrefix(t *testing.T) {
	env := makeScopeEnv()
	opts := ScopeOptions{Prefixes: []string{"DB_", "REDIS_"}}
	result := ScopeEnv(env, opts)
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
}

func TestScopeEnv_Invert(t *testing.T) {
	env := makeScopeEnv()
	opts := ScopeOptions{Prefixes: []string{"DB_", "REDIS_"}, Invert: true}
	result := ScopeEnv(env, opts)
	if len(result) != 3 {
		t.Errorf("expected 3 keys, got %d", len(result))
	}
	if _, ok := result["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in inverted result")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST in inverted result")
	}
}

func TestSummarizeScope(t *testing.T) {
	env := makeScopeEnv()
	opts := ScopeOptions{Prefixes: []string{"DB_"}}
	scoped := ScopeEnv(env, opts)
	summary := SummarizeScope(env, scoped, opts)
	if summary.Total != 6 {
		t.Errorf("expected Total=6, got %d", summary.Total)
	}
	if summary.Included != 2 {
		t.Errorf("expected Included=2, got %d", summary.Included)
	}
	if summary.Excluded != 4 {
		t.Errorf("expected Excluded=4, got %d", summary.Excluded)
	}
}

func TestWriteScopeReport_Output(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	opts := ScopeOptions{Prefixes: []string{"DB_"}}
	summary := SummarizeScope(makeScopeEnv(), env, opts)
	var buf bytes.Buffer
	WriteScopeReport(&buf, env, summary)
	out := buf.String()
	if out == "" {
		t.Error("expected non-empty report")
	}
	if !containsStr(out, "DB_HOST") {
		t.Error("expected DB_HOST in report")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	}())
}
