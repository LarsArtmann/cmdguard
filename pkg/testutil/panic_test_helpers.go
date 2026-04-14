// Package testutil provides shared testing utilities for cmdguard tests.
package testutil

import (
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

// doPanicTest runs fn and returns true if it panicked.
func doPanicTest(fn func()) (didPanic bool) {
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
