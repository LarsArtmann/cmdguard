// Package provider registers cmdguard-lint as a BuildFlow DAG tool via the
// go-finding toolsdk contract. BuildFlow discovers this Spec (toolsdk.All)
// whenever it imports this package, turning the cmdguard usage linter into a
// first-class pipeline tool without BuildFlow-side glue.
package provider

import (
	"context"

	"github.com/larsartmann/cmdguard/lint"
	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
)

// Provider is the BuildFlow tool spec for cmdguard-lint. Detection runs all
// rules in a single analysis pass over the working directory (resolved via
// finding.WorkingDirFromContext, "." when unset). ModuleFanOut declares
// per-module fan-out so multi-module workspaces lint every module; findings
// carry repo-root-relative paths from each module's own pass.
//
//nolint:gochecknoglobals // toolsdk contract: specs self-register as package-level vars
var Provider = toolsdk.Register(toolsdk.Spec{
	Name:         "cmdguard-lint",
	Description:  "cmdguard usage linter: execute bypass, stale majors, constructor panics, version/option misuse, double error display",
	Trigger:      toolsdk.OnGoModule(),
	DependsOn:    nil,
	ModuleFanOut: true,
	Inputs:       []string{"**/*.go"},
	Options:      nil,
	Detect: finding.NamedDetectorFunc("cmdguard-lint", func(ctx context.Context) ([]finding.Finding, error) {
		dir := finding.WorkingDirFromContext(ctx)
		if dir == "" {
			dir = "."
		}

		return lint.Detect(ctx, dir)
	}),
	Repair:      nil,
	HealthCheck: nil,
})
