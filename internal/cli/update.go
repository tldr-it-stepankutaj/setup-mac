package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/config"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/installer"
	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

var (
	updateAll      bool
	updateHomebrew bool
	updateOhMyZsh  bool
	updateDryRun   bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update installed tools and packages",
	Long: `Update installed development tools and packages.

Examples:
  # Update everything
  setup-mac update --all

  # Update specific components
  setup-mac update --homebrew
  setup-mac update --ohmyzsh

  # Dry-run mode
  setup-mac update --all --dry-run`,
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVarP(&updateDryRun, "dry-run", "n", false, "show what would be done without making changes")
	updateCmd.Flags().BoolVarP(&updateAll, "all", "a", false, "update all components")
	updateCmd.Flags().BoolVar(&updateHomebrew, "homebrew", false, "update Homebrew and packages")
	updateCmd.Flags().BoolVar(&updateOhMyZsh, "ohmyzsh", false, "update Oh My Zsh")
}

// Updater interface for components that support updating
type Updater interface {
	Name() string
	Description() string
	Update(ctx context.Context) error
}

func runUpdate(cmd *cobra.Command, args []string) error {
	printBanner()

	// Check if running as root/sudo
	if err := checkNotRoot(); err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override dry-run from flags
	if updateDryRun {
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

	// Determine what to update
	updaters := determineUpdaters(ictx)

	if len(updaters) == 0 {
		ui.PrintWarning("No components selected. Use --all or specific flags like --homebrew, --ohmyzsh")
		return nil
	}

	// Show what will be updated
	ui.PrintInfo(fmt.Sprintf("Updating %d component(s):", len(updaters)))
	for _, u := range updaters {
		fmt.Printf("  - %s\n", u.Description())
	}
	fmt.Println()

	if cfg.Settings.DryRun {
		color.New(color.FgYellow, color.Bold).Println("=== DRY-RUN MODE ===")
		fmt.Println()
	}

	// Run updaters
	var errors []error
	for i, updater := range updaters {
		select {
		case <-ctx.Done():
			return fmt.Errorf("update interrupted")
		default:
			// Print progress
			color.New(color.FgCyan).Printf("[%d/%d] ", i+1, len(updaters))
			fmt.Println(updater.Description())
			fmt.Println("──────────────────────────────────────")

			if err := updater.Update(ctx); err != nil {
				errors = append(errors, fmt.Errorf("%s: %w", updater.Name(), err))
				ui.PrintError(fmt.Sprintf("Failed to update %s: %v", updater.Name(), err))
			}
			fmt.Println()
		}
	}

	// Print summary
	if len(errors) > 0 {
		color.New(color.FgYellow).Println("Update completed with errors:")
		for _, err := range errors {
			color.New(color.FgRed).Printf("  - %v\n", err)
		}
		return fmt.Errorf("%d update(s) failed", len(errors))
	}

	color.New(color.FgGreen, color.Bold).Println("Update completed successfully!")
	return nil
}

func determineUpdaters(ictx *installer.Context) []Updater {
	var updaters []Updater

	if updateAll {
		updaters = append(updaters, installer.NewHomebrewUpdater(ictx))
		updaters = append(updaters, installer.NewOhMyZshUpdater(ictx))
		return updaters
	}

	if updateHomebrew {
		updaters = append(updaters, installer.NewHomebrewUpdater(ictx))
	}

	if updateOhMyZsh {
		updaters = append(updaters, installer.NewOhMyZshUpdater(ictx))
	}

	return updaters
}
