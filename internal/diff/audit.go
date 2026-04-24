package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// AuditEntry records a single key-level change event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Key       string    `json:"key"`
	Event     string    `json:"event"` // added, removed, changed, unchanged
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
}

// AuditOptions controls what gets included in the audit log.
type AuditOptions struct {
	IncludeUnchanged bool
	RedactValues     bool
}

// DefaultAuditOptions returns sensible defaults.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		IncludeUnchanged: false,
		RedactValues:     false,
	}
}

// AuditResults converts diff results into a timestamped audit log.
func AuditResults(results []Result, opts AuditOptions) []AuditEntry {
	now := time.Now().UTC()
	entries := make([]AuditEntry, 0, len(results))

	for _, r := range results {
		if r.Status == StatusMatch && !opts.IncludeUnchanged {
			continue
		}

		entry := AuditEntry{
			Timestamp: now,
			Key:       r.Key,
			Event:     auditEvent(r.Status),
		}

		if !opts.RedactValues {
			entry.OldValue = r.ValueA
			entry.NewValue = r.ValueB
		}

		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// WriteAuditReport writes a human-readable audit report to w.
func WriteAuditReport(w io.Writer, entries []AuditEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No audit events recorded.")
		return
	}
	fmt.Fprintf(w, "Audit Log (%d events)\n", len(entries))
	fmt.Fprintln(w, "---")
	for _, e := range entries {
		ts := e.Timestamp.Format(time.RFC3339)
		fmt.Fprintf(w, "[%s] %-12s %s", ts, e.Event, e.Key)
		if e.OldValue != "" || e.NewValue != "" {
			fmt.Fprintf(w, " (%q -> %q)", e.OldValue, e.NewValue)
		}
		fmt.Fprintln(w)
	}
}

func auditEvent(status Status) string {
	switch status {
	case StatusMissingInB:
		return "removed"
	case StatusMissingInA:
		return "added"
	case StatusMismatch:
		return "changed"
	default:
		return "unchanged"
	}
}
