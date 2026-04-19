package diff

import (
	"encoding/json"
	"os"
	"time"
)

// Baseline represents a saved snapshot of diff results for later comparison.
type Baseline struct {
	CreatedAt time.Time `json:"created_at"`
	FileA     string    `json:"file_a"`
	FileB     string    `json:"file_b"`
	Results   []Result  `json:"results"`
}

// SaveBaseline writes a baseline snapshot to the given path.
func SaveBaseline(path, fileA, fileB string, results []Result) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		FileA:     fileA,
		FileB:     fileB,
		Results:   results,
	}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadBaseline reads a baseline snapshot from the given path.
func LoadBaseline(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// DiffBaseline compares current results against a saved baseline.
// Returns new issues (in current but not baseline) and resolved issues (in baseline but not current).
func DiffBaseline(baseline []Result, current []Result) (newIssues []Result, resolved []Result) {
	baselineMap := make(map[string]Result, len(baseline))
	for _, r := range baseline {
		baselineMap[r.Key] = r
	}
	currentMap := make(map[string]Result, len(current))
	for _, r := range current {
		currentMap[r.Key] = r
	}
	for _, r := range current {
		if _, found := baselineMap[r.Key]; !found {
			newIssues = append(newIssues, r)
		}
	}
	for _, r := range baseline {
		if _, found := currentMap[r.Key]; !found {
			resolved = append(resolved, r)
		}
	}
	return
}
