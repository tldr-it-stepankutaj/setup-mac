package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tldr-it-stepankutaj/setup-mac/configs"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

var (
	configInitForce bool
	configInitPath  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage setup-mac configuration files",
	Long:  `Subcommands for working with the setup-mac YAML configuration.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Write the default config to ~/.config/setup-mac/config.yaml",
	Long: `Write the embedded default configuration to disk so you have a
starting point to edit. By default it goes to ~/.config/setup-mac/config.yaml
and refuses to overwrite an existing file.

Examples:
  # Write to the default location
  setup-mac config init

  # Write to a custom path
  setup-mac config init --output ./my-setup.yaml

  # Overwrite an existing file
  setup-mac config init --force`,
	// User-input errors (e.g. destination exists) shouldn't dump full usage —
	// the warning above the error already explains what to do.
	SilenceUsage: true,
	RunE:         runConfigInit,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)

	configInitCmd.Flags().BoolVarP(&configInitForce, "force", "f", false, "overwrite existing file")
	configInitCmd.Flags().StringVarP(&configInitPath, "output", "o", "", "destination path (default ~/.config/setup-mac/config.yaml)")
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	dest := configInitPath
	if dest == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to resolve home directory: %w", err)
		}
		dest = filepath.Join(homeDir, ".config", "setup-mac", "config.yaml")
	}

	if _, err := os.Stat(dest); err == nil && !configInitForce {
		ui.PrintWarning(fmt.Sprintf("%s already exists — pass --force to overwrite", dest))
		return fmt.Errorf("destination exists")
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(dest, []byte(configs.Default), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	color.New(color.FgGreen).Printf("✓ ")
	fmt.Printf("Wrote default config to %s\n", dest)
	fmt.Println()
	fmt.Println("Edit the file, then run:")
	fmt.Printf("  setup-mac install --all --config %s\n", dest)
	return nil
}
