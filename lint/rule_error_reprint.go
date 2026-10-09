package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// checkExecuteErrorReprint detects fmt Print calls on the error returned by
// CLI Execute, double-displaying it (rule CG006).
//
// Origin: cmdguard's error/exit contract (v4). The error returned by
// cli.Execute has already been displayed exactly once (fang when enabled,
// cobra otherwise); it exists for exit-code mapping. Re-printing it with the
// fmt Print family double-reports every failure.
//
// Detection: an identifier assigned from <cli>.Execute(...) where <cli> was
// assigned from NewCLI in the same file OR as a package-level var anywhere
// in the package (cross-file name tracing, see crossfile.go), later passed as
// an argument to fmt.Print/Printf/Println/Fprint/Fprintf/Fprintln.
// fmt.Errorf wrapping is deliberately NOT flagged: mapping errors into
// wrapped errors is part of the supported exit-code path.
func checkExecuteErrorReprint(proj *project, meta linter.RuleMeta) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		cliVars := proj.cliVarsFor(file)

		if len(cliVars) == 0 {
			continue
		}

		execErrIdents := proj.execErrIdentsFor(file, cliVars)

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

				findings = append(findings, newFinding(
					meta,
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
