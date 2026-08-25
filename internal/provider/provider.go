// Package provider defines the Provider interface and all shared model types.
//
// Every provider adapter must implement the Provider interface.
// The agent runtime must never import adapter packages directly — it
// communicates exclusively through this interface.
package provider

import (
	"context"
	"io"
)

// Provider is the stable interface every provider adapter must implement.
// Provider-specific logic (request marshalling, authentication headers,
// error code mapping) lives entirely inside the adapter, not here.
type Provider interface {
	// ID returns the unique provider identifier (e.g. "openrouter", "nvidia").
	ID() string

	// Name returns the human-readable display name.
	Name() string

	// Authenticate validates credentials without making an inference call.
	// Returns nil on success.
	Authenticate(ctx context.Context) error

	// ListModels returns the models available on this provider.
	ListModels(ctx context.Context) ([]Model, error)

	// Chat sends a request and returns the complete response.
	Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)

	// Stream sends a request and returns a reader for streaming token output.
	// The caller is responsible for closing the StreamReader.
	Stream(ctx context.Context, req ChatRequest) (StreamReader, error)

	// TestModel runs a standard test sequence against the given model ID and
	// returns a detailed result covering credentials, connectivity, inference,
	// streaming, and tool calling.
	TestModel(ctx context.Context, modelID string) TestResult

	// Capabilities returns the capability set for a given model ID.
	// Providers should return a zeroed Capabilities if the model is unknown.
	Capabilities(modelID string) Capabilities
}

// StreamReader provides access to a streaming model response.
// Implementations emit tokens via Read and signal completion with io.EOF.
type StreamReader interface {
	io.ReadCloser
	// Event returns the next structured stream event.
	// Returns io.EOF when the stream is finished.
	Event() (StreamEvent, error)
}

// StreamEvent is a single token or control event from a streaming model response.
type StreamEvent struct {
	// Delta contains the incremental text content.
	Delta string
	// ToolCall is populated when the model emits a tool call during streaming.
	ToolCall *ToolCall
	// FinishReason is non-empty on the final event.
	FinishReason string
	// Done signals the end of the stream.
	Done bool
}

// Capabilities describes what a model supports.
type Capabilities struct {
	Streaming      bool
	ToolCalling    bool
	Vision         bool
	Reasoning      bool
	StructuredJSON bool
	// ContextWindow is the maximum token context the model supports.
	ContextWindow int
}

// Modality describes input/output types supported by a model.
type Modality string

const (
	ModalityText  Modality = "text"
	ModalityImage Modality = "image"
	ModalityAudio Modality = "audio"
)

// Model holds normalised metadata for a single model from any provider.
type Model struct {
	// ProviderID is the ID of the provider that owns this model.
	ProviderID string
	// ID is the model identifier used in API calls (e.g. "gpt-4o").
	ID string
	// Name is the human-readable display name.
	Name string
	// ContextWindow is the maximum token context.
	ContextWindow int
	// Capabilities describes what this model can do.
	Caps Capabilities
	// Modalities lists supported input types.
	Modalities []Modality
	// InputPricePer1M is the USD cost per 1M input tokens (0 if unknown).
	InputPricePer1M float64
	// OutputPricePer1M is the USD cost per 1M output tokens (0 if unknown).
	OutputPricePer1M float64
}

// Role defines the sender of a conversation message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a single turn in a conversation.
type Message struct {
	Role Role
	// Content holds the text content of the message.
	Content string
	// ToolCalls is populated when the assistant requests tool execution.
	ToolCalls []ToolCall
	// ToolCallID links a tool result back to the originating tool call.
	ToolCallID string
	// ToolName is the name of the tool that produced this result (role=tool).
	ToolName string
}

// ToolCall represents a single tool invocation requested by the model.
type ToolCall struct {
	ID       string
	Name     string
	// Arguments is the raw JSON string of tool arguments.
	Arguments string
}

// ToolDefinition describes a tool the model may call.
type ToolDefinition struct {
	Name        string
	Description string
	// Parameters is a JSON Schema object describing the tool's input.
	Parameters map[string]any
}

// ChatRequest is the normalised input to a model inference call.
type ChatRequest struct {
	// Model is the model ID to use. Empty means use the provider default.
	Model string
	// Messages is the conversation history including the new user message.
	Messages []Message
	// Tools are the tools the model may invoke.
	Tools []ToolDefinition
	// MaxTokens caps the output length. 0 means provider default.
	MaxTokens int
	// Temperature controls sampling randomness (0.0–2.0). -1 means default.
	Temperature float64
	// SystemPrompt overrides or prepends to the system message if non-empty.
	SystemPrompt string
}

// ChatResponse is the normalised output from a non-streaming model call.
type ChatResponse struct {
	// Message is the assistant's reply.
	Message Message
	// FinishReason explains why generation stopped (stop, tool_calls, length).
	FinishReason string
	// PromptTokens is the number of input tokens consumed.
	PromptTokens int
	// OutputTokens is the number of output tokens generated.
	OutputTokens int
}

// TestResult summarises the outcome of a provider/model health check.
type TestResult struct {
	ModelID string
	// Overall is true only when all stages passed.
	Overall bool
	// Stages contains per-stage pass/fail results.
	Stages map[string]StageResult
	// LatencyMs is the end-to-end latency in milliseconds.
	LatencyMs int64
}

// StageResult describes the outcome of a single test stage.
type StageResult struct {
	Passed bool
	Error  string
}
