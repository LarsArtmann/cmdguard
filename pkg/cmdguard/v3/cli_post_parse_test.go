package v3

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestWithPostFlagParse_RunsAfterFlagParsing(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `flag:"name" default:"world" help:"name"`
	}

	var seenCfg *cfg
	var seenCmd *cobra.Command

	cli, err := NewCLI[cfg](
		"test", "Test", cfg{},
		WithPostFlagParse(func(cmd *cobra.Command, c *cfg) error {
			seenCfg = c
			seenCmd = cmd

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub", "--name", "alice"})
	testutil.AssertNoError(t, err)

	if seenCfg == nil {
		t.Fatal("post-flag-parse hook was not called")
	}

	if seenCfg.Name != "alice" {
		t.Errorf("config Name = %q, want %q (flag should be parsed before hook)", seenCfg.Name, "alice")
	}

	if seenCmd == nil {
		t.Error("cmd parameter was nil")
	}
}

func TestWithPostFlagParse_MultipleHooksInOrder(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	var order []int

	cli, err := NewCLI[cfg](
		"test", "Test", cfg{},
		WithPostFlagParse(func(*cobra.Command, *cfg) error {
			order = append(order, 1)

			return nil
		}),
		WithPostFlagParse(func(*cobra.Command, *cfg) error {
			order = append(order, 2)

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertNoError(t, err)

	if len(order) != 2 || order[0] != 1 || order[1] != 2 {
		t.Errorf("hooks ran in order %v, want [1 2]", order)
	}
}

func TestWithPostFlagParse_ErrorStopsExecution(t *testing.T) {
	t.Parallel()

	type cfg struct{}

	handlerRan := false

	cli, err := NewCLI[cfg](
		"test", "Test", cfg{},
		WithPostFlagParse(func(*cobra.Command, *cfg) error {
			return errors.New("init failed")
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{
		Use: "sub",
		RunE: func(*cobra.Command, []string) error {
			handlerRan = true

			return nil
		},
	}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertExpectedError(t, err)

	if handlerRan {
		t.Error("command handler should NOT run when post-flag-parse hook errors")
	}
}

func TestWithPostFlagParse_RunsAfterConfigValidation(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Mode string `flag:"mode" default:"full" help:"mode"`
	}

	validateRan := false
	postParseRan := false

	cli, err := NewCLI[cfg](
		"test", "Test", cfg{},
		WithConfigValidation(func(c *cfg) error {
			validateRan = true

			return nil
		}),
		WithPostFlagParse(func(*cobra.Command, *cfg) error {
			postParseRan = true

			return nil
		}),
	)
	testutil.AssertNoError(t, err)

	subCmd := &cobra.Command{Use: "sub", RunE: func(*cobra.Command, []string) error { return nil }}
	cli.RootCommand().AddCommand(subCmd)

	err = cli.ExecuteWithArgs(context.Background(), []string{"sub"})
	testutil.AssertNoError(t, err)

	if !validateRan {
		t.Error("config validation should have run")
	}

	if !postParseRan {
		t.Error("post-flag-parse hook should have run")
	}
}
