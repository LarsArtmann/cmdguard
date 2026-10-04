package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// setVersionRuntimeMeta is the identity header of the SetVersion rule.
//
// Origin: timesheets deep dive F7. cmdguard's WithCLIVersion wires fang's
// version machinery at construction time; calling cli.SetVersion after the
// fact is the half-migrated shape. The audit consumer also hand-rolled a
// version command cmdguard already ships (VersionCommand helper).
//
// Detection is dataflow-lite: the receiver identifier must be assigned from
// a cmdguard NewCLI call in the same file, keeping the rule precise without
// type information.
var setVersionRuntimeMeta = linter.RuleMeta{
	ID:          RuleSetVersionRuntime,
	Name:        "runtime SetVersion",
	Description: "cli.SetVersion mutates after construction; pass WithCLIVersion to NewCLI (and use the VersionCommand helper) instead",
	Cat:         CategoryUsage,
	Sev:         finding.SeverityWarning,
}

var setVersionRuntimeRule = ruleFor(setVersionRuntimeMeta, checkSetVersionRuntime)

func checkSetVersionRuntime(proj *project) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		cliVars := file.cliVars()

		if len(cliVars) == 0 {
			continue
		}

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || sel.Sel.Name != "SetVersion" {
				return true
			}

			for name := range cliVars {
				if !isMethodCallOn(call, name, "SetVersion") {
					continue
				}

				findings = append(findings, newFinding(setVersionRuntimeMeta,
					"cli.SetVersion patches the version after construction: pass WithCLIVersion(version) to NewCLI so fang's version wiring is complete from the start",
					file.pos(call),
				).
					WithConfidence(finding.ConfidenceHigh).
					WithSuggestion("Move the version into NewCLI options: cmdguard.WithCLIVersion(version) (WithCLICommit for the commit stamp)").
					MustBuild())

				break
			}

			return true
		})
	}

	return findings
}
