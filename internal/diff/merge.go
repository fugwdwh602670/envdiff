package diff

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// MergeOptions controls how env files are merged.
type MergeOptions struct {
	// PreferA: when a key exists in both, use A's value
	PreferA bool
	// SkipMissing: omit keys missing in either file
	SkipMissing bool
}

// DefaultMergeOptions returns sensible defaults.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{
		PreferA:     false,
		SkipMissing: false,
	}
}

// Merge combines two env maps into a single map according to opts.
// By default (PreferA=false), values from B win on conflict.
func Merge(a, b map[string]string, opts MergeOptions) map[string]string {
	out := make(map[string]string)

	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		if existing, ok := out[k]; ok {
			if opts.PreferA {
				_ = existing // keep A
			} else {
				out[k] = v
			}
		} else {
			out[k] = v
		}
	}

	if opts.SkipMissing {
		for k := range a {
			if _, ok := b[k]; !ok {
				delete(out, k)
			}
		}
		for k := range b {
			if _, ok := a[k]; !ok {
				delete(out, k)
			}
		}
	}

	return out
}

// WriteMerged writes the merged env map to path in KEY=VALUE format.
func WriteMerged(merged map[string]string, path string) error {
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, merged[k])
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}
