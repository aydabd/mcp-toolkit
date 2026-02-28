---
applyTo: "**"
---

# MCP Toolkit — AI Agent Instructions

> **Single source of truth** for all AI coding agents (GitHub Copilot, Claude Code,
> OpenAI Codex, Cursor, and others).
> Edit **only this file** to update instructions across every agent.

## Project Overview

MCP Toolkit is a Go CLI tool that configures MCP (Model Context Protocol) servers
for VS Code, Cursor, Windsurf, and Zed editors. It detects container runtimes,
manages editor-specific `mcp.json` configs, and handles credential files with
proper permissions.

```text
cmd/mcp-toolkit/     # CLI entry point (main.go)
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

## Golden Rules

1. **Simplicity** — simplest working solution wins.
2. **Go Idioms** — follow standard Go patterns (`gofmt`, `golangci-lint`).
3. **SOLID** — single responsibility, dependency injection, interface-based design.
4. **Long-term** — code must be maintainable for years.
5. **Testability** — every package must be independently testable.

## Package Layout Rules

- `internal/` — all implementation is private. No package is exported outside the module.
- Each `internal/<pkg>/` has a single, focused responsibility.
- `cmd/mcp-toolkit/main.go` — only wires dependencies and calls `cli.Execute()`.
- OS-specific code uses build-tag files: `editor_darwin.go`, `editor_linux.go`, etc.

## Code Style

- Use `gofmt` and `golangci-lint` standards.
- Keep functions small and focused (max 50 lines).
- Use meaningful variable names; avoid abbreviations.
- Add comments only when code is not self-explanatory (explain _why_, not _what_).
- Handle errors explicitly — never ignore them.
- Write tests for new functionality.
- **Formatting**: 4 spaces for Go code (tabs via `gofmt`), 2 spaces for YAML/JSON.
- No trailing whitespace.

## Skills

Specialised agent skills live in `.github/skills/`.
Each skill directory contains a `SKILL.md` (≤ 100 lines) covering one focused topic.

| Skill directory | Purpose                                                     |
| --------------- | ----------------------------------------------------------- |
| `go-patterns`   | Architecture patterns, error handling, dependency injection |
| `testing-go`    | Table-driven tests, mocking, coverage, E2E testing          |
| `common-tasks`  | Build commands, adding servers/editors/flags, checklists    |

## Documentation

Update docs only when:

- CLI interface changes.
- Setup process changes.
- New configuration options are added.
- New MCP servers or editors are supported.

Do **not** update docs for internal refactoring, bug fixes, or performance improvements.

## What to Avoid

- Features added "just in case"
- Complex abstractions or inheritance hierarchies
- Reflection unless absolutely necessary
- Magic numbers (use named constants)
- Global state
- New dependencies without strong justification
- Premature optimization
- Hardcoded secrets or paths
- Silently ignoring errors
- Placeholder or TODO code in production
