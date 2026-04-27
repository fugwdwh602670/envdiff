package diff

// Status represents the comparison result for a single key.
type Status string

const (
	StatusMatch      Status = "match"
	StatusMissingInA Status = "missing_in_a"
	StatusMissingInB Status = "missing_in_b"
	StatusMismatch   Status = "mismatch"
)

// Result holds the comparison outcome for a single environment key.
type Result struct {
	Key      string
	ValueA   string
	ValueB   string
	Status   Status
}

// Compare compares two env maps and returns a slice of Results.
func Compare(a, b map[string]string) []Result {
	seen := make(map[string]bool)
	var results []Result

	for k, va := range a {
		seen[k] = true
		if vb, ok := b[k]; !ok {
			results = append(results, Result{Key: k, ValueA: va, ValueB: "", Status: StatusMissingInB})
		} else if va != vb {
			results = append(results, Result{Key: k, ValueA: va, ValueB: vb, Status: StatusMismatch})
		} else {
			results = append(results, Result{Key: k, ValueA: va, ValueB: vb, Status: StatusMatch})
		}
	}

	for k, vb := range b {
		if !seen[k] {
			results = append(results, Result{Key: k, ValueA: "", ValueB: vb, Status: StatusMissingInA})
		}
	}

	return results
}
