package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	output "github.com/larsartmann/go-output"
	"github.com/spf13/cobra"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func buildCommands(cli *v5.CLI[AppConfig]) error {
	dbActionHandler := func(verb string) func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
		return func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("%s %s (force=%v)\n", verb, flags.Env, flags.Force)

			return nil
		}
	}
	printAndNil := func(msg string) func(_ context.Context, _ *AppConfig, _ v5.NoFlags) error {
		return func(_ context.Context, _ *AppConfig, _ v5.NoFlags) error {
			fmt.Println(msg)

			return nil
		}
	}
	scope := cli.Scope()
	b := v5.NewBuilder(cli)

	// --- list: multi-format output, aliases, filter flags ---
	v5.Add(b,
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
				return v5.NewFlagError("format", err)
			}

			headers := []string{"ID", "Title", "Priority", "Status", "Created"}
			rows := make([][]string, 0, len(tasks))
			for _, t := range tasks {
				rows = append(rows, t.Row())
			}

			return v5.OutputTable(format, headers, rows)
		},
		v5.WithShort("List tasks"),
		v5.WithLong(`# List Tasks

Display all tasks with optional filtering by priority and completion status.

Supports multiple **output formats** for scripting and automation:

- `+"`table`"+` (default) — human-readable columnar output
- `+"`json`"+` — structured JSON for piping into `+"`jq`"+`
- `+"`csv`"+` — comma-separated for spreadsheet import
- `+"`yaml`"+` — structured YAML for config files`),
		v5.WithExample("taskctl list --format json --all"),
		v5.WithAliases("ls"),
		v5.WithGroupID("tasks"),
		v5.WithNoArgs(),
	)

	// --- add: required flags, values tag, PreRunE validation ---
	v5.Add(b,
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
		v5.WithShort("Add a new task"),
		v5.WithLong(`# Add Task

Create a new task with a **title** and *priority*.

The `+"`--title`"+` flag is **required** and will prompt interactively if omitted.

Priority must be one of: `+"`low`"+`, `+"`medium`"+`, `+"`high`"+` (default: `+"`medium`"+`).`),
		v5.WithExample("taskctl add --title \"Fix bug\" --priority high"),
		v5.WithPreRunE(func(_ context.Context, _ *AppConfig, flags *AddFlags) error {
			if _, err := v5.ParseEnum(flags.Priority, strings.Split(allowedPriorities, ",")); err != nil {
				return v5.NewFlagError("priority", err)
			}
			return nil
		}),
		v5.WithGroupID("tasks"),
		v5.WithNoArgs(),
	)

	// --- done: exit codes, PostRunE cleanup, dynamic completion ---
	v5.Add(b,
		"done",
		&DoneFlags{},
		func(_ context.Context, _ *AppConfig, flags *DoneFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			task, err := store.Done(flags.ID)
			if err != nil {
				exitErr, _ := v5.NewExitError(2, v5.NewCommandError("done", err))
				return exitErr
			}

			fmt.Printf("Completed task #%d: %s\n", task.ID, task.Title)
			return nil
		},
		v5.WithShort("Mark a task as done"),
		v5.WithExample("taskctl done --id 1"),
		v5.WithPostRunE(func(_ context.Context, _ *AppConfig, _ *DoneFlags) error {
			fmt.Println("[cleanup] syncing state")
			return nil
		}),
		v5.WithGroupID("tasks"),
		v5.WithCompletion(func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			store, err := v5.Invoke[*TaskStore](scope)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			return store.IDs(), cobra.ShellCompDirectiveNoFileComp
		}),
		v5.WithNoArgs(),
	)

	// --- stats: OutputTable for terminal output ---
	v5.Add(b,
		"stats",
		&StatsFlags{},
		func(_ context.Context, _ *AppConfig, _ *StatsFlags) error {
			store, err := resolveStore(scope)
			if err != nil {
				return err
			}

			total, pending, done, byPriority := store.Stats()

			return v5.OutputTable(
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
		v5.WithShort("Show task statistics"),
		v5.WithExample("taskctl stats --format json"),
		v5.WithGroupID("tasks"),
		v5.WithNoArgs(),
	)

	// --- inspect: ExactArgs, BranchingFlowContext, real task lookup ---
	v5.Add(b,
		"inspect",
		&InspectFlags{},
		func(ctx context.Context, _ *AppConfig, flags *InspectFlags) error {
			bfc, ok := v5.GetBranchingFlowContext(ctx)
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
		v5.WithShort("Inspect a task in detail"),
		v5.WithExample("taskctl inspect 1"),
		v5.WithExactArgs(1),
		v5.WithValidArgs("1", "2", "3"),
		v5.WithGroupID("tasks"),
	)

	// --- db: parent command with shared DBFlags ---
	migrateCmd, err := v5.NewCommand(
		"migrate",
		&DBFlags{},
		dbActionHandler("Running migrations on"),
		v5.WithShort("Run database migrations"),
		v5.WithLong(`# Database Migrations

Run pending database migrations against the configured environment.

Use `+"`--force`"+` to skip confirmation prompts in **CI/CD** pipelines.`),
	)
	b.Check(err)

	seedCmd, err := v5.NewCommand(
		"seed",
		&DBFlags{},
		dbActionHandler("Seeding"),
		v5.WithShort("Seed the database"),
	)
	b.Check(err)

	dbStatusCmd, err := v5.NewCommand(
		"status",
		&DBFlags{},
		func(_ context.Context, _ *AppConfig, flags *DBFlags) error {
			fmt.Printf("DB status on %s: connected\n", flags.Env)
			return nil
		},
		v5.WithShort("Check database status"),
	)
	b.Check(err)

	v5.AddParent(b,
		"db",
		"Database operations",
		&DBFlags{},
		v5.WithSubcommands(migrateCmd, seedCmd, dbStatusCmd),
		v5.WithShort("Database operations"),
		v5.WithGroupID("system"),
	)

	// --- doctor: DoctorCommand, HealthCheckResultsWithContext ---
	doctorCmd, err := v5.DoctorCommand[AppConfig](
		cli,
		v5.WithDoctorGroupID[AppConfig]("system"),
	)
	v5.Try(b, doctorCmd, err)

	// --- config show: demonstrate env/config resolution ---
	configShowCmd, err := v5.NewCommand(
		"show",
		v5.NoFlags{},
		func(_ context.Context, cfg *AppConfig, _ v5.NoFlags) error {
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
		v5.WithShort("Show resolved configuration"),
		v5.WithNoArgs(),
	)
	b.Check(err)

	configEditCmd, err := v5.NewCommand(
		"edit",
		v5.NoFlags{},
		func(_ context.Context, cfg *AppConfig, _ v5.NoFlags) error {
			fmt.Printf("Edit config manually at ~/.config/taskctl/config.json\n")
			fmt.Printf("  LogLevel: %s\n", cfg.LogLevel)
			fmt.Printf("  DataDir:  %s\n", cfg.DataDir)
			return nil
		},
		v5.WithShort("Show config location"),
		v5.WithNoArgs(),
	)
	b.Check(err)

	v5.AddParent(b,
		"config",
		"Configuration management",
		v5.NoFlags{},
		v5.WithSubcommands(configShowCmd, configEditCmd),
		v5.WithShort("Configuration management"),
		v5.WithGroupID("system"),
	)

	// --- version: VersionCommand ---
	versionCmd, err := v5.VersionCommand[AppConfig](cli)
	v5.Try(b, versionCmd, err)

	// --- hidden command ---
	v5.Add(b,
		"secret",
		v5.NoFlags{},
		printAndNil("You found the secret command!"),
		v5.WithShort("Secret command"),
		v5.WithHidden(true),
	)

	// --- deprecated command ---
	v5.Add(b,
		"complete",
		v5.NoFlags{},
		printAndNil("Use 'done' instead."),
		v5.WithShort("Deprecated: use done"),
		v5.WithDeprecated("Use 'done' instead"),
	)

	return b.Build()
}

func resolveStore(scope *v5.Scope) (*TaskStore, error) {
	store, err := v5.Invoke[*TaskStore](scope)
	if err != nil {
		return nil, v5.NewCommandError("task", fmt.Errorf("resolve store: %w", err))
	}
	return store, nil
}

func taskStatusLabel(done bool) string {
	if done {
		return taskStatusDone
	}
	return taskStatusPending
}

func seedTasks(cli *v5.CLI[AppConfig]) {
	store, err := v5.Invoke[*TaskStore](cli.Scope())
	if err != nil {
		return
	}
	store.Add("Set up CI pipeline", PriorityHigh)
	store.Add("Write API documentation", PriorityMedium)
	store.Add("Review pull request #42", PriorityLow)
}
