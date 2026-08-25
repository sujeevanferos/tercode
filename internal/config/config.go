// Package config loads, validates, and exposes Tercode configuration.
//
// Configuration resolution order (lowest to highest precedence):
//   1. Built-in defaults (defaults.go)
//   2. Global user config (~/.config/tercode/config.yaml)
//   3. Project config (.config/config.yaml in workspace root)
//   4. Environment variables (TERCODE_*)
//   5. CLI flags
//   6. Interactive session overrides
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config is the top-level Tercode configuration structure.
type Config struct {
	// Log controls logging behaviour.
	Log LogConfig `mapstructure:"log"`

	// Provider holds active provider and model configuration.
	Provider ProviderConfig `mapstructure:"provider"`

	// Agent holds agent runtime configuration.
	Agent AgentConfig `mapstructure:"agent"`

	// UI holds UI presentation configuration.
	UI UIConfig `mapstructure:"ui"`

	// Workspace holds workspace-level configuration.
	Workspace WorkspaceConfig `mapstructure:"workspace"`

	// Collaboration holds LAN collaboration settings.
	Collaboration CollaborationConfig `mapstructure:"collaboration"`
}

// LogConfig controls the logging subsystem.
type LogConfig struct {
	// Level sets the minimum log level: debug, info, warn, error.
	Level string `mapstructure:"level"`

	// Format controls output format: text or json.
	Format string `mapstructure:"format"`

	// File writes logs to a file path in addition to stderr.
	// Empty means stderr only.
	File string `mapstructure:"file"`
}

// ProviderConfig holds the active provider and model selection.
type ProviderConfig struct {
	// Active is the provider ID to use by default.
	Active string `mapstructure:"active"`

	// Model is the model ID to use by default.
	Model string `mapstructure:"model"`

	// Providers is a map of provider ID → per-provider settings.
	Providers map[string]ProviderEntry `mapstructure:"providers"`
}

// ProviderEntry holds configuration for a single provider.
type ProviderEntry struct {
	// Type identifies the provider adapter: openrouter, openai, nvidia.
	Type string `mapstructure:"type"`

	// BaseURL overrides the default API base URL.
	BaseURL string `mapstructure:"base_url"`

	// APIKeyEnv is the environment variable name holding the API key.
	// When empty, the default env var for the adapter type is used.
	APIKeyEnv string `mapstructure:"api_key_env"`

	// DefaultModel overrides the global default model for this provider.
	DefaultModel string `mapstructure:"default_model"`
}

// AgentConfig holds agent runtime settings.
type AgentConfig struct {
	// MaxIterations caps the agent tool loop per task.
	MaxIterations int `mapstructure:"max_iterations"`

	// TimeoutSeconds is the per-task wall-clock timeout.
	TimeoutSeconds int `mapstructure:"timeout_seconds"`

	// ShowToolOutput controls whether raw tool output is shown in the TUI.
	ShowToolOutput bool `mapstructure:"show_tool_output"`

	// AutoApprove enables automatic approval of low-risk tool calls.
	AutoApprove bool `mapstructure:"auto_approve"`
}

// UIConfig holds presentation-layer settings.
type UIConfig struct {
	// Theme is the name of the active theme (loaded from .config/themes/).
	Theme string `mapstructure:"theme"`

	// Layout is the name of the active layout file.
	Layout string `mapstructure:"layout"`

	// KeybindingsFile is the path to a keybindings override file.
	KeybindingsFile string `mapstructure:"keybindings_file"`
}

// WorkspaceConfig holds workspace detection and behaviour settings.
type WorkspaceConfig struct {
	// Root overrides automatic workspace root detection.
	Root string `mapstructure:"root"`
}

// CollaborationConfig holds LAN collaboration settings.
type CollaborationConfig struct {
	// Port is the default port for hosting a collaboration workspace.
	Port int `mapstructure:"port"`

	// TLSCertFile is the path to a TLS certificate for the collaboration server.
	TLSCertFile string `mapstructure:"tls_cert_file"`

	// TLSKeyFile is the path to the TLS private key.
	TLSKeyFile string `mapstructure:"tls_key_file"`
}

// Loader orchestrates configuration loading from all sources.
type Loader struct {
	v             *viper.Viper
	globalCfgDir  string
	projectCfgDir string
}

// NewLoader creates a Loader. globalCfgDir is the user's ~/.config/tercode path.
// projectCfgDir is the .config/ path inside the current workspace root.
func NewLoader(globalCfgDir, projectCfgDir string) *Loader {
	v := viper.NewWithOptions(viper.KeyDelimiter("."))
	return &Loader{v: v, globalCfgDir: globalCfgDir, projectCfgDir: projectCfgDir}
}

// Load reads configuration from all sources in precedence order and returns
// a fully resolved Config.
func (l *Loader) Load() (*Config, error) {
	// Apply built-in defaults.
	applyDefaults(l.v)

	l.v.SetConfigName("config")
	l.v.SetConfigType("yaml")

	// 1. Global config directory.
	if l.globalCfgDir != "" {
		l.v.AddConfigPath(l.globalCfgDir)
	}

	// 2. Project .config/ directory (higher precedence — merged in after).
	// We handle this manually so project config overlays global config.
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("config: read global config: %w", err)
		}
	}

	// Merge project-level config on top.
	if l.projectCfgDir != "" {
		if err := l.mergeProjectConfig(); err != nil {
			return nil, err
		}
	}

	// 3. Environment variables: TERCODE_LOG_LEVEL → log.level, etc.
	l.v.SetEnvPrefix("TERCODE")
	l.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	l.v.AutomaticEnv()

	cfg := &Config{}
	if err := l.v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}
	return cfg, nil
}

// mergeProjectConfig reads and merges .config/config.yaml from the project root.
func (l *Loader) mergeProjectConfig() error {
	pv := viper.NewWithOptions(viper.KeyDelimiter("."))
	pv.SetConfigName("config")
	pv.SetConfigType("yaml")
	pv.AddConfigPath(l.projectCfgDir)

	if err := pv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil // No project config — not an error.
		}
		return fmt.Errorf("config: read project config: %w", err)
	}
	return l.v.MergeConfigMap(pv.AllSettings())
}

// SetFlag allows CLI flag values to override config values at runtime.
func (l *Loader) SetFlag(key, value string) {
	l.v.Set(key, value)
}

// ProjectConfigDir returns the resolved path of the project .config/ directory
// given a workspace root. It does not require the directory to exist.
func ProjectConfigDir(workspaceRoot string) string {
	return filepath.Join(workspaceRoot, ".config")
}

// EnsureProjectConfigDir creates the project .config/ directory structure.
func EnsureProjectConfigDir(workspaceRoot string) error {
	base := ProjectConfigDir(workspaceRoot)
	dirs := []string{
		base,
		filepath.Join(base, "themes"),
		filepath.Join(base, "layouts"),
		filepath.Join(base, "commands"),
		filepath.Join(base, "agents"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("config: ensure dir %s: %w", d, err)
		}
	}
	return nil
}
