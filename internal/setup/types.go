// Package setup handles MCP server environment configuration.
package setup

// AtlassianConfig holds Atlassian MCP server settings.
type AtlassianConfig struct {
	URL          string
	Email        string
	APIToken     string
	JiraProjects string
	ReadOnly     bool
}

// VaultConfig holds Vault MCP server settings.
type VaultConfig struct {
	Address    string
	Token      string
	Namespace  string
	BinaryPath string
}

// KubernetesConfig holds Kubernetes MCP server settings.
type KubernetesConfig struct{}

// GitHubConfig holds GitHub MCP server settings.
type GitHubConfig struct{}

// SupabaseConfig holds Supabase MCP server settings.
type SupabaseConfig struct {
	ProjectRef string
}

// GetProjectRef implements the SupabaseServerConfig interface.
func (c *SupabaseConfig) GetProjectRef() string {
	return c.ProjectRef
}
