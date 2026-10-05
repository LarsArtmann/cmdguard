package v4

import (
	"context"
	"slices"
	"testing"
)

func TestProbeEqualsInSliceValues(t *testing.T) {
	type repeatFlags struct {
		Answers []string `flag:"answer" default:"" help:"repeatable answer" short:"a"`
	}
	type repeatConfig struct{}

	var got []string

	cli, err := NewCLI[repeatConfig]("probetest", "probe", repeatConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, err := NewCommand(
		"checkin",
		repeatFlags{},
		func(_ context.Context, _ *repeatConfig, flags repeatFlags) error {
			got = flags.Answers
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := AddCommand(cli, cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(t.Context(), []string{
		"checkin", "-a", "yesterday=Wrote", "-a", "blockers=ok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Equal(got, []string{"yesterday=Wrote", "blockers=ok"}) {
		t.Errorf("answers = %q, want [yesterday=Wrote blockers=ok]", got)
	}
}
