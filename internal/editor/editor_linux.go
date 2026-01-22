//go:build linux

package editor

import (
	"path/filepath"
)

func (e Editor) osConfigPath(home string) string {
	switch e {
	case VSCode:
		return filepath.Join(home, ".config", "Code", "User", "mcp.json")
	case Cursor:
		return filepath.Join(home, ".config", "Cursor", "User", "mcp.json")
	case Windsurf:
		return filepath.Join(home, ".config", "Windsurf", "User", "mcp.json")
	case Zed:
		return filepath.Join(home, ".config", "zed", "settings.json")
	default:
		return ""
	}
}
