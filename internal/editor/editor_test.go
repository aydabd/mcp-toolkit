package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEditorString(t *testing.T) {
	tests := []struct {
		editor Editor
		want   string
	}{
		{VSCode, "VS Code"},
		{Cursor, "Cursor"},
		{Windsurf, "Windsurf"},
		{Zed, "Zed"},
		{Editor("unknown"), "unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.editor.String())
	}
}

func TestAll(t *testing.T) {
	all := All()
	assert.Len(t, all, 4)
	assert.Contains(t, all, VSCode)
	assert.Contains(t, all, Cursor)
	assert.Contains(t, all, Windsurf)
	assert.Contains(t, all, Zed)
}

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  Editor
		ok    bool
	}{
		{"vscode", VSCode, true},
		{"code", VSCode, true},
		{"vs-code", VSCode, true},
		{"cursor", Cursor, true},
		{"windsurf", Windsurf, true},
		{"zed", Zed, true},
		{"unknown", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		got, ok := Parse(tt.input)
		assert.Equal(t, tt.ok, ok, "input: %s", tt.input)
		if ok {
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestConfigPath(t *testing.T) {
	// Just verify paths are non-empty for known editors
	for _, e := range All() {
		path := e.ConfigPath()
		assert.NotEmpty(t, path, "editor: %s", e)
		assert.True(t, filepath.IsAbs(path) || path[0] == '~', "path should be absolute: %s", path)
	}

	// Unknown editor returns empty
	unknown := Editor("unknown")
	assert.Empty(t, unknown.ConfigPath())
}

func TestIsInstalled(t *testing.T) {
	// Create a temp dir to simulate installed editor
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create VS Code config dir (darwin path)
	vscodeDir := filepath.Join(tmpDir, "Library", "Application Support", "Code", "User")
	os.MkdirAll(vscodeDir, 0755)

	// Test installed editor
	assert.True(t, VSCode.IsInstalled())

	// Test not installed editor
	assert.False(t, Cursor.IsInstalled())
}

func TestIsInstalledUnknownEditor(t *testing.T) {
	unknown := Editor("unknown")
	assert.False(t, unknown.IsInstalled())
}

func TestInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create one editor config dir
	vscodeDir := filepath.Join(tmpDir, "Library", "Application Support", "Code", "User")
	os.MkdirAll(vscodeDir, 0755)

	installed := Installed()
	assert.NotNil(t, installed)
	assert.Contains(t, installed, VSCode)
}

func TestHomeDir(t *testing.T) {
	home := homeDir()
	assert.NotEmpty(t, home)
}
