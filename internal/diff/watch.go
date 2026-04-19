package diff

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

// WatchOptions configures the watch behavior.
type WatchOptions struct {
	Interval time.Duration
	OnChange func(results []Result)
}

// DefaultWatchOptions returns sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Interval: 2 * time.Second,
	}
}

// fileHash returns the md5 hash of a file's contents.
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// Watch polls two env files and calls opts.OnChange when their diff changes.
// It blocks until stop is closed.
func Watch(fileA, fileB string, opts WatchOptions, stop <-chan struct{}) error {
	if opts.OnChange == nil {
		return fmt.Errorf("watch: OnChange callback must be set")
	}
	lastHash := ""
	tick := time.NewTicker(opts.Interval)
	defer tick.Stop()
	for {
		select {
		case <-stop:
			return nil
		case <-tick.C:
			ha, err1 := fileHash(fileA)
			hb, err2 := fileHash(fileB)
			if err1 != nil || err2 != nil {
				continue
			}
			combined := ha + hb
			if combined == lastHash {
				continue
			}
			lastHash = combined
			results, err := compareFiles(fileA, fileB)
			if err != nil {
				continue
			}
			opts.OnChange(results)
		}
	}
}

// compareFiles is a helper that parses and compares two env files.
func compareFiles(fileA, fileB string) ([]Result, error) {
	import_parser_inline := func(path string) (map[string]string, error) {
		// Delegate to the parser package via a thin re-export.
		return parseEnvFile(path)
	}
	a, err := import_parser_inline(fileA)
	if err != nil {
		return nil, err
	}
	b, err := import_parser_inline(fileB)
	if err != nil {
		return nil, err
	}
	return Compare(a, b), nil
}
