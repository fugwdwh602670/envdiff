package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Lint a .env file for common issues",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]

		env, err := parser.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		opts := diff.DefaultLintOptions()
		issues := diff.LintEnv(env, opts)
		diff.WriteLintReport(cmd.OutOrStdout(), issues, filePath)

		if len(issues) > 0 {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
