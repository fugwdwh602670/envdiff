package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	scopePrefixes []string
	scopeInvert   bool
)

var scopeCmd = &cobra.Command{
	Use:   "scope <file> [prefixes...]",
	Short: "Filter env file keys by prefix scope",
	Long: `Display only keys that match (or don't match) the given prefixes.

Examples:
  envdiff scope .env DB_ REDIS_
  envdiff scope .env APP_ --invert`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		if len(args) > 1 {
			scopePrefixes = append(scopePrefixes, args[1:]...)
		}

		env, err := parser.ParseFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		// Normalise prefixes: ensure they are upper-case for consistency.
		normalised := make([]string, len(scopePrefixes))
		for i, p := range scopePrefixes {
			normalised[i] = strings.ToUpper(p)
		}

		opts := diff.ScopeOptions{
			Prefixes: normalised,
			Invert:   scopeInvert,
		}

		scoped := diff.ScopeEnv(env, opts)
		summary := diff.SummarizeScope(env, scoped, opts)
		diff.WriteScopeReport(os.Stdout, scoped, summary)
		return nil
	},
}

func init() {
	scopeCmd.Flags().BoolVar(&scopeInvert, "invert", false, "Return keys that do NOT match the given prefixes")
	rootCmd.AddCommand(scopeCmd)
}
