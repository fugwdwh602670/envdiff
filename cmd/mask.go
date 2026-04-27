package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	maskPatterns     []string
	maskChar         string
	maskVisibleChars int
	maskShowAll      bool
)

var maskCmd = &cobra.Command{
	Use:   "mask <file>",
	Short: "Mask sensitive values in a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := parser.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		opts := diff.DefaultMaskOptions()
		if len(maskPatterns) > 0 {
			opts.Patterns = maskPatterns
		}
		if maskChar != "" {
			opts.MaskChar = maskChar
		}
		if maskVisibleChars >= 0 {
			opts.VisibleChars = maskVisibleChars
		}

		results := diff.MaskEnv(env, opts)

		// Sort for deterministic output
		sort.Slice(results, func(i, j int) bool {
			return results[i].Key < results[j].Key
		})

		if maskShowAll {
			diff.WriteMaskReport(cmd.OutOrStdout(), results)
			return nil
		}

		// Default: print masked key=value pairs
		for _, r := range results {
			val := r.Masked
			if strings.ContainsAny(val, " \t") {
				val = `"` + val + `"`
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", r.Key, val)
		}
		return nil
	},
}

func init() {
	maskCmd.Flags().StringArrayVar(&maskPatterns, "pattern", nil, "regex patterns to match sensitive key names (repeatable)")
	maskCmd.Flags().StringVar(&maskChar, "char", "", "character used for masking (default: *)")
	maskCmd.Flags().IntVar(&maskVisibleChars, "visible", 0, "number of trailing characters to leave visible")
	maskCmd.Flags().BoolVar(&maskShowAll, "report", false, "show full mask report instead of env output")
	rootCmd.AddCommand(maskCmd)

	_ = os.Stderr
}
