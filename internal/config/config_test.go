package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sujeevanferos/tercode/internal/config"
)

func TestConfigDefaultsAndPrecedence(t *testing.T) {
	tempDir := t.TempDir()
	projectCfgDir := filepath.Join(tempDir, ".config")
	_ = os.MkdirAll(projectCfgDir, 0o755)

	customYAML := `
provider:
  active: "nvidia"
  model: "custom-nemotron"
agent:
  max_iterations: 42
`
	_ = os.WriteFile(filepath.Join(projectCfgDir, "config.yaml"), []byte(customYAML), 0o644)

	loader := config.NewLoader("", projectCfgDir)
	cfg, err := loader.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Provider.Active != "nvidia" {
		t.Errorf("expected provider 'nvidia', got %q", cfg.Provider.Active)
	}
	if cfg.Provider.Model != "custom-nemotron" {
		t.Errorf("expected model 'custom-nemotron', got %q", cfg.Provider.Model)
	}
	if cfg.Agent.MaxIterations != 42 {
		t.Errorf("expected max_iterations 42, got %d", cfg.Agent.MaxIterations)
	}
	// Verify defaults still present for unconfigured keys
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level 'info', got %q", cfg.Log.Level)
	}
}
