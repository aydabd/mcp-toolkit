// Package envvar provides centralized environment variable definitions.
// This is the single source of truth for all environment variables used by mcp-toolkit.
// It integrates with Viper for value resolution from env vars, config files, and flags.
package envvar

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Var represents an environment variable definition.
type Var struct {
	Name        string // Environment variable name
	Key         string // Viper config key (derived from Name if empty)
	Description string // Human-readable description
	Default     string // Default value if not set
	Required    bool   // Whether this variable is required
	Secret      bool   // Whether this is a sensitive value (hidden in help)
}

// key returns the Viper config key for this variable.
func (v Var) key() string {
	if v.Key != "" {
		return v.Key
	}
	// Convert MCP_ATLASSIAN_URL -> atlassian.url
	name := strings.TrimPrefix(v.Name, "MCP_")
	name = strings.ToLower(name)
	parts := strings.SplitN(name, "_", 2)
	if len(parts) == 2 {
		return parts[0] + "." + strings.ReplaceAll(parts[1], "_", ".")
	}
	return name
}

// Get returns the value from Viper (env var, config file, or default).
func (v Var) Get() string {
	val := viper.GetString(v.key())
	if val != "" {
		return val
	}
	return v.Default
}

// GetOrError returns the value or an error if required and not set.
func (v Var) GetOrError() (string, error) {
	val := v.Get()
	if val == "" && v.Required {
		return "", fmt.Errorf("required config %s (env: %s) is not set", v.key(), v.Name)
	}
	return val, nil
}

// IsSet returns true if the variable has a value (from any source).
func (v Var) IsSet() bool {
	return viper.IsSet(v.key()) || v.Default != ""
}

// Set sets the value in Viper.
func (v Var) Set(value string) {
	viper.Set(v.key(), value)
}

// All environment variables used by mcp-toolkit.
var (
	// Atlassian configuration
	AtlassianURL = Var{
		Name:        "MCP_ATLASSIAN_URL",
		Description: "Atlassian instance URL (e.g., https://company.atlassian.net)",
		Required:    true,
	}
	AtlassianEmail = Var{
		Name:        "MCP_ATLASSIAN_EMAIL",
		Description: "Atlassian account email",
		Required:    true,
	}
	AtlassianAPIToken = Var{
		Name:        "MCP_ATLASSIAN_API_TOKEN",
		Description: "Atlassian API token",
		Required:    true,
		Secret:      true,
	}
	AtlassianJiraProjects = Var{
		Name:        "MCP_ATLASSIAN_JIRA_PROJECTS",
		Description: "Comma-separated list of Jira project keys to filter",
	}
	AtlassianReadOnly = Var{
		Name:        "MCP_ATLASSIAN_READ_ONLY",
		Description: "Enable read-only mode (true/false)",
		Default:     "false",
	}

	// Vault configuration
	VaultAddr = Var{
		Name:        "MCP_VAULT_ADDR",
		Description: "HashiCorp Vault server address",
		Required:    true,
	}
	VaultToken = Var{
		Name:        "MCP_VAULT_TOKEN",
		Description: "Vault authentication token",
		Required:    true,
		Secret:      true,
	}
	VaultNamespace = Var{
		Name:        "MCP_VAULT_NAMESPACE",
		Description: "Vault namespace (optional)",
	}
	VaultBinaryPath = Var{
		Name:        "MCP_VAULT_BINARY_PATH",
		Description: "Path to vault-mcp-server binary",
		Default:     "/usr/local/bin/vault-mcp-server",
	}

	// Supabase configuration
	SupabaseProjectRef = Var{
		Name:        "MCP_SUPABASE_PROJECT_REF",
		Description: "Supabase project reference ID",
		Required:    true,
	}

	// General configuration
	ConfigDir = Var{
		Name:        "MCP_CONFIG_DIR",
		Description: "Directory for MCP configuration files",
		Default:     "~/.mcp-server-envs",
	}
	ContainerRuntime = Var{
		Name:        "MCP_CONTAINER_RUNTIME",
		Description: "Container runtime to use (docker, podman, nerdctl)",
		Default:     "docker",
	}
	Editor = Var{
		Name:        "MCP_EDITOR",
		Description: "Target editor(s): vscode, cursor, windsurf, zed, or 'all'",
		Default:     "vscode",
	}
	Verbose = Var{
		Name:        "MCP_VERBOSE",
		Key:         "verbose",
		Description: "Enable verbose logging (true/false)",
		Default:     "false",
	}
)

// BindAll binds all environment variables to Viper.
// Call this during initialization after viper.AutomaticEnv().
func BindAll() {
	viper.SetEnvPrefix("MCP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults in Viper
	for _, v := range All() {
		if v.Default != "" {
			viper.SetDefault(v.key(), v.Default)
		}
	}
}

// All returns all defined environment variables.
func All() []Var {
	return []Var{
		// Atlassian
		AtlassianURL,
		AtlassianEmail,
		AtlassianAPIToken,
		AtlassianJiraProjects,
		AtlassianReadOnly,
		// Vault
		VaultAddr,
		VaultToken,
		VaultNamespace,
		VaultBinaryPath,
		// Supabase
		SupabaseProjectRef,
		// General
		ConfigDir,
		ContainerRuntime,
		Editor,
		Verbose,
	}
}

// ByCategory returns environment variables grouped by category.
func ByCategory() map[string][]Var {
	return map[string][]Var{
		"Atlassian": {
			AtlassianURL,
			AtlassianEmail,
			AtlassianAPIToken,
			AtlassianJiraProjects,
			AtlassianReadOnly,
		},
		"Vault": {
			VaultAddr,
			VaultToken,
			VaultNamespace,
			VaultBinaryPath,
		},
		"Supabase": {
			SupabaseProjectRef,
		},
		"General": {
			ConfigDir,
			ContainerRuntime,
			Editor,
			Verbose,
		},
	}
}

// FormatHelp returns a formatted string for CLI help output.
func FormatHelp() string {
	var sb strings.Builder
	sb.WriteString("Environment Variables:\n")

	for category, vars := range ByCategory() {
		sb.WriteString(fmt.Sprintf("\n  %s:\n", category))
		for _, v := range vars {
			required := ""
			if v.Required {
				required = " (required)"
			}
			defaultVal := ""
			if v.Default != "" && !v.Secret {
				defaultVal = fmt.Sprintf(" [default: %s]", v.Default)
			}
			sb.WriteString(fmt.Sprintf("    %-35s %s%s%s\n", v.Name, v.Description, required, defaultVal))
		}
	}
	return sb.String()
}

// Validate checks that all required environment variables are set.
func Validate(vars ...Var) error {
	var missing []string
	for _, v := range vars {
		if v.Required && v.Get() == "" {
			missing = append(missing, v.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}
	return nil
}
