package diff

import (
	"os"
	"sync"
	"testing"
	"time"
)

func writeTempEnvWatch(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envwatch*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestWatch_DetectsChange(t *testing.T) {
	fileA := writeTempEnvWatch(t, "KEY=foo\n")
	fileB := writeTempEnvWatch(t, "KEY=foo\n")

	var mu sync.Mutex
	var calls [][]Result

	opts := WatchOptions{
		Interval: 50 * time.Millisecond,
		OnChange: func(r []Result) {
			mu.Lock()
			calls = append(calls, r)
			mu.Unlock()
		},
	}

	stop := make(chan struct{})
	done := make(chan error, 1)
	go func() { done <- Watch(fileA, fileB, opts, stop) }()

	// Wait for initial fire
	time.Sleep(120 * time.Millisecond)

	// Mutate fileB to introduce a mismatch
	os.WriteFile(fileB, []byte("KEY=bar\n"), 0644)
	time.Sleep(120 * time.Millisecond)

	close(stop)
	if err := <-done; err != nil {
		t.Fatalf("Watch returned error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(calls) < 2 {
		t.Fatalf("expected at least 2 onChange calls, got %d", len(calls))
	}
}

func TestWatch_NoChangeCallback(t *testing.T) {
	fileA := writeTempEnvWatch(t, "KEY=foo\n")
	fileB := writeTempEnvWatch(t, "KEY=foo\n")
	opts := WatchOptions{Interval: 50 * time.Millisecond}
	stop := make(chan struct{})
	close(stop)
	err := Watch(fileA, fileB, opts, stop)
	if err == nil {
		t.Fatal("expected error when OnChange is nil")
	}
}

func TestFileHash_Consistency(t *testing.T) {
	f := writeTempEnvWatch(t, "A=1\nB=2\n")
	h1, err := fileHash(f)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := fileHash(f)
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("expected consistent hash, got %s and %s", h1, h2)
	}
}
