# AGENTS.md - cmdguard Project Guide

**Last Updated:** 2026-04-17
**Project:** cmdguard - CLI Guard Library
**Go Version:** 1.26
**Status:** v2.1.0 - Production Ready

---

## Quick Start

```bash
# Run tests (all packages, with race detection)
go test ./... -count=1 -timeout 120s -race

# Build all
go build ./...

# Lint (golangci-lint 2.x)
golangci-lint run ./...

# Format
golangci-lint fmt ./...

# Coverage
go test ./... -count=1 -timeout 120s -cover
```

**Important:** `git commit --no-verify` is required (pre-commit hooks have pre-existing errors).

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with type-safe dependency injection.

| API  | Package           | Use Case                         |
| ---- | ----------------- | -------------------------------- |
| v2   | `pkg/cmdguard/v2` | Type-safe, DI-powered, no panics |

**Current Status:** v2.1.0. All packages tested, 0 lint issues.

---

## Project Structure

```
cmdguard/
├── pkg/cmdguard/
│   ├── v2/                       # v2 API (recommended)
│   │   ├── cli.go                # CLI[T] struct, NewCLI, AddCommand, Execute
│   │   ├── cli_accessors.go      # CLI accessor methods (Config, Scope, etc.)
│   │   ├── cli_command.go        # Internal cobra wiring (cliToCobraCommand)
│   │   ├── cli_options.go        # CLI functional options (WithCLIVersion, etc.)
│   │   ├── command.go            # Command[T,F] struct, constructors, options, Validate
│   │   ├── config.go             # Config type constraint
│   │   ├── config_parsing.go     # ParseFlagTags, DefaultValue
│   │   ├── config_setfield.go    # SetField for config structs
│   │   ├── errors.go             # Sentinel errors and error types
│   │   ├── flags.go              # FlagRegistry with struct tags
│   │   ├── flags_parse.go        # Flag parsing logic
│   │   ├── flags_suggest.go      # Typo suggestions (Levenshtein)
│   │   ├── flags_validate.go     # Flag validation
│   │   ├── flag_helpers.go       # Flag type constraints, cloning, parsing helpers
│   │   ├── flow_context.go       # BranchingFlowContext for command path tracking
│   │   ├── middleware.go         # Middleware chain pattern
│   │   ├── scope.go              # DI scope wrapping samber/do/v2
│   │   ├── type_helpers.go       # Generic type helpers
│   │   ├── types_duration.go     # Duration type
│   │   ├── types_email.go        # Email type
│   │   ├── types_enum.go         # Enum[T] type
│   │   ├── types_filepath.go     # FilePath type
│   │   ├── types_hostport.go     # HostPort type
│   │   ├── types_log.go          # LogLevel type
│   │   ├── types_port.go         # Port type
│   │   └── types_url.go          # URL type
├── pkg/testutil/
│   └── panic_test_helpers.go     # Shared test assertions
├── examples/
│   ├── basic/                    # Simple v2 demo
│   ├── typed/                    # v2 API demo with DI and lifecycle
│   ├── di/                       # DI-focused example
│   ├── advanced-flags/           # Advanced flag types
│   ├── validation/               # Validation patterns example
│   └── internal/                 # Shared example helpers
├── benchmarks/                   # Performance benchmarks
├── tests/integration/            # Integration tests
├── docs/                         # Documentation
├── AGENTS.md                     # This file
├── FEATURES.md                   # Feature status
├── TODO_LIST.md                  # Remaining tasks
├── .golangci.yml                 # Lint configuration
└── README.md                     # User documentation
```

### Package Guidelines

| Package           | Purpose          | Importable? | Coverage |
| ----------------- | ---------------- | ----------- | -------- |
| `pkg/cmdguard/v2` | Type-safe API    | Yes         | 82.2%    |
| `pkg/testutil`    | Test helpers     | Yes         | —        |

---

## Key Dependencies

| Library                     | Purpose              | Version |
| --------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`    | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`   | Dependency injection | v2.0.0  |
| `github.com/spf13/pflag`    | Flag parsing         | v1.0.6  |
| `charm.land/fang/v2`        | Cobra styling        | v2.0.1  |

---

## API Reference

### Architecture: CLI[T] + Command[T, F]

`CLI[T]` has one type parameter (config type). Each command gets its own flags type via `Command[T, F]`. Because Go doesn't support additional type parameters on methods, `AddCommand` is a standalone function.

Commands are created via constructors — `NewCommand` for leaf commands, `NewParentCommand` for commands with subcommands. Struct fields are unexported to enforce validation at construction time.

```go
cli, err := v2.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})
cmd, err := v2.NewCommand[AppConfig, *GreetFlags]("greet", greetHandler,
    v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v2.AddCommand(cli, cmd)
```

### Command Constructors

```go
// Leaf command with handler
func NewCommand[T, F any](use string, runE func(ctx context.Context, cfg *T, flags F) error, opts ...CommandOption[T, F]) (Command[T, F], error)

// Parent command with subcommands
func NewParentCommand[T, F any](use string, long string, subcommands []Command[T, F], opts ...CommandOption[T, F]) (Command[T, F], error)

// Panic variants (for compile-time-known config)
func MustNewCommand[T, F any](...) Command[T, F]
func MustNewParentCommand[T, F any](...) Command[T, F]
```

### Command Options

| Option                           | Purpose                                |
| -------------------------------- | -------------------------------------- |
| `WithShort[T, F](short)`         | Short description                      |
| `WithLong[T, F](long)`           | Long description                       |
| `WithAliases[T, F](aliases...)`  | Alternative names                      |
| `WithExample[T, F](example)`     | Example usage                          |
| `WithFlags[T, F](flags)`         | Typed flags struct                     |
| `WithRunE[T, F](runE)`           | Main handler (required for NewCommand) |
| `WithPreRunE[T, F](preRunE)`     | Pre-validation hook                    |
| `WithPostRunE[T, F](postRunE)`   | Post-success cleanup hook              |
| `WithSubcommands[T, F](cmds...)` | Child commands                         |
| `WithHidden[T, F](hidden)`       | Hide from help                         |
| `WithDeprecated[T, F](msg)`      | Deprecation message                    |
| `WithGroupID[T, F](group)`       | Help group name                        |

### CLI[T] Constructor

```go
cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
```

Functional options:

| Option                   | Purpose                                     |
| ------------------------ | ------------------------------------------- |
| `WithCLIVersion[T](v)`   | Set version string                          |
| `WithCLILong[T](desc)`   | Set long description                        |
| `WithCLIScope[T](scope)` | Set custom DI scope                         |
| `WithSilenceErrors[T]()` | Suppress cobra error printing               |
| `WithSilenceUsage[T]()`  | Suppress usage on error                     |
| `WithColor[T](bool)`     | Enable/disable fang styling (default: true) |

### CLI[T] Methods

| Method                        | Returns                 | Purpose                           |
| ----------------------------- | ----------------------- | --------------------------------- |
| `Execute(ctx)`                | `error`                 | Run CLI with context              |
| `ExecuteWithArgs(ctx, args)`  | `error`                 | Run with specific args            |
| `ExecuteAndExit(ctx)`         |                         | Run and os.Exit(1) on error       |
| `Scope()`                     | `*Scope`                | DI scope                          |
| `Injector()`                  | `do.Injector`           | Raw samber/do injector            |
| `Config()`                    | `*T`                    | Typed config                      |
| `SetConfig(cfg)`              |                         | Update config                     |
| `RootCommand()`               | `*cobra.Command`        | Underlying cobra command          |
| `Shutdown(ctx)`               | `error`                 | Graceful shutdown                 |
| `HealthCheck()`               | `error`                 | Run health checks                 |
| `HealthCheckWithContext(ctx)` | `error`                 | Health checks with context        |
| `SetVersion(v)`               |                         | Set version at runtime            |
| `SetLong(desc)`               |                         | Set long description              |
| `FlowContext()`               | `*BranchingFlowContext` | Path tracking (nil until Execute) |
| `AddGlobalFlag(...)`          |                         | Persistent string flag            |
| `AddGlobalBoolFlag(...)`      |                         | Persistent bool flag              |

### Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
    Output  string `flag:"output" short:"o" default:"text" help:"Output format"`
}

func main() {
    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
    if err != nil {
        panic(err)
    }

    cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("hello",
        func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            fmt.Printf("Hello! Verbose: %v\n", cfg.Verbose)
            return nil
        },
        v2.WithShort[AppConfig, v2.NoFlags]("Say hello"),
    )
    if err != nil {
        panic(err)
    }

    if err := v2.AddCommand(cli, cmd); err != nil {
        panic(err)
    }

    if err := cli.Execute(context.Background()); err != nil {
        fmt.Println("Error:", err)
    }
}
```

### Command with Custom Flags

```go
type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Count uint   `flag:"count" short:"c" default:"1"    help:"Number of greetings"`
    Shout bool   `flag:"shout" default:"false"          help:"Shout the greeting"`
}

greetCmd, err := v2.NewCommand[AppConfig, *GreetFlags]("greet",
    func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        for i := uint(0); i < flags.Count; i++ {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
        }
        return nil
    },
    v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v2.AddCommand(cli, greetCmd)
```

### Subcommands

```go
listCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("list",
    listUsersHandler, v2.WithShort[AppConfig, v2.NoFlags]("List users"),
)
createCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("create",
    createUserHandler, v2.WithShort[AppConfig, v2.NoFlags]("Create user"),
)
userCmd, err := v2.NewParentCommand[AppConfig, v2.NoFlags]("user",
    "User management", []v2.Command[AppConfig, v2.NoFlags]{listCmd, createCmd},
    v2.WithShort[AppConfig, v2.NoFlags]("User management"),
)
v2.AddCommand(cli, userCmd)
```

### Dependency Injection

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
scope := cli.Scope()

// Register services
v2.Provide(scope, func(i do.Injector) (*Database, error) {
    cfg, _ := v2.Invoke[*AppConfig](scope)
    return &Database{DSN: cfg.DSN}, nil
})
v2.ProvideValue(scope, &Logger{Level: "info"})

// Invoke in command handlers
db, err := v2.Invoke[*Database](cli.Scope())
```

### Lifecycle Hooks

```go
cmd, err := v2.NewCommand[AppConfig, *Flags]("example", runHandler,
    v2.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // validation
    }),
    v2.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // cleanup (only called on success)
    }),
)
```

### BranchingFlowContext

Automatically created on `Execute`. Access via `GetBranchingFlowContext(ctx)` in handlers.

```go
bfc, ok := v2.GetBranchingFlowContext(ctx)
bfc.PathString()  // "app.subcmd"
bfc.SetValue(key, val)  // propagates to children
bfc.GetValue(key)       // looks up hierarchy
```

### Error Handling

```go
// All v2 functions return errors
cli, err := v2.NewCLI[Config]("app", "My app", Config{})
cmd, err := v2.NewCommand[Config, NoFlags]("test", handler)

// Sentinel errors for errors.Is()
errors.Is(err, v2.ErrInvalidCommand)
errors.Is(err, v2.ErrMissingName)
errors.Is(err, v2.ErrDuplicateCommand)
errors.Is(err, v2.ErrMissingHandler)

// Rich error types
v2.NewCommandError(name, err)    // wraps with command context
v2.NewServiceError(type, err)    // wraps with DI service context
v2.NewFlagError(name, err)       // wraps with flag context
v2.NewFlagErrorWithSuggestion(name, err, suggestion)  // includes typo fix
```

---

## Coding Standards

### Go Conventions

- **Go 1.26** - Use modern Go features
- **gofumpt** formatting (via `golangci-lint fmt`)
- **Error handling** - Always check errors, wrap with `fmt.Errorf("context: %w", err)`
- **No panics** in v2 library code
- **Functional options** pattern for configuration
- **Constructor pattern** - All Command creation via `NewCommand`/`NewParentCommand`, struct fields unexported

### Testing

- `t.Parallel()` in every test function and subtest (paralleltest linter)
- `//nolint:paralleltest` for tests using `t.Setenv` or capturing `os.Stdout`
- `//nolint:fatcontext` at file level for test files with context in closures
- Table-driven tests: `tests := []struct{...}` pattern
- Two test packages: `v2` (internal, access private helpers) and `v2_test` (external)

### Test Commands

```bash
go test ./... -count=1 -timeout 120s -race     # All tests with race detection
go test ./... -count=1 -timeout 120s -cover     # Coverage report
golangci-lint run ./...                          # Lint (0 issues)
go build ./...                                   # Verify build
```

---

## Architecture Decisions

### v2.1 Design Principles

1. **Single type parameter** - `CLI[T]` only parameterizes on config; each command has its own flags type
2. **No Panics** - All operations return errors
3. **DI-Powered** - samber/do/v2 for dependency injection
4. **Typed Flags** - Struct tags for flag definitions
5. **Standalone AddCommand** - Function (not method) to support per-command flag types
6. **Constructor validation** - Commands validated at construction, struct fields unexported

### Key Gotchas

1. `t.Setenv` + `t.Parallel()` = panic — use `//nolint:paralleltest`
2. `PostRunE` is NOT called when `RunE` errors (Cobra behavior)
3. `NoFlags` is `type NoFlags = struct{}` — use `(NoFlags{})` with parens for comparisons
5. fang provides styled output by default; `WithColor(false)` falls back to plain cobra
6. `AddCommand` calls `cmd.Validate()` as defense-in-depth even though constructors already validate

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
