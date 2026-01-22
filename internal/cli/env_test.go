package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	envCmd.SetOut(buf)
	envCmd.SetArgs([]string{})

	err := envCmd.Execute()
	assert.NoError(t, err)
}

func TestShowEnvVars(t *testing.T) {
	// Just verify it doesn't panic
	showEnvVars()
}

func TestShowEnvVarsWithValues(t *testing.T) {
	// Set some env vars to test different branches
	os.Setenv("MCP_ATLASSIAN_URL", "https://test.atlassian.net")
	os.Setenv("MCP_ATLASSIAN_API_TOKEN", "secret-token")
	// Set to default value to cover the "default" branch
	os.Setenv("MCP_VAULT_BINARY", "/usr/local/bin/vault-mcp-server")
	defer os.Unsetenv("MCP_ATLASSIAN_URL")
	defer os.Unsetenv("MCP_ATLASSIAN_API_TOKEN")
	defer os.Unsetenv("MCP_VAULT_BINARY")

	// Should cover the secret masking and value display branches
	showEnvVars()
}
