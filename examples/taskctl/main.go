// Package main demonstrates every major cmdguard feature in a production-grade task manager.
//
// Features shown:
//   - Type-safe config with env var bindings, counting flags, and validated types (Email, URL, LogLevel, Duration, Port)
//   - Dependency injection with lifecycle hooks (HealthCheck, Shutdown)
//   - Per-command typed flags with prompt, required, validate, and values tags
//   - PreRunE validation and PostRunE cleanup
//   - Middleware (spinner + timing + recovery)
//   - Glamour markdown help rendering (glamour.WithHelpTheme sub-module)
//   - Rich output in multiple formats (OutputTable, OutputResult)
//   - Command groups (WithGroup)
//   - Subcommands via NewParentCommand
//   - Error handling with typed errors and exit codes
//   - Graceful shutdown with DI service cleanup (WithGracefulShutdown)
//   - Config file loading (JSON)
//   - Shell completion via WithCompletion
//   - Hidden and deprecated commands
//   - Command aliases
//   - Arg validators (WithNoArgs, WithExactArgs)
//   - BranchingFlowContext for path tracking
//   - Version command
//   - DI audit logging via samber-do-auditlog (WithAuditLog + plugin accessor pattern)
//
// Usage:
//
//	go run examples/taskctl/main.go list
//	go run examples/taskctl/main.go list --format json --all
//	go run examples/taskctl/main.go add --title "Buy groceries" --priority high
//	go run examples/taskctl/main.go done --id 1
//	go run examples/taskctl/main.go stats
//	go run examples/taskctl/main.go doctor
//	go run examples/taskctl/main.go version
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	auditlog "github.com/larsartmann/samber-do-auditlog"

	"github.com/larsartmann/cmdguard/glamour"
	"github.com/larsartmann/cmdguard/spinner"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func main() {
	ctx := context.Background()

	// Audit logging — captures DI lifecycle events for observability
	// Set DO_AUDITLOG_ENABLED=true to enable without changing code.
	auditPlugin, err := auditlog.New(auditlog.Config{
		Enabled:     true,
		ContainerID: "taskctl",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating audit log plugin: %v\n", err)
		os.Exit(1)
	}

	cli, err := v4.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v4.WithCLIVersion("1.0.0"),
		v4.WithEnvPrefix("TASKCTL_"),
		v4.WithAuditLog(auditPlugin),
		v4.WithConfigFile("$HOME/.config/taskctl/config.json"),
		v4.WithConfigValidation(func(cfg *AppConfig) error {
			if cfg.DataDir == "" {
				return fmt.Errorf("data-dir must not be empty")
			}
			return nil
		}),
		v4.WithGracefulShutdown(),
		v4.WithStrictValidation(),
		v4.WithMiddleware(
			spinner.Middleware[AppConfig]("Working..."),
			v4.TimingMiddleware[AppConfig](func(name string, d time.Duration, err error) {
				fmt.Fprintf(os.Stderr, "[timing] %s took %v (err=%v)\n", name, d, err)
			}),
			v4.RecoveryMiddleware[AppConfig](),
		),
		glamour.WithHelpTheme("dark"),
		v4.WithGroup("tasks", "Task Management"),
		v4.WithGroup("system", "System"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Global flags available to all commands
	cli.AddGlobalBoolFlag("debug", "D", false, "Enable debug mode")

	// Register DI services
	if err := v4.Provide(cli.Scope(), NewTaskStore); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Seed demo data
	seedTasks(cli)

	// Build all commands
	if err := buildCommands(cli); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// cmdguard prints command errors exactly once (via fang by default). The
	// error returned by Execute is used only to map the process exit code —
	// re-printing it would duplicate the error output on stderr.
	execErr := cli.Execute(ctx)

	if plugin := cli.AuditLog(); plugin != nil && plugin.EventsCount() > 0 {
		// AUDIT_LOG_FORMAT selects the export format: html, json, ndjson,
		// csv, tsv, mermaid, dot, d2, plantuml, tree, or htmltree.
		// Defaults to html.
		format, err := v4.ParseAuditLogFormat(os.Getenv("AUDIT_LOG_FORMAT"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "audit-log format invalid: %v\n", err)
		} else {
			path := "taskctl-audit." + format.String()
			if err := v4.ExportAuditLog(cli, v4.AuditLogExportConfig{
				Format: format,
				Path:   path,
			}); err != nil {
				fmt.Fprintf(os.Stderr, "audit-log export failed: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "audit-log written to %s\n", path)
			}
		}
	}

	os.Exit(v4.ExitCode(execErr))
}
