package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// HomebrewUpdater handles Homebrew updates
type HomebrewUpdater struct {
	ctx *Context
}

// NewHomebrewUpdater creates a new Homebrew updater
func NewHomebrewUpdater(ctx *Context) *HomebrewUpdater {
	return &HomebrewUpdater{ctx: ctx}
}

// Name returns the updater name
func (h *HomebrewUpdater) Name() string {
	return "homebrew"
}

// Description returns the updater description
func (h *HomebrewUpdater) Description() string {
	return "Homebrew Package Manager"
}

// Update updates Homebrew and all packages
func (h *HomebrewUpdater) Update(ctx context.Context) error {
	// Check if Homebrew is installed
	if !h.ctx.Executor.Exists("brew") {
		return fmt.Errorf("homebrew is not installed")
	}

	// Update Homebrew itself
	ui.PrintStep("Updating Homebrew...")
	if h.ctx.DryRun {
		ui.PrintDryRun("brew update")
	} else {
		spinner := ui.NewSpinner("Running brew update...")
		spinner.Start()
		result, err := h.ctx.Executor.Run(ctx, "brew", "update")
		if err != nil {
			spinner.Fail("Failed to update Homebrew")
			return fmt.Errorf("brew update failed: %w\n%s", err, result.Stderr)
		}
		spinner.Success("Homebrew updated")
	}

	// Upgrade all packages
	ui.PrintStep("Upgrading packages...")
	if h.ctx.DryRun {
		ui.PrintDryRun("brew upgrade")
	} else {
		spinner := ui.NewSpinner("Running brew upgrade...")
		spinner.Start()
		result, err := h.ctx.Executor.Run(ctx, "brew", "upgrade")
		if err != nil {
			spinner.Fail("Failed to upgrade packages")
			return fmt.Errorf("brew upgrade failed: %w\n%s", err, result.Stderr)
		}
		spinner.Success("Packages upgraded")
	}

	// Upgrade casks
	ui.PrintStep("Upgrading casks...")
	if h.ctx.DryRun {
		ui.PrintDryRun("brew upgrade --cask")
	} else {
		spinner := ui.NewSpinner("Running brew upgrade --cask...")
		spinner.Start()
		result, err := h.ctx.Executor.Run(ctx, "brew", "upgrade", "--cask")
		if err != nil {
			// Cask upgrade failures are often non-critical (app already running, etc.)
			spinner.Warning(fmt.Sprintf("Some casks may not have been upgraded: %s", result.Stderr))
		} else {
			spinner.Success("Casks upgraded")
		}
	}

	// Cleanup old versions
	ui.PrintStep("Cleaning up...")
	if h.ctx.DryRun {
		ui.PrintDryRun("brew cleanup")
	} else {
		spinner := ui.NewSpinner("Running brew cleanup...")
		spinner.Start()
		_, err := h.ctx.Executor.Run(ctx, "brew", "cleanup")
		if err != nil {
			spinner.Warning("Cleanup had some issues (non-critical)")
		} else {
			spinner.Success("Cleanup complete")
		}
	}

	return nil
}

// OhMyZshUpdater handles Oh My Zsh updates
type OhMyZshUpdater struct {
	ctx *Context
}

// NewOhMyZshUpdater creates a new Oh My Zsh updater
func NewOhMyZshUpdater(ctx *Context) *OhMyZshUpdater {
	return &OhMyZshUpdater{ctx: ctx}
}

// Name returns the updater name
func (o *OhMyZshUpdater) Name() string {
	return "ohmyzsh"
}

// Description returns the updater description
func (o *OhMyZshUpdater) Description() string {
	return "Oh My Zsh Framework"
}

// Update updates Oh My Zsh
func (o *OhMyZshUpdater) Update(ctx context.Context) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	omzDir := filepath.Join(homeDir, ".oh-my-zsh")

	// Check if Oh My Zsh is installed
	if _, err := os.Stat(omzDir); os.IsNotExist(err) {
		return fmt.Errorf("oh My Zsh is not installed")
	}

	ui.PrintStep("Updating Oh My Zsh...")

	if o.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("cd %s && git pull", omzDir))
		return nil
	}

	// Update using git pull
	spinner := ui.NewSpinner("Pulling latest changes...")
	spinner.Start()

	// Save current directory and change to omz dir
	result, err := o.ctx.Executor.Run(ctx, "git", "-C", omzDir, "pull", "--rebase", "--stat", "origin", "master")
	if err != nil {
		spinner.Fail("Failed to update Oh My Zsh")
		return fmt.Errorf("git pull failed: %w\n%s", err, result.Stderr)
	}

	spinner.Success("Oh My Zsh updated")

	// Update custom plugins if any
	customPluginsDir := filepath.Join(omzDir, "custom", "plugins")
	if entries, err := os.ReadDir(customPluginsDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			pluginDir := filepath.Join(customPluginsDir, entry.Name())
			gitDir := filepath.Join(pluginDir, ".git")

			if _, err := os.Stat(gitDir); os.IsNotExist(err) {
				continue // Not a git repo
			}

			spinner := ui.NewSpinner(fmt.Sprintf("Updating plugin: %s...", entry.Name()))
			spinner.Start()

			result, err := o.ctx.Executor.Run(ctx, "git", "-C", pluginDir, "pull", "--rebase")
			if err != nil {
				spinner.Warning(fmt.Sprintf("Failed to update plugin %s", entry.Name()))
			} else {
				_ = result
				spinner.Success(fmt.Sprintf("Updated plugin: %s", entry.Name()))
			}
		}
	}

	// Update custom themes (like powerlevel10k)
	customThemesDir := filepath.Join(omzDir, "custom", "themes")
	if entries, err := os.ReadDir(customThemesDir); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			themeDir := filepath.Join(customThemesDir, entry.Name())
			gitDir := filepath.Join(themeDir, ".git")

			if _, err := os.Stat(gitDir); os.IsNotExist(err) {
				continue // Not a git repo
			}

			spinner := ui.NewSpinner(fmt.Sprintf("Updating theme: %s...", entry.Name()))
			spinner.Start()

			result, err := o.ctx.Executor.Run(ctx, "git", "-C", themeDir, "pull", "--rebase")
			if err != nil {
				spinner.Warning(fmt.Sprintf("Failed to update theme %s", entry.Name()))
			} else {
				_ = result
				spinner.Success(fmt.Sprintf("Updated theme: %s", entry.Name()))
			}
		}
	}

	return nil
}
