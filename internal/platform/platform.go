// Package platform provides abstractions over operating system differences.
// All platform-specific behavior (paths, keychains, terminal capabilities,
// file watching) must be accessed through this package rather than being
// scattered across the codebase.
package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// OS returns the current operating system identifier.
func OS() string {
	return runtime.GOOS
}

// ConfigDir returns the platform-appropriate global configuration directory
// for Tercode (e.g. ~/.config/tercode on Linux/macOS, %AppData%\tercode on Windows).
func ConfigDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "tercode"), nil
	default:
		// Linux and macOS: use XDG_CONFIG_HOME or ~/.config
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg != "" {
			return filepath.Join(xdg, "tercode"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "tercode"), nil
	}
}

// DataDir returns the platform-appropriate data directory for Tercode
// (e.g. ~/.local/share/tercode on Linux, %LocalAppData%\tercode on Windows).
func DataDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			localAppData = filepath.Join(home, "AppData", "Local")
		}
		return filepath.Join(localAppData, "tercode"), nil
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", "tercode"), nil
	default:
		xdg := os.Getenv("XDG_DATA_HOME")
		if xdg != "" {
			return filepath.Join(xdg, "tercode"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "tercode"), nil
	}
}

// EnsureDir creates a directory and all parents if they do not exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// HomeDir returns the current user's home directory.
func HomeDir() (string, error) {
	return os.UserHomeDir()
}
