package diff

import (
	"fmt"
	"io"
	"sort"
)

// Annotation holds a comment or note attached to a key.
type Annotation struct {
	Key     string `json:"key"`
	Comment string `json:"comment"`
}

// AnnotateOptions controls annotation behaviour.
type AnnotateOptions struct {
	Annotations map[string]string // key -> comment
}

// DefaultAnnotateOptions returns options with an empty annotation map.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Annotations: make(map[string]string),
	}
}

// AnnotateResults attaches annotations to results, returning a slice of
// Annotation for keys that have an entry in opts.Annotations.
func AnnotateResults(results []Result, opts AnnotateOptions) []Annotation {
	var out []Annotation
	for _, r := range results {
		if comment, ok := opts.Annotations[r.Key]; ok {
			out = append(out, Annotation{Key: r.Key, Comment: comment})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out
}

// WriteAnnotationReport writes a human-readable annotation report to w.
func WriteAnnotationReport(w io.Writer, annotations []Annotation) {
	if len(annotations) == 0 {
		fmt.Fprintln(w, "No annotations.")
		return
	}
	fmt.Fprintln(w, "Annotations:")
	for _, a := range annotations {
		fmt.Fprintf(w, "  %-30s %s\n", a.Key, a.Comment)
	}
}
