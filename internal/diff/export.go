package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ExportFormat string

const (
	ExportFormatEnv      ExportFormat = "env"
	ExportFormatJSON     ExportFormat = "json"
	ExportFormatMarkdown ExportFormat = "markdown"
)

type ExportOptions struct {
	Format ExportFormat
	OnlyMissing bool
}

func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:      ExportFormatEnv,
		OnlyMissing: false,
	}
}

// ExportResults writes diff results to a file in the specified format.
func ExportResults(results []Result, path string, opts ExportOptions) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("export: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create: %w", err)
	}
	defer f.Close()

	filtered := results
	if opts.OnlyMissing {
		filtered = []Result{}
		for _, r := range results {
			if r.Status == StatusMissingInB || r.Status == StatusMissingInA {
				filtered = append(filtered, r)
			}
		}
	}

	switch opts.Format {
	case ExportFormatJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(filtered)
	case ExportFormatMarkdown:
		return exportMarkdown(f, filtered)
	default:
		return exportEnv(f, filtered)
	}
}

func exportEnv(f *os.File, results []Result) error {
	for _, r := range results {
		fmt.Fprintf(f, "# status: %s\n%s=%s\n", r.Status, r.Key, r.ValueA)
	}
	return nil
}

func exportMarkdown(f *os.File, results []Result) error {
	fmt.Fprintln(f, "| Key | Status | Value A | Value B |")
	fmt.Fprintln(f, "|-----|--------|---------|---------|")
	for _, r := range results {
		fmt.Fprintf(f, "| %s | %s | %s | %s |\n",
			r.Key,
			strings.ToLower(string(r.Status)),
			r.ValueA,
			r.ValueB,
		)
	}
	return nil
}
