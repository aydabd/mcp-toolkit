---
name: common-tasks
description: >
  Common development tasks for MCP Toolkit. Use this when building, adding new
  servers, editors, CLI flags, or running the change checklist.
---

# Skill: Common Tasks

## Building and Testing

```bash
# Setup environment
micromamba create -f environment.yml
micromamba activate mcp-toolkit

# Development commands
make lint           # Format, vet, and lint
make test           # Run unit tests
make test-cover     # Run tests with coverage
make coverage-check # Verify 90% coverage threshold
make ci             # Full CI pipeline (lint + test + coverage)
make build          # Build for current platform
make build-all      # Build for all platforms
```

## Adding a New MCP Server

1. Add server type to `internal/setup/types.go`.
2. Add writer function to `internal/setup/writers.go`.
3. Add setup subcommand to `internal/cli/setup.go`.
4. Add to quickstart flow in `internal/cli/quickstart.go`.
5. Add environment variables to `internal/envvar/envvar.go`.
6. Write tests for all new code.

## Adding a New Editor

1. Add editor type to `internal/editor/editor.go`.
2. Add OS-specific paths in `editor_darwin.go`, `editor_linux.go`, `editor_windows.go`.
3. Write tests in `editor_test.go`.

## Adding a CLI Flag

```go
cmd.Flags().String("flag-name", "default", "Clear description")
```

## Before Any Change — Checklist

1. Is this the simplest solution?
2. Will this be easy to maintain in two years?
3. Can someone new understand this in five minutes?
4. Does this follow Go idioms?
5. Are all new functions testable and tested?

If any answer is "no", reconsider the approach.
