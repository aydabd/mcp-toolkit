package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtlassianConfig(t *testing.T) {
	cfg := AtlassianConfig{
		URL:          "https://example.atlassian.net",
		Email:        "test@example.com",
		APIToken:     "token123",
		JiraProjects: "PROJ",
		ReadOnly:     false,
	}
	assert.Equal(t, "https://example.atlassian.net", cfg.URL)
	assert.Equal(t, "test@example.com", cfg.Email)
	assert.Equal(t, "token123", cfg.APIToken)
	assert.Equal(t, "PROJ", cfg.JiraProjects)
	assert.False(t, cfg.ReadOnly)
}

func TestVaultConfig(t *testing.T) {
	cfg := VaultConfig{
		Address:    "https://vault.example.com",
		Token:      "hvs.xxx",
		Namespace:  "admin",
		BinaryPath: "/usr/bin/vault-mcp",
	}
	assert.Equal(t, "https://vault.example.com", cfg.Address)
	assert.Equal(t, "hvs.xxx", cfg.Token)
	assert.Equal(t, "admin", cfg.Namespace)
	assert.Equal(t, "/usr/bin/vault-mcp", cfg.BinaryPath)
}

func TestKubernetesConfig(t *testing.T) {
	cfg := KubernetesConfig{}
	assert.NotNil(t, cfg)
}

func TestGitHubConfig(t *testing.T) {
	cfg := GitHubConfig{}
	assert.NotNil(t, cfg)
}

func TestSupabaseConfig(t *testing.T) {
	cfg := SupabaseConfig{ProjectRef: "my-project"}
	assert.Equal(t, "my-project", cfg.ProjectRef)
	assert.Equal(t, "my-project", cfg.GetProjectRef())
}
