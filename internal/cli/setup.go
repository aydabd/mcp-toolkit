package cli

import (
	"fmt"

	"github.com/mcp-toolkit/internal/config"
	"github.com/mcp-toolkit/internal/container"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup [server]",
	Short: "Setup a specific MCP server",
	Long:  "Available: atlassian, kubernetes, vault, github, supabase",
	Args:  cobra.ExactArgs(1),
	RunE:  runSetup,
}

func init() {
	setupCmd.Flags().BoolVar(&skipContainers, "skip-containers", false, "skip container image pulling")
}

func runSetup(cmd *cobra.Command, args []string) error {
	server := args[0]
	if err := config.EnsureEnvDirectory(); err != nil {
		return err
	}

	runtime, _ := container.NewRuntime()

	handlers := map[string]struct {
		fn    func() (any, error)
		image string
	}{
		"atlassian":  {setupAtlassian, "mcp/atlassian"},
		"kubernetes": {setupKubernetes, "mcp/kubernetes"},
		"vault":      {setupVault, ""},
		"github":     {setupGitHub, ""},
		"supabase":   {setupSupabase, ""},
	}

	h, ok := handlers[server]
	if !ok {
		return fmt.Errorf("unknown server: %s", server)
	}

	cfg, err := h.fn()
	if err != nil || cfg == nil {
		return err
	}

	if h.image != "" && runtime != nil && !skipContainers {
		fmt.Printf("📦 Pulling %s...\n", h.image)
		_ = runtime.Pull(cmd.Context(), h.image)
	}

	if err := config.GenerateMCPJSON(map[string]any{server: cfg}); err != nil {
		return err
	}

	fmt.Printf("\n✓ %s setup complete. Restart VS Code.\n", server)
	return nil
}
