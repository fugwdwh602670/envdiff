package diff

// FilterOptions controls which diff results are included in output.
type FilterOptions struct {
	ShowMissingInA bool
	ShowMissingInB bool
	ShowMismatch   bool
}

// DefaultFilterOptions returns a FilterOptions that includes all result types.
func DefaultFilterOptions() FilterOptions {
	return FilterOptions{
		ShowMissingInA: true,
		ShowMissingInB: true,
		ShowMismatch:   true,
	}
}

// Filter returns a subset of results based on the provided FilterOptions.
func Filter(results []Result, opts FilterOptions) []Result {
	filtered := make([]Result, 0, len(results))
	for _, r := range results {
		switch r.Status {
		case MissingInA:
			if opts.ShowMissingInA {
				filtered = append(filtered, r)
			}
		case MissingInB:
			if opts.ShowMissingInB {
				filtered = append(filtered, r)
			}
		case Mismatch:
			if opts.ShowMismatch {
				filtered = append(filtered, r)
			}
		}
	}
	return filtered
}
