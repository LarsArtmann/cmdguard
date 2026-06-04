# Migrating from Plain Cobra to cmdguard

> **Audience:** Go developers with existing Cobra CLI applications who want type-safe flags, DI, and validation without rewriting everything.
> **Time:** 15 minutes for basic migration, 45 minutes for full adoption.

---

## Why Migrate?

| Pain Point (Cobra)                                                  | cmdguard Solution                                               |
| ------------------------------------------------------------------- | --------------------------------------------------------------- |
| `cmd.Flags().GetString("name")` — runtime lookups, easy to misspell | Struct tags: `Name string \`flag:"name"\`` — compile-time typed |
| Flags defined far from where they're used                           | Flags live in a typed struct next to the handler                |
| Manual validation boilerplate                                       | `required:"true"`, custom types, `WithPreRunE` hooks            |
| No built-in DI                                                      | `Provide`/`Invoke` with samber/do/v2                            |
| Missing flags fail at runtime                                       | Missing flags fail at `AddCommand` time                         |
| No typo suggestions for users                                       | "did you mean?" for flags and subcommands                       |

---

## Migration Strategy: Incremental Adoption

You don't need to rewrite your entire CLI. cmdguard wraps Cobra — you can migrate one command at a time.

### Phase 1: Wrap Your Existing Cobra App (5 minutes)

The simplest possible adoption: use cmdguard for the root CLI and keep your existing `*cobra.Command` trees.

**Before (plain Cobra):**

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "My application",
}

var greetCmd = &cobra.Command{
    Use:   "greet",
    Short: "Greet someone",
    Run: func(cmd *cobra.Command, args []string) {
        name, _ := cmd.Flags().GetString("name")
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    greetCmd.Flags().StringP("name", "n", "World", "Name to greet")
    rootCmd.AddCommand(greetCmd)
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**After (cmdguard root, existing commands untouched):**

```go
package main

import (
    "context"
    "fmt"
    "os"

    v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
    "github.com/spf13/cobra"
)

type AppConfig struct {
    Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
}

// Keep your existing cobra commands exactly as they are.
var greetCmd = &cobra.Command{
    Use:   "greet",
    Short: "Greet someone",
    Run: func(cmd *cobra.Command, args []string) {
        name, _ := cmd.Flags().GetString("name")
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    greetCmd.Flags().StringP("name", "n", "World", "Name to greet")
}

func main() {
    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Add your existing cobra command directly.
    cli.RootCommand().AddCommand(greetCmd)

    if err := cli.Execute(context.Background()); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**What you get immediately:**

- Type-safe root config (`AppConfig`)
- Environment variable support (`env:"VERBOSE"`)
- Optional DI scope for new services
- All existing commands continue to work unchanged

---

### Phase 2: Migrate One Command to Typed Flags (10 minutes)

Pick a command with simple flags and convert it to cmdguard's typed struct pattern.

**Before:**

```go
var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy the application",
    Run: func(cmd *cobra.Command, args []string) {
        env, _ := cmd.Flags().GetString("env")
        dryRun, _ := cmd.Flags().GetBool("dry-run")
        timeout, _ := cmd.Flags().GetDuration("timeout")

        if env == "" {
            fmt.Fprintln(os.Stderr, "--env is required")
            os.Exit(1)
        }

        fmt.Printf("Deploying to %s (dry-run=%v, timeout=%v)\n", env, dryRun, timeout)
    },
}

func init() {
    deployCmd.Flags().StringP("env", "e", "", "Target environment")
    deployCmd.Flags().BoolP("dry-run", "d", false, "Simulate deployment")
    deployCmd.Flags().DurationP("timeout", "t", 5*time.Minute, "Deployment timeout")
    _ = deployCmd.MarkFlagRequired("env")
}
```

**After:**

```go
type DeployFlags struct {
    Env     string        `flag:"env"     short:"e" required:"true" help:"Target environment"`
    DryRun  bool          `flag:"dry-run" short:"d" default:"false" help:"Simulate deployment"`
    Timeout time.Duration `flag:"timeout" short:"t" default:"5m"    help:"Deployment timeout"`
}

func main() {
    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    deployCmd, err := v2.NewCommand[AppConfig, *DeployFlags]("deploy",
        func(ctx context.Context, cfg *AppConfig, flags *DeployFlags) error {
            fmt.Printf("Deploying to %s (dry-run=%v, timeout=%v)\n",
                flags.Env, flags.DryRun, flags.Timeout)
            return nil
        },
        v2.WithShort[AppConfig, *DeployFlags]("Deploy the application"),
        v2.WithFlags[AppConfig, *DeployFlags](&DeployFlags{}),
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if err := v2.AddCommand(cli, deployCmd); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Existing cobra commands still work alongside cmdguard commands.
    cli.RootCommand().AddCommand(greetCmd)

    cli.ExecuteAndExit(context.Background())
}
```

**What changed:**

- Flags moved from `init()` + `Flags().Get*` to a typed struct
- `required:"true"` replaces `MarkFlagRequired`
- Handler receives typed `*DeployFlags` instead of querying cobra flags
- Validation happens at `NewCommand` time, not runtime

---

### Phase 3: Add Dependency Injection (15 minutes)

Register shared services on the CLI scope and inject them into handlers.

```go
type Database struct {
    DSN string
}

func (db *Database) Query(ctx context.Context) error {
    fmt.Printf("Querying database at %s\n", db.DSN)
    return nil
}

func main() {
    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{})
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Register services in the DI scope.
    scope := cli.Scope()

    v2.Provide(scope, func(i do.Injector) (*Database, error) {
        cfg, _ := v2.Invoke[*AppConfig](scope)
        return &Database{DSN: cfg.DatabaseURL}, nil
    })

    queryCmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("query",
        func(ctx context.Context, cfg *AppConfig, flags v2.NoFlags) error {
            db, err := v2.Invoke[*Database](scope)
            if err != nil {
                return err
            }
            return db.Query(ctx)
        },
        v2.WithShort[AppConfig, v2.NoFlags]("Run a database query"),
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    v2.AddCommand(cli, queryCmd)
    cli.ExecuteAndExit(context.Background())
}
```

---

### Phase 4: Full Adoption — Subcommands, Validation, Middleware (15 minutes)

```go
func main() {
    cli, err := v2.NewCLI[AppConfig]("myapp", "My application", AppConfig{
        Verbose: false,
    },
        v2.WithCLIVersion[AppConfig]("1.2.3"),
        v2.WithEnvPrefix[AppConfig]("MYAPP_"),
        v2.WithSignalHandling[AppConfig](),
        v2.WithStrictValidation[AppConfig](),
        v2.WithMiddleware[AppConfig](v2.TimingMiddleware[AppConfig]()),
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Add a version command.
    v2.AddCommand(cli, v2.MustVersionCommand[AppConfig](cli))

    // Parent command with subcommands.
    listCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("list", listHandler,
        v2.WithShort[AppConfig, v2.NoFlags]("List resources"),
    )
    createCmd, _ := v2.NewCommand[AppConfig, v2.NoFlags]("create", createHandler,
        v2.WithShort[AppConfig, v2.NoFlags]("Create a resource"),
        v2.WithExactArgs[AppConfig, v2.NoFlags](1),
    )

    resourceCmd, err := v2.NewParentCommand[AppConfig, v2.NoFlags]("resource",
        "Resource management",
        []v2.Command[AppConfig, v2.NoFlags]{listCmd, createCmd},
        v2.WithShort[AppConfig, v2.NoFlags]("Resource management"),
    )
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    v2.AddCommand(cli, resourceCmd)
    cli.ExecuteAndExit(context.Background())
}
```

---

## Side-by-Side Comparison

| Task                   | Plain Cobra                                         | cmdguard                                                            |
| ---------------------- | --------------------------------------------------- | ------------------------------------------------------------------- |
| **Define a flag**      | `cmd.Flags().StringP("name", "n", "World", "Name")` | `Name string \`flag:"name" short:"n" default:"World" help:"Name"\`` |
| **Read a flag**        | `cmd.Flags().GetString("name")`                     | `flags.Name` (typed)                                                |
| **Required flag**      | `cmd.MarkFlagRequired("name")`                      | `required:"true"` struct tag                                        |
| **Add subcommand**     | `parent.AddCommand(child)`                          | `NewParentCommand(..., []Command{child})`                           |
| **Pre-run validation** | `PreRunE` func                                      | `WithPreRunE[T, F](fn)`                                             |
| **Shared services**    | Manual globals / closures                           | `Provide`/`Invoke` in DI scope                                      |
| **Env vars**           | Manual `os.Getenv`                                  | `env:"VAR_NAME"` struct tag                                         |
| **User typos**         | "unknown flag"                                      | "unknown flag: did you mean --name?"                                |
| **Version command**    | Write your own                                      | `MustVersionCommand(cli)`                                           |
| **Exit codes**         | `os.Exit(1)`                                        | `NewExitError(code, err)`                                           |

---

## Common Patterns

### Gradual Migration: Mixing Cobra and cmdguard Commands

You can add raw `*cobra.Command` to a cmdguard CLI at any time:

```go
// Legacy command written in plain Cobra.
legacyCmd := &cobra.Command{Use: "legacy", ...}

// New command written with cmdguard.
newCmd, _ := v2.NewCommand[AppConfig, *NewFlags]("new", ...)

// Both coexist on the same CLI.
cli.RootCommand().AddCommand(legacyCmd)
v2.AddCommand(cli, newCmd)
```

This lets you migrate command-by-command without a big-bang rewrite.

### Preserving Cobra Hooks

cmdguard commands support the same lifecycle hooks as Cobra:

| Cobra Hook          | cmdguard Equivalent                |
| ------------------- | ---------------------------------- |
| `PreRunE`           | `WithPreRunE[T, F](fn)`            |
| `PostRunE`          | `WithPostRunE[T, F](fn)`           |
| `PersistentPreRunE` | Use middleware or register on root |
| `RunE`              | Handler passed to `NewCommand`     |

### Converting Global/Persistent Flags

Move persistent flags to your root `AppConfig` struct. They're automatically available to all commands:

```go
type AppConfig struct {
    Verbose bool   `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
    Output  string `flag:"output"  short:"o" default:"text"    help:"Output format"`
}

// All commands receive *AppConfig automatically.
```

---

## Troubleshooting

### "unknown flag" errors after migration

Make sure you passed `WithFlags[T, F](&MyFlags{})` to `NewCommand`. Without it, cmdguard doesn't know to register flags for that command.

### "command has no handler" error

`NewCommand` requires a handler function. For parent commands (commands with subcommands), use `NewParentCommand` instead.

### Cobra commands lose access to cmdguard DI

Raw `*cobra.Command` added via `cli.RootCommand().AddCommand()` don't receive cmdguard's DI automatically. Either migrate them to `Command[T, F]` or manually invoke services from the scope in their `RunE` functions.

---

## Next Steps

- Explore the [`examples/`](../examples/) directory for complete working programs
- Read the [Quick Start Guide](QUICKSTART.md) for a full API tour
- Review [CLI Design Principles](CLI_DESIGN_PRINCIPLES.md) for best practices
