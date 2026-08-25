// Package git provides native Git inspection and management tools.
package git

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// GitTool handles status, diff, log, and commits.
type GitTool struct {
	ws *workspace.Workspace
}

func New(ws *workspace.Workspace) *GitTool {
	return &GitTool{ws: ws}
}

func (t *GitTool) Name() string        { return "git" }
func (t *GitTool) Description() string { return "Inspect Git status, diffs, log, or commit changes." }
func (t *GitTool) Risk() security.Risk { return security.RiskMedium }

func (t *GitTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"action": map[string]any{
				"type":        "string",
				"enum":        []string{"status", "diff", "log", "branch", "commit"},
				"description": "Git action to perform.",
			},
			"args": map[string]any{
				"type":        "string",
				"description": "Optional arguments or commit message.",
			},
		},
		"required": []string{"action"},
	}
}

func (t *GitTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Action string `json:"action"`
		Args   string `json:"args"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid parameters"}
	}

	var cmdArgs []string
	switch input.Action {
	case "status":
		cmdArgs = []string{"status", "--short"}
	case "diff":
		cmdArgs = []string{"diff"}
		if input.Args != "" {
			cmdArgs = append(cmdArgs, strings.Fields(input.Args)...)
		}
	case "log":
		cmdArgs = []string{"log", "-n", "10", "--oneline"}
	case "branch":
		cmdArgs = []string{"branch", "-a"}
	case "commit":
		if input.Args == "" {
			return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "commit requires a message in 'args'"}
		}
		cmdArgs = []string{"commit", "-m", input.Args}
	default:
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("unsupported action %q", input.Action)}
	}

	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	cmd.Dir = t.ws.RootPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Output: output, Error: err.Error()}
	}

	return tools.ToolResult{ToolName: t.Name(), Success: true, Output: output}
}
