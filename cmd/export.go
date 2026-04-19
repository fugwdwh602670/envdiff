package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var format string
	var onlyMissing bool
	var output string

	cmd := &cobra.Command{
		Use:   "export <fileA> <fileB>",
		Short: "Export diff results to a file in env, json, or markdown format",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			envA, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[0], err)
			}
			envB, err := parser.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[1], err)
			}

			results := diff.Compare(envA, envB)

			opts := diff.DefaultExportOptions()
			opts.Format = diff.ExportFormat(format)
			opts.OnlyMissing = onlyMissing

			if output == "" {
				output = "envdiff-export." + format
			}

			if err := diff.ExportResults(results, output, opts); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Exported %d results to %s\n", len(results), output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "env", "Output format: env, json, markdown")
	cmd.Flags().BoolVar(&onlyMissing, "only-missing", false, "Export only missing keys")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path (default: envdiff-export.<format>)")

	rootCmd.AddCommand(cmd)
}
