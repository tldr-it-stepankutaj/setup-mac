package installer

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

const (
	rosettaPath = "/Library/Apple/usr/share/rosetta/rosetta"
)

// RosettaInstaller handles Rosetta 2 installation on Apple Silicon
type RosettaInstaller struct {
	ctx *Context
}

// NewRosettaInstaller creates a new Rosetta 2 installer
func NewRosettaInstaller(ctx *Context) *RosettaInstaller {
	return &RosettaInstaller{ctx: ctx}
}

// Name returns the installer name
func (r *RosettaInstaller) Name() string {
	return "rosetta"
}

// Description returns the installer description
func (r *RosettaInstaller) Description() string {
	return "Rosetta 2 (x86 compatibility for Apple Silicon)"
}

// IsAppleSilicon returns true if running on Apple Silicon
func (r *RosettaInstaller) IsAppleSilicon() bool {
	return runtime.GOARCH == "arm64"
}

// IsInstalled checks if Rosetta 2 is installed
func (r *RosettaInstaller) IsInstalled(ctx context.Context) bool {
	// Not applicable on Intel Macs
	if !r.IsAppleSilicon() {
		return true // Consider it "installed" on Intel (not needed)
	}

	// Check if Rosetta binary exists
	if _, err := os.Stat(rosettaPath); err == nil {
		return true
	}

	// Alternative check using arch command
	result, err := r.ctx.Executor.Run(ctx, "arch", "-x86_64", "true")
	if err == nil && result.ExitCode == 0 {
		return true
	}

	return false
}

// Install installs Rosetta 2
func (r *RosettaInstaller) Install(ctx context.Context) error {
	// Skip on Intel Macs
	if !r.IsAppleSilicon() {
		ui.PrintInfo("Rosetta 2 not needed (not Apple Silicon)")
		return nil
	}

	if r.IsInstalled(ctx) {
		ui.PrintInfo("Rosetta 2 already installed")
		return nil
	}

	ui.PrintStep("Installing Rosetta 2...")

	if r.ctx.DryRun {
		ui.PrintDryRun("softwareupdate --install-rosetta --agree-to-license")
		return nil
	}

	// Install Rosetta 2 using softwareupdate
	// The --agree-to-license flag accepts the license automatically
	result, err := r.ctx.Executor.Run(ctx, "softwareupdate", "--install-rosetta", "--agree-to-license")
	if err != nil {
		// Check if it's already installed despite the error
		if r.IsInstalled(ctx) {
			ui.PrintInfo("Rosetta 2 already installed")
			return nil
		}
		return fmt.Errorf("failed to install Rosetta 2: %w\nOutput: %s", err, result.Stderr)
	}

	ui.PrintSuccess("Rosetta 2 installed successfully")
	return nil
}
