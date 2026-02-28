# MCP Toolkit

CLI tool to configure MCP (Model Context Protocol) servers for VS Code, Cursor, Windsurf, and Zed editors.

## Quick Start

```bash
make build && ./bin/mcp-toolkit quickstart
```

## Requirements

- Go 1.23+
- Container runtime (optional): Docker, Podman, or Rancher Desktop
- micromamba (for development)

## Installation

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Install to GOPATH/bin
```

## Usage

```bash
mcp-toolkit quickstart           # Interactive setup for all servers
mcp-toolkit quickstart --all     # Configure all servers without prompts
mcp-toolkit setup atlassian      # Setup specific server
mcp-toolkit env                  # Show environment variables
mcp-toolkit version              # Show version
```

## Supported Servers

| Server     | Container      | Description                   |
| ---------- | -------------- | ----------------------------- |
| atlassian  | mcp/atlassian  | Jira & Confluence integration |
| kubernetes | mcp/kubernetes | K8s cluster management        |
| vault      | local binary   | HashiCorp Vault secrets       |
| github     | -              | GitHub Copilot (built-in)     |
| supabase   | -              | Supabase database (hosted)    |

## Supported Editors

| Editor   | Config Location                                                |
| -------- | -------------------------------------------------------------- |
| VS Code  | `~/Library/Application Support/Code/User/mcp.json` (macOS)     |
| Cursor   | `~/Library/Application Support/Cursor/User/mcp.json` (macOS)   |
| Windsurf | `~/Library/Application Support/Windsurf/User/mcp.json` (macOS) |
| Zed      | `~/.config/zed/settings.json`                                  |

## Development

```bash
# Setup environment
micromamba create -f environment.yml
micromamba activate mcp-toolkit

# Development commands
make lint           # Format and vet code
make test           # Run unit tests
make test-cover     # Run tests with coverage
make coverage-check # Verify 90% coverage threshold
make ci             # Full CI pipeline
make pre-commit     # Run pre-commit hooks
```

## Project Structure

```text
cmd/mcp-toolkit/     # CLI entry point
internal/
  cli/               # Cobra commands (quickstart, setup, env, version)
  config/            # Paths and mcp.json generation
  container/         # Docker/Podman/Rancher runtime detection
  editor/            # Multi-editor support with OS-specific paths
  envvar/            # Environment variable management
  prompt/            # Interactive terminal input
  setup/             # Server configuration writers
make/                # Modular Makefile includes
test/e2e/            # End-to-end tests
```

## Configuration

| Type          | Location                                      |
| ------------- | --------------------------------------------- |
| Credentials   | `~/.mcp-server-envs/*.env` (0700 permissions) |
| Editor config | OS and editor-specific `mcp.json`             |

## Environment Variables

All configuration can be set via environment variables with `MCP_` prefix:

```bash
MCP_ATLASSIAN_URL=https://company.atlassian.net
MCP_ATLASSIAN_EMAIL=user@example.com
MCP_ATLASSIAN_API_TOKEN=secret
MCP_VAULT_ADDRESS=https://vault.example.com
```

Run `mcp-toolkit env` to see all available variables.

## License

MIT
