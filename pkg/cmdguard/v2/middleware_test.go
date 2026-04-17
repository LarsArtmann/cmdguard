package v2

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestMiddleware_BasicChaining(t *testing.T) {
	t.Parallel()

	type testConfig struct {
		Name string `default:"test" flag:"name" help:"test"`
	}

	var callOrder []string

	mw1 := func(_ context.Context, _ *testConfig, _ CommandInfo, next func() error) error {
		callOrder = append(callOrder, "mw1-before")
		err := next()

		callOrder = append(callOrder, "mw1-after")

		return err
	}

	mw2 := func(_ context.Context, _ *testConfig, _ CommandInfo, next func() error) error {
		callOrder = append(callOrder, "mw2-before")
		err := next()

		callOrder = append(callOrder, "mw2-after")

		return err
	}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(mw1, mw2),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "run",
		short: "Run command",
		long:  "Run command",
		runE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			callOrder = append(callOrder, "handler")

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"run"})
	testutil.AssertNoError(t, err)

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(callOrder) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(callOrder), callOrder)
	}

	for i, v := range expected {
		if callOrder[i] != v {
			t.Errorf("call[%d]: expected %q, got %q", i, v, callOrder[i])
		}
	}
}

func TestMiddleware_ErrorPropagation(t *testing.T) {
	t.Parallel()

	type testConfig struct {
		Name string `default:"test" flag:"name" help:"test"`
	}

	handlerErr := errors.New("handler failed")

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "fail",
		short: "Fail command",
		long:  "Fail command",
		runE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			return handlerErr
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"fail"})
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, handlerErr)
}

func TestMiddleware_ShortCircuit(t *testing.T) {
	t.Parallel()

	type testConfig struct {
		Name string `default:"test" flag:"name" help:"test"`
	}

	shortCircuitErr := errors.New("blocked")
	handlerCalled := false

	blockingMiddleware := func(_ context.Context, _ *testConfig, _ CommandInfo, _ func() error) error {
		return shortCircuitErr
	}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(blockingMiddleware),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "blocked",
		short: "Blocked command",
		long:  "Blocked command",
		runE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			handlerCalled = true

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"blocked"})
	testutil.AssertExpectedError(t, err)
	testutil.AssertErrorIs(t, err, shortCircuitErr)

	if handlerCalled {
		t.Error("handler should not be called when middleware short-circuits")
	}
}

func TestMiddleware_CommandInfo(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	var capturedInfo CommandInfo

	captureMiddleware := func(_ context.Context, _ *testConfig, info CommandInfo, next func() error) error {
		capturedInfo = info

		return next()
	}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(captureMiddleware),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "mycommand",
		short: "My command",
		long:  "My command",
		runE:  noOpRunE[testConfig, NoFlags],
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"mycommand"})
	testutil.AssertNoError(t, err)

	testutil.AssertFieldEqString(t, capturedInfo.Name, "mycommand", "command name")

	if !capturedInfo.HasRunE {
		t.Error("HasRunE should be true")
	}
}

func TestTimingMiddleware(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	var (
		capturedName     string
		capturedDuration time.Duration
	)

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(TimingMiddleware[testConfig](func(name string, d time.Duration) {
			capturedName = name
			capturedDuration = d
		})),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "timed",
		short: "Timed command",
		long:  "Timed command",
		runE:  noOpRunE[testConfig, NoFlags],
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"timed"})
	testutil.AssertNoError(t, err)

	testutil.AssertFieldEqString(t, capturedName, "timed", "command name")

	if capturedDuration < 0 {
		t.Error("duration should be non-negative")
	}

	if capturedDuration > time.Second {
		t.Errorf("duration %v seems unreasonably large for a no-op command", capturedDuration)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(RecoveryMiddleware[testConfig]()),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "panic",
		short: "Panic command",
		long:  "Panic command",
		runE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			panic("something went wrong")
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"panic"})
	testutil.AssertExpectedError(t, err)

	if !strings.Contains(err.Error(), "panic") {
		t.Errorf("error should mention panic, got: %v", err)
	}
}

func TestMiddleware_NoMiddleware(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	handlerCalled := false

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "plain",
		short: "Plain command",
		long:  "Plain command",
		runE: func(_ context.Context, _ *testConfig, _ NoFlags) error {
			handlerCalled = true

			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"plain"})
	testutil.AssertNoError(t, err)

	if !handlerCalled {
		t.Error("handler should be called without middleware")
	}
}

func TestMiddleware_Subcommands(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	var commandNames []string

	trackingMiddleware := func(_ context.Context, _ *testConfig, info CommandInfo, next func() error) error {
		commandNames = append(commandNames, info.Name)

		return next()
	}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(trackingMiddleware),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		use:   "parent",
		short: "Parent",
		long:  "Parent command with subcommands",
		commands: []Command[testConfig, NoFlags]{
			{
				use:   "child",
				short: "Child",
				runE:  noOpRunE[testConfig, NoFlags],
			},
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"parent", "child"})
	testutil.AssertNoError(t, err)

	if len(commandNames) == 0 {
		t.Error("middleware should track subcommand execution")
	}
}

func TestChainMiddleware(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	var order []string

	mw := func(name string) Middleware[testConfig] {
		return func(_ context.Context, _ *testConfig, _ CommandInfo, next func() error) error {
			order = append(order, name+"-before")

			err := next()

			order = append(order, name+"-after")

			return err
		}
	}

	chain := buildChain(
		context.Background(), &testConfig{}, CommandInfo{Name: "test"},
		[]Middleware[testConfig]{mw("a"), mw("b"), mw("c")},
		func() error {
			order = append(order, "handler")

			return nil
		},
	)

	err := chain()
	testutil.AssertNoError(t, err)

	expected := []string{
		"a-before",
		"b-before",
		"c-before",
		"handler",
		"c-after",
		"b-after",
		"a-after",
	}
	if len(order) != len(expected) {
		t.Fatalf("expected %d entries, got %d: %v", len(expected), len(order), order)
	}

	for i, v := range expected {
		if order[i] != v {
			t.Errorf("entry[%d]: expected %q, got %q", i, v, order[i])
		}
	}
}

func TestChainMiddleware_Empty(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	called := false

	chain := buildChain(
		context.Background(),
		&testConfig{},
		CommandInfo{},
		[]Middleware[testConfig]{},
		func() error {
			called = true

			return nil
		},
	)

	err := chain()
	testutil.AssertNoError(t, err)

	if !called {
		t.Error("handler should be called directly with empty middleware")
	}
}

func TestMiddleware_WithFlags(t *testing.T) {
	t.Parallel()

	type testConfig struct{}

	type greetFlags struct {
		Name string `default:"World" flag:"name" help:"Name" short:"n"`
	}

	var capturedName string

	inspectMiddleware := func(_ context.Context, _ *testConfig, info CommandInfo, next func() error) error {
		capturedName = info.Name

		return next()
	}

	cli, err := NewCLI[testConfig]("test", "Test CLI", testConfig{},
		WithMiddleware(inspectMiddleware),
		WithColor[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, *greetFlags]{
		use:   "greet",
		short: "Greet",
		long:  "Greet someone",
		flags: &greetFlags{},
		runE: func(_ context.Context, _ *testConfig, _ *greetFlags) error {
			return nil
		},
	})
	testutil.AssertNoError(t, err)

	err = cli.ExecuteWithArgs(context.Background(), []string{"greet", "-n", "Middleware"})
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, capturedName, "greet", "command name in middleware")
}
