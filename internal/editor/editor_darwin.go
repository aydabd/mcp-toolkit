//go:build darwin

package editor

import (
	"path/filepath"
)

func (e Editor) osConfigPath(home string) string {
	switch e {
	case VSCode:
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "mcp.json")
	case Cursor:
		return filepath.Join(home, "Library", "Application Support", "Cursor", "User", "mcp.json")
	case Windsurf:
		return filepath.Join(home, "Library", "Application Support", "Windsurf", "User", "mcp.json")
	case Zed:
		return filepath.Join(home, ".config", "zed", "settings.json")
	default:
		return ""
	}
}
