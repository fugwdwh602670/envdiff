package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// GraphOptions controls how the key dependency graph is rendered.
type GraphOptions struct {
	// ShowValues includes values in the graph output.
	ShowValues bool
	// PrefixDepth controls how many prefix segments define a "group node".
	PrefixDepth int
}

// DefaultGraphOptions returns sensible defaults.
func DefaultGraphOptions() GraphOptions {
	return GraphOptions{
		ShowValues:  false,
		PrefixDepth: 1,
	}
}

// GraphNode represents a node in the key graph.
type GraphNode struct {
	Label    string
	Children []*GraphNode
	Keys     []string
}

// BuildGraph groups env keys into a prefix-based tree structure.
func BuildGraph(env map[string]string, opts GraphOptions) *GraphNode {
	root := &GraphNode{Label: "(root)"}
	groups := map[string]*GraphNode{}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		prefix := extractPrefix(k, opts.PrefixDepth)
		if prefix == "" {
			root.Keys = append(root.Keys, k)
			continue
		}
		if _, ok := groups[prefix]; !ok {
			node := &GraphNode{Label: prefix}
			groups[prefix] = node
			root.Children = append(root.Children, node)
		}
		groups[prefix].Keys = append(groups[prefix].Keys, k)
	}

	sort.Slice(root.Children, func(i, j int) bool {
		return root.Children[i].Label < root.Children[j].Label
	})
	return root
}

func extractPrefix(key string, depth int) string {
	parts := strings.SplitN(key, "_", depth+1)
	if len(parts) <= depth {
		return ""
	}
	return strings.Join(parts[:depth], "_")
}

// WriteGraphReport writes an ASCII tree of the graph to w.
func WriteGraphReport(w io.Writer, root *GraphNode, env map[string]string, opts GraphOptions) {
	fmt.Fprintf(w, "%s\n", root.Label)
	for _, k := range root.Keys {
		writeKeyLine(w, "  ", k, env, opts.ShowValues)
	}
	for i, child := range root.Children {
		connector := "├──"
		if i == len(root.Children)-1 {
			connector = "└──"
		}
		fmt.Fprintf(w, "%s %s/\n", connector, child.Label)
		for _, k := range child.Keys {
			writeKeyLine(w, "    ", k, env, opts.ShowValues)
		}
	}
}

func writeKeyLine(w io.Writer, indent, key string, env map[string]string, showValues bool) {
	if showValues {
		fmt.Fprintf(w, "%s- %s = %s\n", indent, key, env[key])
	} else {
		fmt.Fprintf(w, "%s- %s\n", indent, key)
	}
}
