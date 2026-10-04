// Command cmdguard-lint runs the cmdguard usage linter over a Go project:
// it detects anti-patterns from real consumer audits (execute bypass, stale
// majors, constructor panics, version/option misuse, double error display).
//
// It is dogfood: the CLI itself is built with cmdguard v4.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/larsartmann/cmdguard/lint"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/larsartmann/cmdguard/v4/pkg/version"
	"github.com/larsartmann/go-finding"
)

type cliConfig struct {
	Output string `flag:"output" default:"text" help:"Output format for findings" values:"text,json,sarif"`
	Enable string `flag:"enable" default:"" help:"Comma-separated rule IDs to run (default: all)"`
	Disable string `flag:"disable" default:"" help:"Comma-separated rule IDs to skip"`
}

type lintFlags struct {
	Dir string `flag:"dir" default:"." help:"Directory to lint"`
}

func main() {
	cli, err := v4.NewCLI(
		"cmdguard-lint",
		"cmdguard usage linter",
		cliConfig{},
		v4.WithCLIVersion(version.Version),
		v4.WithSignalHandling(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	lintCmd, lintErr := v4.NewCommand(
		"lint",
		lintFlags{},
		runLint,
		v4.WithShort("Lint a directory for cmdguard usage anti-patterns"),
		v4.WithExample("cmdguard-lint lint --dir . --output sarif"),
		v4.WithExample("cmdguard-lint lint --disable CG004,CG006"),
	)
	if lintErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", lintErr)
		os.Exit(1)
	}

	rulesCmd, rulesErr := v4.NewCommand(
		"rules",
		v4.NoFlags{},
		runRules,
		v4.WithShort("List available rules"),
	)
	if rulesErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", rulesErr)
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

// errFindings is returned when the lint run produced findings; cmdguard
// displays it once and ExecuteAndExit maps it to exit code 1.
var errFindings = errors.New("lint run produced findings") //nolint:err113 // sentinel with dynamic count in the wrapping site

func runLint(ctx context.Context, cfg *cliConfig, flags lintFlags) error {
	findings, err := detectFiltered(ctx, flags.Dir, cfg)
	if err != nil {
		return err
	}

	report := finding.NewReportFromFindings(finding.ToolInfo{Name: lint.ToolName, Version: version.Version}, findings)
	report.ComputeSummary()

	var outputErr error

	switch cfg.Output {
	case "json":
		encoded, marshalErr := json.Marshal(report)
		if marshalErr != nil {
			return fmt.Errorf("marshaling report: %w", marshalErr)
		}

		_, outputErr = os.Stdout.Write(append(encoded, '\n'))
	case "sarif":
		outputErr = report.WriteSARIF(ctx, os.Stdout)
	default:
		outputErr = finding.FormatText(os.Stdout, findings)
	}

	if outputErr != nil {
		return fmt.Errorf("writing output: %w", outputErr)
	}

	if len(findings) > 0 {
		return fmt.Errorf("%d cmdguard-lint finding(s) (suppress with //cmdguard-lint:ignore <RULE> <reason>): %w", len(findings), errFindings)
	}

	return nil
}

// detectFiltered runs lint.Detect and removes findings from disabled rules,
// keeping the single-pass fast path while honoring --enable/--disable.
func detectFiltered(ctx context.Context, dir string, cfg *cliConfig) ([]finding.Finding, error) {
	findings, err := lint.Detect(ctx, dir)
	if err != nil {
		return nil, err
	}

	enabled := idSet(cfg.Enable)
	disabled := idSet(cfg.Disable)

	kept := make([]finding.Finding, 0, len(findings))

	for _, f := range findings {
		if disabled[string(f.Rule)] {
			continue
		}

		if len(enabled) > 0 && !enabled[string(f.Rule)] {
			continue
		}

		kept = append(kept, f)
	}

	return kept, nil
}

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

func runRules(_ context.Context, _ *cliConfig, _ v4.NoFlags) error {
	for _, rule := range lint.AllRules() {
		fmt.Printf("%s  %-28s %s\n", rule.Meta.ID, rule.Meta.Name, rule.Meta.Description)
	}

	fmt.Println("\nSuppress with //cmdguard-lint:ignore <RULE> <reason> on the offending line or the line above.")

	return nil
}
