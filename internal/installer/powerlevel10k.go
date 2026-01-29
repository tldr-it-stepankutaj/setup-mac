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
	powerlevel10kRepo = "https://github.com/romkatv/powerlevel10k.git"
)

// P10kStyle represents a Powerlevel10k style option
type P10kStyle struct {
	Name        string
	Description string
	ConfigFile  string
}

var p10kStyles = []P10kStyle{
	{Name: "lean", Description: "Lean - minimal with no icons", ConfigFile: ""},
	{Name: "classic", Description: "Classic - traditional prompt with icons", ConfigFile: ""},
	{Name: "rainbow", Description: "Rainbow - colorful with icons", ConfigFile: ""},
	{Name: "pure", Description: "Pure - minimalist, inspired by sindresorhus/pure", ConfigFile: ""},
}

// Powerlevel10kInstaller handles Powerlevel10k theme installation
type Powerlevel10kInstaller struct {
	ctx *Context
}

// NewPowerlevel10kInstaller creates a new Powerlevel10k installer
func NewPowerlevel10kInstaller(ctx *Context) *Powerlevel10kInstaller {
	return &Powerlevel10kInstaller{ctx: ctx}
}

// Name returns the installer name
func (p *Powerlevel10kInstaller) Name() string {
	return "powerlevel10k"
}

// Description returns the installer description
func (p *Powerlevel10kInstaller) Description() string {
	return "Powerlevel10k Theme"
}

// IsInstalled checks if Powerlevel10k is installed
func (p *Powerlevel10kInstaller) IsInstalled(ctx context.Context) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	p10kDir := filepath.Join(homeDir, ".oh-my-zsh", "custom", "themes", "powerlevel10k")
	_, err = os.Stat(p10kDir)
	return err == nil
}

// Install installs Powerlevel10k theme
func (p *Powerlevel10kInstaller) Install(ctx context.Context) error {
	cfg := p.ctx.Config.Terminal.Powerlevel10k

	if !cfg.Install {
		ui.PrintInfo("Powerlevel10k installation skipped (disabled in config)")
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Install Powerlevel10k if not present
	if !p.IsInstalled(ctx) {
		ui.PrintStep("Installing Powerlevel10k theme...")
		if err := p.installPowerlevel10k(ctx, homeDir); err != nil {
			return fmt.Errorf("failed to install Powerlevel10k: %w", err)
		}
	} else {
		ui.PrintInfo("Powerlevel10k already installed")
	}

	// Configure theme in .zshrc
	if err := p.configureTheme(ctx, homeDir); err != nil {
		return fmt.Errorf("failed to configure theme: %w", err)
	}

	// Prompt for style selection
	style := cfg.Style
	if style == "" && p.ctx.Config.Settings.Interactive {
		selectedStyle, err := p.promptStyleSelection()
		if err != nil {
			ui.PrintWarning("Style selection skipped")
		} else {
			style = selectedStyle
		}
	}

	if style != "" {
		ui.PrintInfo(fmt.Sprintf("Selected style: %s", style))
		ui.PrintInfo("Run 'p10k configure' to customize your prompt after installation")
	}

	return nil
}

func (p *Powerlevel10kInstaller) installPowerlevel10k(ctx context.Context, homeDir string) error {
	themesDir := filepath.Join(homeDir, ".oh-my-zsh", "custom", "themes")
	p10kDir := filepath.Join(themesDir, "powerlevel10k")

	// Ensure themes directory exists
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	spinner := ui.NewSpinner("Cloning Powerlevel10k repository...")
	spinner.Start()

	result, err := p.ctx.Executor.Run(ctx, "git", "clone", "--depth=1", powerlevel10kRepo, p10kDir)
	if err != nil {
		spinner.Fail("Failed to clone Powerlevel10k")
		return err
	}

	if result.DryRun {
		spinner.Info("[DRY-RUN] Would clone Powerlevel10k")
	} else {
		spinner.Success("Powerlevel10k cloned successfully")
	}

	return nil
}

func (p *Powerlevel10kInstaller) configureTheme(ctx context.Context, homeDir string) error {
	zshrcPath := filepath.Join(homeDir, ".zshrc")

	if p.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would set ZSH_THEME to powerlevel10k/powerlevel10k in %s", zshrcPath))
		return nil
	}

	// Read current .zshrc
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		return fmt.Errorf("failed to read .zshrc: %w", err)
	}

	// Replace theme line
	lines := strings.Split(string(content), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "ZSH_THEME=") {
			lines[i] = `ZSH_THEME="powerlevel10k/powerlevel10k"`
			found = true
			break
		}
	}

	if !found {
		// Add theme line at the beginning after any comments
		for i, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "#") && strings.TrimSpace(line) != "" {
				lines = append(lines[:i], append([]string{`ZSH_THEME="powerlevel10k/powerlevel10k"`}, lines[i:]...)...)
				break
			}
		}
	}

	// Write back
	if err := os.WriteFile(zshrcPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write .zshrc: %w", err)
	}

	ui.PrintSuccess("Theme configured: powerlevel10k/powerlevel10k")
	return nil
}

func (p *Powerlevel10kInstaller) promptStyleSelection() (string, error) {
	items := make([]ui.SelectItem, len(p10kStyles))
	for i, style := range p10kStyles {
		items[i] = ui.SelectItem{
			Name:        style.Name,
			Description: style.Description,
			Value:       style.Name,
		}
	}

	_, selected, err := p.ctx.Prompt.SelectWithDescription("Select Powerlevel10k style", items)
	if err != nil {
		return "", err
	}

	return selected.Value, nil
}
