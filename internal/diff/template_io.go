package diff

import "os"

// readFile is a thin wrapper so tests can swap it out if needed.
var readFile = func(path string) ([]byte, error) {
	return os.ReadFile(path)
}
