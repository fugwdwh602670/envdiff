package diff

import (
	"fmt"
	"io"
	"sort"
)

// PromoteOptions controls how values are promoted from one environment to another.
type PromoteOptions struct {
	// OverwriteExisting replaces keys that already exist in the target.
	OverwriteExisting bool
	// OnlyMissing promotes only keys missing in the target (ignores mismatches).
	OnlyMissing bool
	// DryRun reports what would be promoted without modifying the target map.
	DryRun bool
}

// DefaultPromoteOptions returns sensible defaults for promotion.
func DefaultPromoteOptions() PromoteOptions {
	return PromoteOptions{
		OverwriteExisting: false,
		OnlyMissing:       true,
		DryRun:            false,
	}
}

// PromoteResult describes a single key promotion action.
type PromoteResult struct {
	Key      string
	OldValue string // value in target before promotion (empty if missing)
	NewValue string // value from source
	Action   string // "added" | "overwritten" | "skipped"
}

// Promote copies keys from source into target according to opts.
// It returns the (possibly modified) target map and a slice of PromoteResults.
func Promote(source, target map[string]string, opts PromoteOptions) (map[string]string, []PromoteResult) {
	out := make(map[string]string, len(target))
	for k, v := range target {
		out[k] = v
	}

	var results []PromoteResult
	keys := sortedKeys(source)

	for _, k := range keys {
		srcVal := source[k]
		tgtVal, exists := out[k]

		switch {
		case !exists:
			results = append(results, PromoteResult{Key: k, OldValue: "", NewValue: srcVal, Action: "added"})
			if !opts.DryRun {
				out[k] = srcVal
			}
		case exists && tgtVal == srcVal:
			// values are identical — skip silently
		case exists && opts.OnlyMissing:
			results = append(results, PromoteResult{Key: k, OldValue: tgtVal, NewValue: srcVal, Action: "skipped"})
		case exists && opts.OverwriteExisting:
			results = append(results, PromoteResult{Key: k, OldValue: tgtVal, NewValue: srcVal, Action: "overwritten"})
			if !opts.DryRun {
				out[k] = srcVal
			}
		default:
			results = append(results, PromoteResult{Key: k, OldValue: tgtVal, NewValue: srcVal, Action: "skipped"})
		}
	}

	return out, results
}

// WritePromoteReport writes a human-readable promotion summary to w.
func WritePromoteReport(w io.Writer, results []PromoteResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No keys promoted.")
		return
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	for _, r := range results {
		switch r.Action {
		case "added":
			fmt.Fprintf(w, "[added]       %s = %s\n", r.Key, r.NewValue)
		case "overwritten":
			fmt.Fprintf(w, "[overwritten] %s: %s -> %s\n", r.Key, r.OldValue, r.NewValue)
		case "skipped":
			fmt.Fprintf(w, "[skipped]     %s (target=%s, source=%s)\n", r.Key, r.OldValue, r.NewValue)
		}
	}
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
