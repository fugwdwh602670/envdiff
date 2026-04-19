package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var annotations []string

	cmd := &cobra.Command{
		Use:   "annotate <fileA> <fileB>",
		Short: "Annotate diff results with comments for specific keys",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			envA, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}
			envB, err := parser.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[1], err)
			}

			results := diff.Compare(envA, envB)

			opts := diff.DefaultAnnotateOptions()
			for _, a := range annotations {
				parts := strings.SplitN(a, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid annotation %q: use KEY=comment", a)
				}
				opts.Annotations[parts[0]] = parts[1]
			}

			anns := diff.AnnotateResults(results, opts)
			diff.WriteAnnotationReport(os.Stdout, anns)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&annotations, "annotate", "a", nil, "KEY=comment pairs to annotate")
	rootCmd.AddCommand(cmd)
}
