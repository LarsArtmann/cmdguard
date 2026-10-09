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
	"path/filepath"

	"github.com/larsartmann/go-finding"

	"github.com/larsartmann/cmdguard/lint"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/larsartmann/cmdguard/v4/pkg/version"
)

type cliConfig struct {
	Output  string `default:"text" flag:"output"  help:"Output format for findings"                     values:"text,json,sarif"`
	Enable  string `default:""     flag:"enable"  help:"Comma-separated rule IDs to run (default: all)"`
	Disable string `default:""     flag:"disable" help:"Comma-separated rule IDs to skip"`
}

type lintFlags struct {
	Dir           string `default:"."     flag:"dir"            help:"Directory to lint"`
	Baseline      string `default:""      flag:"baseline"       help:"Baseline file for ratchet mode (default: .cmdguard-lint-baseline.json in the linted dir when present)"`
	WriteBaseline bool   `default:"false" flag:"write-baseline" help:"Write current findings as the new baseline, then exit 0"`
}

type rulesFlags struct {
	Markdown bool `default:"false" flag:"markdown" help:"Print the rule table in Markdown (README generation)"`
}

func main() {
	cli, err := newApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}

// newApp builds the CLI. Split from main for in-process testing via
// ExecuteWithArgs.
func newApp() (*v4.CLI[cliConfig], error) {
	cli, err := v4.NewCLI(
		"cmdguard-lint",
		"cmdguard usage linter",
		cliConfig{Output: "text", Enable: "", Disable: ""},
		v4.WithCLIVersion(version.Version),
		v4.WithSignalHandling(),
	)
	if err != nil {
		return nil, fmt.Errorf("building CLI: %w", err)
	}

	lintCmd, err := v4.NewCommand(
		"lint",
		lintFlags{Dir: ".", Baseline: "", WriteBaseline: false},
		runLint,
		v4.WithShort("Lint a directory for cmdguard usage anti-patterns"),
		v4.WithExample("cmdguard-lint lint --dir . --output sarif"),
		v4.WithExample("cmdguard-lint lint --disable CG004,CG006"),
		v4.WithExample("cmdguard-lint lint --write-baseline"),
	)
	if err != nil {
		return nil, fmt.Errorf("building lint command: %w", err)
	}

	rulesCmd, err := v4.NewCommand(
		"rules",
		rulesFlags{Markdown: false},
		runRules,
		v4.WithShort("List available rules"),
		v4.WithExample("cmdguard-lint rules --markdown"),
	)
	if err != nil {
		return nil, fmt.Errorf("building rules command: %w", err)
	}

	if err := v4.AddCommand(cli, lintCmd); err != nil {
		return nil, fmt.Errorf("registering lint command: %w", err)
	}

	if err := v4.AddCommand(cli, rulesCmd); err != nil {
		return nil, fmt.Errorf("registering rules command: %w", err)
	}

	return cli, nil
}

// errFindings is returned when the lint run produced findings; cmdguard
// displays it once and ExecuteAndExit maps it to exit code 1.
var errFindings = errors.New("lint run produced findings")

func runLint(ctx context.Context, cfg *cliConfig, flags lintFlags) error {
	findings, err := detectFiltered(ctx, flags.Dir, cfg)
	if err != nil {
		return err
	}

	baselinePath, baseline, err := resolveBaseline(flags)
	if err != nil {
		return err
	}

	if flags.WriteBaseline {
		if err := lint.WriteBaseline(findings, baselinePath); err != nil {
			return fmt.Errorf("writing baseline: %w", err)
		}

		fmt.Fprintf(os.Stderr, "baseline written: %s (%d accepted finding(s))\n", baselinePath, len(findings))

		return nil
	}

	var stale []lint.BaselineEntry

	if baseline != nil {
		findings, stale = lint.ApplyBaseline(findings, baseline)

		if count := len(stale); count > 0 {
			fmt.Fprintf(os.Stderr,
				"baseline: %d stale entr%s (fixed finding(s)) — tighten with --write-baseline\n",
				count, pluralY(count),
			)
		}
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
		return fmt.Errorf(
			"%d cmdguard-lint finding(s) (suppress with //cmdguard-lint:ignore <RULE> <reason>): %w",
			len(findings),
			errFindings,
		)
	}

	return nil
}

// resolveBaseline loads the baseline for ratchet mode: the explicit --baseline
// path when given (a missing explicit file is an error), otherwise the default
// location inside the linted directory when it exists, otherwise no baseline.
func resolveBaseline(flags lintFlags) (string, *lint.Baseline, error) {
	path := flags.Baseline

	if path == "" {
		path = filepath.Join(flags.Dir, lint.DefaultBaselinePath)
	}

	baseline, found, err := lint.LoadBaseline(path)
	if err != nil {
		return "", nil, fmt.Errorf("loading baseline: %w", err)
	}

	if !found && flags.Baseline != "" {
		return "", nil, fmt.Errorf("%w: %s (create it with --write-baseline)", lint.ErrBaselineNotFound, flags.Baseline)
	}

	return path, baseline, nil
}

// pluralY returns "y" or "ies" for the stale-entry message.
func pluralY(n int) string {
	if n == 1 {
		return "y"
	}

	return "ies"
}

// detectFiltered runs lint.Detect and applies the enable/disable filters,
// keeping the single-pass fast path. Filtering itself lives in
// lint.FilterByIDs so the CLI flags and the provider options share semantics.
func detectFiltered(ctx context.Context, dir string, cfg *cliConfig) ([]finding.Finding, error) {
	findings, err := lint.Detect(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("detecting cmdguard usage findings in %s: %w", dir, err)
	}

	return lint.FilterByIDs(findings, cfg.Enable, cfg.Disable), nil
}

func runRules(_ context.Context, _ *cliConfig, flags rulesFlags) error {
	if flags.Markdown {
		if _, err := fmt.Fprint(os.Stdout, lint.RuleTableMarkdown()); err != nil {
			return fmt.Errorf("printing rule table: %w", err)
		}

		return nil
	}

	for _, rule := range lint.AllRules() {
		if _, err := fmt.Fprintf(
			os.Stdout,
			"%s  %-28s %s\n",
			rule.Meta.ID,
			rule.Meta.Name,
			rule.Meta.Description,
		); err != nil {
			return fmt.Errorf("printing rule: %w", err)
		}
	}

	if _, err := fmt.Fprintln(
		os.Stdout,
		"\nSuppress with //cmdguard-lint:ignore <RULE> <reason> on the offending line or the line above.",
	); err != nil {
		return fmt.Errorf("printing suppression hint: %w", err)
	}

	return nil
}
