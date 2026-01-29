package installer

import (
	"context"
	"fmt"

	"github.com/stepankutaj/setup-mac/internal/config"
	"github.com/stepankutaj/setup-mac/internal/executor"
	"github.com/stepankutaj/setup-mac/internal/ui"
)

// Installer defines the interface for all installers
type Installer interface {
	// Name returns the installer name
	Name() string

	// Description returns a short description
	Description() string

	// IsInstalled checks if the component is already installed
	IsInstalled(ctx context.Context) bool

	// Install performs the installation
	Install(ctx context.Context) error
}

// Context provides shared context for installers
type Context struct {
	Config   *config.Config
	Executor *executor.Executor
	Prompt   *ui.Prompt
	DryRun   bool
	Verbose  bool
}

// NewContext creates a new installer context
func NewContext(cfg *config.Config, dryRun, verbose bool) *Context {
	return &Context{
		Config:   cfg,
		Executor: executor.New(dryRun, verbose),
		Prompt:   ui.NewPrompt(cfg.Settings.Interactive),
		DryRun:   dryRun,
		Verbose:  verbose,
	}
}

// Registry holds all available installers
type Registry struct {
	installers map[string]func(*Context) Installer
}

// NewRegistry creates a new installer registry
func NewRegistry() *Registry {
	return &Registry{
		installers: make(map[string]func(*Context) Installer),
	}
}

// Register adds an installer factory to the registry
func (r *Registry) Register(name string, factory func(*Context) Installer) {
	r.installers[name] = factory
}

// Get returns an installer by name
func (r *Registry) Get(name string, ctx *Context) (Installer, error) {
	factory, ok := r.installers[name]
	if !ok {
		return nil, fmt.Errorf("unknown installer: %s", name)
	}
	return factory(ctx), nil
}

// GetAll returns all registered installers
func (r *Registry) GetAll(ctx *Context) []Installer {
	var installers []Installer
	for _, factory := range r.installers {
		installers = append(installers, factory(ctx))
	}
	return installers
}

// Names returns all registered installer names
func (r *Registry) Names() []string {
	var names []string
	for name := range r.installers {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the global installer registry
var DefaultRegistry = NewRegistry()

func init() {
	// Register all installers
	DefaultRegistry.Register("homebrew", func(ctx *Context) Installer {
		return NewHomebrewInstaller(ctx)
	})
	DefaultRegistry.Register("ohmyzsh", func(ctx *Context) Installer {
		return NewOhMyZshInstaller(ctx)
	})
	DefaultRegistry.Register("powerlevel10k", func(ctx *Context) Installer {
		return NewPowerlevel10kInstaller(ctx)
	})
	DefaultRegistry.Register("shell", func(ctx *Context) Installer {
		return NewShellInstaller(ctx)
	})
	DefaultRegistry.Register("macos", func(ctx *Context) Installer {
		return NewMacOSInstaller(ctx)
	})
	DefaultRegistry.Register("git", func(ctx *Context) Installer {
		return NewGitInstaller(ctx)
	})
	DefaultRegistry.Register("ssh", func(ctx *Context) Installer {
		return NewSSHInstaller(ctx)
	})
}

// RunInstaller runs a single installer
func RunInstaller(ctx context.Context, installer Installer, ictx *Context) error {
	ui.PrintHeader(installer.Description())

	if installer.IsInstalled(ctx) {
		ui.PrintInfo(fmt.Sprintf("%s is already installed", installer.Name()))
		return nil
	}

	spinner := ui.NewSpinner(fmt.Sprintf("Installing %s...", installer.Name()))
	spinner.Start()

	err := installer.Install(ctx)
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to install %s: %v", installer.Name(), err))
		return err
	}

	spinner.Success(fmt.Sprintf("%s installed successfully", installer.Name()))
	return nil
}
