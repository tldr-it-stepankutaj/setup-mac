# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3] - 2026-05-05

### Added
- **Auto-discovery of the user config.** Without `--config`, setup-mac now looks for `~/.config/setup-mac/config.yaml` (canonical name written by `setup-mac config init` and `make install`) and falls back to `~/.config/setup-mac/default.yaml` (the legacy name shipped by older release tarballs). The path actually loaded is printed as `ℹ Using config: …` in `install`/`update` and surfaced in `validate` output. Tests use `LoadDefault()` (embedded only) so a developer's real config doesn't leak into test results.
- **Diagnostics on Xcode CLT failure.** `xcode-select --install` now runs through `Run` (not `RunInteractive`) so stderr is captured. On failure or 30-minute timeout the installer prints macOS version (`sw_vers`), the verbatim `xcode-select` error, and actionable fixes — `softwareupdate -l` / `sudo softwareupdate -ia`, link to full Xcode in the App Store, and the manual retry command. The "not currently available from the Software Update server" failure (macOS too old for Apple's current CLT build) is detected and called out explicitly.

### Changed
- **Dist Makefile (`make dist` and `.github/workflows/build.yml`) writes `~/.config/setup-mac/config.yaml`** instead of `default.yaml`, and is now idempotent — if the file already exists the install target prints `Keeping existing …` and does not overwrite it. Previously every `make install` from a fresh tarball clobbered any edits the user had made.
- **Skip-on-disabled no longer lies about success.** Installers that bail out at the top because of a config flag (`Install: false`, `Configure: false`, `GenerateKey: false`) now return `installer.ErrSkipped`; the runner suppresses the misleading `✓ <name> installed successfully` line. Affects `git`, `homebrew`, `oh-my-zsh`, `powerlevel10k`, `macos`, `ssh`.

### Fixed
- **User config under `~/.config/setup-mac/` was being ignored.** Without an explicit `--config` flag the tool fell back to embedded defaults (which have `git.configure: false`, `ssh.generate_key: false`, etc.), so e.g. `setup-mac install --git` against a user config with `git.configure: true` printed "Git configuration skipped (disabled in config)" — and then falsely "git installed successfully" right after. Auto-discovery (see Added) plus `ErrSkipped` (see Changed) fix both halves of the bug.

## [1.0.2] - 2026-05-04

### Added
- **`setup-mac config init`** command — writes the embedded default config to `~/.config/setup-mac/config.yaml` (or any path with `--output`); refuses to overwrite without `--force`. Gives users a starting point to edit instead of having to copy `configs/default.yaml` from the repo.
- **Powerlevel10k style is now actually applied.** When the interactive style selection (`lean` / `classic` / `rainbow` / `pure`) returns, the corresponding `p10k-{style}.zsh` is copied to `~/.p10k.zsh` and a `source ~/.p10k.zsh` line is appended to `.zshrc` if missing. Existing `~/.p10k.zsh` is **never** overwritten — that's your hand-tuned config from `p10k configure`.
- 6 additional default Oh-My-Zsh plugins: `colored-man-pages`, `colorize`, `brew`, `ssh-agent`, `pip`, `python`.
- `.golangci.yml` for consistent linting; `make lint` now passes with 0 issues.
- `.github/dependabot.yml` for weekly Go module + GitHub Actions bumps.
- Unit tests for semver comparison (`versioncheck_test.go`) and plugin merge logic (`ohmyzsh_test.go`).
- `CLAUDE.md` developer guide covering architecture, safety invariants, and prompt/spinner rules.

### Changed
- **Plugin configuration is now MERGED, not replaced.** `setup-mac install --terminal` previously overwrote the entire `plugins=(…)` line in `.zshrc`, silently wiping plugins the user had added by hand. Now it parses the existing line, unions with the configured list, and only writes if anything actually changed.
- **Single source of truth for defaults.** The embedded default config now lives at `configs/default.yaml` only (was duplicated with `internal/config/defaults.yaml`, which had drifted out of sync). The `dist/` tarball and the embedded binary now ship identical defaults.
- **Embedded defaults match the `configs/default.yaml` shipped in `dist/`**: `eza`, `yq`, `wget`, `curl` are back in formulae and the default `ll`/`la`/`lt` aliases now use `eza --icons`.
- **`git` IsInstalled** now correctly reports installed in `setup-mac status` when global `user.name` AND `user.email` are configured (was hardcoded to `false` — git always showed `✗`).
- **Per-item failures propagate.** `configureAliases` / `configureSettings` (git) and `applyDefaults` (macos) now accumulate per-item errors via `errors.Join` and return them. The install summary no longer reports "completed successfully" when half the config silently failed.

### Fixed
- **Confirm prompt accepted Enter as "no".** `promptui.IsConfirm` ignored the `Default` field and treated Enter as `ErrAbort`, so `Proceed with installation? [Y/n]` + Enter aborted the install. Replaced with a `bufio` reader: Enter respects the default, invalid input re-prompts, and stray characters don't leak into the next prompt.
- **Version comparison was lexicographic.** `1.0.10` was reported as older than `1.0.9` because string compare ranks `"1" < "9"`. Now compares numeric segments per SemVer 2.0 (with pre-release suffix handling).
- **Powerlevel10k `MkdirAll` bypassed `--dry-run`** and created `~/.oh-my-zsh/custom/themes/` even when the user only wanted a preview. Now properly guarded.
- **22 pre-existing `errcheck` lint warnings**: real I/O issues (`resp.Body.Close`, `os.Setenv`) properly handled, `fmt`/`color.Color` print returns excluded by config.

### Removed
- `internal/config/defaults.yaml` (consolidated into `configs/default.yaml` — see Changed).

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

[Unreleased]: https://github.com/tldr-it-stepankutaj/setup-mac/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/tldr-it-stepankutaj/setup-mac/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/tldr-it-stepankutaj/setup-mac/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/tldr-it-stepankutaj/setup-mac/releases/tag/v1.0.0
