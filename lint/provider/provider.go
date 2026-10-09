// Package provider registers cmdguard-lint as a BuildFlow DAG tool via the
// go-finding toolsdk contract. BuildFlow discovers this Spec (toolsdk.All)
// whenever it imports this package, turning the cmdguard usage linter into a
// first-class pipeline tool without BuildFlow-side glue.
package provider

import (
	"context"
	"fmt"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"

	"github.com/larsartmann/cmdguard/lint"
)

// Provider is the BuildFlow tool spec for cmdguard-lint. Detection runs all
// rules in a single analysis pass over the working directory (resolved via
// finding.WorkingDirFromContext, "." when unset). ModuleFanOut declares
// per-module fan-out so multi-module workspaces lint every module; findings
// carry repo-root-relative paths from each module's own pass.
//
// Options (set per run via toolsdk.WithOptions by the consumer):
//
//	enable  (string, default ""): comma-separated rule IDs to run (default all)
//	disable (string, default ""): comma-separated rule IDs to skip
//
// Both use the same comma-separated syntax as the CLI's --enable/--disable
// flags (filtering is lint.FilterByIDs), so a config tuned for one surface
// transfers verbatim to the other.
//
//nolint:gochecknoglobals // toolsdk contract: specs self-register as package-level vars
var Provider = toolsdk.Register(toolsdk.Spec{
	Name:         "cmdguard-lint",
	Description:  "cmdguard usage linter: execute bypass, stale majors, constructor panics, version/option misuse, double error display",
	Trigger:      toolsdk.OnGoModule(),
	DependsOn:    nil,
	ModuleFanOut: true,
	Inputs:       []string{"**/*.go"},
	Options: []toolsdk.Option{
		{
			Name:        "enable",
			Kind:        toolsdk.OptionKindString,
			Default:     "",
			Description: "Comma-separated rule IDs to run (default: all)",
		},
		{
			Name:        "disable",
			Kind:        toolsdk.OptionKindString,
			Default:     "",
			Description: "Comma-separated rule IDs to skip",
		},
	},
	Detect: finding.NamedDetectorFunc("cmdguard-lint", func(ctx context.Context) ([]finding.Finding, error) {
		dir := finding.WorkingDirFromContext(ctx)
		if dir == "" {
			dir = "."
		}

		findings, err := lint.Detect(ctx, dir)
		if err != nil {
			return nil, fmt.Errorf("cmdguard-lint detection: %w", err)
		}

		return lint.FilterByIDs(findings, optionString(ctx, "enable"), optionString(ctx, "disable")), nil
	}),
	Repair:      nil,
	HealthCheck: nil,
})

// optionString reads a declared string option from the run context, returning
// "" for unset options or values of an unexpected kind (defensive: the SDK
// already kind-checks consumer values in ValidateOptions).
func optionString(ctx context.Context, name string) string {
	values, ok := toolsdk.OptionsFromContext(ctx)
	if !ok {
		return ""
	}

	value, _ := values[name].(string)

	return value
}
