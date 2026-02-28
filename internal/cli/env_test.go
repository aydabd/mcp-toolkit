package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	envCmd.SetOut(buf)
	envCmd.SetArgs([]string{})

	err := envCmd.Execute()
	assert.NoError(t, err)
}

func TestShowEnvVars(t *testing.T) {
	t.Helper()
	// Just verify it doesn't panic
	showEnvVars()
}

func TestShowEnvVarsWithValues(t *testing.T) {
	t.Helper()
	// Set some env vars to test different branches
	require.NoError(t, os.Setenv("MCP_ATLASSIAN_URL", "https://test.atlassian.net"))
	require.NoError(t, os.Setenv("MCP_ATLASSIAN_API_TOKEN", "secret-token"))
	// Set to default value to cover the "default" branch
	require.NoError(t, os.Setenv("MCP_VAULT_BINARY", "/usr/local/bin/vault-mcp-server"))
	t.Cleanup(func() {
		_ = os.Unsetenv("MCP_ATLASSIAN_URL")
		_ = os.Unsetenv("MCP_ATLASSIAN_API_TOKEN")
		_ = os.Unsetenv("MCP_VAULT_BINARY")
	})

	// Should cover the secret masking and value display branches
	showEnvVars()
}
