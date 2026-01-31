package installer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tldr-it-stepankutaj/setup-mac/internal/ui"
)

const (
	xcodeSelectPath = "/usr/bin/xcode-select"
)

// XcodeInstaller handles Xcode Command Line Tools installation
type XcodeInstaller struct {
	ctx *Context
}

// NewXcodeInstaller creates a new Xcode CLT installer
func NewXcodeInstaller(ctx *Context) *XcodeInstaller {
	return &XcodeInstaller{ctx: ctx}
}

// Name returns the installer name
func (x *XcodeInstaller) Name() string {
	return "xcode"
}

// Description returns the installer description
func (x *XcodeInstaller) Description() string {
	return "Xcode Command Line Tools"
}

// IsInstalled checks if Xcode CLT is installed
func (x *XcodeInstaller) IsInstalled(ctx context.Context) bool {
	// In dry-run mode, assume not installed to show what would happen
	if x.ctx.DryRun {
		return false
	}

	// Check if xcode-select exists
	if _, err := os.Stat(xcodeSelectPath); os.IsNotExist(err) {
		return false
	}

	// Check if CLT is properly installed by checking the path
	result, err := x.ctx.Executor.Run(ctx, "xcode-select", "-p")
	if err != nil {
		return false
	}

	// If path is empty, not installed
	path := strings.TrimSpace(result.Stdout)
	if path == "" {
		return false
	}

	// Verify the path actually exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Install installs Xcode Command Line Tools
func (x *XcodeInstaller) Install(ctx context.Context) error {
	if x.IsInstalled(ctx) {
		ui.PrintInfo("Xcode Command Line Tools already installed")
		return nil
	}

	ui.PrintStep("Installing Xcode Command Line Tools...")
	ui.PrintInfo("This may take a while and will show a system dialog...")

	if x.ctx.DryRun {
		ui.PrintDryRun("xcode-select --install")
		return nil
	}

	// Start the installation - this triggers a macOS dialog
	err := x.ctx.Executor.RunInteractive(ctx, "xcode-select", "--install")
	if err != nil {
		// Check if it's because CLT is already installed (exit code 1 with specific message)
		if x.IsInstalled(ctx) {
			ui.PrintInfo("Xcode Command Line Tools already installed")
			return nil
		}
		// The command might "fail" but still show the dialog - check again after a moment
		time.Sleep(2 * time.Second)
		if x.IsInstalled(ctx) {
			ui.PrintInfo("Xcode Command Line Tools already installed")
			return nil
		}
	}

	// Wait for installation to complete by polling
	ui.PrintInfo("Waiting for Xcode Command Line Tools installation to complete...")
	ui.PrintInfo("Please follow the installation dialog that appeared.")

	if err := x.waitForInstallation(ctx); err != nil {
		return err
	}

	ui.PrintSuccess("Xcode Command Line Tools installed successfully")
	return nil
}

// waitForInstallation polls until CLT is installed or context is cancelled
func (x *XcodeInstaller) waitForInstallation(ctx context.Context) error {
	spinner := ui.NewSpinner("Waiting for installation to complete (press Ctrl+C to skip)...")
	spinner.Start()
	defer spinner.Stop()

	// Poll every 5 seconds for up to 30 minutes
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(30 * time.Minute)
	defer timeout.Stop()

	for {
		select {
		case <-ctx.Done():
			spinner.Fail("Installation cancelled")
			return ctx.Err()
		case <-timeout.C:
			spinner.Fail("Installation timed out")
			return fmt.Errorf("xcode CLT installation timed out after 30 minutes")
		case <-ticker.C:
			if x.IsInstalled(ctx) {
				spinner.Success("Xcode Command Line Tools installed")
				return nil
			}
		}
	}
}
