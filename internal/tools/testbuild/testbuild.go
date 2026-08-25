// Package testbuild runs workspace test suites and builds.
package testbuild

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// TestBuildTool runs tests and builds.
type TestBuildTool struct {
	ws *workspace.Workspace
}

func New(ws *workspace.Workspace) *TestBuildTool {
	return &TestBuildTool{ws: ws}
}

func (t *TestBuildTool) Name() string        { return "run_tests" }
func (t *TestBuildTool) Description() string { return "Run the project test suite (e.g. go test ./..., npm test, cargo test)." }
func (t *TestBuildTool) Risk() security.Risk { return security.RiskMedium }

func (t *TestBuildTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "Optional specific test command to run. If omitted, uses standard defaults.",
			},
		},
	}
}

func (t *TestBuildTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Command string `json:"command"`
	}
	_ = json.Unmarshal(raw, &input)

	cmdStr := input.Command
	if cmdStr == "" {
		cmdStr = "go test ./..."
	}

	execCtx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	cmd := exec.CommandContext(execCtx, "sh", "-c", cmdStr)
	cmd.Dir = t.ws.RootPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	if stderr.Len() > 0 {
		output += "\n" + stderr.String()
	}

	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Output: output, Error: fmt.Sprintf("tests failed: %v", err)}
	}

	return tools.ToolResult{ToolName: t.Name(), Success: true, Output: output}
}
