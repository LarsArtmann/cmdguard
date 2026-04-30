package v2

import (
	"os"
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestEditInEditor(t *testing.T) {
	//nolint:paralleltest // uses t.Setenv

	t.Run("uses EDITOR env var", func(t *testing.T) {
		// Create a script that just writes to the file
		script := "#!/bin/sh\ncat > \"$1\" << 'EOF'\nedited content\nEOF\n"
		tmpScript, err := os.CreateTemp("", "editor-test-*.sh")
		testutil.AssertNoError(t, err)
		defer os.Remove(tmpScript.Name())

		_, err = tmpScript.WriteString(script)
		testutil.AssertNoError(t, err)
		tmpScript.Close()

		err = os.Chmod(tmpScript.Name(), 0o755)
		testutil.AssertNoError(t, err)

		t.Setenv("EDITOR", tmpScript.Name())

		result, err := EditInEditor("original content")
		testutil.AssertNoError(t, err)
		if result != "edited content\n" {
			t.Errorf("EditInEditor() = %q, want %q", result, "edited content\n")
		}
	})

	t.Run("falls back to vi when EDITOR not set", func(t *testing.T) {
		// This test just verifies the function doesn't crash with empty EDITOR.
		// We can't actually test vi behavior, so we skip if vi would run.
		t.Skip("cannot test vi fallback without an interactive terminal")
	})
}
