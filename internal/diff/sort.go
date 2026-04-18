package diff

import "sort"

// SortOrder defines how results are sorted.
type SortOrder string

const (
	SortByKey    SortOrder = "key"
	SortByStatus SortOrder = "status"
)

// SortResults sorts a slice of Result by the given order.
// Secondary sort is always by key for stability.
func SortResults(results []Result, order SortOrder) []Result {
	out := make([]Result, len(results))
	copy(out, results)

	switch order {
	case SortByStatus:
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].Status != out[j].Status {
				return statusRank(out[i].Status) < statusRank(out[j].Status)
			}
			return out[i].Key < out[j].Key
		})
	default: // SortByKey
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}

func statusRank(s Status) int {
	switch s {
	case StatusMissingInB:
		return 0
	case StatusMissingInA:
		return 1
	case StatusMismatch:
		return 2
	case StatusMatch:
		return 3
	default:
		return 4
	}
}
