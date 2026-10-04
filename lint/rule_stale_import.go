package lint

import (
	"fmt"
	"go/ast"

	"github.com/larsartmann/go-finding"
	linter "github.com/larsartmann/go-linter-sdk"
)

// staleMajorImportMeta is the identity header of the stale-major-import rule.
//
// Origin: timesheets deep dive F2 (critical). Older majors are frozen: they
// receive no fixes, and the cmdguard sub-modules (glamour, prompts, spinner,
// telemetry, flightrecorder) only target the current major, so their features
// are unreachable from a stale import. The v3->v4 migration is mechanical
// (import-path rename), which makes staying behind purely a loss.
var staleMajorImportMeta = linter.RuleMeta{
	ID:          RuleStaleMajorImport,
	Name:        "stale cmdguard major",
	Description: "importing a frozen cmdguard major blocks all fixes and sub-module features; migrate to the current major",
	Cat:         CategoryUsage,
	Sev:         finding.SeverityError,
	ToolName: ToolName,
}

var staleMajorImportRule = ruleFor(staleMajorImportMeta, checkStaleMajorImport)

func checkStaleMajorImport(proj *project) []finding.Finding {
	var findings []finding.Finding

	for i := range proj.files {
		file := &proj.files[i]

		for _, spec := range file.file.Imports {
			major := majorOf(importPathOf(spec))

			// major "" covers sub-module paths, which track the current
			// line; only versioned core imports can go stale.
			if major == "" || major == CurrentMajor {
				continue
			}

			findings = append(findings, newFinding(staleMajorImportMeta,
				fmt.Sprintf("cmdguard %s is frozen (no further fixes; sub-modules unreachable): migrate the import to %s", major, CurrentMajor),
				file.pos(spec),
			).
				WithConfidence(finding.ConfidenceFull).
				WithSuggestion(fmt.Sprintf("go get github.com/larsartmann/cmdguard/%s@latest and rename the import path", CurrentMajor)).
				MustBuild())
		}
	}

	return findings
}

// importPathOf unwraps an ImportSpec's quoted path.
func importPathOf(spec *ast.ImportSpec) string {
	if spec.Path == nil {
		return ""
	}

	return trimQuotes(spec.Path.Value)
}

// trimQuotes removes surrounding double quotes from a Go string literal.
func trimQuotes(value string) string {
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}

	return value
}
