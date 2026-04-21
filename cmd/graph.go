package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var showValues bool
	var prefixDepth int

	graphCmd := &cobra.Command{
		Use:   "graph <file>",
		Short: "Render an ASCII tree grouping keys by prefix",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			opts := diff.DefaultGraphOptions()
			opts.ShowValues = showValues
			if prefixDepth > 0 {
				opts.PrefixDepth = prefixDepth
			}

			root := diff.BuildGraph(env, opts)
			diff.WriteGraphReport(os.Stdout, root, env, opts)
			return nil
		},
	}

	graphCmd.Flags().BoolVar(&showValues, "values", false, "Show key values in the graph")
	graphCmd.Flags().IntVar(&prefixDepth, "depth", 1, "Number of underscore-delimited segments to use as group prefix")

	rootCmd.AddCommand(graphCmd)
}
