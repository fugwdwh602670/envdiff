package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
)

var ignoreFile string

var ignoreCmd = &cobra.Command{
	Use:   "ignore [keys...]",
	Short: "Test which keys would be suppressed by an ignore file",
	Long:  "Loads an ignore file and prints the keys it contains, or checks whether given keys are ignored.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ignoreFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		il, err := diff.LoadIgnoreFile(ignoreFile)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("ignore file not found: %s", ignoreFile)
			}
			return err
		}

		if len(args) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "ignore file loaded: %d key(s) suppressed\n", il.Len())
			return nil
		}

		for _, key := range args {
			if il.Contains(key) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: ignored\n", key)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: not ignored\n", key)
			}
		}
		return nil
	},
}

func init() {
	ignoreCmd.Flags().StringVarP(&ignoreFile, "file", "f", "", "path to ignore file (one key per line)")
	rootCmd.AddCommand(ignoreCmd)
}
