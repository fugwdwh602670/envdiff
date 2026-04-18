package diff

import "fmt"

// Result holds the comparison outcome between two env files.
type Result struct {
	MissingInA  []string          // keys present in B but not in A
	MissingInB  []string          // keys present in A but not in B
	Mismatched  map[string][2]string // keys present in both but with different values
}

// String returns a human-readable summary of the diff result.
func (r Result) String() string {
	var out string
	for _, k := range r.MissingInA {
		out += fmt.Sprintf("- missing in A: %s\n", k)
	}
	for _, k := range r.MissingInB {
		out += fmt.Sprintf("- missing in B: %s\n", k)
	}
	for k, vals := range r.Mismatched {
		out += fmt.Sprintf("~ mismatch %s: %q vs %q\n", k, vals[0], vals[1])
	}
	return out
}

// HasDiff reports whether there are any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInA) > 0 || len(r.MissingInB) > 0 || len(r.Mismatched) > 0
}

// Compare takes two parsed env maps and returns a Result describing their differences.
func Compare(a, b map[string]string) Result {
	res := Result{
		Mismatched: make(map[string][2]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; !ok {
			res.MissingInB = append(res.MissingInB, k)
		} else if va != vb {
			res.Mismatched[k] = [2]string{va, vb}
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			res.MissingInA = append(res.MissingInA, k)
		}
	}

	return res
}
