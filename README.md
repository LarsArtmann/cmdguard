# cmdguard

[![CI](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/larsartmann/cmdguard/v3.svg)](https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/cmdguard)](https://goreportcard.com/report/github.com/larsartmann/cmdguard)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Build production Go CLIs with type-safe flags, dependency injection, and zero panics.**

cmdguard wraps [Cobra](https://github.com/spf13/cobra) with compile-time type safety, struct-tag-driven flags, and built-in dependency injection via [samber/do/v2](https://github.com/samber/do). Your flags are typed structs — no more stringly-typed `Flags().GetString("name")` calls that fail at runtime.

> **API Stability:** v3.0.0 is the current major version. The legacy v2 line is in maintenance at v2.10.4. See [CHANGELOG.md](CHANGELOG.md) and the [v2→v3 Migration Guide](docs/MIGRATION_v2_v3.md).

---

## Why cmdguard?

Raw Cobra is powerful but full of footguns that bite every new project. cmdguard
makes the **correct behaviour the default** so consumers use Cobra properly
without having to learn its traps the hard way.

**What cmdguard fixes by default (the parts raw Cobra gets wrong):**

| Raw Cobra footgun                                                               | cmdguard default                                                                           |
| ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------ |
| Prints the full usage block after _every_ command error (`SilenceUsage: false`) | Usage-on-error is silenced; `--help` still works                                           |
| Errors print twice (Cobra prints _and_ `main()` prints the returned error)      | cmdguard prints the error exactly once — see [Error handling](#error-handling--exit-codes) |
| Failed commands exit `0` (easy to forget `os.Exit` with the right code)         | `ExecuteAndExit` / `ExitCode(err)` map errors to correct exit codes                        |
| `Run` (panics) vs `RunE` (returns error) confusion                              | Only error-returning handlers exist — zero panics, by construction                         |
| Missing handler / duplicate command / invalid name found at runtime             | Caught at `NewCommand` / `AddCommand` time                                                 |

**Plus: flags are typed structs, validated at construction — no stringly-typed lookups:**

**Raw Cobra — flags are strings, validated at runtime:**

```go
var name string
var count int
rootCmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")
rootCmd.Flags().IntVarP(&count, "count", "c", 1, "Number of greetings")
// Oops — forgot to add "count"? You find out at runtime.
```

**cmdguard — flags are typed structs, validated at construction:**

```go
type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Count int    `flag:"count" short:"c" default:"1"    help:"Number of greetings"`
}
// Missing handler? Duplicate command? Invalid name? Caught at AddCommand time.
```

---

## Error handling & exit codes

cmdguard owns the error-output contract so you can't get it wrong:

- **The error is printed exactly once** — styled by [fang](https://github.com/charmbracelet/fang) when enabled (the default), or plain by Cobra when disabled.
- **Usage is never printed on error** (`SilenceUsage: true` by default).
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
os.Exit(v3.ExitCode(err)) // 0 on success, ExitCoder code or 1 on failure
```

> Pitfall to avoid: `if err := cli.Execute(ctx); err != nil { fmt.Fprintln(os.Stderr, err) }`
> re-prints the error that cmdguard already printed.

---

## Quick Start

```bash
go get github.com/larsartmann/cmdguard/v3
```

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
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
    cli, err := v3.NewCLI[AppConfig]("myapp", "My CLI application", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
        os.Exit(1)
    }

    greetCmd, err := v3.NewCommand("greet", &GreetFlags{},
        func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
        v3.WithShort("Greet someone"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create command: %v\n", err)
        os.Exit(1)
    }

    v3.AddCommand(cli, greetCmd)
    cli.ExecuteAndExit(context.Background())
}
```

```bash
$ go run main.go greet -n "cmdguard" --shout
HELLO, CMDGUARD!
```

---

## Features

| Category                   | Highlights                                                                                              |
| -------------------------- | ------------------------------------------------------------------------------------------------------- |
| **Type-safe flags**        | Struct tags (`flag`, `short`, `default`, `help`, `env`, `required`, `count`) — no string lookups        |
| **Per-command flag types** | Each `Command[T, F]` has its own `F` — mix different flag structs freely                                |
| **Dependency injection**   | Built-in [samber/do/v2](https://github.com/samber/do) with `Provide`, `Invoke`, lifecycle hooks         |
| **Environment variables**  | `env:"DB_HOST"` tag with `WithEnvPrefix("MYAPP_")` prefix support                                       |
| **16 output formats**      | table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, JSONL, TOML, PlantUML, and more               |
| **Signal handling**        | `WithSignalHandling()` — Ctrl+C cancels context in all handlers                                         |
| **Typo suggestions**       | "did you mean?" for flags and subcommands (Levenshtein distance)                                        |
| **Constructor validation** | Missing handlers, duplicate names, invalid flags — caught at `AddCommand` time                          |
| **Flow context**           | `BranchingFlowContext` — track command path and share values across hierarchy                           |
| **Config files**           | `WithConfigFile(paths...)` — JSON/YAML/TOML auto-loading with flag override                             |
| **Counting flags**         | `count:"true"` for `-v`/`-vv`/`-vvv` verbosity patterns                                                 |
| **Extensible types**       | `RegisterTypeHandler()` for custom flag types with full parse/validate support                          |
| **Middleware**             | `TimingMiddleware`, `RecoveryMiddleware`, or write your own; spinner/telemetry available as [sub-modules](#optional-sub-modules) |
| **Interactive prompts**    | `WithPromptOnMissing()` with `prompt:"Question?"` tag via huh                                          |
| **Markdown help**          | `glamour.WithHelp()` renders Long/Example as styled markdown via glamour                                |
| **Color control**          | `--no-color` flag + `NO_COLOR` env var + `cli.NoColor()` accessor                                       |
| **Shell completion**       | Dynamic completion via `WithCompletion(fn)`                                                             |
| **Man page generation**    | `manpage.GenerateCommand[T](cli)` for roff output                                                       |
| **Positional args**        | `WithExactArgs`, `WithMinimumArgs`, `WithRangeArgs`, `WithNoArgs`, or custom                            |
| **Zero panics**            | All functions return errors; no Must\* panic variants                                                   |
| **Cobra escape hatch**     | `ConfigFromContext[T]`, `WithPostFlagParse`, `RegisterLocalCommandFlags` — raw cobra + cmdguard runtime |
| **Scoped flags**           | `local:"true"` — root-only flags not inherited by subcommands                                           |
| **Hidden flags**           | `hidden:"true"` — exclude from --help without losing functionality                                      |
| **420+ tests**             | 87.3% coverage, race-detected, fuzz-tested                                                              |

---

## Dependency Injection

Register services on the CLI scope and invoke them in handlers:

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{})
scope := cli.Scope()

// Register (lazy initialization)
v3.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{DSN: "postgres://..."}, nil
})

// Invoke in handlers
v3.NewCommand("query", v3.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
        db, _ := v3.Invoke[*Database](cli.Scope())
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

cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v3.WithEnvPrefix("MYAPP_"), // reads MYAPP_DB_HOST, MYAPP_DB_PORT, etc.
)
```

Priority chain: **explicit flag → env var → config file → default value**.

---

## Rich Output

```go
import "github.com/larsartmann/go-output"

v3.OutputTable(output.FormatTable, headers, rows)  // Aligned terminal table
v3.OutputTable(output.FormatJSON, headers, rows)    // JSON array
v3.OutputTable(output.FormatYAML, headers, rows)    // YAML

format, _ := output.ParseFormat("csv")
v3.OutputTable(format, headers, rows)
```

All 16 formats: `table`, `json`, `csv`, `tsv`, `markdown`, `xml`, `yaml`, `html`, `d2`, `tree`, `mermaid`, `dot`, `jsonl`, `asciidoc`, `toml`, `plantuml`.

---

## Subcommands

```go
listCmd, _ := v3.NewCommand("list", v3.NoFlags{}, listHandler,
    v3.WithShort("List users"),
)
createCmd, _ := v3.NewCommand("create", v3.NoFlags{}, createHandler,
    v3.WithShort("Create a user"),
)
userCmd, _ := v3.NewParentCommand[AppConfig]("user",
    "User management", v3.NoFlags{},
    v3.WithSubcommands(listCmd, createCmd),
    v3.WithShort("User management"),
)
v3.AddCommand(cli, userCmd)
```

---

## Lifecycle Hooks

```go
v3.NewCommand("deploy", &Flags{}, runHandler,
    v3.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return validateConfig(flags)
    }),
    v3.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return cleanup()
    }),
)
```

`PostRunE` only fires on success — Cobra semantics. For cleanup that must run
even on failure, put `defer` directly inside your `RunE` handler.

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
    cfg, ok := v3.ConfigFromContext[AppConfig](cmd.Context())
    if !ok { return errors.New("config not initialized") }
    // use cfg.Field...
}

// 3. Run initialization (DI, logging, session) after flag parsing
cli, _ := v3.NewCLI[AppConfig]("app", "...", AppConfig{},
    v3.WithPostFlagParse[AppConfig](func(cmd *cobra.Command, cfg *AppConfig) error {
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
v3.RegisterTypeHandler(reflect.TypeFor[MyType](), v3.TypeHandlerFunc{
    ParseFunc:    func(value string, _ v3.FlagTag) (any, error) { return MyType{Value: value}, nil },
    DefaultFunc:  func(_ v3.FlagTag) any { return MyType{} },
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
cli, _ := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v3.WithCLIVersion("1.0.0"),
    v3.WithEnvPrefix("MYAPP_"),
    v3.WithSignalHandling(),
    v3.WithFang(true),                  // Styled help output
    v3.WithMiddleware[AppConfig](myMiddleware),     // Wrap all handlers
    v3.WithStrictValidation(),           // Require WithShort on commands
    v3.WithConfigValidation[AppConfig](validateFn), // Validate config after parsing
    v3.WithPostFlagParse[AppConfig](initFn),        // DI init / session storage after flags
)
```

| Option                              | Purpose                                                       |
| ----------------------------------- | ------------------------------------------------------------- |
| `WithCLIVersion(v)`                 | Version string                                                |
| `WithCLILong(desc)`                 | Long description                                              |
| `WithSilenceErrors()`               | Suppress error printing (advanced; fang handles this)         |
| `WithSilenceUsage()`                | Suppress usage on error (**default**; kept for compatibility) |
| `WithFang(bool)`                    | Styled help output                                            |
| `WithEnvPrefix(prefix)`             | Prefix for env vars                                           |
| `WithSignalHandling()`              | Cancel context on SIGINT/SIGTERM                              |
| `WithMiddleware[T](mw...)`          | Middleware for all commands                                   |
| `WithGroup(id, title)`              | Help group on root                                            |
| `WithConfigValidation[T](fn)`       | Validate config after flag parsing                            |
| `WithPostFlagParse[T](fn...)`       | Post-parse hook: DI init, session storage                     |
| `WithCleanup[T](fn...)`             | Post-RunE cleanup that fires even when RunE errors            |
| `WithStrictValidation()`            | Require `WithShort` on all commands                           |
| `WithDraconianValidation()`         | Strict + require `WithExample` on leaf commands               |
| `WithConfigFile(paths...)`          | Auto-load JSON config from first found path                   |
| `WithConfigFileLoader(l, paths...)` | Load config with custom loader (YAML/TOML)                    |
| `glamour.WithHelp()`                | Render markdown in command help text (glamour sub-module)     |
| `telemetry.WithTelemetry[T](tracer)` | OpenTelemetry spans for all commands (telemetry sub-module)   |

---

## Error Handling

```go
// All v3 functions return errors — zero panics in library code
cli, err := v3.NewCLI[Config]("app", "...", Config{})
cmd, err := v3.NewCommand("test", NoFlags{}, handler)

// Sentinel errors for errors.Is()
errors.Is(err, v3.ErrInvalidCommand)
errors.Is(err, v3.ErrMissingHandler)
errors.Is(err, v3.ErrDuplicateCommand)

// Rich error types with context
v3.NewCommandError(name, err)
v3.NewFlagError(name, err)
v3.NewFlagErrorWithSuggestion(name, err, suggestion) // includes typo fix
v3.NewExitError(code, err)                            // custom exit code

// ExitCoder interface — check with errors.As
var exitCoder v3.ExitCoder
errors.As(err, &exitCoder)
exitCoder.ExitCode() // returns custom exit code
```

---

## Config Files

### JSON (built-in)

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v3.WithConfigFile("~/.config/myapp/config.json", "/etc/myapp/config.json"),
)
```

Paths are tried in order; missing files are silently skipped. Supports `$ENV` and `~` expansion.

### YAML / TOML (custom loaders)

```go
import "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3/configload"

cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v3.WithConfigFileLoader(configload.YAML(), "config.yaml"),
)
```

`configload.YAML()` and `configload.TOML()` return `ConfigFileLoader` implementations. See [`pkg/cmdguard/v3/configload/`](pkg/cmdguard/v3/configload/) for available loaders.

**Precedence:** explicit flag → env var → config file → default value (highest to lowest priority).

---

## Man Page Generation

```go
import "github.com/larsartmann/cmdguard/manpage"

manCmd, err := manpage.GenerateCommand[AppConfig](cli)
if err != nil {
    log.Fatal(err)
}
v3.AddCommand(cli, manCmd)
// $ myapp man
```

Generates roff-formatted man pages from your command structure.

---

## Optional Sub-Modules

cmdguard's core is dependency-free. Five optional features live in standalone sub-modules — import only what you need:

| Sub-module                            | Import path                                  | Provides                                            |
| ------------------------------------- | -------------------------------------------- | --------------------------------------------------- |
| **glamour**                           | `github.com/larsartmann/cmdguard/glamour`    | Markdown help rendering (`WithHelp`, `RenderMarkdown`) |
| **manpage**                           | `github.com/larsartmann/cmdguard/manpage`    | Man page generation (`GenerateCommand`, `Write`)    |
| **prompts**                           | `github.com/larsartmann/cmdguard/prompts`    | Interactive prompts via huh (`Register`)            |
| **spinner**                           | `github.com/larsartmann/cmdguard/spinner`    | Terminal spinner middleware (`Middleware`)          |
| **telemetry**                         | `github.com/larsartmann/cmdguard/telemetry`  | OpenTelemetry spans (`WithTelemetry`, `Middleware`) |

```go
import (
    "github.com/larsartmann/cmdguard/spinner"
    "github.com/larsartmann/cmdguard/telemetry"
)

cli, _ := v3.NewCLI[Config]("app", "...", Config{},
    telemetry.WithTelemetry[Config](tracer),
    v3.WithMiddleware[Config](spinner.Middleware[Config]("Working...")),
)
```

---

## BranchingFlowContext

Track the command execution path and share values across the hierarchy:

```go
func handler(ctx context.Context, cfg *AppConfig, flags *Flags) error {
    bfc, ok := v3.GetBranchingFlowContext(ctx)
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
cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v3.WithFang(true),   // styled help (default)
    v3.WithFang(false),  // plain text help
)
```

---

## Version Command

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v3.WithCLIVersion("1.0.0"),
)

versionCmd, err := v3.VersionCommand[AppConfig](cli)
if err != nil {
    log.Fatal(err)
}
v3.AddCommand(cli, versionCmd)
// $ myapp version
```

---

## Test Helpers

The `testutil` subpackage provides a harness for testing cmdguard CLIs:

```go
import "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3/testutil"

result := testutil.TestCLI(t, cli, []string{"greet", "--name", "Alice"})
result.AssertNoError()
result.AssertExitCode(0)
result.AssertOutputContains("Hello, Alice!")
```

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

- [Tutorial](docs/TUTORIAL.md) — Build a task manager CLI step by step
- [Quick Start Guide](docs/QUICKSTART.md) — Learn cmdguard in 5 minutes
- [Migrating from Cobra](docs/MIGRATION_FROM_COBRA.md) — Step-by-step migration guide
- [Framework Comparison](docs/COMPARISON.md) — vs Kong, sflags, go-flags, urfave/cli
- [Performance](docs/PERFORMANCE.md) — Benchmark results and overhead analysis
- [CLI Design Principles](docs/CLI_DESIGN_PRINCIPLES.md) — Design guidelines
- [API Reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3) — Full API docs on pkg.go.dev

---

## License

[MIT](LICENSE)
