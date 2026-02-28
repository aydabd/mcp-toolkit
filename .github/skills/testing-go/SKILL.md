---
name: testing-go
description: >
  Go testing patterns for MCP Toolkit. Use this when writing or reviewing tests —
  table-driven tests, mocking, coverage thresholds, and E2E testing.
---

# Skill: Go Testing Patterns

## Table-Driven Tests

Always prefer table-driven tests with `testify`:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid input", "input", "output", false},
        {"empty input returns error", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

## Mocking

Mock at the **interface** boundary, not the implementation:

```go
type mockRuntime struct{ available bool }

func (m *mockRuntime) Available() bool { return m.available }
func (m *mockRuntime) Name() string    { return "mock" }
```

Keep mocks in the same `_test.go` file that uses them.

## Test Naming

Use `TestUnit_Scenario_ExpectedOutcome`:

```go
func TestDetectRuntime_NoRuntimes_ReturnsError(t *testing.T) { ... }
func TestDetectRuntime_DockerAvailable_ReturnsDocker(t *testing.T) { ... }
```

## Coverage

- Minimum 90% coverage (enforced by `make coverage-check`).
- Mock external dependencies (filesystem, container runtime, editor paths).
- Test error cases, not just happy paths.

## E2E Tests

E2E tests live in `test/e2e/` and exercise the compiled binary:

```go
func TestCLI_Version_PrintsVersion(t *testing.T) {
    out, err := exec.Command(binary, "version").CombinedOutput()
    assert.NoError(t, err)
    assert.Contains(t, string(out), "mcp-toolkit")
}
```

## Corner Cases Checklist

Before marking a feature tested, verify:

- Nil/zero/empty input?
- Boundary values (first, last, at-limit)?
- Dependency fails mid-operation?
- Missing required configuration?
- File permission errors?
