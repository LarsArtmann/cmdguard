# cmdguard Features

This document provides a comprehensive reference for all cmdguard features across v1 and v2 APIs.

## API Versions

| Version | Package           | Description                                            |
| ------- | ----------------- | ------------------------------------------------------ |
| v1      | `pkg/cmdguard`    | Simple Cobra wrapper with panic-on-invalid guards      |
| v2      | `pkg/cmdguard/v2` | Full type-safe API with DI, lifecycle hooks, flag tags |

---

## v2 API Features

### Type-Safe Commands

Commands are generic over two type parameters:

- `T` - Application-level configuration type
- `F` - Command-specific flags type

```go
type Command[T any, F any] struct {
    Use      string
    Short    string
    Long     string
    Flags    F                                    // Typed flags
    RunE     func(ctx context.Context, cfg *T, flags F) error
    Commands []Command[T, F]                      // Subcommands share F type
}
```

### NoFlags Type Alias

For commands without flags, use `v2.NoFlags`:

```go
type NoFlags = struct{}

// Usage
cmd := v2.Command[Config, v2.NoFlags]{
    Use:  "version",
    RunE: func(ctx context.Context, cfg *Config, flags v2.NoFlags) error {
        return nil
    },
}
```

### Flag Tags

Define flags declaratively using struct tags:

| Tag       | Required | Description     | Example                 |
| --------- | -------- | --------------- | ----------------------- |
| `flag`    | Yes      | Flag name       | `flag:"verbose"`        |
| `short`   | No       | Short flag (-v) | `short:"v"`             |
| `default` | No       | Default value   | `default:"false"`       |
| `help`    | No       | Help text       | `help:"Enable verbose"` |

```go
type ServerFlags struct {
    Host    string `flag:"host" short:"H" default:"localhost" help:"Server host"`
    Port    int    `flag:"port" short:"p" default:"8080" help:"Server port"`
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Verbose output"`
    Timeout int    `flag:"timeout" default:"30" help:"Request timeout in seconds"`
}
```

**Supported types:** `string`, `int`, `int64`, `uint`, `uint64`, `float64`, `bool`, `[]string`, `time.Duration`

### AddAnyCommand for Mixed Flag Types

When commands need different flag types than the CLI root, use `AddAnyCommand`:

```go
// Root uses NoFlags
cli, _ := v2.New[Config, v2.NoFlags]("app", "...", Config{})

// Commands can have any F type
v2.AddAnyCommand(cli, v2.Command[Config, *GreetFlags]{...})
v2.AddAnyCommand(cli, v2.Command[Config, *ServerFlags]{...})

// Or use AddCommand for same F type
cli.AddCommand(v2.Command[Config, v2.NoFlags]{...})
```

### Dependency Injection

Built on `samber/do/v2`:

```go
// Register services
func setupDI(scope *v2.Scope, cfg Config) {
    // Provide with constructor
    v2.Provide(scope, func(i do.Injector) (*Database, error) {
        return &Database{URL: cfg.DBURL}, nil
    })

    // Provide value directly
    v2.ProvideValue(scope, &Logger{Level: cfg.LogLevel})
}

// Invoke in handlers
RunE: func(ctx context.Context, cfg *Config, flags *Flags) error {
    db := do.MustInvoke[*Database](cli.Scope())
    logger := do.MustInvoke[*Logger](cli.Scope())
    // ...
}
```

### Lifecycle Hooks

```go
cmd := v2.Command[Config, Flags]{
    PreRunE: func(ctx context.Context, cfg *Config, flags *Flags) error {
        // Called before RunE - use for validation
        return nil
    },
    RunE: func(ctx context.Context, cfg *Config, flags *Flags) error {
        // Main handler
        return nil
    },
    PostRunE: func(ctx context.Context, cfg *Config, flags *Flags) error {
        // Called after RunE (even on error) - use for cleanup
        return nil
    },
}
```

### Health Checks

```go
// Services can implement HealthChecker
type Database struct{}

func (d *Database) HealthCheck() error {
    // Check database connection
    return nil
}

// Register with DI
v2.Provide(scope, func(i do.Injector) (*Database, error) {
    return &Database{}, nil
})

// Run health checks
if err := cli.HealthCheck(); err != nil {
    log.Fatal("Health check failed:", err)
}
```

### Graceful Shutdown

```go
ctx := context.Background()

// Run CLI
cli.Execute(ctx)

// Graceful shutdown with timeout
shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := cli.Shutdown(shutdownCtx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### Functional Options

```go
cmd, err := v2.NewCommand[Config, Flags]("greet",
    v2.WithShort("Greet someone"),
    v2.WithLong("Send a friendly greeting to the specified person"),
    v2.WithAliases("hi", "hello"),
    v2.WithExample("myapp greet --name Alice"),
    v2.WithFlags(&Flags{Name: "World"}),
    v2.WithRunE(greetHandler),
    v2.WithPreRunE(validateHandler),
    v2.WithHidden(false),
    v2.WithDeprecated(""),
)
```

### Validation

Commands are validated at construction time:

| Check                           | Error Type          |
| ------------------------------- | ------------------- |
| Empty `Use` field               | `ErrMissingName`    |
| No `RunE` and no subcommands    | `ErrMissingHandler` |
| Invalid subcommands (recursive) | Wrapped error       |

```go
// Validate explicitly
if err := cmd.Validate(); err != nil {
    // Handle error
}

// Or use MustNewCommand which panics
cmd := v2.MustNewCommand[Config, Flags]("greet", opts...)
```

---

## v1 API Features

### GuardedCommand

Simple Cobra wrapper with panic-on-invalid guards:

```go
root := cmdguard.New("myapp", "My CLI application")
root.AddCommand(&cobra.Command{...})
root.ExecuteAndExit(context.Background())
```

### Panic-at-Construction-Time

Invalid commands panic immediately:

- Commands without `Run` or `RunE`
- Commands with empty `Use` field
- (Optional) Commands with `Run` instead of `RunE` in strict mode

### Strict Mode

Requires `RunE` instead of `Run`:

```go
os.Setenv("CMDGUARD_STRICT_MODE", "true")
root := cmdguard.New("myapp", "My CLI")
// root.AddCommand(&cobra.Command{Use: "bad", Run: ...}) // PANICS!
```

### Built-in Commands

| Command    | Description            |
| ---------- | ---------------------- |
| `version`  | Prints version info    |
| `validate` | Validates command tree |
| `help`     | Cobra's help command   |

### Global Flags (v1)

| Flag              | Env Var                | Default | Description                          |
| ----------------- | ---------------------- | ------- | ------------------------------------ |
| `--config, -c`    | `CMDGUARD_CONFIG_FILE` |         | Config file path                     |
| `--log-level, -l` | `CMDGUARD_LOG_LEVEL`   | `info`  | Log level (debug, info, warn, error) |
| `--strict, -s`    | `CMDGUARD_STRICT_MODE` | `false` | Enable strict mode                   |

---

## Error Types

```go
var (
    ErrInvalidCommand   = errors.New("invalid command")
    ErrMissingHandler   = errors.New("command has no handler")
    ErrMissingName      = errors.New("command has no name")
    ErrFlagParseFailed  = errors.New("flag parsing failed")
)
```

---

## Comparison with Raw Cobra

| Feature                   | Raw Cobra | cmdguard v1  | cmdguard v2      |
| ------------------------- | --------- | ------------ | ---------------- |
| Invalid command detection | Runtime   | Construction | Construction     |
| Type-safe config          | No        | No           | Yes              |
| Type-safe flags           | No        | No           | Yes              |
| Flag binding              | Manual    | Manual       | Automatic (tags) |
| DI integration            | Manual    | No           | Yes              |
| Lifecycle hooks           | Basic     | Basic        | Full (Pre/Post)  |
| Health checks             | No        | No           | Yes              |
| Graceful shutdown         | Manual    | No           | Yes              |

---

## Architecture

```
pkg/cmdguard/
├── guarded_command.go     # v1 API

pkg/cmdguard/v2/
├── command.go             # Command[T, F] type
├── guard.go               # GuardedCommand[T, F] type
├── scope.go               # DI scope wrapper
├── flags.go               # FlagRegistry for tag-based flags
├── errors.go              # Error types
└── *_test.go              # Tests
```

---

## Dependencies

| Library                         | Purpose       | Version |
| ------------------------------- | ------------- | ------- |
| `github.com/spf13/cobra`        | CLI framework | v1.10.2 |
| `github.com/samber/do/v2`       | DI container  | v2.0.0  |
| `github.com/charmbracelet/fang` | CLI styling   | v0.4.4  |

---

## Examples

See `examples/` directory:

| Example    | API | Description                    |
| ---------- | --- | ------------------------------ |
| `basic`    | v1  | Simple CLI with validation     |
| `advanced` | v1  | Nested commands, config        |
| `guarded`  | v1  | Strict mode, custom validation |
| `typed`    | v2  | Full type-safe CLI with DI     |
