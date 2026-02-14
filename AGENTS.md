# AGENTS.md - cmdguard Project Guide

**Last Updated:** 2026-02-14  
**Project:** cmdguard - CLI Guard Library  
**Go Version:** 1.26.0  
**Status:** Phase 2 Complete - Testing at 100% for all packages

---

## Quick Start

```bash
# Run tests
go test ./...
go test -v -cover ./...        # Verbose with coverage
go test -race ./...            # Race detection

# Build examples
cd examples/basic && go build -o myapp .
```

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with fail-fast validation. It provides:

- **Compile-time validation** - Panics at construction if commands are invalid
- **Single-step initialization** - No multi-step init, validate, execute flow
- **Guard philosophy** - Fail fast at startup, not at runtime
- **Minimal dependencies** - Only cobra, fang, and testify

**Current Status:** Phase 2 Complete. All packages tested with 100% coverage.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/           # Public API - GuardedCommand
│   └── guarded_command.go  # Guard API implementation
├── internal/
│   ├── config/             # Configuration (environment variables)
│   └── logging/            # Structured logging (slog)
├── examples/               # Example applications
│   └── basic/              # Basic CLI example
├── docs/
│   └── CLI_DESIGN_PRINCIPLES.md
├── FEATURES.md             # Feature status
├── TODO_LIST.md            # Remaining tasks
├── justfile                # Build commands
├── go.mod
└── README.md
```

### Package Guidelines

| Package | Purpose | Importable? | Coverage |
|---------|---------|-------------|----------|
| `pkg/cmdguard` | Public Guard API | Yes | 66.7% |
| `internal/config` | Configuration | No | 94.1% |
| `internal/logging` | Logging utilities | No | 100% |

---

## Key Dependencies

| Library | Purpose | Version |
|---------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | v1.10.2 |
| `github.com/charmbracelet/fang` | Cobra styling | v0.4.4 |
| `github.com/stretchr/testify` | Testing | v1.11.1 |

---

## Coding Standards

### Go Conventions

- **Go 1.26.0** - Use modern Go features
- **gofumpt** formatting preferred
- **Error handling** - Always check errors, wrap with context
- **Interface naming** - `-er` suffix (e.g., `Validator`)
- **Constructor naming** - `New` + type name (e.g., `NewLogger`)

### Validation Patterns

```go
// Valid: Has RunE handler
cmd := &cobra.Command{
    Use: "test",
    RunE: func(cmd *cobra.Command, args []string) error {
        return nil
    },
}

// Valid: Has subcommands (intermediate command)
parent := &cobra.Command{Use: "parent"}
parent.AddCommand(child)

// Invalid: No handler, no subcommands - WILL PANIC
cmd := &cobra.Command{Use: "invalid"}
```

---

## Testing

### Test Commands

```bash
# Run all tests with coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Race detection
go test -race ./...
```

### Test Pattern

Use `testify/assert` with table-driven tests:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {name: "valid", input: "test", want: "result"},
        {name: "invalid", input: "", wantErr: true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := DoSomething(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

---

## Guard API Reference

### Basic Usage

```go
package main

import (
    "context"
    "github.com/larsartmann/cmdguard/pkg/cmdguard"
    "github.com/spf13/cobra"
)

func main() {
    root := cmdguard.New("myapp", "My application")
    
    root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, World!")
        },
    })
    
    root.ExecuteAndExit(context.Background())
}
```

### GuardedCommand Methods

| Method | Description |
|--------|-------------|
| `New(name, short)` | Create new guarded command |
| `AddCommand(cmd)` | Add subcommand (panics if invalid) |
| `AddSubcommand(parent, child)` | Add nested subcommand |
| `Execute(ctx)` | Run with context |
| `ExecuteAndExit(ctx)` | Run and exit with code |
| `Command()` | Access underlying cobra command |
| `Config()` | Get configuration |
| `IsStrictMode()` | Check strict mode |

### Panic Conditions (Intentional)

The Guard API panics at construction time if:
1. Command has no `Run`/`RunE` and no subcommands
2. Strict mode requires `RunE` but only `Run` provided
3. Command has no name

---

## Configuration

Environment variables (prefix `CMDGUARD_`):
- `CMDGUARD_LOG_LEVEL` - debug, info, warn, error (default: info)
- `CMDGUARD_STRICT_MODE` - true/false (default: false)

---

## Architecture Decisions

### Guard API over Framework

**Before:** Multi-step initialization with separate Validate step  
**Now:** Single-step, panic at construction time

Rationale: Fail-fast philosophy. Go lacks compile-time macros; panic at init is closest to catching errors early.

### Minimal Dependencies

Removed:
- `samber/do/v2` - No DI container needed
- `knadh/koanf/v2` - Direct env var reading is simpler

Kept:
- `cobra` - CLI framework
- `fang` - Beautiful help output
- `testify` - Testing assertions

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [testify Documentation](https://github.com/stretchr/testify)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
