package provider

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
)

func TestProviderDeclaresEnableDisableOptions(t *testing.T) {
	t.Parallel()

	spec, ok := lookupSpec(t)
	if !ok {
		t.Fatal("cmdguard-lint spec not registered")
	}

	if len(spec.Options) != 2 {
		t.Fatalf("Options length = %d, want 2 (enable, disable)", len(spec.Options))
	}

	for _, opt := range spec.Options {
		if opt.Kind != toolsdk.OptionKindString {
			t.Errorf("option %q Kind = %q, want string", opt.Name, opt.Kind)
		}

		if opt.Description == "" {
			t.Errorf("option %q has no Description", opt.Name)
		}
	}

	if err := spec.ValidateOptions(toolsdk.OptionValues{"enable": "CG001", "disable": "CG004"}); err != nil {
		t.Errorf("ValidateOptions(valid values) = %v, want nil", err)
	}

	if err := spec.ValidateOptions(toolsdk.OptionValues{"threshold": 3}); err == nil {
		t.Error("ValidateOptions(unknown option) = nil, want ErrUnknownOption")
	}

	if err := spec.ValidateOptions(toolsdk.OptionValues{"enable": 3}); err == nil {
		t.Error("ValidateOptions(kind mismatch) = nil, want ErrOptionKindMismatch")
	}
}

func TestProviderDetectHonorsOptions(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	// Two findings by construction: a stale v3 import (CG002) and a runtime
	// SetVersion on a NewCLI-assigned variable (CG004).
	fixture := `package x

import (
	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
)

var _ = v3.NoFlags{}

func setup() {
	cli, _ := v4.NewCLI("x", "X", cfg{})
	cli.SetVersion("2.0.0")
}

type cfg struct{}
`

	if err := os.WriteFile(filepath.Join(dir, "fixture.go"), []byte(fixture), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	spec, ok := lookupSpec(t)
	if !ok {
		t.Fatal("cmdguard-lint spec not registered")
	}

	base := finding.WithWorkingDir(context.Background(), dir)

	unfiltered, err := spec.Detect.Detect(base)
	if err != nil {
		t.Fatalf("detect: %v", err)
	}

	if !hasRule(unfiltered, "CG002") || !hasRule(unfiltered, "CG004") {
		t.Fatalf("fixture must produce CG002 and CG004, got rules %v", rulesOf(unfiltered))
	}

	disabled := toolsdk.WithOptions(base, toolsdk.OptionValues{"disable": "CG004"})
	filtered, err := spec.Detect.Detect(disabled)
	if err != nil {
		t.Fatalf("detect(disabled): %v", err)
	}

	if hasRule(filtered, "CG004") {
		t.Error("disable=CG004 still produced a CG004 finding")
	}

	if !hasRule(filtered, "CG002") {
		t.Error("disable=CG004 must not affect CG002")
	}

	enabledOnly := toolsdk.WithOptions(base, toolsdk.OptionValues{"enable": " CG002 "})
	onlyTwo, err := spec.Detect.Detect(enabledOnly)
	if err != nil {
		t.Fatalf("detect(enabled): %v", err)
	}

	if hasRule(onlyTwo, "CG004") {
		t.Error("enable=CG002 still produced a CG004 finding")
	}

	if !hasRule(onlyTwo, "CG002") {
		t.Error("enable=CG002 lost the CG002 finding")
	}
}

func hasRule(findings []finding.Finding, rule string) bool {
	for _, f := range findings {
		if string(f.Rule) == rule {
			return true
		}
	}

	return false
}

func rulesOf(findings []finding.Finding) []string {
	rules := make([]string, 0, len(findings))
	for _, f := range findings {
		rules = append(rules, string(f.Rule))
	}

	return rules
}
