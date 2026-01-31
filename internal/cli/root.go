package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/installer"
)

var (
	cfgFile         string
	verbose         bool
	skipUpdateCheck bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "setup-mac",
	Short: "macOS developer setup CLI tool",
	Long: `setup-mac is a CLI tool for automating macOS developer environment setup.

It installs and configures:
  - Homebrew and packages
  - iTerm2 terminal
  - Oh-My-Zsh with plugins
  - Powerlevel10k theme
  - Shell aliases and environment
  - macOS system defaults
  - Git configuration
  - SSH keys

Configuration is done via YAML files.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip update check for certain commands or if disabled
		if skipUpdateCheck || cmd.Name() == "version" || cmd.Name() == "help" || cmd.Name() == "completion" {
			return
		}

		// Check for updates in background (non-blocking, with short timeout)
		checkForUpdates()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand, show help
		_ = cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default: embedded defaults)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&skipUpdateCheck, "skip-update-check", false, "skip checking for updates")
}

// checkForUpdates checks for new versions on GitHub
func checkForUpdates() {
	// Don't check for dev versions
	if Version == "dev" || Version == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := installer.NewVersionChecker(Version)
	release, isNewer, err := checker.CheckForUpdate(ctx)
	if err != nil {
		// Silently ignore errors
		return
	}

	if isNewer {
		fmt.Println()
		color.New(color.FgYellow).Printf("⬆ ")
		fmt.Printf("New version available: %s (current: %s)\n",
			color.CyanString(release.TagName),
			color.YellowString(Version))
		fmt.Printf("  Download: %s\n", checker.GetDownloadURL(release))
		fmt.Printf("  Run with --skip-update-check to disable this message\n")
		fmt.Println()
	}
}

// printBanner prints the application banner
func printBanner() {
	banner := `
███████╗███████╗████████╗██╗   ██╗██████╗       ███╗   ███╗ █████╗  ██████╗
██╔════╝██╔════╝╚══██╔══╝██║   ██║██╔══██╗      ████╗ ████║██╔══██╗██╔════╝
███████╗█████╗     ██║   ██║   ██║██████╔╝█████╗██╔████╔██║███████║██║
╚════██║██╔══╝     ██║   ██║   ██║██╔═══╝ ╚════╝██║╚██╔╝██║██╔══██║██║
███████║███████╗   ██║   ╚██████╔╝██║           ██║ ╚═╝ ██║██║  ██║╚██████╗
╚══════╝╚══════╝   ╚═╝    ╚═════╝ ╚═╝           ╚═╝     ╚═╝╚═╝  ╚═╝ ╚═════╝
`
	color.New(color.FgCyan).Println(banner)
	fmt.Println("  macOS Developer Setup Tool")
	fmt.Println()
}
