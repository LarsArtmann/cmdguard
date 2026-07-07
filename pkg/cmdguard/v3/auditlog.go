package v3

import (
	"errors"
	"fmt"
	"io"
	"maps"
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
	AuditLogFormatCSV      AuditLogFormat = "csv"
	AuditLogFormatD2       AuditLogFormat = "d2"
	AuditLogFormatDOT      AuditLogFormat = "dot"
	AuditLogFormatHTML     AuditLogFormat = "html"
	AuditLogFormatHTMLTree AuditLogFormat = "htmltree"
	AuditLogFormatJSON     AuditLogFormat = "json"
	AuditLogFormatMermaid  AuditLogFormat = "mermaid"
	AuditLogFormatNDJSON   AuditLogFormat = "ndjson"
	AuditLogFormatPlantUML AuditLogFormat = "plantuml"
	AuditLogFormatTree     AuditLogFormat = "tree"
	AuditLogFormatTSV      AuditLogFormat = "tsv"
)

// auditLogExporter binds a format to its file and writer export implementations.
// The auditLogExporterRegistry map is the single source of truth for which
// formats are supported; Valid(), ParseAuditLogFormat, and both export
// directions all derive from it. Adding a format = one entry there.
type auditLogExporter struct {
	toFile   func(*auditlog.Plugin, string) error
	toWriter func(*auditlog.Plugin, io.Writer) error
}

// auditLogExporterRegistry is the single source of truth for supported audit
// log export formats. Kept as a function (not a package-level var) to satisfy
// gochecknoglobals; the per-call allocation is negligible (export runs once).
func auditLogExporterRegistry() map[AuditLogFormat]auditLogExporter {
	return map[AuditLogFormat]auditLogExporter{
		AuditLogFormatHTML:     {(*auditlog.Plugin).ExportToHTML, (*auditlog.Plugin).WriteHTML},
		AuditLogFormatJSON:     {(*auditlog.Plugin).ExportToFile, (*auditlog.Plugin).WriteReportJSON},
		AuditLogFormatNDJSON:   {(*auditlog.Plugin).ExportEventsToNDJSON, (*auditlog.Plugin).WriteEventsNDJSON},
		AuditLogFormatCSV:      {(*auditlog.Plugin).ExportToCSV, (*auditlog.Plugin).WriteReportCSV},
		AuditLogFormatTSV:      {(*auditlog.Plugin).ExportToTSV, (*auditlog.Plugin).WriteReportTSV},
		AuditLogFormatMermaid:  {(*auditlog.Plugin).ExportToMermaid, (*auditlog.Plugin).WriteMermaid},
		AuditLogFormatDOT:      {(*auditlog.Plugin).ExportToDOT, (*auditlog.Plugin).WriteDOT},
		AuditLogFormatD2:       {(*auditlog.Plugin).ExportToD2, (*auditlog.Plugin).WriteD2},
		AuditLogFormatPlantUML: {(*auditlog.Plugin).ExportToPlantUML, (*auditlog.Plugin).WritePlantUML},
		AuditLogFormatTree:     {(*auditlog.Plugin).ExportToTree, (*auditlog.Plugin).WriteTree},
		AuditLogFormatHTMLTree: {(*auditlog.Plugin).ExportToHTMLTree, (*auditlog.Plugin).WriteHTMLTree},
	}
}

// Valid returns true if the format is one of the supported values.
func (f AuditLogFormat) Valid() bool {
	_, ok := auditLogExporterRegistry()[f]

	return ok
}

// String implements fmt.Stringer.
func (f AuditLogFormat) String() string {
	return string(f)
}

// supportedAuditLogFormatNames returns the names of all supported formats for
// error messages, derived from the single registry so it can never drift.
// Sorted for deterministic output.
func supportedAuditLogFormatNames() []string {
	keys := slices.SortedFunc(maps.Keys(auditLogExporterRegistry()), func(a, b AuditLogFormat) int {
		return strings.Compare(string(a), string(b))
	})

	names := make([]string, len(keys))
	for i, k := range keys {
		names[i] = string(k)
	}

	return names
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
	exporter, ok := auditLogExporterRegistry()[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	err := exporter.toFile(plugin, path)
	if err != nil {
		return fmt.Errorf("exporting %s audit log to %q: %w", format, path, err)
	}

	return nil
}

func exportAuditLogToWriter(plugin *auditlog.Plugin, format AuditLogFormat, w io.Writer) error {
	exporter, ok := auditLogExporterRegistry()[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnsupportedAuditLogFormat, format)
	}

	err := exporter.toWriter(plugin, w)
	if err != nil {
		return fmt.Errorf("writing %s audit log: %w", format, err)
	}

	return nil
}

// AuditLogServiceByName returns the first auditlog ServiceInfo matching the name.
// Returns nil if audit logging is not enabled.
func AuditLogServiceByName[T any](cli *CLI[T], name string) *auditlog.ServiceInfo {
	if cli.spec.auditLog == nil {
		return nil
	}

	return cli.spec.auditLog.Report().ServiceByName(name)
}

// AuditLogFailedServices returns all services with invocation or shutdown errors.
// Returns nil if audit logging is not enabled.
func AuditLogFailedServices[T any](cli *CLI[T]) []auditlog.ServiceInfo {
	if cli.spec.auditLog == nil {
		return nil
	}

	return cli.spec.auditLog.Report().FailedServices()
}
