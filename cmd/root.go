package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	outputFormat string
	strict       bool
)

var rootCmd = &cobra.Command{
	Use:   "envdiff <file1> <file2>",
	Short: "Compare .env files and surface missing or mismatched keys",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		fileA, fileB := args[0], args[1]

		envA, err := parser.ParseFile(fileA)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fileA, err)
		}

		envB, err := parser.ParseFile(fileB)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fileB, err)
		}

		results := diff.Compare(envA, envB)

		if len(results) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "No differences found.")
			return nil
		}

		err = diff.WriteReport(cmd.OutOrStdout(), results, fileA, fileB)
		if err != nil {
			return fmt.Errorf("writing report: %w", err)
		}

		if strict {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text")
	rootCmd.Flags().BoolVarP(&strict, "strict", "s", false, "Exit with code 1 if any differences are found")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
