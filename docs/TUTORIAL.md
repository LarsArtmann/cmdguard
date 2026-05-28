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
go get github.com/larsartmann/cmdguard
```

Create `main.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"

    v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type AppConfig struct {
    DataDir string `flag:"data-dir" short:"d" default:"./tasks" help:"Directory for task storage"`
}

func main() {
    cli, err := v2.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{})
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

Commands in cmdguard are created with `NewCommand` and have their own typed flags:

```go
type AddFlags struct {
    Title    string `flag:"title"    short:"t" required:"true"  help:"Task title"`
    Priority string `flag:"priority" short:"p" default:"medium" help:"low, medium, or high"`
}

addCmd, err := v2.NewCommand[AppConfig, *AddFlags]("add",
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        fmt.Printf("Adding task: %s [%s]\n", flags.Title, flags.Priority)
        return nil
    },
    v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
    v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
)
if err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}

if err := v2.AddCommand(cli, addCmd); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

**Key points:**

- `Command[AppConfig, *AddFlags]` — first type param is the shared config, second is per-command flags
- `required:"true"` makes the flag mandatory — cmdguard enforces this automatically
- Flags are defined as struct tags, not string lookups

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

cli, _ := v2.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{})

v2.Provide(cli.Scope(), func(i do.Injector) (*TaskStore, error) {
    fmt.Println("[store] initialized")
    return &TaskStore{tasks: []string{}}, nil
})
```

Now commands can access the store:

```go
addCmd, _ := v2.NewCommand[AppConfig, *AddFlags]("add",
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        store, err := v2.Invoke[*TaskStore](cli.Scope())
        if err != nil {
            return err
        }
        store.tasks = append(store.tasks, flags.Title)
        fmt.Printf("Added task #%d: %s\n", len(store.tasks), flags.Title)
        return nil
    },
    v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
    v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
)
```

---

## Step 4: Add a List Command with Rich Output

```go
type ListFlags struct {
    Format string `flag:"format" short:"f" default:"table" help:"Output format (table, json, csv, yaml)"`
}

listCmd, _ := v2.NewCommand[AppConfig, *ListFlags]("list",
    func(ctx context.Context, cfg *AppConfig, flags *ListFlags) error {
        store, _ := v2.Invoke[*TaskStore](cli.Scope())

        if len(store.tasks) == 0 {
            fmt.Println("No tasks found.")
            return nil
        }

        format, _ := v2.ParseOutputFormat(flags.Format)

        headers := []string{"#", "Title"}
        rows := make([][]string, len(store.tasks))
        for i, t := range store.tasks {
            rows[i] = []string{fmt.Sprintf("%d", i+1), t}
        }

        return v2.OutputTable(format, headers, rows)
    },
    v2.WithShort[AppConfig, *ListFlags]("List all tasks"),
    v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
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
addCmd, _ := v2.NewCommand[AppConfig, *AddFlags]("add",
    addHandler,
    v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
    v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
    v2.WithPreRunE[AppConfig, *AddFlags](
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

---

## Step 6: Add Cleanup with PostRunE

```go
v2.WithPostRunE[AppConfig, *AddFlags](
    func(ctx context.Context, cfg *AppConfig, flags *AddFlags) error {
        fmt.Println("[cleanup] flushing to disk")
        return nil
    },
),
```

`PostRunE` only fires on **success** — standard Cobra semantics.

---

## Step 7: Add Middleware

Middleware wraps every command handler:

```go
cli, _ := v2.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{},
    v2.WithMiddleware[AppConfig](
        v2.TimingMiddleware[AppConfig](func(name string, d time.Duration) {
            fmt.Fprintf(os.Stderr, "[timing] %s took %v\n", name, d)
        }),
        v2.RecoveryMiddleware[AppConfig](),
    ),
)
```

Now every command is timed and panic-safe.

---

## Step 8: Add Signal Handling

```go
cli, _ := v2.NewCLI[AppConfig]("taskctl", "A task manager CLI", AppConfig{},
    v2.WithSignalHandling[AppConfig](),
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
cli, _ := v2.NewCLI[AppConfig]("taskctl", "...", AppConfig{},
    v2.WithEnvPrefix[AppConfig]("TASKCTL_"),
)
```

Priority chain: **explicit flag → env var → default value**.

---

## Step 10: Add Error Handling with Exit Codes

Return typed errors with exit codes for CI/automation:

```go
doneCmd, _ := v2.NewCommand[AppConfig, *DoneFlags]("done",
    func(ctx context.Context, cfg *AppConfig, flags *DoneFlags) error {
        // task not found? Exit with code 2
        exitErr, _ := v2.NewExitError(2, fmt.Errorf("task %d not found", flags.ID))
        return exitErr
    },
    v2.WithShort[AppConfig, *DoneFlags]("Mark a task as done"),
    v2.WithFlags[AppConfig, *DoneFlags](&DoneFlags{}),
)
```

When using `cli.ExecuteAndExit(ctx)`, the exit code is propagated to `os.Exit`.

---

## Complete Example

See [`examples/kitchen-sink/`](../examples/kitchen-sink/) for a fully working production-grade CLI that uses every feature described above.

---

## Next Steps

- Read the [API Reference](https://pkg.go.dev/github.com/larsartmann/cmdguard/pkg/cmdguard/v2) for the full API
- See the [Quick Start Guide](QUICKSTART.md) for a complete feature tour
- Check [MIGRATION_FROM_COBRA.md](MIGRATION_FROM_COBRA.md) if you're migrating an existing Cobra app
- Browse [examples/](../examples/) for more working demos
