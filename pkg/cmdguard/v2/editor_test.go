package v2

import (
	"context"
	"os"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestEditInEditor(t *testing.T) {
	t.Run("uses EDITOR env var", func(t *testing.T) {
		script := "#!/bin/sh\ncat > \"$1\" << 'EOF'\nedited content\nEOF\n"
		tmpScript, err := os.CreateTemp(t.TempDir(), "editor-test-*.sh")
		testutil.AssertNoError(t, err)

		defer os.Remove(tmpScript.Name())

		_, err = tmpScript.WriteString(script)
		testutil.AssertNoError(t, err)
		tmpScript.Close()

		err = os.Chmod(tmpScript.Name(), 0o755)
		testutil.AssertNoError(t, err)

		t.Setenv("EDITOR", tmpScript.Name())

		result, err := EditInEditor(context.Background(), "original content")
		testutil.AssertNoError(t, err)

		if result != "edited content\n" {
			t.Errorf("EditInEditor() = %q, want %q", result, "edited content\n")
		}
	})

	t.Run("falls back to vi when EDITOR not set", func(t *testing.T) {
		t.Skip("cannot test vi fallback without an interactive terminal")
	})
}
