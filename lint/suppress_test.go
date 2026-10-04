package lint

import (
	"testing"

	"github.com/larsartmann/go-finding"
)

func TestSuppressionSameLine(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"cmd.go": "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\" //cmdguard-lint:ignore CG002 frozen until the v4 migration PR lands\n\nvar _ = v3.NoFlags{}\n",
	})

	for _, f := range findings {
		if string(f.Rule) == RuleStaleMajorImport {
			t.Errorf(
				"expected CG002 suppressed by same-line directive, still reported at %s:%d",
				f.Position.File,
				f.Position.Line,
			)
		}
	}
}

func TestSuppressionLineAboveWithCmdguard(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"guard.go": "package main\n\nimport \"github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4\"\n\nvar _ = v4.NoFlags{}\n",
		"main.go": `package main

import (
	"charm.land/fang/v2"
)

func main() {
	//cmdguard-lint:ignore CG001 errorfamily owns display; verified in AUDIT-12
	_ = fang.Execute(nil, rootOfCLI())
}

func rootOfCLI() *rootCmd { return nil }

type rootCmd struct{}
`,
	})

	for _, f := range findings {
		if string(f.Rule) == RuleExecuteBypass {
			t.Errorf(
				"expected CG001 suppressed by directive above the call, still reported at %s:%d",
				f.Position.File,
				f.Position.Line,
			)
		}
	}
}

func TestSuppressionRequiresReason(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"cmd.go": "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{} //cmdguard-lint:ignore CG002\n",
	})

	assertFinding(t, findings, RuleStaleMajorImport, "cmd.go", 3)
}

func TestSuppressionRequiresMatchingRule(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"cmd.go": "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{} //cmdguard-lint:ignore CG006 wrong rule does not suppress\n",
	})

	assertFinding(t, findings, RuleStaleMajorImport, "cmd.go", 3)
}

func TestParseSuppressionLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		line  string
		rule  string
		valid bool
	}{
		{"x := 1 //cmdguard-lint:ignore CG001 deliberate", "CG001", true},
		{"//cmdguard-lint:ignore CG003 static args, panic unreachable in practice", "CG003", true},
		{"//cmdguard-lint:ignore CG002", "", false},
		{"//cmdguard-lint:ignore", "", false},
		{"//nolint:all", "", false},
		{"plain line", "", false},
	}

	for _, test := range tests {
		rule, ok := parseSuppressionLine(test.line)

		if ok != test.valid || rule != test.rule {
			t.Errorf("parseSuppressionLine(%q) = (%q, %v), want (%q, %v)", test.line, rule, ok, test.rule, test.valid)
		}
	}
}

func TestApplySuppressionsKeepsUnpositionedFindings(t *testing.T) {
	t.Parallel()

	f := finding.Finding{Rule: "CG001"}
	kept := applySuppressions(&project{}, []finding.Finding{f})

	if len(kept) != 1 {
		t.Fatalf("expected unpositioned finding kept, got %d", len(kept))
	}
}
