package diff

import (
	"testing"
)

func TestLintEnv_NoIssues(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"PORT":     "8080",
	}
	issues := LintEnv(env, DefaultLintOptions())
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %d", len(issues))
	}
}

func TestLintEnv_EmptyValue(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "",
	}
	issues := LintEnv(env, DefaultLintOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
}

func TestLintEnv_NamingConvention(t *testing.T) {
	env := map[string]string{
		"myKey":   "val",
		"GOOD_KEY": "val",
	}
	issues := LintEnv(env, DefaultLintOptions())
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "myKey" {
		t.Errorf("expected myKey, got %s", issues[0].Key)
	}
}

func TestLintEnv_DisabledRules(t *testing.T) {
	env := map[string]string{
		"bad_key": "",
	}
	opts := LintOptions{
		CheckEmptyValues:      false,
		CheckNamingConvention: false,
	}
	issues := LintEnv(env, opts)
	if len(issues) != 0 {
		t.Fatalf("expected no issues with rules disabled, got %d", len(issues))
	}
}

func TestIsValidEnvKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"GOOD", true},
		{"GOOD_KEY", true},
		{"KEY123", true},
		{"bad", false},
		{"Mixed", false},
		{"HAS SPACE", false},
		{"", false},
	}
	for _, c := range cases {
		got := isValidEnvKey(c.key)
		if got != c.valid {
			t.Errorf("isValidEnvKey(%q) = %v, want %v", c.key, got, c.valid)
		}
	}
}
