package v2

import (
	"context"
	"fmt"
	"io"
	"os"

	auditlog "github.com/larsartmann/samber-do-auditlog"
)

// auditLogFlags holds the flags for the audit-log command.
type auditLogFlags struct {
	Format string `flag:"format" short:"f" default:"html" help:"Output format"                         values:"html,json,ndjson,mermaid"`
	Output string `flag:"output" short:"o" default:""     help:"Output file path (defaults to stdout)"`
}

// auditLogConfig holds the configuration for the audit-log command.
type auditLogConfig[T any] struct {
	short   string
	long    string
	groupID string
}

// AuditLogOption configures the audit-log command.
type AuditLogOption[T any] func(*auditLogConfig[T])

// WithAuditLogShort sets a custom short description for the audit-log command.
func WithAuditLogShort[T any](short string) AuditLogOption[T] {
	return func(cfg *auditLogConfig[T]) {
		cfg.short = short
	}
}

// WithAuditLogLong sets a custom long description for the audit-log command.
func WithAuditLogLong[T any](long string) AuditLogOption[T] {
	return func(cfg *auditLogConfig[T]) {
		cfg.long = long
	}
}

// WithAuditLogGroupID sets the command group ID for the audit-log command.
func WithAuditLogGroupID[T any](groupID string) AuditLogOption[T] {
	return func(cfg *auditLogConfig[T]) {
		cfg.groupID = groupID
	}
}

// AuditLogCommand creates a typed "audit-log" subcommand that exports the DI audit log.
// Returns an error if audit logging is not enabled (use WithAuditLog first).
//
// The command supports four output formats:
//   - html (default): self-contained HTML visualization with dependency graph and timeline
//   - json: full structured report
//   - ndjson: streaming event log
//   - mermaid: dependency graph flowchart
//
// Usage:
//
//	plugin := auditlog.New(auditlog.Config{Enabled: true})
//	cli, _ := v2.NewCLI[Config]("myapp", "My app", Config{},
//	    v2.WithAuditLog[Config](plugin),
//	)
//	auditCmd, err := v2.AuditLogCommand[Config](cli)
//	if err != nil { log.Fatal(err) }
//	v2.AddCommand(cli, auditCmd)
func AuditLogCommand[T any](cli *CLI[T], opts ...AuditLogOption[T]) (Command[T, auditLogFlags], error) {
	if cli.auditLog == nil {
		return Command[T, auditLogFlags]{}, fmt.Errorf(
			"%w: audit-log command requires WithAuditLog",
			ErrAuditLogNotEnabled,
		)
	}

	cfg := auditLogConfig[T]{
		short: "Export DI audit log",
		long:  "Export a snapshot of the DI container audit log with dependency graph, service timeline, and event stream.",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	plugin := cli.auditLog
	appName := cli.name

	return NewCommand[T, auditLogFlags](
		"audit-log",
		func(_ context.Context, _ *T, flags auditLogFlags) error {
			w := cli.rootCmd.OutOrStdout()

			switch flags.Format {
			case "html", "":
				return writeAuditToFileOrWriter(
					appName, flags.Output,
					plugin.ExportToHTML, plugin.WriteHTML, w,
				)
			case "json":
				return writeAuditToFileOrWriter(
					appName, flags.Output,
					plugin.ExportToFile, plugin.WriteReportJSON, w,
				)
			case "ndjson":
				return writeAuditToFileOrWriter(
					appName, flags.Output,
					plugin.ExportEventsToNDJSON, plugin.WriteEventsNDJSON, w,
				)
			case "mermaid":
				return writeAuditMermaid(appName, flags.Output, plugin, w)
			default:
				return fmt.Errorf(
					"%w: %q (use html, json, ndjson, or mermaid)",
					ErrInvalidOutputFormat,
					flags.Format,
				)
			}
		},
		WithShort[T, auditLogFlags](cfg.short),
		WithLong[T, auditLogFlags](cfg.long),
		withAuditLogGroupID[T](cfg.groupID),
	)
}

func writeAuditToFileOrWriter(
	appName, path string,
	exportFn func(string) error,
	writeFn func(io.Writer) error,
	w io.Writer,
) error {
	if path != "" {
		err := exportFn(path)
		if err != nil {
			return fmt.Errorf("exporting audit log for %q: %w", appName, err)
		}

		return nil
	}

	err := writeFn(w)
	if err != nil {
		return fmt.Errorf("writing audit log for %q: %w", appName, err)
	}

	return nil
}

func writeAuditMermaid(appName, path string, plugin *auditlog.Plugin, w io.Writer) error {
	report := plugin.Report()

	if path != "" {
		f, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("creating mermaid file %q: %w", path, err)
		}

		defer f.Close()

		err = report.WriteMermaid(f)
		if err != nil {
			return fmt.Errorf("writing mermaid audit log for %q: %w", appName, err)
		}

		return nil
	}

	err := report.WriteMermaid(w)
	if err != nil {
		return fmt.Errorf("writing mermaid audit log for %q: %w", appName, err)
	}

	return nil
}

// withAuditLogGroupID conditionally adds WithGroupID if non-empty.
func withAuditLogGroupID[T any](groupID string) CommandOption[T, auditLogFlags] {
	return func(cmd *Command[T, auditLogFlags]) {
		if groupID != "" {
			WithGroupID[T, auditLogFlags](groupID)(cmd)
		}
	}
}

// AuditLogServiceByName returns the first auditlog ServiceInfo matching the name.
// Returns nil if audit logging is not enabled.
func AuditLogServiceByName[T any](cli *CLI[T], name string) *auditlog.ServiceInfo {
	if cli.auditLog == nil {
		return nil
	}

	return cli.auditLog.Report().ServiceByName(name)
}

// AuditLogFailedServices returns all services with invocation or shutdown errors.
// Returns nil if audit logging is not enabled.
func AuditLogFailedServices[T any](cli *CLI[T]) []auditlog.ServiceInfo {
	if cli.auditLog == nil {
		return nil
	}

	return cli.auditLog.Report().FailedServices()
}
