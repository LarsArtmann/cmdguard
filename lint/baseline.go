package lint

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/larsartmann/go-finding"
)

// Sentinel errors for baseline handling. Match with errors.Is; wrapped
// details name the offending path.
var (
	// ErrBaselineVersionMismatch reports a baseline file whose schema version
	// this cmdguard-lint cannot interpret (regenerating fixes it).
	ErrBaselineVersionMismatch = errors.New("baseline schema version mismatch")
	// ErrBaselineNotFound reports an explicitly requested baseline file that
	// does not exist (the default location being absent is normal, not an
	// error: LoadBaseline signals that via its found return).
	ErrBaselineNotFound = errors.New("baseline file not found")
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
// baseline was written. Line numbers are match hints, not exact contracts:
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

// LoadBaseline reads and validates a baseline file. found is false (with nil
// baseline and nil error) when the file does not exist: "no baseline" is the
// normal pre-adoption state. Use the explicit path when the caller asked for
// one and wants absence reported as ErrBaselineNotFound.
func LoadBaseline(path string) (baseline *Baseline, found bool, err error) {
	raw, readErr := os.ReadFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("reading baseline %s: %w", path, readErr)
	}

	if unmarshalErr := json.Unmarshal(raw, &baseline); unmarshalErr != nil {
		return nil, false, fmt.Errorf("parsing baseline %s: %w", path, unmarshalErr)
	}

	if baseline.Version != baselineVersion {
		return nil, false, fmt.Errorf(
			"%w: %s has version %d, this cmdguard-lint expects %d (regenerate the baseline)",
			ErrBaselineVersionMismatch, path, baseline.Version, baselineVersion,
		)
	}

	return baseline, true, nil
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

	if err := os.WriteFile(path, append(encoded, '\n'), 0o600); err != nil {
		return fmt.Errorf("writing baseline %s: %w", path, err)
	}

	return nil
}

// ApplyBaseline splits findings into fresh (not covered by the baseline:
// these gate the run) and returns the stale baseline entries (cover nothing
// anymore: fixed findings, so the baseline can be tightened). A nil baseline
// marks every finding fresh.
//
// Matching is per rule+file with one baseline entry consuming at most one
// finding: an exact line match wins first, then any remaining entry for the
// same rule+file absorbs line drift. The per-file count is what ultimately
// gates: more findings than entries for a rule+file means the surplus is new.
func ApplyBaseline(
	findings []finding.Finding,
	baseline *Baseline,
) ([]finding.Finding, []BaselineEntry) {
	if baseline == nil {
		return findings, nil
	}

	consumed := map[int]bool{}
	kept := map[int]bool{}

	for i := range findings {
		f := &findings[i]

		entry := matchBaselineEntry(
			baseline, consumed,
			string(f.Rule), string(f.Position.File), f.Position.Line,
		)

		if entry >= 0 {
			consumed[entry] = true
			kept[i] = true
		}
	}

	fresh := make([]finding.Finding, 0, len(findings)-len(kept))

	for i := range findings {
		if !kept[i] {
			fresh = append(fresh, findings[i])
		}
	}

	stale := make([]BaselineEntry, 0, len(baseline.Entries)-len(consumed))

	for i, entry := range baseline.Entries {
		if !consumed[i] {
			stale = append(stale, entry)
		}
	}

	return fresh, stale
}

// matchBaselineEntry finds the index of the best unconsumed entry for the
// rule+file: an exact line match if present, otherwise the first entry for
// the same rule+file (drift absorption). File paths compare exactly: the
// baseline is written from this tool's own findings, so paths always agree.
// Returns -1 when nothing matches.
func matchBaselineEntry(baseline *Baseline, consumed map[int]bool, rule, file string, line int) int {
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
