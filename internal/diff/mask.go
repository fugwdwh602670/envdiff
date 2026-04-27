package diff

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// Patterns is a list of key name patterns (regex) to mask.
	Patterns []string
	// MaskChar is the character used to replace masked value characters.
	MaskChar string
	// VisibleChars is the number of trailing characters to leave visible (0 = mask all).
	VisibleChars int
}

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Patterns:     []string{`(?i)(secret|password|token|key|api_key|private)`},
		MaskChar:     "*",
		VisibleChars: 0,
	}
}

// MaskResult holds a key, its original value, and the masked value.
type MaskResult struct {
	Key     string
	Original string
	Masked  string
	WasMasked bool
}

// MaskEnv applies masking rules to the given env map and returns MaskResults.
func MaskEnv(env map[string]string, opts MaskOptions) []MaskResult {
	if opts.MaskChar == "" {
		opts.MaskChar = "*"
	}
	compiledPatterns := make([]*regexp.Regexp, 0, len(opts.Patterns))
	for _, p := range opts.Patterns {
		if re, err := regexp.Compile(p); err == nil {
			compiledPatterns = append(compiledPatterns, re)
		}
	}

	results := make([]MaskResult, 0, len(env))
	for k, v := range env {
		masked := false
		for _, re := range compiledPatterns {
			if re.MatchString(k) {
				masked = true
				break
			}
		}
		maskedVal := v
		if masked {
			maskedVal = maskValue(v, opts.MaskChar, opts.VisibleChars)
		}
		results = append(results, MaskResult{
			Key:       k,
			Original:  v,
			Masked:    maskedVal,
			WasMasked: masked,
		})
	}
	return results
}

func maskValue(val, char string, visible int) string {
	if len(val) == 0 {
		return val
	}
	if visible <= 0 || visible >= len(val) {
		return strings.Repeat(char, len(val))
	}
	return strings.Repeat(char, len(val)-visible) + val[len(val)-visible:]
}

// WriteMaskReport writes a human-readable mask report to w.
func WriteMaskReport(w io.Writer, results []MaskResult) {
	maskedCount := 0
	for _, r := range results {
		if r.WasMasked {
			maskedCount++
		}
	}
	fmt.Fprintf(w, "Mask report: %d key(s) masked\n\n", maskedCount)
	for _, r := range results {
		if r.WasMasked {
			fmt.Fprintf(w, "  [MASKED] %s = %s\n", r.Key, r.Masked)
		} else {
			fmt.Fprintf(w, "  [plain]  %s = %s\n", r.Key, r.Original)
		}
	}
}
