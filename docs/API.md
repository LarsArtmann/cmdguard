# cmdguard API Reference

Extracted from AGENTS.md for conciseness. See the [pkg.go.dev reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3) for the complete API.

## API Reference

### Architecture: CLI[T] + Command[T, F]

`CLI[T]` has one type parameter (config type). Each command gets its own flags type via `Command[T, F]`. Because Go doesn't support additional type parameters on methods, `AddCommand` is a standalone function.

Commands are created via constructors — `NewCommand` for leaf commands, `NewParentCommand` for commands with subcommands. Struct fields are unexported to enforce validation at construction time.

```go
cli, err := v3.NewCLI[AppConfig]("myapp", "My CLI", AppConfig{})
cmd, err := v3.NewCommand[AppConfig, *GreetFlags]("greet", greetHandler,
    v3.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v3.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v3.AddCommand(cli, cmd)
```

### Command Constructors

```go
// Leaf command with handler
func NewCommand[T, F any](use string, runE func(ctx context.Context, cfg *T, flags F) error, opts ...CommandOption[T, F]) (Command[T, F], error)

// Parent command with subcommands
func NewParentCommand[T, F any](use string, long string, subcommands []Command[T, F], opts ...CommandOption[T, F]) (Command[T, F], error)
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
cli, err := v3.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
```

Functional options:

| Option                             | Purpose                                    |
| ---------------------------------- | ------------------------------------------ |
| `WithCLIVersion[T](v)`             | Set version string (auto-pipes to fang)    |
| `WithCLICommit[T](commit)`         | Set git commit hash (auto-pipes to fang)   |
| `WithCLILong[T](desc)`             | Set long description                       |
| `WithCLIScope[T](scope)`           | Set custom DI scope                        |
| `WithSilenceErrors[T]()`           | Suppress cobra error printing (advanced)   |
| `WithSilenceUsage[T]()`            | Suppress usage on error (**default**)      |
| `WithFang[T](bool)`                | Enable/disable fang styling (preferred)    |
| `WithFangOptions[T](opts...)`      | Custom fang options                        |
| `WithFangErrorHandler[T](handler)` | Custom fang error handler                  |
| `WithFangColorScheme[T](cs)`       | Custom fang color scheme                   |
| `WithMiddleware[T](mw...)`         | Middleware wrapping every handler          |
| `WithGroup[T](id, title)`          | Register command group on root             |
| `WithEnvPrefix[T](prefix)`         | Prefix for env var lookups                 |
| `WithSignalHandling[T]()`          | Cancel context on SIGINT/SIGTERM           |
| `WithGracefulShutdown[T]()`        | Graceful DI shutdown on SIGINT/SIGTERM     |
| `WithDILogging[T](logf)`           | Internal DI container logging              |
| `WithConfigValidation[T](fn)`      | Validate config after flag parsing         |
| `WithStrictValidation[T]()`        | Require short descriptions on all commands |
| `WithDraconianValidation[T]()`     | Strict + examples on leaf commands         |
| `WithConfigFile[T](paths...)`      | Load JSON config file before flags         |
| `WithConfigFileLoader[T](l,p)`     | Load config with custom loader (YAML/TOML) |
| `WithGlamourHelp[T]()`             | Render markdown in command help text       |
| `WithTelemetry[T](tracer)`         | OpenTelemetry spans for all commands       |
| `WithPostFlagParse[T](fns...)`     | Post-parse hook: DI init, session storage  |
| `WithCleanup[T](fns...)`           | Post-RunE cleanup (fires on error too)     |

### CLI[T] Methods

| Method                               | Returns                 | Purpose                              |
| ------------------------------------ | ----------------------- | ------------------------------------ |
| `Execute(ctx)`                       | `error`                 | Run CLI with context                 |
| `ExecuteWithArgs(ctx, args)`         | `error`                 | Run with specific args               |
| `ExecuteAndExit(ctx)`                |                         | Run and os.Exit (respects ExitCoder) |
| `Scope()`                            | `*Scope`                | DI scope                             |
| `Injector()`                         | `do.Injector`           | Raw samber/do injector               |
| `Config()`                           | `*T`                    | Typed config                         |
| `SetConfig(cfg)`                     |                         | Update config                        |
| `RootCommand()`                      | `*cobra.Command`        | Underlying cobra command             |
| `Shutdown(ctx)`                      | `error`                 | Graceful shutdown                    |
| `HealthCheck()`                      | `error`                 | Run health checks                    |
| `HealthCheckWithContext(ctx)`        | `error`                 | Health checks with context           |
| `HealthCheckResults()`               | `map[string]error`      | Per-service health map               |
| `HealthCheckResultsWithContext(ctx)` | `map[string]error`      | Per-service health map with context  |
| `SetVersion(v)`                      |                         | Set version at runtime               |
| `SetLong(desc)`                      |                         | Set long description                 |
| `FlowContext()`                      | `*BranchingFlowContext` | Path tracking (nil until Execute)    |
| `AddGlobalFlag(...)`                 |                         | Persistent string flag               |
| `AddGlobalBoolFlag(...)`             |                         | Persistent bool flag                 |
| `NoColor()`                          | `bool`                  | True if `--no-color` was passed      |

### Basic Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
    Output  string `flag:"output" short:"o" default:"text" help:"Output format"`
}

func main() {
    cli, err := v3.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
    if err != nil {
        panic(err)
    }

    cmd, err := v3.NewCommand[AppConfig, v3.NoFlags]("hello",
        func(ctx context.Context, cfg *AppConfig, flags v3.NoFlags) error {
            fmt.Printf("Hello! Verbose: %v\n", cfg.Verbose)
            return nil
        },
        v3.WithShort[AppConfig, v3.NoFlags]("Say hello"),
    )
    if err != nil {
        panic(err)
    }

    if err := v3.AddCommand(cli, cmd); err != nil {
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

greetCmd, err := v3.NewCommand[AppConfig, *GreetFlags]("greet",
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
    v3.WithShort[AppConfig, *GreetFlags]("Greet someone"),
    v3.WithFlags[AppConfig, *GreetFlags](&GreetFlags{}),
)
v3.AddCommand(cli, greetCmd)
```

### Subcommands

```go
listCmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("list",
    listUsersHandler, v3.WithShort[AppConfig, v3.NoFlags]("List users"),
)
createCmd, _ := v3.NewCommand[AppConfig, v3.NoFlags]("create",
    createUserHandler, v3.WithShort[AppConfig, v3.NoFlags]("Create user"),
)
userCmd, err := v3.NewParentCommand[AppConfig, v3.NoFlags]("user",
    "User management", []v3.Command[AppConfig, v3.NoFlags]{listCmd, createCmd},
    v3.WithShort[AppConfig, v3.NoFlags]("User management"),
)
v3.AddCommand(cli, userCmd)
```

### Dependency Injection

```go
cli, _ := v3.NewCLI[AppConfig]("myapp", "My app", AppConfig{})
scope := cli.Scope()

// Register services
v3.Provide(scope, func(i do.Injector) (*Database, error) {
    cfg, _ := v3.Invoke[*AppConfig](scope)
    return &Database{DSN: cfg.DSN}, nil
})
v3.ProvideValue(scope, &Logger{Level: "info"})

// Invoke in command handlers
db, err := v3.Invoke[*Database](cli.Scope())

// Testing — clone scope and override services
cloned := v3.CloneScope(scope)
v3.OverrideValue(cloned, &MockDatabase{})
mockDB, _ := v3.Invoke[*Database](cloned) // returns mock
```

#### DI Scope Functions

| Function                           | Purpose                                 |
| ---------------------------------- | --------------------------------------- |
| `Provide[T](scope, provider)`      | Lazy singleton registration             |
| `ProvideNamed[T](scope, name, fn)` | Named service registration              |
| `ProvideValue[T](scope, value)`    | Eager value registration                |
| `Invoke[T](scope)`                 | Retrieve singleton service              |
| `InvokeNamed[T](scope, name)`      | Retrieve named service                  |
| `Override[T](scope, provider)`     | Replace service provider (testing)      |
| `OverrideValue[T](scope, value)`   | Replace pre-constructed value (testing) |
| `CloneScope(scope)`                | Clone scope for test isolation          |
| `NewScopeWithOpts(name, opts)`     | Create scope with custom InjectorOpts   |
| `Scope.Child(name)`                | Create child scope                      |
| `Scope.Shutdown(ctx)`              | Graceful shutdown of scope services     |
| `Scope.ShutdownAll(ctx)`           | Shutdown scope + all parents            |

### Lifecycle Hooks

```go
cmd, err := v3.NewCommand[AppConfig, *Flags]("example", runHandler,
    v3.WithPreRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // validation
    }),
    v3.WithPostRunE[AppConfig, *Flags](func(ctx context.Context, cfg *AppConfig, flags *Flags) error {
        return nil // cleanup (only called on success)
    }),
)
```

### Raw Cobra Subcommands (Escape Hatch)

Consumers migrating from plain Cobra can register raw `*cobra.Command` subcommands
via `cli.RootCommand().AddCommand(...)`. cmdguard stores the resolved config in the
command context during `PersistentPreRunE`, so raw handlers retrieve it without a
parallel context-key system:

```go
cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    // Initialise DI / store session once flags are parsed
    v3.WithPostFlagParse[Config](func(cmd *cobra.Command, cfg *Config) error {
        return initDI(cfg) // runs after parsing + config validation
    }),
)

root := cli.RootCommand()
root.AddCommand(&cobra.Command{
    Use: "raw",
    RunE: func(cmd *cobra.Command, _ []string) error {
        cfg, ok := v3.ConfigFromContext[Config](cmd.Context())
        if !ok {
            return errors.New("config not initialised")
        }
        // use cfg...
        return nil
    },
})
```

`WithCleanup[T]` registers hooks that run after every command's `RunE` —
including when `RunE` errors, which Cobra's `PostRunE`/`PersistentPostRunE`
silently skip. The hook receives `(cmd, cfg, runErr)`, so a single hook can
flush buffers or release resources on both success and failure. It wraps every
`RunE` in the tree at execute time, so it covers raw cobra subcommands too.
The original `RunE` error is never swallowed; cleanup errors are joined to it.

Scoped flags (`local:"true"`) keep root-only flags out of subcommand `--help`. Use
`cli.RegisterLocalCommandFlags(cmd)` to re-register a local flag group on a
subcommand that needs it (e.g. one that re-runs the root pipeline).

### Middleware

```go
// Timing middleware — logs execution duration
cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    v3.WithMiddleware(v3.TimingMiddleware[Config](func(name string, d time.Duration, err error) {
        log.Printf("%s took %v (err=%v)", name, d, err)
    })),
)

// Recovery middleware — catches panics
cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    v3.WithMiddleware(v3.RecoveryMiddleware[Config]()),
)

// Spinner middleware — shows terminal spinner
cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    v3.WithMiddleware(v3.SpinnerMiddleware[Config]("Loading...")),
)

// Telemetry middleware — OpenTelemetry spans
import "go.opentelemetry.io/otel"
tracer := otel.Tracer("myapp")
cli, _ := v3.NewCLI[Config]("app", "My app", Config{},
    v3.WithTelemetry[Config](tracer), // or WithMiddleware(TelemetryMiddleware[Config](tracer))
)
```

### BranchingFlowContext

Automatically created on `Execute`. Access via `GetBranchingFlowContext(ctx)` in handlers.

```go
bfc, ok := v3.GetBranchingFlowContext(ctx)
bfc.PathString()  // "app.subcmd"
bfc.SetValue(key, val)  // propagates to children
bfc.GetValue(key)       // looks up hierarchy
```

### Error Handling

```go
// All v2 functions return errors
cli, err := v3.NewCLI[Config]("app", "My app", Config{})
cmd, err := v3.NewCommand[Config, NoFlags]("test", handler)

// Sentinel errors for errors.Is()
errors.Is(err, v3.ErrInvalidCommand)
errors.Is(err, v3.ErrMissingName)
errors.Is(err, v3.ErrDuplicateCommand)
errors.Is(err, v3.ErrMissingHandler)

// Rich error types
v3.NewCommandError(name, err)    // wraps with command context
v3.NewServiceError(type, err)    // wraps with DI service context
v3.NewFlagError(name, err)       // wraps with flag context
v3.NewFlagErrorWithSuggestion(name, err, suggestion)  // includes typo fix

// Exit codes
v3.NewExitError(code, err)       // error with custom exit code for ExecuteAndExit
errors.As(err, &exitCoder)       // check if error implements ExitCoder
```

### Version Command

```go
cli, _ := v3.NewCLI[Config]("myapp", "My app", Config{},
    v3.WithCLIVersion[Config]("1.0.0"),
)
cmd, err := v3.VersionCommand[Config](cli)
if err != nil {
    panic(err)
}
v3.AddCommand(cli, cmd)
```

### Doctor Command

```go
// Simple — just DI health checks
cmd, err := v3.DoctorCommand[Config](cli)
if err != nil {
    panic(err)
}
v3.AddCommand(cli, cmd)

// With custom diagnostic checks and group
cmd, err := v3.DoctorCommand[Config](cli,
    v3.WithDoctorCheck[Config]("database", func(ctx context.Context) error {
        return db.Ping(ctx)
    }),
    v3.WithDoctorGroupID[Config]("system"),
)
if err != nil {
    panic(err)
}
v3.AddCommand(cli, cmd)

// Per-service results programmatically
results := cli.HealthCheckResultsWithContext(ctx)
for name, err := range results {
    if err != nil { fmt.Printf("✗ %s: %v\n", name, err) }
}
```

### Markdown Help (glamour)

```go
cli, _ := v3.NewCLI[Config]("myapp", "My app", Config{},
    v3.WithGlamourHelp[Config](),
)
// Command Long and Example fields are rendered as markdown in terminal help
```

### Strict Validation

```go
cli, _ := v3.NewCLI[Config]("myapp", "My app", Config{},
    v3.WithStrictValidation[Config](),  // requires WithShort on all commands
)
```

### Config Validation

```go
cli, _ := v3.NewCLI[Config]("myapp", "My app", Config{},
    v3.WithConfigValidation[Config](func(cfg *Config) error {
        if cfg.Port < 1 { return fmt.Errorf("invalid port") }
        return nil
    }),
)
```

---

### Audit Log (samber-do-auditlog)

Wire DI audit logging into the CLI's injector, then export or query the captured snapshot.

```go
import auditlog "github.com/larsartmann/samber-do-auditlog"

plugin, _ := auditlog.New(auditlog.Config{Enabled: true, ContainerID: "myapp"})

cli, _ := v3.NewCLI[Config]("myapp", "My app", Config{},
    v3.WithAuditLog[Config](plugin),
)

// after Execute, export the snapshot
format, _ := v3.ParseAuditLogFormat(os.Getenv("AUDIT_LOG_FORMAT")) // "" -> html
_ = v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
    Format: format,            // html | json | ndjson | csv | tsv | mermaid | dot | d2 | plantuml | tree | htmltree
    Path:   "myapp-audit." + format.String(),
})
```

Accessors and query helpers:

```go
cli.AuditLog()                 // *auditlog.Plugin (nil if not configured)
cli.AuditLogReport()           // *auditlog.Report snapshot (nil if not configured)
cli.RecordAuditHealthCheck(ctx) // map[string]error

v3.AuditLogServiceByName[Config](cli, "taskStore")  // *auditlog.ServiceInfo
v3.AuditLogFailedServices[Config](cli)              // []auditlog.ServiceInfo
```

`AuditLogFormat` is a validated enum &mdash; build it via `ParseAuditLogFormat` so an
invalid value surfaces as `ErrUnsupportedAuditLogFormat` rather than a silent failure.

For the full 16-format service summary table (table, json, csv, tsv, markdown, xml, d2,
yaml, html, tree, mermaid, dot, jsonl, asciidoc, toml, plantuml), call the plugin directly:

```go
cli.AuditLog().ExportToTable(path, output.FormatCSV, output.RenderOptions{})
```

---
