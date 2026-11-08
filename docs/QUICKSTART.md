# Quickstart Guide

**Learn cmdguard in 5 minutes.**

cmdguard is a Go library for building validated CLI applications with type-safe flags, dependency injection, and 16 output formats.

---

## Installation

```bash
go get github.com/larsartmann/cmdguard/v5
```

---

## Your First CLI (2 minutes)

### 1. Create main.go

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type AppConfig struct{}

func main() {
    cli, err := v5.NewCLI[AppConfig]("hello", "A simple CLI", AppConfig{})
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cmd, err := v5.NewCommand("greet", v5.NoFlags{},
        func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
            fmt.Println("Hello, World!")
            return nil
        },
        v5.WithShort("Greet someone"),
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    if err := v5.AddCommand(cli, cmd); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    cli.ExecuteAndExit(context.Background())
}
```

### 2. Run It

```bash
go run main.go greet
# Output: Hello, World!

go run main.go --help
# Shows styled usage help
```

---

## Adding Flags (1 minute)

### 1. Define Flag Struct

```go
type GreetFlags struct {
    Name  string `flag:"name"  short:"n" default:"World" help:"Name to greet"`
    Shout bool   `flag:"shout" short:"s" default:"false" help:"Shout the greeting"`
}
```

### 2. Use in Command

```go
cmd, err := v5.NewCommand("greet", &GreetFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *GreetFlags) error {
        msg := fmt.Sprintf("Hello, %s!", flags.Name)
        if flags.Shout {
            msg = strings.ToUpper(msg)
        }
        fmt.Println(msg)
        return nil
    },
    v5.WithShort("Greet someone"),
)
v5.AddCommand(cli, cmd)
```

### 3. Run It

```bash
go run main.go greet --name=Alice
# Output: Hello, Alice!

go run main.go greet -n Bob -s
# Output: HELLO, BOB!

go run main.go greet --help
# Shows flag help
```

---

## Subcommands (1 minute)

```go
migrateCmd, _ := v5.NewCommand("migrate", v5.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
        fmt.Println("Running migrations...")
        return nil
    },
    v5.WithShort("Run migrations"),
)

rollbackCmd, _ := v5.NewCommand("rollback", v5.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
        fmt.Println("Rolling back...")
        return nil
    },
    v5.WithShort("Rollback last migration"),
)

dbCmd, _ := v5.NewParentCommand[AppConfig]("db", "Database operations", v5.NoFlags{},
    v5.WithSubcommands(migrateCmd, rollbackCmd),
    v5.WithShort("Database operations"),
)
v5.AddCommand(cli, dbCmd)
```

```bash
go run main.go db migrate
go run main.go db rollback
```

---

## Environment Variables

Flags can read from environment variables with the `env` tag:

```go
type DBFlags struct {
    Host     string `flag:"host"     env:"DB_HOST"     default:"localhost" help:"Database host"`
    Port     int    `flag:"port"     env:"DB_PORT"     default:"5432"      help:"Database port"`
    Password string `flag:"password" env:"DB_PASSWORD" default:""          help:"Database password"`
}
```

Priority: explicit flag > env var > default.

Add a prefix with `WithEnvPrefix`:

```go
cli, _ := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v5.WithEnvPrefix("MYAPP_"),
)
// Now reads MYAPP_DB_HOST, MYAPP_DB_PORT, etc.
```

---

## Counting Flags

Use `count:"true"` for `-v`/`-vv`/`-vvv` verbosity:

```go
type Flags struct {
    Verbose int `flag:"verbose" short:"v" help:"Verbosity level" count:"true"`
}
```

```bash
myapp run          # Verbose = 0
myapp run -v       # Verbose = 1
myapp run -vvv     # Verbose = 3
```

---

## Dependency Injection (1 minute)

### 1. Register Services

```go
cli, _ := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{})

v5.ProvideValue(cli.Scope(), &Logger{Level: "info"})
```

### 2. Use in Handler

```go
cmd, _ := v5.NewCommand("log", v5.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
        logger, err := v5.Invoke[*Logger](cli.Scope())
        if err != nil {
            return err
        }
        logger.Log("Hello from DI!")
        return nil
    },
    v5.WithShort("Log a message"),
)
v5.AddCommand(cli, cmd)
```

---

## Signal Handling and Graceful Shutdown

On by default in v5: SIGINT/SIGTERM cancel the context in handlers, and DI
services implementing `do.ShutdownerWithError` shut down in reverse creation
order. Ctrl+C just works:

```go
cli, _ := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
// Ctrl+C cancels the context in handlers and shuts down DI services
```

Opt out with `WithoutSignalHandling()` (drops both) or
`WithoutGracefulShutdown()` (keeps cancellation, skips DI teardown). A
normally completing run never tears down the scope, so repeated
`ExecuteWithArgs` is safe.

---

## Batteries: The Production Defaults

v5 ships opinionated defaults — no configuration needed for a serious CLI:

| Default | What you get | Opt-out |
| ------- | ------------ | ------- |
| Strict validation | every command needs `WithShort` | `WithoutStrictValidation()` |
| Signals + graceful shutdown | Ctrl+C handled, DI shut down | `WithoutSignalHandling()` |
| Strict config keys | config typos fail loudly | `WithoutStrictConfigKeys()` |
| Error hints | a `Hint:` line under every known error | `WithoutErrorHints()` |
| `doctor` command | DI health + config diagnostics | `WithoutBatteries()` |
| `version` command | when `WithCLIVersion` is set | consumer command wins |

```go
cli, _ := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v5.WithCLIVersion("1.0.0"),
)
// app doctor   → health report
// app version  → version output
```

---

## The describe Contract and Error Hints

Every v5 CLI exposes a machine-readable contract of its full command tree:

```bash
myapp describe   # JSON: commands, flags, env vars, defaults, required, hints
```

Agents and scripts use it to discover the CLI without reading help text. A
consumer-defined `describe` command wins; `WithoutDescribeCommand()` opts
out (`CLI.Describe()` always works).

Known errors carry one actionable hint, printed under the message and
emitted as `hint` in the JSON error envelope:

```
error: required flag not set
Hint: pass the flag on the command line, or set its env var or config-file value
```

---

## Runtime Tag Enforcement

`required`, `values`, and `validate` tags are enforced after parsing — at
run time, not just construction:

```go
type DeployFlags struct {
	Env  string `flag:"env" required:"true"        help:"Target environment"`
	Mode string `flag:"mode" values:"debug,release" help:"Build mode"`
	Age  int    `flag:"age"  min:"1" max:"120"        help:"Age"`
}
```

- `required` — satisfied by the flag OR its env var; missing fails with a hint
- `values` — the value must be one of the listed options
- `min`/`max` — numeric bounds (length bounds on strings)

Invalid defaults fail at construction; invalid env/config values fail before
the handler runs. Introspection commands (`describe`, `version`, `doctor`,
`help`) skip required enforcement so discovery always works.

---

## Rich Output (16 Formats)

```go
import "github.com/larsartmann/go-output"

// Use go-output directly for format parsing
format := output.FormatJSON

// Render table data
headers := []string{"Name", "Age"}
rows := [][]string{{"Alice", "30"}, {"Bob", "25"}}
v5.OutputTable(format, headers, rows)
```

Supported formats: table, json, csv, tsv, markdown, xml, d2, yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml.

---

## Config Files and Strict Keys

```go
cli, err := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v5.WithConfigFile("myapp.yaml", "$HOME/.config/myapp.yaml"),
)
```

Every field in your config struct accepts several spellings in the file: the
field name (`LogLevel`), the flag name (`log-level`), and case/hyphen/underscore
variants (`loglevel`, `log_level`). Nested keys are dotted (`database.host`);
arrays are indexed (`servers[0].port`).

Unknown keys **fail the load by default** (strict config keys) — a typo is
caught at startup, not silently ignored. To accept unknown keys (a config
file shared with other tools, for example):

```go
v5.WithoutStrictConfigKeys(),
```

A typo now fails with the key, the file path, and a suggestion:

```
error: ...: failed to load config file: loading config file: "/home/me/.config/myapp.yaml": unknown config key "log-leve1" (did you mean "LogLevel"?)
```

`app doctor` reports config problems (missing file, parse failure, unknown
keys) alongside DI health.

---

## Key Concepts

### Constructor Pattern

All commands are created via constructors, not struct literals:

```go
// Leaf command — flags passed positionally, options are non-generic
cmd, err := v5.NewCommand("name", &Flags{}, handler, opts...)

// Parent command — subcommands added via WithSubcommands option
parent, err := v5.NewParentCommand[T]("name", "long desc", v5.NoFlags{},
    v5.WithSubcommands(childCmd1, childCmd2),
)

// Add to CLI
v5.AddCommand(cli, cmd)
```

### Type Parameters

| Parameter | Description     | Example                     |
| --------- | --------------- | --------------------------- |
| `T`       | App config type | `AppConfig`                 |
| `F`       | Flags type      | `*GreetFlags`, `v5.NoFlags` |

Type parameters are inferred from the flags argument — you rarely need to specify them explicitly.

### NoFlags

Use `v5.NoFlags` for commands without flags:

```go
cmd, _ := v5.NewCommand("version", v5.NoFlags{},
    func(ctx context.Context, cfg *AppConfig, flags v5.NoFlags) error {
        fmt.Println("v5.0.0")
        return nil
    },
)
```

### Error Handling

All functions return errors — no panics in library code:

```go
cli, err := v5.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
if err != nil {
    return err
}
```

---

## Next Steps

- **[API Reference](../README.md)** - Full API documentation
- **[examples/](../examples/)** - Complete working examples
- **[Feature Status](../FEATURES.md)** - All features and their status
