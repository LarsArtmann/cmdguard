package provider

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/go-finding"
	"github.com/larsartmann/go-finding/toolsdk"
)

func TestProviderRegisteredWithSpecContract(t *testing.T) {
	t.Parallel()

	spec, ok := lookupSpec(t)
	if !ok {
		t.Fatal("cmdguard-lint spec not registered in toolsdk.All()")
	}

	if spec.Description == "" {
		t.Error("spec has no description")
	}

	if spec.Detect == nil {
		t.Fatal("spec has no Detect capability")
	}

	if spec.Repair != nil {
		t.Error("linter must not register a Repairer (no autofixes)")
	}

	if !spec.ModuleFanOut {
		t.Error("spec must declare ModuleFanOut for multi-module workspaces")
	}

	if spec.Trigger.Language != "go" {
		t.Errorf("trigger language %q, want go", spec.Trigger.Language)
	}
}

func TestProviderDetectUsesWorkingDirFromContext(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	stale := "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{}\n"

	if err := os.WriteFile(filepath.Join(dir, "stale.go"), []byte(stale), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	spec, ok := lookupSpec(t)
	if !ok {
		t.Fatal("cmdguard-lint spec not registered")
	}

	ctx := finding.WithWorkingDir(context.Background(), dir)

	findings, err := spec.Detect.Detect(ctx)
	if err != nil {
		t.Fatalf("detect: %v", err)
	}

	found := false

	for _, f := range findings {
		if string(f.Rule) == "CG002" {
			found = true
		}
	}

	if !found {
		t.Errorf("expected CG002 finding from detector, got %d findings", len(findings))
	}
}

// lookupSpec returns the registered cmdguard-lint spec.
func lookupSpec(t *testing.T) (toolsdk.Spec, bool) {
	t.Helper()

	for _, spec := range toolsdk.All() {
		if spec.Name == "cmdguard-lint" {
			return spec, true
		}
	}

	return toolsdk.Spec{}, false
}
