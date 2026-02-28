package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("mcp-toolkit %s (commit: %s, built: %s)\n", appVersion, appCommit, appBuildTime)
	},
}
