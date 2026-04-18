package diff

import (
	"bufio"
	"os"
	"strings"
)

// IgnoreList holds a set of keys to exclude from diff results.
type IgnoreList struct {
	keys map[string]struct{}
}

// NewIgnoreList creates an empty IgnoreList.
func NewIgnoreList() *IgnoreList {
	return &IgnoreList{keys: make(map[string]struct{})}
}

// Add adds a key to the ignore list.
func (il *IgnoreList) Add(key string) {
	il.keys[strings.TrimSpace(key)] = struct{}{}
}

// Contains returns true if the key should be ignored.
func (il *IgnoreList) Contains(key string) bool {
	_, ok := il.keys[key]
	return ok
}

// Len returns the number of ignored keys.
func (il *IgnoreList) Len() int {
	return len(il.keys)
}

// LoadIgnoreFile reads a file with one key per line (# comments supported).
func LoadIgnoreFile(path string) (*IgnoreList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	il := NewIgnoreList()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		il.Add(line)
	}
	return il, scanner.Err()
}

// ApplyIgnoreList removes results whose keys appear in the ignore list.
func ApplyIgnoreList(results []Result, il *IgnoreList) []Result {
	if il == nil || il.Len() == 0 {
		return results
	}
	filtered := make([]Result, 0, len(results))
	for _, r := range results {
		if !il.Contains(r.Key) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
