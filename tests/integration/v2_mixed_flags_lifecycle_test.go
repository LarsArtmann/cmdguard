// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

func TestV2_MixedFlagTypes_WithLifecycleHooks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v5.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		preRunCalled  bool
		runCalled     bool
		postRunCalled bool
		receivedFlags *GreetFlags
	)

	greetCmd, err := v5.NewCommand(
		"greet",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			runCalled = true
			receivedFlags = flags

			return nil
		},
		v5.WithShort("Greet with lifecycle"),
		v5.WithPreRunE(
			func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				preRunCalled = true

				return nil
			},
		),
		v5.WithPostRunE(
			func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				postRunCalled = true

				return nil
			},
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v5.AddCommand(cli, greetCmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(ctx, []string{cmdNameGreet, "--name=TestUser", flagShout})
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

	_, err := v5.NewCommand(
		"",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, _ *GreetFlags) error { return nil },
		v5.WithShort("Invalid command"),
	)
	if err == nil {
		t.Error("expected error for empty Use field")
	}

	_, err = v5.NewCommand[RootConfig](
		"invalid",
		v5.NoFlags{},
		nil,
		v5.WithShort("No handler"),
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

	cli, err := v5.NewCLI[RootConfig]("testapp", "Test application", defaultConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var receivedConfig *RootConfig

	checkCmd, err := v5.NewCommand(
		"check",
		&GreetFlags{},
		func(_ context.Context, cfg *RootConfig, _ *GreetFlags) error {
			receivedConfig = cfg

			return nil
		},
		v5.WithShort("Check config access"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v5.AddCommand(cli, checkCmd)
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

	cli, err := v5.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executedFlags *MigrateFlags

	migrateUpCmd, err := v5.NewCommand(
		"up",
		&MigrateFlags{},
		func(_ context.Context, _ *RootConfig, flags *MigrateFlags) error {
			executedFlags = flags

			return nil
		},
		v5.WithShort("Run up migrations"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	migrateCmd, err := v5.NewParentCommand[RootConfig](
		"migrate",
		"Database migration management commands", &MigrateFlags{},
		v5.WithSubcommands(migrateUpCmd),
		v5.WithShort("Migration commands"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v5.AddCommand(cli, migrateCmd)
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
