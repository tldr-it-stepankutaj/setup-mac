package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stepankutaj/setup-mac/internal/config"
	"github.com/stepankutaj/setup-mac/internal/installer"
	"github.com/stepankutaj/setup-mac/internal/ui"
)

var (
	dryRun          bool
	installAll      bool
	installHomebrew bool
	installTerminal bool
	installShell    bool
	installMacOS    bool
	installGit      bool
	installSSH      bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure development tools",
	Long: `Install and configure development tools on your macOS system.

Examples:
  # Install everything
  setup-mac install --all

  # Install specific components
  setup-mac install --homebrew
  setup-mac install --terminal
  setup-mac install --shell

  # Dry-run mode (show what would be done)
  setup-mac install --all --dry-run

  # Use custom config
  setup-mac install --all --config my-config.yaml`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would be done without making changes")
	installCmd.Flags().BoolVarP(&installAll, "all", "a", false, "install all components")
	installCmd.Flags().BoolVar(&installHomebrew, "homebrew", false, "install Homebrew and packages")
	installCmd.Flags().BoolVar(&installTerminal, "terminal", false, "install Oh-My-Zsh and Powerlevel10k")
	installCmd.Flags().BoolVar(&installShell, "shell", false, "configure shell aliases and environment")
	installCmd.Flags().BoolVar(&installMacOS, "macos", false, "configure macOS defaults")
	installCmd.Flags().BoolVar(&installGit, "git", false, "configure Git")
	installCmd.Flags().BoolVar(&installSSH, "ssh", false, "generate SSH key")
}

func runInstall(cmd *cobra.Command, args []string) error {
	printBanner()

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override dry-run from flags
	if dryRun {
		cfg.Settings.DryRun = true
	}

	// Create installer context
	ictx := installer.NewContext(cfg, cfg.Settings.DryRun, verbose)

	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		color.New(color.FgYellow).Println("\nInterrupted, cleaning up...")
		cancel()
	}()

	// Determine what to install
	installersToRun := determineInstallers(ictx)

	if len(installersToRun) == 0 {
		ui.PrintWarning("No components selected. Use --all or specific flags like --homebrew, --terminal, etc.")
		return nil
	}

	// Show what will be installed
	ui.PrintInfo(fmt.Sprintf("Installing %d component(s):", len(installersToRun)))
	for _, i := range installersToRun {
		fmt.Printf("  - %s\n", i.Description())
	}
	fmt.Println()

	if cfg.Settings.DryRun {
		color.New(color.FgYellow, color.Bold).Println("=== DRY-RUN MODE ===")
		fmt.Println()
	}

	// Confirm if interactive
	if cfg.Settings.Interactive && !cfg.Settings.DryRun {
		confirm, err := ictx.Prompt.Confirm("Proceed with installation?", true)
		if err != nil || !confirm {
			ui.PrintInfo("Installation cancelled")
			return nil
		}
		fmt.Println()
	}

	// Run installers
	var errors []error
	for _, inst := range installersToRun {
		select {
		case <-ctx.Done():
			return fmt.Errorf("installation interrupted")
		default:
			if err := installer.RunInstaller(ctx, inst, ictx); err != nil {
				errors = append(errors, fmt.Errorf("%s: %w", inst.Name(), err))
			}
		}
	}

	// Print summary
	fmt.Println()
	if len(errors) > 0 {
		color.New(color.FgYellow).Println("Installation completed with errors:")
		for _, err := range errors {
			color.New(color.FgRed).Printf("  - %v\n", err)
		}
		return fmt.Errorf("%d installer(s) failed", len(errors))
	}

	color.New(color.FgGreen, color.Bold).Println("Installation completed successfully!")

	if !cfg.Settings.DryRun {
		fmt.Println()
		ui.PrintInfo("You may need to restart your terminal for all changes to take effect.")
	}

	return nil
}

func determineInstallers(ictx *installer.Context) []installer.Installer {
	var installers []installer.Installer

	if installAll {
		// Install all in order
		installers = append(installers, installer.NewHomebrewInstaller(ictx))
		installers = append(installers, installer.NewOhMyZshInstaller(ictx))
		installers = append(installers, installer.NewPowerlevel10kInstaller(ictx))
		installers = append(installers, installer.NewShellInstaller(ictx))
		installers = append(installers, installer.NewMacOSInstaller(ictx))
		installers = append(installers, installer.NewGitInstaller(ictx))
		installers = append(installers, installer.NewSSHInstaller(ictx))
		return installers
	}

	if installHomebrew {
		installers = append(installers, installer.NewHomebrewInstaller(ictx))
	}

	if installTerminal {
		installers = append(installers, installer.NewOhMyZshInstaller(ictx))
		installers = append(installers, installer.NewPowerlevel10kInstaller(ictx))
	}

	if installShell {
		installers = append(installers, installer.NewShellInstaller(ictx))
	}

	if installMacOS {
		installers = append(installers, installer.NewMacOSInstaller(ictx))
	}

	if installGit {
		installers = append(installers, installer.NewGitInstaller(ictx))
	}

	if installSSH {
		installers = append(installers, installer.NewSSHInstaller(ictx))
	}

	return installers
}
