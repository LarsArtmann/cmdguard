# cmdguard Features

**Updated:** 2026-06-10
**Status:** v2.4.0 — 0 lint issues, 0 race conditions

---

## API Versions

| Version | Package           | Description                                            | Status |
| ------- | ----------------- | ------------------------------------------------------ | ------ |
| v2      | `pkg/cmdguard/v2` | Full type-safe API with DI, lifecycle hooks, flag tags | Active |

> **Note:** v1 (`pkg/cmdguard`) has been deprecated. Use v2 for all new projects.

---

## v2 API Features

### Type-Safe Commands

Commands are generic over two type parameters via the constructor pattern:

- `T` — Application-level configuration type
- `F` — Command-specific flags type

```go
cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})

cmd, err := v2.NewCommand[AppConfig, *GreetFlags]("greet", greetHandler,
    v2.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v2.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v2.AddCommand(cli, cmd)
```

### NoFlags Type Alias

For commands without flags, use `v2.NoFlags`:

```go
cmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("version", handler,
    v2.WithShort[AppConfig, v2.NoFlags]("Print version"),
)
```

### Flag Tags

Define flags declaratively using struct tags:

| Tag       | Required | Description          | Example                 |
| --------- | -------- | -------------------- | ----------------------- |
| `flag`    | Yes      | Flag name            | `flag:"verbose"`        |
| `short`   | No       | Short flag (-v)      | `short:"v"`             |
| `default` | No       | Default value        | `default:"false"`       |
| `help`    | No       | Help text            | `help:"Enable verbose"` |
| `env`     | No       | Environment variable | `env:"VERBOSE"`         |
| `count`   | No       | Counting flag        | `count:"true"`          |
| `prompt`  | No       | Interactive prompt   | `prompt:"Server host?"` |
| `values`  | No       | Enum values          | `values:"json,text"`    |

```go
type ServerFlags struct {
    Host    string `flag:"host" short:"H" default:"localhost" help:"Server host" env:"HOST"`
    Port    int    `flag:"port" short:"p" default:"8080" help:"Server port"`
    Verbose int    `flag:"verbose" short:"v" default:"0" help:"Verbosity level" count:"true"`
}
```

**Supported types:** `string`, `int`, `int64`, `uint`, `uint64`, `float64`, `bool`, `[]string`, `time.Duration`

### Config File Loading

Load configuration from JSON, YAML, or TOML files before flag parsing:

```go
cli, _ := v2.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v2.WithConfigFile[AppConfig]("/etc/myapp/config.json", "~/.myapp/config.json"),
)

// Or use YAML/TOML with custom loader
cli, _ := v2.NewCLI[AppConfig]("myapp", "My app", AppConfig{},
    v2.WithConfigFileLoader[AppConfig](configload.YAML(), "config.yaml"),
)
```

**Precedence:** explicit flag → env var → config file → default value (highest to lowest priority).

### Dependency Injection

Built on `samber/do/v2`:

```go
scope := cli.Scope()

v2.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{DSN: cfg.DSN}, nil
})
v2.ProvideValue(scope, &Logger{Level: "info"})

// Invoke in command handlers
db, err := v2.Invoke[*Database](scope)
```

### Lifecycle Hooks

```go
cmd, _ := v2.NewCommand[AppConfig, *Flags]("example", runHandler,
    v2.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // validation
    }),
    v2.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // cleanup (only called on success)
    }),
)
```

### Middleware

Chain middleware around every command handler:

```go
cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
    v2.WithMiddleware(v2.TimingMiddleware[Config](logTiming)),
    v2.WithMiddleware(v2.RecoveryMiddleware[Config]()),
    v2.WithMiddleware(v2.SpinnerMiddleware[Config]("Loading...")),
    v2.WithTelemetry[Config](tracer), // OpenTelemetry spans
)
```

| Middleware               | Purpose                                    |
| ------------------------ | ------------------------------------------ |
| `TimingMiddleware[T]`    | Logs execution duration with error context |
| `RecoveryMiddleware[T]`  | Catches panics, returns error              |
| `SpinnerMiddleware[T]`   | Terminal spinner (auto-skips non-TTY)      |
| `TelemetryMiddleware[T]` | OpenTelemetry span per command             |

### Interactive Prompts

Automatically prompt for missing flags using `huh`:

```go
type Flags struct {
    Name    string `flag:"name" prompt:"What is your name?"`
    Role    string `flag:"role" prompt:"Select role" values:"admin,user,viewer"`
    Confirm bool   `flag:"confirm" prompt:"Proceed?"`
}

cmd, _ := v2.NewCommand[Config, *Flags]("setup", handler,
    v2.WithFlags[Config, *Flags](&Flags{}),
    v2.WithPromptOnMissing[Config, *Flags](),
)
```

### Rich Output

Output tables and results in 12+ formats:

```go
headers := []string{"Name", "Status"}
rows := [][]string{{"db-1", "running"}}
v2.OutputTable("json", headers, rows)
```

**Formats:** table, json, csv, yaml, toml, html, markdown, d2, and more.

### Markdown Help Rendering

Render command help text as styled markdown via `glamour`:

```go
cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
    v2.WithGlamourHelp[Config](),
)
```

### Color Control

`--no-color` flag and `NO_COLOR` env var are registered by default:

```go
if cli.NoColor() {
    // disable color output
}
```

### Custom Value Types

9 built-in validated types:

| Type        | File                | Validation            |
| ----------- | ------------------- | --------------------- |
| `Duration`  | `types_duration.go` | Valid time.Duration   |
| `Email`     | `types_email.go`    | RFC 5322 email format |
| `Enum[T]`   | `types_enum.go`     | Allowed values        |
| `FilePath`  | `types_filepath.go` | Path validation       |
| `HostPort`  | `types_hostport.go` | host:port format      |
| `LogFormat` | `types_log.go`      | text/json format      |
| `LogLevel`  | `types_log.go`      | Standard log levels   |
| `Port`      | `types_port.go`     | 1–65535 range         |
| `URL`       | `types_url.go`      | Valid URL             |

### Branching Flow Context

Track command execution paths and propagate values through the command hierarchy:

```go
bfc, ok := v2.GetBranchingFlowContext(ctx)
bfc.PathString()           // "app.subcmd"
bfc.SetValue(key, val)     // propagates to children
bfc.GetValue(key)          // looks up hierarchy
bfc.BranchWithDuration("slow-op", 5*time.Second)
```

### Command Validation

| Mode      | CLI Option                     | Enforces                                |
| --------- | ------------------------------ | --------------------------------------- |
| Lenient   | (default)                      | Basic validation                        |
| Strict    | `WithStrictValidation[T]()`    | `WithShort` on all commands             |
| Draconian | `WithDraconianValidation[T]()` | Strict + `WithExample` on leaf commands |

### Error Handling

60 sentinel errors for `errors.Is()` chainability:

```go
errors.Is(err, v2.ErrInvalidCommand)
errors.Is(err, v2.ErrMissingName)
errors.Is(err, v2.ErrDuplicateCommand)

// Rich error constructors
v2.NewCommandError(name, err)
v2.NewServiceError(type, err)
v2.NewFlagError(name, err)
v2.NewFlagErrorWithSuggestion(name, err, suggestion)
v2.NewExitError(code, err)  // custom exit codes
```

### Shell Completion

Built-in shell completion for bash, zsh, fish, and PowerShell.

### Man Page Generation

Generate man pages from command definitions:

```go
cmd, err := v2.GenerateManPageCommand[Config](cli)
if err != nil {
    panic(err)
}
v2.AddCommand(cli, cmd)
```

### Health Checks & Graceful Shutdown

```go
if err := cli.HealthCheck(); err != nil {
    log.Fatal("Health check failed:", err)
}

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
cli.Shutdown(shutdownCtx)
```

### $EDITOR Support

Open user's editor for multi-line input:

```go
content, err := v2.EditInEditor(ctx, []byte("template"))
```

---

## CLI Options Reference

| Option                         | Purpose                                    |
| ------------------------------ | ------------------------------------------ |
| `WithCLIVersion[T](v)`         | Set version string                         |
| `WithCLILong[T](desc)`         | Set long description                       |
| `WithCLIScope[T](scope)`       | Set custom DI scope                        |
| `WithSilenceErrors[T]()`       | Suppress cobra error printing              |
| `WithSilenceUsage[T]()`        | Suppress usage on error                    |
| `WithFang[T](bool)`            | Enable/disable fang styling                |
| `WithMiddleware[T](mw...)`     | Middleware wrapping every handler          |
| `WithGroup[T](id, title)`      | Register command group on root             |
| `WithEnvPrefix[T](prefix)`     | Prefix for env var lookups                 |
| `WithSignalHandling[T]()`      | Cancel context on SIGINT/SIGTERM           |
| `WithConfigValidation[T](fn)`  | Validate config after flag parsing         |
| `WithStrictValidation[T]()`    | Require short descriptions on all commands |
| `WithDraconianValidation[T]()` | Strict + examples on leaf commands         |
| `WithConfigFile[T](paths...)`  | Load JSON config file before flags         |
| `WithConfigFileLoader[T](l,p)` | Load config with custom loader (YAML/TOML) |
| `WithGlamourHelp[T]()`         | Render markdown in command help text       |
| `WithTelemetry[T](tracer)`     | OpenTelemetry spans for all commands       |

---

## Command Options Reference

| Option                           | Purpose                              |
| -------------------------------- | ------------------------------------ |
| `WithShort[T, F](short)`         | Short description                    |
| `WithLong[T, F](long)`           | Long description                     |
| `WithAliases[T, F](aliases...)`  | Alternative names                    |
| `WithExample[T, F](example)`     | Example usage                        |
| `WithFlags[T, F](flags)`         | Typed flags struct                   |
| `WithRunE[T, F](runE)`           | Main handler                         |
| `WithPreRunE[T, F](preRunE)`     | Pre-validation hook                  |
| `WithPostRunE[T, F](postRunE)`   | Post-success cleanup hook            |
| `WithSubcommands[T, F](cmds...)` | Child commands                       |
| `WithHidden[T, F](hidden)`       | Hide from help                       |
| `WithDeprecated[T, F](msg)`      | Deprecation message                  |
| `WithGroupID[T, F](group)`       | Help group name                      |
| `WithExactArgs[T, F](n)`         | Require exactly n positional args    |
| `WithMinimumArgs[T, F](n)`       | Require at least n positional args   |
| `WithMaximumArgs[T, F](n)`       | Allow at most n positional args      |
| `WithRangeArgs[T, F](min, max)`  | Require between min and max args     |
| `WithNoArgs[T, F]()`             | Reject any positional args           |
| `WithArgs[T, F](fn)`             | Custom positional args validator     |
| `WithPromptOnMissing[T, F]()`    | Interactive prompt for missing flags |

---

## Comparison with Raw Cobra

| Feature                   | Raw Cobra | cmdguard v2      |
| ------------------------- | --------- | ---------------- |
| Invalid command detection | Runtime   | Construction     |
| Type-safe config          | No        | Yes              |
| Type-safe flags           | No        | Yes              |
| Flag binding              | Manual    | Automatic (tags) |
| DI integration            | Manual    | Yes              |
| Lifecycle hooks           | Basic     | Full (Pre/Post)  |
| Health checks             | No        | Yes              |
| Graceful shutdown         | Manual    | Yes              |
| Middleware                | No        | Yes              |
| Interactive prompts       | No        | Yes              |
| Config file loading       | Manual    | Built-in         |
| Output formatting         | Manual    | 12+ formats      |
| Shell completion          | Manual    | Built-in         |
| Man page generation       | Manual    | Built-in         |
| Markdown help             | No        | Built-in         |
| Telemetry                 | Manual    | Built-in         |

---

## Architecture

```
pkg/cmdguard/v2/
├── cli.go                # CLI[T] struct, NewCLI, AddCommand, Execute
├── cli_accessors.go      # CLI accessor methods
├── cli_command.go         # Internal cobra wiring
├── cli_options.go         # CLI functional options
├── command.go             # Command[T,F] struct, constructors, options
├── config.go              # Config type constraint
├── config_file.go         # ConfigFileLoader, JSON loader
├── configload/            # Optional YAML/TOML loaders
├── errors.go              # Error types (CommandError, FlagError, etc.)
├── errors_command.go      # Command sentinel errors
├── errors_config.go       # Config sentinel errors
├── errors_di.go           # DI sentinel errors
├── errors_flags.go        # Flag sentinel errors
├── flags.go               # FlagRegistry with struct tags
├── flow_context.go        # BranchingFlowContext
├── glamour.go             # Markdown help rendering
├── middleware.go           # Middleware chain (Timing, Recovery, Spinner, Telemetry)
├── output.go              # Rich output (12 formats)
├── prompts.go             # Interactive prompts (huh)
├── scope.go               # DI scope wrapping samber/do/v2
├── spinner.go             # Terminal spinner
├── telemetry.go           # OpenTelemetry integration
├── type_handler.go        # Extensible type registry
└── types_*.go             # Custom value types
```

---

## Dependencies

| Library                            | Purpose              | Version |
| ---------------------------------- | -------------------- | ------- |
| `github.com/spf13/cobra`           | CLI framework        | v1.10.2 |
| `github.com/samber/do/v2`          | Dependency injection | v2.0.0  |
| `charm.land/fang/v2`               | CLI styling          | v2.0.1  |
| `charm.land/huh/v2`                | Interactive prompts  | v2.0.3  |
| `charm.land/glamour/v2`            | Markdown rendering   | v2.0.0  |
| `go.opentelemetry.io/otel/trace`   | OpenTelemetry spans  | v1.44.0 |
| `github.com/larsartmann/go-output` | Rich output formats  | v0.7.2  |
