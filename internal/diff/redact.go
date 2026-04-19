package diff

import "strings"

// RedactOptions controls which keys have their values redacted in output.
type RedactOptions struct {
	// Patterns is a list of key substrings that trigger redaction (case-insensitive).
	Patterns []string
}

// DefaultRedactOptions returns sensible default patterns for sensitive keys.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Patterns: []string{"secret", "password", "passwd", "token", "apikey", "api_key", "private"},
	}
}

// ShouldRedact returns true if the key matches any redaction pattern.
func (r RedactOptions) ShouldRedact(key string) bool {
	lower := strings.ToLower(key)
	for _, p := range r.Patterns {
		if strings.Contains(lower, strings.ToLower(p)) {
			return true
		}
	}
	return false
}

const redactedPlaceholder = "***REDACTED***"

// RedactResults returns a copy of results with sensitive values replaced.
func RedactResults(results []Result, opts RedactOptions) []Result {
	out := make([]Result, len(results))
	for i, r := range results {
		if opts.ShouldRedact(r.Key) {
			r.ValueA = redactedPlaceholder
			r.ValueB = redactedPlaceholder
		}
		out[i] = r
	}
	return out
}
