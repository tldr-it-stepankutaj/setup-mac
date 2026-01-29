package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// ShellInstaller handles shell configuration
type ShellInstaller struct {
	ctx *Context
}

// NewShellInstaller creates a new shell installer
func NewShellInstaller(ctx *Context) *ShellInstaller {
	return &ShellInstaller{ctx: ctx}
}

// Name returns the installer name
func (s *ShellInstaller) Name() string {
	return "shell"
}

// Description returns the installer description
func (s *ShellInstaller) Description() string {
	return "Shell Configuration"
}

// IsInstalled checks if shell configuration exists
func (s *ShellInstaller) IsInstalled(ctx context.Context) bool {
	// Shell config is always "installable" (can be updated)
	return false
}

// Install configures shell aliases, environment variables, and extras
func (s *ShellInstaller) Install(ctx context.Context) error {
	cfg := s.ctx.Config.Shell

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	zshrcPath := filepath.Join(homeDir, ".zshrc")

	// Backup existing .zshrc if configured (skip in dry-run mode)
	if s.ctx.Config.Settings.BackupDotfiles && !s.ctx.DryRun {
		if err := s.backupFile(zshrcPath); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to backup .zshrc: %v", err))
		}
	} else if s.ctx.Config.Settings.BackupDotfiles && s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would backup %s", zshrcPath))
	}

	// Configure aliases
	if len(cfg.Aliases) > 0 {
		ui.PrintStep("Configuring shell aliases...")
		if err := s.configureAliases(ctx, zshrcPath, cfg.Aliases); err != nil {
			return fmt.Errorf("failed to configure aliases: %w", err)
		}
	}

	// Configure environment variables
	if len(cfg.Environment) > 0 {
		ui.PrintStep("Configuring environment variables...")
		if err := s.configureEnvironment(ctx, zshrcPath, cfg.Environment); err != nil {
			return fmt.Errorf("failed to configure environment: %w", err)
		}
	}

	// Add zshrc extras
	if len(cfg.ZshrcExtras) > 0 {
		ui.PrintStep("Adding .zshrc extras...")
		if err := s.addExtras(ctx, zshrcPath, cfg.ZshrcExtras); err != nil {
			return fmt.Errorf("failed to add extras: %w", err)
		}
	}

	return nil
}

func (s *ShellInstaller) backupFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup.%s", path, timestamp)

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return err
	}

	ui.PrintInfo(fmt.Sprintf("Backed up %s to %s", path, backupPath))
	return nil
}

func (s *ShellInstaller) configureAliases(ctx context.Context, zshrcPath string, aliases map[string]string) error {
	if s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would configure %d aliases", len(aliases)))
		for name, cmd := range aliases {
			ui.PrintDryRun(fmt.Sprintf("  alias %s='%s'", name, cmd))
		}
		return nil
	}

	// Build aliases block
	var aliasLines []string
	aliasLines = append(aliasLines, "# Custom aliases (managed by setup-mac)")
	for name, cmd := range aliases {
		aliasLines = append(aliasLines, fmt.Sprintf("alias %s='%s'", name, cmd))
	}
	aliasLines = append(aliasLines, "# End custom aliases")

	aliasBlock := strings.Join(aliasLines, "\n")

	return s.updateZshrcBlock(zshrcPath, "# Custom aliases (managed by setup-mac)", "# End custom aliases", aliasBlock)
}

func (s *ShellInstaller) configureEnvironment(ctx context.Context, zshrcPath string, env map[string]string) error {
	if s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would configure %d environment variables", len(env)))
		for name, value := range env {
			ui.PrintDryRun(fmt.Sprintf("  export %s=\"%s\"", name, value))
		}
		return nil
	}

	// Build environment block
	var envLines []string
	envLines = append(envLines, "# Environment variables (managed by setup-mac)")
	for name, value := range env {
		envLines = append(envLines, fmt.Sprintf("export %s=\"%s\"", name, value))
	}
	envLines = append(envLines, "# End environment variables")

	envBlock := strings.Join(envLines, "\n")

	return s.updateZshrcBlock(zshrcPath, "# Environment variables (managed by setup-mac)", "# End environment variables", envBlock)
}

func (s *ShellInstaller) addExtras(ctx context.Context, zshrcPath string, extras []string) error {
	if s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would add %d extra lines to .zshrc", len(extras)))
		return nil
	}

	// Build extras block
	var extraLines []string
	extraLines = append(extraLines, "# Extra configuration (managed by setup-mac)")
	extraLines = append(extraLines, extras...)
	extraLines = append(extraLines, "# End extra configuration")

	extraBlock := strings.Join(extraLines, "\n")

	return s.updateZshrcBlock(zshrcPath, "# Extra configuration (managed by setup-mac)", "# End extra configuration", extraBlock)
}

func (s *ShellInstaller) updateZshrcBlock(zshrcPath, startMarker, endMarker, newBlock string) error {
	// Read current content
	content, err := os.ReadFile(zshrcPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new file with block
			return os.WriteFile(zshrcPath, []byte(newBlock+"\n"), 0644)
		}
		return err
	}

	contentStr := string(content)

	// Check if block already exists
	startIdx := strings.Index(contentStr, startMarker)
	endIdx := strings.Index(contentStr, endMarker)

	if startIdx != -1 && endIdx != -1 && endIdx > startIdx {
		// Replace existing block
		newContent := contentStr[:startIdx] + newBlock + contentStr[endIdx+len(endMarker):]
		return os.WriteFile(zshrcPath, []byte(newContent), 0644)
	}

	// Append new block
	if !strings.HasSuffix(contentStr, "\n") {
		contentStr += "\n"
	}
	contentStr += "\n" + newBlock + "\n"

	return os.WriteFile(zshrcPath, []byte(contentStr), 0644)
}
