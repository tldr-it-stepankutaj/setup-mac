package installer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// MacOSInstaller handles macOS defaults configuration
type MacOSInstaller struct {
	ctx *Context
}

// NewMacOSInstaller creates a new macOS installer
func NewMacOSInstaller(ctx *Context) *MacOSInstaller {
	return &MacOSInstaller{ctx: ctx}
}

// Name returns the installer name
func (m *MacOSInstaller) Name() string {
	return "macos"
}

// Description returns the installer description
func (m *MacOSInstaller) Description() string {
	return "macOS System Defaults"
}

// IsInstalled returns false (macOS defaults are always "installable")
func (m *MacOSInstaller) IsInstalled(ctx context.Context) bool {
	return false
}

// Install configures macOS defaults
func (m *MacOSInstaller) Install(ctx context.Context) error {
	cfg := m.ctx.Config.MacOS

	if !cfg.Configure {
		ui.PrintInfo("macOS defaults configuration skipped (disabled in config)")
		return nil
	}

	// Configure Dock
	ui.PrintStep("Configuring Dock...")
	if err := m.configureDock(ctx); err != nil {
		ui.PrintWarning(fmt.Sprintf("Some Dock settings failed: %v", err))
	}

	// Configure Finder
	ui.PrintStep("Configuring Finder...")
	if err := m.configureFinder(ctx); err != nil {
		ui.PrintWarning(fmt.Sprintf("Some Finder settings failed: %v", err))
	}

	// Configure Keyboard
	ui.PrintStep("Configuring Keyboard...")
	if err := m.configureKeyboard(ctx); err != nil {
		ui.PrintWarning(fmt.Sprintf("Some Keyboard settings failed: %v", err))
	}

	// Restart affected apps
	if !m.ctx.DryRun {
		ui.PrintStep("Restarting affected applications...")
		m.restartApps(ctx)
	}

	return nil
}

func (m *MacOSInstaller) configureDock(ctx context.Context) error {
	dock := m.ctx.Config.MacOS.Defaults.Dock

	defaults := []struct {
		domain string
		key    string
		typ    string
		value  string
	}{
		{"com.apple.dock", "autohide", "bool", strconv.FormatBool(dock.Autohide)},
		{"com.apple.dock", "autohide-delay", "float", strconv.Itoa(dock.AutohideDelay)},
		{"com.apple.dock", "tilesize", "int", strconv.Itoa(dock.TileSize)},
		{"com.apple.dock", "magnification", "bool", strconv.FormatBool(dock.Magnification)},
		{"com.apple.dock", "minimize-to-application", "bool", strconv.FormatBool(dock.MinimizeToApp)},
		{"com.apple.dock", "show-recents", "bool", strconv.FormatBool(dock.ShowRecents)},
	}

	return m.applyDefaults(ctx, defaults)
}

func (m *MacOSInstaller) configureFinder(ctx context.Context) error {
	finder := m.ctx.Config.MacOS.Defaults.Finder

	defaults := []struct {
		domain string
		key    string
		typ    string
		value  string
	}{
		{"com.apple.finder", "AppleShowAllFiles", "bool", strconv.FormatBool(finder.ShowHiddenFiles)},
		{"NSGlobalDomain", "AppleShowAllExtensions", "bool", strconv.FormatBool(finder.ShowExtensions)},
		{"com.apple.finder", "ShowPathbar", "bool", strconv.FormatBool(finder.ShowPathBar)},
		{"com.apple.finder", "ShowStatusBar", "bool", strconv.FormatBool(finder.ShowStatusBar)},
	}

	// View style
	if finder.DefaultViewStyle != "" {
		viewStyles := map[string]string{
			"icon":    "icnv",
			"list":    "Nlsv",
			"column":  "clmv",
			"gallery": "glyv",
		}
		if style, ok := viewStyles[finder.DefaultViewStyle]; ok {
			defaults = append(defaults, struct {
				domain string
				key    string
				typ    string
				value  string
			}{"com.apple.finder", "FXPreferredViewStyle", "string", style})
		}
	}

	return m.applyDefaults(ctx, defaults)
}

func (m *MacOSInstaller) configureKeyboard(ctx context.Context) error {
	keyboard := m.ctx.Config.MacOS.Defaults.Keyboard

	defaults := []struct {
		domain string
		key    string
		typ    string
		value  string
	}{
		{"NSGlobalDomain", "KeyRepeat", "int", strconv.Itoa(keyboard.KeyRepeat)},
		{"NSGlobalDomain", "InitialKeyRepeat", "int", strconv.Itoa(keyboard.InitialKeyRepeat)},
		{"NSGlobalDomain", "NSAutomaticQuoteSubstitutionEnabled", "bool", strconv.FormatBool(!keyboard.DisableSmartQuotes)},
		{"NSGlobalDomain", "NSAutomaticDashSubstitutionEnabled", "bool", strconv.FormatBool(!keyboard.DisableSmartDashes)},
	}

	return m.applyDefaults(ctx, defaults)
}

func (m *MacOSInstaller) applyDefaults(ctx context.Context, defaults []struct {
	domain string
	key    string
	typ    string
	value  string
}) error {
	for _, d := range defaults {
		var args []string
		switch d.typ {
		case "bool":
			args = []string{"write", d.domain, d.key, "-bool", d.value}
		case "int":
			args = []string{"write", d.domain, d.key, "-int", d.value}
		case "float":
			args = []string{"write", d.domain, d.key, "-float", d.value}
		case "string":
			args = []string{"write", d.domain, d.key, "-string", d.value}
		}

		if m.ctx.DryRun {
			ui.PrintDryRun(fmt.Sprintf("defaults %s", joinArgs(args)))
			continue
		}

		result, err := m.ctx.Executor.Run(ctx, "defaults", args...)
		if err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to set %s %s: %v", d.domain, d.key, err))
			continue
		}

		if result.ExitCode == 0 {
			ui.PrintSuccess(fmt.Sprintf("Set %s %s = %s", d.domain, d.key, d.value))
		}
	}

	return nil
}

func (m *MacOSInstaller) restartApps(ctx context.Context) {
	apps := []string{"Dock", "Finder"}

	for _, app := range apps {
		if _, err := m.ctx.Executor.Run(ctx, "killall", app); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to restart %s", app))
		} else {
			ui.PrintSuccess(fmt.Sprintf("Restarted %s", app))
		}
	}
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		result += arg
	}
	return result
}
