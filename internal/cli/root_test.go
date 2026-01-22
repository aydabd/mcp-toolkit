package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetVersion(t *testing.T) {
	SetVersion("1.0.0", "abc123", "2026-01-01")
	assert.Equal(t, "1.0.0", appVersion)
	assert.Equal(t, "abc123", appCommit)
	assert.Equal(t, "2026-01-01", appBuildTime)
}

func TestRootCommand(t *testing.T) {
	assert.Equal(t, "mcp-toolkit", rootCmd.Use)
	assert.NotEmpty(t, rootCmd.Short)
}

func TestInitConfig(t *testing.T) {
	// Test default config path
	initConfig()
}

func TestInitConfigWithFile(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "test-config.yaml")
	os.WriteFile(cfgPath, []byte("test: value"), 0644)

	oldCfgFile := cfgFile
	cfgFile = cfgPath
	defer func() { cfgFile = oldCfgFile }()

	initConfig()
}

func TestInitLogger(t *testing.T) {
	// Test default level
	verbose = false
	initLogger()

	// Test verbose level
	verbose = true
	initLogger()
	verbose = false
}

func TestExecute(t *testing.T) {
	// Execute runs the cobra command - just verify it doesn't panic
	// We can't fully test it without capturing os.Exit
	rootCmd.SetArgs([]string{"--help"})
	err := Execute()
	assert.NoError(t, err)
}
