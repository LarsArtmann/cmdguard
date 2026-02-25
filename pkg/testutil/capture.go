// Package testutil provides testing utilities for cmdguard examples and tests.
package testutil

import (
	"io"
	"os"
)

// CaptureOutput captures stdout output during the execution of function f.
// It temporarily redirects os.Stdout to a pipe, executes f, then restores
// the original stdout and returns the captured output as a string.
func CaptureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()

	os.Stdout = old

	out, _ := io.ReadAll(r)

	return string(out)
}
