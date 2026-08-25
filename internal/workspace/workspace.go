// Package workspace encapsulates the active workspace (project root, sandboxed filesystem, Git state).
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Workspace represents an active project directory.
type Workspace struct {
	ID       string
	Name     string
	RootPath string
}

// Detect detects the workspace root by searching upwards for .git, .config, or go.mod.
// If none found, uses the current working directory.
func Detect(startDir string) (*Workspace, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("workspace: getwd: %w", err)
		}
	}

	curr := startDir
	for {
		if hasIndicator(curr) {
			return &Workspace{
				ID:       filepath.Base(curr),
				Name:     filepath.Base(curr),
				RootPath: curr,
			}, nil
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	// Fallback to startDir
	return &Workspace{
		ID:       filepath.Base(startDir),
		Name:     filepath.Base(startDir),
		RootPath: startDir,
	}, nil
}

func hasIndicator(dir string) bool {
	indicators := []string{".git", ".config", "go.mod", "package.json", "Cargo.toml", "pyproject.toml"}
	for _, ind := range indicators {
		if _, err := os.Stat(filepath.Join(dir, ind)); err == nil {
			return true
		}
	}
	return false
}

// ResolvePath ensures a path is safely within the workspace boundary.
func (w *Workspace) ResolvePath(relOrAbs string) (string, error) {
	var target string
	if filepath.IsAbs(relOrAbs) {
		target = filepath.Clean(relOrAbs)
	} else {
		target = filepath.Clean(filepath.Join(w.RootPath, relOrAbs))
	}

	rel, err := filepath.Rel(w.RootPath, target)
	if err != nil || len(rel) >= 2 && rel[:2] == ".." {
		return "", fmt.Errorf("path %q escapes workspace root %q", relOrAbs, w.RootPath)
	}
	return target, nil
}
