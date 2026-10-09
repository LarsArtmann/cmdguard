package lint

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/larsartmann/go-finding"
)

// baselineVersion is the schema version of the baseline file. Bump on
// incompatible changes so LoadBaseline can reject future formats loudly
// instead of misinterpreting them.
const baselineVersion = 1

// DefaultBaselinePath is the file a lint run picks up automatically (relative
// to the linted directory) when no explicit --baseline path is given. Keep in
// sync with the CLI flag default documented in README.
const DefaultBaselinePath = ".cmdguard-lint-baseline.json"

// BaselineEntry records one accepted finding: its rule and location when the
// baseline was written. Line numbers are match hints, not exact contracts —
// ApplyBaseline tolerates drift within the same rule+file (edits above a
// finding shift its line without changing its identity).
type BaselineEntry struct {
	Rule string `json:"rule"`
	File string `json:"file"`
	Line int    `json:"line"`
}

// Baseline is the on-disk ratchet state: the findings a codebase already had
// when cmdguard-lint was adopted. The lint run fails only on findings NOT
// covered by the baseline; fixing a baselined finding and rewriting the
// baseline ratchets the count down. A baseline never grows without a rewrite.
type Baseline struct {
	Version int             `json:"version"`
	Entries []BaselineEntry `json:"entries"`
}

// LoadBaseline reads and validates a baseline file. A missing file returns
// (nil, nil): "no baseline" is the normal pre-adoption state, not an error.
func LoadBaseline(path string) (*Baseline, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("reading baseline %s: %w", path, err)
	}

	var baseline Baseline
	if err := json.Unmarshal(raw, &baseline); err != nil {
		return nil, fmt.Errorf("parsing baseline %s: %w", path, err)
	}

	if baseline.Version != baselineVersion {
		return nil, fmt.Errorf(
			"baseline %s has version %d, this cmdguard-lint expects %d — regenerate the baseline",
			path, baseline.Version, baselineVersion,
		)
	}

	return &baseline, nil
}

// WriteBaseline serializes findings into a baseline file at path. Used at
// adoption time (--write-baseline) and whenever the ratchet should tighten.
func WriteBaseline(findings []finding.Finding, path string) error {
	baseline := Baseline{
		Version: baselineVersion,
		Entries: make([]BaselineEntry, 0, len(findings)),
	}

	for _, f := range findings {
		baseline.Entries = append(baseline.Entries, BaselineEntry{
			Rule: string(f.Rule),
			File: string(f.Position.File),
			Line: f.Position.Line,
		})
	}

	encoded, err := json.MarshalIndent(baseline, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling baseline: %w", err)
	}

	if err := os.WriteFile(path, append(encoded, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing baseline %s: %w", path, err)
	}

	return nil
}

// ApplyBaseline splits findings into fresh (not covered by the baseline —
// these gate the run) and returns the stale baseline entries (cover nothing
// anymore — fixed findings, so the baseline can be tightened). A nil baseline
// marks every finding fresh.
//
// Matching is per rule+file with one baseline entry consuming at most one
// finding: an exact line match wins first, then any remaining entry for the
// same rule+file absorbs line drift. The per-file count is what ultimately
// gates: more findings than entries for a rule+file means the surplus is new.
func ApplyBaseline(
	findings []finding.Finding,
	baseline *Baseline,
) (fresh []finding.Finding, stale []BaselineEntry) {
	if baseline == nil {
		return findings, nil
	}

	consumed := make([]bool, len(baseline.Entries))
	kept := make([]bool, len(findings))

	for i := range findings {
		f := &findings[i]
		entry := matchBaselineEntry(baseline, consumed, string(f.Rule), string(f.Position.File), f.Position.Line)
		if entry >= 0 {
			consumed[entry] = true
			kept[i] = true
		}
	}

	for i, wasKept := range kept {
		if !wasKept {
			fresh = append(fresh, findings[i])
		}
	}

	for i, used := range consumed {
		if !used {
			stale = append(stale, baseline.Entries[i])
		}
	}

	return fresh, stale
}

// matchBaselineEntry finds the index of the best unconsumed entry for the
// rule+file: an exact line match if present, otherwise the first entry for
// the same rule+file (drift absorption). File paths compare exactly — the
// baseline is written from this tool's own findings, so paths always agree.
// Returns -1 when nothing matches.
func matchBaselineEntry(baseline *Baseline, consumed []bool, rule, file string, line int) int {
	fallback := -1

	for i, entry := range baseline.Entries {
		if consumed[i] || entry.Rule != rule || entry.File != file {
			continue
		}

		if entry.Line == line {
			return i
		}

		if fallback < 0 {
			fallback = i
		}
	}

	return fallback
}
