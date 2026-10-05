package v4

import (
	"context"
	"slices"
	"testing"
)

type probeAppConfig struct {
	Notifier    string `default:"stdout" flag:"notifier" help:"Notifier backend"`
	PapURL      string `default:""       flag:"pap-url"  help:"Pap URL"`
	Store       string `default:"memory" flag:"store"    help:"Store backend"`
	DBPath      string `default:"sk.db"  flag:"db-path"  help:"DB path"`
	GitHubToken string `default:""       flag:"github-token" help:"token"`
}

type probeCheckinFlags struct {
	TeamID   string   `default:"" flag:"team"   help:"Team ID (required)"   short:"t"`
	MemberID string   `default:"" flag:"member" help:"Member ID (required)" short:"m"`
	Answers  []string `default:"" flag:"answer" help:"Answers as questionID=text (repeatable)" short:"a"`
}

func TestProbeExactMirror(t *testing.T) {
	var got []string

	opts := []CLIOption{
		WithFang(false),
		WithCLIVersion("0.1.0"),
		WithCLILong("Async standup automation — replace synchronous daily standup meetings with asynchronous check-ins."),
	}

	cli, err := NewCLI[probeAppConfig]("standup-killer", "Kill your daily standup meeting", probeAppConfig{}, opts...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, err := NewCommand[probeAppConfig](
		"checkin",
		&probeCheckinFlags{TeamID: "", MemberID: "", Answers: nil},
		func(_ context.Context, _ *probeAppConfig, flags *probeCheckinFlags) error {
			got = flags.Answers
			return nil
		},
		WithShort("Submit a standup check-in"),
		WithLong("Submit a check-in with answers for today's standup. Use --answer flag repeatedly with questionID=text format, e.g. --answer yesterday='Fixed bug' --answer today='Code review' --answer blockers='None'"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := AddCommand(cli, cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(t.Context(), []string{
		"--store=sqlite", "--db-path=/tmp/x.db",
		"checkin", "-t", "s2", "-m", "bob", "-a", "yesterday=Wrote", "-a", "blockers=ok",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Equal(got, []string{"yesterday=Wrote", "blockers=ok"}) {
		t.Errorf("answers = %q, want [yesterday=Wrote blockers=ok]", got)
	}
}
