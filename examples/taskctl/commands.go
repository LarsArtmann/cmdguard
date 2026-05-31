package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

//nolint:gocyclo // example file: linear command registration reads best as a single function
func buildCommands(cli *v2.CLI[AppConfig]) error {
	scope := cli.Scope()

	// --- list: multi-format output, aliases, filter flags ---
	listCmd, err := v2.NewCommand[AppConfig, *ListFlags](
		"list",
		func(_ context.Context, _ *AppConfig, flags *ListFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			tasks := store.List(flags.Priority, flags.All)
			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return nil
			}

			format, err := v2.ParseOutputFormat(flags.Format)
			if err != nil {
				return v2.NewFlagError("format", err)
			}

			headers := []string{"ID", "Title", "Priority", "Status", "Created"}
			rows := make([][]string, len(tasks))
			for i, t := range tasks {
				rows[i] = t.Row()
			}

			return v2.OutputTable(format, headers, rows)
		},
		v2.WithShort[AppConfig, *ListFlags]("List tasks"),
		v2.WithLong[AppConfig, *ListFlags](`# List Tasks

Display all tasks with optional filtering by priority and completion status.

Supports multiple **output formats** for scripting and automation:

- `+"`table`"+` (default) — human-readable columnar output
- `+"`json`"+` — structured JSON for piping into `+"`jq`"+`
- `+"`csv`"+` — comma-separated for spreadsheet import
- `+"`yaml`"+` — structured YAML for config files`),
		v2.WithExample[AppConfig, *ListFlags]("taskctl list --format json --all"),
		v2.WithAliases[AppConfig, *ListFlags]("ls"),
		v2.WithFlags[AppConfig, *ListFlags](&ListFlags{}),
		v2.WithGroupID[AppConfig, *ListFlags]("tasks"),
		v2.WithNoArgs[AppConfig, *ListFlags](),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, listCmd); err != nil {
		return err
	}

	// --- add: required flags, values tag, PreRunE validation ---
	addCmd, err := v2.NewCommand[AppConfig, *AddFlags](
		"add",
		func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task := store.Add(flags.Title, flags.Priority)
			fmt.Printf("Created task #%d: %s [%s]\n", task.ID, task.Title, task.Priority)
			return nil
		},
		v2.WithShort[AppConfig, *AddFlags]("Add a new task"),
		v2.WithLong[AppConfig, *AddFlags](`# Add Task

Create a new task with a **title** and *priority*.

The `+"`--title`"+` flag is **required** and will prompt interactively if omitted.

Priority must be one of: `+"`low`"+`, `+"`medium`"+`, `+"`high`"+` (default: `+"`medium`"+`).`),
		v2.WithExample[AppConfig, *AddFlags]("taskctl add --title \"Fix bug\" --priority high"),
		v2.WithFlags[AppConfig, *AddFlags](&AddFlags{}),
		v2.WithPreRunE[AppConfig, *AddFlags](
			func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
				if _, err := v2.ParseEnum(flags.Priority, strings.Split(allowedPriorities, ",")); err != nil {
					return v2.NewFlagError("priority", err)
				}
				return nil
			},
		),
		v2.WithGroupID[AppConfig, *AddFlags]("tasks"),
		v2.WithNoArgs[AppConfig, *AddFlags](),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, addCmd); err != nil {
		return err
	}

	// --- done: exit codes, PostRunE cleanup, dynamic completion ---
	doneCmd, err := v2.NewCommand[AppConfig, *DoneFlags](
		"done",
		func(_ context.Context, _ *AppConfig, flags *DoneFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task, err := store.Done(flags.ID)
			if err != nil {
				exitErr, _ := v2.NewExitError(2, v2.NewCommandError("done", err))
				return exitErr
			}

			fmt.Printf("Completed task #%d: %s\n", task.ID, task.Title)
			return nil
		},
		v2.WithShort[AppConfig, *DoneFlags]("Mark a task as done"),
		v2.WithExample[AppConfig, *DoneFlags]("taskctl done --id 1"),
		v2.WithFlags[AppConfig, *DoneFlags](&DoneFlags{}),
		v2.WithPostRunE[AppConfig, *DoneFlags](
			func(_ context.Context, _ *AppConfig, _ *DoneFlags) error {
				fmt.Println("[cleanup] syncing state")
				return nil
			},
		),
		v2.WithGroupID[AppConfig, *DoneFlags]("tasks"),
		v2.WithCompletion[AppConfig, *DoneFlags](
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				store, err := v2.Invoke[*TaskStore](scope)
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				return store.IDs(), cobra.ShellCompDirectiveNoFileComp
			},
		),
		v2.WithNoArgs[AppConfig, *DoneFlags](),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, doneCmd); err != nil {
		return err
	}

	// --- stats: OutputStyledTable for terminal-pretty output ---
	statsCmd, err := v2.NewCommand[AppConfig, *StatsFlags](
		"stats",
		func(_ context.Context, _ *AppConfig, _ *StatsFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			total, pending, done, byPriority := store.Stats()

			return v2.OutputStyledTable(
				[]string{"Metric", "Value"},
				[][]string{
					{"Total", strconv.Itoa(total)},
					{"Pending", strconv.Itoa(pending)},
					{"Done", strconv.Itoa(done)},
					{"High Priority", strconv.Itoa(byPriority[PriorityHigh])},
					{"Medium Priority", strconv.Itoa(byPriority[PriorityMedium])},
					{"Low Priority", strconv.Itoa(byPriority[PriorityLow])},
				},
			)
		},
		v2.WithShort[AppConfig, *StatsFlags]("Show task statistics"),
		v2.WithFlags[AppConfig, *StatsFlags](&StatsFlags{}),
		v2.WithGroupID[AppConfig, *StatsFlags]("tasks"),
		v2.WithNoArgs[AppConfig, *StatsFlags](),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, statsCmd); err != nil {
		return err
	}

	// --- inspect: ExactArgs, BranchingFlowContext, real task lookup ---
	inspectCmd, err := v2.NewCommand[AppConfig, *InspectFlags](
		"inspect",
		func(ctx context.Context, _ *AppConfig, flags *InspectFlags) error {
			bfc, ok := v2.GetBranchingFlowContext(ctx)
			if ok {
				fmt.Printf("Flow path: %s\n", bfc.PathString())
			}

			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			// Positional arg (task ID) is validated by WithExactArgs(1)
			// but cobra doesn't pass it to our RunE. We look up via store.Get
			// with the ID from the first valid arg.
			task, found := store.Get(1) // demo: always shows task #1
			if !found {
				fmt.Println("No task found at ID")
				return nil
			}

			fmt.Printf("Task #%d: %s\n", task.ID, task.Title)
			fmt.Printf("  Priority: %s\n", task.Priority)
			fmt.Printf("  Status:   %s\n", taskStatusLabel(task.Done))
			fmt.Printf("  Created:  %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))

			if flags.ShowMetadata {
				fmt.Printf("  ID (raw): %d\n", task.ID)
			}

			return nil
		},
		v2.WithShort[AppConfig, *InspectFlags]("Inspect a task in detail"),
		v2.WithExample[AppConfig, *InspectFlags]("taskctl inspect 1"),
		v2.WithFlags[AppConfig, *InspectFlags](&InspectFlags{}),
		v2.WithExactArgs[AppConfig, *InspectFlags](1),
		v2.WithValidArgs[AppConfig, *InspectFlags]("1", "2", "3"),
		v2.WithGroupID[AppConfig, *InspectFlags]("tasks"),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, inspectCmd); err != nil {
		return err
	}

	// --- db: NewParentCommand with shared DBFlags ---
	migrateCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"migrate",
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("Running migrations on %s (force=%v)\n", flags.Env, flags.Force)
			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Run database migrations"),
		v2.WithLong[AppConfig, *DBFlags](`# Database Migrations

Run pending database migrations against the configured environment.

Use `+"`--force`"+` to skip confirmation prompts in **CI/CD** pipelines.`),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		return err
	}

	seedCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"seed",
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("Seeding %s (force=%v)\n", flags.Env, flags.Force)
			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Seed the database"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		return err
	}

	dbStatusCmd, err := v2.NewCommand[AppConfig, *DBFlags](
		"status",
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("DB status on %s: connected\n", flags.Env)
			return nil
		},
		v2.WithShort[AppConfig, *DBFlags]("Check database status"),
		v2.WithFlags[AppConfig, *DBFlags](&DBFlags{}),
	)
	if err != nil {
		return err
	}

	dbCmd, err := v2.NewParentCommand[AppConfig, *DBFlags](
		"db",
		"Database operations",
		[]v2.Command[AppConfig, *DBFlags]{migrateCmd, seedCmd, dbStatusCmd},
		v2.WithShort[AppConfig, *DBFlags]("Database operations"),
		v2.WithGroupID[AppConfig, *DBFlags]("system"),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, dbCmd); err != nil {
		return err
	}

	// --- health: HealthCheckWithContext, MustNewCommand ---
	healthCmd := v2.MustNewCommand[AppConfig, v2.NoFlags](
		"health",
		func(ctx context.Context, _ *AppConfig, _ v2.NoFlags) error {
			if err := cli.HealthCheckWithContext(ctx); err != nil {
				fmt.Printf("UNHEALTHY: %v\n", err)
				exitErr, _ := v2.NewExitError(1, err)
				return exitErr
			}
			fmt.Println("All systems healthy")
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Check system health"),
		v2.WithGroupID[AppConfig, v2.NoFlags]("system"),
	)
	if err := v2.AddCommand(cli, healthCmd); err != nil {
		return err
	}

	// --- config show: demonstrate env/config resolution ---
	configShowCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"show",
		func(_ context.Context, cfg *AppConfig, _ v2.NoFlags) error {
			fmt.Println("Current configuration:")
			fmt.Printf("  LogLevel:  %s\n", cfg.LogLevel)
			fmt.Printf("  DataDir:     %s\n", cfg.DataDir)
			fmt.Printf("  Timeout:     %s\n", cfg.Timeout)
			fmt.Printf("  Port:        %s\n", cfg.Port)
			fmt.Printf("  AdminEmail:  %s\n", cfg.AdminEmail)
			fmt.Printf("  APIUrl:      %s\n", cfg.APIUrl)
			fmt.Printf("  Verbose:     %d\n", cfg.Verbose)
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Show resolved configuration"),
		v2.WithNoArgs[AppConfig, v2.NoFlags](),
	)
	if err != nil {
		return err
	}

	configEditCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"edit",
		func(ctx context.Context, cfg *AppConfig, _ v2.NoFlags) error {
			content := fmt.Sprintf("log-level: %s\ndata-dir: %s\n", cfg.LogLevel, cfg.DataDir)
			edited, err := v2.EditInEditor(ctx, content)
			if err != nil {
				return fmt.Errorf("editor: %w", err)
			}
			fmt.Printf("Edited config:\n%s", edited)
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Edit config in $EDITOR"),
		v2.WithNoArgs[AppConfig, v2.NoFlags](),
	)
	if err != nil {
		return err
	}

	configCmd, err := v2.NewParentCommand[AppConfig, v2.NoFlags](
		"config",
		"Configuration management",
		[]v2.Command[AppConfig, v2.NoFlags]{configShowCmd, configEditCmd},
		v2.WithShort[AppConfig, v2.NoFlags]("Configuration management"),
		v2.WithGroupID[AppConfig, v2.NoFlags]("system"),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, configCmd); err != nil {
		return err
	}

	// --- version: MustVersionCommand ---
	versionCmd := v2.MustVersionCommand[AppConfig](cli)
	if err := v2.AddCommand(cli, versionCmd); err != nil {
		return err
	}

	// --- hidden command ---
	hiddenCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"secret",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			fmt.Println("You found the secret command!")
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Secret command"),
		v2.WithHidden[AppConfig, v2.NoFlags](true),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, hiddenCmd); err != nil {
		return err
	}

	// --- deprecated command ---
	deprecatedCmd, err := v2.NewCommand[AppConfig, v2.NoFlags](
		"complete",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			fmt.Println("Use 'done' instead.")
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Deprecated: use done"),
		v2.WithDeprecated[AppConfig, v2.NoFlags]("Use 'done' instead"),
	)
	if err != nil {
		return err
	}
	if err := v2.AddCommand(cli, deprecatedCmd); err != nil {
		return err
	}

	return nil
}

func resolveStore(scope *v2.Scope) (*TaskStore, error) {
	store, err := v2.Invoke[*TaskStore](scope)
	if err != nil {
		return nil, v2.NewCommandError("task", fmt.Errorf("resolve store: %w", err))
	}
	return store, nil
}

func taskStatusLabel(done bool) string {
	if done {
		return taskStatusDone
	}
	return taskStatusPending
}

func seedTasks(cli *v2.CLI[AppConfig]) {
	store, err := v2.Invoke[*TaskStore](cli.Scope())
	if err != nil {
		return
	}
	store.Add("Set up CI pipeline", PriorityHigh)
	store.Add("Write API documentation", PriorityMedium)
	store.Add("Review pull request #42", PriorityLow)
}
