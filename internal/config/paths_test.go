package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".mcp-server-envs")
	assert.Equal(t, expected, EnvDir())
}

func TestKubeconfigPath(t *testing.T) {
	path := KubeconfigPath()
	assert.Contains(t, path, ".kube")
	assert.Contains(t, path, "config")
}

func TestAtlassianEnvPath(t *testing.T) {
	path := AtlassianEnvPath()
	assert.Contains(t, path, ".mcp-atlassian.env")
}

func TestVaultEnvPath(t *testing.T) {
	path := VaultEnvPath()
	assert.Contains(t, path, ".mcp-vault.env")
}

func TestVaultWrapperPath(t *testing.T) {
	path := VaultWrapperPath()
	assert.Contains(t, path, "vault-mcp-wrapper.sh")
}

func TestEnsureEnvDirectory(t *testing.T) {
	// This creates real directory, so just test it doesn't error
	err := EnsureEnvDirectory()
	assert.NoError(t, err)
	assert.DirExists(t, EnvDir())
}

func TestFileExists(t *testing.T) {
	assert.False(t, FileExists("/nonexistent/path"))

	f, _ := os.CreateTemp("", "test")
	defer os.Remove(f.Name())
	assert.True(t, FileExists(f.Name()))
}

func TestKubernetesEnvPath(t *testing.T) {
	path := KubernetesEnvPath()
	assert.Contains(t, path, ".mcp-kubernetes.env")
	assert.Contains(t, path, ".mcp-server-envs")
}

func TestGitHubEnvPath(t *testing.T) {
	path := GitHubEnvPath()
	assert.Contains(t, path, ".mcp-github.env")
	assert.Contains(t, path, ".mcp-server-envs")
}

func TestSupabaseEnvPath(t *testing.T) {
	path := SupabaseEnvPath()
	assert.Contains(t, path, ".mcp-supabase.env")
	assert.Contains(t, path, ".mcp-server-envs")
}
