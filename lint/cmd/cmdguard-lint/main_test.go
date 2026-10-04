package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLintCommandCleanTreeExitsZero(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "clean.go"), []byte("package clean\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"lint", "--dir", dir}); execErr != nil {
		t.Errorf("expected clean tree to pass, got: %v", execErr)
	}
}

func TestLintCommandFindingsFailTheRun(t *testing.T) {
	dir := t.TempDir()

	stale := "package x\n\nimport \"github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3\"\n\nvar _ = v3.NoFlags{}\n"

	if err := os.WriteFile(filepath.Join(dir, "stale.go"), []byte(stale), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	execErr := cli.ExecuteWithArgs(context.Background(), []string{"lint", "--dir", dir})
	if execErr == nil {
		t.Fatal("expected findings to fail the run")
	}

	if !strings.Contains(execErr.Error(), "finding") {
		t.Errorf("error should mention findings, got: %v", execErr)
	}
}

func TestRulesCommandListsAllRules(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	cli, err := newApp()
	if err != nil {
		t.Fatalf("newApp: %v", err)
	}

	if execErr := cli.ExecuteWithArgs(context.Background(), []string{"rules"}); execErr != nil {
		t.Errorf("rules command failed: %v", execErr)
	}
}
