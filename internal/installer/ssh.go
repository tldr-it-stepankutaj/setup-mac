package installer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

// SSHInstaller handles SSH key generation
type SSHInstaller struct {
	ctx *Context
}

// NewSSHInstaller creates a new SSH installer
func NewSSHInstaller(ctx *Context) *SSHInstaller {
	return &SSHInstaller{ctx: ctx}
}

// Name returns the installer name
func (s *SSHInstaller) Name() string {
	return "ssh"
}

// Description returns the installer description
func (s *SSHInstaller) Description() string {
	return "SSH Key Generation"
}

// IsInstalled checks if SSH key exists
func (s *SSHInstaller) IsInstalled(ctx context.Context) bool {
	keyFile := s.expandKeyPath(s.ctx.Config.SSH.KeyFile)
	_, err := os.Stat(keyFile)
	return err == nil
}

// Install generates an SSH key
func (s *SSHInstaller) Install(ctx context.Context) error {
	cfg := s.ctx.Config.SSH

	if !cfg.GenerateKey {
		ui.PrintInfo("SSH key generation skipped (disabled in config)")
		return nil
	}

	keyFile := s.expandKeyPath(cfg.KeyFile)

	// Check if key already exists
	if _, err := os.Stat(keyFile); err == nil {
		ui.PrintInfo(fmt.Sprintf("SSH key already exists: %s", keyFile))
		return nil
	}

	// Get comment from config or prompt (skip prompting in dry-run mode)
	comment := cfg.Comment
	if comment == "" && s.ctx.Config.Settings.Interactive && !s.ctx.DryRun {
		var err error
		comment, err = s.ctx.Prompt.Input("SSH key comment (e.g., your email)", "")
		if err != nil {
			comment = ""
		}
	}

	// Ensure .ssh directory exists
	sshDir := filepath.Dir(keyFile)
	if s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("Would create directory: %s", sshDir))
	} else {
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return fmt.Errorf("failed to create .ssh directory: %w", err)
		}
	}

	// Generate SSH key
	ui.PrintStep(fmt.Sprintf("Generating %s SSH key...", cfg.KeyType))

	args := []string{
		"-t", cfg.KeyType,
		"-f", keyFile,
		"-N", "", // Empty passphrase (user can change later)
	}

	if comment != "" {
		args = append(args, "-C", comment)
	}

	if s.ctx.DryRun {
		ui.PrintDryRun(fmt.Sprintf("ssh-keygen %s", strings.Join(args, " ")))
		return nil
	}

	result, err := s.ctx.Executor.Run(ctx, "ssh-keygen", args...)
	if err != nil {
		return fmt.Errorf("failed to generate SSH key: %w", err)
	}

	if result.ExitCode == 0 {
		ui.PrintSuccess(fmt.Sprintf("SSH key generated: %s", keyFile))

		// Display public key
		pubKeyFile := keyFile + ".pub"
		pubKey, err := os.ReadFile(pubKeyFile)
		if err == nil {
			ui.PrintInfo("Public key:")
			fmt.Println(string(pubKey))
		}

		// Add to ssh-agent
		ui.PrintStep("Adding key to ssh-agent...")
		if err := s.addToAgent(ctx, keyFile); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to add key to ssh-agent: %v", err))
		}
	}

	return nil
}

func (s *SSHInstaller) expandKeyPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func (s *SSHInstaller) addToAgent(ctx context.Context, keyFile string) error {
	// Start ssh-agent if not running
	_, _ = s.ctx.Executor.RunShell(ctx, "eval $(ssh-agent -s)")

	// Add key to agent
	result, err := s.ctx.Executor.Run(ctx, "ssh-add", keyFile)
	if err != nil {
		return err
	}

	if result.ExitCode == 0 {
		ui.PrintSuccess("Key added to ssh-agent")
	}

	return nil
}
