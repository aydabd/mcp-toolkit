//go:build integration

// Package e2e provides end-to-end tests in isolated environments.
package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E runs the CLI in an isolated temp HOME directory.
func TestE2E(t *testing.T) {
	// Build the binary
	binary := buildBinary(t)

	t.Run("version", func(t *testing.T) {
		out := runCmd(t, binary, nil, "version")
		assert.Contains(t, out, "mcp-toolkit")
	})

	t.Run("help", func(t *testing.T) {
		out := runCmd(t, binary, nil, "--help")
		assert.Contains(t, out, "MCP Toolkit")
		assert.Contains(t, out, "Environment Variables")
		assert.Contains(t, out, "MCP_EDITOR")
	})

	t.Run("env_shows_defaults", func(t *testing.T) {
		out := runCmd(t, binary, nil, "env")
		assert.Contains(t, out, "MCP_ATLASSIAN_URL")
		assert.Contains(t, out, "vscode (default)")
	})

	t.Run("env_override", func(t *testing.T) {
		env := map[string]string{
			"MCP_ATLASSIAN_URL": "https://test.atlassian.net",
		}
		out := runCmd(t, binary, env, "env")
		assert.Contains(t, out, "https://test.atlassian.net")
	})

	t.Run("isolated_home", func(t *testing.T) {
		// Create isolated home
		home := t.TempDir()
		env := map[string]string{"HOME": home}

		// Run env command
		out := runCmd(t, binary, env, "env")
		assert.Contains(t, out, "MCP Toolkit")

		// Verify no files leaked to real home
		realHome, _ := os.UserHomeDir()
		testMarker := filepath.Join(realHome, ".mcp-toolkit-test-marker")
		assert.NoFileExists(t, testMarker)
	})

	t.Run("quickstart_help", func(t *testing.T) {
		out := runCmd(t, binary, nil, "quickstart", "--help")
		assert.Contains(t, out, "Interactive setup")
		assert.Contains(t, out, "--all")
		assert.Contains(t, out, "--skip-containers")
	})

	t.Run("setup_help", func(t *testing.T) {
		out := runCmd(t, binary, nil, "setup", "--help")
		assert.Contains(t, out, "atlassian")
		assert.Contains(t, out, "kubernetes")
	})

	t.Run("completion", func(t *testing.T) {
		out := runCmd(t, binary, nil, "completion", "bash")
		assert.Contains(t, out, "bash completion")
	})
}

func buildBinary(t *testing.T) string {
	t.Helper()

	// Find project root
	wd, err := os.Getwd()
	require.NoError(t, err)
	projectRoot := filepath.Join(wd, "..", "..")

	// Build to temp location
	binary := filepath.Join(t.TempDir(), "mcp-toolkit")

	cmd := exec.Command("go", "build", "-o", binary, "./cmd/mcp-toolkit")
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "build failed: %s", string(out))

	return binary
}

func runCmd(t *testing.T, binary string, env map[string]string, args ...string) string {
	t.Helper()

	cmd := exec.Command(binary, args...)

	// Start with minimal environment
	baseEnv := []string{
		"PATH=" + os.Getenv("PATH"),
	}

	// Use temp home by default for isolation
	if env == nil || env["HOME"] == "" {
		baseEnv = append(baseEnv, "HOME="+t.TempDir())
	}

	// Add custom env vars
	for k, v := range env {
		baseEnv = append(baseEnv, k+"="+v)
	}
	cmd.Env = baseEnv

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stderr: %s", stderr.String())
	}
	require.NoError(t, err, "command failed: %s %s", binary, strings.Join(args, " "))

	return stdout.String()
}
