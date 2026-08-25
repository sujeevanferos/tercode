// Package runtime provides the top-level AgentRuntime coordination layer.
package runtime

import (
	"context"
	"fmt"

	"github.com/sujeevanferos/tercode/internal/agent/modes"
	"github.com/sujeevanferos/tercode/internal/agent/session"
	"github.com/sujeevanferos/tercode/internal/agent/task"
	"github.com/sujeevanferos/tercode/internal/context/manager"
	"github.com/sujeevanferos/tercode/internal/events"
	"github.com/sujeevanferos/tercode/internal/provider"
	providerreg "github.com/sujeevanferos/tercode/internal/provider/registry"
	toolreg "github.com/sujeevanferos/tercode/internal/tools/registry"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// Runtime orchestrates sessions, tasks, context budgeting, and tool execution.
type Runtime struct {
	providers *providerreg.Registry
	tools     *toolreg.Registry
	sessions  *session.Manager
	executor  *task.Executor
	context   *manager.Manager
	events    *events.Bus
	workspace *workspace.Workspace
}

func New(
	providers *providerreg.Registry,
	tools *toolreg.Registry,
	sessions *session.Manager,
	ctxMgr *manager.Manager,
	bus *events.Bus,
	ws *workspace.Workspace,
) *Runtime {
	return &Runtime{
		providers: providers,
		tools:     tools,
		sessions:  sessions,
		executor:  task.NewExecutor(tools, bus),
		context:   ctxMgr,
		events:    bus,
		workspace: ws,
	}
}

// CreateSession creates a new agent session.
func (r *Runtime) CreateSession(providerID, modelID string, mode session.Mode) *session.Session {
	wsID := "default"
	if r.workspace != nil {
		wsID = r.workspace.ID
	}
	return r.sessions.Create(wsID, providerID, modelID, mode)
}

// Run executes a prompt turn in a session.
func (r *Runtime) Run(ctx context.Context, sessionID, prompt string) (string, error) {
	s, ok := r.sessions.Get(sessionID)
	if !ok {
		return "", fmt.Errorf("session %q not found", sessionID)
	}

	p, err := r.providers.Get(s.ProviderID)
	if err != nil {
		return "", fmt.Errorf("provider %q: %w", s.ProviderID, err)
	}

	sysPrompt := modes.GetSystemPrompt(s.Mode)
	if r.workspace != nil {
		sysPrompt += fmt.Sprintf("\n\nWorkspace root: %s", r.workspace.RootPath)
	}

	result, err := r.executor.ExecuteTask(ctx, s, p, s.ModelID, sysPrompt, prompt, 30)
	if err != nil {
		return "", err
	}

	// Update session history
	r.sessions.Append(sessionID,
		provider.Message{Role: provider.RoleUser, Content: prompt},
		provider.Message{Role: provider.RoleAssistant, Content: result},
	)

	return result, nil
}
