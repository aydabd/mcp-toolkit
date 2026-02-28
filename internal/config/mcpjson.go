package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcp-toolkit/internal/editor"
)

// MCPConfig represents the mcp.json structure.
type MCPConfig struct {
	Inputs  []any                   `json:"inputs"`
	Servers map[string]ServerConfig `json:"servers"`
}

// ServerConfig represents an MCP server entry.
type ServerConfig struct {
	Type    string   `json:"type"`
	URL     string   `json:"url,omitempty"`
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// SupabaseServerConfig is used to extract project ref from supabase config.
type SupabaseServerConfig interface {
	GetProjectRef() string
}

// GenerateMCPJSON creates or updates mcp.json for the default editor (VS Code).
func GenerateMCPJSON(serverConfigs map[string]any) error {
	return GenerateMCPJSONForEditor(editor.VSCode, serverConfigs)
}

// GenerateMCPJSONForEditor creates or updates mcp.json for a specific editor.
func GenerateMCPJSONForEditor(e editor.Editor, serverConfigs map[string]any) error {
	path := e.ConfigPath()
	if path == "" {
		return fmt.Errorf("unsupported editor: %s", e)
	}
	return GenerateMCPJSONToPath(path, serverConfigs)
}

// GenerateMCPJSONToPath creates or updates mcp.json at a specific path.
func GenerateMCPJSONToPath(path string, serverConfigs map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil { //nolint:gosec // editor config dir
		return fmt.Errorf("create config dir: %w", err)
	}

	cfg := &MCPConfig{Inputs: []any{}, Servers: make(map[string]ServerConfig)}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, cfg)
	}

	home := homeDir()
	for name, serverCfg := range serverConfigs {
		switch name {
		case "atlassian":
			cfg.Servers["atlassian"] = ServerConfig{
				Type:    "stdio",
				Command: "docker",
				Args:    []string{"run", "--rm", "-i", "--env-file", AtlassianEnvPath(), "mcp/atlassian"},
			}
		case "kubernetes":
			cfg.Servers["kubernetes"] = ServerConfig{
				Type:    "stdio",
				Command: "docker",
				Args:    []string{"run", "--rm", "-i", "-v", home + "/.kube:/root/.kube:ro", "mcp/kubernetes"},
			}
		case "vault":
			cfg.Servers["vault-mcp-server"] = ServerConfig{
				Type:    "stdio",
				Command: VaultWrapperPath(),
				Args:    []string{},
			}
		case "github":
			cfg.Servers["github"] = ServerConfig{
				Type: "http",
				URL:  "https://api.githubcopilot.com/mcp/",
			}
		case "supabase":
			// Extract project ref using type assertion with interface
			if sc, ok := serverCfg.(SupabaseServerConfig); ok {
				cfg.Servers["supabase"] = ServerConfig{
					Type: "http",
					URL:  fmt.Sprintf("https://mcp.supabase.com/mcp?project_ref=%s", sc.GetProjectRef()),
				}
			} else if m, ok := serverCfg.(map[string]string); ok {
				if ref, exists := m["projectRef"]; exists {
					cfg.Servers["supabase"] = ServerConfig{
						Type: "http",
						URL:  fmt.Sprintf("https://mcp.supabase.com/mcp?project_ref=%s", ref),
					}
				}
			}
		}
	}

	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0600) //nolint:gosec // mcp.json needs user read for editor
}

// GenerateMCPJSONForEditors creates or updates mcp.json for multiple editors.
func GenerateMCPJSONForEditors(editors []editor.Editor, serverConfigs map[string]any) error {
	for _, e := range editors {
		if err := GenerateMCPJSONForEditor(e, serverConfigs); err != nil {
			return fmt.Errorf("%s: %w", e.String(), err)
		}
	}
	return nil
}
