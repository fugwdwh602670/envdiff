package diff

import (
	"strings"
)

// NormalizeOptions controls how env values are normalized before comparison.
type NormalizeOptions struct {
	TrimSpace    bool
	LowerKeys    bool
	LowerValues  bool
	StripQuotes  bool
}

// DefaultNormalizeOptions returns sensible defaults: trim whitespace and strip
// surrounding quotes, but leave key/value casing unchanged.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:   true,
		LowerKeys:   false,
		LowerValues: false,
		StripQuotes: true,
	}
}

// NormalizeEnv returns a new map with keys and values normalized according to
// the provided options. The original map is never mutated.
func NormalizeEnv(env map[string]string, opts NormalizeOptions) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.StripQuotes {
			v = stripSurroundingQuotes(v)
		}
		if opts.LowerValues {
			v = strings.ToLower(v)
		}
		if opts.LowerKeys {
			k = strings.ToLower(k)
		}
		out[k] = v
	}
	return out
}

func stripSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// WriteNormalizeReport writes a human-readable summary of normalization changes
// to the provided writer.
func WriteNormalizeReport(original, normalized map[string]string, w interface{ WriteString(string) (int, error) }) {
	changed := 0
	for k, orig := range original {
		if norm, ok := normalized[k]; ok && norm != orig {
			changed++
			w.WriteString("  ~ " + k + ": " + repr(orig) + " -> " + repr(norm) + "\n")
		}
	}
	if changed == 0 {
		w.WriteString("No normalization changes.\n")
	}
}

func repr(s string) string {
	return "\"" + s + "\""
}
