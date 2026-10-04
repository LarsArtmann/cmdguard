package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// executeErrorReprintMeta is the identity header of the error-reprint rule.
//
// Origin: cmdguard's error/exit contract (v4). The error returned by
// cli.Execute has already been displayed exactly once (fang when enabled,
// cobra otherwise); it exists for exit-code mapping. Re-printing it with the
// fmt Print family double-reports every failure.
//
// Detection: an identifier assigned from <cli>.Execute(...) where <cli> was
// assigned from NewCLI in the same file, later passed as an argument to
// fmt.Print/Printf/Println/Fprint/Fprintf/Fprintln. fmt.Errorf wrapping is
// deliberately NOT flagged: mapping errors into wrapped errors is part of
// the supported exit-code path.
var executeErrorReprintMeta = linter.RuleMeta{
	ID:          RuleExecuteErrorReprint,
	Name:        "execute error reprinted",
	Description: "the error returned by cli.Execute is already displayed by cmdguard; re-printing it double-reports the failure",
	Cat:         CategoryUsage,
	Sev:         finding.SeverityWarning,
}

var executeErrorReprintRule = ruleFor(executeErrorReprintMeta, checkExecuteErrorReprint)

func checkExecuteErrorReprint(proj *project) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		cliVars := file.cliVars()

		if len(cliVars) == 0 {
			continue
		}

		execErrIdents := collectExecuteErrIdents(file, cliVars)

		if len(execErrIdents) == 0 {
			continue
		}

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isFmtPrintCall(file, call) {
				return true
			}

			for _, arg := range call.Args {
				ident, ok := arg.(*ast.Ident)
				if !ok {
					continue
				}

				if !execErrIdents[ident.Name] {
					continue
				}

				findings = append(findings, newFinding(executeErrorReprintMeta,
					"the error returned by cli.Execute has already been displayed exactly once by cmdguard; printing it again double-reports the failure",
					file.pos(call),
				).
					WithConfidence(finding.ConfidenceHigh).
					WithSuggestion("Map the error to an exit code instead: cli.ExecuteAndExit(ctx), or ExitCode(err) after your post-execution work").
					MustBuild())

				break
			}

			return true
		})
	}

	return findings
}

// collectExecuteErrIdents returns identifiers assigned from
// `<cliVar>.Execute(...)` calls for the given CLI variables.
func collectExecuteErrIdents(file *sourceFile, cliVars map[string]bool) map[string]bool {
	errs := map[string]bool{}

	ast.Inspect(file.file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}

		for i, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}

			for name := range cliVars {
				if !isMethodCallOn(call, name, "Execute") {
					continue
				}

				if ident, ok := assign.Lhs[i].(*ast.Ident); ok {
					errs[ident.Name] = true
				}
			}
		}

		return true
	})

	return errs
}

// isFmtPrintCall reports whether call is one of the fmt display functions
// (Print, Println, Printf, Fprint, Fprintln, Fprintf).
func isFmtPrintCall(file *sourceFile, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return false
	}

	switch sel.Sel.Name {
	case "Print", "Println", "Printf", "Fprint", "Fprintln", "Fprintf":
		return file.selectorFrom(sel, "fmt")
	default:
		return false
	}
}
