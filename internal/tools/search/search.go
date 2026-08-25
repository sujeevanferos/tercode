// Package search implements workspace search using ripgrep or standard directory walking.
package search

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// SearchTool executes regex search across workspace files.
type SearchTool struct {
	ws *workspace.Workspace
}

func New(ws *workspace.Workspace) *SearchTool {
	return &SearchTool{ws: ws}
}

func (t *SearchTool) Name() string        { return "search" }
func (t *SearchTool) Description() string { return "Search for a pattern or keyword across files in the workspace (using ripgrep if installed, or git grep)." }
func (t *SearchTool) Risk() security.Risk { return security.RiskLow }

func (t *SearchTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"query": map[string]any{
				"type":        "string",
				"description": "Regex pattern or exact text to search for.",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Optional subdirectory or file glob to filter.",
			},
		},
		"required": []string{"query"},
	}
}

func (t *SearchTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Query string `json:"query"`
		Path  string `json:"path"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid parameters"}
	}

	// Try ripgrep (rg) first, fallback to git grep
	args := []string{"-n", "--max-count=50", input.Query}
	if input.Path != "" {
		args = append(args, input.Path)
	}

	cmd := exec.CommandContext(ctx, "rg", args...)
	cmd.Dir = t.ws.RootPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Fallback to git grep
		gitArgs := []string{"grep", "-n", "-I", input.Query}
		if input.Path != "" {
			gitArgs = append(gitArgs, "--", input.Path)
		}
		cmd2 := exec.CommandContext(ctx, "git", gitArgs...)
		cmd2.Dir = t.ws.RootPath
		stdout.Reset()
		stderr.Reset()
		cmd2.Stdout = &stdout
		cmd2.Stderr = &stderr
		err = cmd2.Run()
	}

	output := stdout.String()
	if output == "" {
		output = "No matches found."
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   output,
	}
}
