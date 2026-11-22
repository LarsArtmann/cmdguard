package v5

import (
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

// TestNoDeadSentinels fails when a package-level sentinel error is declared
// but never referenced outside its own declaration and the canonical facts
// table. The v5 scaffold mechanically re-added five editor/nil sentinels
// that v3 had already cut as dead weight; they sat in the public API and
// the generated ERROR_REFERENCE for weeks while no code path could ever
// produce them (found via `erraudit tree ./... --max` "unreferenced"
// badges, 2026-11-22). This tripwire keeps the declared sentinel set equal
// to the reachable one: a sentinel must be wired (returned, wrapped, or
// classified by non-facts code) or deleted. The detection logic itself is
// regression-tested against temp-dir fixtures in pkg/testutil
// (TestFindDeadSentinels*), and the same helper guards the
// flightrecorder and lint modules.
func TestNoDeadSentinels(t *testing.T) {
	t.Parallel()

	scan, err := testutil.FindDeadSentinels(".", "errors_facts.go")
	if err != nil {
		t.Fatalf("scanning package for dead sentinels: %v", err)
	}
	if scan.Declarations == 0 {
		t.Fatal("no sentinel declarations found; the tripwire's declaration scan is broken")
	}
	for _, f := range scan.Findings {
		t.Errorf(
			"sentinel %s (%s) is declared but never referenced outside its declaration and the facts table; remove it or wire it",
			f.Name,
			f.File,
		)
	}
}
