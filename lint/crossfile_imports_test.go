package lint

import (
	"testing"
)

// H9 (2026-11-08): fixture arms for the import-resolution edges —
// dot imports and blank imports. Before these fixtures both arms had zero
// coverage.
//
// KNOWN GAP (recorded, not fixed here): dot-imported CONSTRUCTORS are not
// traced. `isSelectorCall` requires a selector (`pkg.NewCLI`), and a dot
// import produces a bare `NewCLI(...)` call — so a dot-import-only file gets
// no CLI-variable tracing and no cross-file rules. The `selectorFrom` dot arm
// (bare selector treated as from the dot-imported package) covers method
// calls on TRACED vars only. Fixing constructor tracing for dot imports is
// Tracked as the H9b follow-up in the cmdguard TODO_LIST. Dot imports are
// also non-idiomatic for cmdguard consumers (the aliased `v5` import is the
// documented pattern), which is why this is a recorded gap rather than a
// release blocker.

func TestDotImport_BareConstructorsAreTraced(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"main.go": `package app

import (
	. "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	cli, err := NewCLI("app", "App", cfg{})
	if err != nil {
		panic(err)
	}
	_ = cli.Execute(nil)
}
`,
	})

	// H9b closed (2026-11-11): a dot import puts cmdguard's symbols in file
	// scope, so bare constructor calls are traced like qualified ones and
	// the panic on the constructor error is flagged (CG003). Before the
	// fix this pinned zero findings as a known gap.
	if got := ruleIDsOf(t, findings); len(got) != 1 || got[0] != "CG003" {
		t.Errorf("rules=%v, want exactly [CG003] (dot-imported bare constructors must be traced)", got)
	}
}

func TestBlankImport_IsSkippedWithoutCrash(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"main.go": `package app

import (
	_ "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	cli, err := newCLI()
	if err != nil {
		panic(err)
	}
	_ = cli
}

func newCLI() (*int, error) { return nil, nil }
`,
	})

	// No cmdguard construction happens through a blank import; the pass must
	// neither crash nor fabricate findings.
	if got := ruleIDsOf(t, findings); len(got) != 0 {
		t.Errorf("rules=%v, want none (blank import registers nothing)", got)
	}
}
