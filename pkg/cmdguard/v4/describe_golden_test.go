package v4

import (
	"bytes"
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// updateGolden rewrites the golden describe fixtures. Any intentional change
// to the describe contract must regenerate them and be documented in the
// CHANGELOG; see the schema stability rule on CLIDescription.
var updateGolden = flag.Bool("update-golden", false, "rewrite golden describe fixture files")

// describeFixtureConfig exercises every FlagDescription shape on the root:
// enum values, env, count, defaults, and a hidden flag.
type describeFixtureConfig struct {
	Environment string `flag:"environment" short:"e" default:"dev" env:"FIXTURE_ENV"       values:"dev,staging,prod" help:"Target environment"`
	Verbose     int    `flag:"verbose"     short:"v" default:"0"                                                     help:"Verbosity level; repeat for more"  count:"true"`
	APIToken    string `flag:"api-token"             default:""    env:"FIXTURE_API_TOKEN"                           help:"API token for authenticated calls"`
	Quiet       bool   `flag:"quiet"                                                                                 help:"Suppress non-error output"                      hidden:"true"`
}

// deployFixtureFlags exercises required, validate, and typed-duration flags
// on a leaf command.
type deployFixtureFlags struct {
	Region   string        `flag:"region"   env:"FIXTURE_REGION" required:"true" help:"Region to deploy to"`
	Replicas int           `flag:"replicas"                                      help:"Replica count"       default:"1"   validate:"min=1,max=10"`
	Timeout  time.Duration `flag:"timeout"                                       help:"Deploy timeout"      default:"30s"`
}

// newDescribeFixtureCLI builds a deterministic CLI that exercises the whole
// describe contract: nested commands, aliases, inherited root flags, leaf
// flag metadata, and a version.
func newDescribeFixtureCLI(t *testing.T, opts ...CLIOption) *CLI[describeFixtureConfig] {
	t.Helper()

	cliOpts := append([]CLIOption{
		WithFang(false),
		WithCLIVersion("9.9.9"),
		WithCLILong("Fixture CLI whose describe output is pinned byte-for-byte."),
	}, opts...)

	cli, err := NewCLI("fixture", "Golden describe fixture", describeFixtureConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	deploy, err := NewCommand(
		"deploy",
		deployFixtureFlags{},
		func(_ context.Context, _ *describeFixtureConfig, _ deployFixtureFlags) error { return nil },
		WithShort("Deploy the service"),
		WithExample("fixture deploy --region eu-central-1"),
		WithAliases("release"),
	)
	testutil.AssertNoError(t, err)

	list, err := NewCommand(
		"list",
		NoFlags{},
		func(_ context.Context, _ *describeFixtureConfig, _ NoFlags) error { return nil },
		WithShort("List remotes"),
		WithExample("fixture remote list"),
	)
	testutil.AssertNoError(t, err)

	remote, err := NewParentCommand[describeFixtureConfig, NoFlags](
		"remote",
		"Remote operations",
		NoFlags{},
		WithShort("Remote operations"),
		WithSubcommands[describeFixtureConfig, NoFlags](list),
	)
	testutil.AssertNoError(t, err)

	testutil.AssertNoError(t, AddCommands(cli, remote, deploy))

	return cli
}

func TestDescribe_GoldenFixtureByteStable(t *testing.T) {
	t.Parallel()

	cli := newDescribeFixtureCLI(t)

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))

	assertGoldenBytes(t, "describe_golden.json", out.Bytes())
}

func TestDescribe_OptOutGoldenFixture(t *testing.T) {
	t.Parallel()

	cli := newDescribeFixtureCLI(t, WithoutDescribeCommand())

	var buf bytes.Buffer

	testutil.AssertNoError(t, cli.WriteDescription(&buf))

	assertGoldenBytes(t, "describe_optout_golden.json", buf.Bytes())

	if bytes.Contains(buf.Bytes(), []byte(`"name": "describe"`)) {
		t.Error("opt-out fixture must not contain a describe command entry")
	}
}

func TestDescribe_OptOutExecuteRejectsDescribeCommand(t *testing.T) {
	t.Parallel()

	cli := newDescribeFixtureCLI(t, WithoutDescribeCommand())

	cli.RootCommand().SetOut(&bytes.Buffer{})

	err := cli.ExecuteWithArgs(context.Background(), []string{"describe"})
	if err == nil {
		t.Fatal("expected an unknown-command error for describe under WithoutDescribeCommand")
	}
}

func TestDescribe_ConsumerCommandWinsGoldenFixture(t *testing.T) {
	t.Parallel()

	cli := newDescribeFixtureCLI(t)

	mine, err := NewCommand(
		"describe",
		NoFlags{},
		func(_ context.Context, _ *describeFixtureConfig, _ NoFlags) error { return nil },
		WithShort("Consumer-owned describe"),
		WithExample("fixture describe"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, mine))

	var buf bytes.Buffer

	testutil.AssertNoError(t, cli.WriteDescription(&buf))

	assertGoldenBytes(t, "describe_consumer_wins_golden.json", buf.Bytes())

	if !bytes.Contains(buf.Bytes(), []byte("Consumer-owned describe")) {
		t.Error("consumer-wins fixture must carry the consumer's describe command")
	}
}

// assertGoldenBytes compares got against the pinned fixture file, rewriting it
// when -update-golden is passed.
func assertGoldenBytes(t *testing.T, name string, got []byte) {
	t.Helper()

	path := filepath.Join("testdata", name)

	if *updateGolden {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatalf("creating testdata dir: %v", err)
		}

		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("writing golden fixture: %v", err)
		}

		t.Logf("updated golden fixture %s", path)

		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading golden fixture (run with -update-golden to create it): %v", err)
	}

	if !bytes.Equal(want, got) {
		t.Errorf(
			"describe output drifted from %s; run go test ./pkg/cmdguard/v4 -run TestDescribe -update-golden and review the diff before committing",
			path,
		)
	}
}
