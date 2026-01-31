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
		// Check if already installed (exact match or base name match for versioned packages)
		if h.isFormulaInstalled(formula, installed) {
			ui.PrintInfo(fmt.Sprintf("Formula already installed: %s", formula))
			continue
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Installing: %s", formula))
		spinner.Start()

		result, err := h.ctx.Executor.Run(ctx, "brew", "install", formula)
		if err != nil {
			// Check if it's actually installed despite the error (e.g., already installed warning)
			if h.isFormulaInstalled(formula, h.getInstalledFormulae(ctx)) {
				spinner.Success(fmt.Sprintf("Already installed: %s", formula))
				continue
			}
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

// isFormulaInstalled checks if a formula is installed, handling versioned packages
func (h *HomebrewInstaller) isFormulaInstalled(formula string, installed map[string]bool) bool {
	// Exact match
	if installed[formula] {
		return true
	}

	// Check for versioned variants (e.g., "node" matches "node@18", "node@20")
	baseName := formula
	if idx := strings.Index(formula, "@"); idx != -1 {
		baseName = formula[:idx]
	}

	for pkg := range installed {
		pkgBase := pkg
		if idx := strings.Index(pkg, "@"); idx != -1 {
			pkgBase = pkg[:idx]
		}
		if pkgBase == baseName {
			return true
		}
	}

	return false
}

func (h *HomebrewInstaller) installCasks(ctx context.Context, casks []string) error {
	// Check which casks are already installed
	installed := h.getInstalledCasks(ctx)

	// Also check Applications folder for already installed apps
	installedApps := h.getInstalledApplications()

	for _, cask := range casks {
		if installed[cask] {
			ui.PrintInfo(fmt.Sprintf("Cask already installed: %s", cask))
			continue
		}

		// Check if app is already in /Applications (manually installed)
		if h.isCaskAppInstalled(cask, installedApps) {
			ui.PrintInfo(fmt.Sprintf("Application already installed (not via Homebrew): %s", cask))
			continue
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Installing cask: %s", cask))
		spinner.Start()

		result, err := h.ctx.Executor.Run(ctx, "brew", "install", "--cask", cask)
		if err != nil {
			// Check if it failed because already installed
			if result != nil && strings.Contains(result.Stderr, "already installed") {
				spinner.Success(fmt.Sprintf("Already installed: %s", cask))
				continue
			}
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

// getInstalledApplications returns a list of apps in /Applications
func (h *HomebrewInstaller) getInstalledApplications() map[string]bool {
	apps := make(map[string]bool)

	entries, err := os.ReadDir("/Applications")
	if err != nil {
		return apps
	}

	for _, entry := range entries {
		name := entry.Name()
		// Remove .app suffix and lowercase for comparison
		if strings.HasSuffix(name, ".app") {
			name = strings.TrimSuffix(name, ".app")
			apps[strings.ToLower(name)] = true
		}
	}

	return apps
}

// isCaskAppInstalled checks if a cask's app is already installed
func (h *HomebrewInstaller) isCaskAppInstalled(cask string, installedApps map[string]bool) bool {
	// Common cask name to app name mappings
	caskToApp := map[string]string{
		"visual-studio-code": "visual studio code",
		"google-chrome":      "google chrome",
		"sublime-text":       "sublime text",
		"intellij-idea":      "intellij idea",
		"intellij-idea-ce":   "intellij idea ce",
	}

	// Check mapping first
	if appName, ok := caskToApp[cask]; ok {
		return installedApps[appName]
	}

	// Try direct match (replace hyphens with spaces)
	normalizedCask := strings.ReplaceAll(cask, "-", " ")
	return installedApps[strings.ToLower(normalizedCask)] || installedApps[cask]
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
