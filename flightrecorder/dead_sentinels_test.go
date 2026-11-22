package flightrecorder

import (
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

// TestNoDeadSentinels guards this module's sentinel errors against the
// scaffold-rot bug class the 2026-11-22 core audit caught (the v5 scaffold
// had mechanically re-added five declared-but-unwired sentinels): every
// package-level Err*/err* var initialized by a call — or typed error — must
// be referenced outside its declaration somewhere in the module, or be
// deleted. Detection logic is shared with the core tripwire via
// testutil.FindDeadSentinels and regression-tested there.
func TestNoDeadSentinels(t *testing.T) {
	t.Parallel()

	scan, err := testutil.FindDeadSentinels(".", "")
	if err != nil {
		t.Fatalf("scanning module for dead sentinels: %v", err)
	}
	if scan.Declarations == 0 {
		t.Fatal("no sentinel declarations found; the tripwire's declaration scan is broken")
	}
	for _, f := range scan.Findings {
		t.Errorf(
			"sentinel %s (%s) is declared but never referenced outside its declaration; remove it or wire it",
			f.Name,
			f.File,
		)
	}
}
