package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// GroupOptions controls how results are grouped.
type GroupOptions struct {
	// GroupBy can be "prefix" (split key on separator) or "status"
	GroupBy   string
	Separator string
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		GroupBy:   "prefix",
		Separator: "_",
	}
}

// GroupResults groups a slice of Result values according to opts.
// It returns a map of group label -> results in that group.
func GroupResults(results []Result, opts GroupOptions) map[string][]Result {
	groups := make(map[string][]Result)

	for _, r := range results {
		var label string
		switch strings.ToLower(opts.GroupBy) {
		case "status":
			label = string(r.Status)
		default: // "prefix"
			sep := opts.Separator
			if sep == "" {
				sep = "_"
			}
			parts := strings.SplitN(r.Key, sep, 2)
			if len(parts) > 1 {
				label = parts[0]
			} else {
				label = "(ungrouped)"
			}
		}
		groups[label] = append(groups[label], r)
	}
	return groups
}

// WriteGroupReport writes a human-readable grouped report to w.
func WriteGroupReport(w io.Writer, groups map[string][]Result) {
	// Sort group labels for deterministic output.
	labels := make([]string, 0, len(groups))
	for l := range groups {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	for _, label := range labels {
		fmt.Fprintf(w, "[%s]\n", label)
		results := groups[label]
		sort.Slice(results, func(i, j int) bool {
			return results[i].Key < results[j].Key
		})
		for _, r := range results {
			switch r.Status {
			case StatusMissingInB:
				fmt.Fprintf(w, "  - %s (missing in B)\n", r.Key)
			case StatusMissingInA:
				fmt.Fprintf(w, "  + %s (missing in A)\n", r.Key)
			case StatusMismatch:
				fmt.Fprintf(w, "  ~ %s (%q != %q)\n", r.Key, r.ValueA, r.ValueB)
			default:
				fmt.Fprintf(w, "    %s (ok)\n", r.Key)
			}
		}
		fmt.Fprintln(w)
	}
}
