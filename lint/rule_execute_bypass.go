package lint

import (
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// executeBypassRule detects fang.Execute receiving a cmdguard RootCommand().
//
// Origin: timesheets deep dive F1 (critical). fang only installs a signal
// context when the caller passes fang.WithNotifySignal; routing the cmdguard
// root command through raw fang.Execute means signal handling, graceful DI
// shutdown, cleanup hooks, NO_COLOR restore, and the single-error-display
// contract of cli.Execute never engage — carefully written shutdown paths
// become dead code on the signal path.
//
// Signals: (1) the project imports cmdguard anywhere, (2) a fang.Execute
// call, (3) an argument that is a RootCommand() call or an identifier
// assigned from one. All three together are deterministic; confidence is
// full.
var executeBypassRule = ruleFor(
	linter.RuleMeta{
		ID:          RuleExecuteBypass,
		Name:        "execute bypass",
		Description: "fang.Execute runs the raw cobra tree, bypassing cli.Execute and its signal handling, graceful shutdown, cleanup hooks, and single-error-display contract",
		Cat:         linter.CategoryCorrectness,
		Sev:         finding.SeverityCritical,
	},
	checkExecuteBypass,
)

func checkExecuteBypass(proj *project) []finding.Finding {
	if !proj.importsCmdguard {
		return nil
	}

	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		rootCommandIdents := collectRootCommandIdents(file.file)

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !file.isSelectorCall(call, fangImportPath, "Execute") {
				return true
			}

			if !file.callArgIsRootCommand(call, rootCommandIdents) {
				return true
			}

			findings = append(findings, executeBypassRule.NewFinding(
				"fang.Execute runs the raw cobra tree: cmdguard's cli.Execute is bypassed, so signal handling, graceful shutdown, cleanup hooks, and the single-error-display contract never engage",
				file.pos(call),
			).
				WithConfidence(finding.ConfidenceFull).
				WithSuggestion("Call cli.Execute(ctx) (add WithSignalHandling/WithGracefulShutdown as needed) so the cmdguard lifecycle layer runs").
				MustBuild())

			return true
		})
	}

	return findings
}
