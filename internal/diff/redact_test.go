package diff

import (
	"testing"
)

func TestShouldRedact_Matches(t *testing.T) {
	opts := DefaultRedactOptions()
	cases := []string{"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN", "STRIPE_APIKEY", "PRIVATE_KEY"}
	for _, key := range cases {
		if !opts.ShouldRedact(key) {
			t.Errorf("expected %q to be redacted", key)
		}
	}
}

func TestShouldRedact_NoMatch(t *testing.T) {
	opts := DefaultRedactOptions()
	cases := []string{"APP_ENV", "PORT", "DATABASE_URL", "LOG_LEVEL"}
	for _, key := range cases {
		if opts.ShouldRedact(key) {
			t.Errorf("expected %q NOT to be redacted", key)
		}
	}
}

func TestRedactResults_SensitiveValues(t *testing.T) {
	input := []Result{
		{Key: "DB_PASSWORD", ValueA: "hunter2", ValueB: "s3cr3t", Status: StatusMismatch},
		{Key: "APP_ENV", ValueA: "prod", ValueB: "staging", Status: StatusMismatch},
		{Key: "AUTH_TOKEN", ValueA: "abc123", ValueB: "", Status: StatusMissingInB},
	}

	opts := DefaultRedactOptions()
	out := RedactResults(input, opts)

	if out[0].ValueA != redactedPlaceholder || out[0].ValueB != redactedPlaceholder {
		t.Errorf("DB_PASSWORD values should be redacted")
	}
	if out[1].ValueA != "prod" || out[1].ValueB != "staging" {
		t.Errorf("APP_ENV values should not be redacted")
	}
	if out[2].ValueA != redactedPlaceholder {
		t.Errorf("AUTH_TOKEN valueA should be redacted")
	}
}

func TestRedactResults_DoesNotMutateOriginal(t *testing.T) {
	input := []Result{
		{Key: "DB_PASSWORD", ValueA: "original", ValueB: "value", Status: StatusMismatch},
	}
	opts := DefaultRedactOptions()
	_ = RedactResults(input, opts)
	if input[0].ValueA != "original" {
		t.Errorf("original slice should not be mutated")
	}
}

func TestRedactResults_CustomPattern(t *testing.T) {
	opts := RedactOptions{Patterns: []string{"internal"}}
	input := []Result{
		{Key: "INTERNAL_HOST", ValueA: "10.0.0.1", ValueB: "10.0.0.2", Status: StatusMismatch},
		{Key: "PUBLIC_HOST", ValueA: "a.com", ValueB: "b.com", Status: StatusMismatch},
	}
	out := RedactResults(input, opts)
	if out[0].ValueA != redactedPlaceholder {
		t.Errorf("INTERNAL_HOST should be redacted with custom pattern")
	}
	if out[1].ValueA == redactedPlaceholder {
		t.Errorf("PUBLIC_HOST should not be redacted")
	}
}
