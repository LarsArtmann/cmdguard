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

	"github.com/larsartmann/cmdguard/flightrecorder"
	"github.com/larsartmann/cmdguard/glamour"
	"github.com/larsartmann/cmdguard/spinner"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/larsartmann/cmdguard/v4/pkg/version"
)

func main() {
	os.Exit(v4.ExitCode(run(context.Background(), os.Args[1:])))
}

// run builds the full production CLI and executes args against it: audit
// logging, middleware, flight recorder (started lazily, stopped before
// return), and the post-run audit export. Split from main so tests can drive
// the exact production composition. Construction errors are printed here and
// returned; the execution error is printed by cmdguard exactly once and
// returned unprinted — main only maps it to the process exit code.
func run(ctx context.Context, args []string) error {
	rec := newProductionRecorder()

	cli, err := buildApp(rec) //nolint:contextcheck // DI construction is context-free by design; ctx arrives at ExecuteWithArgs below
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)

		return err
	}

	execErr := cli.ExecuteWithArgs(ctx, args)

	exportAuditLog(cli)

	rec.Stop()

	return execErr
}

// newProductionRecorder returns the production flight recorder: snapshots
// on commands slower than 5s or on errors. Split from main for testability.
func newProductionRecorder() *flightrecorder.Recorder {
	return flightrecorder.New(flightrecorder.Config{
		CaptureOnSlow:  true,
		SlowThreshold:  5 * time.Second,
		CaptureOnError: true,
	})
}

// buildApp composes the full production CLI: audit logging, config file
// loading, validation, graceful shutdown, middleware (spinner + timing +
// recovery), flight recorder, glamour help, and command groups. Split from
// main so tests can exercise the exact production composition; rec is started
// lazily by the middleware and stopped by the caller (see main).
func buildApp(rec *flightrecorder.Recorder) (*v4.CLI[AppConfig], error) {
	// Audit logging — captures DI lifecycle events for observability
	// Set DO_AUDITLOG_ENABLED=true to enable without changing code.
	auditPlugin, err := auditlog.New(auditlog.Config{
		Enabled:     true,
		ContainerID: "taskctl",
	})
	if err != nil {
		return nil, fmt.Errorf("creating audit log plugin: %w", err)
	}

	cli, err := v4.NewCLI[AppConfig](
		"taskctl", "A production-grade task manager CLI", AppConfig{},
		v4.WithCLIVersion(version.Version),
		v4.WithEnvPrefix("TASKCTL_"),
		v4.WithAuditLog(auditPlugin),
		v4.WithAuditMiddleware[AppConfig](auditPlugin),
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
		// Flight recorder — captures execution traces for slow or failing
		// commands. Snapshots are written to /tmp and analyzed with:
		//   go tool trace /tmp/cmdguard-*.trace
		flightrecorder.WithFlightRecorderRecorder[AppConfig](rec),
		glamour.WithHelpTheme("dark"),
		v4.WithGroup("tasks", "Task Management"),
		v4.WithGroup("system", "System"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating CLI: %w", err)
	}

	// Global flags available to all commands
	cli.AddGlobalBoolFlag("debug", "D", false, "Enable debug mode")

	// Register DI services
	if err := v4.Provide(cli.Scope(), NewTaskStore); err != nil {
		return nil, fmt.Errorf("registering TaskStore: %w", err)
	}

	// Seed demo data
	seedTasks(cli)

	// Build all commands
	if err := buildCommands(cli); err != nil {
		return nil, fmt.Errorf("building commands: %w", err)
	}

	return cli, nil
}

// exportAuditLog writes the DI audit log to disk when any events were
// captured. AUDIT_LOG_FORMAT selects the export format: html, json, ndjson,
// csv, tsv, mermaid, dot, d2, plantuml, tree, or htmltree (default: html).
func exportAuditLog(cli *v4.CLI[AppConfig]) {
	plugin := cli.AuditLog()
	if plugin == nil || plugin.EventsCount() == 0 {
		return
	}

	format, err := v4.ParseAuditLogFormat(os.Getenv("AUDIT_LOG_FORMAT"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "audit-log format invalid: %v\n", err)

		return
	}

	path := "taskctl-audit." + format.String()
	if err := v4.ExportAuditLog(cli, v4.AuditLogExportConfig{
		Format: format,
		Path:   path,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "audit-log export failed: %v\n", err)

		return
	}

	fmt.Fprintf(os.Stderr, "audit-log written to %s\n", path)
}
