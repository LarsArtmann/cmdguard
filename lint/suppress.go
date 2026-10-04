package lint

import (
	"regexp"
	"strings"

	"github.com/larsartmann/go-finding"
)

// suppressionDirective is the in-source suppression comment format:
//
//	//cmdguard-lint:ignore CG001 <required reason>
//
// The directive suppresses the given rule for a finding positioned on the
// same line, or on the line directly below the directive (the common
// "directive above the call" style). A directive without a reason is invalid
// and suppresses nothing.
const suppressionDirective = "//cmdguard-lint:ignore"

var suppressionPattern = regexp.MustCompile( //nolint:gochecknoglobals // compiled once, immutable
	`^//cmdguard-lint:ignore\s+([A-Z]+[0-9]+)\s+(\S.*)$`,
)

// applySuppressions drops findings whose position carries an active
// suppression directive.
func applySuppressions(proj *project, findings []finding.Finding) []finding.Finding {
	if len(findings) == 0 {
		return nil
	}

	directives := map[string]map[int]string{} // relPath -> line -> ruleID

	for i := range proj.files {
		file := &proj.files[i]

		for lineIdx, line := range file.lines {
			ruleID, ok := parseSuppressionLine(line)
			if !ok {
				continue
			}

			if directives[file.relPath] == nil {
				directives[file.relPath] = map[int]string{}
			}

			// Directive applies to its own line (0-based index) and the line
			// below it.
			directives[file.relPath][lineIdx+1] = ruleID
			directives[file.relPath][lineIdx+2] = ruleID
		}
	}

	kept := make([]finding.Finding, 0, len(findings))

	for _, f := range findings {
		if f.Position.File == "" || f.Position.Line <= 0 {
			kept = append(kept, f)

			continue
		}

		ruleID, active := directives[string(f.Position.File)][f.Position.Line]
		if active && ruleID == string(f.Rule) {
			continue
		}

		kept = append(kept, f)
	}

	return kept
}

// parseSuppressionLine extracts (ruleID, ok) from a source line carrying a
// valid suppression directive. Lines without the directive prefix, without a
// rule ID, or without a reason do not suppress.
func parseSuppressionLine(line string) (string, bool) {
	idx := strings.Index(line, suppressionDirective)
	if idx < 0 {
		return "", false
	}

	match := suppressionPattern.FindStringSubmatch(strings.TrimSpace(line[idx:]))
	if match == nil {
		return "", false
	}

	return match[1], true
}
