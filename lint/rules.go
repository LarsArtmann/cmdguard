package lint

import (
	"context"
	"fmt"
	"go/ast"
	"strings"

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

// CategoryUsage is this linter's custom rule category: findings about how a
// codebase consumes (or under-consumes) cmdguard, as opposed to defects in
// the code itself. The go-linter-sdk explicitly supports domain categories.
const CategoryUsage = linter.Category("usage")

// ruleDef couples a rule's declarative identity with its check. The table in
// allRuleDefs is the single source of truth: AllRules, NewRegistry, and
// Detect all derive from it.
type ruleDef struct {
	meta  linter.RuleMeta
	check func(proj *project, meta linter.RuleMeta) []finding.Finding
}

// allRuleDefs returns every rule with its identity header and detection
// closure. Each check's doc comment carries the rule's origin and trust
// signals; the meta is passed into the check so findings are stamped without
// package-level rule variables.
func allRuleDefs() []ruleDef {
	return []ruleDef{
		{
			meta: linter.RuleMeta{
				ID:          RuleExecuteBypass,
				Name:        "execute bypass",
				Description: "fang.Execute runs the raw cobra tree, bypassing cli.Execute and its signal handling, graceful shutdown, cleanup hooks, and single-error-display contract",
				Cat:         linter.CategoryCorrectness,
				Sev:         finding.SeverityCritical,
				ToolName:    ToolName,
			},
			check: checkExecuteBypass,
		},
		{
			meta: linter.RuleMeta{
				ID:          RuleStaleMajorImport,
				Name:        "stale cmdguard major",
				Description: "importing a frozen cmdguard major blocks all fixes and sub-module features; migrate to the current major",
				Cat:         CategoryUsage,
				Sev:         finding.SeverityError,
				ToolName:    ToolName,
			},
			check: checkStaleMajorImport,
		},
		{
			meta: linter.RuleMeta{
				ID:          RulePanicOnConstructor,
				Name:        "panic on constructor error",
				Description: "panicking on a cmdguard constructor error re-introduces the panic cmdguard removed by design; return the error instead",
				Cat:         linter.CategoryCorrectness,
				Sev:         finding.SeverityError,
				ToolName:    ToolName,
			},
			check: checkPanicOnConstructor,
		},
		{
			meta: linter.RuleMeta{
				ID:          RuleSetVersionRuntime,
				Name:        "runtime SetVersion",
				Description: "cli.SetVersion mutates after construction; pass WithCLIVersion to NewCLI (and use the VersionCommand helper) instead",
				Cat:         CategoryUsage,
				Sev:         finding.SeverityWarning,
				ToolName:    ToolName,
			},
			check: checkSetVersionRuntime,
		},
		{
			meta: linter.RuleMeta{
				ID:          RuleDuplicateVersionOpts,
				Name:        "duplicate version options",
				Description: "WithCLIVersion and WithFangOptions(fang.WithVersion(...)) in one NewCLI call pass duplicate fang version options",
				Cat:         linter.CategoryCorrectness,
				Sev:         finding.SeverityError,
				ToolName:    ToolName,
			},
			check: checkDuplicateVersionOpts,
		},
		{
			meta: linter.RuleMeta{
				ID:          RuleExecuteErrorReprint,
				Name:        "execute error reprinted",
				Description: "the error returned by cli.Execute is already displayed by cmdguard; re-printing it double-reports the failure",
				Cat:         CategoryUsage,
				Sev:         finding.SeverityWarning,
				ToolName:    ToolName,
			},
			check: checkExecuteErrorReprint,
		},
	}
}

// AllRules returns every rule in table order. Filter with linter.FilterRules
// for --enable/--disable semantics.
func AllRules() []linter.RuleFunc {
	defs := allRuleDefs()

	rules := make([]linter.RuleFunc, 0, len(defs))

	for _, def := range defs {
		rules = append(rules, ruleFor(def.meta, def.check))
	}

	return rules
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

// Detect runs every rule against dir in a single analysis pass and returns
// findings with in-source suppressions applied. This is the canonical fast
// path used by the provider package and the CLI; the per-rule registry path
// re-analyzes per rule but yields the same findings.
func Detect(ctx context.Context, dir string) ([]finding.Finding, error) {
	proj, err := analyze(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("cmdguard-lint: analyzing %s: %w", dir, err)
	}

	var findings []finding.Finding

	for _, def := range allRuleDefs() {
		findings = append(findings, def.check(proj, def.meta)...)
	}

	return applySuppressions(proj, findings), nil
}

// FilterByIDs applies --enable/--disable semantics to a finding set: findings
// from disabled rules are dropped; when enable is non-empty, only findings
// from enabled rules survive. Both take comma-separated rule IDs (spaces
// tolerated, empty segments ignored) — the same syntax as the CLI flags and
// the provider's enable/disable options, so all three surfaces stay uniform.
func FilterByIDs(findings []finding.Finding, enable, disable string) []finding.Finding {
	enabled := idSet(enable)
	disabled := idSet(disable)

	kept := make([]finding.Finding, 0, len(findings))

	for _, candidate := range findings {
		if disabled[string(candidate.Rule)] {
			continue
		}

		if len(enabled) > 0 && !enabled[string(candidate.Rule)] {
			continue
		}

		kept = append(kept, candidate)
	}

	return kept
}

// idSet splits a comma-separated rule-ID list into a set.
func idSet(value string) map[string]bool {
	set := map[string]bool{}

	for part := range strings.SplitSeq(value, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			set[trimmed] = true
		}
	}

	return set
}

// ruleFor wraps a check into a go-linter-sdk RuleFunc with stable identity.
// Checks build findings via newFinding (not the rule variable) to avoid Go
// initialization cycles.
func ruleFor(meta linter.RuleMeta, check func(*project, linter.RuleMeta) []finding.Finding) linter.RuleFunc {
	return linter.RuleFunc{
		Meta: meta,
		Run: func(ctx context.Context, dir string) ([]finding.Finding, error) {
			proj, err := analyze(ctx, dir)
			if err != nil {
				return nil, fmt.Errorf("cmdguard-lint: analyzing %s: %w", dir, err)
			}

			return applySuppressions(proj, check(proj, meta)), nil
		},
	}
}

// newFinding builds a finding builder pre-stamped with the rule's identity.
// Checks use this instead of the rule variable so package initialization
// stays acyclic.
func newFinding(meta linter.RuleMeta, message string, pos finding.Position) *finding.Builder {
	builder := finding.NewBuilder(finding.RuleName(meta.ID), meta.ToolName, message, meta.Sev, pos)

	if meta.Cat != "" {
		builder = builder.WithCategory(finding.Category(meta.Cat))
	}

	return builder
}

// pos builds a finding position for a node in a scanned file.
func (f *sourceFile) pos(node ast.Node) finding.Position {
	position := f.fset.Position(node.Pos())

	return finding.Pos(finding.FilePath(f.relPath), position.Line, position.Column)
}
