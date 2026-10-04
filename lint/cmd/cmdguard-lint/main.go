// Command cmdguard-lint runs the cmdguard usage linter over a Go project:
// it detects anti-patterns from real consumer audits (execute bypass, stale
// majors, constructor panics, version/option misuse, double error display).
//
// It is dogfood: the CLI itself is built with cmdguard v4.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/larsartmann/cmdguard/v4/pkg/version"
	"github.com/larsartmann/cmdguard/lint"
	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

type cliConfig struct {
	Output  string `flag:"output" default:"text" help:"Output format: text, json, or sarif" validate:"enum=text,json,sarif"`
	Enable  string `flag:"enable" default:"" help:"Comma-separated rule IDs to run (default: all)"`
	Disable string `flag:"disable" default:"" help:"Comma-separated rule IDs to skip"`
}

type lintFlags struct {
	Dir string `flag:"dir" default:"." help:"Directory to lint"`
}

type rulesFlags = v4.NoFlags

func main() {
	cli, err := v4.NewCLI(
		"cmdguard-lint",
		"cmdguard usage linter",
		cliConfig{},
		v4.WithCLIVersion(version.Version),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	lintCmd, err := v4.NewCommand(
		"lint",
		lintFlags{},
		runLint,
		v4.WithShort("Lint a directory for cmdguard usage anti-patterns"),
		v4.WithExample("cmdguard-lint lint --dir . --output sarif"),
		v4.WithExample("cmdguard-lint lint --disable CG004,CG006"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rulesCmd, err := v4.NewCommand(
		"rules",
		v4.NoFlags{},
		runRules,
		v4.WithShort("List available rules"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v4.AddCommand(cli, lintCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := v4.AddCommand(cli, rulesCmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}

func runLint(_ context.Context, cfg *cliConfig, flags lintFlags) error {
	ctx := context.Background()

	findings, err := detectFiltered(ctx, cfg)
	if err != nil {
		return err
	}

	report := finding.NewReportFromFindings(finding.ToolInfo{Name: lint.ToolName, Version: version.Version}, findings)
	report.ComputeSummary()

	switch cfg.Output {
	case "json":
		return writeJSON(report)
	case "sarif":
		return report.WriteSARIF(ctx, os.Stdout)
	default:
		return finding.FormatText(os.Stdout, findings)
	}
}

// detectFiltered runs lint.Detect and removes findings from disabled rules,
// keeping the single-pass fast path while honoring --enable/--disable.
func detectFiltered(ctx context.Context, cfg *cliConfig) ([]finding.Finding, error) {
	findings, err := lint.Detect(ctx, flagsDir(cfg))
	if err != nil {
		return nil, err
	}

	enabled := idSet(cfg.Enable)
	disabled := idSet(cfg.Disable)

	kept := make([]finding.Finding, 0, len(findings))

	for _, f := range findings {
		switch {
		case disabled[string(f.Rule)]:
			continue
		case len(enabled) > 0 && !enabled[string(f.Rule)]:
			continue
		}

		kept = append(kept, f)
	}

	return kept, nil
}

// flagsDir is a placeholder kept for symmetry; the dir flag lives on lintFlags.
func flagsDir(_ *cliConfig) string { return "." }

// idSet splits a comma-separated flag value into an ID set.
func idSet(value string) map[string]bool {
	set := map[string]bool{}

	for _, part := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			set[trimmed] = true
		}
	}

	return set
}

// writeJSON prints the report as JSON.
func writeJSON(report *finding.Report) error {
	json, err := report.ToJSON()
	if err != nil {
		return fmt.Errorf("marshaling report: %w", err)
	}

	_, err = os.Stdout.Write(json)

	return err
}

func runRules(_ context.Context, _ *cliConfig, _ v4.NoFlags) error {
	for _, rule := range lint.AllRules() {
		fmt.Printf("%s  %-28s %s\n", rule.Meta.ID, rule.Meta.Name, rule.Meta.Description)
	}

	fmt.Println("\nSuppress with //cmdguard-lint:ignore <RULE> <reason> on the offending line or the line above.")

	return nil
}

// ExitCode maps findings to the process exit code: 0 clean, 1 findings.
// Called after runLint via the error path; kept for library consumers.
func ExitCode(findings []finding.Finding) int {
	report := finding.NewReportFromFindings(finding.ToolInfo{Name: lint.ToolName}, findings)
	report.ComputeSummary()

	return linter.ExitCodeFromReport(report)
}
