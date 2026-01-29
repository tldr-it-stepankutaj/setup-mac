# setup-mac

[![CI](https://github.com/stepankutaj/setup-mac/actions/workflows/ci.yml/badge.svg)](https://github.com/stepankutaj/setup-mac/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/stepankutaj/setup-mac)](https://github.com/stepankutaj/setup-mac/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/stepankutaj/setup-mac)](https://github.com/stepankutaj/setup-mac)
[![License](https://img.shields.io/github/license/stepankutaj/setup-mac)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/stepankutaj/setup-mac)](https://goreportcard.com/report/github.com/stepankutaj/setup-mac)

CLI tool for automating macOS developer environment setup.

## Features

- **Homebrew** - Install Homebrew, formulae, casks, and taps
- **Oh-My-Zsh** - Install with plugins (zsh-autosuggestions, zsh-syntax-highlighting)
- **Powerlevel10k** - Theme with interactive style selection
- **Shell Config** - Aliases, environment variables, .zshrc customization
- **macOS Defaults** - Dock, Finder, Keyboard settings
- **Git Config** - User info, aliases, settings
- **SSH Key** - Generate ed25519 key

## Installation

### From Release

Download the latest release from [Releases](https://github.com/stepankutaj/setup-mac/releases).

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
- macOS (tested on Sonoma 14.x)
- Xcode Command Line Tools (`xcode-select --install`)

```bash
git clone https://github.com/stepankutaj/setup-mac.git
cd setup-mac
make build
make install
```

## Usage

### Install Everything

```bash
setup-mac install --all
```

### Install Specific Components

```bash
setup-mac install --homebrew    # Homebrew and packages
setup-mac install --terminal    # Oh-My-Zsh + Powerlevel10k
setup-mac install --shell       # Shell aliases and environment
setup-mac install --macos       # macOS defaults
setup-mac install --git         # Git configuration
setup-mac install --ssh         # SSH key generation
```

### Dry-Run Mode

Preview changes without executing:

```bash
setup-mac install --all --dry-run
```

### Custom Configuration

```bash
setup-mac install --all --config my-config.yaml
```

### Version

```bash
setup-mac version
```

## Configuration

Configuration is in YAML format. See `configs/default.yaml` for full schema.

### Example

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
    - eza
  casks:
    - iterm2
    - visual-studio-code
    - docker

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
    ll: "eza -la --icons"
    gs: "git status"
    k: "kubectl"
  environment:
    EDITOR: "code --wait"

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
  user:
    name: "Your Name"
    email: "your@email.com"

ssh:
  generate_key: true
  key_type: "ed25519"
  key_file: "~/.ssh/id_ed25519"
```

## Development

### Setup

```bash
git clone https://github.com/stepankutaj/setup-mac.git
cd setup-mac
make deps
make install-tools
```

### Build

```bash
make build
```

### Test

```bash
make test
```

### Lint

```bash
make lint
```

### Run Without Installing

```bash
make build
./bin/setup-mac install --all --dry-run
```

### Project Structure

```
setup-mac/
├── cmd/setup-mac/main.go       # Entry point
├── internal/
│   ├── cli/                    # Cobra commands
│   ├── config/                 # Configuration loading
│   ├── installer/              # Component installers
│   ├── executor/               # Command execution
│   └── ui/                     # Spinners and prompts
├── configs/
│   └── default.yaml            # Default configuration
├── go.mod
├── Makefile
└── README.md
```

### Available Make Targets

```
make help
```

| Target | Description |
|--------|-------------|
| `all` | Download dependencies and build |
| `build` | Build the binary |
| `build-release` | Build for amd64 and arm64 |
| `test` | Run tests |
| `lint` | Run linter |
| `fmt` | Format code |
| `deps` | Download dependencies |
| `install` | Install to /usr/local/bin |
| `uninstall` | Remove from /usr/local/bin |
| `install-tools` | Install golangci-lint |
| `dry-run` | Build and run with --dry-run |

## Uninstall

```bash
make uninstall
# Or
sudo rm /usr/local/bin/setup-mac
```

## Author

Stepan Kutaj <stepan.kutaj@tldr-it.com>

## License

Apache-2.0 [License](LICENSE)
