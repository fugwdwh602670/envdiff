package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	interpolateAllowMissing bool
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <file>",
	Short: "Resolve variable references within a .env file",
	Long: `Scans a .env file for values containing variable references in the form
${KEY} and resolves them using other keys defined in the same file.

By default an error is returned for any unresolved reference. Use
--allow-missing to suppress errors and emit a warning instead.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := parser.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("interpolate: failed to parse %q: %w", args[0], err)
		}

		opts := diff.DefaultInterpolateOptions()
		opts.AllowMissing = interpolateAllowMissing

		results, err := diff.InterpolateEnv(env, opts)
		if err != nil {
			return err
		}

		diff.WriteInterpolateReport(os.Stdout, results)
		return nil
	},
}

func init() {
	interpolateCmd.Flags().BoolVar(
		&interpolateAllowMissing,
		"allow-missing",
		false,
		"continue even if a referenced variable is not defined",
	)
	rootCmd.AddCommand(interpolateCmd)
}
