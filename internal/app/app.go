// Package app constructs and coordinates all Tercode subsystems.
package app

import (
	"context"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sujeevanferos/tercode/internal/agent/runtime"
	"github.com/sujeevanferos/tercode/internal/agent/session"
	"github.com/sujeevanferos/tercode/internal/config"
	"github.com/sujeevanferos/tercode/internal/context/manager"
	"github.com/sujeevanferos/tercode/internal/events"
	"github.com/sujeevanferos/tercode/internal/persistence"
	"github.com/sujeevanferos/tercode/internal/platform"
	"github.com/sujeevanferos/tercode/internal/provider/adapters/nvidia"
	"github.com/sujeevanferos/tercode/internal/provider/adapters/openai"
	"github.com/sujeevanferos/tercode/internal/provider/adapters/openrouter"
	"github.com/sujeevanferos/tercode/internal/provider/registry"
	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools/filesystem"
	"github.com/sujeevanferos/tercode/internal/tools/git"
	toolreg "github.com/sujeevanferos/tercode/internal/tools/registry"
	"github.com/sujeevanferos/tercode/internal/tools/search"
	"github.com/sujeevanferos/tercode/internal/tools/shell"
	"github.com/sujeevanferos/tercode/internal/tools/testbuild"
	tuiapp "github.com/sujeevanferos/tercode/internal/tui/app"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// App is the root Tercode application container.
type App struct {
	Config    *config.Config
	Workspace *workspace.Workspace
	Events    *events.Bus
	Providers *registry.Registry
	Tools     *toolreg.Registry
	Agent     *runtime.Runtime
	Secrets   *security.SecretManager
	DB        *persistence.DB
}

// New initializes all services with clean dependency injection.
func New(cfg *config.Config, ws *workspace.Workspace) (*App, error) {
	bus := events.NewBus()
	secrets := security.NewSecretManager()

	// Persistence
	dataDir, _ := platform.DataDir()
	platform.EnsureDir(dataDir)
	db, _ := persistence.Open(filepath.Join(dataDir, "tercode.db"))

	// Providers
	providers := registry.NewRegistry()
	openrouterKey := secrets.Get("OPENROUTER_API_KEY")
	providers.Register(openrouter.New(openrouterKey))

	nvidiaKey := secrets.Get("NVIDIA_API_KEY")
	providers.Register(nvidia.New(nvidiaKey, ""))

	openaiKey := secrets.Get("OPENAI_API_KEY")
	providers.Register(openai.New("openai", "OpenAI Compatible", "", openaiKey))

	// Permissions & Tools
	permMgr := security.NewPermissionManager(nil, func(ctx context.Context, req security.PermissionRequest) (bool, error) {
		return true, nil // Auto-approve in MVP default
	})

	tools := toolreg.New(permMgr, bus)
	tools.Register(filesystem.NewReadFileTool(ws))
	tools.Register(filesystem.NewWriteFileTool(ws))
	tools.Register(filesystem.NewEditFileTool(ws))
	tools.Register(filesystem.NewListDirectoryTool(ws))
	tools.Register(shell.New(ws))
	tools.Register(git.New(ws))
	tools.Register(search.New(ws))
	tools.Register(testbuild.New(ws))

	// Context & Sessions & Agent Runtime
	ctxMgr := manager.New(ws)
	sessionMgr := session.NewManager()
	agentRt := runtime.New(providers, tools, sessionMgr, ctxMgr, bus, ws)

	return &App{
		Config:    cfg,
		Workspace: ws,
		Events:    bus,
		Providers: providers,
		Tools:     tools,
		Agent:     agentRt,
		Secrets:   secrets,
		DB:        db,
	}, nil
}

// RunTUI launches the interactive terminal interface.
func (a *App) RunTUI() error {
	sess := a.Agent.CreateSession(a.Config.Provider.Active, a.Config.Provider.Model, session.ModeCode)
	model := tuiapp.New(a.Agent, a.Workspace, a.Events, sess)

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
