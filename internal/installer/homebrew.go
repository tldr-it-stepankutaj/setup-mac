package installer

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

const (
	homebrewInstallScript = "https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh"
	homebrewPath          = "/opt/homebrew/bin/brew"
	homebrewPathIntel     = "/usr/local/bin/brew"
)

// HomebrewInstaller handles Homebrew installation
type HomebrewInstaller struct {
	ctx *Context
}

// NewHomebrewInstaller creates a new Homebrew installer
func NewHomebrewInstaller(ctx *Context) *HomebrewInstaller {
	return &HomebrewInstaller{ctx: ctx}
}

// Name returns the installer name
func (h *HomebrewInstaller) Name() string {
	return "homebrew"
}

// Description returns the installer description
func (h *HomebrewInstaller) Description() string {
	return "Homebrew Package Manager"
}

// IsInstalled checks if Homebrew is installed
func (h *HomebrewInstaller) IsInstalled(ctx context.Context) bool {
	return h.ctx.Executor.Exists("brew")
}

// Install installs Homebrew and configured packages
func (h *HomebrewInstaller) Install(ctx context.Context) error {
	cfg := h.ctx.Config.Homebrew

	if !cfg.Install {
		ui.PrintInfo("Homebrew installation skipped (disabled in config)")
		return nil
	}

	// Install Homebrew if not present
	if !h.IsInstalled(ctx) {
		ui.PrintStep("Installing Homebrew...")
		if err := h.installHomebrew(ctx); err != nil {
			return fmt.Errorf("failed to install Homebrew: %w", err)
		}
	} else {
		ui.PrintInfo("Homebrew already installed")
	}

	// Add taps
	if len(cfg.Taps) > 0 {
		ui.PrintStep("Adding taps...")
		for _, tap := range cfg.Taps {
			if err := h.addTap(ctx, tap); err != nil {
				ui.PrintWarning(fmt.Sprintf("Failed to add tap %s: %v", tap, err))
			}
		}
	}

	// Install formulae
	if len(cfg.Formulae) > 0 {
		ui.PrintStep("Installing formulae...")
		if err := h.installFormulae(ctx, cfg.Formulae); err != nil {
			return fmt.Errorf("failed to install formulae: %w", err)
		}
	}

	// Install casks
	if len(cfg.Casks) > 0 {
		ui.PrintStep("Installing casks...")
		if err := h.installCasks(ctx, cfg.Casks); err != nil {
			return fmt.Errorf("failed to install casks: %w", err)
		}
	}

	return nil
}

func (h *HomebrewInstaller) installHomebrew(ctx context.Context) error {
	cmd := fmt.Sprintf(`/bin/bash -c "$(curl -fsSL %s)"`, homebrewInstallScript)

	if h.ctx.DryRun {
		ui.PrintDryRun(cmd)
		return nil
	}

	// Run Homebrew installer interactively
	if err := h.ctx.Executor.RunInteractive(ctx, "bash", "-c", cmd); err != nil {
		return err
	}

	// Add brew to PATH for current session
	brewPath := h.getBrewPath()
	if _, err := os.Stat(brewPath); err == nil {
		os.Setenv("PATH", fmt.Sprintf("%s:%s", brewPath, os.Getenv("PATH")))
	}

	return nil
}

func (h *HomebrewInstaller) getBrewPath() string {
	if runtime.GOARCH == "arm64" {
		return "/opt/homebrew/bin"
	}
	return "/usr/local/bin"
}

func (h *HomebrewInstaller) addTap(ctx context.Context, tap string) error {
	spinner := ui.NewSpinner(fmt.Sprintf("Adding tap: %s", tap))
	spinner.Start()
	defer spinner.Stop()

	result, err := h.ctx.Executor.Run(ctx, "brew", "tap", tap)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to add tap: %s", tap))
		return err
	}

	if result.DryRun {
		spinner.Info(fmt.Sprintf("[DRY-RUN] Would add tap: %s", tap))
	} else {
		spinner.Success(fmt.Sprintf("Added tap: %s", tap))
	}
	return nil
}

func (h *HomebrewInstaller) installFormulae(ctx context.Context, formulae []string) error {
	// Check which formulae are already installed
	installed := h.getInstalledFormulae(ctx)

	for _, formula := range formulae {
		if installed[formula] {
			ui.PrintInfo(fmt.Sprintf("Formula already installed: %s", formula))
			continue
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Installing: %s", formula))
		spinner.Start()

		result, err := h.ctx.Executor.Run(ctx, "brew", "install", formula)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to install: %s", formula))
			continue
		}

		if result.DryRun {
			spinner.Info(fmt.Sprintf("[DRY-RUN] Would install: %s", formula))
		} else {
			spinner.Success(fmt.Sprintf("Installed: %s", formula))
		}
	}

	return nil
}

func (h *HomebrewInstaller) installCasks(ctx context.Context, casks []string) error {
	// Check which casks are already installed
	installed := h.getInstalledCasks(ctx)

	for _, cask := range casks {
		if installed[cask] {
			ui.PrintInfo(fmt.Sprintf("Cask already installed: %s", cask))
			continue
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Installing cask: %s", cask))
		spinner.Start()

		result, err := h.ctx.Executor.Run(ctx, "brew", "install", "--cask", cask)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to install cask: %s", cask))
			continue
		}

		if result.DryRun {
			spinner.Info(fmt.Sprintf("[DRY-RUN] Would install cask: %s", cask))
		} else {
			spinner.Success(fmt.Sprintf("Installed cask: %s", cask))
		}
	}

	return nil
}

func (h *HomebrewInstaller) getInstalledFormulae(ctx context.Context) map[string]bool {
	installed := make(map[string]bool)

	if h.ctx.DryRun {
		return installed
	}

	result, err := h.ctx.Executor.Run(ctx, "brew", "list", "--formula")
	if err != nil {
		return installed
	}

	for _, line := range strings.Split(result.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			installed[line] = true
		}
	}

	return installed
}

func (h *HomebrewInstaller) getInstalledCasks(ctx context.Context) map[string]bool {
	installed := make(map[string]bool)

	if h.ctx.DryRun {
		return installed
	}

	result, err := h.ctx.Executor.Run(ctx, "brew", "list", "--cask")
	if err != nil {
		return installed
	}

	for _, line := range strings.Split(result.Stdout, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			installed[line] = true
		}
	}

	return installed
}
