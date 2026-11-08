package lint

import (
	"testing"
)

// H9 (2026-11-08): fixture arms for the import-resolution edges —
// dot imports (package symbols land in file scope; selectorFrom treats any
// bare selector as potentially from the dot-imported package) and blank
// imports (skipped as identifiers). Before these fixtures both arms had zero
// coverage.

func TestDotImport_StillDetectsExecuteBypass(t *testing.T) {
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

	if got := ruleIDsOf(t, findings); !containsRule(got, "CG001") {
		t.Errorf("rules=%v, want CG001 through a dot import", got)
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
