package v2

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	AuditLogFormatCSV     AuditLogFormat = "csv"
	AuditLogFormatTSV     AuditLogFormat = "tsv"
	AuditLogFormatDOT     AuditLogFormat = "dot"
)

// Valid returns true if the format is one of the supported values.
func (f AuditLogFormat) Valid() bool {
	switch f {
	case AuditLogFormatHTML, AuditLogFormatJSON, AuditLogFormatNDJSON,
		AuditLogFormatMermaid, AuditLogFormatCSV, AuditLogFormatTSV, AuditLogFormatDOT:
		return true
	}

	return false
}

// String implements fmt.Stringer.
func (f AuditLogFormat) String() string {
	return string(f)
}

// supportedAuditLogFormatNames returns the names of all supported formats for
// error messages. Kept as a function (not a global var) to satisfy gochecknoglobals.
func supportedAuditLogFormatNames() []string {
	return []string{
		string(AuditLogFormatHTML),
		string(AuditLogFormatJSON),
		string(AuditLogFormatNDJSON),
		string(AuditLogFormatMermaid),
		string(AuditLogFormatCSV),
		string(AuditLogFormatTSV),
		string(AuditLogFormatDOT),
	}
}

// ParseAuditLogFormat converts a string to an AuditLogFormat.
// Empty string defaults to HTML.
func ParseAuditLogFormat(s string) (AuditLogFormat, error) {
	if s == "" {
		return AuditLogFormatHTML, nil
	}

	format := AuditLogFormat(s)
	if !format.Valid() {
		allowed := strings.Join(supportedAuditLogFormatNames(), ", ")

		return "", fmt.Errorf("%w: %q (use %s)", ErrUnsupportedAuditLogFormat, s, allowed)
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
	exporters := map[AuditLogFormat]func(*auditlog.Plugin, string) error{
		AuditLogFormatHTML:    (*auditlog.Plugin).ExportToHTML,
		AuditLogFormatJSON:    (*auditlog.Plugin).ExportToFile,
		AuditLogFormatNDJSON:  (*auditlog.Plugin).ExportEventsToNDJSON,
		AuditLogFormatCSV:     (*auditlog.Plugin).ExportToCSV,
		AuditLogFormatTSV:     (*auditlog.Plugin).ExportToTSV,
		AuditLogFormatMermaid: exportMermaidReportToFile,
		AuditLogFormatDOT:     exportDOTReportToFile,
	}

	exporter, ok := exporters[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	if err := exporter(plugin, path); err != nil {
		return fmt.Errorf("exporting %s audit log to %q: %w", format, path, err)
	}

	return nil
}

func exportAuditLogToWriter(plugin *auditlog.Plugin, format AuditLogFormat, w io.Writer) error {
	exporters := map[AuditLogFormat]func(*auditlog.Plugin, io.Writer) error{
		AuditLogFormatHTML:    (*auditlog.Plugin).WriteHTML,
		AuditLogFormatJSON:    (*auditlog.Plugin).WriteReportJSON,
		AuditLogFormatNDJSON:  (*auditlog.Plugin).WriteEventsNDJSON,
		AuditLogFormatCSV:     (*auditlog.Plugin).WriteReportCSV,
		AuditLogFormatTSV:     (*auditlog.Plugin).WriteReportTSV,
		AuditLogFormatMermaid: writeMermaidReport,
		AuditLogFormatDOT:     writeDOTReport,
	}

	exporter, ok := exporters[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	if err := exporter(plugin, w); err != nil {
		return fmt.Errorf("writing %s audit log: %w", format, err)
	}

	return nil
}

// exportMermaidReportToFile and exportDOTReportToFile adapt Report-level WriteX
// methods to the file-exporter signature for formats lacking Plugin.ExportToX.
func exportMermaidReportToFile(plugin *auditlog.Plugin, path string) error {
	return writeReportToFile(path, plugin.Report().WriteMermaid)
}

func exportDOTReportToFile(plugin *auditlog.Plugin, path string) error {
	return writeReportToFile(path, plugin.Report().WriteDOT)
}

func writeMermaidReport(plugin *auditlog.Plugin, w io.Writer) error {
	if err := plugin.Report().WriteMermaid(w); err != nil {
		return fmt.Errorf("mermaid report: %w", err)
	}

	return nil
}

func writeDOTReport(plugin *auditlog.Plugin, w io.Writer) error {
	if err := plugin.Report().WriteDOT(w); err != nil {
		return fmt.Errorf("DOT report: %w", err)
	}

	return nil
}

// writeReportToFile creates path, invokes write on the file, then closes it.
// Used for report formats (mermaid, dot) that lack a Plugin-level ExportToX method.
func writeReportToFile(path string, write func(io.Writer) error) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	if err := write(file); err != nil {
		return fmt.Errorf("writing report content: %w", err)
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
