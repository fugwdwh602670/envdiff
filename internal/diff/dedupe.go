package diff

import (
	"fmt"
	"io"
	"sort"
)

// DedupeOptions controls deduplication behavior.
type DedupeOptions struct {
	// PreferLast keeps the last occurrence of a duplicate key instead of the first.
	PreferLast bool
	// ReportOnly returns duplicates without modifying the map.
	ReportOnly bool
}

// DefaultDedupeOptions returns sensible defaults.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{
		PreferLast: false,
		ReportOnly: false,
	}
}

// DuplicateEntry records a key that appeared more than once.
type DuplicateEntry struct {
	Key    string
	Values []string // all values seen, in order
	Kept   string   // the value that was retained
}

// DedupeResult is the output of DedupeEnv.
type DedupeResult struct {
	Env        map[string]string
	Duplicates []DuplicateEntry
}

// DedupeEnv scans an ordered list of raw key-value pairs (allowing duplicate
// keys) and collapses them into a single map, recording which keys were
// duplicated. Use ParseFileRaw (or build the pairs slice yourself) to preserve
// original ordering before deduplication.
func DedupeEnv(pairs [][2]string, opts DedupeOptions) DedupeResult {
	type entry struct {
		values []string
	}
	tracking := make(map[string]*entry)
	order := make([]string, 0)

	for _, p := range pairs {
		key, val := p[0], p[1]
		if e, ok := tracking[key]; ok {
			e.values = append(e.values, val)
		} else {
			tracking[key] = &entry{values: []string{val}}
			order = append(order, key)
		}
	}

	result := DedupeResult{
		Env: make(map[string]string),
	}

	for _, key := range order {
		e := tracking[key]
		var kept string
		if opts.PreferLast {
			kept = e.values[len(e.values)-1]
		} else {
			kept = e.values[0]
		}
		if len(e.values) > 1 {
			result.Duplicates = append(result.Duplicates, DuplicateEntry{
				Key:    key,
				Values: e.values,
				Kept:   kept,
			})
		}
		if !opts.ReportOnly {
			result.Env[key] = kept
		}
	}

	sort.Slice(result.Duplicates, func(i, j int) bool {
		return result.Duplicates[i].Key < result.Duplicates[j].Key
	})
	return result
}

// WriteDedupeReport writes a human-readable summary of duplicate keys.
func WriteDedupeReport(w io.Writer, result DedupeResult) {
	if len(result.Duplicates) == 0 {
		fmt.Fprintln(w, "No duplicate keys found.")
		return
	}
	fmt.Fprintf(w, "Found %d duplicate key(s):\n\n", len(result.Duplicates))
	for _, d := range result.Duplicates {
		fmt.Fprintf(w, "  %s (%d occurrences, kept: %q)\n", d.Key, len(d.Values), d.Kept)
		for i, v := range d.Values {
			fmt.Fprintf(w, "    [%d] %q\n", i+1, v)
		}
	}
}
