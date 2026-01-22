package container

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedRuntimes(t *testing.T) {
	assert.Contains(t, supportedRuntimes, "docker")
	assert.Contains(t, supportedRuntimes, "podman")
	assert.Contains(t, supportedRuntimes, "nerdctl")
	assert.Len(t, supportedRuntimes, 3)
}

func TestRuntimeBinary(t *testing.T) {
	r := &Runtime{binary: "/usr/bin/docker"}
	assert.Equal(t, "/usr/bin/docker", r.Binary())
}

func TestNewRuntime(t *testing.T) {
	// This test may succeed or fail depending on environment
	r, err := NewRuntime()
	if err != nil {
		assert.Contains(t, err.Error(), "no container runtime found")
	} else {
		assert.NotEmpty(t, r.Binary())
	}
}

func TestPullImageTagging(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"nginx", "nginx:latest"},
		{"nginx:1.0", "nginx:1.0"},
		{"registry.io/image", "registry.io/image:latest"},
		{"registry.io/image:tag", "registry.io/image:tag"},
	}

	for _, tt := range tests {
		image := tt.input
		if !contains(image, ":") {
			image += ":latest"
		}
		assert.Equal(t, tt.expected, image)
	}
}

func TestPullWithoutRuntime(t *testing.T) {
	r := &Runtime{binary: "/nonexistent/binary"}
	err := r.Pull(context.Background(), "test:latest")
	assert.Error(t, err)
}

func TestPullWithoutTag(t *testing.T) {
	r := &Runtime{binary: "/nonexistent/binary"}
	// Test that :latest is appended
	err := r.Pull(context.Background(), "nginx")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nginx:latest")
}

func TestCheckRuntimeFails(t *testing.T) {
	err := checkRuntime("/nonexistent/binary")
	assert.Error(t, err)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
