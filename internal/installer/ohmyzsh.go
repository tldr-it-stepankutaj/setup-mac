package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stepankutaj/setup-mac/internal/ui"
)

const (
	ohMyZshInstallScript      = "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
	zshAutosuggestionsRepo    = "https://github.com/zsh-users/zsh-autosuggestions"
	zshSyntaxHighlightingRepo = "https://github.com/zsh-users/zsh-syntax-highlighting"
)

// OhMyZshInstaller handles Oh-My-Zsh installation
type OhMyZshInstaller struct {
	ctx *Context
}

// NewOhMyZshInstaller creates a new Oh-My-Zsh installer
func NewOhMyZshInstaller(ctx *Context) *OhMyZshInstaller {
	return &OhMyZshInstaller{ctx: ctx}
}

// Name returns the installer name
func (o *OhMyZshInstaller) Name() string {
	return "oh-my-zsh"
}

// Description returns the installer description
func (o *OhMyZshInstaller) Description() string {
	return "Oh My Zsh Framework"
}

// IsInstalled checks if Oh-My-Zsh is installed
func (o *OhMyZshInstaller) IsInstalled(ctx context.Context) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	omzDir := filepath.Join(homeDir, ".oh-my-zsh")
	_, err = os.Stat(omzDir)
	return err == nil
}

// Install installs Oh-My-Zsh and configured plugins
func (o *OhMyZshInstaller) Install(ctx context.Context) error {
	cfg := o.ctx.Config.Terminal.OhMyZsh

	if !cfg.Install {
		ui.PrintInfo("Oh-My-Zsh installation skipped (disabled in config)")
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Install Oh-My-Zsh if not present
	if !o.IsInstalled(ctx) {
		ui.PrintStep("Installing Oh-My-Zsh...")
		if err := o.installOhMyZsh(ctx); err != nil {
			return fmt.Errorf("failed to install Oh-My-Zsh: %w", err)
		}
	} else {
		ui.PrintInfo("Oh-My-Zsh already installed")
	}

	// Install custom plugins
	if len(cfg.Plugins) > 0 {
		ui.PrintStep("Installing Oh-My-Zsh plugins...")
		if err := o.installPlugins(ctx, homeDir, cfg.Plugins); err != nil {
			return fmt.Errorf("failed to install plugins: %w", err)
		}
	}

	// Configure plugins in .zshrc
	if err := o.configurePlugins(ctx, homeDir, cfg.Plugins); err != nil {
		return fmt.Errorf("failed to configure plugins: %w", err)
	}

	return nil
}

func (o *OhMyZshInstaller) installOhMyZsh(ctx context.Context) error {
	cmd := fmt.Sprintf(`sh -c "$(curl -fsSL %s)" "" --unattended`, ohMyZshInstallScript)

	if o.ctx.DryRun {
		ui.PrintDryRun(cmd)
		return nil
	}

	_, err := o.ctx.Executor.RunShell(ctx, cmd)
	return err
}

func (o *OhMyZshInstaller) installPlugins(ctx context.Context, homeDir string, plugins []string) error {
	customPluginsDir := filepath.Join(homeDir, ".oh-my-zsh", "custom", "plugins")

	// External plugins that need to be cloned
	externalPlugins := map[string]string{
		"zsh-autosuggestions":     zshAutosuggestionsRepo,
		"zsh-syntax-highlighting": zshSyntaxHighlightingRepo,
	}

	for _, plugin := range plugins {
		repo, isExternal := externalPlugins[plugin]
		if !isExternal {
			continue
		}

		pluginDir := filepath.Join(customPluginsDir, plugin)
		if _, err := os.Stat(pluginDir); err == nil {
			ui.PrintInfo(fmt.Sprintf("Plugin already installed: %s", plugin))
			continue
		}

		spinner := ui.NewSpinner(fmt.Sprintf("Installing plugin: %s", plugin))
		spinner.Start()

		result, err := o.ctx.Executor.Run(ctx, "git", "clone", repo, pluginDir)
		if err != nil {
			spinner.Fail(fmt.Sprintf("Failed to install plugin: %s", plugin))
			continue
		}

		if result.DryRun {
			spinner.Info(fmt.Sprintf("[DRY-RUN] Would install plugin: %s", plugin))
		} else {
			spinner.Success(fmt.Sprintf("Installed plugin: %s", plugin))
		}
	}

	return nil
}

func (o *OhMyZshInstaller) configurePlugins(ctx context.Context, homeDir string, plugins []string) error {
	zshrcPath := filepath.Join(homeDir, ".zshrc")

	if o.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would configure plugins in %s: %v", zshrcPath, plugins))
		return nil
	}

	// Read current .zshrc
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return fmt.Errorf("failed to read .zshrc: %w", err)
	}

	// Build plugins line
	pluginsLine := fmt.Sprintf("plugins=(%s)", strings.Join(plugins, " "))

	// Replace existing plugins line or add it
	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "plugins=") {
			lines[i] = pluginsLine
			found = true
			break
		}
	}

	if !found {
		// Add plugins line before source oh-my-zsh.sh
		for i, line := range lines {
			if strings.Contains(line, "source $ZSH/oh-my-zsh.sh") {
				lines = append(lines[:i], append([]string{pluginsLine, ""}, lines[i:]...)...)
				break
			}
		}
	}

	// Write back
	if err := os.WriteFile(zshrcPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write .zshrc: %w", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Configured plugins: %v", plugins))
	return nil
}
