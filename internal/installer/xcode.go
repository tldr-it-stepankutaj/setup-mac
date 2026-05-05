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

	// Capture stdout/stderr — `xcode-select --install` writes errors there
	// (e.g. "Can't install the software because it is not currently
	// available from the Software Update server"), and `RunInteractive`
	// would discard them so we couldn't surface anything actionable on
	// failure. The actual installer UI is a macOS-level dialog, not a TTY,
	// so we lose nothing by not piping stdin/stdout.
	result, runErr := x.ctx.Executor.Run(ctx, "xcode-select", "--install")
	if runErr != nil {
		// `xcode-select --install` exits non-zero in two harmless cases:
		// (1) CLT is already installed, (2) the dialog was successfully
		// triggered. Re-check before treating the exit as a real failure.
		if x.IsInstalled(ctx) {
			ui.PrintInfo("Xcode Command Line Tools already installed")
			return nil
		}
		time.Sleep(2 * time.Second)
		if x.IsInstalled(ctx) {
			ui.PrintInfo("Xcode Command Line Tools already installed")
			return nil
		}
		// Not installed and not in progress — diagnose so the user knows why.
		stderr := strings.TrimSpace(result.Stderr)
		if stderr != "" || result.ExitCode != 0 {
			x.diagnoseFailure(ctx, stderr)
			return fmt.Errorf("xcode-select --install failed: %s", firstLine(stderr))
		}
	}

	// Wait for installation to complete by polling
	ui.PrintInfo("Waiting for Xcode Command Line Tools installation to complete...")
	ui.PrintInfo("Please follow the installation dialog that appeared.")

	if err := x.waitForInstallation(ctx); err != nil {
		x.diagnoseFailure(ctx, "")
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

// diagnoseFailure prints actionable context after a failed CLT install. The
// most common cause we've seen is a macOS that's too old for the latest CLT
// build Apple's Software Update server is offering — `xcode-select --install`
// reports "Can't install the software because it is not currently available
// from the Software Update server" and exits without telling the user that
// the fix is to run a system update first. Surface enough info that the
// user can act on it without having to dig through Console.app.
func (x *XcodeInstaller) diagnoseFailure(ctx context.Context, stderr string) {
	fmt.Println()
	ui.PrintError("Xcode Command Line Tools installation did not complete.")

	if stderr != "" {
		fmt.Println()
		ui.PrintInfo("Output from xcode-select:")
		for _, line := range strings.Split(stderr, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			fmt.Printf("    %s\n", line)
		}
	}

	fmt.Println()
	ui.PrintInfo("System info:")
	if v := x.captureOutput(ctx, "sw_vers", "-productName"); v != "" {
		fmt.Printf("    macOS:        %s\n", v)
	}
	if v := x.captureOutput(ctx, "sw_vers", "-productVersion"); v != "" {
		fmt.Printf("    Version:      %s\n", v)
	}
	if v := x.captureOutput(ctx, "sw_vers", "-buildVersion"); v != "" {
		fmt.Printf("    Build:        %s\n", v)
	}

	// Stderr containing this phrase is the canonical "macOS is too old for
	// the CLT build Apple is offering today" failure mode.
	tooOld := strings.Contains(strings.ToLower(stderr),
		"not currently available from the software update server")

	fmt.Println()
	ui.PrintInfo("Common causes and fixes:")
	if tooOld {
		fmt.Println("    • Apple is no longer publishing a Command Line Tools build for this")
		fmt.Println("      macOS version. Update macOS first, then re-run setup-mac.")
		fmt.Println("        System Settings → General → Software Update")
	}
	fmt.Println("    • If a system update is pending, install it first:")
	fmt.Println("        softwareupdate -l           # list available updates")
	fmt.Println("        sudo softwareupdate -ia     # install all available updates")
	fmt.Println("    • As a fallback, install the full Xcode from the App Store:")
	fmt.Println("        https://apps.apple.com/app/xcode/id497799835")
	fmt.Println("    • To install CLT manually after the OS is up to date:")
	fmt.Println("        xcode-select --install")
	fmt.Println()
}

// captureOutput runs a command and returns its trimmed stdout, or "" on error.
// Used for diagnostics where a failure to capture is itself non-fatal.
func (x *XcodeInstaller) captureOutput(ctx context.Context, name string, args ...string) string {
	result, err := x.ctx.Executor.Run(ctx, name, args...)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(result.Stdout)
}

// firstLine returns the first non-empty line of s, useful for boiling down a
// multi-line stderr into a single-line error message.
func firstLine(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			return line
		}
	}
	return s
}
