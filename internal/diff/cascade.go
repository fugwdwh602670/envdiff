package diff

import (
	"fmt"
	"io"
	"sort"
)

// CascadeOptions controls how cascading merge behaves.
type CascadeOptions struct {
	// Overwrite allows higher-priority values to overwrite lower ones.
	Overwrite bool
	// SkipEmpty ignores keys with empty values in source layers.
	SkipEmpty bool
}

// DefaultCascadeOptions returns sensible defaults.
func DefaultCascadeOptions() CascadeOptions {
	return CascadeOptions{
		Overwrite: true,
		SkipEmpty: false,
	}
}

// CascadeEntry records how a key was resolved across layers.
type CascadeEntry struct {
	Key        string
	Value      string
	SourceLayer int // index of the winning layer (0 = base)
}

// CascadeEnv merges multiple env layers in order, where later layers have
// higher priority. It returns the merged map and a trace of source layers.
func CascadeEnv(layers []map[string]string, opts CascadeOptions) (map[string]string, []CascadeEntry) {
	merged := make(map[string]string)
	trace := make(map[string]CascadeEntry)

	for i, layer := range layers {
		for k, v := range layer {
			if opts.SkipEmpty && v == "" {
				continue
			}
			if _, exists := merged[k]; !exists || opts.Overwrite {
				merged[k] = v
				trace[k] = CascadeEntry{Key: k, Value: v, SourceLayer: i}
			}
		}
	}

	entries := make([]CascadeEntry, 0, len(trace))
	for _, e := range trace {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return merged, entries
}

// WriteCascadeReport writes a human-readable cascade trace to w.
func WriteCascadeReport(w io.Writer, entries []CascadeEntry, layerNames []string) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No keys resolved.")
		return
	}
	for _, e := range entries {
		name := fmt.Sprintf("layer[%d]", e.SourceLayer)
		if e.SourceLayer < len(layerNames) {
			name = layerNames[e.SourceLayer]
		}
		fmt.Fprintf(w, "%-30s = %-30s (from %s)\n", e.Key, e.Value, name)
	}
}
