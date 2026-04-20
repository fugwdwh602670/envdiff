package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/diff"
	"envdiff/internal/parser"
)

func init() {
	var passphrase string
	var patterns []string
	var outputFile string

	encryptCmd := &cobra.Command{
		Use:   "encrypt <file>",
		Short: "Encrypt sensitive values in a .env file",
		Long: `Encrypt values for keys matching sensitive patterns (e.g. PASSWORD, TOKEN, KEY).
Encrypted values are prefixed with "enc:" and can be decrypted with the same passphrase.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if passphrase == "" {
				passphrase = os.Getenv("ENVDIFF_PASSPHRASE")
			}
			if passphrase == "" {
				return fmt.Errorf("passphrase required: use --passphrase or set ENVDIFF_PASSPHRASE")
			}

			env, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse %s: %w", args[0], err)
			}

			opts := diff.DefaultEncryptOptions()
			opts.Passphrase = passphrase
			if len(patterns) > 0 {
				opts.Patterns = patterns
			}

			encrypted, err := diff.EncryptEnv(env, opts)
			if err != nil {
				return fmt.Errorf("encrypt: %w", err)
			}

			w := cmd.OutOrStdout()
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("create output: %w", err)
				}
				defer f.Close()
				w = f
			}

			for k, v := range encrypted {
				fmt.Fprintf(w, "%s=%s\n", k, v)
			}

			diff.WriteEncryptReport(cmd.ErrOrStderr(), env, encrypted)
			return nil
		},
	}

	encryptCmd.Flags().StringVar(&passphrase, "passphrase", "", "Passphrase for AES-256 encryption")
	encryptCmd.Flags().StringSliceVar(&patterns, "patterns", nil, "Key patterns to encrypt (overrides defaults)")
	encryptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write encrypted env to file instead of stdout")

	rootCmd.AddCommand(encryptCmd)
}
