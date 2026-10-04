package lint

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/larsartmann/go-finding"
)

// detectOn runs the full Detect pass over a fixture tree.
func detectOn(t *testing.T, files map[string]string) []finding.Finding {
	t.Helper()

	return detectOnDir(t, writeTree(t, files))
}

func detectOnDir(t *testing.T, dir string) []finding.Finding {
	t.Helper()

	findings, err := Detect(context.Background(), dir)
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}

	return findings
}

// assertFinding verifies that a finding with the given rule exists at the
// expected file and line.
func assertFinding(t *testing.T, findings []finding.Finding, rule, file string, line int) {
	t.Helper()

	for _, f := range findings {
		if string(f.Rule) == rule && string(f.Position.File) == file && f.Position.Line == line {
			return
		}
	}

	t.Errorf("expected %s finding at %s:%d, got %v", rule, file, line, summarize(findings))
}

func summarize(findings []finding.Finding) string {
	parts := make([]string, 0, len(findings))

	for _, f := range findings {
		parts = append(parts, string(f.Rule)+"@"+string(f.Position.File)+":"+strconv.Itoa(f.Position.Line))
	}

	return "[" + strings.Join(parts, ", ") + "]"
}

func TestRuleExecuteBypassPositive(t *testing.T) {
	t.Parallel()

	// The exact timesheets shape: CLI dismantled via RootCommand, handed to
	// raw fang.Execute, in a project that imports cmdguard (other file).
	findings := detectOn(t, map[string]string{
		"guard.go": "package cmd\n\nimport cmdguard \"github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4\"\n\nvar _ = cmdguard.NoFlags{}\n",
		"main.go": `package main

import (
	"context"
	"io"

	"charm.land/fang/v2"
)

func main() {
	ctx := context.Background()
	cli := newCLI()
	err := fang.Execute(ctx, cli.RootCommand(),
		fang.WithErrorHandler(func(_ io.Writer, _ error) {}))
	_ = err
}
`,
	})

	assertFinding(t, findings, RuleExecuteBypass, "main.go", 13)
}

func TestRuleExecuteBypassPositiveViaVariable(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"a.go": "package a\n\nimport \"github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4\"\n\nvar _ = v4.NoFlags{}\n",
		"main.go": `package main

import (
	"charm.land/fang/v2"
)

func main() {
	cli := newCLI()
	root := cli.RootCommand()
	_ = fang.Execute(nil, root)
}
`,
	})

	assertFinding(t, findings, RuleExecuteBypass, "main.go", 10)
}

func TestRuleExecuteBypassNegative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		files map[string]string
	}{
		{
			name: "cli.Execute is the blessed path",
			files: map[string]string{
				"main.go": `package main

import (
	"context"

	"charm.land/fang/v2"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func main() {
	ctx := context.Background()
	cli, _ := v4.NewCLI("app", "app", struct{}{})
	_ = fang.Execute(ctx, nil) // fang without RootCommand: other runner use
	_ = cli.Execute(ctx)
}
`,
			},
		},
		{
			name: "no cmdguard anywhere silences the rule",
			files: map[string]string{
				"main.go": `package main

import "charm.land/fang/v2"

type fakeCLI struct{}

func (fakeCLI) RootCommand() *fakeCmd { return nil }

type fakeCmd struct{}

func main() {
	cli := fakeCLI{}
	_ = fang.Execute(nil, cli.RootCommand())
}
`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			for _, f := range detectOn(t, test.files) {
				if string(f.Rule) == RuleExecuteBypass {
					t.Errorf("unexpected %s finding at %s:%d", RuleExecuteBypass, f.Position.File, f.Position.Line)
				}
			}
		})
	}
}

func TestRuleStaleMajorImportPositive(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"cmd.go": "package cmd\n\nimport cmdguard \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = cmdguard.NoFlags{}\n",
	})

	assertFinding(t, findings, RuleStaleMajorImport, "cmd.go", 3)
}

func TestRuleStaleMajorImportNegative(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"cmd.go": "package cmd\n\nimport v4 \"github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4\"\n\nvar _ = v4.NoFlags{}\n",
		"sub.go": "package cmd\n\nimport \"github.com/larsartmann/cmdguard/spinner\"\n\nvar _ = spinner.Spin\n",
	})

	for _, f := range findings {
		if string(f.Rule) == RuleStaleMajorImport {
			t.Errorf("unexpected %s finding at %s:%d", RuleStaleMajorImport, f.Position.File, f.Position.Line)
		}
	}
}

func TestRulePanicOnConstructorPositive(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"helpers.go": `package cmd

import (
	"fmt"

	cmdguard "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func mustNewCommand(use string) cmdguard.Command[struct{}, struct{}] {
	cmd, err := cmdguard.NewCommand(use, struct{}{}, nil)
	if err != nil {
		panic(fmt.Sprintf("mustNewCommand(%q): %v", use, err))
	}

	return cmd
}
`,
	})

	assertFinding(t, findings, RulePanicOnConstructor, "helpers.go", 12)
}

func TestRulePanicOnConstructorPositiveAddCommand(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"reg.go": `package cmd

import (
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func register(cli *v4.CLI[struct{}], cmd v4.Command[struct{}, struct{}]) {
	if err := v4.AddCommand(cli, cmd); err != nil {
		panic(err)
	}
}
`,
	})

	assertFinding(t, findings, RulePanicOnConstructor, "reg.go", 9)
}

func TestRulePanicOnConstructorNegative(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"ok.go": `package cmd

import (
	"fmt"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

// Errors are returned, not panicked on.
func newCommand(use string) (v4.Command[struct{}, struct{}], error) {
	cmd, err := v4.NewCommand(use, struct{}{}, nil)
	if err != nil {
		return cmd, fmt.Errorf("building %q: %w", use, err)
	}

	return cmd, nil
}

// Panic on an unrelated error stays unflagged.
func other(path string) {
	if _, err := osStat(path); err != nil {
		panic(err)
	}
}
`,
	})

	for _, f := range findings {
		if string(f.Rule) == RulePanicOnConstructor {
			t.Errorf("unexpected %s finding at %s:%d", RulePanicOnConstructor, f.Position.File, f.Position.Line)
		}
	}
}

func TestRuleSetVersionPositive(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"guard.go": `package cmd

import (
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func setup() {
	cli, err := v4.NewCLI("app", "app", struct{}{})
	if err != nil {
		return
	}

	cli.SetVersion("1.2.3")
}
`,
	})

	assertFinding(t, findings, RuleSetVersionRuntime, "guard.go", 13)
}

func TestRuleSetVersionNegativeUnrelatedReceiver(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"other.go": `package other

type thing struct{}

func (thing) SetVersion(v string) {}

func setup() {
	t := thing{}
	t.SetVersion("1.0.0")
}
`,
	})

	for _, f := range findings {
		if string(f.Rule) == RuleSetVersionRuntime {
			t.Errorf("unexpected %s finding at %s:%d", RuleSetVersionRuntime, f.Position.File, f.Position.Line)
		}
	}
}

func TestRuleDuplicateVersionOptsPositive(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"main.go": `package main

import (
	"charm.land/fang/v2"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func build() {
	_, _ = v4.NewCLI("app", "app", struct{}{},
		v4.WithCLIVersion("1.2.3"),
		v4.WithFangOptions(fang.WithVersion("1.2.3")),
	)
}
`,
	})

	assertFinding(t, findings, RuleDuplicateVersionOpts, "main.go", 9)
}

func TestRuleDuplicateVersionOptsNegative(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"main.go": `package main

import (
	"charm.land/fang/v2"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func build() {
	_, _ = v4.NewCLI("app", "app", struct{}{},
		v4.WithCLIVersion("1.2.3"),
		v4.WithFangOptions(fang.WithColors()),
	)
}
`,
	})

	for _, f := range findings {
		if string(f.Rule) == RuleDuplicateVersionOpts {
			t.Errorf("unexpected %s finding at %s:%d", RuleDuplicateVersionOpts, f.Position.File, f.Position.Line)
		}
	}
}

func TestRuleExecuteErrorReprintPositive(t *testing.T) {
	t.Parallel()

	findings := detectOn(t, map[string]string{
		"main.go": `package main

import (
	"context"
	"fmt"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func main() {
	ctx := context.Background()
	cli, _ := v4.NewCLI("app", "app", struct{}{})
	if err := cli.Execute(ctx); err != nil {
		fmt.Println(err)
	}
}
`,
	})

	assertFinding(t, findings, RuleExecuteErrorReprint, "main.go", 14)
}

func TestRuleExecuteErrorReprintNegative(t *testing.T) {
	t.Parallel()

	// fmt.Errorf wrapping is the supported exit-code mapping path; fang
	// errors are not cmdguard Execute errors.
	findings := detectOn(t, map[string]string{
		"main.go": `package main

import (
	"context"
	"fmt"

	"charm.land/fang/v2"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

func main() {
	ctx := context.Background()
	cli, _ := v4.NewCLI("app", "app", struct{}{})
	if err := cli.Execute(ctx); err != nil {
		wrapped := fmt.Errorf("run failed: %w", err)
		_ = wrapped
	}

	if fangErr := fangHelp(); fangErr != nil {
		fmt.Println(fangErr)
	}
}

func fangHelp() error { return nil }

var _ = fang.Version
`,
	})

	for _, f := range findings {
		if string(f.Rule) == RuleExecuteErrorReprint {
			t.Errorf("unexpected %s finding at %s:%d", RuleExecuteErrorReprint, f.Position.File, f.Position.Line)
		}
	}
}

func TestDetectRunsAllRulesWithUniqueIDs(t *testing.T) {
	t.Parallel()

	all := AllRules()

	if len(all) != 6 {
		t.Fatalf("expected 6 rules, got %d", len(all))
	}

	ids := make([]string, 0, len(all))

	for _, rule := range all {
		ids = append(ids, rule.Meta.ID)

		if rule.Meta.ToolName != ToolName {
			t.Errorf("rule %s: tool name %q, want %q", rule.Meta.ID, rule.Meta.ToolName, ToolName)
		}

		if rule.Meta.Sev == "" {
			t.Errorf("rule %s: severity must be set", rule.Meta.ID)
		}
	}

	slices.Sort(ids)

	for i := 1; i < len(ids); i++ {
		if ids[i] == ids[i-1] {
			t.Errorf("duplicate rule ID %s", ids[i])
		}
	}
}

func TestRegistryRunsRulesAndMatchesDetect(t *testing.T) {
	t.Parallel()

	dir := writeTree(t, map[string]string{
		"stale.go": "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{}\n",
	})

	viaRegistry, err := NewRegistry().Run(context.Background(), dir)
	if err != nil {
		t.Fatalf("registry run: %v", err)
	}

	viaDetect := detectOnDir(t, dir)

	if viaRegistry.Len() != len(viaDetect) {
		t.Errorf("registry findings %d != Detect findings %d", viaRegistry.Len(), len(viaDetect))
	}
}
