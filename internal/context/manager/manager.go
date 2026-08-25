// Package manager manages context assembly, token budgeting, and compaction.
package manager

import (
	"context"
	"strings"

	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// Manager constructs focused context budgets for LLM calls.
type Manager struct {
	ws *workspace.Workspace
}

// New creates a ContextManager.
func New(ws *workspace.Workspace) *Manager {
	return &Manager{ws: ws}
}

// BuildContext prepares a ChatRequest with the system prompt, workspace awareness, and compacted history.
func (m *Manager) BuildContext(ctx context.Context, systemPrompt string, history []provider.Message, maxTokens int) provider.ChatRequest {
	var msgs []provider.Message

	// Always prepend workspace root info in system context
	sys := systemPrompt
	if m.ws != nil {
		wsInfo := "\n\nWorkspace root: " + m.ws.RootPath
		sys += wsInfo
	}

	msgs = append(msgs, provider.Message{
		Role:    provider.RoleSystem,
		Content: sys,
	})

	// Add conversation history
	msgs = append(msgs, history...)

	return provider.ChatRequest{
		SystemPrompt: sys,
		Messages:     msgs,
		MaxTokens:    maxTokens,
	}
}

// Compact summarizes or truncates very long message histories when exceeding token budgets.
func (m *Manager) Compact(history []provider.Message, maxTurns int) []provider.Message {
	if len(history) <= maxTurns {
		return history
	}
	// Keep the first message and the last (maxTurns-1) messages
	compacted := make([]provider.Message, 0, maxTurns)
	compacted = append(compacted, history[0])
	compacted = append(compacted, provider.Message{
		Role:    provider.RoleSystem,
		Content: "[Earlier conversation history compacted for token budget]",
	})
	compacted = append(compacted, history[len(history)-(maxTurns-2):]...)
	return compacted
}

// EstimateTokens is a fast, conservative heuristic (1 token ~= 4 chars).
func EstimateTokens(text string) int {
	return (len(strings.TrimSpace(text)) + 3) / 4
}
