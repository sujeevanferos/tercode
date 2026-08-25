package filesystem_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/sujeevanferos/tercode/internal/tools/filesystem"
	"github.com/sujeevanferos/tercode/internal/workspace"
)

func TestFilesystemTools(t *testing.T) {
	tempDir := t.TempDir()
	ws := &workspace.Workspace{
		ID:       "test",
		Name:     "test",
		RootPath: tempDir,
	}

	writeTool := filesystem.NewWriteFileTool(ws)
	readTool := filesystem.NewReadFileTool(ws)
	editTool := filesystem.NewEditFileTool(ws)

	ctx := context.Background()

	// 1. Write file
	writePayload, _ := json.Marshal(map[string]string{
		"path":    "hello.txt",
		"content": "Hello World",
	})
	res := writeTool.Execute(ctx, writePayload)
	if !res.Success {
		t.Fatalf("write failed: %s", res.Error)
	}

	// 2. Read file
	readPayload, _ := json.Marshal(map[string]string{"path": "hello.txt"})
	readRes := readTool.Execute(ctx, readPayload)
	if !readRes.Success || readRes.Output != "Hello World" {
		t.Fatalf("expected 'Hello World', got output: %q, error: %s", readRes.Output, readRes.Error)
	}

	// 3. Edit file
	editPayload, _ := json.Marshal(map[string]string{
		"path":        "hello.txt",
		"target":      "World",
		"replacement": "Tercode",
	})
	editRes := editTool.Execute(ctx, editPayload)
	if !editRes.Success {
		t.Fatalf("edit failed: %s", editRes.Error)
	}

	readRes2 := readTool.Execute(ctx, readPayload)
	if readRes2.Output != "Hello Tercode" {
		t.Fatalf("expected 'Hello Tercode', got %q", readRes2.Output)
	}

	// 4. Sandbox violation test (path escaping workspace root)
	escapePayload, _ := json.Marshal(map[string]string{"path": "../secret.txt"})
	escapeRes := readTool.Execute(ctx, escapePayload)
	if escapeRes.Success {
		t.Fatal("expected path escape to be rejected by workspace sandbox")
	}

	_ = os.RemoveAll(filepath.Join(tempDir, "hello.txt"))
}
