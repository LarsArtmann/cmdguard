package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	output "github.com/larsartmann/go-output"
	"github.com/spf13/cobra"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

//nolint:gocyclo // example file: linear command registration reads best as a single function
func buildCommands(cli *v3.CLI[AppConfig]) error {
	dbActionHandler := func(verb string) func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
		return func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("%s %s (force=%v)\n", verb, flags.Env, flags.Force)

			return nil
		}
	}
	printAndNil := func(msg string) func(_ context.Context, _ *AppConfig, _ v3.NoFlags) error {
		return func(_ context.Context, _ *AppConfig, _ v3.NoFlags) error {
			fmt.Println(msg)

			return nil
		}
	}
	scope := cli.Scope()

	// --- list: multi-format output, aliases, filter flags ---
	listCmd, err := v3.NewCommand(
		"list",
		&ListFlags{},
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

			format, err := output.ParseFormat(flags.Format)
			if err != nil {
				return v3.NewFlagError("format", err)
			}

			headers := []string{"ID", "Title", "Priority", "Status", "Created"}
			rows := make([][]string, 0, len(tasks))
			for _, t := range tasks {
				rows = append(rows, t.Row())
			}

			return v3.OutputTable(format, headers, rows)
		},
		v3.WithShort("List tasks"),
		v3.WithLong(`# List Tasks

Display all tasks with optional filtering by priority and completion status.

Supports multiple **output formats** for scripting and automation:

- `+"`table`"+` (default) — human-readable columnar output
- `+"`json`"+` — structured JSON for piping into `+"`jq`"+`
- `+"`csv`"+` — comma-separated for spreadsheet import
- `+"`yaml`"+` — structured YAML for config files`),
		v3.WithExample("taskctl list --format json --all"),
		v3.WithAliases("ls"),
		v3.WithGroupID("tasks"),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, listCmd); err != nil {
		return err
	}

	// --- add: required flags, values tag, PreRunE validation ---
	addCmd, err := v3.NewCommand(
		"add",
		&AddFlags{},
		func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task := store.Add(flags.Title, flags.Priority)
			fmt.Printf("Created task #%d: %s [%s]\n", task.ID, task.Title, task.Priority)
			return nil
		},
		v3.WithShort("Add a new task"),
		v3.WithLong(`# Add Task

Create a new task with a **title** and *priority*.

The `+"`--title`"+` flag is **required** and will prompt interactively if omitted.

Priority must be one of: `+"`low`"+`, `+"`medium`"+`, `+"`high`"+` (default: `+"`medium`"+`).`),
		v3.WithExample("taskctl add --title \"Fix bug\" --priority high"),
		v3.WithPreRunE(func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			if _, err := v3.ParseEnum(flags.Priority, strings.Split(allowedPriorities, ",")); err != nil {
				return v3.NewFlagError("priority", err)
			}
			return nil
		}),
		v3.WithGroupID("tasks"),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, addCmd); err != nil {
		return err
	}

	// --- done: exit codes, PostRunE cleanup, dynamic completion ---
	doneCmd, err := v3.NewCommand(
		"done",
		&DoneFlags{},
		func(_ context.Context, _ *AppConfig, flags *DoneFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task, err := store.Done(flags.ID)
			if err != nil {
				exitErr, _ := v3.NewExitError(2, v3.NewCommandError("done", err))
				return exitErr
			}

			fmt.Printf("Completed task #%d: %s\n", task.ID, task.Title)
			return nil
		},
		v3.WithShort("Mark a task as done"),
		v3.WithExample("taskctl done --id 1"),
		v3.WithPostRunE(func(_ context.Context, _ *AppConfig, _ *DoneFlags) error {
			fmt.Println("[cleanup] syncing state")
			return nil
		}),
		v3.WithGroupID("tasks"),
		v3.WithCompletion(func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			store, err := v3.Invoke[*TaskStore](scope)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return store.IDs(), cobra.ShellCompDirectiveNoFileComp
		}),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, doneCmd); err != nil {
		return err
	}

	// --- stats: OutputTable for terminal output ---
	statsCmd, err := v3.NewCommand(
		"stats",
		&StatsFlags{},
		func(_ context.Context, _ *AppConfig, _ *StatsFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			total, pending, done, byPriority := store.Stats()

			return v3.OutputTable(
				output.FormatTable,
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
		v3.WithShort("Show task statistics"),
		v3.WithExample("taskctl stats --format json"),
		v3.WithGroupID("tasks"),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, statsCmd); err != nil {
		return err
	}

	// --- inspect: ExactArgs, BranchingFlowContext, real task lookup ---
	inspectCmd, err := v3.NewCommand(
		"inspect",
		&InspectFlags{},
		func(ctx context.Context, _ *AppConfig, flags *InspectFlags) error {
			bfc, ok := v3.GetBranchingFlowContext(ctx)
			if ok {
				fmt.Printf("Flow path: %s\n", bfc.PathString())
			}

			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task, found := store.Get(1)
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
		v3.WithShort("Inspect a task in detail"),
		v3.WithExample("taskctl inspect 1"),
		v3.WithExactArgs(1),
		v3.WithValidArgs("1", "2", "3"),
		v3.WithGroupID("tasks"),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, inspectCmd); err != nil {
		return err
	}

	// --- db: parent command with shared DBFlags ---
	migrateCmd, err := v3.NewCommand(
		"migrate",
		&DBFlags{},
		dbActionHandler("Running migrations on"),
		v3.WithShort("Run database migrations"),
		v3.WithLong(`# Database Migrations

Run pending database migrations against the configured environment.

Use `+"`--force`"+` to skip confirmation prompts in **CI/CD** pipelines.`),
	)
	if err != nil {
		return err
	}

	seedCmd, err := v3.NewCommand(
		"seed",
		&DBFlags{},
		dbActionHandler("Seeding"),
		v3.WithShort("Seed the database"),
	)
	if err != nil {
		return err
	}

	dbStatusCmd, err := v3.NewCommand(
		"status",
		&DBFlags{},
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("DB status on %s: connected\n", flags.Env)
			return nil
		},
		v3.WithShort("Check database status"),
	)
	if err != nil {
		return err
	}

	dbCmd, err := v3.NewParentCommand[AppConfig](
		"db",
		"Database operations",
		&DBFlags{},
		v3.WithSubcommands(migrateCmd, seedCmd, dbStatusCmd),
		v3.WithShort("Database operations"),
		v3.WithGroupID("system"),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, dbCmd); err != nil {
		return err
	}

	// --- doctor: DoctorCommand, HealthCheckResultsWithContext ---
	doctorCmd, err := v3.DoctorCommand[AppConfig](
		cli,
		v3.WithDoctorGroupID[AppConfig]("system"),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, doctorCmd); err != nil {
		return err
	}

	// --- config show: demonstrate env/config resolution ---
	configShowCmd, err := v3.NewCommand(
		"show",
		v3.NoFlags{},
		func(_ context.Context, cfg *AppConfig, _ v3.NoFlags) error {
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
		v3.WithShort("Show resolved configuration"),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}

	configEditCmd, err := v3.NewCommand(
		"edit",
		v3.NoFlags{},
		func(_ context.Context, cfg *AppConfig, _ v3.NoFlags) error {
			fmt.Printf("Edit config manually at ~/.config/taskctl/config.json\n")
			fmt.Printf("  LogLevel: %s\n", cfg.LogLevel)
			fmt.Printf("  DataDir:  %s\n", cfg.DataDir)
			return nil
		},
		v3.WithShort("Show config location"),
		v3.WithNoArgs(),
	)
	if err != nil {
		return err
	}

	configCmd, err := v3.NewParentCommand[AppConfig](
		"config",
		"Configuration management",
		v3.NoFlags{},
		v3.WithSubcommands(configShowCmd, configEditCmd),
		v3.WithShort("Configuration management"),
		v3.WithGroupID("system"),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, configCmd); err != nil {
		return err
	}

	// --- version: VersionCommand ---
	versionCmd, err := v3.VersionCommand[AppConfig](cli)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, versionCmd); err != nil {
		return err
	}

	// --- hidden command ---
	hiddenCmd, err := v3.NewCommand(
		"secret",
		v3.NoFlags{},
		printAndNil("You found the secret command!"),
		v3.WithShort("Secret command"),
		v3.WithHidden(true),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, hiddenCmd); err != nil {
		return err
	}

	// --- deprecated command ---
	deprecatedCmd, err := v3.NewCommand(
		"complete",
		v3.NoFlags{},
		printAndNil("Use 'done' instead."),
		v3.WithShort("Deprecated: use done"),
		v3.WithDeprecated("Use 'done' instead"),
	)
	if err != nil {
		return err
	}
	if err := v3.AddCommand(cli, deprecatedCmd); err != nil {
		return err
	}

	return nil
}

func resolveStore(scope *v3.Scope) (*TaskStore, error) {
	store, err := v3.Invoke[*TaskStore](scope)
	if err != nil {
		return nil, v3.NewCommandError("task", fmt.Errorf("resolve store: %w", err))
	}
	return store, nil
}

func taskStatusLabel(done bool) string {
	if done {
		return taskStatusDone
	}
	return taskStatusPending
}

func seedTasks(cli *v3.CLI[AppConfig]) {
	store, err := v3.Invoke[*TaskStore](cli.Scope())
	if err != nil {
		return
	}
	store.Add("Set up CI pipeline", PriorityHigh)
	store.Add("Write API documentation", PriorityMedium)
	store.Add("Review pull request #42", PriorityLow)
}
