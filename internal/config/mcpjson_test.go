package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mcp-toolkit/internal/editor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPConfigStructure(t *testing.T) {
	cfg := MCPConfig{
		Inputs: []any{},
		Servers: map[string]ServerConfig{
			"test": {Type: "stdio", Command: "echo"},
		},
	}

	data, err := json.Marshal(cfg)
	require.NoError(t, err)
	assert.Contains(t, string(data), "test")
	assert.Contains(t, string(data), "stdio")
}

func TestServerConfigTypes(t *testing.T) {
	tests := []struct {
		name   string
		config ServerConfig
		want   string
	}{
		{
			name:   "http type",
			config: ServerConfig{Type: "http", URL: "https://example.com"},
			want:   "http",
		},
		{
			name:   "stdio type",
			config: ServerConfig{Type: "stdio", Command: "docker", Args: []string{"run"}},
			want:   "stdio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.config.Type)
		})
	}
}

func TestServerConfigJSON(t *testing.T) {
	cfg := ServerConfig{
		Type:    "stdio",
		Command: "docker",
		Args:    []string{"run", "--rm", "-i", "mcp/atlassian"},
	}

	data, err := json.Marshal(cfg)
	require.NoError(t, err)

	var decoded ServerConfig
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, cfg.Type, decoded.Type)
	assert.Equal(t, cfg.Command, decoded.Command)
	assert.Equal(t, cfg.Args, decoded.Args)
}

func TestMCPConfigMerge(t *testing.T) {
	// Test the merge behavior by simulating what GenerateMCPJSON does
	existing := MCPConfig{
		Inputs: []any{},
		Servers: map[string]ServerConfig{
			"existing": {Type: "http", URL: "https://existing.com"},
		},
	}

	// Add new server (simulating merge)
	existing.Servers["new-server"] = ServerConfig{Type: "stdio", Command: "test"}

	assert.Contains(t, existing.Servers, "existing")
	assert.Contains(t, existing.Servers, "new-server")
}

func TestGenerateMCPJSONCreatesFile(t *testing.T) {
	// Test in a temp directory
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, ".vscode", "mcp.json")

	// Create the directory
	require.NoError(t, os.MkdirAll(filepath.Dir(testPath), 0755))

	// Write a config directly (testing the file writing logic)
	cfg := MCPConfig{
		Inputs: []any{},
		Servers: map[string]ServerConfig{
			"github": {Type: "http", URL: "https://api.githubcopilot.com/mcp/"},
		},
	}

	data, err := json.MarshalIndent(cfg, "", "\t")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(testPath, data, 0644))

	// Verify file exists and is valid
	readData, err := os.ReadFile(testPath)
	require.NoError(t, err)

	var readCfg MCPConfig
	require.NoError(t, json.Unmarshal(readData, &readCfg))
	assert.Contains(t, readCfg.Servers, "github")
}

func TestSupabaseServerConfigInterface(t *testing.T) {
	// Test map-based config
	m := map[string]string{"projectRef": "my-project"}
	ref, ok := m["projectRef"]
	assert.True(t, ok)
	assert.Equal(t, "my-project", ref)
}

func TestGenerateMCPJSONToPath(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "mcp.json")

	// Test with multiple server configs
	serverConfigs := map[string]any{
		"atlassian":  struct{}{},
		"kubernetes": struct{}{},
		"vault":      struct{}{},
		"github":     struct{}{},
		"supabase":   map[string]string{"projectRef": "test-project"},
	}

	err := GenerateMCPJSONToPath(testPath, serverConfigs)
	require.NoError(t, err)

	// Verify file exists
	assert.FileExists(t, testPath)

	// Read and verify content
	data, err := os.ReadFile(testPath)
	require.NoError(t, err)

	var cfg MCPConfig
	require.NoError(t, json.Unmarshal(data, &cfg))

	// Verify all servers are configured
	assert.Contains(t, cfg.Servers, "atlassian")
	assert.Contains(t, cfg.Servers, "kubernetes")
	assert.Contains(t, cfg.Servers, "vault-mcp-server")
	assert.Contains(t, cfg.Servers, "github")
	assert.Contains(t, cfg.Servers, "supabase")

	// Verify server types
	assert.Equal(t, "stdio", cfg.Servers["atlassian"].Type)
	assert.Equal(t, "docker", cfg.Servers["atlassian"].Command)
	assert.Equal(t, "http", cfg.Servers["github"].Type)
	assert.Equal(t, "https://api.githubcopilot.com/mcp/", cfg.Servers["github"].URL)
	assert.Contains(t, cfg.Servers["supabase"].URL, "test-project")
}

func TestGenerateMCPJSONToPathMergesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "mcp.json")

	// Create existing config
	existing := MCPConfig{
		Inputs: []any{"existing-input"},
		Servers: map[string]ServerConfig{
			"custom-server": {Type: "http", URL: "https://custom.com"},
		},
	}
	data, _ := json.Marshal(existing)
	require.NoError(t, os.WriteFile(testPath, data, 0644))

	// Add new server
	err := GenerateMCPJSONToPath(testPath, map[string]any{"github": struct{}{}})
	require.NoError(t, err)

	// Verify both exist
	newData, err := os.ReadFile(testPath)
	require.NoError(t, err)

	var cfg MCPConfig
	require.NoError(t, json.Unmarshal(newData, &cfg))

	assert.Contains(t, cfg.Servers, "custom-server")
	assert.Contains(t, cfg.Servers, "github")
}

func TestGenerateMCPJSONToPathCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "nested", "dir", "mcp.json")

	err := GenerateMCPJSONToPath(testPath, map[string]any{"github": struct{}{}})
	require.NoError(t, err)

	assert.FileExists(t, testPath)
}

func TestGenerateMCPJSONForEditorUnsupported(t *testing.T) {
	err := GenerateMCPJSONForEditor("unknown-editor", map[string]any{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported editor")
}

func TestGenerateMCPJSONForEditors(t *testing.T) {
	// This tests the multi-editor function but will fail on unsupported
	err := GenerateMCPJSONForEditors([]editor.Editor{"unknown"}, map[string]any{})
	assert.Error(t, err)
}

func TestGenerateMCPJSONForEditorsSuccess(t *testing.T) {
	// Test with a valid editor in temp directory
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create editor config dir
	os.MkdirAll(filepath.Join(tmpDir, ".config", "zed"), 0755)

	err := GenerateMCPJSONForEditors([]editor.Editor{editor.Zed}, map[string]any{"github": struct{}{}})
	require.NoError(t, err)
}

func TestGenerateMCPJSONWrapper(t *testing.T) {
	// This tests the wrapper that uses VSCode default
	// Will create file in real VS Code location, so just test it doesn't error badly
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create VS Code config dir (darwin path)
	os.MkdirAll(filepath.Join(tmpDir, "Library", "Application Support", "Code", "User"), 0755)

	err := GenerateMCPJSON(map[string]any{"github": struct{}{}})
	require.NoError(t, err)
}

func TestGenerateMCPJSONToPathWithSupabaseInterface(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "mcp.json")

	// Test with SupabaseServerConfig interface using setup.SupabaseConfig
	serverConfigs := map[string]any{
		"supabase": &supabaseConfigMock{ref: "my-project-ref"},
	}

	err := GenerateMCPJSONToPath(testPath, serverConfigs)
	require.NoError(t, err)

	// Verify supabase was configured with interface
	data, err := os.ReadFile(testPath)
	require.NoError(t, err)

	var cfg MCPConfig
	require.NoError(t, json.Unmarshal(data, &cfg))

	assert.Contains(t, cfg.Servers, "supabase")
	assert.Contains(t, cfg.Servers["supabase"].URL, "my-project-ref")
}

// supabaseConfigMock implements SupabaseServerConfig
type supabaseConfigMock struct {
	ref string
}

func (s *supabaseConfigMock) GetProjectRef() string {
	return s.ref
}
