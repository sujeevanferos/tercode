// Package tools defines the Tool interface, ToolRegistry, and ToolResult.
package tools

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sujeevanferos/tercode/internal/security"
)

// Tool is the universal interface that all native tools and extensions implement.
type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]any
	Risk() security.Risk
	Execute(ctx context.Context, input json.RawMessage) ToolResult
}

// ToolResult represents the structured outcome of a tool execution.
type ToolResult struct {
	ToolName  string        `json:"tool_name"`
	Success   bool          `json:"success"`
	Output    string        `json:"output"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// FormatForModel returns a clean string representation suitable for feeding back to an LLM.
func (r ToolResult) FormatForModel() string {
	if !r.Success {
		return "Error: " + r.Error + "\nOutput: " + r.Output
	}
	return r.Output
}
