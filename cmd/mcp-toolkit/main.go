// Package main is the entry point for the MCP Toolkit CLI.
package main

import (
	"fmt"
	"os"

	"github.com/mcp-toolkit/internal/cli"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	cli.SetVersion(version, commit, buildTime)
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
