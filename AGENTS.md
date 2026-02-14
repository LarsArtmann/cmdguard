# AGENTS.md - cmdguard Project Guide

**Last Updated:** 2026-02-14  
**Project:** cmdguard - CLI Validation Library  
**Go Version:** 1.26.0

---

## Quick Start

```bash
# Build the CLI application
go build -o cmdguard ./cmd/cmdguard

# Run tests
go test ./...
go test -v ./internal/...          # Verbose output
go test -cover ./...               # With coverage

# Run the application
./cmdguard --help
./cmdguard validate
./cmdguard version
```

---

## Project Overview

**cmdguard** is a Go library/framework for building validated Cobra CLI applications. It provides:

- **Lifecycle management** - Structured init → validate → execute → shutdown flow
- **Dependency injection** - Built on `samber/do/v2`
- **Configuration management** - Integrated with `knadh/koanf`
- **Runtime validation** - Ensures commands have handlers and flags are properly bound

**Current Status:** Work in progress. The project is transitioning from a framework approach to a "guard library" approach with compile-time enforcement. See `/docs/planning/` for transformation plan.

---

## Project Structure

```
cmdguard/
├── cmd/cmdguard/           # CLI application entry point (planned for removal)
├── pkg/cmdguard/           # Public API - main library entry point
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

### Package Guidelines

| Package | Purpose | Importable? |
|---------|---------|-------------|
| `pkg/cmdguard` | Public API | Yes |
| `internal/commands` | Command registry | No |
| `internal/config` | Configuration | No |
| `internal/di` | DI container setup | No |
| `internal/logging` | Logging utilities | No |
| `internal/validation` | Validation logic | No |

**Critical:** Never import `internal/` packages from outside the module.

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

### DI Patterns (samber/do/v2)

**Correct:** Constructor injection
```go
func NewRegistry(i do.Injector) (*Registry, error) {
    cfg, err := do.Invoke[*config.Config](i)
    if err != nil {
        return nil, fmt.Errorf("failed to invoke config: %w", err)
    }
    return &Registry{cfg: cfg}, nil
}
```

**Avoid:** Manual wiring after creation
```go
// DON'T DO THIS
registry := module.MustInvokeRegistry()
validator := module.MustInvokeValidator()
registry.SetValidator(validator)  // Manual wiring!
```

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

1. **Unused import in root.go** - `log/slog` imported but not used (L7)
2. **Missing dependency** - `internal/logging/logger.go` imports `github.com/charmbracelet/log` which is not in go.mod
3. **Errcheck violations** - Several `fmt.Fprintln`/`fmt.Fprintf` errors not checked
4. **Code duplication** - Version command logic duplicated in `main.go` and `root.go`

---

## Planned Changes

Per `docs/planning/2026-02-14_04-21_CMDGUARD_TRANSFORMATION_PLAN.md`:

1. **Remove `cmd/` folder** - Establish as library, not application
2. **Redesign public API** - Single-step initialization, panic on invalid
3. **Add compile-time validation** - Intercept Cobra calls to enforce correctness
4. **Fix DI usage** - Remove manual service linking
5. **Improve test coverage** - Target 80%+ for all packages
6. **Add justfile** - Standardize build commands

---

## Common Tasks

### Adding a New Command

```go
// In registry.go or setupCommands()
registry.AddCommand(&cobra.Command{
    Use:   "newcmd",
    Short: "Brief description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
})
```

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

**Current:** Runtime validation after initialization  
**Planned:** Panic at construction time (guard approach)  

Rationale: Go lacks compile-time macros; panic at init is closest to "fail fast" philosophy.

### Koanf over Viper

- Lighter weight
- Explicit provider loading
- Better separation of concerns
- No global state

---

## Links & References

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [koanf Documentation](https://github.com/knadh/koanf)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Transformation Plan](./docs/planning/2026-02-14_04-21_CMDGUARD_TRANSFORMATION_PLAN.md)
