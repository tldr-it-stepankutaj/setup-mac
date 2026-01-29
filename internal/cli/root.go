package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
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
