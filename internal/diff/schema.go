package diff

import (
	"fmt"
	"strings"
)

// SchemaRule defines an expected key and its constraints.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  string // optional: expected value pattern (prefix match)
}

// SchemaOptions controls schema validation behavior.
type SchemaOptions struct {
	Rules []SchemaRule
}

// SchemaIssue represents a single schema violation.
type SchemaIssue struct {
	Key     string
	Message string
}

// ValidateSchema checks a parsed env map against a set of schema rules.
func ValidateSchema(env map[string]string, opts SchemaOptions) []SchemaIssue {
	var issues []SchemaIssue
	for _, rule := range opts.Rules {
		val, exists := env[rule.Key]
		if rule.Required && !exists {
			issues = append(issues, SchemaIssue{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}
		if exists && rule.Pattern != "" && !strings.HasPrefix(val, rule.Pattern) {
			issues = append(issues, SchemaIssue{
				Key:     rule.Key,
				Message: fmt.Sprintf("value does not match expected pattern %q", rule.Pattern),
			})
		}
	}
	return issues
}

// WriteSchemaReport writes schema issues in a human-readable format.
func WriteSchemaReport(issues []SchemaIssue, w interface{ WriteString(string) (int, error) }) {
	if len(issues) == 0 {
		w.WriteString("schema: no issues found\n")
		return
	}
	for _, issue := range issues {
		w.WriteString(fmt.Sprintf("schema [%s]: %s\n", issue.Key, issue.Message))
	}
}
