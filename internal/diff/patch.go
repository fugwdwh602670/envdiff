package diff

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// PatchOptions controls how a patch is generated or applied.
type PatchOptions struct {
	// OnlyMissing limits the patch to keys missing in the target.
	OnlyMissing bool
	// OnlyMismatch limits the patch to keys whose values differ.
	OnlyMismatch bool
	// DryRun prints what would change without writing anything.
	DryRun bool
}

// DefaultPatchOptions returns sensible defaults for patch generation.
func DefaultPatchOptions() PatchOptions {
	return PatchOptions{
		OnlyMissing:  false,
		OnlyMismatch: false,
		DryRun:       false,
	}
}

// PatchEntry represents a single key/value change in a patch.
type PatchEntry struct {
	Key      string
	Value    string
	Status   string // "add" or "update"
}

// BuildPatch produces a list of PatchEntry values from diff results.
// It describes the changes needed to bring the target env up to date with
// the source env based on the supplied options.
func BuildPatch(results []Result, opts PatchOptions) []PatchEntry {
	var entries []PatchEntry

	for _, r := range results {
		switch r.Status {
		case StatusMissingInB:
			if !opts.OnlyMismatch {
				entries = append(entries, PatchEntry{
					Key:    r.Key,
					Value:  r.ValueA,
					Status: "add",
				})
			}
		case StatusMismatch:
			if !opts.OnlyMissing {
				entries = append(entries, PatchEntry{
					Key:    r.Key,
					Value:  r.ValueA,
					Status: "update",
				})
			}
		}
	}

	// Stable output order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// ApplyPatch writes the patched env content to dst by merging entries into
// the existing env map. Keys present in entries overwrite or extend the map.
func ApplyPatch(existing map[string]string, entries []PatchEntry, dst io.Writer) error {
	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		merged[k] = v
	}
	for _, e := range entries {
		merged[e.Key] = e.Value
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w := bufio.NewWriter(dst)
	for _, k := range keys {
		_, err := fmt.Fprintf(w, "%s=%s\n", k, merged[k])
		if err != nil {
			return err
		}
	}
	return w.Flush()
}

// WritePatchReport writes a human-readable summary of patch entries to w.
func WritePatchReport(entries []PatchEntry, w io.Writer) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No patch entries — target is up to date.")
		return
	}

	var adds, updates []PatchEntry
	for _, e := range entries {
		if e.Status == "add" {
			adds = append(adds, e)
		} else {
			updates = append(updates, e)
		}
	}

	if len(adds) > 0 {
		fmt.Fprintln(w, "Keys to add:")
		for _, e := range adds {
			fmt.Fprintf(w, "  + %s=%s\n", e.Key, e.Value)
		}
	}
	if len(updates) > 0 {
		fmt.Fprintln(w, "Keys to update:")
		for _, e := range updates {
			fmt.Fprintf(w, "  ~ %s=%s\n", e.Key, e.Value)
		}
	}

	total := len(entries)
	fmt.Fprintf(w, "\nTotal: %d change(s) (%d add, %d update)\n", total, len(adds), len(updates))
}

// FilterPatchEntries returns only the entries whose keys match the given
// prefix. This is useful when callers want to scope a patch to a subset of
// keys sharing a common naming convention (e.g. "DB_", "AWS_").
func FilterPatchEntries(entries []PatchEntry, prefix string) []PatchEntry {
	var filtered []PatchEntry
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
