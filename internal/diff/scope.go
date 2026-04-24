package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ScopeOptions controls how scoping is applied.
type ScopeOptions struct {
	// Prefixes is the list of key prefixes to include (e.g. ["DB_", "REDIS_"]).
	Prefixes []string
	// Invert returns keys that do NOT match any prefix.
	Invert bool
}

// DefaultScopeOptions returns ScopeOptions with no filtering.
func DefaultScopeOptions() ScopeOptions {
	return ScopeOptions{}
}

// ScopeEnv filters an env map to only keys matching (or not matching) the given prefixes.
func ScopeEnv(env map[string]string, opts ScopeOptions) map[string]string {
	if len(opts.Prefixes) == 0 {
		result := make(map[string]string, len(env))
		for k, v := range env {
			result[k] = v
		}
		return result
	}

	result := make(map[string]string)
	for k, v := range env {
		matched := matchesAnyPrefix(k, opts.Prefixes)
		if (!opts.Invert && matched) || (opts.Invert && !matched) {
			result[k] = v
		}
	}
	return result
}

func matchesAnyPrefix(key string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// ScopeSummary holds the result of a scope operation.
type ScopeSummary struct {
	Total    int
	Included int
	Excluded int
	Prefixes []string
	Inverted bool
}

// SummarizeScope returns a summary of what was kept/dropped.
func SummarizeScope(original, scoped map[string]string, opts ScopeOptions) ScopeSummary {
	return ScopeSummary{
		Total:    len(original),
		Included: len(scoped),
		Excluded: len(original) - len(scoped),
		Prefixes: opts.Prefixes,
		Inverted: opts.Invert,
	}
}

// WriteScopeReport writes a human-readable scope report to w.
func WriteScopeReport(w io.Writer, env map[string]string, summary ScopeSummary) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	if len(summary.Prefixes) > 0 {
		mode := "include"
		if summary.Inverted {
			mode = "exclude"
		}
		fmt.Fprintf(w, "Scope (%s prefixes: %s)\n", mode, strings.Join(summary.Prefixes, ", "))
	} else {
		fmt.Fprintln(w, "Scope (all keys)")
	}
	fmt.Fprintf(w, "  Total: %d | Included: %d | Excluded: %d\n\n", summary.Total, summary.Included, summary.Excluded)

	for _, k := range keys {
		fmt.Fprintf(w, "  %s=%s\n", k, env[k])
	}
	if len(keys) == 0 {
		fmt.Fprintln(w, "  (no keys matched)")
	}
}
