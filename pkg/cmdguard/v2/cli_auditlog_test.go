package v2_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	auditlog "github.com/larsartmann/samber-do-auditlog"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

// newTestPlugin returns a fresh auditlog.Plugin with the same config every
// test case in this file uses (Enabled, ContainerID="test"). Centralised so a
// future Config field change updates 11 call sites in one place.
func newTestPlugin() *auditlog.Plugin {
	return auditlog.New(auditlog.Config{
		Enabled:     true,
		ContainerID: "test",
	})
}

func TestWithAuditLog(t *testing.T) {
	t.Parallel()

	t.Run("auditlog plugin is accessible via accessor", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin()

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

		plugin := newTestPlugin()

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

		plugin := newTestPlugin()

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

		plugin := newTestPlugin()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test", testCLIConfig{},
			v2.WithAuditLog[testCLIConfig](plugin),
			v2.WithDILogging[testCLIConfig](func(format string, args ...any) {
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

func TestAuditLogCommand(t *testing.T) {
	t.Parallel()

	t.Run("returns error when audit logging not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		_, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err == nil {
			t.Fatal("expected error when audit logging not enabled")
		}
	})

	t.Run("exports JSON to stdout", func(t *testing.T) {
		var buf bytes.Buffer

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, auditCmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.RootCommand().SetOut(&buf)

		err = cli.ExecuteWithArgs(t.Context(), []string{"audit-log", "--format", "json"})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		var report map[string]any
		err = json.Unmarshal(buf.Bytes(), &report)
		if err != nil {
			t.Fatalf("output is not valid JSON: %v\noutput: %s", err, buf.String())
		}

		if cid, _ := report["container_id"].(string); cid != "test" {
			t.Errorf("container_id = %q, want %q", cid, "test")
		}
	})

	t.Run("exports NDJSON to stdout", func(t *testing.T) {
		var buf bytes.Buffer

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, auditCmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.RootCommand().SetOut(&buf)

		err = cli.ExecuteWithArgs(t.Context(), []string{"audit-log", "--format", "ndjson"})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
		if len(lines) == 0 {
			t.Fatal("expected NDJSON output with at least one event")
		}

		var firstEvent map[string]any
		err = json.Unmarshal([]byte(lines[0]), &firstEvent)
		if err != nil {
			t.Fatalf("first line is not valid JSON: %v", err)
		}
	})

	t.Run("exports Mermaid to stdout", func(t *testing.T) {
		var buf bytes.Buffer

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, auditCmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.RootCommand().SetOut(&buf)

		err = cli.ExecuteWithArgs(t.Context(), []string{"audit-log", "--format", "mermaid"})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if !strings.Contains(buf.String(), "flowchart TD") {
			t.Errorf("expected mermaid output to contain 'flowchart TD', got: %s", buf.String())
		}
	})

	t.Run("exports HTML to stdout", func(t *testing.T) {
		var buf bytes.Buffer

		plugin := newTestPlugin()

		cli, err := v2.NewCLI[testCLIConfig](
			"test", "Test", testCLIConfig{},
			v2.WithAuditLog[testCLIConfig](plugin),
			v2.WithFang[testCLIConfig](false),
		)
		if err != nil {
			t.Fatalf("NewCLI failed: %v", err)
		}

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, auditCmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		cli.RootCommand().SetOut(&buf)

		err = cli.ExecuteWithArgs(t.Context(), []string{"audit-log"})
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "<!DOCTYPE html>") && !strings.Contains(output, "<html") {
			t.Errorf("expected HTML output to contain HTML elements, got: %s", output[:min(len(output), 200)])
		}

		if !strings.Contains(output, "do-auditlog") {
			t.Errorf("expected HTML output to contain do-auditlog title, got: %s", output[:min(len(output), 200)])
		}
	})

	t.Run("rejects invalid format", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](cli)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		err = v2.AddCommand(cli, auditCmd)
		if err != nil {
			t.Fatalf("AddCommand failed: %v", err)
		}

		err = cli.ExecuteWithArgs(t.Context(), []string{"audit-log", "--format", "xml"})
		if err == nil {
			t.Fatal("expected error for invalid format")
		}
	})

	t.Run("accepts custom descriptions", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		auditCmd, err := v2.AuditLogCommand[testCLIConfig](
			cli,
			v2.WithAuditLogShort[testCLIConfig]("Custom short"),
			v2.WithAuditLogLong[testCLIConfig]("Custom long description"),
		)
		if err != nil {
			t.Fatalf("AuditLogCommand failed: %v", err)
		}

		if auditCmd.Short() != "Custom short" {
			t.Errorf("Short = %q, want %q", auditCmd.Short(), "Custom short")
		}
	})
}

func TestAuditLogConvenienceHelpers(t *testing.T) {
	t.Parallel()

	t.Run("AuditLogServiceByName returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		svc := v2.AuditLogServiceByName(cli, "anything")
		if svc != nil {
			t.Error("expected nil when audit logging not enabled")
		}
	})

	t.Run("AuditLogFailedServices returns nil when not enabled", func(t *testing.T) {
		t.Parallel()

		cli := newTestCLI(t)

		failed := v2.AuditLogFailedServices(cli)
		if failed != nil {
			t.Error("expected nil when audit logging not enabled")
		}
	})

	t.Run("AuditLogServiceByName returns service when enabled", func(t *testing.T) {
		t.Parallel()

		plugin := newTestPlugin()

		cli := newTestCLIWithAuditLog(t, plugin)

		_ = cli.ExecuteWithArgs(t.Context(), []string{})

		report := cli.AuditLogReport()
		if report == nil || report.ServiceCount == 0 {
			t.Fatal("expected audit log to capture services")
		}

		name := report.Services[0].ServiceName
		svc := v2.AuditLogServiceByName(cli, name)
		if svc == nil {
			t.Fatalf("expected to find service %q in audit log", name)
		}

		if svc.Status != "registered" && svc.Status != "active" {
			t.Errorf("Status = %q, want registered or active", svc.Status)
		}
	})
}
