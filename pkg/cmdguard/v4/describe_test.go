package v4

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type describeConfig struct {
	Name      string `flag:"name"       short:"n" default:"world" help:"Name"      env:"NAME"`
	Verbose   int    `flag:"verbose"    short:"v"                 help:"Verbosity"            count:"true"`
	Secret    string `flag:"secret"                               help:"Hidden"                            hidden:"true"`
	LocalOnly string `flag:"local-only"                           help:"Root only"                                       local:"true"`
}

type describeScanFlags struct {
	Mode string `flag:"mode" default:"fast" help:"Scan mode"   values:"fast,slow" env:"MODE"`
	Tag  string `flag:"tag"                 help:"Tag to scan"                               required:"true"`
}

func newDescribeCLI(t *testing.T, opts ...CLIOption) *CLI[describeConfig] {
	t.Helper()

	cliOpts := append([]CLIOption{
		WithFang(false),
		WithCLIVersion("1.2.3"),
		WithEnvPrefix("APP_"),
		WithSilenceErrors(),
	}, opts...)

	cli, err := NewCLI("app", "App short", describeConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	return cli
}

func TestDescribe_Tree(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t)

	scan, err := NewCommand(
		"scan",
		describeScanFlags{},
		func(_ context.Context, _ *describeConfig, _ describeScanFlags) error { return nil },
		WithShort("Scan things"),
		WithGroupID("ops"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, scan))

	desc := cli.Describe()

	testutil.AssertFieldEqString(t, desc.SchemaVersion, descriptionSchemaVersion, "SchemaVersion")
	testutil.AssertFieldEqString(t, desc.Name, "app", "Name")
	testutil.AssertFieldEqString(t, desc.Short, "App short", "Short")
	testutil.AssertFieldEqString(t, desc.Version, "1.2.3", "Version")

	rootFlags := flagsByName(desc.Flags)

	nameFlag, ok := rootFlags["name"]
	if !ok {
		t.Fatal("expected root flag \"name\" in description")
	}

	testutil.AssertFieldEqString(t, nameFlag.Env, "APP_NAME", "name env")
	testutil.AssertFieldEqString(t, nameFlag.Shorthand, "n", "name shorthand")
	testutil.AssertFieldEqString(t, nameFlag.Default, "world", "name default")
	testutil.AssertFieldEqString(t, nameFlag.Type, "string", "name type")

	verboseFlag, ok := rootFlags["verbose"]
	if !ok {
		t.Fatal("expected root flag \"verbose\" in description")
	}

	if !verboseFlag.Count {
		t.Error("expected verbose to be a counting flag")
	}

	if secret := rootFlags["secret"]; !secret.Hidden {
		t.Error("expected secret flag to be hidden")
	}

	if len(desc.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(desc.Commands))
	}

	scanDesc := desc.Commands[0]
	testutil.AssertFieldEqString(t, scanDesc.Name, "scan", "command name")
	testutil.AssertFieldEqString(t, scanDesc.Path, "app scan", "command path")
	testutil.AssertFieldEqString(t, scanDesc.Short, "Scan things", "command short")
	testutil.AssertFieldEqString(t, scanDesc.Group, "ops", "command group")

	scanFlags := flagsByName(scanDesc.Flags)

	modeFlag, ok := scanFlags["mode"]
	if !ok {
		t.Fatal("expected scan flag \"mode\" in description")
	}

	testutil.AssertFieldEqString(t, modeFlag.Env, "APP_MODE", "mode env")
	testutil.AssertStringSlicesEqual(t, modeFlag.Values, []string{"fast", "slow"}, "mode values")

	tagFlag, ok := scanFlags["tag"]
	if !ok {
		t.Fatal("expected scan flag \"tag\" in description")
	}

	if !tagFlag.Required {
		t.Error("expected tag flag to be required")
	}

	// Root persistent flags are inherited by subcommands; root-local flags are not.
	if inherited, ok := scanFlags["name"]; !ok || !inherited.Inherited {
		t.Error("expected root persistent flag \"name\" to be inherited by scan")
	}

	if _, ok := scanFlags["local-only"]; ok {
		t.Error("root-local flag \"local-only\" must not appear on scan")
	}
}

func TestWriteDescription_EmitsValidSchemaJSON(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t)

	var buf bytes.Buffer

	testutil.AssertNoError(t, cli.WriteDescription(&buf))

	var decoded map[string]any
	testutil.AssertNoError(t, json.Unmarshal(buf.Bytes(), &decoded))

	if decoded["schema_version"] != descriptionSchemaVersion {
		t.Errorf("schema_version = %v, want %q", decoded["schema_version"], descriptionSchemaVersion)
	}

	if decoded["name"] != "app" {
		t.Errorf("name = %v, want app", decoded["name"])
	}

	if !strings.HasSuffix(buf.String(), "\n") {
		t.Error("description JSON must end with a newline")
	}
}

func TestDescribeCommand_AutoRegistered(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t)

	list, err := NewCommand(
		"list",
		NoFlags{},
		func(_ context.Context, _ *describeConfig, _ NoFlags) error { return nil },
		WithShort("List things"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, list))

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))

	var desc CLIDescription
	testutil.AssertNoError(t, json.Unmarshal(out.Bytes(), &desc))

	testutil.AssertFieldEqString(t, desc.Name, "app", "described name")

	names := make([]string, 0, len(desc.Commands))
	for _, cmd := range desc.Commands {
		names = append(names, cmd.Name)
	}

	if !strings.Contains(strings.Join(names, ","), "list") {
		t.Errorf("expected describe output to include %q, got %v", "list", names)
	}

	if !strings.Contains(strings.Join(names, ","), "describe") {
		t.Errorf("expected describe command to document itself, got %v", names)
	}
}

func TestDescribeCommand_OptOut(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t, WithoutDescribeCommand())

	cli.RootCommand().SetOut(&bytes.Buffer{})

	// "describe" is not a registered command, so cobra treats it as an
	// argument to the root (which prints help) rather than dispatching it.
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))

	for _, cmd := range cli.RootCommand().Commands() {
		if cmd.Name() == describeCommandName {
			t.Fatal("describe command must not be registered when disabled")
		}
	}

	// The programmatic contract remains available.
	if cli.Describe() == nil {
		t.Fatal("Describe must work even when the describe command is disabled")
	}
}

func TestDescribeCommand_ConsumerCommandWins(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t)

	mine, err := NewCommand(
		"describe",
		NoFlags{},
		func(_ context.Context, _ *describeConfig, _ NoFlags) error {
			_, writeErr := fmt.Fprintln(cli.RootCommand().OutOrStdout(), "mine")

			return writeErr
		},
		WithShort("Consumer describe"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, mine))

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))

	testutil.AssertFieldEqString(t, out.String(), "mine\n", "consumer describe output")
}

func TestDescribeCommand_IdempotentAcrossExecutes(t *testing.T) {
	t.Parallel()

	cli := newDescribeCLI(t)

	var out bytes.Buffer

	cli.RootCommand().SetOut(&out)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))
	testutil.AssertNoError(t, cli.ExecuteWithArgs(context.Background(), []string{"describe"}))
}

func TestAddCommands_MixedFlagTypes(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	type alphaFlags struct {
		Alpha string `flag:"alpha" default:"a" help:"Alpha"`
	}

	type betaFlags struct {
		Beta int `flag:"beta" default:"2" help:"Beta"`
	}

	alpha, err := NewCommand(
		"alpha",
		alphaFlags{},
		func(_ context.Context, _ *testConfig, _ alphaFlags) error { return nil },
		WithShort("Alpha command"),
	)
	testutil.AssertNoError(t, err)

	beta, err := NewCommand(
		"beta",
		betaFlags{},
		func(_ context.Context, _ *testConfig, _ betaFlags) error { return nil },
		WithShort("Beta command"),
	)
	testutil.AssertNoError(t, err)

	testutil.AssertNoError(t, AddCommands(cli, alpha, beta))

	testutil.AssertEqual(t, len(cli.RootCommand().Commands()), 2)
}

func TestAddCommands_AtomicOnDuplicate(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	first, err := NewCommand(
		"dup",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("First"),
	)
	testutil.AssertNoError(t, err)

	second, err := NewCommand(
		"dup",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("Second"),
	)
	testutil.AssertNoError(t, err)

	err = AddCommands(cli, first, second)
	if !errors.Is(err, ErrDuplicateCommand) {
		t.Fatalf("expected ErrDuplicateCommand, got: %v", err)
	}

	if got := len(cli.RootCommand().Commands()); got != 0 {
		t.Errorf("failed batch must not attach commands, got %d attached", got)
	}

	// The failed batch must not reserve names either.
	testutil.AssertNoError(t, AddCommands(cli, first))
}

func TestAddCommands_DuplicateUseLineNames(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	withArgs, err := NewCommand(
		"list [path]",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("List a path"),
	)
	testutil.AssertNoError(t, err)

	plain, err := NewCommand(
		"list",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("List everything"),
	)
	testutil.AssertNoError(t, err)

	err = AddCommands(cli, withArgs, plain)
	if !errors.Is(err, ErrDuplicateCommand) {
		t.Fatalf("expected ErrDuplicateCommand for name collision, got: %v", err)
	}
}

func TestAddCommand_DuplicateUseLineNames(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	withArgs, err := NewCommand(
		"list [path]",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("List a path"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, withArgs))

	plain, err := NewCommand(
		"list",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("List everything"),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, plain)
	if !errors.Is(err, ErrDuplicateCommand) {
		t.Fatalf("expected ErrDuplicateCommand, got: %v", err)
	}
}

func TestAddCommands_NilCommand(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	err = AddCommands[testConfig](cli, nil)
	if !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("expected ErrInvalidCommand, got: %v", err)
	}
}

func TestAddCommand_FailedConversionReleasesName(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", testConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	type unsupportedFlags struct {
		Value string `flag:"value" required:"not-a-bool" help:"Invalid bool tag"`
	}

	bad, err := NewCommand(
		"bad",
		unsupportedFlags{},
		func(_ context.Context, _ *testConfig, _ unsupportedFlags) error { return nil },
		WithShort("Bad"),
	)
	testutil.AssertNoError(t, err)

	if addErr := AddCommand(cli, bad); addErr == nil {
		t.Fatal("expected flag registration to fail for complex128 flag")
	}

	good, err := NewCommand(
		"bad",
		NoFlags{},
		func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
		WithShort("Good"),
	)
	testutil.AssertNoError(t, err)

	if addErr := AddCommand(cli, good); addErr != nil {
		t.Fatalf("failed conversion must not reserve the command name, got: %v", addErr)
	}
}

// flagsByName indexes flag descriptions by flag name for assertions.
func flagsByName(flags []FlagDescription) map[string]FlagDescription {
	out := make(map[string]FlagDescription, len(flags))

	for _, flag := range flags {
		out[flag.Name] = flag
	}

	return out
}
