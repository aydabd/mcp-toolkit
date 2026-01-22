package cli

import (
	"fmt"

	"github.com/mcp-toolkit/internal/envvar"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment variables",
	Long:  "Display all environment variables used by mcp-toolkit with their current values.",
	Run: func(cmd *cobra.Command, args []string) {
		showEnvVars()
	},
}

func showEnvVars() {
	fmt.Println("MCP Toolkit Configuration")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("Values can be set via environment variables (MCP_*) or config file (~/.mcp-toolkit.yaml)")
	fmt.Println()

	for category, vars := range envvar.ByCategory() {
		fmt.Printf("%s:\n", category)
		for _, v := range vars {
			val := v.Get()
			status := "not set"

			if val != "" {
				if v.Secret {
					status = "***"
				} else if val == v.Default {
					status = fmt.Sprintf("%s (default)", val)
				} else {
					status = val
				}
			}

			required := ""
			if v.Required {
				required = "*"
			}

			fmt.Printf("  %-35s %s\n", v.Name+required, status)
		}
		fmt.Println()
	}

	fmt.Println("* = required")
}
