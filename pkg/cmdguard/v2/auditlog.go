package v2

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	auditlog "github.com/larsartmann/samber-do-auditlog"
)

// ErrUnsupportedAuditLogFormat indicates an unsupported audit log export format.
var ErrUnsupportedAuditLogFormat = errors.New("unsupported audit log format")

// AuditLogFormat is a validated output format for audit log exports.
// Construct via ParseAuditLogFormat to guarantee a valid value.
type AuditLogFormat string

const (
	AuditLogFormatHTML    AuditLogFormat = "html"
	AuditLogFormatJSON    AuditLogFormat = "json"
	AuditLogFormatNDJSON  AuditLogFormat = "ndjson"
	AuditLogFormatMermaid AuditLogFormat = "mermaid"
)

var supportedAuditLogFormats = []AuditLogFormat{
	AuditLogFormatHTML,
	AuditLogFormatJSON,
	AuditLogFormatNDJSON,
	AuditLogFormatMermaid,
}

// Valid returns true if the format is one of the supported values.
func (f AuditLogFormat) Valid() bool {
	return slices.Contains(supportedAuditLogFormats, f)
}

// String implements fmt.Stringer.
func (f AuditLogFormat) String() string {
	return string(f)
}

// ParseAuditLogFormat converts a string to an AuditLogFormat.
// Empty string defaults to HTML.
func ParseAuditLogFormat(s string) (AuditLogFormat, error) {
	if s == "" {
		return AuditLogFormatHTML, nil
	}

	format := AuditLogFormat(s)
	if !format.Valid() {
		names := make([]string, len(supportedAuditLogFormats))
		for i, f := range supportedAuditLogFormats {
			names[i] = f.String()
		}

		return "", fmt.Errorf("%w: %q (use %s)", ErrUnsupportedAuditLogFormat, s, strings.Join(names, ", "))
	}

	return format, nil
}

// AuditLogExportConfig configures where and how to write an audit log export.
type AuditLogExportConfig struct {
	Format AuditLogFormat

	// Path is the output file path. If set, the audit log is written to this file.
	Path string

	// Writer receives the output when Path is empty.
	// If both Path and Writer are empty/nil, defaults to os.Stdout.
	Writer io.Writer
}

// ExportAuditLog writes the CLI's audit log snapshot in the configured format.
// Returns nil (no-op) if audit logging is not enabled or no events were captured.
func ExportAuditLog[T any](cli *CLI[T], cfg AuditLogExportConfig) error {
	plugin := cli.AuditLog()
	if plugin == nil || plugin.EventsCount() == 0 {
		return nil
	}

	if cfg.Path != "" {
		return exportAuditLogToFile(plugin, cfg.Format, cfg.Path)
	}

	writer := cfg.Writer
	if writer == nil {
		writer = os.Stdout
	}

	return exportAuditLogToWriter(plugin, cfg.Format, writer)
}

func exportAuditLogToFile(plugin *auditlog.Plugin, format AuditLogFormat, path string) error {
	switch format {
	case AuditLogFormatHTML:
		if err := plugin.ExportToHTML(path); err != nil {
			return fmt.Errorf("exporting HTML audit log to %q: %w", path, err)
		}
	case AuditLogFormatJSON:
		if err := plugin.ExportToFile(path); err != nil {
			return fmt.Errorf("exporting JSON audit log to %q: %w", path, err)
		}
	case AuditLogFormatNDJSON:
		if err := plugin.ExportEventsToNDJSON(path); err != nil {
			return fmt.Errorf("exporting NDJSON audit log to %q: %w", path, err)
		}
	case AuditLogFormatMermaid:
		if err := exportMermaidToFile(plugin, path); err != nil {
			return fmt.Errorf("exporting mermaid audit log to %q: %w", path, err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	return nil
}

func exportAuditLogToWriter(plugin *auditlog.Plugin, format AuditLogFormat, w io.Writer) error {
	switch format {
	case AuditLogFormatHTML:
		if err := plugin.WriteHTML(w); err != nil {
			return fmt.Errorf("writing HTML audit log: %w", err)
		}
	case AuditLogFormatJSON:
		if err := plugin.WriteReportJSON(w); err != nil {
			return fmt.Errorf("writing JSON audit log: %w", err)
		}
	case AuditLogFormatNDJSON:
		if err := plugin.WriteEventsNDJSON(w); err != nil {
			return fmt.Errorf("writing NDJSON audit log: %w", err)
		}
	case AuditLogFormatMermaid:
		if err := plugin.Report().WriteMermaid(w); err != nil {
			return fmt.Errorf("writing mermaid audit log: %w", err)
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	return nil
}

func exportMermaidToFile(plugin *auditlog.Plugin, path string) error {
	report := plugin.Report()

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating mermaid file %q: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	if err := report.WriteMermaid(file); err != nil {
		return fmt.Errorf("writing mermaid content: %w", err)
	}

	return nil
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
