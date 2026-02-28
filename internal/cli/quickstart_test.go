package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcp-toolkit/internal/config"
	"github.com/mcp-toolkit/internal/prompt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuickstartCommand(t *testing.T) {
	assert.Equal(t, "quickstart", quickstartCmd.Use)
	assert.NotEmpty(t, quickstartCmd.Short)
}

func TestPrintBanner(t *testing.T) {
	// Verify it doesn't panic
	printBanner()
}

func TestPrintSummaryEmpty(t *testing.T) {
	// Test empty servers
	printSummary(map[string]any{})
}

func TestPrintSummaryWithServers(t *testing.T) {
	servers := map[string]any{
		"atlassian": struct{}{},
		"vault":     struct{}{},
	}
	printSummary(servers)
}

func TestSetupAtlassian(t *testing.T) {
	// Ensure env directory exists
	require.NoError(t, os.MkdirAll(config.EnvDir(), 0700))

	// Mock prompt input
	input := "https://test.atlassian.net\nuser@test.com\nsecret-token\nPROJ,DEV\ny\n"
	prompt.SetReader(strings.NewReader(input))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupAtlassian()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestSetupAtlassianEmpty(t *testing.T) {
	// Empty URL should return nil
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupAtlassian()
	require.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestSetupKubernetes(t *testing.T) {
	cfg, err := setupKubernetes()
	require.NoError(t, err)
	// May or may not return config depending on kubeconfig existence
	_ = cfg
}

func TestSetupKubernetesWithKubeconfig(t *testing.T) {
	// Create temp HOME with kubeconfig
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create kubeconfig
	kubeDir := filepath.Join(tmpDir, ".kube")
	os.MkdirAll(kubeDir, 0755)
	os.WriteFile(filepath.Join(kubeDir, "config"), []byte("test"), 0644)

	cfg, err := setupKubernetes()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestSetupKubernetesNoKubeconfig(t *testing.T) {
	// Create temp HOME without kubeconfig
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := setupKubernetes()
	require.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestSetupVault(t *testing.T) {
	require.NoError(t, os.MkdirAll(config.EnvDir(), 0700))

	input := "https://vault.test.com\nhvs.secret-token\nadmin\n/usr/bin/vault-mcp-server\n"
	prompt.SetReader(strings.NewReader(input))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupVault()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestSetupVaultEmpty(t *testing.T) {
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupVault()
	require.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestSetupGitHub(t *testing.T) {
	cfg, err := setupGitHub()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestSetupSupabase(t *testing.T) {
	input := "my-project-ref\n"
	prompt.SetReader(strings.NewReader(input))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupSupabase()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestSetupSupabaseEmpty(t *testing.T) {
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	cfg, err := setupSupabase()
	require.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestRunQuickstartAllServersNoConfirm(t *testing.T) {
	// Test with --all flag but all prompts empty (skip all)
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create env directory
	os.MkdirAll(filepath.Join(tmpDir, ".mcp-server-envs"), 0700)

	// Empty prompts for all servers
	input := "\n\n\n\n\n"
	prompt.SetReader(strings.NewReader(input))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	allServers = true
	skipContainers = true
	defer func() {
		allServers = false
		skipContainers = false
	}()

	err := runQuickstart(quickstartCmd, []string{})
	require.NoError(t, err)
}
