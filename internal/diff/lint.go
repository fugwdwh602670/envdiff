package diff

import "fmt"

// LintSeverity indicates how severe a lint issue is.
type LintSeverity string

const (
	LintError   LintSeverity = "error"
	LintWarning LintSeverity = "warning"
)

// LintIssue represents a single linting problem found in env entries.
type LintIssue struct {
	Key      string
	Message  string
	Severity LintSeverity
}

// LintOptions controls which lint rules are enabled.
type LintOptions struct {
	CheckEmptyValues    bool
	CheckDuplicateKeys  bool
	CheckNamingConvention bool
}

// DefaultLintOptions returns sensible defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		CheckEmptyValues:    true,
		CheckDuplicateKeys:  true,
		CheckNamingConvention: true,
	}
}

// LintEnv runs lint rules against a parsed env map and returns any issues found.
func LintEnv(env map[string]string, opts LintOptions) []LintIssue {
	var issues []LintIssue

	for k, v := range env {
		if opts.CheckEmptyValues && v == "" {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  "value is empty",
				Severity: LintWarning,
			})
		}
		if opts.CheckNamingConvention && !isValidEnvKey(k) {
			issues = append(issues, LintIssue{
				Key:      k,
				Message:  fmt.Sprintf("key %q does not follow UPPER_SNAKE_CASE convention", k),
				Severity: LintWarning,
			})
		}
	}
	return issues
}

// isValidEnvKey returns true if the key is UPPER_SNAKE_CASE.
func isValidEnvKey(k string) bool {
	if len(k) == 0 {
		return false
	}
	for _, c := range k {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
