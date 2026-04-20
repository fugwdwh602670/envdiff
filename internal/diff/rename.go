package diff

import (
	"fmt"
	"io"
	"sort"
)

// RenameMap maps old key names to new key names.
type RenameMap map[string]string

// RenameOptions controls how key renames are applied.
type RenameOptions struct {
	// Renames maps old key -> new key.
	Renames RenameMap
	// IgnoreUnmatched skips keys not found in the rename map.
	IgnoreUnmatched bool
}

// DefaultRenameOptions returns sensible defaults.
func DefaultRenameOptions() RenameOptions {
	return RenameOptions{
		Renames:         RenameMap{},
		IgnoreUnmatched: true,
	}
}

// ApplyRenames rewrites keys in an env map according to the rename map.
// Keys not present in the rename map are kept as-is (unless IgnoreUnmatched
// is false, in which case they are dropped).
func ApplyRenames(env map[string]string, opts RenameOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if newKey, ok := opts.Renames[k]; ok {
			out[newKey] = v
		} else if opts.IgnoreUnmatched {
			out[k] = v
		}
	}
	return out
}

// WriteRenameReport writes a human-readable summary of applied renames to w.
func WriteRenameReport(w io.Writer, original, renamed map[string]string, renames RenameMap) {
	type entry struct{ from, to string }
	var applied []entry
	for from, to := range renames {
		if _, existed := original[from]; existed {
			applied = append(applied, entry{from, to})
		}
	}
	sort.Slice(applied, func(i, j int) bool { return applied[i].from < applied[j].from })

	if len(applied) == 0 {
		fmt.Fprintln(w, "No renames applied.")
		return
	}
	fmt.Fprintf(w, "Applied %d rename(s):\n", len(applied))
	for _, e := range applied {
		v := renamed[e.to]
		fmt.Fprintf(w, "  %s -> %s (value: %q)\n", e.from, e.to, v)
	}
}
