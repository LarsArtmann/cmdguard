package lint

import (
	"context"
	"fmt"
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// Rule ID constants. IDs are stable once shipped: suppressions
// (//cmdguard-lint:ignore <ID> <reason>) and --enable/--disable filters key on
// them. Never renumber.
const (
	RuleExecuteBypass        = "CG001"
	RuleStaleMajorImport     = "CG002"
	RulePanicOnConstructor   = "CG003"
	RuleSetVersionRuntime    = "CG004"
	RuleDuplicateVersionOpts = "CG005"
	RuleExecuteErrorReprint  = "CG006"
)

// ToolName is the finding attribution for every rule in this linter.
const ToolName = "cmdguard-lint"

// AllRules returns every rule in registration order. Filter with
// linter.FilterRules for --enable/--disable semantics.
func AllRules() []linter.RuleFunc {
	return []linter.RuleFunc{
		executeBypassRule,
		staleMajorImportRule,
		panicOnConstructorRule,
		setVersionRuntimeRule,
		duplicateVersionOptsRule,
		executeErrorReprintRule,
	}
}

// NewRegistry builds a registry with every rule registered, stamped with the
// cmdguard-lint tool name. This is the entry point for the go-linter-sdk
// execution paths (Registry.Run, DetectorFromRegistry).
func NewRegistry() *linter.Registry {
	registry := linter.NewRegistry(linter.WithToolName(ToolName))

	for _, rule := range AllRules() {
		registry.Register(rule)
	}

	return registry
}

// Detect runs every default rule against dir in a single analysis pass and
// returns findings with in-source suppressions applied. This is the canonical
// fast path used by the provider package and the CLI; the per-rule registry
// path re-analyzes per rule but yields the same findings.
func Detect(ctx context.Context, dir string) ([]finding.Finding, error) {
	proj, err := analyze(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("cmdguard-lint: analyzing %s: %w", dir, err)
	}

	checks := []func(*project) []finding.Finding{
		checkExecuteBypass,
		checkStaleMajorImport,
		checkPanicOnConstructor,
		checkSetVersionRuntime,
		checkDuplicateVersionOpts,
		checkExecuteErrorReprint,
	}

	var findings []finding.Finding

	for _, check := range checks {
		findings = append(findings, check(proj)...)
	}

	return applySuppressions(proj, findings), nil
}

// ruleFor wraps a check into a go-linter-sdk RuleFunc with stable identity.
func ruleFor(meta linter.RuleMeta, check func(*project) []finding.Finding) linter.RuleFunc {
	meta.ToolName = ToolName

	return linter.RuleFunc{
		Meta: meta,
		Run: func(ctx context.Context, dir string) ([]finding.Finding, error) {
			proj, err := analyze(ctx, dir)
			if err != nil {
				return nil, fmt.Errorf("cmdguard-lint: analyzing %s: %w", dir, err)
			}

			return applySuppressions(proj, check(proj)), nil
		},
	}
}

// pos builds a finding position for a node in a scanned file.
func (f *sourceFile) pos(node ast.Node) finding.Position {
	position := f.fset.Position(node.Pos())

	return finding.Pos(finding.FilePath(f.relPath), position.Line, position.Column)
}
