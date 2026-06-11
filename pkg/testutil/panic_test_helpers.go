// Package testutil provides shared testing utilities for cmdguard tests.
package testutil

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// AssertEqual fails the test if got != want.
func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// AssertEqualf fails the test if got != want with formatted message.
func AssertEqualf[T comparable](t *testing.T, got, want T, format string, args ...any) {
	t.Helper()

	if got != want {
		t.Errorf("%s: got %v, want %v", fmt.Sprintf(format, args...), got, want)
	}
}

// AssertNotEqual fails the test if got == want.
func AssertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got == want {
		t.Errorf("got %v, expected not to equal %v", got, want)
	}
}

// AssertNil fails the test if v is not nil.
func AssertNil[T any](t *testing.T, v *T) {
	t.Helper()

	if v != nil {
		t.Errorf("expected nil, got %v", *v)
	}
}

// AssertNotNil fails the test if v is nil.
func AssertNotNil[T any](t *testing.T, v *T) {
	t.Helper()

	if v == nil {
		t.Error("expected non-nil")
	}
}

// AssertErrorIs fails the test if err does not match target via errors.Is.
func AssertErrorIs(t *testing.T, err, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Errorf("error does not match %v", target)
	}
}

// AssertErrorIsf fails the test if err does not match target with formatted message.
func AssertErrorIsf(t *testing.T, err, target error, format string, args ...any) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Errorf("%s: error does not match %v", fmt.Sprintf(format, args...), target)
	}
}

// AssertErrorContains fails the test if err does not contain all given substrings.
func AssertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	errMsg := err.Error()
	for _, s := range substrings {
		if !strings.Contains(errMsg, s) {
			t.Errorf("error should contain %q, got %q", s, errMsg)
		}
	}
}

// AssertStderrContains fails the test if stderr does not contain all given substrings (case-insensitive).
func AssertStderrContains(t *testing.T, stderr string, substrings ...string) {
	t.Helper()

	stderrLower := strings.ToLower(stderr)
	for _, s := range substrings {
		if !strings.Contains(stderrLower, strings.ToLower(s)) {
			t.Errorf("stderr should contain %q, got %q", s, stderr)
		}
	}
}

// ContainsString checks if a slice contains a specific string.
func ContainsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

// StringSliceContains is an alias for ContainsString for better readability.
func StringSliceContains(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

// NoOpCobraRun returns a no-op Run function for cobra.Command.
func NoOpCobraRun() func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {}
}

// NoOpCobraRunE returns a no-op RunE function for cobra.Command.
func NoOpCobraRunE() func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error { return nil }
}

// NoOpRunE returns a no-op RunE handler matching cmdguard's generic handler
// signature. Useful in tests and benchmarks that need a benign command body.
func NoOpRunE[T, F any](_ context.Context, _ *T, _ F) error {
	return nil
}

// doPanicTest runs fn and returns true if it panicked.
func doPanicTest(fn func()) bool {
	didPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		fn()
	}()

	return didPanic
}

// AssertPanics runs fn and fails the test if it doesn't panic.
func AssertPanics(t *testing.T, fn func()) {
	t.Helper()

	if !doPanicTest(fn) {
		t.Error("expected panic, got none")
	}
}

// ExpectPanics runs fn and returns true if it panicked.
func ExpectPanics(t *testing.T, fn func()) bool {
	t.Helper()

	return doPanicTest(fn)
}

// AssertDoesNotPanic runs fn and fails the test if it panics.
func AssertDoesNotPanic(t *testing.T, fn func()) {
	t.Helper()

	if doPanicTest(fn) {
		t.Error("expected no panic, but it panicked")
	}
}

// assertFieldEq is the internal generic assertion for field equality.
func assertFieldEq[T comparable](t *testing.T, field, expected T, msg string) {
	t.Helper()

	if field != expected {
		t.Errorf(msg, expected, field)
	}
}

// AssertFieldEq fails the test if field != expected.
func AssertFieldEq[T comparable](t *testing.T, field, expected T, fieldName string) {
	assertFieldEq(t, field, expected, "expected "+fieldName+" %v, got %v")
}

// AssertFieldEqString fails the test if string field != expected (uses %q for formatting).
func AssertFieldEqString(t *testing.T, field, expected, fieldName string) {
	assertFieldEq(t, field, expected, "expected "+fieldName+" %q, got %q")
}

// AssertBoolField fails the test if bool field != expected.
func AssertBoolField(t *testing.T, field, expected bool, fieldName string) {
	t.Helper()

	if field != expected {
		t.Errorf("expected %s to be %v, got %v", fieldName, expected, field)
	}
}

// AssertBoolTrue fails the test if field is not true.
func AssertBoolTrue(t *testing.T, field bool, fieldName string) {
	t.Helper()

	if !field {
		t.Error("expected " + fieldName + " to be true")
	}
}

// AssertBoolFalse fails the test if field is not false.
func AssertBoolFalse(t *testing.T, field bool, fieldName string) {
	t.Helper()

	if field {
		t.Error("expected " + fieldName + " to be false")
	}
}

// AssertStringSlicesEqual fails the test if got and want have different lengths
// or differ at any index. The context label appears in the length-mismatch message.
func AssertStringSlicesEqual(t *testing.T, got, want []string, context string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("%s: got %v, want %v", context, got, want)

		return
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("%s[%d]: got %q, want %q", context, i, got[i], want[i])
		}
	}
}

// assertFlagRegistered is the internal helper for flag registration checks.
func assertFlagRegistered(t *testing.T, cmd *cobra.Command, flagName string, shouldExist bool) {
	t.Helper()

	exists := cmd.Flags().Lookup(flagName) != nil
	if exists != shouldExist {
		if shouldExist {
			t.Errorf("expected %q flag to be registered", flagName)
		} else {
			t.Errorf("expected %q flag to not be registered", flagName)
		}
	}
}

// AssertFlagRegistered fails the test if the flag is not registered.
func AssertFlagRegistered(t *testing.T, cmd *cobra.Command, flagName string) {
	assertFlagRegistered(t, cmd, flagName, true)
}

// AssertFlagNotRegistered fails the test if the flag is registered.
func AssertFlagNotRegistered(t *testing.T, cmd *cobra.Command, flagName string) {
	assertFlagRegistered(t, cmd, flagName, false)
}

// AssertStringFieldContains fails the test if string field does not contain substring.
func AssertStringFieldContains(t *testing.T, field, substr, fieldName string) {
	t.Helper()

	if !strings.Contains(field, substr) {
		t.Errorf("%s should contain %q, got %q", fieldName, substr, field)
	}
}

// AssertOutputContains fails the test if output does not contain substring.
// Use for captured CLI command output buffers.
func AssertOutputContains(t *testing.T, output, substr string) {
	t.Helper()

	if !strings.Contains(output, substr) {
		t.Errorf("output should contain %q, got: %s", substr, output)
	}
}

// assertFieldLen is the internal helper for length assertions.
func assertFieldLen[T any](t *testing.T, slice []T, expected int, msg string) {
	t.Helper()

	if len(slice) != expected {
		t.Errorf(msg, len(slice), expected)
	}
}

// AssertFieldLen fails the test if slice length != expected.
func AssertFieldLen[T any](t *testing.T, slice []T, expected int, fieldName string) {
	assertFieldLen(t, slice, expected, "len("+fieldName+") = %d, want %d")
}

// AssertLen is an alias for AssertFieldLen using the field name as description.
func AssertLen[T any](t *testing.T, slice []T, expected int) {
	assertFieldLen(t, slice, expected, "len() = %d, want %d")
}

// AssertNoError fails the test if err != nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

// AssertExpectedError fails the test if err == nil.
func AssertExpectedError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// assertContainsString is the internal helper for slice containment checks.
func assertContainsString(t *testing.T, slice []string, s string, shouldContain bool) {
	t.Helper()

	contains := slices.Contains(slice, s)
	if contains != shouldContain {
		if shouldContain {
			t.Errorf("should contain %q, got %v", s, slice)
		} else {
			t.Errorf("should not contain %q, got %v", s, slice)
		}
	}
}

// AssertContainsString fails the test if slice does not contain s.
func AssertContainsString(t *testing.T, slice []string, s string) {
	assertContainsString(t, slice, s, true)
}

// AssertNotContainsString fails the test if slice contains s.
func AssertNotContainsString(t *testing.T, slice []string, s string) {
	assertContainsString(t, slice, s, false)
}

// AssertPointerEq fails the test if pointer addresses are not equal.
func AssertPointerEq[T any](t *testing.T, got, want *T) {
	t.Helper()

	if got != want {
		t.Errorf("got %p, want %p", got, want)
	}
}

// AssertJSONMarshal fails the test if json marshal result != expected.
func AssertJSONMarshal(t *testing.T, got []byte, expected string) {
	t.Helper()

	if string(got) != expected {
		t.Errorf("json.Marshal() = %q, want %q", string(got), expected)
	}
}

// AssertFieldEqQuote fails the test if got != want for string field (uses %q formatting).
func AssertFieldEqQuote(t *testing.T, got, want, fieldName string) {
	assertFieldEq(t, got, want, fieldName+" = %q, want %q")
}

// AssertStringerEq fails the test if got.String() != want.
func AssertStringerEq[T fmt.Stringer](t *testing.T, got T, want string) {
	t.Helper()

	if got.String() != want {
		t.Errorf("String() = %q, want %q", got.String(), want)
	}
}
