package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var baselinePath string

func init() {
	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage diff baselines",
	}

	saveCmd := &cobra.Command{
		Use:   "save <fileA> <fileB>",
		Short: "Save current diff as a baseline",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}
			b, err := parser.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[1], err)
			}
			results := diff.Compare(a, b)
			if err := diff.SaveBaseline(baselinePath, args[0], args[1], results); err != nil {
				return fmt.Errorf("saving baseline: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Baseline saved to %s\n", baselinePath)
			return nil
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check <fileA> <fileB>",
		Short: "Compare current diff against saved baseline",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := diff.LoadBaseline(baselinePath)
			if err != nil {
				return fmt.Errorf("loading baseline: %w", err)
			}
			a, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[0], err)
			}
			curB, err := parser.ParseFile(args[1])
			if err != nil {
				return fmt.Errorf("parsing %s: %w", args[1], err)
			}
			current := diff.Compare(a, curB)
			newIssues, resolved := diff.DiffBaseline(b.Results, current)
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "New issues: %d\n", len(newIssues))
			for _, r := range newIssues {
				fmt.Fprintf(w, "  + %s (%s)\n", r.Key, r.Status)
			}
			fmt.Fprintf(w, "Resolved: %d\n", len(resolved))
			for _, r := range resolved {
				fmt.Fprintf(w, "  - %s (%s)\n", r.Key, r.Status)
			}
			if len(newIssues) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	for _, c := range []*cobra.Command{saveCmd, checkCmd} {
		c.Flags().StringVar(&baselinePath, "baseline", ".envdiff-baseline.json", "Path to baseline file")
	}
	baselineCmd.AddCommand(saveCmd, checkCmd)
	rootCmd.AddCommand(baselineCmd)
}
