// Package testutil provides shared testing utilities for cmdguard tests.
package testutil

import (
	"slices"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// NoOpCobraRun returns a no-op Run function for cobra.Command.
// This reduces duplication in tests that need a valid command with a handler.
func NoOpCobraRun() func(*cobra.Command, []string) {
	return func(*cobra.Command, []string) {}
}

// NoOpCobraRunE returns a no-op RunE function for cobra.Command.
func NoOpCobraRunE() func(*cobra.Command, []string) error {
	return func(*cobra.Command, []string) error { return nil }
}

// doPanicTest runs fn and returns true if it panicked.
// This is the shared implementation for panic detection.
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
// Unlike AssertPanics, this does NOT fail the test if no panic occurs.
func ExpectPanics(t *testing.T, fn func()) bool {
	t.Helper()

	return doPanicTest(fn)
}

// AssertErrorContains fails the test if the error does not contain all the given substrings.
func AssertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()

	errMsg := err.Error()
	for _, s := range substrings {
		if !strings.Contains(errMsg, s) {
			t.Errorf("error should contain %q, got %q", s, errMsg)
		}
	}
}

// AssertStderrContains fails the test if the stderr does not contain all given substrings (case-insensitive).
func AssertStderrContains(t *testing.T, stderr string, substrings ...string) {
	t.Helper()

	for _, s := range substrings {
		if !strings.Contains(strings.ToLower(stderr), strings.ToLower(s)) {
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
