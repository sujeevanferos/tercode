package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sujeevanferos/tercode/internal/events"
	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
)

// Registry manages tool registration, permissions checking, and dispatch.
type Registry struct {
	mu          sync.RWMutex
	tools       map[string]tools.Tool
	permissions *security.PermissionManager
	events      *events.Bus
}

// New creates an initialized ToolRegistry.
func New(permissions *security.PermissionManager, bus *events.Bus) *Registry {
	return &Registry{
		tools:       make(map[string]tools.Tool),
		permissions: permissions,
		events:      bus,
	}
}

// Register registers a tool.
func (r *Registry) Register(t tools.Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[t.Name()]; exists {
		return fmt.Errorf("tool registry: %q already registered", t.Name())
	}
	r.tools[t.Name()] = t
	return nil
}

// Get returns the tool by name.
func (r *Registry) Get(name string) (tools.Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

// Definitions returns tool definitions in the format expected by provider adapters.
func (r *Registry) Definitions() []provider.ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	defs := make([]provider.ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		defs = append(defs, provider.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters:  t.Parameters(),
		})
	}
	return defs
}

// Dispatch evaluates permissions and executes the requested tool.
func (r *Registry) Dispatch(ctx context.Context, sessionID, toolName string, input json.RawMessage) tools.ToolResult {
	start := time.Now()
	t, ok := r.Get(toolName)
	if !ok {
		return tools.ToolResult{
			ToolName: toolName,
			Success:  false,
			Error:    fmt.Sprintf("tool %q not found", toolName),
			Duration: time.Since(start),
		}
	}

	if r.events != nil {
		r.events.Publish(events.New(events.ToolRequest, events.ToolRequestPayload{
			ToolName: toolName,
			Input:    string(input),
		}))
	}

	// Permission check
	if r.permissions != nil {
		permRes, err := r.permissions.Check(ctx, security.PermissionRequest{
			ToolName:    toolName,
			Action:      "execute",
			Description: fmt.Sprintf("Execute %s with %s", toolName, string(input)),
			Risk:        t.Risk(),
			SessionID:   sessionID,
		})
		if err != nil || !permRes.Granted {
			reason := "permission denied"
			if err != nil {
				reason = err.Error()
			} else if permRes.Reason != "" {
				reason = permRes.Reason
			}
			return tools.ToolResult{
				ToolName: toolName,
				Success:  false,
				Error:    reason,
				Duration: time.Since(start),
			}
		}
	}

	if r.events != nil {
		r.events.Publish(events.New(events.ToolStarted, events.ToolRequestPayload{
			ToolName: toolName,
		}))
	}

	res := t.Execute(ctx, input)
	res.Duration = time.Since(start)

	if r.events != nil {
		if res.Success {
			r.events.Publish(events.New(events.ToolCompleted, events.ToolCompletedPayload{
				ToolName: toolName,
				Success:  true,
			}))
		} else {
			r.events.Publish(events.New(events.ToolFailed, events.ToolCompletedPayload{
				ToolName: toolName,
				Success:  false,
				Error:    res.Error,
			}))
		}
	}

	return res
}
