//go:build windows

package editor

import (
	"os"
	"path/filepath"
)

func (e Editor) osConfigPath(_ string) string {
	appdata := os.Getenv("APPDATA")
	switch e {
	case VSCode:
		return filepath.Join(appdata, "Code", "User", "mcp.json")
	case Cursor:
		return filepath.Join(appdata, "Cursor", "User", "mcp.json")
	case Windsurf:
		return filepath.Join(appdata, "Windsurf", "User", "mcp.json")
	case Zed:
		return filepath.Join(appdata, "Zed", "settings.json")
	default:
		return ""
	}
}
