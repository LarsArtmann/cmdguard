package v4

import (
	"context"
	"slices"
	"testing"
)

type probeConfig struct {
	Store  string `default:"memory" flag:"store"  help:"store backend"`
	DBPath string `default:"x.db"   flag:"db-path" help:"db path"`
}

func TestProbePointerFlags(t *testing.T) {
	type checkinFlags struct {
		TeamID   string   `default:"" flag:"team"   help:"team"  short:"t"`
		MemberID string   `default:"" flag:"member" help:"member" short:"m"`
		Answers  []string `default:"" flag:"answer" help:"answers" short:"a"`
	}

	var got []string

	cli, err := NewCLI[probeConfig]("probetest3", "probe3", probeConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, err := NewCommand(
		"checkin",
		&checkinFlags{},
		func(_ context.Context, _ *probeConfig, flags *checkinFlags) error {
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
		"checkin", "-t", "s2", "-m", "bob", "-a", "yesterday=Wrote", "-a", "blockers=ok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Equal(got, []string{"yesterday=Wrote", "blockers=ok"}) {
		t.Errorf("answers = %q, want [yesterday=Wrote blockers=ok]", got)
	}
}
