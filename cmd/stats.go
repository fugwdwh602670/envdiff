package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats <fileA> [fileB]",
		Short: "Show statistics for one or two .env files",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			envA, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("error reading %s: %w", args[0], err)
			}
			statsA := diff.StatEnv(envA)
			diff.WriteStatsReport(os.Stdout, args[0], statsA)

			if len(args) == 2 {
				envB, err := parser.ParseFile(args[1])
				if err != nil {
					return fmt.Errorf("error reading %s: %w", args[1], err)
				}
				statsB := diff.StatEnv(envB)
				diff.WriteStatsReport(os.Stdout, args[1], statsB)
			}
			return nil
		},
	}
	rootCmd.AddCommand(statsCmd)
}
