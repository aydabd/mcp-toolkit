// Package editor provides MCP configuration for multiple editors/IDEs.
package editor

import (
	"os"
	"path/filepath"
)

// Editor represents a supported editor/IDE.
type Editor string

const (
	VSCode   Editor = "vscode"
	Cursor   Editor = "cursor"
	Windsurf Editor = "windsurf"
	Zed      Editor = "zed"
)

// All returns all supported editors.
func All() []Editor {
	return []Editor{VSCode, Cursor, Windsurf, Zed}
}

// String returns the display name.
func (e Editor) String() string {
	switch e {
	case VSCode:
		return "VS Code"
	case Cursor:
		return "Cursor"
	case Windsurf:
		return "Windsurf"
	case Zed:
		return "Zed"
	default:
		return string(e)
	}
}

// ConfigPath returns the MCP config file path for this editor.
// Uses osConfigPath which is defined in OS-specific files.
func (e Editor) ConfigPath() string {
	return e.osConfigPath(homeDir())
}

// IsInstalled checks if the editor appears to be installed.
func (e Editor) IsInstalled() bool {
	path := e.ConfigPath()
	if path == "" {
		return false
	}
	// Check if config directory exists
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	return err == nil
}

// Parse converts a string to Editor.
func Parse(s string) (Editor, bool) {
	switch s {
	case "vscode", "code", "vs-code":
		return VSCode, true
	case "cursor":
		return Cursor, true
	case "windsurf":
		return Windsurf, true
	case "zed":
		return Zed, true
	default:
		return "", false
	}
}

// Installed returns all editors that appear to be installed.
func Installed() []Editor {
	var installed []Editor
	for _, e := range All() {
		if e.IsInstalled() {
			installed = append(installed, e)
		}
	}
	return installed
}

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}
