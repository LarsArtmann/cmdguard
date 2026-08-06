package testutil

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestAssertEqual(t *testing.T) {
	t.Parallel()
	AssertEqual(t, 42, 42)
	AssertEqual(t, "hello", "hello")
}

func TestAssertEqualf(t *testing.T) {
	t.Parallel()
	AssertEqualf(t, "a", "a", "context: %d", 1)
}

func TestAssertNotEqual(t *testing.T) {
	t.Parallel()
	AssertNotEqual(t, 1, 2)
}

func TestAssertNil(t *testing.T) {
	t.Parallel()
	var p *int
	AssertNil(t, p)
}

func TestAssertNotNil(t *testing.T) {
	t.Parallel()
	v := 42
	AssertNotNil(t, &v)
}

func TestAssertErrorIs(t *testing.T) {
	t.Parallel()
	err := errors.New("test error")
	AssertErrorIs(t, err, err)
}

func TestAssertNoError(t *testing.T) {
	t.Parallel()
	AssertNoError(t, nil)
}

func TestAssertErrorContains(t *testing.T) {
	t.Parallel()
	AssertErrorContains(t, errors.New("file not found: config.yaml"), "not found", "config.yaml")
}

func TestAssertBoolTrue(t *testing.T) {
	t.Parallel()
	AssertBoolTrue(t, true, "flag")
}

func TestAssertBoolField(t *testing.T) {
	t.Parallel()
	AssertBoolField(t, true, true, "enabled")
	AssertBoolField(t, false, false, "disabled")
}

func TestAssertBoolFalse(t *testing.T) {
	t.Parallel()
	AssertBoolFalse(t, false, "flag")
}

func TestContainsString(t *testing.T) {
	t.Parallel()
	if !ContainsString([]string{"a", "b", "c"}, "b") {
		t.Error("expected to find 'b'")
	}
	if ContainsString([]string{"a", "b", "c"}, "d") {
		t.Error("expected not to find 'd'")
	}
}

func TestAssertPanics(t *testing.T) {
	t.Parallel()
	AssertPanics(t, func() { panic("test") })
}

func TestExpectPanics(t *testing.T) {
	t.Parallel()
	if !ExpectPanics(t, func() { panic("test") }) {
		t.Error("expected panic to be detected")
	}
}

func TestAssertDoesNotPanic(t *testing.T) {
	t.Parallel()
	AssertDoesNotPanic(t, func() {})
}

func TestAssertLen(t *testing.T) {
	t.Parallel()
	AssertLen(t, []int{1, 2, 3}, 3)
}

func TestAssertStringSlicesEqual(t *testing.T) {
	t.Parallel()
	AssertStringSlicesEqual(t, []string{"a", "b"}, []string{"a", "b"}, "test")
}

func TestAssertContainsString(t *testing.T) {
	t.Parallel()
	AssertContainsString(t, []string{"x", "y"}, "x")
}

func TestAssertNotContainsString(t *testing.T) {
	t.Parallel()
	AssertNotContainsString(t, []string{"x", "y"}, "z")
}

func TestAssertPointerEq(t *testing.T) {
	t.Parallel()
	v := 42
	AssertPointerEq(t, &v, &v)
}

func TestNoOpCobraRun(t *testing.T) {
	t.Parallel()
	fn := NoOpCobraRun()
	if fn == nil {
		t.Fatal("NoOpCobraRun returned nil")
	}
}

func TestNoOpCobraRunE(t *testing.T) {
	t.Parallel()
	fn := NoOpCobraRunE()
	if fn == nil {
		t.Fatal("NoOpCobraRunE returned nil")
	}
}

func TestAssertOutputContains(t *testing.T) {
	t.Parallel()
	AssertOutputContains(t, "hello world", "world")
}

func TestAssertFieldLen(t *testing.T) {
	t.Parallel()
	AssertFieldLen(t, []int{1, 2}, 2, "items")
}

func TestAssertStringerEq(t *testing.T) {
	t.Parallel()
	type str string

	AssertStringerEq(t, testStringer("test"), "test")
}

type testStringer string

func (s testStringer) String() string { return string(s) }

func TestAssertErrorIsf(t *testing.T) {
	t.Parallel()
	err := errors.New("test error")
	AssertErrorIsf(t, err, err, "context: %d", 1)
}

func TestAssertStderrContains(t *testing.T) {
	t.Parallel()
	AssertStderrContains(t, "Error: something failed", "error", "failed")
}

func TestNoOpRunE(t *testing.T) {
	t.Parallel()
	type cfg struct{ Debug bool }
	type flags struct{ Name string }
	err := NoOpRunE[cfg, flags](context.Background(), &cfg{}, flags{})
	AssertNoError(t, err)
}

func TestAssertFieldEq(t *testing.T) {
	t.Parallel()
	AssertFieldEq(t, 42, 42, "count")
}

func TestAssertFieldEqString(t *testing.T) {
	t.Parallel()
	AssertFieldEqString(t, "hello", "hello", "name")
}

func TestAssertFlagRegistered(t *testing.T) {
	t.Parallel()
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("name", "", "name flag")
	AssertFlagRegistered(t, cmd, "name")
}

func TestAssertFlagNotRegistered(t *testing.T) {
	t.Parallel()
	cmd := &cobra.Command{Use: "test"}
	AssertFlagNotRegistered(t, cmd, "nonexistent")
}

func TestAssertStringFieldContains(t *testing.T) {
	t.Parallel()
	AssertStringFieldContains(t, "hello world", "world", "greeting")
}

func TestAssertExpectedError(t *testing.T) {
	t.Parallel()
	AssertExpectedError(t, errors.New("expected"))
}

func TestAssertJSONMarshal(t *testing.T) {
	t.Parallel()
	AssertJSONMarshal(t, []byte(`{"key":"value"}`), `{"key":"value"}`)
}

func TestAssertFieldEqQuote(t *testing.T) {
	t.Parallel()
	AssertFieldEqQuote(t, "value", "value", "field")
}
