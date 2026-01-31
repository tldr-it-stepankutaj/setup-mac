# setup-mac

[![CI](https://github.com/tldr-it-stepankutaj/setup-mac/actions/workflows/ci.yml/badge.svg)](https://github.com/tldr-it-stepankutaj/setup-mac/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/tldr-it-stepankutaj/setup-mac)](https://github.com/tldr-it-stepankutaj/setup-mac/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tldr-it-stepankutaj/setup-mac)](https://github.com/tldr-it-stepankutaj/setup-mac)
[![License](https://img.shields.io/github/license/tldr-it-stepankutaj/setup-mac)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tldr-it-stepankutaj/setup-mac)](https://goreportcard.com/report/github.com/tldr-it-stepankutaj/setup-mac)

CLI tool for automating macOS developer environment setup.

## Features

- **Xcode Command Line Tools** - Automatically installs if missing
- **Rosetta 2** - Installs on Apple Silicon Macs for x86 compatibility
- **Homebrew** - Install Homebrew, formulae, casks, and taps
- **Oh-My-Zsh** - Install with plugins (zsh-autosuggestions, zsh-syntax-highlighting)
- **Powerlevel10k** - Theme with interactive style selection
- **Shell Config** - Aliases, environment variables, .zshrc customization
- **macOS Defaults** - Dock, Finder, Keyboard settings
- **Git Config** - User info (interactive prompt), aliases, settings
- **SSH Key** - Generate ed25519 key
- **Network Check** - Verifies internet connectivity before installation
- **Progress Indication** - Shows [1/9], [2/9], etc. during installation
- **Auto-Update Check** - Notifies when new version is available on GitHub

## Installation

### From Release

Download the latest release from [Releases](https://github.com/tldr-it-stepankutaj/setup-mac/releases).

```bash
# Extract
tar -xzf setup-mac-darwin-arm64.tar.gz
cd setup-mac

# Install
make install

# Or manually copy binary
sudo cp bin/setup-mac /usr/local/bin/
```

### From Source

Requirements:
- Go 1.24+
- macOS (tested on Sonoma 14.x+)

```bash
git clone https://github.com/tldr-it-stepankutaj/setup-mac.git
cd setup-mac
make build
make install
```

## Usage

### Quick Start

```bash
# Check what's already installed
setup-mac status

# Preview what would be installed (safe, no changes)
setup-mac install --all --dry-run

# Install everything
setup-mac install --all
```

### Commands

| Command | Description |
|---------|-------------|
| `install` | Install and configure development tools |
| `status` | Show installation status of all components |
| `update` | Update installed tools (Homebrew, Oh-My-Zsh) |
| `validate` | Validate configuration file |
| `version` | Print version information |

### Install Options

```bash
setup-mac install --all         # Install everything
setup-mac install --xcode       # Xcode Command Line Tools
setup-mac install --rosetta     # Rosetta 2 (Apple Silicon only)
setup-mac install --homebrew    # Homebrew and packages
setup-mac install --terminal    # Oh-My-Zsh + Powerlevel10k
setup-mac install --shell       # Shell aliases and environment
setup-mac install --macos       # macOS defaults
setup-mac install --git         # Git configuration (prompts for name/email)
setup-mac install --ssh         # SSH key generation
```

### Update Installed Tools

```bash
setup-mac update --all          # Update everything
setup-mac update --homebrew     # brew update && brew upgrade
setup-mac update --ohmyzsh      # Update Oh-My-Zsh and plugins
```

### Dry-Run Mode

Preview changes without executing:

```bash
setup-mac install --all --dry-run
setup-mac update --all --dry-run
```

### Custom Configuration

```bash
setup-mac install --all --config my-config.yaml
setup-mac validate --config my-config.yaml
```

### Auto-Update Check

The tool automatically checks for new versions on GitHub when you run any command.

```
⬆ New version available: v1.1.0 (current: v1.0.0)
  Download: https://github.com/.../setup-mac-darwin-arm64.tar.gz
  Run with --skip-update-check to disable this message
```

To disable the update check:

```bash
setup-mac status --skip-update-check
```

### Global Flags

| Flag | Description |
|------|-------------|
| `-c, --config` | Custom config file path |
| `-v, --verbose` | Verbose output |
| `--skip-update-check` | Skip checking for new versions |

## Example Output

### Status Command

```
$ setup-mac status

System Information
──────────────────────────────────────
  OS:            darwin
  Architecture:  arm64
  Apple Silicon: Yes
  macOS Version: 14.5

Components
──────────────────────────────────────
  ✓ xcode                Xcode Command Line Tools
  ✓ rosetta              Rosetta 2 (x86 compatibility for Apple Silicon)
  ✓ homebrew             Homebrew Package Manager
  ✓ oh-my-zsh            Oh My Zsh Framework
  ✓ powerlevel10k        Powerlevel10k Theme
  ✗ shell                Shell Configuration
  ✗ macos                macOS System Defaults
  ✗ git                  Git Configuration
  ✓ ssh                  SSH Key Generation

  6/9 components installed
```

### Validate Command

```
$ setup-mac validate

Validating embedded default config

Configuration Summary
──────────────────────────────────────
  Homebrew:      enabled (11 formulae, 5 casks, 1 taps)
  Oh My Zsh:     enabled (6 plugins)
  Powerlevel10k: enabled
  Git:           enabled (user: )
  SSH:           enabled (type: ed25519)
  Shell:         13 aliases, 4 env vars
  macOS:         enabled

Warnings
──────────────────────────────────────
  ⚠ Git user.name is not set
  ⚠ Git user.email is not set

✓ Configuration is valid
```

### Install with Progress

```
$ setup-mac install --all --dry-run

ℹ Installing 9 component(s):
  - Xcode Command Line Tools
  - Rosetta 2 (x86 compatibility for Apple Silicon)
  - Homebrew Package Manager
  - Oh My Zsh Framework
  - Powerlevel10k Theme
  - Shell Configuration
  - macOS System Defaults
  - Git Configuration
  - SSH Key Generation

=== DRY-RUN MODE ===

═══════════════════════════════════════
  [1/9] Xcode Command Line Tools
═══════════════════════════════════════

→ Installing Xcode Command Line Tools...
ℹ This may take a while and will show a system dialog...
[DRY-RUN] xcode-select --install
✓ xcode installed successfully

═══════════════════════════════════════
  [2/9] Rosetta 2 (x86 compatibility for Apple Silicon)
═══════════════════════════════════════

ℹ rosetta is already installed

═══════════════════════════════════════
  [3/9] Homebrew Package Manager
═══════════════════════════════════════

ℹ homebrew is already installed
...
```

### Update Command

```
$ setup-mac update --all --dry-run

ℹ Updating 2 component(s):
  - Homebrew Package Manager
  - Oh My Zsh Framework

=== DRY-RUN MODE ===

[1/2] Homebrew Package Manager
──────────────────────────────────────
→ Updating Homebrew...
[DRY-RUN] brew update
→ Upgrading packages...
[DRY-RUN] brew upgrade
→ Upgrading casks...
[DRY-RUN] brew upgrade --cask
→ Cleaning up...
[DRY-RUN] brew cleanup

[2/2] Oh My Zsh Framework
──────────────────────────────────────
→ Updating Oh My Zsh...
[DRY-RUN] cd ~/.oh-my-zsh && git pull

Update completed successfully!
```

### JSON Output (for scripting)

```bash
$ setup-mac status --json
{
  "system": {
    "os": "darwin",
    "arch": "arm64",
    "apple_silicon": true,
    "macos_version": "14.5"
  },
  "components": [
    {"name": "xcode", "description": "Xcode Command Line Tools", "installed": true},
    {"name": "rosetta", "description": "Rosetta 2", "installed": true},
    {"name": "homebrew", "description": "Homebrew Package Manager", "installed": true},
    ...
  ]
}
```

## Configuration

Configuration is in YAML format. The tool includes sensible defaults, but you can customize everything.

### Example Configuration

```yaml
version: "1.0"

settings:
  dry_run: false
  interactive: true
  backup_dotfiles: true

homebrew:
  install: true
  taps:
    - homebrew/cask-fonts
  formulae:
    - git
    - gh
    - fzf
    - ripgrep
    - bat
    - jq
    - htop
    - tree
  casks:
    - iterm2
    - visual-studio-code
    - docker
    - rectangle

terminal:
  oh_my_zsh:
    install: true
    plugins:
      - git
      - docker
      - fzf
      - zsh-autosuggestions
      - zsh-syntax-highlighting
  powerlevel10k:
    install: true

shell:
  aliases:
    ll: "ls -la"
    la: "ls -a"
    lt: "tree"
    cat: "bat"
    gs: "git status"
    k: "kubectl"
  environment:
    EDITOR: "code --wait"
    LANG: "en_US.UTF-8"

macos:
  configure: true
  defaults:
    dock:
      autohide: true
      tile_size: 48
    finder:
      show_hidden_files: true
      show_extensions: true
    keyboard:
      key_repeat: 2
      disable_smart_quotes: true

git:
  configure: true
  # user.name and user.email will be prompted interactively
  aliases:
    st: "status"
    co: "checkout"
    lg: "log --oneline --graph --decorate"
  settings:
    init.defaultBranch: "main"
    pull.rebase: "true"

ssh:
  generate_key: true
  key_type: "ed25519"
```

## Safety Features

- **Root/Sudo Detection** - Refuses to run as root to prevent permission issues
- **Network Check** - Verifies connectivity before starting installations
- **Dry-Run Mode** - Preview all changes before applying
- **Backup** - Automatically backs up dotfiles before modification
- **Idempotent** - Safe to run multiple times, skips already installed components

## Testing Safely

```bash
# Always start with dry-run
setup-mac install --all --dry-run

# Check current status
setup-mac status

# Validate your config
setup-mac validate --config my-config.yaml

# Test on a separate user account (recommended for full test)
# See documentation for creating a test user
```

## Development

### Setup

```bash
git clone https://github.com/tldr-it-stepankutaj/setup-mac.git
cd setup-mac
make deps
make install-tools
```

### Build & Test

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linter
make dry-run    # Build and run with --dry-run
```

### Project Structure

```
setup-mac/
├── cmd/setup-mac/main.go       # Entry point
├── internal/
│   ├── cli/                    # Cobra commands (install, status, update, validate)
│   ├── config/                 # Configuration loading and schema
│   ├── installer/              # Component installers
│   ├── executor/               # Command execution with dry-run support
│   └── ui/                     # Spinners, prompts, and output formatting
├── configs/
│   └── default.yaml            # Default configuration
├── go.mod
├── Makefile
└── README.md
```

### Make Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the binary |
| `make test` | Run tests |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make deps` | Download dependencies |
| `make install` | Install to /usr/local/bin |
| `make dry-run` | Build and run with --dry-run |
| `make dist` | Create distribution packages |

## Uninstall

```bash
make uninstall
# Or
sudo rm /usr/local/bin/setup-mac
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes.

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting a pull request.

## Author

Stepan Kutaj <stepan.kutaj@tldr-it.com>

## License

Apache-2.0 [License](LICENSE)
