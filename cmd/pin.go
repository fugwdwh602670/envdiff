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
	pinFilePath  string
	pinComment   string
	pinOverwrite bool
	pinReportOnly bool
)

var pinCmd = &cobra.Command{
	Use:   "pin [file] [KEY...]",
	Short: "Pin specific key values from a .env file",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		envFile := args[0]
		keys := args[1:]

		env, err := parser.ParseFile(envFile)
		if err != nil {
			return fmt.Errorf("parse %s: %w", envFile, err)
		}

		opts := diff.DefaultPinOptions()
		opts.Comment = pinComment
		opts.Overwrite = pinOverwrite
		opts.ReportOnly = pinReportOnly

		added, err := diff.PinEnv(env, keys, pinFilePath, opts)
		if err != nil {
			return err
		}
		diff.WritePinReport(os.Stdout, added, nil)
		return nil
	},
}

var pinCheckCmd = &cobra.Command{
	Use:   "check [file]",
	Short: "Check pinned keys against a .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := parser.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("parse %s: %w", args[0], err)
		}
		violations, err := diff.CheckPins(env, pinFilePath)
		if err != nil {
			return err
		}
		diff.WritePinReport(os.Stdout, nil, violations)
		if len(violations) > 0 {
			lines := make([]string, len(violations))
			for i, v := range violations {
				lines[i] = v.Key
			}
			return fmt.Errorf("pin violations: %s", strings.Join(lines, ", "))
		}
		return nil
	},
}

func init() {
	pinCmd.PersistentFlags().StringVar(&pinFilePath, "pin-file", ".envpins", "path to pin file")
	pinCmd.Flags().StringVar(&pinComment, "comment", "", "comment to attach to pinned entries")
	pinCmd.Flags().BoolVar(&pinOverwrite, "overwrite", false, "overwrite existing pins")
	pinCmd.Flags().BoolVar(&pinReportOnly, "dry-run", false, "report only, do not write")
	pinCmd.AddCommand(pinCheckCmd)
	rootCmd.AddCommand(pinCmd)
}
