package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// PinEntry represents a pinned key-value pair with metadata.
type PinEntry struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	PinnedAt  time.Time `json:"pinned_at"`
	Comment   string    `json:"comment,omitempty"`
}

// PinFile holds all pinned entries.
type PinFile struct {
	Version int        `json:"version"`
	Pins    []PinEntry `json:"pins"`
}

// DefaultPinOptions returns default options for pin operations.
type PinOptions struct {
	Comment    string
	Overwrite  bool
	ReportOnly bool
}

func DefaultPinOptions() PinOptions {
	return PinOptions{
		Overwrite:  false,
		ReportOnly: false,
	}
}

// PinEnv pins the given keys from env into the pin file at path.
func PinEnv(env map[string]string, keys []string, path string, opts PinOptions) ([]PinEntry, error) {
	pf, err := LoadPinFile(path)
	if err != nil {
		return nil, err
	}

	pinned := make(map[string]*PinEntry, len(pf.Pins))
	for i := range pf.Pins {
		pinned[pf.Pins[i].Key] = &pf.Pins[i]
	}

	var added []PinEntry
	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			continue
		}
		if _, exists := pinned[k]; exists && !opts.Overwrite {
			continue
		}
		entry := PinEntry{Key: k, Value: v, PinnedAt: time.Now().UTC(), Comment: opts.Comment}
		pinned[k] = &entry
		added = append(added, entry)
	}

	if opts.ReportOnly {
		return added, nil
	}

	pf.Pins = make([]PinEntry, 0, len(pinned))
	for _, e := range pinned {
		pf.Pins = append(pf.Pins, *e)
	}
	sort.Slice(pf.Pins, func(i, j int) bool { return pf.Pins[i].Key < pf.Pins[j].Key })

	return added, savePinFile(path, pf)
}

// CheckPins verifies that pinned keys in the pin file match the given env.
func CheckPins(env map[string]string, path string) ([]PinViolation, error) {
	pf, err := LoadPinFile(path)
	if err != nil {
		return nil, err
	}
	var violations []PinViolation
	for _, pin := range pf.Pins {
		actual, ok := env[pin.Key]
		if !ok {
			violations = append(violations, PinViolation{Key: pin.Key, Expected: pin.Value, Actual: "", Missing: true})
		} else if actual != pin.Value {
			violations = append(violations, PinViolation{Key: pin.Key, Expected: pin.Value, Actual: actual})
		}
	}
	return violations, nil
}

// PinViolation describes a mismatch between a pinned value and the actual env.
type PinViolation struct {
	Key      string
	Expected string
	Actual   string
	Missing  bool
}

// LoadPinFile loads a pin file from disk, returning an empty one if not found.
func LoadPinFile(path string) (PinFile, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return PinFile{Version: 1}, nil
	}
	if err != nil {
		return PinFile{}, fmt.Errorf("pin: read %s: %w", path, err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return PinFile{}, fmt.Errorf("pin: parse %s: %w", path, err)
	}
	return pf, nil
}

func savePinFile(path string, pf PinFile) error {
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("pin: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
