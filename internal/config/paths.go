// Package config handles configuration paths and file generation.
package config

import (
	"os"
	"path/filepath"
)

// EnvDir returns ~/.mcp-server-envs path.
func EnvDir() string {
	return filepath.Join(homeDir(), ".mcp-server-envs")
}

// KubeconfigPath returns ~/.kube/config path.
func KubeconfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

// AtlassianEnvPath returns Atlassian env file path.
func AtlassianEnvPath() string {
	return filepath.Join(EnvDir(), ".mcp-atlassian.env")
}

// VaultEnvPath returns Vault env file path.
func VaultEnvPath() string {
	return filepath.Join(EnvDir(), ".mcp-vault.env")
}

// VaultWrapperPath returns Vault wrapper script path.
func VaultWrapperPath() string {
	return filepath.Join(EnvDir(), "vault-mcp-wrapper.sh")
}

// EnsureEnvDirectory creates env directory with 0700 permissions.
func EnsureEnvDirectory() error {
	dir := EnvDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.Chmod(dir, 0700)
}

// FileExists checks if file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func homeDir() string {
	h, _ := os.UserHomeDir()
	return h
}

// KubernetesEnvPath returns Kubernetes env file path.
func KubernetesEnvPath() string {
	return filepath.Join(EnvDir(), ".mcp-kubernetes.env")
}

// GitHubEnvPath returns GitHub env file path.
func GitHubEnvPath() string {
	return filepath.Join(EnvDir(), ".mcp-github.env")
}

// SupabaseEnvPath returns Supabase env file path.
func SupabaseEnvPath() string {
	return filepath.Join(EnvDir(), ".mcp-supabase.env")
}
