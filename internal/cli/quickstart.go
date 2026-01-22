package cli

import (
	"fmt"
	"log/slog"

	"github.com/mcp-toolkit/internal/config"
	"github.com/mcp-toolkit/internal/container"
	"github.com/mcp-toolkit/internal/editor"
	"github.com/mcp-toolkit/internal/prompt"
	"github.com/mcp-toolkit/internal/setup"
	"github.com/spf13/cobra"
)

var (
	skipContainers bool
	allServers     bool
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Interactive setup for all MCP servers",
	RunE:  runQuickstart,
}

func init() {
	quickstartCmd.Flags().BoolVar(&skipContainers, "skip-containers", false, "skip container image pulling")
	quickstartCmd.Flags().BoolVar(&allServers, "all", false, "configure all servers")
}

func runQuickstart(cmd *cobra.Command, args []string) error { //nolint:gocyclo // orchestration function
	printBanner()

	runtime, err := container.NewRuntime()
	if err != nil {
		slog.Warn("container runtime not available", "error", err)
		fmt.Println("⚠️  No container runtime found (Docker/Podman/Rancher)")
	}

	if err := config.EnsureEnvDirectory(); err != nil {
		return fmt.Errorf("failed to create env directory: %w", err)
	}
	fmt.Printf("✓ Environment directory: %s\n\n", config.EnvDir())

	servers := []serverSetup{
		{"atlassian", setupAtlassian, "mcp/atlassian"},
		{"kubernetes", setupKubernetes, "mcp/kubernetes"},
		{"vault", setupVault, ""},
		{"github", setupGitHub, ""},
		{"supabase", setupSupabase, ""},
	}

	configured := make(map[string]any)
	for _, s := range servers {
		if !allServers && !prompt.Confirm(fmt.Sprintf("Configure %s?", s.name)) {
			continue
		}
		fmt.Printf("\n━━━ %s ━━━\n\n", s.name)
		cfg, err := s.fn()
		if err != nil {
			slog.Error("setup failed", "server", s.name, "error", err)
			continue
		}
		if cfg == nil {
			continue
		}
		configured[s.name] = cfg

		if s.image != "" && runtime != nil && !skipContainers {
			fmt.Printf("📦 Pulling %s...\n", s.image)
			if err := runtime.Pull(cmd.Context(), s.image); err != nil {
				slog.Warn("pull failed", "image", s.image, "error", err)
			} else {
				fmt.Printf("✓ %s ready\n", s.image)
			}
		}
	}

	if len(configured) > 0 {
		if err := config.GenerateMCPJSON(configured); err != nil {
			return err
		}
		fmt.Printf("\n✓ mcp.json: %s\n", editor.VSCode.ConfigPath())
	}

	printSummary(configured)
	return nil
}

type serverSetup struct {
	name  string
	fn    func() (any, error)
	image string
}

func setupAtlassian() (any, error) {
	fmt.Println("Get API token: https://id.atlassian.com/manage-profile/security/api-tokens")
	fmt.Println()
	url := prompt.String("Atlassian URL (e.g., https://company.atlassian.net)")
	if url == "" {
		return nil, nil
	}
	cfg := &setup.AtlassianConfig{
		URL:          url,
		Email:        prompt.String("Email"),
		APIToken:     prompt.Secret("API Token"),
		JiraProjects: prompt.String("Jira projects filter (optional)"),
		ReadOnly:     !prompt.ConfirmDefault("Enable write operations?", true),
	}
	if err := setup.WriteAtlassianEnv(cfg); err != nil {
		return nil, err
	}
	fmt.Println("✓ Atlassian configured")
	return cfg, nil
}

func setupKubernetes() (any, error) {
	fmt.Println("Uses ~/.kube/config automatically")
	fmt.Println()
	if !config.FileExists(config.KubeconfigPath()) {
		fmt.Println("⚠️  No kubeconfig found")
		return nil, nil
	}
	fmt.Println("✓ Kubeconfig found")
	return &setup.KubernetesConfig{}, nil
}

func setupVault() (any, error) {
	addr := prompt.String("Vault address")
	if addr == "" {
		return nil, nil
	}
	cfg := &setup.VaultConfig{
		Address:    addr,
		Token:      prompt.Secret("Vault token"),
		Namespace:  prompt.String("Namespace (optional)"),
		BinaryPath: prompt.Default("Binary path", "/usr/local/bin/vault-mcp-server"),
	}
	if err := setup.WriteVaultEnv(cfg); err != nil {
		return nil, err
	}
	if err := setup.WriteVaultWrapper(cfg); err != nil {
		return nil, err
	}
	fmt.Println("✓ Vault configured")
	return cfg, nil
}

func setupGitHub() (any, error) {
	fmt.Println("Uses GitHub Copilot authentication - no config needed")
	fmt.Println()
	fmt.Println("✓ GitHub ready")
	return &setup.GitHubConfig{}, nil
}

func setupSupabase() (any, error) {
	fmt.Println("Get project ref from: https://supabase.com/dashboard")
	fmt.Println()
	ref := prompt.String("Project reference")
	if ref == "" {
		return nil, nil
	}
	fmt.Println("✓ Supabase configured")
	return &setup.SupabaseConfig{ProjectRef: ref}, nil
}

func printBanner() {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════╗")
	fmt.Println("║           MCP Toolkit Quickstart                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")
	fmt.Println()
}

func printSummary(servers map[string]any) {
	fmt.Println("\n╔═══════════════════════════════════════════════════╗")
	fmt.Println("║              Setup Complete!                      ║")
	fmt.Println("╚═══════════════════════════════════════════════════╝")
	if len(servers) == 0 {
		fmt.Println("\nNo servers configured.")
		return
	}
	fmt.Println("\nConfigured:")
	for name := range servers {
		fmt.Printf("  ✓ %s\n", name)
	}
	fmt.Println()
	fmt.Println("Next: Restart VS Code")
	fmt.Println()
}
