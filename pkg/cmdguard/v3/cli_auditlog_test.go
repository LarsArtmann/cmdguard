package v3_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	auditlog "github.com/larsartmann/samber-do-auditlog"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// newTestPlugin returns a fresh auditlog.Plugin with the same config every
// test case in this file uses (Enabled, ContainerID="test"). Centralised so a
// future Config field change updates 11 call sites in one place.
func newTestPlugin(t *testing.T) *auditlog.Plugin {
	t.Helper()

	plugin, err := auditlog.New(auditlog.Config{
		Enabled:     true,
		ContainerID: "test",
	})
	if err != nil {
		t.Fatalf("creating auditlog plugin: %v", err)
	}

	return plugin
}

func TestWithAuditLog(t *testing.T) {
	t.Parallel()

	t.Run("auditlog plugin is accessible via accessor", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli := newTestCLIWithAuditLog(t, plugin)

		if cli.AuditLog() == nil {
			t.Fatal("AuditLog() returned nil, expected the plugin")
		}

		if cli.AuditLog() != plugin {
			t.Error("AuditLog() returned wrong plugin instance")
		}
	})

	t.Run("nil plugin leaves AuditLog nil", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		if cli.AuditLog() != nil {
			t.Error("AuditLog() should be nil when not configured")
		}
	})

	t.Run("auditlog captures service events", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli := newTestCLIWithAuditLog(t, plugin)

		err := cli.ExecuteWithArgs(t.Context(), []string{})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		report := plugin.Report()
		if report.EventCount == 0 {
			t.Error("expected audit log to capture events")
		}

		if report.ServiceCount == 0 {
			t.Error("expected audit log to capture services")
		}
	})

	t.Run("AuditLogReport returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		if cli.AuditLogReport() != nil {
			t.Error("AuditLogReport() should be nil when audit logging is not enabled")
		}
	})

	t.Run("AuditLogReport returns report when enabled", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli := newTestCLIWithAuditLog(t, plugin)

		_ = cli.ExecuteWithArgs(t.Context(), []string{})

		report := cli.AuditLogReport()
		if report == nil {
			t.Fatal("AuditLogReport() returned nil")
		}

		if report.ContainerID != "test" {
			t.Errorf("ContainerID = %q, want %q", report.ContainerID, "test")
		}
	})

	t.Run("RecordAuditHealthCheck returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		results := cli.RecordAuditHealthCheck(t.Context())
		if results != nil {
			t.Error("RecordAuditHealthCheck() should return nil when not enabled")
		}
	})

	t.Run("combines with DILogging", func(t *testing.T) {
		t.Parallel()

		var logs []string

		plugin := newTestPlugin(t)

		cli, err := v3.NewCLI(
			"test", "Test", testCLIConfig{},
			v3.WithAuditLog(plugin),
			v3.WithDILogging(func(format string, args ...any) {
				logs = append(logs, format)
			}),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if len(logs) == 0 {
			t.Error("expected DI log output to be captured alongside audit log")
		}

		report := plugin.Report()
		if report.EventCount == 0 {
			t.Error("expected audit log to capture events alongside DI logging")
		}
	})
}

func TestAuditLogConvenienceHelpers(t *testing.T) {
	t.Parallel()

	t.Run("AuditLogServiceByName returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		svc := v3.AuditLogServiceByName(cli, "anything")
		if svc != nil {
			t.Error("expected nil when audit logging not enabled")
		}
	})

	t.Run("AuditLogFailedServices returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		failed := v3.AuditLogFailedServices(cli)
		if failed != nil {
			t.Error("expected nil when audit logging not enabled")
		}
	})

	t.Run("AuditLogServiceByName returns service when enabled", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli := newTestCLIWithAuditLog(t, plugin)

		_ = cli.ExecuteWithArgs(t.Context(), []string{})

		report := cli.AuditLogReport()
		if report == nil || report.ServiceCount == 0 {
			t.Fatal("expected audit log to capture services")
		}

		name := report.Services[0].ServiceName
		svc := v3.AuditLogServiceByName(cli, name)
		if svc == nil {
			t.Fatalf("expected to find service %q in audit log", name)
		}

		if svc.Status != "registered" && svc.Status != "active" {
			t.Errorf("Status = %q, want registered or active", svc.Status)
		}
	})
}

func TestParseAuditLogFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    v3.AuditLogFormat
		wantErr bool
	}{
		{"", v3.AuditLogFormatHTML, false},
		{"html", v3.AuditLogFormatHTML, false},
		{"htmltree", v3.AuditLogFormatHTMLTree, false},
		{"json", v3.AuditLogFormatJSON, false},
		{"ndjson", v3.AuditLogFormatNDJSON, false},
		{"mermaid", v3.AuditLogFormatMermaid, false},
		{"csv", v3.AuditLogFormatCSV, false},
		{"tsv", v3.AuditLogFormatTSV, false},
		{"dot", v3.AuditLogFormatDOT, false},
		{"d2", v3.AuditLogFormatD2, false},
		{"plantuml", v3.AuditLogFormatPlantUML, false},
		{"tree", v3.AuditLogFormatTree, false},
		{"xml", "", true},
		{"HTML", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := v3.ParseAuditLogFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if !errors.Is(err, v3.ErrUnsupportedAuditLogFormat) {
					t.Errorf("error = %v, want ErrUnsupportedAuditLogFormat", err)
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ParseAuditLogFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExportAuditLog(t *testing.T) {
	t.Parallel()

	t.Run("noop when audit log not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
			Format: v3.AuditLogFormatHTML,
		})
		if err != nil {
			t.Fatalf("ExportAuditLog returned error: %v", err)
		}
	})

	t.Run("noop when no events captured", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)

		cli := newTestCLIWithAuditLog(t, plugin)

		// Don't execute the CLI, so no events are captured
		err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
			Format: v3.AuditLogFormatHTML,
		})
		if err != nil {
			t.Fatalf("ExportAuditLog returned error: %v", err)
		}
	})

	formats := []v3.AuditLogFormat{
		v3.AuditLogFormatJSON,
		v3.AuditLogFormatNDJSON,
		v3.AuditLogFormatMermaid,
		v3.AuditLogFormatCSV,
		v3.AuditLogFormatTSV,
		v3.AuditLogFormatDOT,
		v3.AuditLogFormatD2,
		v3.AuditLogFormatPlantUML,
		v3.AuditLogFormatTree,
		v3.AuditLogFormatHTMLTree,
	}

	for _, format := range formats {
		t.Run("writer/"+format.String(), func(t *testing.T) {
			t.Parallel()

			plugin := newTestPlugin(t)
			cli := newTestCLIWithAuditLog(t, plugin)

			_ = cli.ExecuteWithArgs(t.Context(), []string{})

			var buf bytes.Buffer

			err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
				Format: format,
				Writer: &buf,
			})
			if err != nil {
				t.Fatalf("ExportAuditLog(%s) returned error: %v", format, err)
			}

			if buf.Len() == 0 {
				t.Errorf("expected non-empty %s output", format)
			}
		})

		t.Run("file/"+format.String(), func(t *testing.T) {
			t.Parallel()

			plugin := newTestPlugin(t)
			cli := newTestCLIWithAuditLog(t, plugin)

			_ = cli.ExecuteWithArgs(t.Context(), []string{})

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "audit-"+format.String()+".txt")

			err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
				Format: format,
				Path:   path,
			})
			if err != nil {
				t.Fatalf("ExportAuditLog(%s, file) returned error: %v", format, err)
			}

			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("output file not created: %v", err)
			}
			if info.Size() == 0 {
				t.Errorf("output file %s is empty", path)
			}
		})
	}

	t.Run("HTML file export produces valid document", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)
		cli := newTestCLIWithAuditLog(t, plugin)

		_ = cli.ExecuteWithArgs(t.Context(), []string{})

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "audit.html")

		err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
			Format: v3.AuditLogFormatHTML,
			Path:   path,
		})
		if err != nil {
			t.Fatalf("ExportAuditLog(html, file) returned error: %v", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading output file: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "<html") && !strings.Contains(content, "<!DOCTYPE") {
			t.Errorf("expected HTML document, got: %s", content[:min(len(content), 200)])
		}
	})

	t.Run("rejects unsupported format", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin(t)
		cli := newTestCLIWithAuditLog(t, plugin)

		_ = cli.ExecuteWithArgs(t.Context(), []string{})

		err := v3.ExportAuditLog(cli, v3.AuditLogExportConfig{
			Format: v3.AuditLogFormat("xml"),
			Writer: &bytes.Buffer{},
		})
		if err == nil {
			t.Fatal("expected error for unsupported format")
		}
		if !errors.Is(err, v3.ErrUnsupportedAuditLogFormat) {
			t.Errorf("error = %v, want ErrUnsupportedAuditLogFormat", err)
		}
	})
}
