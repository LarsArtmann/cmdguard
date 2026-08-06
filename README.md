# cmdguard

[![CI](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml)
[![Website](https://github.com/larsartmann/cmdguard/actions/workflows/website.yml/badge.svg)](https://github.com/larsartmann/cmdguard/actions/workflows/website.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/larsartmann/cmdguard/v4.svg)](https://pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4)
[![Coverage](https://img.shields.io/badge/coverage-87.8%25-brightgreen)](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**[Website](https://cmdguard.lars.software)** · **[Docs](https://cmdguard.lars.software/getting-started/installation/)** · **[pkg.go.dev](https://pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4)**

---

**The Go CLI framework that catches missing handlers, duplicate commands, and invalid flags at construction — not at 2am in production.**

cmdguard is the only Go CLI framework that unifies **type-safe flags**, **dependency injection with lifecycle management**, and a **zero-panic error contract** into a single system validated at construction. It wraps [Cobra](https://github.com/spf13/cobra) so you keep full compatibility while eliminating its footguns.

**Get started in 30 seconds:** `go get github.com/larsartmann/cmdguard/v4` · [Quick Start](#quick-start) · [Full Docs](https://cmdguard.lars.software)

> **API Stability:** v4.0.0 is the current release on the v4 major line. The legacy v2 line is in maintenance at v2.10.4. See [CHANGELOG.md](CHANGELOG.md) for the v3→v4 migration notes and the [v2→v3 Migration Guide](docs/MIGRATION_v2_v3.md) for the earlier migration.

---

## Why cmdguard?

Other Go CLI frameworks give you flags. cmdguard gives you flags **plus** everything production CLIs need: dependency injection, service lifecycle, health checks, graceful shutdown, styled output, and error handling that won't bite you.

### The trinity — what no other CLI framework offers together

| Capability                                        | Cobra | Kong                | urfave/cli | **cmdguard** |
| ------------------------------------------------- | ----- | ------------------- | ---------- | ------------ |
| Struct-tag flags (no string lookups)              | —     | Yes                 | —          | **Yes**      |
| Dependency injection (lazy services, lifecycle)   | —     | —                   | —          | **Yes**      |
| Graceful shutdown (reverse-order on SIGINT)       | —     | —                   | —          | **Yes**      |
| Health checks (`DoctorCommand`, `Healthchecker`)  | —     | —                   | —          | **Yes**      |
| Zero panics by construction (no `Run`, no `Must`) | —     | —                   | —          | **Yes**      |
| Validated at construction (not at runtime)        | —     | Some <sup>[1]</sup> | —          | **Yes**      |
| Error printed exactly once + correct exit codes   | —     | —                   | —          | **Yes**      |
| Styled output by default (fang + lipgloss)        | —     | —                   | —          | **Yes**      |
| Gradual migration (raw cobra + typed runtime)     | —     | —                   | —          | **Yes**      |

<sup>[1]</sup> Kong validates struct tags at parse time but does not validate command structure (missing handlers, duplicates) at registration.

### Dependency injection — the real differentiator

Register services, invoke them in handlers, and manage their entire lifecycle:

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "My production CLI", AppConfig{},
    v4.WithGracefulShutdown(), // SIGINT → reverse-order shutdown of all services
)

// Register a database service (lazy — created on first invoke)
v4.Provide(cli.Scope(), func(i do.Injector) (*Database, error) {
    return &Database{DSN: "postgres://..."}, nil
})

// Use it in any command handler
v4.NewCommand("query", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, _ v4.NoFlags) error {
        db, _ := v4.Invoke[*Database](cli.Scope())
        return db.Query(ctx)
    },
)
```

Services can implement `HealthCheck` (wired into `DoctorCommand`) and `Shutdown` (called on SIGINT/SIGTERM in reverse invocation order).

### Type-safe flags — validated at construction

**Raw Cobra — flags are strings, validated at runtime:**

```go
var name string
var count int
rootCmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")
rootCmd.Flags().IntVarP(&count, "count", "c", 1, "Number of greetings")
// Forgot to add "count"? Missing handler? Duplicate command? Runtime surprise.
```

**cmdguard — flags are typed structs, validated at construction:**

```go
type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Count int    `flag:"count" short:"c" default:"1"    help:"Number of greetings"`
}
// Missing handler? Duplicate command? Invalid name? Caught at AddCommand time.
```

### Zero panics — by construction

Every function returns errors. No `Run` (panics), no `Must*` variants. The error is printed exactly once (styled by [fang](https://github.com/charmbracelet/fang)), usage is silenced on error by default, and exit codes are handled automatically:

```go
cli.ExecuteAndExit(context.Background()) // one line, correct exit code
```

---

## Error handling & exit codes

cmdguard owns the error-output contract so you can't get it wrong:

- **The error is printed exactly once** — styled by [fang](https://github.com/charmbracelet/fang) when enabled (the default), or plain by Cobra when disabled.
- **Usage is never printed on error** (`SilenceUsage: true` by default). Use `WithoutSilenceUsage()` to re-enable usage-on-error.
- **The error returned by `Execute` is for exit-code mapping only — do not re-print it**, or you'll duplicate the output.

**Recommended — one line, correct exit code:**

```go
cli.ExecuteAndExit(context.Background())
```

**When you need to run code before exiting** (flush logs, export an audit log, tear
down resources), use `ExitCode` instead of `ExecuteAndExit`:

```go
err := cli.Execute(ctx)
// ...flush / export audit log / teardown...
os.Exit(v4.ExitCode(err)) // 0 on success, ExitCoder code or 1 on failure
```

> Pitfall to avoid: `if err := cli.Execute(ctx); err != nil { fmt.Fprintln(os.Stderr, err) }`
> re-prints the error that cmdguard already printed.

---

## Quick Start

```bash
go get github.com/larsartmann/cmdguard/v4
```

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
    Output  string `flag:"output" short:"o" default:"text" help:"Output format"`
}

type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Uppercase output"`
}

func main() {
    cli, err := v4.NewCLI[AppConfig]("myapp", "My CLI application", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
        os.Exit(1)
    }

    greetCmd, err := v4.NewCommand("greet", &GreetFlags{},
        func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
        v4.WithShort("Greet someone"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create command: %v\n", err)
        os.Exit(1)
    }

    v4.AddCommand(cli, greetCmd)
    cli.ExecuteAndExit(context.Background())
}
```

```bash
$ go run main.go greet -n "cmdguard" --shout
HELLO, CMDGUARD!
```

---

## Features

| Category                   | Highlights                                                                                                                                       |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| **Type-safe flags**        | Struct tags (`flag`, `short`, `default`, `help`, `env`, `required`, `count`) — no string lookups                                                 |
| **Per-command flag types** | Each `Command[T, F]` has its own `F` — mix different flag structs freely                                                                         |
| **Dependency injection**   | Built-in [samber/do/v2](https://github.com/samber/do) with `Provide`, `Invoke`, lifecycle hooks                                                  |
| **Environment variables**  | `env:"DB_HOST"` tag with `WithEnvPrefix("MYAPP_")` prefix support                                                                                |
| **16 output formats**      | table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, JSONL, TOML, PlantUML, and more                                                        |
| **Signal handling**        | `WithSignalHandling()` — Ctrl+C cancels context in all handlers                                                                                  |
| **Typo suggestions**       | "did you mean?" for flags and subcommands (Levenshtein distance)                                                                                 |
| **Constructor validation** | Missing handlers, duplicate names, invalid flags — caught at `AddCommand` time                                                                   |
| **Flow context**           | `BranchingFlowContext` — track command path and share values across hierarchy                                                                    |
| **Config files**           | `WithConfigFile(paths...)` — JSON/YAML/TOML auto-loading with flag override                                                                      |
| **Counting flags**         | `count:"true"` for `-v`/`-vv`/`-vvv` verbosity patterns                                                                                          |
| **Extensible types**       | `RegisterTypeHandler()` for custom flag types with full parse/validate support                                                                   |
| **Middleware**             | `TimingMiddleware`, `RecoveryMiddleware`, or write your own; spinner/telemetry/flight-recorder available as [sub-modules](#optional-sub-modules) |
| **Interactive prompts**    | `WithPromptOnMissing()` with `prompt:"Question?"` tag via huh                                                                                    |
| **Markdown help**          | `glamour.WithHelp()` renders Long/Example as styled markdown via glamour                                                                         |
| **Color control**          | `--no-color` flag + `NO_COLOR` env var + `cli.NoColor()` accessor                                                                                |
| **Shell completion**       | Dynamic completion via `WithCompletion(fn)`                                                                                                      |
| **Positional args**        | `WithExactArgs`, `WithMinimumArgs`, `WithRangeArgs`, `WithNoArgs`, or custom                                                                     |
| **Zero panics**            | All functions return errors; no Must\* panic variants                                                                                            |
| **Cobra escape hatch**     | `ConfigFromContext[T]`, `WithPostFlagParse`, `RegisterLocalCommandFlags` — raw cobra + cmdguard runtime                                          |
| **Scoped flags**           | `local:"true"` — root-only flags not inherited by subcommands                                                                                    |
| **Hidden flags**           | `hidden:"true"` — exclude from --help without losing functionality                                                                               |
| **1434 test runs**         | 87.8% coverage, race-detected, fuzz-tested                                                                                                       |

---

## Dependency Injection

Register services on the CLI scope and invoke them in handlers:

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "...", AppConfig{})
scope := cli.Scope()

// Register (lazy initialization)
v4.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{DSN: "postgres://..."}, nil
})

// Invoke in handlers
v4.NewCommand("query", v4.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v4.NoFlags) error {
        db, _ := v4.Invoke[*Database](cli.Scope())
        return db.Query(ctx)
    },
)
```

Services can implement `HealthCheck` and `Shutdown` for lifecycle management.

---

## Environment Variables

```go
type DBFlags struct {
    Host     string `flag:"host"     env:"DB_HOST"     default:"localhost" help:"Database host"`
    Port     int    `flag:"port"     env:"DB_PORT"     default:"5432"      help:"Database port"`
    Password string `flag:"password" env:"DB_PASSWORD"                     help:"Database password"`
}

cli, _ := v4.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v4.WithEnvPrefix("MYAPP_"), // reads MYAPP_DB_HOST, MYAPP_DB_PORT, etc.
)
```

Priority chain: **explicit flag → env var → config file → default value**.

---

## Rich Output

```go
import "github.com/larsartmann/go-output"

v4.OutputTable(output.FormatTable, headers, rows)  // Aligned terminal table
v4.OutputTable(output.FormatJSON, headers, rows)    // JSON array
v4.OutputTable(output.FormatYAML, headers, rows)    // YAML

format, _ := output.ParseFormat("csv")
v4.OutputTable(format, headers, rows)
```

All 16 formats: `table`, `json`, `csv`, `tsv`, `markdown`, `xml`, `yaml`, `html`, `d2`, `tree`, `mermaid`, `dot`, `jsonl`, `asciidoc`, `toml`, `plantuml`.

---

## Subcommands

```go
listCmd, _ := v4.NewCommand("list", v4.NoFlags{}, listHandler,
    v4.WithShort("List users"),
)
createCmd, _ := v4.NewCommand("create", v4.NoFlags{}, createHandler,
    v4.WithShort("Create a user"),
)
userCmd, _ := v4.NewParentCommand[AppConfig]("user",
    "User management", v4.NoFlags{},
    v4.WithSubcommands(listCmd, createCmd),
    v4.WithShort("User management"),
)
v4.AddCommand(cli, userCmd)
```

---

## Lifecycle Hooks

```go
v4.NewCommand("deploy", &Flags{}, runHandler,
    v4.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return validateConfig(flags)
    }),
    v4.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return cleanup()
    }),
)
```

`PostRunE` only fires on success — Cobra semantics. For cleanup that must run
even on failure, use [`WithCleanup[T]`](https://cmdguard.lars.software/guides/lifecycle/)
which fires after every `RunE` regardless of outcome.

---

## Raw Cobra Subcommands (Escape Hatch)

cmdguard's `Command[T,F]` is great for new commands, but real apps often mix
raw `*cobra.Command` subcommands (gradual migration, third-party commands, or
commands that don't fit the typed-flags pattern). cmdguard provides three APIs
to bridge raw cobra commands with the cmdguard runtime:

```go
// 1. Register raw subcommands on cmdguard's root
cli.RootCommand().AddCommand(myRawCmd)

// 2. Access resolved config from any cobra command context
func(cmd *cobra.Command, _ []string) error {
    cfg, ok := v4.ConfigFromContext[AppConfig](cmd.Context())
    if !ok { return errors.New("config not initialized") }
    // use cfg.Field...
}

// 3. Run initialization (DI, logging, session) after flag parsing
cli, _ := v4.NewCLI[AppConfig]("app", "...", AppConfig{},
    v4.WithPostFlagParse[AppConfig](func(cmd *cobra.Command, cfg *AppConfig) error {
        // Flags are parsed, config is resolved, context is stored.
        // Initialize DI, set up logging, store session for subcommands.
        return initDI(cfg)
    }),
)
```

**Scoped flags** (`local:"true"`) prevent root-only flags from polluting every
subcommand's `--help`. Use `cli.RegisterLocalCommandFlags(cmd)` on subcommands
that need the root's execution-flag group.

---

## Built-in Value Types

| Type       | Validation                         |
| ---------- | ---------------------------------- |
| `Duration` | Wraps `time.Duration`              |
| `Enum[T]`  | Validated against allowed values   |
| `LogLevel` | debug / info / warn / error        |
| `URL`      | Validated URL string               |
| `Email`    | RFC 5322 email validation          |
| `Port`     | 1–65535 range                      |
| `FilePath` | Path cleaning and existence checks |
| `HostPort` | `host:port` validation             |

Add your own with `RegisterTypeHandler()`:

```go
v4.RegisterTypeHandler(reflect.TypeFor[MyType](), v4.TypeHandlerFunc{
    ParseFunc:    func(value string, _ v4.FlagTag) (any, error) { return MyType{Value: value}, nil },
    DefaultFunc:  func(_ v4.FlagTag) any { return MyType{} },
})
```

---

## Flag Tags Reference

```go
type Flags struct {
    Name    string `flag:"name"    short:"n" default:"World"  help:"Name"`
    Verbose int    `flag:"verbose" short:"v" help:"Verbosity" count:"true"`
    Host    string `flag:"host"             default:"localhost" env:"DB_HOST" help:"DB host"`
    Mode    string `flag:"mode"  required:"true"                help:"Required!"`
    Build   string `flag:"build" local:"true"  default:"full"   help:"Root-only flag"`
    Debug   string `flag:"debug" hidden:"true"                  help:"Hidden from --help"`
}
```

| Tag        | Purpose                                 | Example                |
| ---------- | --------------------------------------- | ---------------------- |
| `flag`     | Flag name (required)                    | `flag:"name"`          |
| `short`    | Short flag                              | `short:"n"`            |
| `default`  | Default value                           | `default:"World"`      |
| `help`     | Help text                               | `help:"Name to greet"` |
| `env`      | Environment variable                    | `env:"DB_HOST"`        |
| `required` | Mark as required                        | `required:"true"`      |
| `count`    | Counting flag                           | `count:"true"`         |
| `local`    | Root-only, not inherited by subcommands | `local:"true"`         |
| `hidden`   | Exclude from --help but stay functional | `hidden:"true"`        |

---

## Command Options

| Option                                        | Purpose                            |
| --------------------------------------------- | ---------------------------------- |
| `WithShort(short)`                            | Short description                  |
| `WithLong(long)`                              | Long description                   |
| `WithExample(example)`                        | Example usage                      |
| `WithAliases(aliases...)`                     | Alternative names                  |
| _(flags passed positionally to `NewCommand`)_ | —                                  |
| `WithPreRunE[T, F](fn)`                       | Pre-validation hook                |
| `WithPostRunE[T, F](fn)`                      | Post-success cleanup               |
| `WithHidden(bool)`                            | Hide from help                     |
| `WithDeprecated(msg)`                         | Deprecation message                |
| `WithGroupID(id)`                             | Help group name                    |
| `WithExactArgs(n)`                            | Require exactly n positional args  |
| `WithMinimumArgs(n)`                          | Require at least n positional args |
| `WithMaximumArgs(n)`                          | Allow at most n positional args    |
| `WithValidArgs(args...)`                      | Restrict args to allowed values    |
| `WithSubcommands(cmds...)`                    | Attach child commands (parent)     |
| `WithRangeArgs(min, max)`                     | Require between min and max args   |
| `WithNoArgs()`                                | Reject any positional args         |
| `WithCompletion(fn)`                          | Dynamic shell completion           |

---

## CLI Options

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v4.WithCLIVersion("1.0.0"),
    v4.WithEnvPrefix("MYAPP_"),
    v4.WithSignalHandling(),
    v4.WithFang(true),                  // Styled help output
    v4.WithMiddleware[AppConfig](myMiddleware),     // Wrap all handlers
    v4.WithStrictValidation(),           // Require WithShort on commands
    v4.WithConfigValidation[AppConfig](validateFn), // Validate config after parsing
    v4.WithPostFlagParse[AppConfig](initFn),        // DI init / session storage after flags
)
```

| Option                                      | Purpose                                                           |
| ------------------------------------------- | ----------------------------------------------------------------- |
| `WithCLIVersion(v)`                         | Version string                                                    |
| `WithCLILong(desc)`                         | Long description                                                  |
| `WithSilenceErrors()`                       | Suppress error printing (advanced; fang handles this)             |
| `WithSilenceUsage()`                        | Suppress usage on error (**default**; kept for compatibility)     |
| `WithoutSilenceUsage()`                     | Re-enable usage-on-error (opt-in)                                 |
| `WithFang(bool)`                            | Styled help output                                                |
| `WithEnvPrefix(prefix)`                     | Prefix for env vars                                               |
| `WithSignalHandling()`                      | Cancel context on SIGINT/SIGTERM                                  |
| `WithMiddleware[T](mw...)`                  | Middleware for all commands                                       |
| `WithGroup(id, title)`                      | Help group on root                                                |
| `WithConfigValidation[T](fn)`               | Validate config after flag parsing                                |
| `WithPostFlagParse[T](fn...)`               | Post-parse hook: DI init, session storage                         |
| `WithCleanup[T](fn...)`                     | Post-RunE cleanup that fires even when RunE errors                |
| `WithStrictValidation()`                    | Require `WithShort` on all commands                               |
| `WithDraconianValidation()`                 | Strict + require `WithExample` on leaf commands                   |
| `WithConfigFile(paths...)`                  | Auto-load JSON config from first found path                       |
| `WithConfigFileLoader(l, paths...)`         | Load config with custom loader (YAML/TOML)                        |
| `glamour.WithHelp()`                        | Render markdown in command help text (glamour sub-module)         |
| `telemetry.WithTelemetry[T](tracer)`        | OpenTelemetry spans for all commands (telemetry sub-module)       |
| `flightrecorder.WithFlightRecorder[T](cfg)` | Runtime trace snapshots on slow/error (flightrecorder sub-module) |

---

## Error Handling

```go
// All v4 functions return errors — zero panics in library code
cli, err := v4.NewCLI[Config]("app", "...", Config{})
cmd, err := v4.NewCommand("test", NoFlags{}, handler)

// Sentinel errors for errors.Is()
errors.Is(err, v4.ErrInvalidCommand)
errors.Is(err, v4.ErrMissingHandler)
errors.Is(err, v4.ErrDuplicateCommand)

// Rich error types with context
v4.NewCommandError(name, err)
v4.NewFlagError(name, err)
v4.NewFlagErrorWithSuggestion(name, err, suggestion) // includes typo fix
v4.NewExitError(code, err)                            // custom exit code

// ExitCoder interface — check with errors.As
var exitCoder v4.ExitCoder
errors.As(err, &exitCoder)
exitCoder.ExitCode() // returns custom exit code
```

---

## Config Files

`WithConfigFile` auto-detects JSON, YAML, and TOML by file extension. No extra imports needed.

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v4.WithConfigFile("~/.config/myapp/config.yaml", "/etc/myapp/config.json"),
)
```

Paths are tried in order; missing files are silently skipped. Supports `$ENV` and `~` expansion. Nested config structs are fully supported (case-insensitive Go field name matching).

**Precedence:** explicit flag → env var → config file → default value (highest to lowest priority).

---

## Optional Sub-Modules

cmdguard's core stays lean — five optional features live in standalone sub-modules so you import only what you need:

| Sub-module         | Import path                                      | Provides                                                                   |
| ------------------ | ------------------------------------------------ | -------------------------------------------------------------------------- |
| **glamour**        | `github.com/larsartmann/cmdguard/glamour`        | Markdown help rendering (`WithHelp`, `RenderMarkdown`)                     |
| **prompts**        | `github.com/larsartmann/cmdguard/prompts`        | Interactive prompts via huh (`Register`)                                   |
| **spinner**        | `github.com/larsartmann/cmdguard/spinner`        | Terminal spinner middleware (`Middleware`)                                 |
| **telemetry**      | `github.com/larsartmann/cmdguard/telemetry`      | OpenTelemetry spans (`WithTelemetry`, `Middleware`)                        |
| **flightrecorder** | `github.com/larsartmann/cmdguard/flightrecorder` | Runtime trace snapshots on slow/error (`WithFlightRecorder`, `Middleware`) |

```go
import (
    "time"

    "github.com/larsartmann/cmdguard/flightrecorder"
    "github.com/larsartmann/cmdguard/spinner"
    "github.com/larsartmann/cmdguard/telemetry"
)

cli, _ := v4.NewCLI[Config]("app", "...", Config{},
    telemetry.WithTelemetry[Config](tracer),
    v4.WithMiddleware[Config](spinner.Middleware[Config]("Working...")),
    flightrecorder.WithFlightRecorder[Config](flightrecorder.Config{
        CaptureOnSlow:  true,
        SlowThreshold:  200 * time.Millisecond,
        CaptureOnError: true,
    }),
)
```

---

## BranchingFlowContext

Track the command execution path and share values across the hierarchy:

```go
func handler(ctx context.Context, cfg *AppConfig, flags *Flags) error {
    bfc, ok := v4.GetBranchingFlowContext(ctx)
    if ok {
        fmt.Println("Path:", bfc.PathString()) // "myapp.resource.list"
        bfc.SetValue("key", "value")              // propagates to children
        val, _ := bfc.GetValue("key")             // looks up hierarchy
        _ = val
    }
    return nil
}
```

---

## Color Output

cmdguard uses [fang](https://github.com/charmbracelet/fang) for styled help output via [lipgloss](https://github.com/charmbracelet/lipgloss). A `--no-color` flag is registered by default — pass it to disable color output. Lipgloss also respects the [`NO_COLOR`](https://no-color.org/) environment variable automatically.

```go
// Check if color is disabled
if cli.NoColor() {
    // use plain output
}
```

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v4.WithFang(true),   // styled help (default)
    v4.WithFang(false),  // plain text help
)
```

---

## Version Command

```go
cli, _ := v4.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v4.WithCLIVersion("1.0.0"),
)

versionCmd, err := v4.VersionCommand[AppConfig](cli)
if err != nil {
    log.Fatal(err)
}
v4.AddCommand(cli, versionCmd)
// $ myapp version
```

---

## Test Helpers

The `testutil` package provides panic and assertion helpers for testing cmdguard CLIs:

```go
import "github.com/larsartmann/cmdguard/v4/pkg/testutil"

testutil.AssertNoError(t, err)
testutil.AssertErrorIs(t, err, v4.ErrInvalidCommand)
testutil.AssertPanics(t, func() { /* ... */ })
```

See [`examples/taskctl/main_test.go`](examples/taskctl/main_test.go) for integration test patterns using `ExecuteWithArgs`.

---

## Examples

See [`examples/taskctl/`](examples/taskctl/) — a production-grade task manager CLI demonstrating all features: DI, typed flags, middleware, subcommands, config files, rich output, and more.

---

## Development

```bash
# Enter dev shell (Go 1.26, gopls, golangci-lint)
nix develop

# Run tests
go test ./... -count=1 -timeout 120s -race

# Lint
golangci-lint run ./...

# Format (Nix + Go via treefmt)
nix fmt

# Check everything
nix flake check
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for full contribution guidelines.

---

## Documentation

**Full docs at [cmdguard.lars.software](https://cmdguard.lars.software)**

- [Installation](https://cmdguard.lars.software/getting-started/installation/) — Get started in 30 seconds
- [Quick Start](https://cmdguard.lars.software/getting-started/quick-start/) — Build a CLI in under a minute
- [Type-Safe Flags](https://cmdguard.lars.software/guides/type-safe-flags/) — Struct tags reference
- [Custom Value Types](https://cmdguard.lars.software/guides/custom-types/) — Built-in types and RegisterTypeHandler
- [Dependency Injection](https://cmdguard.lars.software/guides/dependency-injection/) — samber/do/v2 integration
- [Lifecycle & Signals](https://cmdguard.lars.software/guides/lifecycle/) — Graceful shutdown, WithCleanup, Doctor
- [Error Handling](https://cmdguard.lars.software/guides/error-handling/) — Zero panics, exit codes, sentinel errors
- [Audit Log](https://cmdguard.lars.software/guides/audit-log/) — DI audit trail in 11 export formats
- [Migrating from Cobra](https://cmdguard.lars.software/guides/migrating-from-cobra/) — Step-by-step guide
- [API Reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4) — Full API on pkg.go.dev

Local docs: [Tutorial](docs/TUTORIAL.md), [Quick Start](docs/QUICKSTART.md), [Framework Comparison](docs/COMPARISON.md), [Performance](docs/PERFORMANCE.md), [CLI Design Principles](docs/CLI_DESIGN_PRINCIPLES.md).

---

## License

[MIT](LICENSE)
