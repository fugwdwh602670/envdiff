package diff

// Summary holds aggregate counts from a diff result set.
type Summary struct {
	Total        int
	MissingInB   int
	MissingInA   int
	Mismatch     int
	Match        int
}

// Summarize computes a Summary from a slice of Result.
func Summarize(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusMissingInB:
			s.MissingInB++
		case StatusMissingInA:
			s.MissingInA++
		case StatusMismatch:
			s.Mismatch++
		case StatusMatch:
			s.Match++
		}
	}
	return s
}

// HasDiff returns true if any non-match results exist.
func (s Summary) HasDiff() bool {
	return s.MissingInA > 0 || s.MissingInB > 0 || s.Mismatch > 0
}
