package diff

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a formatted diff report to w.
func WriteReport(w io.Writer, fileA, fileB string, r Result) {
	fmt.Fprintf(w, "Comparing %s <-> %s\n", fileA, fileB)

	if !r.HasDiff() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(r.MissingInA) > 0 {
		sorted := sortedCopy(r.MissingInA)
		fmt.Fprintf(w, "\nKeys missing in %s (%d):\n", fileA, len(sorted))
		for _, k := range sorted {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.MissingInB) > 0 {
		sorted := sortedCopy(r.MissingInB)
		fmt.Fprintf(w, "\nKeys missing in %s (%d):\n", fileB, len(sorted))
		for _, k := range sorted {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(r.Mismatched) > 0 {
		keys := make([]string, 0, len(r.Mismatched))
		for k := range r.Mismatched {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Fprintf(w, "\nMismatched values (%d):\n", len(keys))
		for _, k := range keys {
			v := r.Mismatched[k]
			fmt.Fprintf(w, "  ~ %s\n    %s: %q\n    %s: %q\n", k, fileA, v[0], fileB, v[1])
		}
	}
}

func sortedCopy(s []string) []string {
	cp := make([]string, len(s))
	copy(cp, s)
	sort.Strings(cp)
	return cp
}
