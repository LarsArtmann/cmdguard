package lint

import (
	"slices"
	"testing"

	"github.com/larsartmann/go-finding"
)

func TestFilterByIDs(t *testing.T) {
	t.Parallel()

	findings := []finding.Finding{
		findingAt("CG001", "main.go", 1),
		findingAt("CG002", "main.go", 2),
		findingAt("CG004", "run.go", 3),
	}

	t.Run("no filters keep everything", func(t *testing.T) {
		t.Parallel()

		if kept := FilterByIDs(findings, "", ""); len(kept) != 3 {
			t.Errorf("kept=%d, want 3", len(kept))
		}
	})

	t.Run("disable drops matching rules", func(t *testing.T) {
		t.Parallel()

		kept := FilterByIDs(findings, "", "CG002, CG004")

		got := rulesOfFindings(kept)
		want := []string{"CG001"}

		if !slices.Equal(got, want) {
			t.Errorf("rules=%v, want %v (spaces between IDs tolerated)", got, want)
		}
	})

	t.Run("enable keeps only matching rules", func(t *testing.T) {
		t.Parallel()

		kept := FilterByIDs(findings, "CG002", "")

		got := rulesOfFindings(kept)
		want := []string{"CG002"}

		if !slices.Equal(got, want) {
			t.Errorf("rules=%v, want %v", got, want)
		}
	})

	t.Run("enable wins over disable for the same rule", func(t *testing.T) {
		t.Parallel()

		// disable removes first; an enabled rule that is also disabled is
		// gone — disable has veto semantics (matches the CLI flag order).
		kept := FilterByIDs(findings, "CG002", "CG002")

		if len(kept) != 0 {
			t.Errorf("kept=%d, want 0 (disable vetoes enable)", len(kept))
		}
	})

	t.Run("empty segments are ignored", func(t *testing.T) {
		t.Parallel()

		if kept := FilterByIDs(findings, "", " , CG001 ,"); len(kept) != 2 {
			t.Errorf("kept=%d, want 2", len(kept))
		}
	})
}

func rulesOfFindings(findings []finding.Finding) []string {
	rules := make([]string, 0, len(findings))
	for _, f := range findings {
		rules = append(rules, string(f.Rule))
	}

	return rules
}
