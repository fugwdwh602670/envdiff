package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// TransformFunc is a function that transforms an env value.
type TransformFunc func(value string) string

// TransformOptions controls how transformations are applied.
type TransformOptions struct {
	// Transforms maps key patterns (prefix match) to transform functions.
	Transforms map[string]TransformFunc
	// ApplyToAll applies a single transform to every key if set.
	ApplyToAll TransformFunc
}

// DefaultTransformOptions returns an empty TransformOptions.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		Transforms: make(map[string]TransformFunc),
	}
}

// TransformResult holds the outcome of transforming a single key.
type TransformResult struct {
	Key      string
	Original string
	Result   string
	Changed  bool
}

// TransformEnv applies transforms to the given env map and returns
// a new map with transformed values along with a change log.
func TransformEnv(env map[string]string, opts TransformOptions) (map[string]string, []TransformResult) {
	out := make(map[string]string, len(env))
	var results []TransformResult

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := env[k]
		newVal := v

		if opts.ApplyToAll != nil {
			newVal = opts.ApplyToAll(newVal)
		}

		for pattern, fn := range opts.Transforms {
			if strings.HasPrefix(k, pattern) {
				newVal = fn(newVal)
				break
			}
		}

		out[k] = newVal
		results = append(results, TransformResult{
			Key:      k,
			Original: v,
			Result:   newVal,
			Changed:  newVal != v,
		})
	}

	return out, results
}

// WriteTransformReport writes a human-readable transform report to w.
func WriteTransformReport(w io.Writer, results []TransformResult) {
	changed := 0
	for _, r := range results {
		if r.Changed {
			changed++
		}
	}
	fmt.Fprintf(w, "Transform report: %d key(s) changed\n", changed)
	for _, r := range results {
		if r.Changed {
			fmt.Fprintf(w, "  ~ %s: %q => %q\n", r.Key, r.Original, r.Result)
		}
	}
	if changed == 0 {
		fmt.Fprintln(w, "  (no changes)")
	}
}
