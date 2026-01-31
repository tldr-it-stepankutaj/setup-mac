package installer

import (
	"context"
	"fmt"
	"strings"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// GitInstaller handles Git configuration
type GitInstaller struct {
	ctx *Context
}

// NewGitInstaller creates a new Git installer
func NewGitInstaller(ctx *Context) *GitInstaller {
	return &GitInstaller{ctx: ctx}
}

// Name returns the installer name
func (g *GitInstaller) Name() string {
	return "git"
}

// Description returns the installer description
func (g *GitInstaller) Description() string {
	return "Git Configuration"
}

// IsInstalled returns false (Git config is always "installable")
func (g *GitInstaller) IsInstalled(ctx context.Context) bool {
	return false
}

// Install configures Git
func (g *GitInstaller) Install(ctx context.Context) error {
	cfg := g.ctx.Config.Git

	if !cfg.Configure {
		ui.PrintInfo("Git configuration skipped (disabled in config)")
		return nil
	}

	// Check if git is available
	if !g.ctx.Executor.Exists("git") {
		return fmt.Errorf("git is not installed")
	}

	// Configure user
	ui.PrintStep("Configuring Git user...")
	if err := g.configureUser(ctx); err != nil {
		return fmt.Errorf("failed to configure user: %w", err)
	}

	// Configure aliases
	if len(cfg.Aliases) > 0 {
		ui.PrintStep("Configuring Git aliases...")
		if err := g.configureAliases(ctx); err != nil {
			return fmt.Errorf("failed to configure aliases: %w", err)
		}
	}

	// Configure settings
	if len(cfg.Settings) > 0 {
		ui.PrintStep("Configuring Git settings...")
		if err := g.configureSettings(ctx); err != nil {
			return fmt.Errorf("failed to configure settings: %w", err)
		}
	}

	return nil
}

func (g *GitInstaller) configureUser(ctx context.Context) error {
	user := g.ctx.Config.Git.User

	// Try to get existing git config values as defaults
	existingName := g.getExistingConfig(ctx, "user.name")
	existingEmail := g.getExistingConfig(ctx, "user.email")

	// Get name: config -> existing -> prompt
	name := user.Name
	if name == "" {
		name = existingName
	}
	if name == "" && g.ctx.Config.Settings.Interactive && !g.ctx.DryRun {
		var err error
		name, err = g.ctx.Prompt.Input("Git user name (required for commits)", "")
		if err != nil {
			return fmt.Errorf("failed to get user name: %w", err)
		}
	}

	// Get email: config -> existing -> prompt
	email := user.Email
	if email == "" {
		email = existingEmail
	}
	if email == "" && g.ctx.Config.Settings.Interactive && !g.ctx.DryRun {
		var err error
		email, err = g.ctx.Prompt.Input("Git user email (required for commits)", "")
		if err != nil {
			return fmt.Errorf("failed to get user email: %w", err)
		}
	}

	// Set user name
	if name != "" {
		if err := g.setConfig(ctx, "user.name", name); err != nil {
			return err
		}
	} else if g.ctx.DryRun {
		ui.PrintDryRun("Would prompt for Git user name")
	} else {
		ui.PrintWarning("Git user.name not set - commits will fail without this")
	}

	// Set user email
	if email != "" {
		if err := g.setConfig(ctx, "user.email", email); err != nil {
			return err
		}
	} else if g.ctx.DryRun {
		ui.PrintDryRun("Would prompt for Git user email")
	} else {
		ui.PrintWarning("Git user.email not set - commits will fail without this")
	}

	return nil
}

// getExistingConfig gets an existing git config value
func (g *GitInstaller) getExistingConfig(ctx context.Context, key string) string {
	if g.ctx.DryRun {
		return ""
	}
	result, err := g.ctx.Executor.Run(ctx, "git", "config", "--global", "--get", key)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(result.Stdout)
}

func (g *GitInstaller) configureAliases(ctx context.Context) error {
	for alias, command := range g.ctx.Config.Git.Aliases {
		key := fmt.Sprintf("alias.%s", alias)
		if err := g.setConfig(ctx, key, command); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to set alias %s: %v", alias, err))
		}
	}
	return nil
}

func (g *GitInstaller) configureSettings(ctx context.Context) error {
	for key, value := range g.ctx.Config.Git.Settings {
		if err := g.setConfig(ctx, key, value); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to set %s: %v", key, err))
		}
	}
	return nil
}

func (g *GitInstaller) setConfig(ctx context.Context, key, value string) error {
	if g.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("git config --global %s %q", key, value))
		return nil
	}

	result, err := g.ctx.Executor.Run(ctx, "git", "config", "--global", key, value)
	if err != nil {
		return err
	}

	if result.ExitCode == 0 {
		ui.PrintSuccess(fmt.Sprintf("Set git config: %s = %s", key, value))
	}

	return nil
}
