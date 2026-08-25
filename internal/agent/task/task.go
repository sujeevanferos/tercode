// Package task defines tasks and the autonomous tool execution loop.
package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sujeevanferos/tercode/internal/agent/session"
	"github.com/sujeevanferos/tercode/internal/events"
	"github.com/sujeevanferos/tercode/internal/provider"
	"github.com/sujeevanferos/tercode/internal/tools/registry"
)

// Status represents the state of a Task.
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
	StatusCancelled Status = "CANCELLED"
)

// Task represents an autonomous unit of work.
type Task struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Request   string    `json:"request"`
	Status    Status    `json:"status"`
	Result    string    `json:"result,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Executor runs the multi-turn agent tool execution loop.
type Executor struct {
	tools  *registry.Registry
	events *events.Bus
}

func NewExecutor(tools *registry.Registry, bus *events.Bus) *Executor {
	return &Executor{tools: tools, events: bus}
}

// ExecuteTask runs the agent loop:
// 1. Sends conversation history + tool definitions to the provider.
// 2. If the model emits tool calls, executes them via ToolRegistry.
// 3. Feeds tool results back to the model.
// 4. Loops until the model yields a final text response or maxIterations is hit.
func (e *Executor) ExecuteTask(
	ctx context.Context,
	s *session.Session,
	p provider.Provider,
	modelID string,
	systemPrompt string,
	userPrompt string,
	maxIterations int,
) (string, error) {
	if e.events != nil {
		e.events.Publish(events.New(events.TaskStarted, events.TaskPayload{
			Mode:    string(s.Mode),
			Request: userPrompt,
		}))
	}

	if maxIterations <= 0 {
		maxIterations = 25
	}

	// Prepare messages
	var messages []provider.Message
	messages = append(messages, provider.Message{
		Role:    provider.RoleSystem,
		Content: systemPrompt,
	})
	messages = append(messages, s.History...)
	messages = append(messages, provider.Message{
		Role:    provider.RoleUser,
		Content: userPrompt,
	})

	var toolDefs []provider.ToolDefinition
	if s.Mode != session.ModePlan && s.Mode != session.ModeChat {
		toolDefs = e.tools.Definitions()
	}

	for i := 0; i < maxIterations; i++ {
		select {
		case <-ctx.Done():
			if e.events != nil {
				e.events.Publish(events.New(events.TaskCancelled, events.TaskPayload{Mode: string(s.Mode)}))
			}
			return "", ctx.Err()
		default:
		}

		if e.events != nil {
			e.events.Publish(events.New(events.ModelRequest, events.ModelRequestPayload{
				Provider: p.ID(),
				Model:    modelID,
			}))
		}

		// Execute chat completion
		resp, err := p.Chat(ctx, provider.ChatRequest{
			Model:    modelID,
			Messages: messages,
			Tools:    toolDefs,
		})
		if err != nil {
			if e.events != nil {
				e.events.Publish(events.New(events.ModelError, events.ErrorPayload{
					Component: "provider",
					Code:      "CHAT_FAILED",
					Message:   err.Error(),
				}))
				e.events.Publish(events.New(events.TaskFailed, events.TaskPayload{
					Mode:  string(s.Mode),
					Error: err.Error(),
				}))
			}
			return "", fmt.Errorf("agent executor: chat: %w", err)
		}

		if e.events != nil {
			e.events.Publish(events.New(events.ModelResponse, events.ModelResponsePayload{
				Provider:     p.ID(),
				Model:        modelID,
				FinishReason: resp.FinishReason,
				PromptTokens: resp.PromptTokens,
				OutputTokens: resp.OutputTokens,
			}))
		}

		messages = append(messages, resp.Message)

		// Check if tools were called
		if len(resp.Message.ToolCalls) == 0 {
			// Done! Return final message
			finalText := resp.Message.Content
			if e.events != nil {
				e.events.Publish(events.New(events.TaskCompleted, events.TaskPayload{
					Mode:    string(s.Mode),
					Request: userPrompt,
				}))
			}
			return finalText, nil
		}

		// Execute all tool calls
		for _, tc := range resp.Message.ToolCalls {
			toolRes := e.tools.Dispatch(ctx, s.ID, tc.Name, json.RawMessage(tc.Arguments))

			// Append tool output as role=tool message
			messages = append(messages, provider.Message{
				Role:       provider.RoleTool,
				Content:    toolRes.FormatForModel(),
				ToolCallID: tc.ID,
				ToolName:   tc.Name,
			})
		}
	}

	return "", fmt.Errorf("agent executor: exceeded max iterations (%d)", maxIterations)
}
