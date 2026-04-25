package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an env file's parsed state.
type Snapshot struct {
	CreatedAt time.Time         `json:"created_at"`
	Label     string            `json:"label"`
	Env       map[string]string `json:"env"`
}

// DefaultSnapshotOptions returns a Snapshot with the current time and optional label.
func DefaultSnapshotOptions(label string) Snapshot {
	return Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Env:       make(map[string]string),
	}
}

// SaveSnapshot writes a snapshot of the given env map to the specified file path.
func SaveSnapshot(path string, env map[string]string, label string) error {
	snap := DefaultSnapshotOptions(label)
	for k, v := range env {
		snap.Env[k] = v
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadSnapshot reads a snapshot from disk.
func LoadSnapshot(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read failed: %w", err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: invalid JSON: %w", err)
	}
	return snap, nil
}

// DiffSnapshot compares a live env map against a stored snapshot and returns diff results.
func DiffSnapshot(snap Snapshot, live map[string]string) []Result {
	return Compare(snap.Env, live)
}
