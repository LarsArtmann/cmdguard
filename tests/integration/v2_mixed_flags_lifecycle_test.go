// Package integration provides end-to-end tests for cmdguard.
package integration

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

func TestV2_MixedFlagTypes_WithLifecycleHooks(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	cli, err := v2.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var (
		preRunCalled  bool
		runCalled     bool
		postRunCalled bool
		receivedFlags *GreetFlags
	)

	greetCmd, err := v2.NewCommand(
		"greet",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, flags *GreetFlags) error {
			runCalled = true
			receivedFlags = flags

			return nil
		},
		v2.WithShort("Greet with lifecycle"),
		v2.WithPreRunE(
			func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				preRunCalled = true

				return nil
			},
		),
		v2.WithPostRunE(
			func(_ context.Context, _ *RootConfig, _ *GreetFlags) error {
				postRunCalled = true

				return nil
			},
		),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, greetCmd)
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

	_, err := v2.NewCommand(
		"",
		&GreetFlags{},
		func(_ context.Context, _ *RootConfig, _ *GreetFlags) error { return nil },
		v2.WithShort("Invalid command"),
	)
	if err == nil {
		t.Error("expected error for empty Use field")
	}

	_, err = v2.NewCommand[RootConfig](
		"invalid",
		v2.NoFlags{},
		nil,
		v2.WithShort("No handler"),
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

	cli, err := v2.NewCLI[RootConfig]("testapp", "Test application", defaultConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var receivedConfig *RootConfig

	checkCmd, err := v2.NewCommand(
		"check",
		&GreetFlags{},
		func(_ context.Context, cfg *RootConfig, _ *GreetFlags) error {
			receivedConfig = cfg

			return nil
		},
		v2.WithShort("Check config access"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, checkCmd)
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

	cli, err := v2.NewCLI[RootConfig]("testapp", "Test application", RootConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var executedFlags *MigrateFlags

	migrateUpCmd, err := v2.NewCommand(
		"up",
		&MigrateFlags{},
		func(_ context.Context, _ *RootConfig, flags *MigrateFlags) error {
			executedFlags = flags

			return nil
		},
		v2.WithShort("Run up migrations"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	migrateCmd, err := v2.NewParentCommand[RootConfig](
		"migrate",
		"Database migration management commands", &MigrateFlags{},
		v2.WithSubcommands(migrateUpCmd),
		v2.WithShort("Migration commands"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = v2.AddCommand(cli, migrateCmd)
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
