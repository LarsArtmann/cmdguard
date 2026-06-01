# AGENTS.md - cmdguard Contributor & AI Assistant Guide

> **Note:** This file serves as both a contributor guide and context for AI-assisted development. It documents architecture decisions, API reference, coding standards, and known gotchas.

**Last Updated:** 2026-06-01
**Project:** cmdguard - CLI Guard Library
**Go Version:** 1.26
**Status:** v2.3.0-dev - 333 tests (1084 cases), 83.6% coverage, 0 lint issues, 0 race conditions

---

## Quick Start

```bash
# Enter dev shell (Go 1.26, gopls, golangci-lint)
nix develop

# Run tests (all packages, with race detection)
go test ./... -count=1 -timeout 120s -race

# Build all
go build ./...

# Lint (golangci-lint 2.x)
golangci-lint run ./...

# Format (nix + go via treefmt)
nix fmt

# Coverage
go test ./... -count=1 -timeout 120s -cover

# Check everything (format check)
nix flake check
```

**Important:** `git commit --no-verify` is required (pre-commit hooks have pre-existing errors).

---

## Project Overview

**cmdguard** is a Go library for building validated Cobra CLI applications with type-safe dependency injection.

| API | Package           | Use Case                         |
| --- | ----------------- | -------------------------------- |
| v2  | `pkg/cmdguard/v2` | Type-safe, DI-powered, no panics |

**Current Status:** v2.3.0-dev. 272 tests passing, 83.3% coverage, 0 build errors.

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
│   │   ├── command_suggest.go    # Command typo suggestions
│   │   ├── config.go             # Config type constraint
│   │   ├── config_file.go        # ConfigFileLoader, JSON loader, WithConfigFile
│   │   ├── config_parsing.go     # ParseFlagTags, DefaultValue
│   │   ├── config_setfield.go    # SetField for config structs
│   │   ├── configload/           # Optional YAML/TOML loaders
│   │   ├── counting_flag.go      # Counting flag support (count:"true")
│   │   ├── editor.go             # EditInEditor ($EDITOR support)
│   │   ├── errors.go             # Sentinel errors and error types
│   │   ├── flags.go              # FlagRegistry with struct tags
│   │   ├── flags_parse.go        # Flag parsing logic
│   │   ├── flags_suggest.go      # Typo suggestions (Levenshtein)
│   │   ├── flags_validate.go     # Flag validation
│   │   ├── flag_helpers.go       # Flag type constraints, cloning, parsing helpers
│   │   ├── flow_context.go       # BranchingFlowContext for command path tracking
│   │   ├── middleware.go         # Middleware chain pattern
│   │   ├── scope.go              # DI scope wrapping samber/do/v2
│   │   ├── output.go             # Rich output (table/json/csv/yaml)
│   │   ├── type_handler.go       # Extensible type registry
│   │   ├── type_helpers.go       # Generic type helpers
│   │   ├── version.go            # VersionCommand helper
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
│   └── taskctl/                   # Single superb example: production task manager CLI
│       ├── main.go                # CLI construction, DI setup, all CLI options
│       ├── commands.go            # All command definitions with options
│       ├── types.go               # Config, flags, domain types, TaskStore service
│       ├── main_test.go           # Comprehensive integration tests (~66 tests)
│       └── README.md              # Feature matrix and usage guide
├── benchmarks/                   # Performance benchmarks
├── tests/integration/            # Integration tests
├── docs/                         # Documentation
├── AGENTS.md                     # This file
├── FEATURES.md                   # Feature status
├── TODO_LIST.md                  # Remaining tasks
├── .golangci.yml                 # Lint configuration
├── flake.nix                     # Nix dev shell, formatter, checks
├── flake.lock                    # Nix lock file
└── README.md                     # User documentation
```

### Package Guidelines

| Package           | Purpose       | Importable? | Coverage |
| ----------------- | ------------- | ----------- | -------- |
| `pkg/cmdguard/v2` | Type-safe API | Yes         | ~84%     |
| `pkg/testutil`    | Test helpers  | Yes         | —        |

---

## Key Dependencies

| Library                            | Purpose              | Version |
| ---------------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`           | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`          | Dependency injection | v2.0.0  |
| `github.com/spf13/pflag`           | Flag parsing         | v1.0.10 |
| `charm.land/fang/v2`               | Cobra styling        | v2.0.1  |
| `charm.land/huh/v2`                | Interactive prompts  | v2.0.3  |
| `charm.land/glamour/v2`             | Markdown rendering   | v2.0.0  |
| `go.opentelemetry.io/otel/trace`   | OpenTelemetry spans  | v1.44.0 |
| `github.com/larsartmann/go-output` | Rich output formats  | latest  |

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

| Option                           | Purpose                                              |
| -------------------------------- | ---------------------------------------------------- |
| `WithShort[T, F](short)`         | Short description                                    |
| `WithLong[T, F](long)`           | Long description                                     |
| `WithAliases[T, F](aliases...)`  | Alternative names                                    |
| `WithExample[T, F](example)`     | Example usage                                        |
| `WithFlags[T, F](flags)`         | Typed flags struct                                   |
| `WithRunE[T, F](runE)`           | Main handler (required for NewCommand)               |
| `WithPreRunE[T, F](preRunE)`     | Pre-validation hook                                  |
| `WithPostRunE[T, F](postRunE)`   | Post-success cleanup hook                            |
| `WithSubcommands[T, F](cmds...)` | Child commands                                       |
| `WithHidden[T, F](hidden)`       | Hide from help                                       |
| `WithDeprecated[T, F](msg)`      | Deprecation message                                  |
| `WithGroupID[T, F](group)`       | Help group name                                      |
| `WithExactArgs[T, F](n)`         | Require exactly n positional args                    |
| `WithMinimumArgs[T, F](n)`       | Require at least n positional args                   |
| `WithMaximumArgs[T, F](n)`       | Allow at most n positional args                      |
| `WithRangeArgs[T, F](min, max)`  | Require between min and max args                     |
| `WithNoArgs[T, F]()`             | Reject any positional args                           |
| `WithArgs[T, F](fn)`             | Custom cobra.PositionalArgs validator                |
| `WithPromptOnMissing[T, F]()`    | Interactive prompt for missing `prompt`-tagged flags |

### CLI[T] Constructor

```go
cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
```

Functional options:

| Option                         | Purpose                                     |
| ------------------------------ | ------------------------------------------- |
| `WithCLIVersion[T](v)`         | Set version string                          |
| `WithCLILong[T](desc)`         | Set long description                        |
| `WithCLIScope[T](scope)`       | Set custom DI scope                         |
| `WithSilenceErrors[T]()`       | Suppress cobra error printing               |
| `WithSilenceUsage[T]()`        | Suppress usage on error                     |
| `WithColor[T](bool)`           | Enable/disable fang styling (default: true) |
| `WithFang[T](bool)`            | Enable/disable fang styling (preferred)     |
| `WithFangOptions[T](opts...)`  | Custom fang options                         |
| `WithMiddleware[T](mw...)`     | Middleware wrapping every handler           |
| `WithGroup[T](id, title)`      | Register command group on root              |
| `WithEnvPrefix[T](prefix)`     | Prefix for env var lookups                  |
| `WithSignalHandling[T]()`      | Cancel context on SIGINT/SIGTERM            |
| `WithConfigValidation[T](fn)`  | Validate config after flag parsing          |
| `WithStrictValidation[T]()`    | Require short descriptions on all commands  |
| `WithDraconianValidation[T]()` | Strict + examples on leaf commands          |
| `WithConfigFile[T](paths...)`  | Load JSON config file before flags          |
| `WithConfigFileLoader[T](l,p)` | Load config with custom loader (YAML/TOML)  |
| `WithGlamourHelp[T]()`         | Render markdown in command help text        |
| `WithTelemetry[T](tracer)`     | OpenTelemetry spans for all commands        |

### CLI[T] Methods

| Method                        | Returns                 | Purpose                              |
| ----------------------------- | ----------------------- | ------------------------------------ |
| `Execute(ctx)`                | `error`                 | Run CLI with context                 |
| `ExecuteWithArgs(ctx, args)`  | `error`                 | Run with specific args               |
| `ExecuteAndExit(ctx)`         |                         | Run and os.Exit (respects ExitCoder) |
| `Scope()`                     | `*Scope`                | DI scope                             |
| `Injector()`                  | `do.Injector`           | Raw samber/do injector               |
| `Config()`                    | `*T`                    | Typed config                         |
| `SetConfig(cfg)`              |                         | Update config                        |
| `RootCommand()`               | `*cobra.Command`        | Underlying cobra command             |
| `Shutdown(ctx)`               | `error`                 | Graceful shutdown                    |
| `HealthCheck()`               | `error`                 | Run health checks                    |
| `HealthCheckWithContext(ctx)` | `error`                 | Health checks with context           |
| `SetVersion(v)`               |                         | Set version at runtime               |
| `SetLong(desc)`               |                         | Set long description                 |
| `FlowContext()`               | `*BranchingFlowContext` | Path tracking (nil until Execute)    |
| `AddGlobalFlag(...)`          |                         | Persistent string flag               |
| `AddGlobalBoolFlag(...)`      |                         | Persistent bool flag                 |

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

### Middleware

```go
// Timing middleware — logs execution duration
cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
    v2.WithMiddleware(v2.TimingMiddleware[Config](func(name string, d time.Duration) {
        log.Printf("%s took %v", name, d)
    })),
)

// Recovery middleware — catches panics
cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
    v2.WithMiddleware(v2.RecoveryMiddleware[Config]()),
)

// Spinner middleware — shows terminal spinner
cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
    v2.WithMiddleware(v2.SpinnerMiddleware[Config]("Loading...")),
)

// Telemetry middleware — OpenTelemetry spans
import "go.opentelemetry.io/otel"
tracer := otel.Tracer("myapp")
cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
    v2.WithTelemetry[Config](tracer), // or WithMiddleware(TelemetryMiddleware[Config](tracer))
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

// Exit codes
v2.NewExitError(code, err)       // error with custom exit code for ExecuteAndExit
errors.As(err, &exitCoder)       // check if error implements ExitCoder
```

### Version Command

```go
cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
    v2.WithCLIVersion[Config]("1.0.0"),
)
v2.AddCommand(cli, v2.MustVersionCommand[Config](cli))
```

### Markdown Help (glamour)

```go
cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
    v2.WithGlamourHelp[Config](),
)
// Command Long and Example fields are rendered as markdown in terminal help
```

### Strict Validation

```go
cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
    v2.WithStrictValidation[Config](),  // requires WithShort on all commands
)
```

### Config Validation

```go
cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
    v2.WithConfigValidation[Config](func(cfg *Config) error {
        if cfg.Port < 1 { return fmt.Errorf("invalid port") }
        return nil
    }),
)
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

### v2.3 Design Principles

1. **Single type parameter** - `CLI[T]` only parameterizes on config; each command has its own flags type
2. **No Panics** - All operations return errors
3. **DI-Powered** - samber/do/v2 for dependency injection
4. **Typed Flags** - Struct tags for flag definitions
5. **Standalone AddCommand** - Function (not method) to support per-command flag types
6. **Env tags** - `env:"VAR_NAME"` struct tag reads from environment
7. **Counting flags** - `count:"true"` tag enables -v/-vv/-vvv pattern
8. **Signal handling** - `WithSignalHandling[T]()` for graceful shutdown
9. **Rich output** - OutputTable/OutputResult with 12+ formats
10. **Instance-scoped registries** — Each `FlagRegistry` clones from package-level defaults; `RegisterTypeHandler()`/`RegisterValidator()` write to global template; `FlagRegistry.RegisterTypeHandler()`/`FlagRegistry.RegisterFlagValidator()` write to instance
11. **$EDITOR support** - `EditInEditor()` for user input editing
12. **Typo suggestions** - `SuggestFlag`/`SuggestCommand` with Levenshtein
13. **ValidationMode enum** - `Lenient`/`Strict`/`Draconian` spectrum, `>=` comparison
14. **Full sentinel coverage** - All 40+ errors identifiable via `errors.Is()`
15. **Generic helpers** - `textMarshal[T]`/`textUnmarshal[T]`, `renderAndWrite`/`marshalAndWrite`, `branchWithCtx`
16. **Spinner middleware** — `SpinnerMiddleware[T](title)` shows a lipgloss-styled spinner on stderr; skips when not a terminal
17. **Glamour help** — `WithGlamourHelp[T]()` renders command `Long` and `Example` fields via `charm.land/glamour/v2` markdown; uses `RenderWithEnvironmentConfig` (checks `GLAMOUR_STYLE` env var, defaults to `"dark"`); applied recursively to all commands at registration time
18. **Telemetry middleware** — `TelemetryMiddleware[T](tracer)` creates an OpenTelemetry span per command; requires `go.opentelemetry.io/otel/trace.Tracer`; `WithTelemetry[T](tracer)` is the convenience CLI option

### Key Gotchas

1. `t.Setenv` + `t.Parallel()` = panic — use `//nolint:paralleltest`
2. `PostRunE` is NOT called when `RunE` errors (Cobra behavior)
3. `NoFlags` is `type NoFlags = struct{}` — use `(NoFlags{})` with parens for comparisons
4. fang provides styled output by default; `WithFang(false)` falls back to plain cobra
5. `AddCommand` calls `cmd.Validate()` as defense-in-depth even though constructors already validate
6. **envPrefix propagation** — `WithEnvPrefix` sets prefix on root AND command-level flags (fixed in v2.2)
7. **Counting flags** — must use `int` type with `count:"true"` tag; don't reuse flag names from root config
8. **Prompt tag** — `prompt:"Question?"` on a struct field enables interactive prompting when the flag is missing and `WithPromptOnMissing` is set on the command. Bool fields use `huh.NewConfirm`, enum fields (with `values` tag) use `huh.NewSelect`, all others use `huh.NewInput`
9. **SuggestFlag API** — returns `(string, bool)` since v2.2 (breaking change from string-only)
10. **Instance-scoped registries** — `FlagRegistry` clones `typeRegistry` and `validatorRegistry` from globals at creation time; package-level `RegisterTypeHandler()`/`RegisterValidator()` write to the global defaults template, not to existing instances. Use `FlagRegistry.RegisterTypeHandler()` for per-instance customization.
11. **go-output local replace** — uses absolute local path in go.mod, blocks CI/other developers
12. **Deprecated APIs (remove in v3)** — `IsExecutable()` → use `HasHandler()`. `FlowContextAccessor` was removed in v2.3.0 — use `GetBranchingFlowContext(ctx)` directly
13. **Typed branching alternatives** — prefer `BranchWithDuration(name, time.Duration)` and `BranchWithDeadlineTime(name, time.Time)` over string-based `BranchWithTimeout`/`BranchWithDeadline`
14. **Regex validation cache** — `validateRegex` caches compiled patterns in `sync.Map`; global state, tests must not run in parallel
15. **Exit codes** — `ExecuteAndExit` checks for `ExitCoder` interface; use `NewExitError(code, err)` for custom exit codes
16. **Strict validation** — `WithStrictValidation[T]()` requires `WithShort` on all commands; enforced at `AddCommand` time
17. **Draconian validation** — `WithDraconianValidation[T]()` is superset of strict + requires `WithExample` on leaf commands; parent commands are exempt
18. **Config validation** — `WithConfigValidation[T](fn)` runs after root flag parsing but before any command handler; blocks execution on error
19. **Args validation** — `WithExactArgs`/`WithMinimumArgs`/etc. use cobra's built-in arg validators; runs during command execution, not at registration
20. **Spinner non-TTY** — `SpinnerMiddleware` auto-skips when `os.Stderr` is not a terminal; use `SpinnerConfig{Writer: ...}` to override
21. **Glamour v2 env-based theme** — `WithGlamourHelp[T]()` now uses `RenderWithEnvironmentConfig` which checks `GLAMOUR_STYLE` env var, defaulting to "dark"; the string `"auto"` is no longer a valid glamour theme name in v2
22. **Telemetry context propagation** — `TelemetryMiddleware` starts a span but cannot propagate the new context to the handler due to the `next func() error` middleware API signature; child spans must use the original context passed to the handler
23. **FullPath populated at execution time** — `CommandInfo.FullPath` is set via `cobra.CommandPath()` inside the handler closure, NOT at command registration; it's empty in unit tests unless you call the handler through a cobra execution
24. **Glamour idempotent** — `applyGlamourIfEnabled` resets `glamourHelp=false` after applying to prevent double-rendering (which would wrap ANSI codes inside ANSI codes); calling Execute twice is safe
25. **outputEnabled removed** — The unused `outputEnabled` field was removed from `CLI[T]`; use `outputFormat != ""` to detect if output formatting is configured
26. **NewExitError returns (\***ExitError**, **error\**) — validates 0-255 range; breaking change from `*ExitError`
27. **NewScopeFromInjector returns (\***Scope**, **error\*\*) — nil injector returns error; breaking change from nil dereference
26. **Sentinel wrapping** — All 40+ errors use `fmt.Errorf("%w: ...", sentinel)` for `errors.Is()` chainability
27. **Config file precedence** — `WithConfigFile[T](paths...)` loads config BEFORE flag registration; config values become the new tag defaults, so flags/env still override them
28. **Config file paths** — supports `$ENV` expansion and `~` expansion; missing files are silently skipped
29. **Config file `--config` override** — if the config struct has a `config` flag, its value overrides the default search paths from `WithConfigFile`
30. **Config file flat only (v1)** — JSON/YAML/TOML loaders detect top-level keys matching `flag` tag names; nested structs in config files are not yet supported
31. **Nix sandbox vs local replace** — `go.mod` has `replace` directives pointing to `../go-output`; Nix sandboxed checks (`buildGoModule`, `go build` in derivations) cannot resolve these, so `flake.nix` only provides devShell + formatter + format check (no build/vet checks)
32. **Glamour v2 no `"auto"` theme** — `charm.land/glamour/v2` removed the `"auto"` theme; use empty string (env-based via `GLAMOUR_STYLE`) or explicit theme like `"dark"`; `WithGlamourHelp` now sets theme to `""` for env-based detection

---

## Links

- [Cobra Documentation](https://github.com/spf13/cobra)
- [samber/do/v2 Documentation](https://github.com/samber/do)
- [fang Documentation](https://github.com/charmbracelet/fang)
- [CLI Design Principles](./docs/CLI_DESIGN_PRINCIPLES.md)
- [Feature Status](./FEATURES.md)
- [TODO List](./TODO_LIST.md)
