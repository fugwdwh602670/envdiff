package diff

import (
	"bytes"
	"testing"
)

func makeValidateEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestValidateEnv_NoIssues(t *testing.T) {
	env := makeValidateEnv("HOST", "localhost", "PORT", "8080")
	issues := ValidateEnv(env, DefaultValidateOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestValidateEnv_LowercaseKey(t *testing.T) {
	env := makeValidateEnv("host", "localhost")
	opts := DefaultValidateOptions()
	issues := ValidateEnv(env, opts)
	if len(issues) != 1 || issues[0].Key != "host" {
		t.Errorf("expected 1 uppercase issue for 'host', got %v", issues)
	}
}

func TestValidateEnv_EmptyValue(t *testing.T) {
	env := makeValidateEnv("API_KEY", "")
	opts := DefaultValidateOptions()
	issues := ValidateEnv(env, opts)
	if len(issues) != 1 || issues[0].Message != "value is empty" {
		t.Errorf("expected empty value issue, got %v", issues)
	}
}

func TestValidateEnv_DisabledRules(t *testing.T) {
	env := makeValidateEnv("lower_key", "")
	opts := ValidateOptions{
		RequireUppercase: false,
		ForbidEmpty:      false,
		ForbidDuplicates: false,
	}
	issues := ValidateEnv(env, opts)
	if len(issues) != 0 {
		t.Errorf("expected no issues with all rules disabled, got %v", issues)
	}
}

func TestValidateEnv_MultipleIssues(t *testing.T) {
	env := makeValidateEnv("bad_key", "", "GOOD", "val")
	opts := DefaultValidateOptions()
	issues := ValidateEnv(env, opts)
	// bad_key triggers uppercase + empty
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %v", issues)
	}
}

func TestWriteValidateReport_NoIssues(t *testing.T) {
	var buf bytes.Buffer
	WriteValidateReport(&buf, nil, "test.env")
	if got := buf.String(); got == "" {
		t.Error("expected non-empty output")
	}
	if got := buf.String(); got[:1] != "✔"[:3] {
		// just check it contains "no validation issues"
		if !bytes.Contains(buf.Bytes(), []byte("no validation issues")) {
			t.Errorf("unexpected output: %s", got)
		}
	}
}

func TestWriteValidateReport_WithIssues(t *testing.T) {
	issues := []ValidationIssue{
		{Key: "foo", Message: "key should be uppercase"},
	}
	var buf bytes.Buffer
	WriteValidateReport(&buf, issues, "test.env")
	if !bytes.Contains(buf.Bytes(), []byte("foo")) {
		t.Errorf("expected key 'foo' in output: %s", buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("1 issue(s)")) {
		t.Errorf("expected issue count in output: %s", buf.String())
	}
}
