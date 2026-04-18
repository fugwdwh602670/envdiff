package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Env holds the key-value pairs parsed from a .env file.
type Env map[string]string

// ParseFile reads a .env file and returns an Env map.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or a line is malformed.
func ParseFile(path string) (Env, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(Env)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional leading "export "
		line = strings.TrimPrefix(line, "export ")

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			return nil, fmt.Errorf("parser: %q line %d: missing '=' in %q", path, lineNum, line)
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// Strip surrounding quotes from value.
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if key == "" {
			return nil, fmt.Errorf("parser: %q line %d: empty key", path, lineNum)
		}

		env[key] = val
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scan %q: %w", path, err)
	}

	return env, nil
}
