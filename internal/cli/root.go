// Package cli implements Cobra subcommands (root, models, provider, host, connect).
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sujeevanferos/tercode/internal/app"
	"github.com/sujeevanferos/tercode/internal/config"
	"github.com/sujeevanferos/tercode/internal/platform"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "tercode",
		Short: "Tercode — Local-first, provider-agnostic agentic terminal development environment",
		Long: `Tercode is an open-source, provider-agnostic coding environment for the terminal.
Powerful by default. Limitless by configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Auto-detect workspace root
			ws, err := workspace.Detect("")
			if err != nil {
				return err
			}

			// Load multi-source config
			globalCfg, _ := platform.ConfigDir()
			projectCfg := config.ProjectConfigDir(ws.RootPath)
			loader := config.NewLoader(globalCfg, projectCfg)

			cfg, err := loader.Load()
			if err != nil {
				return fmt.Errorf("config error: %w", err)
			}

			// Construct Application
			application, err := app.New(cfg, ws)
			if err != nil {
				return fmt.Errorf("app init: %w", err)
			}

			// Run Interactive TUI
			return application.RunTUI()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .config/config.yaml)")

	// Subcommands
	rootCmd.AddCommand(modelsCmd)
	rootCmd.AddCommand(providerCmd)
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List discovered models from active provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		ws, _ := workspace.Detect("")
		globalCfg, _ := platform.ConfigDir()
		loader := config.NewLoader(globalCfg, config.ProjectConfigDir(ws.RootPath))
		cfg, _ := loader.Load()

		application, err := app.New(cfg, ws)
		if err != nil {
			return err
		}

		p, err := application.Providers.Get(cfg.Provider.Active)
		if err != nil {
			return err
		}

		models, err := p.ListModels(cmd.Context())
		if err != nil {
			return fmt.Errorf("list models from %s: %w", p.Name(), err)
		}

		fmt.Printf("Discovered models from %s (%d models):\n\n", p.Name(), len(models))
		for _, m := range models {
			fmt.Printf("  %-40s (Context: %d tokens)\n", m.ID, m.ContextWindow)
		}
		return nil
	},
}

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Inspect and test configured AI providers",
}

func init() {
	providerCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all registered provider adapters",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Registered Provider Adapters:")
			fmt.Println("  - openrouter  (OpenRouter Gateway - 200+ models)")
			fmt.Println("  - nvidia      (NVIDIA NIM - Enterprise GPU Inference)")
			fmt.Println("  - openai      (Generic OpenAI Compatible - Ollama, vLLM, DeepSeek)")
		},
	})
}

// Execute runs the root CLI command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
