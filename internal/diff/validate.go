package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ValidateOptions controls which validation checks are performed.
type ValidateOptions struct {
	RequireUppercase bool
	ForbidEmpty      bool
	ForbidDuplicates bool
}

// DefaultValidateOptions returns sensible defaults.
func DefaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		RequireUppercase: true,
		ForbidEmpty:      true,
		ForbidDuplicates: true,
	}
}

// ValidationIssue describes a single validation problem.
type ValidationIssue struct {
	Key     string
	Message string
}

// ValidateEnv runs validation rules against a parsed env map.
func ValidateEnv(env map[string]string, opts ValidateOptions) []ValidationIssue {
	var issues []ValidationIssue
	seen := make(map[string]int)

	for k, v := range env {
		seen[k]++
		if opts.RequireUppercase && k != strings.ToUpper(k) {
			issues = append(issues, ValidationIssue{
				Key:     k,
				Message: "key should be uppercase",
			})
		}
		if opts.ForbidEmpty && strings.TrimSpace(v) == "" {
			issues = append(issues, ValidationIssue{
				Key:     k,
				Message: "value is empty",
			})
		}
	}

	if opts.ForbidDuplicates {
		for k, count := range seen {
			if count > 1 {
				issues = append(issues, ValidationIssue{
					Key:     k,
					Message: fmt.Sprintf("key appears %d times", count),
				})
			}
		}
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Key != issues[j].Key {
			return issues[i].Key < issues[j].Key
		}
		return issues[i].Message < issues[j].Message
	})
	return issues
}

// WriteValidateReport writes a human-readable validation report to w.
func WriteValidateReport(w io.Writer, issues []ValidationIssue, filename string) {
	if len(issues) == 0 {
		fmt.Fprintf(w, "✔ %s: no validation issues\n", filename)
		return
	}
	fmt.Fprintf(w, "✖ %s: %d issue(s) found\n", filename, len(issues))
	for _, iss := range issues {
		fmt.Fprintf(w, "  [%s] %s\n", iss.Key, iss.Message)
	}
}
