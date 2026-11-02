package v4

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type reuseConfig struct {
	Apply     bool   `flag:"apply"     help:"Apply changes"`
	Directory string `flag:"directory" help:"Working directory" default:"/tmp"`
}

func TestExecuteWithArgs_FlagValuesDoNotLeakAcrossExecutes(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("reuse", "Reuse", reuseConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cli.RootCommand().RunE = func(_ *cobra.Command, _ []string) error { return nil }
	cli.RootCommand().SetErr(io.Discard)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"--apply"}))
	if !cli.config.Apply {
		t.Fatal("first execute: --apply should be true")
	}

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), nil))
	if cli.config.Apply {
		t.Fatal("second execute: --apply leaked from the first parse, want false")
	}
	if cli.config.Directory != "/tmp" {
		t.Fatalf("second execute: directory = %q, want default %q", cli.config.Directory, "/tmp")
	}
}

func TestExecuteWithArgs_SecondExecuteSeesFreshValues(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("reuse2", "Reuse2", reuseConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cli.RootCommand().RunE = func(_ *cobra.Command, _ []string) error { return nil }
	cli.RootCommand().SetErr(io.Discard)

	testutil.AssertNoError(t,
		cli.ExecuteWithArgs(context.Background(), []string{"--directory", "/first"}))
	if cli.config.Directory != "/first" {
		t.Fatalf("first execute: directory = %q, want /first", cli.config.Directory)
	}

	testutil.AssertNoError(t,
		cli.ExecuteWithArgs(context.Background(), []string{"--directory", "/second"}))
	if cli.config.Directory != "/second" {
		t.Fatalf("second execute: directory = %q, want /second", cli.config.Directory)
	}
}

func TestExecuteWithArgs_SliceFlagDoesNotAccumulateAcrossExecutes(t *testing.T) {
	t.Parallel()

	type sliceConfig struct {
		IDs []string `flag:"ids" help:"repeatable IDs"`
	}

	cli, err := NewCLI("reuseslice", "ReuseSlice", sliceConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cli.RootCommand().RunE = func(_ *cobra.Command, _ []string) error { return nil }
	cli.RootCommand().SetErr(io.Discard)

	testutil.AssertNoError(t,
		cli.ExecuteWithArgs(context.Background(), []string{"--ids", "1", "--ids", "2"}))
	if len(cli.config.IDs) != 2 || cli.config.IDs[0] != "1" || cli.config.IDs[1] != "2" {
		t.Fatalf("first execute: ids = %q, want [1 2]", cli.config.IDs)
	}

	testutil.AssertNoError(t,
		cli.ExecuteWithArgs(context.Background(), []string{"--ids", "3"}))
	if len(cli.config.IDs) != 1 || cli.config.IDs[0] != "3" {
		t.Fatalf("second execute: ids = %q, want [3] (pflag slice Set appends across parses)", cli.config.IDs)
	}
}

func TestExecuteWithArgs_RequiredFlagNotSatisfiedByEarlierParse(t *testing.T) {
	t.Parallel()

	type requiredConfig struct {
		Token string `flag:"token" required:"true" help:"Token"`
	}

	cli, err := NewCLI("reusereq", "ReuseReq", requiredConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	cli.RootCommand().RunE = func(_ *cobra.Command, _ []string) error { return nil }
	cli.RootCommand().SetErr(io.Discard)

	testutil.AssertNoError(t,
		cli.ExecuteWithArgs(context.Background(), []string{"--token", "secret"}))

	err = cli.ExecuteWithArgs(context.Background(), nil)
	if !errors.Is(err, ErrRequiredFlag) {
		t.Fatalf("second execute without --token: expected ErrRequiredFlag, got: %v", err)
	}
}
