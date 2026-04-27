package diff

import (
	"fmt"
	"io"
)

// WritePinReport writes a human-readable report of pin operations or violations.
func WritePinReport(w io.Writer, added []PinEntry, violations []PinViolation) {
	if len(added) > 0 {
		fmt.Fprintf(w, "Pinned %d key(s):\n", len(added))
		for _, e := range added {
			if e.Comment != "" {
				fmt.Fprintf(w, "  + %s = %s  # %s\n", e.Key, e.Value, e.Comment)
			} else {
				fmt.Fprintf(w, "  + %s = %s\n", e.Key, e.Value)
			}
		}
	}

	if len(violations) == 0 {
		if len(added) == 0 {
			fmt.Fprintln(w, "No pin violations found.")
		}
		return
	}

	fmt.Fprintf(w, "\n%d pin violation(s):\n", len(violations))
	for _, v := range violations {
		if v.Missing {
			fmt.Fprintf(w, "  MISSING  %s (expected: %s)\n", v.Key, v.Expected)
		} else {
			fmt.Fprintf(w, "  MISMATCH %s\n    expected: %s\n    actual:   %s\n", v.Key, v.Expected, v.Actual)
		}
	}
}
