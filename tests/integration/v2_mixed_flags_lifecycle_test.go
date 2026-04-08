// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestV2_MixedFlagTypes_WithLifecycleHooks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		preRunCalled  bool
		runCalled     bool
		postRunCalled bool
		receivedFlags *GreetFlags
	)

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "greet",
			Short: "Greet with lifecycle",
			Flags: &GreetFlags{Name: "World", Shout: false},
			PreRunE: func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				preRunCalled = true

				return nil
			},
			RunE: func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
				runCalled = true
				receivedFlags = flags

				return nil
			},
			PostRunE: func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				postRunCalled = true

				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"greet", "--name=TestUser", "--shout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !preRunCalled {
		t.Error("preRunCalled should be true")
	}

	if !runCalled {
		t.Error("runCalled should be true")
	}

	if !postRunCalled {
		t.Error("postRunCalled should be true")
	}

	if receivedFlags.Name != "TestUser" {
		t.Errorf("receivedFlags.Name = %q, want %q", receivedFlags.Name, "TestUser")
	}

	if !receivedFlags.Shout {
		t.Error("receivedFlags.Shout should be true")
	}
}

func TestV2_MixedFlagTypes_ValidationErrors(t *testing.T) {
	t.Parallel()
	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "",
			Short: "Invalid command",
			RunE:  func(_ context.Context, _ *RootConfig, _ *GreetFlags) error { return nil },
		},
	)
	if err == nil {
		t.Error("expected error for empty Use field")
	}

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "invalid",
			Short: "No handler",
		},
	)
	if err == nil {
		t.Error("expected error for missing RunE")
	}
}

func TestV2_MixedFlagTypes_ConfigAccess(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	defaultConfig := RootConfig{
		Debug:   false,
		Verbose: true,
		Level:   "debug",
	}

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", defaultConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var receivedConfig *RootConfig

	err = v2.AddAnyCommand(cli,
		v2.Command[RootConfig, *GreetFlags]{
			Use:   "check",
			Short: "Check config access",
			Flags: &GreetFlags{},
			RunE: func(_ context.Context, cfg *RootConfig, _ *GreetFlags) error {
				receivedConfig = cfg

				return nil
			},
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"check"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedConfig == nil {
		t.Fatal("receivedConfig is nil")
	}

	if !receivedConfig.Verbose {
		t.Error("receivedConfig.Verbose should be true")
	}

	if receivedConfig.Level != "debug" {
		t.Errorf("receivedConfig.Level = %q, want %q", receivedConfig.Level, "debug")
	}
}

func TestV2_MixedFlagTypes_DeeplyNested(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.New[RootConfig, v2.NoFlags]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executedFlags *MigrateFlags

	migrateUpCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:   "up",
		Short: "Run up migrations",
		RunE: func(_ context.Context, _ *RootConfig, flags *MigrateFlags) error {
			executedFlags = flags

			return nil
		},
	}

	migrateCmd := v2.Command[RootConfig, *MigrateFlags]{
		Use:      "migrate",
		Short:    "Migration commands",
		Flags:    &MigrateFlags{Steps: 0, Direction: "up"},
		Commands: []v2.Command[RootConfig, *MigrateFlags]{migrateUpCmd},
	}

	err = v2.AddAnyCommand(cli, migrateCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{"migrate", "up", "--steps=3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if executedFlags.Steps != 3 {
		t.Errorf("executedFlags.Steps = %d, want %d", executedFlags.Steps, 3)
	}
}
