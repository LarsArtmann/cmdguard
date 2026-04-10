package cmdguard

import (
	"testing"

	"github.com/spf13/cobra"
)

// assertPanics runs fn and returns true if it panicked.
func assertPanics(t *testing.T, fn func()) bool {
	t.Helper()

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

func TestGuardedCommand_AddCommand(t *testing.T) {
	t.Parallel()
	t.Run("accepts valid command with Run", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(*cobra.Command, []string) {},
		}

		if assertPanics(t, func() { g.AddCommand(cmd) }) {
			t.Error("AddCommand should not panic for valid command")
		}
	})

	t.Run("accepts valid command with RunE", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		cmd := newTestCommand("sub")

		if assertPanics(t, func() { g.AddCommand(cmd) }) {
			t.Error("AddCommand should not panic for valid command")
		}
	})

	t.Run("panics on command without handler", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		cmd := &cobra.Command{
			Use: "invalid",
		}

		if !assertPanics(t, func() { g.AddCommand(cmd) }) {
			t.Error("AddCommand should panic for command without handler")
		}
	})

	t.Run("panics on command without name", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		cmd := &cobra.Command{}

		if !assertPanics(t, func() { g.AddCommand(cmd) }) {
			t.Error("AddCommand should panic for command without name")
		}
	})

	t.Run("panics after Execute called", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")
		g.validated = true // Simulate post-execute state

		cmd := &cobra.Command{
			Use: "sub",
			Run: func(*cobra.Command, []string) {},
		}

		if !assertPanics(t, func() { g.AddCommand(cmd) }) {
			t.Error("AddCommand should panic after Execute called")
		}
	})
}

func TestGuardedCommand_AddSubcommand(t *testing.T) {
	t.Parallel()
	t.Run("adds subcommand to parent", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(*cobra.Command, []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "child",
			Run: func(*cobra.Command, []string) {},
		}

		if assertPanics(t, func() { g.AddSubcommand(parent, child) }) {
			t.Error("AddSubcommand should not panic for valid child")
		}

		found := false

		for _, c := range parent.Commands() {
			if c.Name() == "child" {
				found = true

				break
			}
		}

		if !found {
			t.Error("child command should be added to parent")
		}
	})

	t.Run("panics on invalid child", func(t *testing.T) {
		t.Parallel()

		g := New("testapp", "Test")

		parent := &cobra.Command{
			Use: "parent",
			Run: func(*cobra.Command, []string) {},
		}
		g.AddCommand(parent)

		child := &cobra.Command{
			Use: "invalid-child",
		}

		if !assertPanics(t, func() { g.AddSubcommand(parent, child) }) {
			t.Error("AddSubcommand should panic for invalid child")
		}
	})
}
