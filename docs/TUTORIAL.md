# Building a Task Manager CLI with cmdguard

> A step-by-step tutorial that builds a real CLI from scratch, covering every major cmdguard feature.
> **Time:** 30 minutes | **Difficulty:** Intermediate

---

## What You'll Build

A task manager CLI called `taskctl` with:

- Typed config with environment variable support
- Multiple commands with different flag types
- Dependency injection for shared state
- Pre-run validation and post-run cleanup
- Rich output in multiple formats
- Signal handling for graceful shutdown

---

## Step 1: Project Setup

Create a new Go module:

```bash
mkdir taskctl && cd taskctl
go mod init taskctl
go get github.com/larsartmann/cmdguard/v4
```

Create `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

type AppConfig struct {
    DataDir string `flag:"data-dir" short:"d" default:"./tasks" help:"Directory for task storage"`
}

func main() {
    cli, err := v4.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{})
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    if err := cli.Execute(context.Background()); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

Run it:

```bash
go run main.go --help
```

You now have a CLI with auto-generated help, version flag, and styled output.

---

## Step 2: Add Your First Command

Commands in cmdguard are created with `NewCommand` — flags are passed positionally and options are non-generic:

```go
type AddFlags struct {
    Title    string `flag:"title"    short:"t" required:"true"  help:"Task title"`
    Priority string `flag:"priority" short:"p" default:"medium" help:"low, medium, or high"`
}

addCmd, err := v4.NewCommand("add", &AddFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        fmt.Printf("Adding task: %s [%s]\n", flags.Title, flags.Priority)
        return nil
    },
    v4.WithShort("Add a new task"),
)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}

if err := v4.AddCommand(cli, addCmd); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

**Key points:**

- Flags are passed as the second argument to `NewCommand` — type parameters are inferred
- `required:"true"` makes the flag mandatory — cmdguard enforces this automatically
- Flags are defined as struct tags, not string lookups
- Options like `WithShort` are non-generic — no type parameters needed

Run it:

```bash
go run main.go add --title "Buy groceries" --priority high
# Output: Adding task: Buy groceries [high]

go run main.go add
# Error: required flag(s) "title" not set
```

---

## Step 3: Add Dependency Injection

Register a shared task store as a DI service:

```go
import "github.com/samber/do/v2"

type TaskStore struct {
    tasks []string
}

cli, _ := v4.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{})

v4.Provide(cli.Scope(), func(i do.Injector) (*TaskStore, error) {
    fmt.Println("[store] initialized")
    return &TaskStore{tasks: []string{}}, nil
})
```

Now commands can access the store:

```go
addCmd, _ := v4.NewCommand("add", &AddFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        store, err := v4.Invoke[*TaskStore](cli.Scope())
        if err != nil {
            return err
        }
        store.tasks = append(store.tasks, flags.Title)
        fmt.Printf("Added task #%d: %s\n", len(store.tasks), flags.Title)
        return nil
    },
    v4.WithShort("Add a new task"),
)
```

---

## Step 4: Add a List Command with Rich Output

```go
import "github.com/larsartmann/go-output"

type ListFlags struct {
    Format string `flag:"format" short:"f" default:"table" help:"Output format (table, json, csv, yaml)"`
}

listCmd, _ := v4.NewCommand("list", &ListFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *ListFlags) error {
        store, _ := v4.Invoke[*TaskStore](cli.Scope())

        if len(store.tasks) == 0 {
            fmt.Println("No tasks found.")
            return nil
        }

        format, _ := output.ParseFormat(flags.Format)

        headers := []string{"#", "Title"}
        rows := make([][]string, len(store.tasks))
        for i, t := range store.tasks {
            rows[i] = []string{fmt.Sprintf("%d", i+1), t}
        }

        return v4.OutputTable(format, headers, rows)
    },
    v4.WithShort("List all tasks"),
)
```

Try different formats:

```bash
go run main.go add --title "Write docs"
go run main.go list                # Pretty table
go run main.go list --format json  # JSON array
go run main.go list --format csv   # CSV
go run main.go list --format yaml  # YAML
```

---

## Step 5: Add Validation with PreRunE

Validate input before the handler runs:

```go
addCmd, _ := v4.NewCommand("add", &AddFlags{},
    addHandler,
    v4.WithShort("Add a new task"),
    v4.WithPreRunE[AppConfig, *AddFlags](
        func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
            if len(flags.Title) < 3 {
                return fmt.Errorf("title must be at least 3 characters")
            }

            switch flags.Priority {
            case "low", "medium", "high":
                // valid
            default:
                return fmt.Errorf("priority must be low, medium, or high, got %q", flags.Priority)
            }

            return nil
        },
    ),
)
```

`PreRunE` fires before `RunE`. If it returns an error, the handler never runs.

Note: `WithPreRunE` and `WithPostRunE` are generic functions that return a non-generic `CommandOption` — they're the only options that require type parameters.

---

## Step 6: Add Cleanup with PostRunE

```go
v4.WithPostRunE[AppConfig, *AddFlags](
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        fmt.Println("[cleanup] flushing to disk")
        return nil
    ),
```

`PostRunE` only fires on **success** — standard Cobra semantics.

---

## Step 7: Add Middleware

Middleware wraps every command handler:

```go
cli, _ := v4.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{},
    v4.WithMiddleware[AppConfig](
        v4.TimingMiddleware[AppConfig](func(name string, d time.Duration, err error) {
            fmt.Fprintf(os.Stderr, "[timing] %s took %v (err=%v)\n", name, d, err)
        }),
        v4.RecoveryMiddleware[AppConfig](),
    ),
)
```

Now every command is timed and panic-safe.

---

## Step 8: Add Signal Handling

```go
cli, _ := v4.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{},
    v4.WithSignalHandling(),
)
```

Ctrl+C now cancels the `context.Context` in all handlers. Use it for graceful cleanup:

```go
func handler(ctx context.Context, cfg *AppConfig, flags *Flags) error {
    select {
    case <-ctx.Done():
        return ctx.Err() // cancelled by signal
    default:
        // do work
    }
    return nil
}
```

---

## Step 9: Add Environment Variable Support

Config fields with `env` tags read from the environment automatically:

```go
type AppConfig struct {
    DataDir  string `flag:"data-dir"  env:"TASKCTL_DATA_DIR"  default:"./tasks" help:"Task storage directory"`
    LogLevel string `flag:"log-level" env:"TASKCTL_LOG_LEVEL" default:"info"     help:"Log level"`
}
```

Or use `WithEnvPrefix` to prefix all env lookups:

```go
cli, _ := v4.NewCLI[AppConfig]("taskctl", "...", AppConfig{},
    v4.WithEnvPrefix("TASKCTL_"),
)
```

Priority chain: **explicit flag → env var → default value**.

---

## Step 10: Add Error Handling with Exit Codes

Return typed errors with exit codes for CI/automation:

```go
doneCmd, _ := v4.NewCommand("done", &DoneFlags{},
    func(ctx context.Context, cfg *AppConfig, flags *DoneFlags) error {
        // task not found? Exit with code 2
        exitErr, _ := v4.NewExitError(2, fmt.Errorf("task %d not found", flags.ID))
        return exitErr
    },
    v4.WithShort("Mark a task as done"),
)
```

When using `cli.ExecuteAndExit(ctx)`, the exit code is propagated to `os.Exit`.

---

## Complete Example

See [`examples/taskctl/`](../examples/taskctl/) for a fully working production-grade CLI that uses every feature described above.

---

## Next Steps

- Read the [API Reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4) for the full API
- See the [Quick Start Guide](QUICKSTART.md) for a complete feature tour
- Check [MIGRATION_FROM_COBRA.md](MIGRATION_FROM_COBRA.md) if you're migrating an existing Cobra app
- Browse [examples/](../examples/) for more working demos
