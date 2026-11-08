package v5

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

// docgenRootConfig exercises the metadata a generated doc must not drop.
type docgenRootConfig struct {
	Environment string `flag:"environment" short:"e" default:"dev" env:"DOCGEN_ENV"    values:"dev,staging,prod" help:"Target environment"`
	Verbose     int    `flag:"verbose"     short:"v" default:"0"                                                 help:"Verbosity level"    count:"true"`
	Secret      string `flag:"secret"                default:""    env:"DOCGEN_SECRET"                           help:"Credential"                      secret:"true"`
	Internal    bool   `flag:"internal"                                                                          help:"Hidden flag"                                   hidden:"true"`
}

// docgenLeafFlags exercises required, bounds, and typed defaults on a leaf.
type docgenLeafFlags struct {
	Region   string        `flag:"region"   env:"DOCGEN_REGION" required:"true" help:"Region"`
	Replicas int           `flag:"replicas"                                     help:"Replica count"  default:"1"   validate:"min=1,max=10"`
	Timeout  time.Duration `flag:"timeout"                                      help:"Deploy timeout" default:"30s"`
}

// newDocgenContractCLI builds a deterministic CLI used by the drift test.
func newDocgenContractCLI(t *testing.T) *CLI[docgenRootConfig] {
	t.Helper()

	cli, err := NewCLI("docapp", "Doc fixture", docgenRootConfig{}, WithFang(false))
	testutil.AssertNoError(t, err)

	deploy, err := NewCommand(
		"deploy",
		docgenLeafFlags{},
		func(_ context.Context, _ *docgenRootConfig, _ docgenLeafFlags) error { return nil },
		WithShort("Deploy the service"),
		WithExample("docapp deploy --region eu-central-1"),
	)
	testutil.AssertNoError(t, err)
	testutil.AssertNoError(t, AddCommand(cli, deploy))

	return cli
}

// flagRows returns every markdown table row whose first cell is the flag.
// A flag can be rendered once per command, so the result is a slice.
func flagRows(docs, name string) []string {
	var rows []string

	for line := range strings.SplitSeq(docs, "\n") {
		if strings.HasPrefix(line, "| `--"+name+"`") {
			rows = append(rows, line)
		}
	}

	return rows
}

// collectDescriptions flattens the describe contract into one list.
func collectDescriptions(descriptions []CommandDescription) []CommandDescription {
	out := make([]CommandDescription, 0, len(descriptions))

	for _, description := range descriptions {
		out = append(out, description)
		out = append(out, collectDescriptions(description.Commands)...)
	}

	return out
}

// TestGenerateDocs_MatchesDescribeContract is the drift gate: every flag fact
// the machine contract carries must reach the generated markdown. Adding a
// metadata field to FlagDescription without teaching docgen about it fails here.
func TestGenerateDocs_MatchesDescribeContract(t *testing.T) {
	t.Parallel()

	cli := newDocgenContractCLI(t)

	var buf strings.Builder

	err := cli.GenerateDocs(&buf)
	testutil.AssertNoError(t, err)

	docs := buf.String()

	description := cli.Describe()
	commands := append([]CommandDescription{{
		Name:  description.Name,
		Path:  description.Name,
		Flags: description.Flags,
	}}, collectDescriptions(description.Commands)...)

	checked := 0

	for _, command := range commands {
		for _, flag := range command.Flags {
			if flag.Hidden {
				if len(flagRows(docs, flag.Name)) > 0 {
					t.Errorf("hidden flag %q must not appear in generated docs", flag.Name)
				}

				continue
			}

			rows := flagRows(docs, flag.Name)
			if len(rows) == 0 {
				t.Errorf("command %q: flag %q missing from generated docs", command.Path, flag.Name)

				continue
			}

			checked++

			for _, want := range documentedFacts(flag) {
				found := false

				for _, row := range rows {
					if strings.Contains(row, want) {
						found = true

						break
					}
				}

				if !found {
					t.Errorf(
						"command %q: flag %q row does not document %q; row: %q",
						command.Path,
						flag.Name,
						want,
						rows[0],
					)
				}
			}
		}
	}

	if checked == 0 {
		t.Fatal("drift test checked no flags; the fixture stopped exercising the contract")
	}
}

// documentedFacts lists the strings a flag's row must contain for the given
// contract entry. It is written from FlagDescription, not from docgen, so a
// renderer that quietly drops a column is caught.
func documentedFacts(flag FlagDescription) []string {
	facts := []string{"`--" + flag.Name + "`", flag.Type}

	if flag.Shorthand != "" {
		facts = append(facts, "`-"+flag.Shorthand+"`")
	}

	if flag.Default != "" {
		facts = append(facts, "`"+flag.Default+"`")
	}

	if flag.Env != "" {
		facts = append(facts, "`"+flag.Env+"`")
	}

	for _, value := range flag.Values {
		facts = append(facts, "`"+value+"`")
	}

	if flag.Required {
		facts = append(facts, "| yes |")
	}

	if flag.Secret {
		facts = append(facts, "| yes |")
	}

	if flag.Usage != "" {
		facts = append(facts, docCellText(flag.Usage))
	}

	return facts
}

func TestGenerateDocs_InheritedFlagsGetTheirOwnSection(t *testing.T) {
	t.Parallel()

	cli := newDocgenContractCLI(t)

	var buf strings.Builder

	err := cli.GenerateDocs(&buf)
	testutil.AssertNoError(t, err)

	docs := buf.String()

	if !strings.Contains(docs, "### Inherited Flags") {
		t.Fatal("expected an Inherited Flags section")
	}

	deployIndex := strings.Index(docs, "## docapp deploy")
	if deployIndex < 0 {
		t.Fatal("expected a deploy command section")
	}

	deploySection := docs[deployIndex:]
	inheritedIndex := strings.Index(deploySection, "### Inherited Flags")
	if inheritedIndex < 0 {
		t.Fatal("deploy section should list the root's inherited flags")
	}

	if !strings.Contains(deploySection[inheritedIndex:], "`--environment`") {
		t.Error("deploy's inherited section should contain the root's --environment flag")
	}

	localIndex := strings.Index(deploySection, "### Flags")
	if localIndex < 0 || localIndex > inheritedIndex {
		t.Error("a command's own flags must come before its inherited flags")
	}
}

func TestGenerateDocs_EscapesTableBreakingText(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Pattern string `flag:"pattern" help:"a | b" default:"x|y"`
	}

	cli, err := NewCLI("app", "app", cfg{}, WithFang(false))
	testutil.AssertNoError(t, err)

	var buf strings.Builder

	err = cli.GenerateDocs(&buf)
	testutil.AssertNoError(t, err)

	row := flagRows(buf.String(), "pattern")
	testutil.AssertFieldLen(t, row, 1, "pattern row")

	if !strings.Contains(row[0], `x\|y`) {
		t.Errorf("a pipe in a default must be escaped in the table cell: %q", row[0])
	}

	// Eight columns means nine table delimiters once the escaped pipes are
	// neutralized.
	columns := strings.Count(strings.ReplaceAll(row[0], `\|`, ""), "|")
	if columns != 9 {
		t.Errorf("row has %d table delimiters, want 9: %q", columns, row[0])
	}
}
