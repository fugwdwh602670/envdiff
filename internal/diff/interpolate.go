package diff

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// DefaultInterpolateOptions returns a sensible default configuration.
func DefaultInterpolateOptions() InterpolateOptions {
	return InterpolateOptions{
		AllowMissing: false,
		Prefix:       "${",
		Suffix:       "}",
	}
}

// InterpolateOptions controls how variable interpolation is performed.
type InterpolateOptions struct {
	// AllowMissing suppresses errors for unresolved references.
	AllowMissing bool
	// Prefix is the opening delimiter for variable references (default "${").
	Prefix string
	// Suffix is the closing delimiter for variable references (default "}").
	Suffix string
}

// InterpolateResult holds the result of interpolating a single key.
type InterpolateResult struct {
	Key      string
	Original string
	Resolved string
	Refs     []string
	Missing  []string
}

// InterpolateEnv resolves variable references within env values using other
// keys in the same map. References use the form ${KEY}.
func InterpolateEnv(env map[string]string, opts InterpolateOptions) ([]InterpolateResult, error) {
	escapedPrefix := regexp.QuoteMeta(opts.Prefix)
	escapedSuffix := regexp.QuoteMeta(opts.Suffix)
	pattern := regexp.MustCompile(escapedPrefix + `([A-Z_][A-Z0-9_]*)` + escapedSuffix)

	var results []InterpolateResult
	for _, key := range sortedKeys(env) {
		original := env[key]
		matches := pattern.FindAllStringSubmatch(original, -1)
		if len(matches) == 0 {
			continue
		}
		var refs, missing []string
		resolved := original
		for _, m := range matches {
			refKey := m[1]
			refs = append(refs, refKey)
			if val, ok := env[refKey]; ok {
				resolved = strings.ReplaceAll(resolved, m[0], val)
			} else {
				missing = append(missing, refKey)
				if !opts.AllowMissing {
					return nil, fmt.Errorf("interpolate: key %q references undefined variable %q", key, refKey)
				}
			}
		}
		results = append(results, InterpolateResult{
			Key:      key,
			Original: original,
			Resolved: resolved,
			Refs:     refs,
			Missing:  missing,
		})
	}
	return results, nil
}

// WriteInterpolateReport writes a human-readable interpolation report to w.
func WriteInterpolateReport(w io.Writer, results []InterpolateResult) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No interpolated variables found.")
		return
	}
	for _, r := range results {
		fmt.Fprintf(w, "[%s]\n", r.Key)
		fmt.Fprintf(w, "  original: %s\n", r.Original)
		fmt.Fprintf(w, "  resolved: %s\n", r.Resolved)
		if len(r.Missing) > 0 {
			fmt.Fprintf(w, "  missing refs: %s\n", strings.Join(r.Missing, ", "))
		}
	}
}
