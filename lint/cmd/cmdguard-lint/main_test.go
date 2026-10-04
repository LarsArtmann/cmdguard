package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
)

// captureStdout runs fn with os.Stdout redirected into a pipe and returns
// everything fn wrote. Callers must not run in parallel (os.Stdout swap).
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = writer

	defer func() { os.Stdout = original }()

	fn()

	if err := writer.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read captured output: %v", err)
	}

	return string(out)
}

// writeStaleFixture writes a one-file tree that deterministically triggers
// CG002 (stale major import) at stale.go:3.
func writeStaleFixture(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	stale := "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{}\n"

	if err := os.WriteFile(filepath.Join(dir, "stale.go"), []byte(stale), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	return dir
}

func TestLintCommandCleanTreeExitsZero(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "clean.go"), []byte("package clean\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"lint", "--dir", dir}); execErr != nil {
		t.Errorf("expected clean tree to pass, got: %v", execErr)
	}
}

func TestLintCommandFindingsFailTheRun(t *testing.T) {
	dir := t.TempDir()

	stale := "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{}\n"

	if err := os.WriteFile(filepath.Join(dir, "stale.go"), []byte(stale), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	execErr := cli.ExecuteWithArgs(context.Background(), []string{"lint", "--dir", dir})
	if execErr == nil {
		t.Fatal("expected findings to fail the run")
	}

	if !strings.Contains(execErr.Error(), "finding") {
		t.Errorf("error should mention findings, got: %v", execErr)
	}
}

func TestRulesCommandListsAllRules(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"rules"}); execErr != nil {
		t.Errorf("rules command failed: %v", execErr)
	}
}

func TestLintCommandJSONOutput(t *testing.T) {
	dir := writeStaleFixture(t)

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	var execErr error

	out := captureStdout(t, func() {
		execErr = cli.ExecuteWithArgs(
			context.Background(),
			[]string{"lint", "--dir", dir, "--output", "json"},
		)
	})

	if execErr == nil {
		t.Fatal("expected findings to fail the run even in JSON mode")
	}

	var report struct {
		Tool struct {
			Name string `json:"name"`
		} `json:"tool"`
		Findings []struct {
			Rule     string `json:"rule"`
			Position struct {
				File string `json:"file"`
				Line int    `json:"line"`
			} `json:"position"`
		} `json:"findings"`
	}

	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("stdout is not a valid JSON report: %v\nstdout: %s", err, out)
	}

	if report.Tool.Name != "cmdguard-lint" {
		t.Errorf("tool name = %q, want cmdguard-lint", report.Tool.Name)
	}

	if len(report.Findings) != 1 {
		t.Fatalf("findings = %d, want 1 (CG002 at stale.go:3)\nstdout: %s", len(report.Findings), out)
	}

	got := report.Findings[0]
	if got.Rule != "CG002" || got.Position.File != "stale.go" || got.Position.Line != 3 {
		t.Errorf("finding = %s@%s:%d, want CG002@stale.go:3", got.Rule, got.Position.File, got.Position.Line)
	}
}

func TestLintCommandSARIFOutput(t *testing.T) {
	dir := writeStaleFixture(t)

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	var execErr error

	out := captureStdout(t, func() {
		execErr = cli.ExecuteWithArgs(
			context.Background(),
			[]string{"lint", "--dir", dir, "--output", "sarif"},
		)
	})

	if execErr == nil {
		t.Fatal("expected findings to fail the run even in SARIF mode")
	}

	var sarif struct {
		Version string `json:"version"`
		Runs    []struct {
			Tool struct {
				Driver struct {
					Name string `json:"name"`
				} `json:"driver"`
			} `json:"tool"`
			Results []struct {
				RuleID string `json:"ruleId"`
			} `json:"results"`
		} `json:"runs"`
	}

	if err := json.Unmarshal([]byte(out), &sarif); err != nil {
		t.Fatalf("stdout is not valid SARIF JSON: %v\nstdout: %s", err, out)
	}

	if sarif.Version != "2.1.0" || len(sarif.Runs) != 1 {
		t.Fatalf("SARIF version/runs = %s/%d, want 2.1.0/1", sarif.Version, len(sarif.Runs))
	}

	run := sarif.Runs[0]
	if run.Tool.Driver.Name != "cmdguard-lint" {
		t.Errorf("driver name = %q, want cmdguard-lint", run.Tool.Driver.Name)
	}

	if len(run.Results) != 1 || run.Results[0].RuleID != "CG002" {
		t.Fatalf("results = %+v, want one CG002 result", run.Results)
	}

	roundTripped, err := finding.FindingsFromSARIF(context.Background(), []byte(out))
	if err != nil {
		t.Fatalf("FindingsFromSARIF round-trip: %v\nstdout: %s", err, out)
	}

	if len(roundTripped) != 1 {
		t.Fatalf("round-tripped findings = %d, want 1", len(roundTripped))
	}

	first := roundTripped[0]
	if string(first.Rule) != "CG002" || string(first.Position.File) != "stale.go" || first.Position.Line != 3 {
		t.Errorf(
			"round-tripped finding = %s@%s:%d, want CG002@stale.go:3",
			first.Rule,
			first.Position.File,
			first.Position.Line,
		)
	}
}
