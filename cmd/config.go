package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourorg/envdiff/internal/config"
)

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envdiff.json"
	}
	return filepath.Join(home, ".envdiff.json")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage envdiff configuration",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(defaultConfigPath())
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "format:        %s\n", cfg.DefaultFormat)
		fmt.Fprintf(cmd.OutOrStdout(), "show_missing:  %v\n", cfg.ShowMissing)
		fmt.Fprintf(cmd.OutOrStdout(), "show_mismatch: %v\n", cfg.ShowMismatch)
		fmt.Fprintf(cmd.OutOrStdout(), "ignore_keys:   %v\n", cfg.IgnoreKeys)
		return nil
	},
}

var configSetFormatCmd = &cobra.Command{
	Use:   "set-format <text|json|csv>",
	Short: "Set the default output format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt := args[0]
		if fmt != "text" && fmt != "json" && fmt != "csv" {
			return fmt_err(fmt)
		}
		cfgPath := defaultConfigPath()
		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}
		cfg.DefaultFormat = fmt
		if err := config.Save(cfgPath, cfg); err != nil {
			return err
		}
		cmd.Printf("Default format set to %q\n", fmt)
		return nil
	},
}

func fmt_err(f string) error {
	return fmt.Errorf("invalid format %q: must be text, json, or csv", f)
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetFormatCmd)
	rootCmd.AddCommand(configCmd)
}
