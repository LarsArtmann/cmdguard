package lint

import (
	"os"
	"strings"
	"testing"
)

// TestReadmeRuleTableInSync pins the README rule table to allRuleDefs(): the
// section between the BEGIN/END markers must equal RuleTableMarkdown()
// byte-for-byte. A failure means a rule was added, removed, or reworded
// without refreshing the README — regenerate with:
//
//	go run ./cmd/cmdguard-lint rules --markdown
//
// and replace the marked section.
func TestReadmeRuleTableInSync(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("reading README.md: %v", err)
	}

	content := string(raw)

	begin := strings.Index(content, readmeTableBegin)
	end := strings.Index(content, readmeTableEnd)

	if begin < 0 {
		t.Fatalf("README.md is missing the BEGIN marker %q", readmeTableBegin)
	}

	if end < begin {
		t.Fatalf("README.md is missing the END marker %q after the BEGIN marker", readmeTableEnd)
	}

	section := content[begin+len(readmeTableBegin) : end]
	want := strings.TrimSpace(RuleTableMarkdown())
	got := strings.TrimSpace(section)

	// Compared trimmed: dprint owns the blank lines around the section and
	// rewrites README.md whenever they differ, so pinning them here made the
	// two guards fight (generator regenerated raw, dprint repadded, guard red
	// again). The generator owns the table content byte-for-byte — padded
	// columns included (RuleTableMarkdown emits dprint's exact table format) —
	// and this test still fails on any rule add/remove/reword.
	if got != want {
		t.Errorf(`README.md rule table is stale (L7 drift guard).
Regenerate with: go run ./cmd/cmdguard-lint rules --markdown
and replace the section between the markers.

want:
%s
got:
%s`, want, got)
	}
}
