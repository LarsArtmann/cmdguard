package v2

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestCommandGroups_BasicGrouping(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithGroup[testConfig]("core", "Core Commands:"),
		WithGroup[testConfig]("utils", "Utilities:"),
		WithColor[testConfig](false),
		WithSilenceUsage[testConfig](),
		WithSilenceErrors[testConfig](),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "serve",
		Short: "Start the server",
		Long:  "Start the server",
		Group: "core",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "migrate",
		Short: "Run migrations",
		Long:  "Run migrations",
		Group: "core",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "version",
		Short: "Print version",
		Long:  "Print version",
		Group: "utils",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	rootCmd := cli.RootCommand()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	_ = cli.Execute(context.Background())

	helpOutput := buf.String()

	if !strings.Contains(helpOutput, "Core Commands:") {
		t.Error("help output should contain 'Core Commands:' group title")
	}

	if !strings.Contains(helpOutput, "Utilities:") {
		t.Error("help output should contain 'Utilities:' group title")
	}

	if !strings.Contains(helpOutput, "serve") {
		t.Error("help output should contain 'serve' command")
	}

	if !strings.Contains(helpOutput, "migrate") {
		t.Error("help output should contain 'migrate' command")
	}

	if !strings.Contains(helpOutput, "version") {
		t.Error("help output should contain 'version' command")
	}
}

func TestCommandGroups_NoGroup(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithGroup[testConfig]("core", "Core Commands:"),
		WithColor[testConfig](false),
		WithSilenceUsage[testConfig](),
		WithSilenceErrors[testConfig](),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "ungrouped",
		Short: "No group assigned",
		Long:  "No group assigned",
		Group: "",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"ungrouped"})
	testutil.AssertNoError(t, err)
}

func TestCommandGroups_CommandExecutionStillWorks(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	executed := false

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithGroup[testConfig]("main", "Main Commands:"),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "run",
		Short: "Run command",
		Long:  "Run command",
		Group: "main",
		RunE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			executed = true

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	if !executed {
		t.Error("command should still execute when assigned to a group")
	}
}

func TestCommandGroups_SubcommandsInheritFromParent(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	childExecuted := false

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithGroup[testConfig]("core", "Core Commands:"),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use:   "db",
		Short: "Database",
		Long:  "Database operations",
		Group: "core",
		Commands: []Command[testConfig, NoFlags]{
			{
				Use:   "migrate",
				Short: "Run migrations",
				RunE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
					childExecuted = true

					return nil
				},
			},
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"db", "migrate"})
	testutil.AssertNoError(t, err)

	if !childExecuted {
		t.Error("child command should execute")
	}
}

func TestWithGroup_RegistersMultipleGroups(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithGroup[testConfig]("alpha", "Alpha Group:"),
		WithGroup[testConfig]("beta", "Beta Group:"),
		WithGroup[testConfig]("gamma", "Gamma Group:"),
		WithColor[testConfig](false),
		WithSilenceUsage[testConfig](),
		WithSilenceErrors[testConfig](),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use: "a", Short: "A", Long: "A",
		Group: "alpha",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use: "b", Short: "B", Long: "B",
		Group: "beta",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		Use: "c", Short: "C", Long: "C",
		Group: "gamma",
		RunE:  func(_ context.Context, _ *testConfig, _ NoFlags) error { return nil },
	})
	testutil.AssertNoError(t, err)

	rootCmd := cli.RootCommand()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	_ = cli.Execute(context.Background())

	helpOutput := buf.String()

	for _, title := range []string{"Alpha Group:", "Beta Group:", "Gamma Group:"} {
		if !strings.Contains(helpOutput, title) {
			t.Errorf("help output should contain %q", title)
		}
	}
}
