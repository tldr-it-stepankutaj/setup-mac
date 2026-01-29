package config

// Config represents the root configuration structure
type Config struct {
	Version  string         `yaml:"version" mapstructure:"version"`
	Settings SettingsConfig `yaml:"settings" mapstructure:"settings"`
	Homebrew HomebrewConfig `yaml:"homebrew" mapstructure:"homebrew"`
	Terminal TerminalConfig `yaml:"terminal" mapstructure:"terminal"`
	Shell    ShellConfig    `yaml:"shell" mapstructure:"shell"`
	MacOS    MacOSConfig    `yaml:"macos" mapstructure:"macos"`
	Git      GitConfig      `yaml:"git" mapstructure:"git"`
	SSH      SSHConfig      `yaml:"ssh" mapstructure:"ssh"`
}

// SettingsConfig contains global settings
type SettingsConfig struct {
	DryRun         bool `yaml:"dry_run" mapstructure:"dry_run"`
	Interactive    bool `yaml:"interactive" mapstructure:"interactive"`
	BackupDotfiles bool `yaml:"backup_dotfiles" mapstructure:"backup_dotfiles"`
}

// HomebrewConfig contains Homebrew installation settings
type HomebrewConfig struct {
	Install  bool     `yaml:"install" mapstructure:"install"`
	Formulae []string `yaml:"formulae" mapstructure:"formulae"`
	Casks    []string `yaml:"casks" mapstructure:"casks"`
	Taps     []string `yaml:"taps" mapstructure:"taps"`
}

// TerminalConfig contains terminal-related settings
type TerminalConfig struct {
	OhMyZsh       OhMyZshConfig       `yaml:"oh_my_zsh" mapstructure:"oh_my_zsh"`
	Powerlevel10k Powerlevel10kConfig `yaml:"powerlevel10k" mapstructure:"powerlevel10k"`
}

// OhMyZshConfig contains Oh-My-Zsh settings
type OhMyZshConfig struct {
	Install bool     `yaml:"install" mapstructure:"install"`
	Plugins []string `yaml:"plugins" mapstructure:"plugins"`
	Theme   string   `yaml:"theme" mapstructure:"theme"`
}

// Powerlevel10kConfig contains Powerlevel10k settings
type Powerlevel10kConfig struct {
	Install bool   `yaml:"install" mapstructure:"install"`
	Style   string `yaml:"style" mapstructure:"style"`
}

// ShellConfig contains shell customization settings
type ShellConfig struct {
	Aliases     map[string]string `yaml:"aliases" mapstructure:"aliases"`
	Environment map[string]string `yaml:"environment" mapstructure:"environment"`
	ZshrcExtras []string          `yaml:"zshrc_extras" mapstructure:"zshrc_extras"`
}

// MacOSConfig contains macOS system settings
type MacOSConfig struct {
	Configure bool          `yaml:"configure" mapstructure:"configure"`
	Defaults  MacOSDefaults `yaml:"defaults" mapstructure:"defaults"`
}

// MacOSDefaults contains macOS defaults settings
type MacOSDefaults struct {
	Dock     DockDefaults     `yaml:"dock" mapstructure:"dock"`
	Finder   FinderDefaults   `yaml:"finder" mapstructure:"finder"`
	Keyboard KeyboardDefaults `yaml:"keyboard" mapstructure:"keyboard"`
}

// DockDefaults contains Dock settings
type DockDefaults struct {
	Autohide      bool `yaml:"autohide" mapstructure:"autohide"`
	AutohideDelay int  `yaml:"autohide_delay" mapstructure:"autohide_delay"`
	TileSize      int  `yaml:"tile_size" mapstructure:"tile_size"`
	Magnification bool `yaml:"magnification" mapstructure:"magnification"`
	MinimizeToApp bool `yaml:"minimize_to_app" mapstructure:"minimize_to_app"`
	ShowRecents   bool `yaml:"show_recents" mapstructure:"show_recents"`
}

// FinderDefaults contains Finder settings
type FinderDefaults struct {
	ShowHiddenFiles  bool   `yaml:"show_hidden_files" mapstructure:"show_hidden_files"`
	ShowExtensions   bool   `yaml:"show_extensions" mapstructure:"show_extensions"`
	ShowPathBar      bool   `yaml:"show_path_bar" mapstructure:"show_path_bar"`
	ShowStatusBar    bool   `yaml:"show_status_bar" mapstructure:"show_status_bar"`
	DefaultViewStyle string `yaml:"default_view_style" mapstructure:"default_view_style"`
}

// KeyboardDefaults contains Keyboard settings
type KeyboardDefaults struct {
	KeyRepeat          int  `yaml:"key_repeat" mapstructure:"key_repeat"`
	InitialKeyRepeat   int  `yaml:"initial_key_repeat" mapstructure:"initial_key_repeat"`
	DisableSmartQuotes bool `yaml:"disable_smart_quotes" mapstructure:"disable_smart_quotes"`
	DisableSmartDashes bool `yaml:"disable_smart_dashes" mapstructure:"disable_smart_dashes"`
}

// GitConfig contains Git settings
type GitConfig struct {
	Configure bool              `yaml:"configure" mapstructure:"configure"`
	User      GitUser           `yaml:"user" mapstructure:"user"`
	Aliases   map[string]string `yaml:"aliases" mapstructure:"aliases"`
	Settings  map[string]string `yaml:"settings" mapstructure:"settings"`
}

// GitUser contains Git user settings
type GitUser struct {
	Name  string `yaml:"name" mapstructure:"name"`
	Email string `yaml:"email" mapstructure:"email"`
}

// SSHConfig contains SSH settings
type SSHConfig struct {
	GenerateKey bool   `yaml:"generate_key" mapstructure:"generate_key"`
	KeyType     string `yaml:"key_type" mapstructure:"key_type"`
	KeyFile     string `yaml:"key_file" mapstructure:"key_file"`
	Comment     string `yaml:"comment" mapstructure:"comment"`
}
