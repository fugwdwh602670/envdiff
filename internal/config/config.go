package config

import (
	"encoding/json"
	"os"
)

// Config holds persistent envdiff configuration.
type Config struct {
	DefaultFormat string   `json:"default_format"` // text, json, csv
	IgnoreKeys    []string `json:"ignore_keys"`
	ShowMissing   bool     `json:"show_missing"`
	ShowMismatch  bool     `json:"show_mismatch"`
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		DefaultFormat: "text",
		IgnoreKeys:    []string{},
		ShowMissing:   true,
		ShowMismatch:  true,
	}
}

// Load reads a config file from path. If the file does not exist,
// Default() is returned without error.
func Load(path string) (Config, error) {
	cfg := Default()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Save writes cfg to path as JSON.
func Save(path string, cfg Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
