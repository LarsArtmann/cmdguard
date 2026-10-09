package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// checkPanicOnConstructor detects panic calls inside `if err != nil`
// branches where err traces to a cmdguard constructor (rule CG003).
//
// Origin: timesheets deep dive F6. cmdguard eliminated Must* constructors by
// design: every function returns errors and the failure mode of a panic at
// process start from deep inside a factory is exactly what v4 removed. A
// must-style wrapper re-introduces it. Constructors only fail on invalid
// static arguments, so probability is low — but the failure mode is a startup
// panic instead of a checkable error.
func checkPanicOnConstructor(proj *project, meta linter.RuleMeta) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]

		// Package-level constructor errors (var x, err = NewCLI(...) in any
		// file of the package) are visible to every function body; body-local
		// assignments stay primary. See crossfile.go. A file without its own
		// cmdguard import still qualifies when the package declares one.
		packageErrs := proj.constructorErrIdentsFor(file)

		if firstCmdguardPath(file.imports) == "" && len(packageErrs) == 0 {
			continue
		}

		forFuncs(file.file, func(body *ast.BlockStmt) {
			errIdents := file.cmdguardConstructorErrIdents(body)

			for name := range packageErrs {
				errIdents[name] = true
			}

			ast.Inspect(body, func(n ast.Node) bool {
				ifStmt, ok := n.(*ast.IfStmt)
				if !ok {
					return true
				}

				errName, ok := condIsErrNotNil(ifStmt.Cond)
				if !ok || !errIdents[errName] {
					return true
				}

				stmt := panicStatementIn(ifStmt.Body.List)
				if stmt == nil {
					return true
				}

				findings = append(findings, newFinding(
					meta,
					"panic on a cmdguard constructor error re-introduces the panic cmdguard removed by design: constructors return errors so registration failures surface as checkable errors",
					file.pos(stmt),
				).
					WithConfidence(finding.ConfidenceFull).
					WithSuggestion("Return the error (the command registration table can check it once) instead of panicking").
					MustBuild())

				return true
			})
		})
	}

	return findings
}

// panicStatementIn returns the first bare `panic(...)` expression statement
// in stmts, or nil when the branch does anything else.
func panicStatementIn(stmts []ast.Stmt) ast.Stmt {
	for _, stmt := range stmts {
		expr, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}

		call, ok := expr.X.(*ast.CallExpr)
		if !ok || call.Fun == nil {
			continue
		}

		ident, ok := call.Fun.(*ast.Ident)
		if ok && ident.Name == "panic" {
			return stmt
		}
	}

	return nil
}
