# MCP Toolkit

CLI tool to configure MCP servers for VS Code. Supports Docker, Podman, and Rancher.

## Quick Start

```bash
make build && ./bin/mcp-toolkit quickstart
```

## Requirements

- Go 1.25+
- Container runtime: Docker, Podman, or Rancher Desktop

## Installation

```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make install        # Install to GOPATH/bin
```

## Usage

```bash
mcp-toolkit quickstart           # Interactive setup
mcp-toolkit setup atlassian      # Setup specific server
mcp-toolkit version              # Show version
```

## Supported Servers

| Server | Container | Description |
|--------|-----------|-------------|
| atlassian | mcp/atlassian | Jira & Confluence |
| kubernetes | mcp/kubernetes | K8s cluster management |
| vault | - | HashiCorp Vault |
| github | - | GitHub Copilot (hosted) |
| supabase | - | Supabase (hosted) |

## Development

```bash
micromamba create -f environment.yml
micromamba activate mcp-toolkit
make pre-commit
make test-cover
make ci
```

## Project Structure

```
cmd/mcp-toolkit/     # Entry point
internal/
  cli/               # Cobra commands
  config/            # Path and mcp.json handling
  container/         # Docker/Podman/Rancher runtime
  prompt/            # Terminal input
  setup/             # Server env writers
```

## Configuration

- Credentials: `~/.mcp-server-envs/` (secure, not committed)
- VS Code config: platform-specific `mcp.json` location

## License

MIT
