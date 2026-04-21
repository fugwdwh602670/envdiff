package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var requireUppercase bool
	var forbidEmpty bool
	var forbidDuplicates bool

	validateCmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate keys and values in a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			env, err := parser.ParseFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", filePath, err)
			}

			opts := diff.ValidateOptions{
				RequireUppercase: requireUppercase,
				ForbidEmpty:      forbidEmpty,
				ForbidDuplicates: forbidDuplicates,
			}

			issues := diff.ValidateEnv(env, opts)
			diff.WriteValidateReport(os.Stdout, issues, filePath)

			if len(issues) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	validateCmd.Flags().BoolVar(&requireUppercase, "require-uppercase", true, "Require all keys to be uppercase")
	validateCmd.Flags().BoolVar(&forbidEmpty, "forbid-empty", true, "Forbid empty values")
	validateCmd.Flags().BoolVar(&forbidDuplicates, "forbid-duplicates", true, "Forbid duplicate keys")

	rootCmd.AddCommand(validateCmd)
}
