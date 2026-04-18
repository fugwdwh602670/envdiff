package diff

import (
	"fmt"
	"strings"
)

// OutputFormat represents the output format for diff results.
type OutputFormat string

const (
	FormatText OutputFormat = "text"
	FormatJSON OutputFormat = "json"
	FormatCSV  OutputFormat = "csv"
)

// FormatResults formats a slice of Result into a string using the given format.
func FormatResults(results []Result, format OutputFormat, fileA, fileB string) (string, error) {
	switch format {
	case FormatText:
		return formatText(results, fileA, fileB), nil
	case FormatJSON:
		return formatJSON(results, fileA, fileB), nil
	case FormatCSV:
		return formatCSV(results), nil
	default:
		return "", fmt.Errorf("unknown format: %q", format)
	}
}

func formatText(results []Result, fileA, fileB string) string {
	if len(results) == 0 {
		return "No differences found.\n"
	}
	var sb strings.Builder
	for _, r := range results {
		switch r.Kind {
		case MissingInB:
			fmt.Fprintf(&sb, "[MISSING in %s] %s\n", fileB, r.Key)
		case MissingInA:
			fmt.Fprintf(&sb, "[MISSING in %s] %s\n", fileA, r.Key)
		case Mismatch:
			fmt.Fprintf(&sb, "[MISMATCH]      %s: %q (in %s) vs %q (in %s)\n", r.Key, r.ValueA, fileA, r.ValueB, fileB)
		}
	}
	return sb.String()
}

func formatJSON(results []Result, fileA, fileB string) string {
	if len(results) == 0 {
		return "[]\n"
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, r := range results {
		sb.WriteString(fmt.Sprintf(
			"  {\"key\": %q, \"kind\": %q, \"value_a\": %q, \"value_b\": %q}",
			r.Key, r.Kind, r.ValueA, r.ValueB,
		))
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_ = fileA
	_ = fileB
	return sb.String()
}

func formatCSV(results []Result) string {
	var sb strings.Builder
	sb.WriteString("key,kind,value_a,value_b\n")
	for _, r := range results {
		fmt.Fprintf(&sb, "%s,%s,%s,%s\n", r.Key, r.Kind, r.ValueA, r.ValueB)
	}
	return sb.String()
}
