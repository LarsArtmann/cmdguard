package v5

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

const styledMarker = "\x1b[38;5;228m-styled-"

// styleTestTree mimics what terminal help transforms (glamour, lipgloss) do:
// overwrite the live cobra Long and Example fields with styled bytes.
func styleTestTree(cmd *cobra.Command) {
	for _, sub := range cmd.Commands() {
		styleTestTree(sub)
	}

	if cmd.Long != "" {
		cmd.Long = styledMarker + cmd.Long
	}

	if cmd.Example != "" {
		cmd.Example = styledMarker + cmd.Example
	}
}

type pristineConfig struct {
	Name string `flag:"name" default:"world" help:"Name"`
}

// newPristineCLI registers a leaf whose Long and Example are plain markdown;
// opts decide whether a styling transform is attached.
func newPristineCLI(t *testing.T, opts ...CLIOption) (*CLI[pristineConfig], string, string) {
	t.Helper()

	const wantLong = "# Scan\n\nSweep every registered host."
	const wantExample = "app scan --mode fast"

	cliOpts := append([]CLIOption{
		WithFang(false),
		WithCLIVersion("1.2.3"),
		WithSilenceErrors(),
	}, opts...)

	cli, err := NewCLI("app", "App short", pristineConfig{}, cliOpts...)
	testutil.AssertNoError(t, err)

	scan, err := NewCommand(
		"scan",
		NoFlags{},
		func(_ context.Context, _ *pristineConfig, _ NoFlags) error { return nil },
		WithShort("Scan things"),
		WithLong(wantLong),
		WithExample(wantExample),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, scan))

	return cli, wantLong, wantExample
}

func scanDescription(desc *CLIDescription) *CommandDescription {
	for i := range desc.Commands {
		if desc.Commands[i].Name == "scan" {
			return &desc.Commands[i]
		}
	}

	return nil
}

func TestPristineHelp_DescribeServesPreTransformText(t *testing.T) {
	t.Parallel()

	cli, wantLong, wantExample := newPristineCLI(t, WithHelpTransform(func(cmd *cobra.Command) {
		styleTestTree(cmd)
	}))

	testutil.AssertNoError(t, cli.ExecuteWithArgs(t.Context(), []string{"describe"}))

	scan := scanDescription(cli.Describe())
	if scan == nil {
		t.Fatal("expected scan command in description")
	}

	if scan.Long != wantLong {
		t.Errorf("describe Long must serve the pre-transform text: got %q, want %q", scan.Long, wantLong)
	}

	if scan.Example != wantExample {
		t.Errorf("describe Example must serve the pre-transform text: got %q, want %q", scan.Example, wantExample)
	}

	live := cli.rootCmd

	for _, sub := range live.Commands() {
		if sub.Name() == "scan" && !strings.Contains(sub.Long, styledMarker) {
			t.Errorf("transform must still style live fields for terminal help: got %q", sub.Long)
		}
	}
}

func TestPristineHelp_RepeatExecuteKeepsOriginalSnapshot(t *testing.T) {
	t.Parallel()

	cli, wantLong, wantExample := newPristineCLI(t, WithHelpTransform(func(cmd *cobra.Command) {
		styleTestTree(cmd)
	}))

	for range 3 {
		testutil.AssertNoError(t, cli.ExecuteWithArgs(t.Context(), []string{"describe"}))
	}

	scan := scanDescription(cli.Describe())
	if scan == nil {
		t.Fatal("expected scan command in description")
	}

	if scan.Long != wantLong {
		t.Errorf("repeated Execute must not re-stash styled text: got %q, want %q", scan.Long, wantLong)
	}

	if scan.Example != wantExample {
		t.Errorf("unstyled Example must round-trip unchanged: got %q, want %q", scan.Example, wantExample)
	}
}

func TestPristineHelp_GenerateDocsServesPreTransformText(t *testing.T) {
	t.Parallel()

	cli, wantLong, _ := newPristineCLI(t, WithHelpTransform(func(cmd *cobra.Command) {
		styleTestTree(cmd)
	}))

	testutil.AssertNoError(t, cli.ExecuteWithArgs(t.Context(), []string{"describe"}))

	var docs strings.Builder

	testutil.AssertNoError(t, cli.GenerateDocs(&docs))

	if strings.Contains(docs.String(), styledMarker) {
		t.Error("GenerateDocs must not carry styled bytes")
	}

	if !strings.Contains(docs.String(), wantLong) {
		t.Errorf("GenerateDocs must serve the pre-transform Long: missing %q", wantLong)
	}
}

func TestPristineHelp_NoTransformPassesThrough(t *testing.T) {
	t.Parallel()

	cli, wantLong, wantExample := newPristineCLI(t)

	testutil.AssertNoError(t, cli.ExecuteWithArgs(t.Context(), []string{"describe"}))

	scan := scanDescription(cli.Describe())
	if scan == nil {
		t.Fatal("expected scan command in description")
	}

	if scan.Long != wantLong || scan.Example != wantExample {
		t.Errorf("untouched CLI must round-trip: got long=%q example=%q", scan.Long, scan.Example)
	}
}
