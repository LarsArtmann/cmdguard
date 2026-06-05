# cmdguard

[![CI](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml/badge.svg)](https://github.com/larsartmann/cmdguard/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/larsartmann/cmdguard/v2.svg)](https://pkg.go.dev/github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2)
[![Go Report Card](https://goreportcard.com/badge/github.com/larsartmann/cmdguard)](https://goreportcard.com/report/github.com/larsartmann/cmdguard)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Build production Go CLIs with type-safe flags, dependency injection, and zero panics.**

cmdguard wraps [Cobra](https://github.com/spf13/cobra) with compile-time type safety, struct-tag-driven flags, and built-in dependency injection via [samber/do/v2](https://github.com/samber/do). Your flags are typed structs — no more stringly-typed `Flags().GetString("name")` calls that fail at runtime.

> **API Stability:** The v2 API is stable and will only receive additive changes until v3. See [CHANGELOG.md](CHANGELOG.md) for deprecation policy.

---

## Why cmdguard?

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

## Quick Start

```bash
go get github.com/larsartmann/cmdguard/v2
```

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strings"

    "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
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
    cli, err := v2.NewCLI[AppConfig]("myapp", "My CLI application", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create CLI: %v\n", err)
        os.Exit(1)
    }

    greetCmd, err := v2.NewCommand[AppConfig, *GreetFlags]("greet",
        func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
            msg := fmt.Sprintf("Hello, %s!", flags.Name)
            if flags.Shout {
                msg = strings.ToUpper(msg)
            }
            fmt.Println(msg)
            return nil
        },
        v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
        v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create command: %v\n", err)
        os.Exit(1)
    }

    v2.AddCommand(cli, greetCmd)
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
| -------------------------- | ------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| **Type-safe flags**        | Struct tags (`flag`, `short`, `default`, `help`, `env`, `required`, `count`) — no string lookups        |
| **Per-command flag types** | Each `Command[T, F]` has its own `F` — mix different flag structs freely                                |
| **Dependency injection**   | Built-in [samber/do/v2](https://github.com/samber/do) with `Provide`, `Invoke`, lifecycle hooks         |
| **Environment variables**  | `env:"DB_HOST"` tag with `WithEnvPrefix("MYAPP_")` prefix support                                       |
| **12 output formats**      | table, JSON, CSV, YAML, Markdown, XML, HTML, D2, Mermaid, and more                                      |
| **Signal handling**        | `WithSignalHandling[T]()` — Ctrl+C cancels context in all handlers                                      |
| **Typo suggestions**       | "did you mean?" for flags and subcommands (Levenshtein distance)                                        |
| **Constructor validation** | Missing handlers, duplicate names, invalid flags — caught at `AddCommand` time                          |
| **Flow context**           | `BranchingFlowContext` — track command path and share values across hierarchy                           |
| **Editor support**         | `EditInEditor()` — open `$EDITOR` for user input                                                        |
| **Config files**           | `WithConfigFile[T](paths...)` — JSON/YAML/TOML auto-loading with flag override                          |
| **Counting flags**         | `count:"true"` for `-v`/`-vv`/`-vvv` verbosity patterns                                                 |
| **Extensible types**       | `RegisterTypeHandler()` for custom flag types with full parse/validate support                          |
| **Middleware**             | `TimingMiddleware`, `RecoveryMiddleware`, `SpinnerMiddleware`, `TelemetryMiddleware`, or write your own |
|                            | **Interactive prompts**                                                                                 | `WithPromptOnMissing[T,F]()` with `prompt:"Question?"` tag via huh         |
|                            | **Markdown help**                                                                                       | `WithGlamourHelp[T]()` renders Long/Example as styled markdown via glamour |
|                            | **Color control**                                                                                       | `--no-color` flag + `NO_COLOR` env var + `cli.NoColor()` accessor          |
| **Shell completion**       | Dynamic completion via `WithCompletion[T, F](fn)`                                                       |
| **Man page generation**    | `GenerateManPageCommand[T](cli)` for roff output                                                        |
| **Positional args**        | `WithExactArgs`, `WithMinimumArgs`, `WithRangeArgs`, `WithNoArgs`, or custom                            |
| **Zero panics**            | Every v2 API function returns errors — never panics in library code                                     |
| **356 tests** (706 cases)  | 82.8% coverage, race-detected, fuzz-tested                                                              |

---

## Dependency Injection

Register services on the CLI scope and invoke them in handlers:

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{})
scope := cli.Scope()

// Register (lazy initialization)
v2.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{DSN: "postgres://..."}, nil
})

// Invoke in handlers
v2.NewCommand[AppConfig, v2.NoFlags]("query",
    func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
        db, _ := v2.Invoke[*Database](cli.Scope())
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

cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v2.WithEnvPrefix[AppConfig]("MYAPP_"), // reads MYAPP_DB_HOST, MYAPP_DB_PORT, etc.
)
```

Priority chain: **explicit flag → env var → default value**.

---

## Rich Output

```go
v2.OutputTable(v2.FormatTable, headers, rows)  // Aligned terminal table
v2.OutputTable(v2.FormatJSON, headers, rows)    // JSON array
v2.OutputTable(v2.FormatYAML, headers, rows)    // YAML

format, _ := v2.ParseOutputFormat("csv")
v2.OutputTable(format, headers, rows)
```

All 12 formats: `table`, `json`, `csv`, `tsv`, `markdown`, `xml`, `yaml`, `html`, `d2`, `tree`, `mermaid`, `dot`.

---

## Subcommands

```go
listCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("list", listHandler,
    v2.WithShort[AppConfig, v2.NoFlags]("List users"),
)
createCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("create", createHandler,
    v2.WithShort[AppConfig, v2.NoFlags]("Create a user"),
)
userCmd, _ := v2.NewParentCommand[AppConfig, v2.NoFlags]("user",
    "User management", []v2.Command[AppConfig, v2.NoFlags]{listCmd, createCmd},
    v2.WithShort[AppConfig, v2.NoFlags]("User management"),
)
v2.AddCommand(cli, userCmd)
```

---

## Lifecycle Hooks

```go
v2.NewCommand[AppConfig, *Flags]("deploy", runHandler,
    v2.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return validateConfig(flags)
    }),
    v2.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return cleanup()
    }),
)
```

`PostRunE` only fires on success — Cobra semantics.

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
v2.RegisterTypeHandler(reflect.TypeFor[MyType](), v2.TypeHandlerFunc{
    ParseFunc:    func(value string, _ v2.FlagTag) (any, error) { return MyType{Value: value}, nil },
    DefaultFunc:  func(_ v2.FlagTag) any { return MyType{} },
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
}
```

| Tag        | Purpose              | Example                |
| ---------- | -------------------- | ---------------------- |
| `flag`     | Flag name (required) | `flag:"name"`          |
| `short`    | Short flag           | `short:"n"`            |
| `default`  | Default value        | `default:"World"`      |
| `help`     | Help text            | `help:"Name to greet"` |
| `env`      | Environment variable | `env:"DB_HOST"`        |
| `required` | Mark as required     | `required:"true"`      |
| `count`    | Counting flag        | `count:"true"`         |

---

## Command Options

| Option                          | Purpose                            |
| ------------------------------- | ---------------------------------- | ------------------------------- |
| `WithShort[T, F](short)`        | Short description                  |
| `WithLong[T, F](long)`          | Long description                   |
| `WithExample[T, F](example)`    | Example usage                      |
| `WithAliases[T, F](aliases...)` | Alternative names                  |
| `WithFlags[T, F](flags)`        | Typed flags struct                 |
| `WithPreRunE[T, F](fn)`         | Pre-validation hook                |
| `WithPostRunE[T, F](fn)`        | Post-success cleanup               |
| `WithHidden[T, F](bool)`        | Hide from help                     |
| `WithDeprecated[T, F](msg)`     | Deprecation message                |
| `WithGroupID[T, F](id)`         | Help group name                    |
| `WithExactArgs[T, F](n)`        | Require exactly n positional args  |
| `WithMinimumArgs[T, F](n)`      | Require at least n positional args |
| `WithMaximumArgs[T, F](n)`      | Allow at most n positional args    |
|                                 | `WithValidArgs[T, F](args...)`     | Restrict args to allowed values |
|                                 | `WithSubcommands[T, F](cmds...)`   | Attach child commands (parent)  |
| `WithRangeArgs[T, F](min, max)` | Require between min and max args   |
| `WithNoArgs[T, F]()`            | Reject any positional args         |
| `WithCompletion[T, F](fn)`      | Dynamic shell completion           |

---

## CLI Options

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v2.WithCLIVersion[AppConfig]("1.0.0"),
    v2.WithEnvPrefix[AppConfig]("MYAPP_"),
    v2.WithSignalHandling[AppConfig](),
    v2.WithFang[AppConfig](true),                  // Styled help output
    v2.WithMiddleware[AppConfig](myMiddleware),     // Wrap all handlers
    v2.WithStrictValidation[AppConfig](),           // Require WithShort on commands
    v2.WithConfigValidation[AppConfig](validateFn), // Validate config after parsing
)
```

| Option                                 | Purpose                                         |
| -------------------------------------- | ----------------------------------------------- |
| `WithCLIVersion[T](v)`                 | Version string                                  |
| `WithCLILong[T](desc)`                 | Long description                                |
| `WithSilenceErrors[T]()`               | Suppress error printing                         |
| `WithSilenceUsage[T]()`                | Suppress usage on error                         |
| `WithFang[T](bool)`                    | Styled help output                              |
| `WithEnvPrefix[T](prefix)`             | Prefix for env vars                             |
| `WithSignalHandling[T]()`              | Cancel context on SIGINT/SIGTERM                |
| `WithMiddleware[T](mw...)`             | Middleware for all commands                     |
| `WithGroup[T](id, title)`              | Help group on root                              |
| `WithConfigValidation[T](fn)`          | Validate config after flag parsing              |
| `WithStrictValidation[T]()`            | Require `WithShort` on all commands             |
| `WithDraconianValidation[T]()`         | Strict + require `WithExample` on leaf commands |
| `WithConfigFile[T](paths...)`          | Auto-load JSON config from first found path     |
| `WithConfigFileLoader[T](l, paths...)` | Load config with custom loader (YAML/TOML)      |
| `WithGlamourHelp[T]()`                 | Render markdown in command help text            |
| `WithTelemetry[T](tracer)`             | OpenTelemetry spans for all commands            |

---

## Error Handling

```go
// All v2 functions return errors — no panics
cli, err := v2.NewCLI[Config]("app", "...", Config{})
cmd, err := v2.NewCommand[Config, NoFlags]("test", handler)

// Sentinel errors for errors.Is()
errors.Is(err, v2.ErrInvalidCommand)
errors.Is(err, v2.ErrMissingHandler)
errors.Is(err, v2.ErrDuplicateCommand)

// Rich error types with context
v2.NewCommandError(name, err)
v2.NewFlagError(name, err)
v2.NewFlagErrorWithSuggestion(name, err, suggestion) // includes typo fix
v2.NewExitError(code, err)                            // custom exit code

// ExitCoder interface — check with errors.As
var exitCoder v2.ExitCoder
errors.As(err, &exitCoder)
exitCoder.ExitCode() // returns custom exit code
```

---

### Must Constructors

`MustNewCommand` and `MustNewParentCommand` panic on error — use when configuration is known at compile time:

```go
greetCmd := v2.MustNewCommand[AppConfig, *GreetFlags]("greet", greetHandler,
    v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)

parentCmd := v2.MustNewParentCommand[AppConfig, v2.NoFlags]("user",
    "User management", []v2.Command[AppConfig, v2.NoFlags]{listCmd, createCmd},
    v2.WithShort[AppConfig, v2.NoFlags]("User management"),
)
```

---

## Config Files

### JSON (built-in)

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v2.WithConfigFile[AppConfig]("~/.config/myapp/config.json", "/etc/myapp/config.json"),
)
```

Paths are tried in order; missing files are silently skipped. Supports `$ENV` and `~` expansion.

### YAML / TOML (custom loaders)

```go
import "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2/configload"

cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v2.WithConfigFileLoader[AppConfig](configload.YAML(), "config.yaml"),
)
```

`configload.YAML()` and `configload.TOML()` return `ConfigFileLoader` implementations. See [`pkg/cmdguard/v2/configload/`](pkg/cmdguard/v2/configload/) for available loaders.

**Precedence:** config file → environment variables → explicit flags → defaults.

---

## Man Page Generation

```go
manCmd, err := v2.GenerateManPageCommand[AppConfig](cli)
if err != nil {
    log.Fatal(err)
}
v2.AddCommand(cli, manCmd)
// $ myapp man
```

Generates roff-formatted man pages from your command structure.

---

## BranchingFlowContext

Track the command execution path and share values across the hierarchy:

```go
func handler(ctx context.Context, cfg *AppConfig, flags *Flags) error {
    bfc, ok := v2.GetBranchingFlowContext(ctx)
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
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v2.WithFang[AppConfig](true),   // styled help (default)
    v2.WithFang[AppConfig](false),  // plain text help
)
```

---

## EditInEditor

Open the user's `$EDITOR` to edit content interactively:

```go
edited, err := v2.EditInEditor(ctx, "# Edit your message here\n")
if err != nil {
    return err
}
fmt.Println("User wrote:", edited)
```

---

## Version Command

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "...", AppConfig{},
    v2.WithCLIVersion[AppConfig]("1.0.0"),
)

versionCmd := v2.MustVersionCommand[AppConfig](cli)
v2.AddCommand(cli, versionCmd)
// $ myapp version
```

---

## Test Helpers

The `testutil` subpackage provides a harness for testing cmdguard CLIs:

```go
import "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2/testutil"

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
- [API Reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2) — Full API docs on pkg.go.dev

---

## License

[MIT](LICENSE)
