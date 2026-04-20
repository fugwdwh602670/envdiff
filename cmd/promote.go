package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var overwrite bool
	var dryRun bool
	var output string

	promoteCmd := &cobra.Command{
		Use:   "promote <source.env> <target.env>",
		Short: "Promote missing (or all) keys from source into target env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			srcPath, tgtPath := args[0], args[1]

			src, err := parser.ParseFile(srcPath)
			if err != nil {
				return fmt.Errorf("parsing source: %w", err)
			}
			tgt, err := parser.ParseFile(tgtPath)
			if err != nil {
				return fmt.Errorf("parsing target: %w", err)
			}

			opts := diff.DefaultPromoteOptions()
			opts.OverwriteExisting = overwrite
			opts.OnlyMissing = !overwrite
			opts.DryRun = dryRun

			merged, results := diff.Promote(src, tgt, opts)

			diff.WritePromoteReport(cmd.OutOrStdout(), results)

			if dryRun {
				return nil
			}

			dest := tgtPath
			if output != "" {
				dest = output
			}
			if err := diff.WriteMerged(merged, dest); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Written to %s\n", dest)
			return nil
		},
	}

	promoteCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in target with source values")
	promoteCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing to disk")
	promoteCmd.Flags().StringVarP(&output, "output", "o", "", "Write result to this file instead of target (implies no in-place edit)")

	_ = os.Getenv // suppress unused import if needed
	rootCmd.AddCommand(promoteCmd)
}
