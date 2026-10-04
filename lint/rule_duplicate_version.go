package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// duplicateVersionOptsMeta is the identity header of the duplicate-version
// rule.
//
// Origin: cmdguard ADR-001 (fang integration). WithCLIVersion already pipes
// into fang.WithVersion; combining it with WithFangOptions(fang.WithVersion(...))
// in the same NewCLI call produces duplicate fang options, whose resolution
// order is undefined from cmdguard's perspective.
var duplicateVersionOptsMeta = linter.RuleMeta{
	ID:          RuleDuplicateVersionOpts,
	Name:        "duplicate version options",
	Description: "WithCLIVersion and WithFangOptions(fang.WithVersion(...)) in one NewCLI call pass duplicate fang version options",
	Cat:         linter.CategoryCorrectness,
	Sev:         finding.SeverityError,
	ToolName: ToolName,
}

var duplicateVersionOptsRule = ruleFor(duplicateVersionOptsMeta, checkDuplicateVersionOpts)

func checkDuplicateVersionOpts(proj *project) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		cmdguardPath := firstCmdguardPath(file.imports)

		if cmdguardPath == "" || file.imports[fangLocalName(file.imports)] != fangImportPath {
			continue
		}

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !file.isSelectorCall(call, cmdguardPath, "NewCLI") {
				return true
			}

			hasCLIVersion := false
			fangVersionInsideFangOpts := false

			for _, arg := range call.Args {
				optCall, ok := arg.(*ast.CallExpr)
				if !ok {
					continue
				}

				if file.isSelectorCall(optCall, cmdguardPath, "WithCLIVersion") {
					hasCLIVersion = true
				}

				if file.isSelectorCall(optCall, cmdguardPath, "WithFangOptions") && containsFangVersion(file, optCall.Args) {
					fangVersionInsideFangOpts = true
				}
			}

			if !hasCLIVersion || !fangVersionInsideFangOpts {
				return true
			}

			findings = append(findings, newFinding(duplicateVersionOptsMeta,
				"WithCLIVersion already pipes into fang.WithVersion: combining it with WithFangOptions(fang.WithVersion(...)) passes the version twice",
				file.pos(call),
			).
				WithConfidence(finding.ConfidenceFull).
				WithSuggestion("Keep exactly one version source: WithCLIVersion (preferred) or the fang option, not both (see cmdguard ADR-001)").
				MustBuild())

			return true
		})
	}

	return findings
}

// containsFangVersion reports whether any expression in exprs is a
// fang.WithVersion(...) call.
func containsFangVersion(file *sourceFile, exprs []ast.Expr) bool {
	for _, expr := range exprs {
		call, ok := expr.(*ast.CallExpr)
		if !ok {
			continue
		}

		if file.isSelectorCall(call, fangImportPath, "WithVersion") {
			return true
		}
	}

	return false
}

// fangLocalName resolves the local package name fang is imported under in
// this file ("" when fang is not imported).
func fangLocalName(imports map[string]string) string {
	for localName, importPath := range imports {
		if importPath == fangImportPath {
			return localName
		}
	}

	return ""
}
