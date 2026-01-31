# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-01-31

### Added
- **Status command** - `setup-mac status` shows installation status of all components
  - JSON output support with `--json` flag for scripting
  - Shows system info (OS, architecture, Apple Silicon, macOS version)
- **Update command** - `setup-mac update` updates installed tools
  - `--homebrew` - runs brew update, upgrade, and cleanup
  - `--ohmyzsh` - updates Oh My Zsh and custom plugins/themes
- **Validate command** - `setup-mac validate` validates configuration files
  - Shows configuration summary
  - Reports errors and warnings
- **Xcode Command Line Tools installer** - automatically installs if missing
- **Rosetta 2 installer** - installs on Apple Silicon Macs for x86 compatibility
- **Network connectivity check** - verifies internet before starting installation
- **Version check** - automatically checks GitHub for new releases
  - Shows download URL for current platform
  - `--skip-update-check` flag to disable
- **Progress indication** - shows `[1/9]`, `[2/9]`, etc. during installation
- **Root/sudo detection** - refuses to run as root to prevent permission issues

### Changed
- **Git installer** - now prompts interactively for user.name and user.email
  - Uses existing git config values as defaults
  - Shows warning if not set
- **Shell aliases** - replaced `eza` with standard `ls` commands
  - `ll` = `ls -la`
  - `la` = `ls -a`
  - `lt` = `tree`
- **Homebrew installer** - improved detection of installed packages
  - Handles versioned packages (e.g., `node@18`)
  - Checks `/Applications` for manually installed apps
  - Better error handling for "already installed" cases

### Fixed
- Interactive prompts now display correctly (removed spinner interference)
- Fixed newline in macOS version output in status command

### Removed
- `eza` from default formulae (users can add it in custom config)

## [1.0.0] - 2026-01-29

### Added
- Initial release
- **Homebrew installer** - install Homebrew, formulae, casks, and taps
- **Oh-My-Zsh installer** - install with custom plugins
  - zsh-autosuggestions
  - zsh-syntax-highlighting
- **Powerlevel10k installer** - theme installation
- **Shell configuration** - aliases, environment variables, .zshrc customization
- **macOS defaults** - Dock, Finder, Keyboard settings
- **Git configuration** - aliases and settings
- **SSH key generation** - ed25519 key generation
- **Dry-run mode** - preview changes without executing
- **Custom configuration** - YAML-based configuration
- **Interactive mode** - confirmation prompts
- **Backup dotfiles** - automatic backup before modification
- CLI with Cobra framework
- Embedded default configuration
- Colorful terminal output with spinners

[Unreleased]: https://github.com/tldr-it-stepankutaj/setup-mac/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/tldr-it-stepankutaj/setup-mac/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/tldr-it-stepankutaj/setup-mac/releases/tag/v1.0.0
