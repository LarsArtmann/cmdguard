# cmdguard API Reference

Extracted from AGENTS.md for conciseness. See the [pkg.go.dev reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2) for the complete API.

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
| `WithCLIVersion[T](v)`            | Set version string (auto-pipes to fang)      |
| `WithCLICommit[T](commit)`        | Set git commit hash (auto-pipes to fang)     |
| `WithCLILong[T](desc)`            | Set long description                         |
| `WithCLIScope[T](scope)`          | Set custom DI scope                          |
| `WithSilenceErrors[T]()`          | Suppress cobra error printing                |
| `WithSilenceUsage[T]()`           | Suppress usage on error                      |
| `WithFang[T](bool)`               | Enable/disable fang styling (preferred)      |
| `WithFangOptions[T](opts...)`     | Custom fang options                          |
| `WithFangErrorHandler[T](handler)` | Custom fang error handler                   |
| `WithFangColorScheme[T](cs)`      | Custom fang color scheme                     |
| `WithMiddleware[T](mw...)`        | Middleware wrapping every handler            |
| `WithGroup[T](id, title)`         | Register command group on root               |
| `WithEnvPrefix[T](prefix)`        | Prefix for env var lookups                   |
| `WithSignalHandling[T]()`         | Cancel context on SIGINT/SIGTERM             |
| `WithGracefulShutdown[T]()`       | Graceful DI shutdown on SIGINT/SIGTERM       |
| `WithDILogging[T](logf)`          | Internal DI container logging                |
| `WithConfigValidation[T](fn)`     | Validate config after flag parsing           |
| `WithStrictValidation[T]()`       | Require short descriptions on all commands   |
| `WithDraconianValidation[T]()`    | Strict + examples on leaf commands           |
| `WithConfigFile[T](paths...)`     | Load JSON config file before flags           |
| `WithConfigFileLoader[T](l,p)`    | Load config with custom loader (YAML/TOML)   |
| `WithGlamourHelp[T]()`            | Render markdown in command help text         |
| `WithTelemetry[T](tracer)`        | OpenTelemetry spans for all commands         |

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

    "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
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

// Testing — clone scope and override services
cloned := v2.CloneScope(scope)
v2.OverrideValue(cloned, &MockDatabase{})
mockDB, _ := v2.Invoke[*Database](cloned) // returns mock
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
    v2.WithMiddleware(v2.TimingMiddleware[Config](func(name string, d time.Duration, err error) {
        log.Printf("%s took %v (err=%v)", name, d, err)
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
cmd, err := v2.VersionCommand[Config](cli)
if err != nil {
    panic(err)
}
v2.AddCommand(cli, cmd)
```

### Doctor Command

```go
// Simple — just DI health checks
cmd, err := v2.DoctorCommand[Config](cli)
if err != nil {
    panic(err)
}
v2.AddCommand(cli, cmd)

// With custom diagnostic checks and group
cmd, err := v2.DoctorCommand[Config](cli,
    v2.WithDoctorCheck[Config]("database", func(ctx context.Context) error {
        return db.Ping(ctx)
    }),
    v2.WithDoctorGroupID[Config]("system"),
)
if err != nil {
    panic(err)
}
v2.AddCommand(cli, cmd)

// Per-service results programmatically
results := cli.HealthCheckResultsWithContext(ctx)
for name, err := range results {
    if err != nil { fmt.Printf("✗ %s: %v\n", name, err) }
}
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
