package config

import "github.com/spf13/viper"

// applyDefaults sets built-in default values for all configuration keys.
// These values are the lowest-precedence source and are overridden by all
// other configuration sources.
func applyDefaults(v *viper.Viper) {
	// Logging
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")
	v.SetDefault("log.file", "")

	// Provider
	v.SetDefault("provider.active", "openrouter")
	v.SetDefault("provider.model", "")

	// Agent
	v.SetDefault("agent.max_iterations", 50)
	v.SetDefault("agent.timeout_seconds", 300)
	v.SetDefault("agent.show_tool_output", true)
	v.SetDefault("agent.auto_approve", false)

	// UI
	v.SetDefault("ui.theme", "default")
	v.SetDefault("ui.layout", "default")
	v.SetDefault("ui.keybindings_file", "keybindings.yaml")

	// Collaboration
	v.SetDefault("collaboration.port", 7800)
	v.SetDefault("collaboration.tls_cert_file", "")
	v.SetDefault("collaboration.tls_key_file", "")
}
