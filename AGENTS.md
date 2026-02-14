# AGENTS.md - cmdguard Project Guide

**Last Updated:** 2026-02-14 09:35 UTC  
**Project:** cmdguard - CLI Guard Library  
**Go Version:** 1.26.0  
**Status:** Phase 1 Complete (1% → 51% Foundation)

---

## Quick Start

```bash
# Run tests
go test ./...
go test -v ./internal/...          # Verbose output
go test -cover ./...               # With coverage
go test -race ./...                # Race detection

# Build a binary (using examples/)
cd examples/basic && go build -o myapp .
```

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with compile-time enforcement. It provides:

- **Compile-time validation** - Panics at construction if commands are invalid
- **Single-step initialization** - No multi-step init, validate, execute flow
- **Guard philosophy** - Fail fast at startup, not at runtime
- **Dependency injection** - Built on `samber/do/v2`
- **Configuration management** - Integrated with `knadh/koanf`

**Current Status:** Phase 1 Complete (Foundation). Guard API implemented, cmd/ folder removed.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/           # Public API - GuardedCommand (main entry point)
│   └── guarded_command.go  # Guard API implementation
├── internal/
│   ├── config/             # Configuration management (Koanf)
│   └── logging/            # Structured logging (slog)
├── docs/
│   ├── planning/           # Execution plans
│   ├── status/             # Status reports
│   └── CLI_DESIGN_PRINCIPLES.md
├── FEATURES.md             # Feature status documentation
├── AGENTS.md               # This file
├── go.mod
└── README.md
```

### Package Guidelines

| Package | Purpose | Importable? | Status |
|---------|---------|-------------|--------|
| `pkg/cmdguard` | Public Guard API | Yes | ✅ Active |
| `internal/config` | Configuration | No | ✅ Active |
| `internal/logging` | Logging utilities | No | ✅ Active |
├── internal/
│   ├── commands/           # Cobra command registry and setup
│   ├── config/             # Koanf-based configuration management
│   ├── di/                 # Dependency injection module setup
│   ├── logging/            # slog-based logging utilities
│   └── validation/         # Command and flag validation logic
├── docs/
│   ├── planning/           # Transformation planning documents
│   ├── status/             # Implementation status updates
│   └── CLI_DESIGN_PRINCIPLES.md  # UX guidelines
├── go.mod
└── README.md
```



---

## Key Dependencies

| Library | Purpose | Version |
|---------|---------|---------|
| `github.com/spf13/cobra` | CLI framework | v1.10.2 |
| `github.com/samber/do/v2` | Dependency injection | v2.0.0 |
| `github.com/knadh/koanf/v2` | Configuration management | v2.3.2 |
| `github.com/charmbracelet/fang` | Cobra styling | v0.4.4 |
| `github.com/stretchr/testify` | Testing | v1.11.1 |

---

## Coding Standards

### Go Conventions

- **Go 1.26.0** - Use modern Go features
- **gofumpt** formatting preferred (stricter than gofmt)
- **Error handling** - Always check errors, wrap with context using `fmt.Errorf("...: %w", err)`
- **Interface naming** - `-er` suffix (e.g., `Validator`, `Healthchecker`)
- **Constructor naming** - `New` + type name (e.g., `NewValidator`)

### Validation Patterns

Commands must have handlers:
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

// Invalid: No handler, no subcommands
cmd := &cobra.Command{Use: "invalid"}  // Will fail validation
```

---

## Testing

### Test Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/validation/...

# Verbose output
go test -v ./internal/config/...

# Race detection
go test -race ./...
```

### Test Patterns

Use `testify/assert` and `testify/require`:

```go
func TestNewValidator(t *testing.T) {
    validator, _, _ := setupTestValidator(t)
    assert.NotNil(t, validator)           // Continue on failure
    require.NoError(t, err)               // Stop on failure
}
```

Table-driven tests preferred:
```go
func TestValidator_ValidateCommands(t *testing.T) {
    tests := []struct {
        name    string
        setup   func()
        wantErr bool
        errMsg  string
    }{
        {name: "valid command", ...},
        {name: "invalid command", ...},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            err := validator.ValidateCommands()
            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            }
        })
    }
}
```

### Test Setup Pattern

```go
func setupTestValidator(t *testing.T) (*Validator, *Registry, do.Injector) {
    injector := do.New()
    do.Provide(injector, func(i do.Injector) (*config.Config, error) {
        return &config.Config{}, nil
    })
    do.Provide(injector, NewRegistry)
    do.Provide(injector, NewValidator)
    // ...
}
```

---

## Configuration

Configuration loading order (highest priority last):
1. Default values
2. Config file (`config.yaml`)
3. Environment variables (`CMDGUARD_*` prefix)
4. Command-line flags

```go
// Access config via DI
func NewService(i do.Injector) (*Service, error) {
    cfg, err := do.Invoke[*config.Config](i)
    // Use cfg.StrictMode, cfg.LogLevel
}
```

**Config struct fields:**
- `StrictMode bool` - Enable strict validation
- `LogLevel string` - One of: debug, info, warn, error
- `ConfigFile string` - Path to config file

---

## CLI Design Principles

See `docs/CLI_DESIGN_PRINCIPLES.md` for detailed UX guidelines.

**Key rules:**
- Boolean flags use `BoolP()` - no string values
- Short flags for common options (`-s` for `--strict`)
- Validate enum values in `PreRunE`
- Error messages should suggest fixes
- Every `--help` example must be copy-pasteable

---

## Known Issues & Gotchas

1. **Test coverage** - Only config package has tests (~48% coverage)
2. **No pkg/cmdguard tests** - Guard API needs unit tests
3. **No integration tests** - Need end-to-end tests with actual CLI execution
4. **No examples** - Users need working examples in examples/ directory
5. **gopls warnings** - Stale cached references to deleted files (harmless)

---

## Transformation Progress

Per `docs/planning/2026-02-14_09-27-COMPREHENSIVE_EXECUTION_PLAN.md`:

### Phase 1: Foundation (1% → 51%) ✅ COMPLETE
- ✅ Remove `cmd/` folder - Establish as library
- ✅ Redesign public API - Guard API with single-step initialization
- ✅ Add compile-time validation - Panic on invalid commands
- ✅ Update AGENTS.md - Document new API

### Phase 2: Core (4% → 64%) ✅ COMPLETE
- ✅ Fix errcheck violations - All fmt errors checked
- ✅ Implement Guard API - Single-step initialization
- ✅ Add compile-time validation - Panic on invalid commands

### Phase 3: Polish (20% → 80%) 🔄 IN PROGRESS
- ✅ Remove orphaned packages - Simplified architecture
- 🔄 Improve test coverage - Target 80%+ for config
- ⏳ Add integration tests - End-to-end validation
- ⏳ Create examples directory - Working examples
- ⏳ Add justfile - Standardize build commands
- ⏳ Fix code duplication - N/A after cleanup

### Phase 3: Polish (20% → 80%) ⏳ PENDING
- ⏳ Add integration tests - End-to-end validation
- ⏳ Create examples directory - Working examples
- ⏳ Add justfile - Standardize build commands
- ⏳ Fix code duplication - 7 clone groups detected

---

## Common Tasks

### Adding a New Command

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
        Use:   "newcmd",
        Short: "Brief description",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation
            return nil
        },
    })
    
    root.ExecuteAndExit(context.Background())
}
```

**Note:** `AddCommand` will panic if the command has no handler and no subcommands.
This is intentional - errors are caught at startup, not at runtime.

### Adding a Service

1. Create constructor: `func NewService(i do.Injector) (*Service, error)`
2. Register in `di/module.go`: `do.Provide(m.injector, NewService)`
3. Add invoke method: `func (m *Module) InvokeService() (*Service, error)`
4. Add health check if needed: `func (s *Service) HealthCheck() error`

### Adding a Validation Rule

1. Add method to `Validator` struct in `internal/validation/validator.go`
2. Call from `ValidateAll()` or create new public method
3. Add tests in `validator_test.go`

---

## Architecture Decision Records

### Using samber/do/v2 over other DI frameworks

- Lightweight, Go-idiomatic
- Constructor injection pattern
- Built-in lifecycle hooks (HealthCheck, Shutdown)
- No code generation required

### Validation at Runtime vs Compile-time

**Before:** Runtime validation after initialization (framework approach)  
**Now:** Panic at construction time (guard approach) ✅ IMPLEMENTED

Rationale: Go lacks compile-time macros; panic at init is closest to "fail fast" philosophy.

The Guard API (`GuardedCommand`) panics if:
- Command has no handler and no subcommands
- Strict mode requires RunE but Run provided
- Commands added after execution begins

### Koanf over Viper

- Lighter weight
- Explicit provider loading
- Better separation of concerns
- No global state

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
    // Single-step initialization
    root := cmdguard.New("myapp", "My application")
    
    // Add command (panics if invalid)
    root.AddCommand(&cobra.Command{
        Use:   "hello",
        Short: "Say hello",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Hello, World!")
        },
    })
    
    // Execute
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
3. Commands added after `Execute()` called
4. Command has no name

This ensures errors are caught immediately at startup, not when users run commands.

---

## Links & References

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [koanf Documentation](https://github.com/knadh/koanf)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Original Transformation Plan](./docs/planning/2026-02-14_04-21_CMDGUARD_TRANSFORMATION_PLAN.md)
- [Comprehensive Execution Plan](./docs/planning/2026-02-14_09-27-COMPREHENSIVE_EXECUTION_PLAN.md)
- [Feature Status](./FEATURES.md)
