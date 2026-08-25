// Package events defines all structured event types that flow through
// Tercode's internal event bus.
//
// Events are immutable facts that have already occurred. They are published
// by domain subsystems (agent, workspace, tools, providers, collaboration)
// and consumed by the TUI, logging, and collaboration subsystems.
//
// Producers must never mutate an event after calling Publish.
// Consumers must never mutate a received event.
package events

import "time"

// Type identifies the category of a structured event.
type Type string

const (
	// Session lifecycle.
	SessionStarted Type = "SESSION_STARTED"
	SessionEnded   Type = "SESSION_ENDED"

	// Task lifecycle.
	TaskCreated   Type = "TASK_CREATED"
	TaskStarted   Type = "TASK_STARTED"
	TaskCompleted Type = "TASK_COMPLETED"
	TaskFailed    Type = "TASK_FAILED"
	TaskCancelled Type = "TASK_CANCELLED"

	// Model interaction.
	ModelRequest  Type = "MODEL_REQUEST"
	ModelResponse Type = "MODEL_RESPONSE"
	ModelStream   Type = "MODEL_STREAM"
	ModelError    Type = "MODEL_ERROR"

	// Tool lifecycle.
	ToolRequest   Type = "TOOL_REQUEST"
	ToolStarted   Type = "TOOL_STARTED"
	ToolOutput    Type = "TOOL_OUTPUT"
	ToolCompleted Type = "TOOL_COMPLETED"
	ToolFailed    Type = "TOOL_FAILED"

	// Approval flow.
	ApprovalRequired Type = "APPROVAL_REQUIRED"
	ApprovalGranted  Type = "APPROVAL_GRANTED"
	ApprovalDenied   Type = "APPROVAL_DENIED"

	// Workspace / filesystem.
	FileCreated  Type = "FILE_CREATED"
	FileModified Type = "FILE_MODIFIED"
	FileDeleted  Type = "FILE_DELETED"
	GitChanged   Type = "GIT_CHANGED"

	// LAN collaboration.
	UserJoined Type = "USER_JOINED"
	UserLeft   Type = "USER_LEFT"

	// Application.
	AppStarted  Type = "APP_STARTED"
	AppShutdown Type = "APP_SHUTDOWN"
	Error       Type = "ERROR"
)

// Event is the universal event envelope transported on the EventBus.
// Every field except Payload is populated by the publisher.
type Event struct {
	ID          string    `json:"id"`
	Type        Type      `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
	WorkspaceID string    `json:"workspace_id,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	TaskID      string    `json:"task_id,omitempty"`
	// Payload carries the event-specific data. Consumers should type-assert
	// to the expected concrete payload type for the given event Type.
	Payload any `json:"payload,omitempty"`
}

// --- Payload types ---

// TaskPayload carries task lifecycle information.
type TaskPayload struct {
	Mode    string `json:"mode"`
	Request string `json:"request,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ModelRequestPayload carries summary information about a model call.
type ModelRequestPayload struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	// TokensEstimate is a best-effort token count before the call.
	TokensEstimate int `json:"tokens_estimate,omitempty"`
}

// ModelResponsePayload carries model response summary.
type ModelResponsePayload struct {
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	FinishReason  string `json:"finish_reason"`
	PromptTokens  int    `json:"prompt_tokens,omitempty"`
	OutputTokens  int    `json:"output_tokens,omitempty"`
	LatencyMs     int64  `json:"latency_ms,omitempty"`
}

// ModelStreamPayload carries a single streaming chunk from the model.
type ModelStreamPayload struct {
	Delta string `json:"delta"`
	Done  bool   `json:"done"`
}

// ToolRequestPayload describes a tool call initiated by the agent.
type ToolRequestPayload struct {
	ToolName string `json:"tool_name"`
	Input    any    `json:"input,omitempty"`
}

// ToolOutputPayload carries incremental output from a running tool.
type ToolOutputPayload struct {
	ToolName string `json:"tool_name"`
	Output   string `json:"output"`
}

// ToolCompletedPayload carries the final result of a tool execution.
type ToolCompletedPayload struct {
	ToolName string `json:"tool_name"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// ApprovalPayload describes an approval request from the permission system.
type ApprovalPayload struct {
	ToolName    string `json:"tool_name"`
	Description string `json:"description"`
	Risk        string `json:"risk"` // low, medium, high
}

// FilePayload carries information about a workspace file change.
type FilePayload struct {
	Path string `json:"path"`
}

// UserPresencePayload carries LAN user join/leave information.
type UserPresencePayload struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
}

// ErrorPayload carries a structured application error.
type ErrorPayload struct {
	Component string `json:"component"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}
