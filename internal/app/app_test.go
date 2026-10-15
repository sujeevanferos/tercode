package app

import (
	"context"
	"testing"

	"github.com/sujeevanferos/tercode/internal/security"
)

func TestDefaultToolPermissions(t *testing.T) {
	tests := []struct {
		name    string
		allowed bool
	}{
		{"read_file", true},
		{"list_directory", true},
		{"write_file", false},
		{"edit_file", false},
		{"shell", false},
		{"git", false},
		{"search", false},
		{"test_build", false},
		{"unregistered_tool", false},
	}

	pm := newToolPermissions()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pm.Check(context.Background(), security.PermissionRequest{
				ToolName: tc.name,
			})
			if err != nil {
				t.Fatalf("permission check failed: %v", err)
			}
			if got.Granted != tc.allowed {
				t.Errorf("tool %q granted = %t, want %t (reason: %s)", tc.name, got.Granted, tc.allowed, got.Reason)
			}
		})
	}
}
