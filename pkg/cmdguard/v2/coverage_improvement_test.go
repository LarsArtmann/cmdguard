package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestGet_GenericTypedAccessor(t *testing.T) {
	t.Parallel()

	t.Run("retrieves typed string value", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		root.SetValue("greeting", "hello")
		ctx := WithBranchingFlowContext(context.Background(), root)

		val, ok := Get[string](ctx, "greeting")
		if !ok {
			t.Fatal("expected to find value")
		}

		testutil.AssertFieldEqString(t, val, "hello", "value")
	})

	t.Run("retrieves typed int value", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		root.SetValue("count", 42)
		ctx := WithBranchingFlowContext(context.Background(), root)

		val, ok := Get[int](ctx, "count")
		if !ok {
			t.Fatal("expected to find value")
		}

		testutil.AssertEqual(t, val, 42)
	})

	t.Run("returns zero and false for missing key", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		ctx := WithBranchingFlowContext(context.Background(), root)

		val, ok := Get[string](ctx, "nonexistent")
		if ok {
			t.Error("expected ok=false for missing key")
		}

		testutil.AssertFieldEqString(t, val, "", "zero value")
	})

	t.Run("returns zero and false for wrong type", func(t *testing.T) {
		t.Parallel()

		root := NewBranchingFlowContext(context.Background())
		root.SetValue("value", "string-value")
		ctx := WithBranchingFlowContext(context.Background(), root)

		val, ok := Get[int](ctx, "value")
		if ok {
			t.Error("expected ok=false for wrong type")
		}

		testutil.AssertEqual(t, val, 0)
	})

	t.Run("returns zero and false for nil context", func(t *testing.T) {
		t.Parallel()

		val, ok := Get[string](nil, "key") //nolint:staticcheck // testing nil context handling
		if ok {
			t.Error("expected ok=false for nil context")
		}

		testutil.AssertFieldEqString(t, val, "", "zero value")
	})

	t.Run("returns zero and false for context without flow context", func(t *testing.T) {
		t.Parallel()

		val, ok := Get[string](context.Background(), "key")
		if ok {
			t.Error("expected ok=false for context without flow context")
		}

		testutil.AssertFieldEqString(t, val, "", "zero value")
	})
}

func TestCLI_HealthCheckResults(t *testing.T) {
	t.Parallel()

	t.Run("returns map from scope health checks", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cli, err := NewCLI[cfg]("test", "Test", cfg{})
		testutil.AssertNoError(t, err)

		results := cli.HealthCheckResults()

		scopeResults := cli.Scope().HealthCheckResults()
		testutil.AssertEqual(t, len(results), len(scopeResults))
	})

	t.Run("returns per-service results", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cli, err := NewCLI[cfg]("test", "Test", cfg{})
		testutil.AssertNoError(t, err)

		mustProvideValue(t, cli.Scope(), "healthy-service")

		results := cli.HealthCheckResults()
		if len(results) == 0 {
			t.Error("expected at least one result")
		}
	})
}

func TestWithArgs_SetsCustomValidator(t *testing.T) {
	t.Parallel()

	t.Run("sets args validator on command", func(t *testing.T) {
		t.Parallel()

		called := false
		validator := func(_ *cobra.Command, _ []string) error {
			called = true

			return nil
		}

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithArgs[testConfig, NoFlags](validator)(&cmd)

		if cmd.args == nil {
			t.Fatal("expected args to be set")
		}

		if err := cmd.args(&cobra.Command{}, []string{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !called {
			t.Error("expected validator to be called")
		}
	})
}

func TestWithCompletion_SetsCompletionFunc(t *testing.T) {
	t.Parallel()

	t.Run("sets completion function on command", func(t *testing.T) {
		t.Parallel()

		fn := func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"apple", "banana"}, cobra.ShellCompDirectiveNoFileComp
		}

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithCompletion[testConfig, NoFlags](fn)(&cmd)

		if cmd.completionFn == nil {
			t.Fatal("expected completionFn to be set")
		}

		results, directive := cmd.completionFn(&cobra.Command{}, nil, "")
		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}

		if directive != cobra.ShellCompDirectiveNoFileComp {
			t.Errorf("expected ShellCompDirectiveNoFileComp, got %v", directive)
		}
	})
}

func TestWithValidArgs_SetsStaticArgs(t *testing.T) {
	t.Parallel()

	t.Run("sets valid args on command", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithValidArgs[testConfig, NoFlags]("one", "two", "three")(&cmd)

		if len(cmd.validArgs) != 3 {
			t.Fatalf("expected 3 valid args, got %d", len(cmd.validArgs))
		}

		testutil.AssertFieldEqString(t, cmd.validArgs[0], "one", "validArgs[0]")
		testutil.AssertFieldEqString(t, cmd.validArgs[1], "two", "validArgs[1]")
		testutil.AssertFieldEqString(t, cmd.validArgs[2], "three", "validArgs[2]")
	})

	t.Run("sets empty valid args", func(t *testing.T) {
		t.Parallel()

		cmd := Command[testConfig, NoFlags]{use: "test"}
		WithValidArgs[testConfig, NoFlags]()(&cmd)

		if cmd.validArgs != nil {
			t.Errorf("expected nil validArgs, got %v", cmd.validArgs)
		}
	})
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptString_DelegatesToRunner(t *testing.T) {
	setFakePromptRunner(t, &fakePromptRunner{
		stringResults: map[string]string{
			"Name?": "Alice",
		},
	})

	result, err := PromptString("Name?", "")
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, result, "Alice", "result")
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptString_WrapsError(t *testing.T) {
	promptErr := errors.New("terminal error")
	setFakePromptRunner(t, &fakePromptRunner{
		stringErrors: map[string]error{
			"Name?": promptErr,
		},
	})

	_, err := PromptString("Name?", "")
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, promptErr) {
		t.Errorf("expected error to wrap promptErr, got: %v", err)
	}
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptSelect_DelegatesToRunner(t *testing.T) {
	setFakePromptRunner(t, &fakePromptRunner{
		selectResults: map[string]string{
			"Color?": "blue",
		},
	})

	result, err := PromptSelect("Color?", []string{"red", "green", "blue"})
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, result, "blue", "result")
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptSelect_WrapsError(t *testing.T) {
	promptErr := errors.New("select error")
	setFakePromptRunner(t, &fakePromptRunner{
		selectErrors: map[string]error{
			"Color?": promptErr,
		},
	})

	_, err := PromptSelect("Color?", []string{"red"})
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, promptErr) {
		t.Errorf("expected error to wrap promptErr, got: %v", err)
	}
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptConfirm_DelegatesToRunner(t *testing.T) {
	setFakePromptRunner(t, &fakePromptRunner{
		confirmResults: map[string]bool{
			"Continue?": true,
		},
	})

	result, err := PromptConfirm("Continue?")
	testutil.AssertNoError(t, err)

	if !result {
		t.Error("expected true")
	}
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestPromptConfirm_WrapsError(t *testing.T) {
	promptErr := errors.New("confirm error")
	setFakePromptRunner(t, &fakePromptRunner{
		confirmErrors: map[string]error{
			"Continue?": promptErr,
		},
	})

	_, err := PromptConfirm("Continue?")
	testutil.AssertExpectedError(t, err)

	if !errors.Is(err, promptErr) {
		t.Errorf("expected error to wrap promptErr, got: %v", err)
	}
}
