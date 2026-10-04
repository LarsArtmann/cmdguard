package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// panicOnConstructorMeta is the identity header of the panic rule.
//
// Origin: timesheets deep dive F6. cmdguard eliminated Must* constructors by
// design: every function returns errors and the failure mode of a panic at
// process start from deep inside a factory is exactly what v4 removed. A
// must-style wrapper re-introduces it. Constructors only fail on invalid
// static arguments, so probability is low — but the failure mode is a startup
// panic instead of a checkable error.
var panicOnConstructorMeta = linter.RuleMeta{
	ID:          RulePanicOnConstructor,
	Name:        "panic on constructor error",
	Description: "panicking on a cmdguard constructor error re-introduces the panic cmdguard removed by design; return the error instead",
	Cat:         linter.CategoryCorrectness,
	Sev:         finding.SeverityError,
}

var panicOnConstructorRule = ruleFor(panicOnConstructorMeta, checkPanicOnConstructor)

func checkPanicOnConstructor(proj *project) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]

		if firstCmdguardPath(file.imports) == "" {
			continue
		}

		forFuncs(file.file, func(body *ast.BlockStmt) {
			errIdents := file.cmdguardConstructorErrIdents(body)

			ast.Inspect(body, func(n ast.Node) bool {
				ifStmt, ok := n.(*ast.IfStmt)
				if !ok {
					return true
				}

				errName, ok := condIsErrNotNil(ifStmt.Cond)
				if !ok || !errIdents[errName] {
					return true
				}

				for _, stmt := range ifStmt.Body.List {
					expr, ok := stmt.(*ast.ExprStmt)
					if !ok {
						continue
					}

					call, ok := expr.X.(*ast.CallExpr)
					if !ok || call.Fun == nil {
						continue
					}

					ident, ok := call.Fun.(*ast.Ident)
					if !ok || ident.Name != "panic" {
						continue
					}

					findings = append(findings, newFinding(panicOnConstructorMeta,
						"panic on a cmdguard constructor error re-introduces the panic cmdguard removed by design: constructors return errors so registration failures surface as checkable errors",
						file.pos(stmt),
					).
						WithConfidence(finding.ConfidenceFull).
						WithSuggestion("Return the error (the command registration table can check it once) instead of panicking").
						MustBuild())
				}

				return true
			})
		})
	}

	return findings
}
