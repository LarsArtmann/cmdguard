package lint

import (
	"testing"
)

// crossFileFixtures returns the two-file package used by the cross-file
// tests: setup.go owns the package-level CLI construction, use.go owns the
// misuse that only the directory-group name tracing can connect.
func crossFileFixtures(useGo string) map[string]string {
	return map[string]string{
		"setup.go": `package app

import (
	"context"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

var cli, setupErr = v4.NewCLI("app", "App", cfg{})

type cfg struct{}

func background() context.Context { return context.Background() }
`,
		"use.go": useGo,
	}
}

func TestCrossFileTracing(t *testing.T) {
	t.Parallel()

	t.Run("CG004 constructor in one file, SetVersion in another", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, crossFileFixtures(`package app

func patchVersion() {
	cli.SetVersion("2.0.0")
}
`))

		if got := ruleIDsOf(t, findings); !containsRule(got, "CG004") {
			t.Errorf("rules=%v, want CG004 (package-level cli var + cross-file SetVersion)", got)
		}
	})

	t.Run("CG006 package-level Execute error printed in another file", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"setup.go": `package app

import (
	"context"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

var cli, setupErr = v4.NewCLI("app", "App", cfg{})

type cfg struct{}

func background() context.Context { return context.Background() }

var execErr = cli.Execute(background())
`,
			"use.go": `package app

import "fmt"

func report() {
	fmt.Println(execErr)
}
`,
		})

		if got := ruleIDsOf(t, findings); !containsRule(got, "CG006") {
			t.Errorf("rules=%v, want CG006 (package-level execErr + cross-file Println)", got)
		}
	})

	t.Run("CG003 package-level constructor error panicked in another file", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, crossFileFixtures(`package app

func mustApp() {
	if setupErr != nil {
		panic(setupErr)
	}
}
`))

		if got := ruleIDsOf(t, findings); !containsRule(got, "CG003") {
			t.Errorf("rules=%v, want CG003 (package-level setupErr + cross-file panic)", got)
		}
	})

	t.Run("function-local names do not leak across files", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"setup.go": `package app

import v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"

func build() {
	cli, err := v4.NewCLI("app", "App", cfg{})
	_, _ = cli, err
}

type cfg struct{}
`,
			"other.go": `package app

type ownCLI struct{}

func (ownCLI) SetVersion(string) {}

func unrelated(msg string) {
	var cli ownCLI
	cli.SetVersion(msg)
}
`,
		})

		if got := ruleIDsOf(t, findings); containsRule(got, "CG004") {
			t.Errorf("rules=%v, function-local cli from setup.go must not leak into other.go", got)
		}
	})

	t.Run("test and prod files do not merge", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"setup.go": `package app

import v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"

var cli, _ = v4.NewCLI("app", "App", cfg{})

type cfg struct{}
`,
			"setup_test.go": `package app

type fakeCLI struct{}

func (fakeCLI) SetVersion(string) {}

func testPatch() {
	var cli fakeCLI
	cli.SetVersion("test")
}
`,
		})

		if got := ruleIDsOf(t, findings); containsRule(got, "CG004") {
			t.Errorf("rules=%v, prod cliVars must not merge into the test package namespace", got)
		}
	})
}

func containsRule(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
