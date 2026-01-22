package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mcp-toolkit/internal/prompt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupCommand(t *testing.T) {
	assert.Contains(t, setupCmd.Use, "setup")
	assert.NotEmpty(t, setupCmd.Short)
}

func TestRunSetupUnknownServer(t *testing.T) {
	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"unknown-server"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown server")
}

func TestRunSetupGitHub(t *testing.T) {
	// GitHub doesn't need prompts
	skipContainers = true
	defer func() { skipContainers = false }()

	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"github"})
	require.NoError(t, err)
}

func TestRunSetupAtlassianSkipped(t *testing.T) {
	// Empty input skips configuration
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	skipContainers = true
	defer func() { skipContainers = false }()

	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"atlassian"})
	require.NoError(t, err)
}

func TestRunSetupVaultSkipped(t *testing.T) {
	// Empty input skips configuration
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	skipContainers = true
	defer func() { skipContainers = false }()

	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"vault"})
	require.NoError(t, err)
}

func TestRunSetupSupabaseSkipped(t *testing.T) {
	prompt.SetReader(strings.NewReader("\n"))
	prompt.SetOutput(&bytes.Buffer{})
	defer prompt.Reset()

	skipContainers = true
	defer func() { skipContainers = false }()

	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"supabase"})
	require.NoError(t, err)
}

func TestRunSetupKubernetesWithConfig(t *testing.T) {
	// Create temp HOME with kubeconfig
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create kubeconfig and env directory
	kubeDir := filepath.Join(tmpDir, ".kube")
	os.MkdirAll(kubeDir, 0755)
	os.WriteFile(filepath.Join(kubeDir, "config"), []byte("test"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, ".mcp-server-envs"), 0700)

	// Create VS Code config dir for mcp.json
	os.MkdirAll(filepath.Join(tmpDir, "Library", "Application Support", "Code", "User"), 0755)

	skipContainers = true
	defer func() { skipContainers = false }()

	cmd := &cobra.Command{}
	err := runSetup(cmd, []string{"kubernetes"})
	require.NoError(t, err)
}
