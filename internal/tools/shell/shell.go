// Package shell provides sandboxed terminal execution with timeout and cancellation.
package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// ShellTool executes commands in the workspace environment.
type ShellTool struct {
	ws *workspace.Workspace
}

func New(ws *workspace.Workspace) *ShellTool {
	return &ShellTool{ws: ws}
}

func (t *ShellTool) Name() string        { return "shell" }
func (t *ShellTool) Description() string { return "Execute a shell command inside the workspace root directory." }
func (t *ShellTool) Risk() security.Risk { return security.RiskHigh }

func (t *ShellTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "Shell command line to execute.",
			},
			"timeout_seconds": map[string]any{
				"type":        "integer",
				"description": "Max execution time in seconds (default: 60).",
			},
		},
		"required": []string{"command"},
	}
}

func (t *ShellTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Command        string `json:"command"`
		TimeoutSeconds int    `json:"timeout_seconds"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid parameters"}
	}

	timeout := 60 * time.Second
	if input.TimeoutSeconds > 0 {
		timeout = time.Duration(input.TimeoutSeconds) * time.Second
	}

	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(execCtx, "powershell", "-NoProfile", "-Command", input.Command)
	} else {
		cmd = exec.CommandContext(execCtx, "sh", "-c", input.Command)
	}
	cmd.Dir = t.ws.RootPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n--- stderr ---\n"
		}
		output += stderr.String()
	}

	if err != nil {
		return tools.ToolResult{
			ToolName: t.Name(),
			Success:  false,
			Output:   output,
			Error:    fmt.Sprintf("command failed: %v", err),
		}
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   output,
	}
}
