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
	transformUpper   bool
	transformLower   bool
	transformTrimmed bool
	transformPrefix  string
	transformOutput  string
)

func init() {
	transformCmd := &cobra.Command{
		Use:   "transform <file>",
		Short: "Apply value transformations to a .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse error: %w", err)
			}

			opts := diff.DefaultTransformOptions()

			if transformUpper {
				opts.ApplyToAll = strings.ToUpper
			} else if transformLower {
				opts.ApplyToAll = strings.ToLower
			} else if transformTrimmed {
				opts.ApplyToAll = strings.TrimSpace
			}

			if transformPrefix != "" {
				parts := strings.SplitN(transformPrefix, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("--prefix-transform must be in format PREFIX=ACTION")
				}
				pfx, action := parts[0], parts[1]
				switch action {
				case "upper":
					opts.Transforms[pfx] = strings.ToUpper
				case "lower":
					opts.Transforms[pfx] = strings.ToLower
				case "trim":
					opts.Transforms[pfx] = strings.TrimSpace
				default:
					return fmt.Errorf("unknown action %q; use upper, lower, or trim", action)
				}
			}

			_, results := diff.TransformEnv(env, opts)

			w := os.Stdout
			if transformOutput != "" {
				f, err := os.Create(transformOutput)
				if err != nil {
					return fmt.Errorf("cannot create output file: %w", err)
				}
				defer f.Close()
				w = f
			}

			diff.WriteTransformReport(w, results)
			return nil
		},
	}

	transformCmd.Flags().BoolVar(&transformUpper, "upper", false, "Convert all values to uppercase")
	transformCmd.Flags().BoolVar(&transformLower, "lower", false, "Convert all values to lowercase")
	transformCmd.Flags().BoolVar(&transformTrimmed, "trim", false, "Trim whitespace from all values")
	transformCmd.Flags().StringVar(&transformPrefix, "prefix-transform", "", "Apply transform to keys with prefix (e.g. DB_=upper)")
	transformCmd.Flags().StringVar(&transformOutput, "output", "", "Write report to file instead of stdout")

	rootCmd.AddCommand(transformCmd)
}
