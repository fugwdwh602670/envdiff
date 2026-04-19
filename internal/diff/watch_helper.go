package diff

import (
	"github.com/user/envdiff/internal/parser"
)

// parseEnvFile is an internal helper used by Watch to parse an env file.
func parseEnvFile(path string) (map[string]string, error) {
	return parser.ParseFile(path)
}
