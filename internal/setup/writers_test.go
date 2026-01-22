package setup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteAtlassianEnv(t *testing.T) {
	// Ensure env directory exists
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &AtlassianConfig{
		URL:          "https://test.atlassian.net",
		Email:        "user@test.com",
		APIToken:     "secret-token",
		JiraProjects: "PROJ,DEV",
		ReadOnly:     false,
	}

	err := WriteAtlassianEnv(cfg)
	require.NoError(t, err)

	envPath := filepath.Join(envDir(), ".mcp-atlassian.env")
	data, err := os.ReadFile(envPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "CONFLUENCE_URL=https://test.atlassian.net/wiki")
	assert.Contains(t, content, "CONFLUENCE_USERNAME=user@test.com")
	assert.Contains(t, content, "JIRA_PROJECTS_FILTER=PROJ,DEV")
	assert.Contains(t, content, "READ_ONLY_MODE=false")

	_ = os.Remove(envPath)
}

func TestWriteAtlassianEnvReadOnly(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &AtlassianConfig{
		URL:      "https://test.atlassian.net",
		Email:    "user@test.com",
		APIToken: "token",
		ReadOnly: true,
	}

	err := WriteAtlassianEnv(cfg)
	require.NoError(t, err)

	envPath := filepath.Join(envDir(), ".mcp-atlassian.env")
	data, err := os.ReadFile(envPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "READ_ONLY_MODE=true")

	_ = os.Remove(envPath)
}

func TestWriteVaultEnv(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &VaultConfig{
		Address:   "https://vault.test.com",
		Token:     "hvs.token",
		Namespace: "admin",
	}

	err := WriteVaultEnv(cfg)
	require.NoError(t, err)

	data, err := os.ReadFile(vaultEnvPath())
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "VAULT_ADDR=https://vault.test.com")
	assert.Contains(t, content, "VAULT_TOKEN=hvs.token")
	assert.Contains(t, content, "VAULT_NAMESPACE=admin")

	_ = os.Remove(vaultEnvPath())
}

func TestWriteVaultWrapper(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &VaultConfig{
		Address:    "https://vault.test.com",
		Token:      "token",
		BinaryPath: "/usr/bin/vault-mcp-server",
	}

	err := WriteVaultWrapper(cfg)
	require.NoError(t, err)

	data, err := os.ReadFile(vaultWrapperPath())
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "#!/bin/bash")
	assert.Contains(t, content, "/usr/bin/vault-mcp-server")

	info, err := os.Stat(vaultWrapperPath())
	require.NoError(t, err)
	// Check that execute bit is set (may vary by umask)
	assert.True(t, info.Mode().Perm()&0100 != 0, "owner execute bit should be set")

	_ = os.Remove(vaultWrapperPath())
}

func TestWriteKubernetesEnv(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &KubernetesConfig{}
	err := WriteKubernetesEnv(cfg)
	require.NoError(t, err)

	envPath := filepath.Join(envDir(), ".mcp-kubernetes.env")
	data, err := os.ReadFile(envPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "# Kubernetes MCP Server")
	assert.Contains(t, content, "KUBECONFIG")

	_ = os.Remove(envPath)
}

func TestWriteGitHubEnv(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &GitHubConfig{}
	err := WriteGitHubEnv(cfg)
	require.NoError(t, err)

	envPath := filepath.Join(envDir(), ".mcp-github.env")
	data, err := os.ReadFile(envPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "# GitHub MCP Server")
	assert.Contains(t, content, "GitHub Copilot")

	_ = os.Remove(envPath)
}

func TestWriteSupabaseEnv(t *testing.T) {
	require.NoError(t, os.MkdirAll(envDir(), 0700))

	cfg := &SupabaseConfig{
		ProjectRef: "abc123xyz",
	}
	err := WriteSupabaseEnv(cfg)
	require.NoError(t, err)

	envPath := filepath.Join(envDir(), ".mcp-supabase.env")
	data, err := os.ReadFile(envPath)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "# Supabase MCP Server")
	assert.Contains(t, content, "SUPABASE_PROJECT_REF=abc123xyz")

	_ = os.Remove(envPath)
}
