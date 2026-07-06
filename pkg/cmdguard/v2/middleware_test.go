package v2

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func makeMiddlewareCommand[T any](name, desc string, handlerCalled *bool) Command[T, NoFlags] {
	return Command[T, NoFlags]{
		spec: commandSpec{use: name, short: desc + " command", long: desc + " command"},
		runE: func(_ context.Context, _ *T, _ NoFlags) error {
			*handlerCalled = true

			return nil
		},
	}
}

// addNoOpCommand builds a NoFlags command with use/short/long and no-op runE, then registers it.
func addNoOpCommand[T any](t *testing.T, cli *CLI[T], use, short string) {
	t.Helper()

	err := AddCommand(cli, Command[T, NoFlags]{
		spec: commandSpec{use: use, short: short, long: short},
		runE:  noOpRunE[T, NoFlags],
	})
	testutil.AssertNoError(t, err)
}

// beforeAfterMiddleware returns a Middleware[T] that records "<name>-before" before
// invoking next() and "<name>-after" after, appending both to order. Used to verify
// middleware chain ordering in tests.
func beforeAfterMiddleware[T any](order *[]string, name string) Middleware[T] {
	return func(_ context.Context, _ *T, _ CommandInfo, next func() error) error {
		*order = append(*order, name+"-before")

		err := next()

		*order = append(*order, name+"-after")

		return err
	}
}

// captureInfoMiddleware returns a Middleware[T] that captures the incoming CommandInfo
// into *captured before delegating to next(). Used to assert the info passed to middleware.
func captureInfoMiddleware[T any](captured *CommandInfo) Middleware[T] {
	return func(_ context.Context, _ *T, info CommandInfo, next func() error) error {
		*captured = info

		return next()
	}
}

// captureNameMiddleware returns a Middleware[T] that appends info.Name to *names
// before delegating to next(). Used to verify which commands the chain visited.
func captureNameMiddleware[T any](names *[]string) Middleware[T] {
	return func(_ context.Context, _ *T, info CommandInfo, next func() error) error {
		*names = append(*names, info.Name)

		return next()
	}
}

func TestMiddleware_BasicChaining(t *testing.T) {
	t.Parallel()

	type testConfig struct {
		Name string `default:"test" flag:"name" help:"test"`
	}

	var callOrder []string

	mw1 := beforeAfterMiddleware[testConfig](&callOrder, "mw1")
	mw2 := beforeAfterMiddleware[testConfig](&callOrder, "mw2")

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(mw1, mw2),
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		spec: commandSpec{use: "run", short: "Run command", long: "Run command"},
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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		spec: commandSpec{use: "fail", short: "Fail command", long: "Fail command"},
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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(blockingMiddleware),
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	cmd := makeMiddlewareCommand[testConfig]("blocked", "Blocked", &handlerCalled)
	err = AddCommand(cli, cmd)
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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(captureInfoMiddleware[testConfig](&capturedInfo)),
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	addNoOpCommand(t, cli, "mycommand", "My command")

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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(TimingMiddleware[testConfig](func(name string, d time.Duration, err error) {
			capturedName = name
			capturedDuration = d
		})),
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	addNoOpCommand(t, cli, "timed", "Timed command")

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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(RecoveryMiddleware[testConfig]()),
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	err = AddCommand(cli, Command[testConfig, NoFlags]{
		spec: commandSpec{use: "panic", short: "Panic command", long: "Panic command"},
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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithFang[testConfig](false),
	)
	testutil.AssertNoError(t, err)

	handlerCalled := false

	cmd := makeMiddlewareCommand[testConfig]("plain", "Plain", &handlerCalled)
	err = AddCommand(cli, cmd)
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

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(captureNameMiddleware[testConfig](&commandNames)),
		WithFang[testConfig](false),
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

	chain := buildChain(
		context.Background(), &testConfig{}, CommandInfo{Name: "test"},
		[]Middleware[testConfig]{
			beforeAfterMiddleware[testConfig](&order, "a"),
			beforeAfterMiddleware[testConfig](&order, "b"),
			beforeAfterMiddleware[testConfig](&order, "c"),
		},
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

	testutil.AssertStringSlicesEqual(t, order, expected, "entry")
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

	var capturedNames []string

	cli, err := NewCLI[testConfig](
		"test", "Test CLI", testConfig{},
		WithMiddleware(captureNameMiddleware[testConfig](&capturedNames)),
		WithFang[testConfig](false),
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
	testutil.AssertFieldEqString(t, capturedNames[0], "greet", "command name in middleware")
}
