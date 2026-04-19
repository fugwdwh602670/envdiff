package diff

import "io"
import "fmt"

// EnvStats holds aggregate statistics for a single env file.
type EnvStats struct {
	TotalKeys    int
	EmptyValues  int
	UniqueKeys   int
	DuplicateKeys []string
}

// StatEnv computes statistics from a parsed env map.
func StatEnv(env map[string]string) EnvStats {
	seen := make(map[string]int)
	for k := range env {
		seen[k]++
	}

	stats := EnvStats{
		TotalKeys:  len(env),
		UniqueKeys: len(seen),
	}

	for k, v := range env {
		if v == "" {
			stats.EmptyValues++
		}
		if seen[k] > 1 {
			stats.DuplicateKeys = append(stats.DuplicateKeys, k)
		}
	}

	return stats
}

// WriteStatsReport writes a human-readable stats report to w.
func WriteStatsReport(w io.Writer, label string, stats EnvStats) {
	fmt.Fprintf(w, "=== Stats: %s ===\n", label)
	fmt.Fprintf(w, "  Total keys   : %d\n", stats.TotalKeys)
	fmt.Fprintf(w, "  Unique keys  : %d\n", stats.UniqueKeys)
	fmt.Fprintf(w, "  Empty values : %d\n", stats.EmptyValues)
	if len(stats.DuplicateKeys) > 0 {
		fmt.Fprintf(w, "  Duplicates   : %v\n", stats.DuplicateKeys)
	} else {
		fmt.Fprintf(w, "  Duplicates   : none\n")
	}
}
