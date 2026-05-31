package main

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestDeployCommand_Compiles(t *testing.T) {
	t.Parallel()

	// Verify the command with WithPromptOnMissing compiles and validates.
	_, err := v2.NewCommand[AppConfig, *DeployFlags]("deploy",
		func(_ context.Context, _ *AppConfig, flags *DeployFlags) error {
			return nil
		},
		v2.WithShort[AppConfig, *DeployFlags]("Deploy an application"),
		v2.WithFlags[AppConfig, *DeployFlags](&DeployFlags{}),
		v2.WithPromptOnMissing[AppConfig, *DeployFlags](),
	)
	if err != nil {
		t.Fatalf("creating command: %v", err)
	}
}
