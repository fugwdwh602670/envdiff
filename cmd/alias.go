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
	aliasFile    string
	aliasRawPairs []string
	aliasAudit   bool
)

var aliasCmd = &cobra.Command{
	Use:   "alias <file>",
	Short: "Resolve key aliases in a .env file to their canonical names",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := parser.ParseFile(args[0])
		if err != nil {
			return fmt.Errorf("parse error: %w", err)
		}

		aliases := diff.AliasMap{}

		// Parse --alias canonical=ALT1,ALT2 flags
		for _, pair := range aliasRawPairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid alias format %q: expected canonical=ALT1,ALT2", pair)
			}
			canonical := strings.TrimSpace(parts[0])
			alts := strings.Split(parts[1], ",")
			for i, a := range alts {
				alts[i] = strings.TrimSpace(a)
			}
			aliases[canonical] = append(aliases[canonical], alts...)
		}

		if aliasAudit {
			reports := diff.AuditAliases(env, aliases)
			diff.WriteAliasReport(os.Stdout, reports)
			return nil
		}

		resolved := diff.ResolveAliases(env, aliases)
		for k, v := range resolved {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	aliasCmd.Flags().StringArrayVar(&aliasRawPairs, "alias", nil, "Alias mapping in form canonical=ALT1,ALT2 (repeatable)")
	aliasCmd.Flags().BoolVar(&aliasAudit, "audit", false, "Print alias resolution report instead of resolved env")
	rootCmd.AddCommand(aliasCmd)
}
