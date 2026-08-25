// Package filesystem implements sandboxed file operations (read, write, edit, delete, list).
package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sujeevanferos/tercode/internal/security"
	"github.com/sujeevanferos/tercode/internal/tools"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

// ReadFileTool reads the contents of a file within the workspace.
type ReadFileTool struct {
	ws *workspace.Workspace
}

func NewReadFileTool(ws *workspace.Workspace) *ReadFileTool {
	return &ReadFileTool{ws: ws}
}

func (t *ReadFileTool) Name() string        { return "read_file" }
func (t *ReadFileTool) Description() string { return "Read the text contents of a file within the workspace." }
func (t *ReadFileTool) Risk() security.Risk { return security.RiskLow }

func (t *ReadFileTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file relative to the workspace root.",
			},
		},
		"required": []string{"path"},
	}
}

func (t *ReadFileTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid input parameters"}
	}

	target, err := t.ws.ResolvePath(input.Path)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: err.Error()}
	}

	content, err := os.ReadFile(target)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("read error: %v", err)}
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   string(content),
	}
}

// WriteFileTool creates or overwrites a file within the workspace.
type WriteFileTool struct {
	ws *workspace.Workspace
}

func NewWriteFileTool(ws *workspace.Workspace) *WriteFileTool {
	return &WriteFileTool{ws: ws}
}

func (t *WriteFileTool) Name() string        { return "write_file" }
func (t *WriteFileTool) Description() string { return "Create or overwrite a file within the workspace." }
func (t *WriteFileTool) Risk() security.Risk { return security.RiskMedium }

func (t *WriteFileTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file relative to workspace root.",
			},
			"content": map[string]any{
				"type":        "string",
				"description": "Text content to write.",
			},
		},
		"required": []string{"path", "content"},
	}
}

func (t *WriteFileTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid input parameters"}
	}

	target, err := t.ws.ResolvePath(input.Path)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: err.Error()}
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("mkdir error: %v", err)}
	}

	if err := os.WriteFile(target, []byte(input.Content), 0o644); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("write error: %v", err)}
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   fmt.Sprintf("Successfully wrote %d bytes to %s", len(input.Content), input.Path),
	}
}

// EditFileTool performs search and replace on a target file.
type EditFileTool struct {
	ws *workspace.Workspace
}

func NewEditFileTool(ws *workspace.Workspace) *EditFileTool {
	return &EditFileTool{ws: ws}
}

func (t *EditFileTool) Name() string        { return "edit_file" }
func (t *EditFileTool) Description() string { return "Replace a specific substring or block of lines in a file with new content." }
func (t *EditFileTool) Risk() security.Risk { return security.RiskMedium }

func (t *EditFileTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to the file to edit.",
			},
			"target": map[string]any{
				"type":        "string",
				"description": "Exact text substring/block to replace.",
			},
			"replacement": map[string]any{
				"type":        "string",
				"description": "New content to replace the target text.",
			},
		},
		"required": []string{"path", "target", "replacement"},
	}
}

func (t *EditFileTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Path        string `json:"path"`
		Target      string `json:"target"`
		Replacement string `json:"replacement"`
	}
	if err := json.Unmarshal(raw, &input); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "invalid parameters"}
	}

	targetFile, err := t.ws.ResolvePath(input.Path)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: err.Error()}
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("read error: %v", err)}
	}

	s := string(content)
	if !strings.Contains(s, input.Target) {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: "target text was not found in the file"}
	}

	newContent := strings.Replace(s, input.Target, input.Replacement, 1)
	if err := os.WriteFile(targetFile, []byte(newContent), 0o644); err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("write error: %v", err)}
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   fmt.Sprintf("Successfully edited %s", input.Path),
	}
}

// ListDirectoryTool lists files and subdirectories.
type ListDirectoryTool struct {
	ws *workspace.Workspace
}

func NewListDirectoryTool(ws *workspace.Workspace) *ListDirectoryTool {
	return &ListDirectoryTool{ws: ws}
}

func (t *ListDirectoryTool) Name() string        { return "list_directory" }
func (t *ListDirectoryTool) Description() string { return "List files and subdirectories at a path relative to the workspace." }
func (t *ListDirectoryTool) Risk() security.Risk { return security.RiskLow }

func (t *ListDirectoryTool) Parameters() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Directory path relative to workspace (empty for root).",
			},
		},
	}
}

func (t *ListDirectoryTool) Execute(ctx context.Context, raw json.RawMessage) tools.ToolResult {
	var input struct {
		Path string `json:"path"`
	}
	_ = json.Unmarshal(raw, &input)

	targetDir, err := t.ws.ResolvePath(input.Path)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: err.Error()}
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return tools.ToolResult{ToolName: t.Name(), Success: false, Error: fmt.Sprintf("read dir error: %v", err)}
	}

	var sb strings.Builder
	for _, entry := range entries {
		kind := "FILE"
		if entry.IsDir() {
			kind = "DIR "
		}
		sb.WriteString(fmt.Sprintf("%s  %s\n", kind, entry.Name()))
	}

	return tools.ToolResult{
		ToolName: t.Name(),
		Success:  true,
		Output:   sb.String(),
	}
}
