package lint

import (
	"testing"

	"github.com/larsartmann/go-finding"
)

// hasRuleInFile reports whether findings contain a finding of rule in file.
func hasRuleInFile(findings []finding.Finding, rule, file string) bool {
	for _, f := range findings {
		if string(f.Rule) == rule && string(f.Position.File) == file {
			return true
		}
	}

	return false
}

func TestRuleMissingShortDescription(t *testing.T) {
	t.Parallel()

	t.Run("flags NewCommand without WithShort", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"context"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func main() {
	cmd, err := v5.NewCommand("scan", v5.NoFlags{},
		func(_ context.Context, _ *struct{}, _ v5.NoFlags) error { return nil })
	_, _ = cmd, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingDescription, "main.go") {
			t.Fatalf("expected %s, got %v", RuleMissingDescription, summarize(findings))
		}
	})

	t.Run("accepts WithShort", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"context"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func main() {
	cmd, err := v5.NewCommand("scan", v5.NoFlags{},
		func(_ context.Context, _ *struct{}, _ v5.NoFlags) error { return nil },
		v5.WithShort("Scan things"))
	_, _ = cmd, err
}
`,
		})

		if hasRuleInFile(findings, RuleMissingDescription, "main.go") {
			t.Fatalf("unexpected %s: %v", RuleMissingDescription, summarize(findings))
		}
	})

	t.Run("rejects an empty WithShort", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"context"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func main() {
	cmd, err := v5.NewCommand("scan", v5.NoFlags{},
		func(_ context.Context, _ *struct{}, _ v5.NoFlags) error { return nil },
		v5.WithShort("  "))
	_, _ = cmd, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingDescription, "main.go") {
			t.Fatalf("expected %s for blank short, got %v", RuleMissingDescription, summarize(findings))
		}
	})

	t.Run("skips spread options", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"context"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func main() {
	opts := []v5.CommandOption{}
	cmd, err := v5.NewCommand("scan", v5.NoFlags{},
		func(_ context.Context, _ *struct{}, _ v5.NoFlags) error { return nil },
		opts...)
	_, _ = cmd, err
}
`,
		})

		if hasRuleInFile(findings, RuleMissingDescription, "main.go") {
			t.Fatalf("unexpected %s for spread options: %v", RuleMissingDescription, summarize(findings))
		}
	})

	t.Run("ignores test files", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main_test.go": `package main

import (
	"context"
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestApp(t *testing.T) {
	_, _ = v5.NewCommand("scan", v5.NoFlags{},
		func(_ context.Context, _ *struct{}, _ v5.NoFlags) error { return nil })
}
`,
		})

		if hasRuleInFile(findings, RuleMissingDescription, "main_test.go") {
			t.Fatalf("unexpected %s in test file: %v", RuleMissingDescription, summarize(findings))
		}
	})
}

func TestRuleMissingCLIVersion(t *testing.T) {
	t.Parallel()

	t.Run("flags NewCLI without a version", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	cli, err := v5.NewCLI("app", "App", cfg{})
	_, _ = cli, err
}
`,
		})

		if !hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("expected %s, got %v", RuleMissingVersion, summarize(findings))
		}
	})

	t.Run("accepts WithCLIVersion", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	cli, err := v5.NewCLI("app", "App", cfg{}, v5.WithCLIVersion("1.0.0"))
	_, _ = cli, err
}
`,
		})

		if hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("unexpected %s: %v", RuleMissingVersion, summarize(findings))
		}
	})

	t.Run("accepts a fang version option", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	"charm.land/fang/v2"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	cli, err := v5.NewCLI("app", "App", cfg{},
		v5.WithFangOptions(fang.WithVersion("1.0.0")))
	_, _ = cli, err
}
`,
		})

		if hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("unexpected %s: %v", RuleMissingVersion, summarize(findings))
		}
	})

	t.Run("skips spread options", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main.go": `package main

import (
	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func main() {
	opts := []v5.CLIOption{}
	cli, err := v5.NewCLI("app", "App", cfg{}, opts...)
	_, _ = cli, err
}
`,
		})

		if hasRuleInFile(findings, RuleMissingVersion, "main.go") {
			t.Fatalf("unexpected %s for spread options: %v", RuleMissingVersion, summarize(findings))
		}
	})

	t.Run("ignores test files", func(t *testing.T) {
		t.Parallel()

		findings := detectOn(t, map[string]string{
			"main_test.go": `package main

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

type cfg struct{}

func TestBuild(t *testing.T) {
	_, _ = v5.NewCLI("app", "App", cfg{})
}
`,
		})

		if hasRuleInFile(findings, RuleMissingVersion, "main_test.go") {
			t.Fatalf("unexpected %s in test file: %v", RuleMissingVersion, summarize(findings))
		}
	})
}
