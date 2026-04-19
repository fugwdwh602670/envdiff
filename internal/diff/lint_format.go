package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// WriteLintReport writes lint issues to w in a human-readable format.
func WriteLintReport(w io.Writer, issues []LintIssue, file string) {
	if len(issues) == 0 {
		fmt.Fprintf(w, "✔ No lint issues found in %s\n", file)
		return
	}

	sorted := make([]LintIssue, len(issues))
	copy(sorted, issues)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Message < sorted[j].Message
	})

	fmt.Fprintf(w, "Lint results for %s (%d issue(s)):\n", file, len(sorted))
	fmt.Fprintln(w, strings.Repeat("-", 40))
	for _, issue := range sorted {
		icon := "⚠"
		if issue.Severity == LintError {
			icon = "✗"
		}
		fmt.Fprintf(w, "%s [%s] %s: %s\n", icon, strings.ToUpper(string(issue.Severity)), issue.Key, issue.Message)
	}
}
