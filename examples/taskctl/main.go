// Package main demonstrates every major cmdguard feature in a production-grade task manager.
//
// Features shown:
//   - Type-safe config with env var bindings, counting flags, and validated types (Email, URL, LogLevel, Duration, Port)
//   - Dependency injection with lifecycle hooks (HealthCheck, Shutdown)
//   - Per-command typed flags with prompt, required, validate, and values tags
//   - PreRunE validation and PostRunE cleanup
//   - Middleware (spinner + timing + recovery)
//   - Glamour markdown help rendering (WithGlamourHelpTheme)
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
//   - EditInEditor for config editing
//   - Version command
//   - DI audit logging via samber-do-auditlog (WithAuditLog + AuditLogCommand)
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

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func main() {
	ctx := context.Background()

	// Audit logging — captures DI lifecycle events for observability
	// Set DO_AUDITLOG_ENABLED=true to enable without changing code.
	auditPlugin := auditlog.New(auditlog.Config{
		Enabled:     true,
		ContainerID: "taskctl",
	})

	cli, err := v2.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v2.WithCLIVersion[AppConfig]("1.0.0"),
		v2.WithEnvPrefix[AppConfig]("TASKCTL_"),
		v2.WithAuditLog[AppConfig](auditPlugin),
		v2.WithConfigFile[AppConfig]("$HOME/.config/taskctl/config.json"),
		v2.WithConfigValidation[AppConfig](func(cfg *AppConfig) error {
			if cfg.DataDir == "" {
				return fmt.Errorf("data-dir must not be empty")
			}
			return nil
		}),
		v2.WithGracefulShutdown[AppConfig](),
		v2.WithStrictValidation[AppConfig](),
		v2.WithMiddleware[AppConfig](
			v2.SpinnerMiddleware[AppConfig]("Working..."),
			v2.TimingMiddleware[AppConfig](func(name string, d time.Duration, err error) {
				fmt.Fprintf(os.Stderr, "[timing] %s took %v (err=%v)\n", name, d, err)
			}),
			v2.RecoveryMiddleware[AppConfig](),
		),
		v2.WithGlamourHelpTheme[AppConfig]("dark"),
		v2.WithGroup[AppConfig]("tasks", "Task Management"),
		v2.WithGroup[AppConfig]("system", "System"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Global flags available to all commands
	cli.AddGlobalBoolFlag("debug", "D", false, "Enable debug mode")

	// Register DI services
	if err := v2.Provide(cli.Scope(), NewTaskStore); err != nil {
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

	cli.ExecuteAndExit(ctx)
}
