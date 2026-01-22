// Package cli provides Cobra CLI commands for MCP Toolkit.
package cli

import (
	"log/slog"
	"os"

	"github.com/mcp-toolkit/internal/envvar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	verbose      bool
	appVersion   = "dev"
	appCommit    = "unknown"
	appBuildTime = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-toolkit",
	Short: "Configure MCP servers for VS Code",
	Long: `MCP Toolkit automates setup of Model Context Protocol servers for VS Code.

` + envvar.FormatHelp(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

// Execute runs the root command.
func Execute() error { return rootCmd.Execute() } //nolint:wrapcheck // cobra entry point

// SetVersion sets build information.
func SetVersion(version, commit, buildTime string) {
	appVersion, appCommit, appBuildTime = version, commit, buildTime
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(quickstartCmd, setupCmd, versionCmd)
}

func initConfig() {
	// Bind all environment variables first
	envvar.BindAll()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, _ := os.UserHomeDir()
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mcp-toolkit")
	}
	_ = viper.ReadInConfig()
}

func initLogger() {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
}
