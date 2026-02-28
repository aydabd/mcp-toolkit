// Package container provides a unified interface for Docker/Podman/Rancher.
package container

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runtime represents a container runtime (Docker, Podman, or Rancher Desktop).
type Runtime struct {
	binary string
}

// supportedRuntimes lists container runtimes in preference order.
var supportedRuntimes = []string{"docker", "podman", "nerdctl"}

// NewRuntime detects and returns an available container runtime.
func NewRuntime() (*Runtime, error) {
	for _, bin := range supportedRuntimes {
		if path, err := exec.LookPath(bin); err == nil {
			if err := checkRuntime(path); err == nil {
				return &Runtime{binary: path}, nil
			}
		}
	}
	return nil, fmt.Errorf("no container runtime found (tried: %s)", strings.Join(supportedRuntimes, ", "))
}

// Binary returns the runtime binary path.
func (r *Runtime) Binary() string {
	return r.binary
}

// Pull pulls a container image.
func (r *Runtime) Pull(ctx context.Context, image string) error {
	if !strings.Contains(image, ":") {
		image += ":latest"
	}
	cmd := exec.CommandContext(ctx, r.binary, "pull", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pull %s: %w\n%s", image, err, output)
	}
	return nil
}

// checkRuntime verifies runtime is working by running version command.
func checkRuntime(binary string) error {
	cmd := exec.Command(binary, "version")
	return cmd.Run()
}
