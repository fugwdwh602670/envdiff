package config

import (
	"os"
	"path/filepath"
)

// ResolveConfigPath returns the config file path to use.
// Priority: ENVDIFF_CONFIG env var > $HOME/.envdiff.json
func ResolveConfigPath() string {
	if p := os.Getenv("ENVDIFF_CONFIG"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envdiff.json"
	}
	return filepath.Join(home, ".envdiff.json")
}

// LoadResolved loads config from the resolved path.
func LoadResolved() (Config, error) {
	return Load(ResolveConfigPath())
}

// IgnoreSet returns a set of keys to ignore from the config.
func (c Config) IgnoreSet() map[string]struct{} {
	s := make(map[string]struct{}, len(c.IgnoreKeys))
	for _, k := range c.IgnoreKeys {
		s[k] = struct{}{}
	}
	return s
}
