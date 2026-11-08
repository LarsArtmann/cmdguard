package lint

import (
	"testing"
)

// Generic instantiations wrap the called selector in IndexExpr (one type
// argument) or IndexListExpr (two or more). cmdguard's constructors are
// generic — `v5.NewCLI[Config](...)`, `v5.NewCommand[C, F](...)` — so every
// selector-matching rule must see through the type arguments. These tests
// pin the unwrap in selectorExprOf (regression for the v0.2.1 fix where
// explicit instantiation syntax silenced the rules).
func TestGenericInstantiationDetection(t *testing.T) {
	t.Parallel()

	t.Run("CG007 flags NewCommand[T, F] without WithShort", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"context"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type scanFlags struct{}

func main() {
	cmd, err := v5.NewCommand[struct{}, scanFlags]("scan", scanFlags{},
		func(_ context.Context, _ *struct{}, _ scanFlags) error { return nil })
	_, _ = cmd, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingDescription, "main.go") {
			t.Fatalf("expected %s on a generic NewCommand call, got %v", RuleMissingDescription, summarize(findings))
		}
	})

	t.Run("CG008 flags NewCLI[T] without WithCLIVersion", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type demoConfig struct{}

func main() {
	cli, err := v5.NewCLI[demoConfig]("demo", "demo CLI", demoConfig{})
	_, _ = cli, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("expected %s on a generic NewCLI call, got %v", RuleMissingVersion, summarize(findings))
		}
	})

	t.Run("cross-file: generic package-level constructor error feeds CG003", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"cli.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type demoConfig struct{}

var appCLI, appErr = v5.NewCLI[demoConfig]("demo", "demo CLI", demoConfig{})
`,
			"main.go": `package main

import "fmt"

func main() {
	if appErr != nil {
		panic(appErr)
	}

	fmt.Print(appCLI)
}
`,
		})

		if !hasRuleInFile(findings, RulePanicOnConstructor, "main.go") {
			t.Fatalf(
				"expected %s via the generic package-level constructor union, got %v",
				RulePanicOnConstructor,
				summarize(findings),
			)
		}
	})

	t.Run("non-generic calls still detected", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type demoConfig struct{}

func main() {
	cli, err := v5.NewCLI("demo", "demo CLI", demoConfig{})
	_, _ = cli, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("expected %s on a plain NewCLI call, got %v", RuleMissingVersion, summarize(findings))
		}
	})
}
