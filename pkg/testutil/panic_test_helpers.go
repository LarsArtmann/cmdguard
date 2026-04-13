// Package testutil provides shared testing utilities for cmdguard tests.
package testutil

import (
	"testing"
)

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
