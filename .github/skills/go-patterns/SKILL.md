---
name: go-patterns
description: >
  Architecture patterns for MCP Toolkit Go code. Use this when writing or reviewing
  Go functions — composable design, error handling, and dependency injection.
---

# Skill: Go Architecture Patterns

## Composable Functions

Prefer small, composable functions with clear responsibilities:

```go
func ConfigureServer(srv *Server, editor Editor) error {
    if err := validate(srv); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    return editor.WriteMCPConfig(srv)
}
```

## Error Handling

Always wrap errors with context using `fmt.Errorf` and `%w`:

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

- Never ignore errors — handle or propagate every one.
- Wrap with enough context to trace the call chain.
- Use `errors.Is()` and `errors.As()` for checking wrapped errors.

## Dependency Injection

Depend on interfaces, not concrete types:

```go
type Runtime interface {
    Available() bool
    Name() string
}

func DetectRuntime(runtimes []Runtime) (Runtime, error) {
    for _, r := range runtimes {
        if r.Available() {
            return r, nil
        }
    }
    return nil, fmt.Errorf("no container runtime found")
}
```

## Design Principles

- **Single responsibility** — each function and package does one thing.
- **Open/closed** — extend via interfaces, not modification.
- **Interface segregation** — small, focused interfaces (1–3 methods).
- **Composition over inheritance** — embed structs; avoid deep hierarchies.
- **No global state** — pass dependencies explicitly.
