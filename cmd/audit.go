package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	auditIncludeUnchanged bool
	auditRedact           bool
)

var auditCmd = &cobra.Command{
	Use:   "audit <file-a> <file-b>",
	Short: "Generate a timestamped audit log of changes between two .env files",
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

		opts := diff.DefaultAuditOptions()
		opts.IncludeUnchanged = auditIncludeUnchanged
		opts.RedactValues = auditRedact

		entries := diff.AuditResults(results, opts)
		diff.WriteAuditReport(os.Stdout, entries)
		return nil
	},
}

func init() {
	auditCmd.Flags().BoolVar(&auditIncludeUnchanged, "include-unchanged", false, "Include keys with identical values in the audit log")
	auditCmd.Flags().BoolVar(&auditRedact, "redact", false, "Redact values from audit output")
	rootCmd.AddCommand(auditCmd)
}
