package lint

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// checkMissingShortDescription detects NewCommand / NewParentCommand calls
// without a WithShort option (rule CG007).
//
// Origin: the v5 end-goal plan (docs/planning/2026-11-01_v5-end-goal-paved-road.md,
// workstream W5). A command without a short description renders as a blank
// line in `--help`, is invisible in command lists, and fails
// WithStrictValidation. cmdguard's constructors cannot reject it in Lenient
// mode without breaking existing consumers, so the linter carries the
// baseline until v5 flips the default.
//
// Detection is syntactic: the command options are literal arguments of the
// constructor call. A spread option argument (`opts...`) makes the call
// unknowable and is skipped rather than guessed.
func checkMissingShortDescription(proj *project, meta linter.RuleMeta) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		if strings.HasSuffix(file.relPath, "_test.go") {
			continue
		}

		cmdguardPath := firstCmdguardPath(file.imports)
		if cmdguardPath == "" {
			continue
		}

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if !file.isSelectorCall(call, cmdguardPath, "NewCommand") &&
				!file.isSelectorCall(call, cmdguardPath, "NewParentCommand") {
				return true
			}

			if callHasSpreadOption(call) || hasShortOption(file, cmdguardPath, call) {
				return true
			}

			findings = append(findings, newFinding(
				meta,
				"command "+commandLabel(call)+" has no WithShort: it renders as a blank line in --help and fails WithStrictValidation",
				file.pos(call),
			).
				WithConfidence(finding.ConfidenceHigh).
				WithSuggestion("Add cmdguard.WithShort(\"one-line description\") to the command options").
				MustBuild())

			return true
		})
	}

	return findings
}

// checkMissingCLIVersion detects NewCLI calls without any version source
// (rule CG008).
//
// Origin: the v5 end-goal plan, workstream W5. A CLI without a version cannot
// answer `--version`, cannot be identified in bug reports, and cannot be
// matched to a release in operator logs. cmdguard cannot require it at
// construction without breaking every existing consumer, so the linter reports
// the baseline.
//
// Detection accepts either version source: WithCLIVersion, or
// WithFangOptions(fang.WithVersion(...)). A spread option argument (`opts...`)
// makes the call unknowable and is skipped rather than guessed.
func checkMissingCLIVersion(proj *project, meta linter.RuleMeta) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]
		if strings.HasSuffix(file.relPath, "_test.go") {
			continue
		}

		cmdguardPath := firstCmdguardPath(file.imports)
		if cmdguardPath == "" {
			continue
		}

		ast.Inspect(file.file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !file.isSelectorCall(call, cmdguardPath, "NewCLI") {
				return true
			}

			if callHasSpreadOption(call) || callHasAnyVersionSource(file, cmdguardPath, call) {
				return true
			}

			findings = append(findings, newFinding(
				meta,
				"NewCLI has no version source: the binary cannot report which build it is, and --version is unavailable",
				file.pos(call),
			).
				WithConfidence(finding.ConfidenceHigh).
				WithSuggestion("Pass cmdguard.WithCLIVersion(version.Version) to NewCLI (and register VersionCommand for a subcommand)").
				MustBuild())

			return true
		})
	}

	return findings
}

// callHasSpreadOption reports whether the call passes a variadic spread
// argument (`opts...`), which hides the option set from syntactic analysis.
func callHasSpreadOption(call *ast.CallExpr) bool {
	return call.Ellipsis.IsValid()
}

// hasShortOption reports whether the constructor call passes WithShort with a
// non-empty description.
func hasShortOption(file *sourceFile, cmdguardPath string, call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		optCall, ok := arg.(*ast.CallExpr)
		if !ok || !file.isSelectorCall(optCall, cmdguardPath, "WithShort") {
			continue
		}

		if len(optCall.Args) == 0 {
			continue
		}

		if lit, ok := optCall.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if description, err := strconv.Unquote(lit.Value); err == nil && strings.TrimSpace(description) == "" {
				continue
			}
		}

		return true
	}

	return false
}

// callHasAnyVersionSource reports whether the NewCLI call supplies a version
// through WithCLIVersion or a fang WithVersion option.
func callHasAnyVersionSource(file *sourceFile, cmdguardPath string, call *ast.CallExpr) bool {
	for _, arg := range call.Args {
		optCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}

		if file.isSelectorCall(optCall, cmdguardPath, "WithCLIVersion") {
			return true
		}

		if file.isSelectorCall(optCall, cmdguardPath, "WithFangOptions") &&
			containsFangVersion(file, optCall.Args) {
			return true
		}
	}

	return false
}

// commandLabel renders a quoted command name for a finding message, falling
// back to a generic label when the Use argument is not a string literal.
func commandLabel(call *ast.CallExpr) string {
	if len(call.Args) == 0 {
		return "command"
	}

	lit, ok := call.Args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "command"
	}

	use, err := strconv.Unquote(lit.Value)
	if err != nil || use == "" {
		return "command"
	}

	return strconv.Quote(use)
}
