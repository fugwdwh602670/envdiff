package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

var (
	cascadeOverwrite bool
	cascadeSkipEmpty bool
	cascadeOutput    string
)

var cascadeCmd = &cobra.Command{
	Use:   "cascade <base.env> [layer.env...]",
	Short: "Merge multiple .env files in priority order (last wins)",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var layers []map[string]string
		for _, path := range args {
			env, err := parser.ParseFile(path)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}
			layers = append(layers, env)
		}

		opts := diff.DefaultCascadeOptions()
		opts.Overwrite = cascadeOverwrite
		opts.SkipEmpty = cascadeSkipEmpty

		_, entries := diff.CascadeEnv(layers, opts)

		w := cmd.OutOrStdout()
		if cascadeOutput != "" {
			f, err := os.Create(cascadeOutput)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		diff.WriteCascadeReport(w, entries, args)
		return nil
	},
}

func init() {
	cascadeCmd.Flags().BoolVar(&cascadeOverwrite, "overwrite", true, "later layers overwrite earlier ones")
	cascadeCmd.Flags().BoolVar(&cascadeSkipEmpty, "skip-empty", false, "ignore empty values from source layers")
	cascadeCmd.Flags().StringVar(&cascadeOutput, "output", "", "write report to file instead of stdout")
	rootCmd.AddCommand(cascadeCmd)
}
