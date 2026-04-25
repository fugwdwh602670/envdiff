package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// AliasMap maps canonical key names to one or more alternate names.
type AliasMap map[string][]string

// DefaultAliasOptions returns a no-op alias map.
func DefaultAliasOptions() AliasMap {
	return AliasMap{}
}

// ResolveAliases rewrites keys in env according to the alias map.
// If a key matches an alias for a canonical key, it is renamed to the canonical key.
// The original key is removed. If the canonical key already exists, it is not overwritten.
func ResolveAliases(env map[string]string, aliases AliasMap) map[string]string {
	// Build reverse map: alias -> canonical
	reverse := make(map[string]string)
	for canonical, alts := range aliases {
		for _, alt := range alts {
			reverse[strings.ToUpper(alt)] = canonical
		}
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if canonical, ok := reverse[strings.ToUpper(k)]; ok {
			if _, exists := result[canonical]; !exists {
				result[canonical] = v
			}
		} else {
			result[k] = v
		}
	}
	return result
}

// AliasReport describes a single alias resolution event.
type AliasReport struct {
	Alias     string
	Canonical string
	Value     string
}

// AuditAliases returns a list of keys that were resolved from aliases.
func AuditAliases(env map[string]string, aliases AliasMap) []AliasReport {
	reverse := make(map[string]string)
	for canonical, alts := range aliases {
		for _, alt := range alts {
			reverse[strings.ToUpper(alt)] = canonical
		}
	}

	var reports []AliasReport
	for k, v := range env {
		if canonical, ok := reverse[strings.ToUpper(k)]; ok {
			reports = append(reports, AliasReport{Alias: k, Canonical: canonical, Value: v})
		}
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Alias < reports[j].Alias
	})
	return reports
}

// WriteAliasReport writes a human-readable alias resolution report to w.
func WriteAliasReport(w io.Writer, reports []AliasReport) {
	if len(reports) == 0 {
		fmt.Fprintln(w, "No alias resolutions applied.")
		return
	}
	fmt.Fprintf(w, "Alias resolutions (%d):\n", len(reports))
	for _, r := range reports {
		fmt.Fprintf(w, "  %s -> %s\n", r.Alias, r.Canonical)
	}
}
