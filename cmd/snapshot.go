package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	snapshotLabel  string
	snapshotOutput string
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage env snapshots",
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <file>",
	Short: "Save a snapshot of an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := parser.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		out := snapshotOutput
		if out == "" {
			out = args[0] + ".snapshot.json"
		}
		if err := diff.SaveSnapshot(out, env, snapshotLabel); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s\n", out)
		return nil
	},
}

var snapshotCheckCmd = &cobra.Command{
	Use:   "check <snapshot> <file>",
	Short: "Compare an env file against a saved snapshot",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		snap, err := diff.LoadSnapshot(args[0])
		if err != nil {
			return fmt.Errorf("load snapshot: %w", err)
		}
		live, err := parser.ParseFile(args[1])
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}
		results := diff.DiffSnapshot(snap, live)
		diff.WriteReport(cmd.OutOrStdout(), results, args[0], args[1])
		for _, r := range results {
			if r.Status != diff.StatusMatch {
				os.Exit(1)
			}
		}
		return nil
	},
}

func init() {
	snapshotSaveCmd.Flags().StringVarP(&snapshotLabel, "label", "l", "", "Label for the snapshot")
	snapshotSaveCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "", "Output path for snapshot file")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotCheckCmd)
	rootCmd.AddCommand(snapshotCmd)
}
