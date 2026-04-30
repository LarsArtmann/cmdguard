package v2

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EditInEditor opens a temporary file in the user's $EDITOR with the given content,
// waits for the editor to close, and returns the edited content.
// Falls back to "vi" if $EDITOR is not set.
func EditInEditor(content string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpFile, err := os.CreateTemp("", "cmdguard-edit-*.txt")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}

	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	_, err = tmpFile.WriteString(content)
	if err != nil {
		tmpFile.Close()

		return "", fmt.Errorf("writing to temp file: %w", err)
	}

	tmpFile.Close()

	parts := strings.Fields(editor)
	cmd := exec.Command(
		parts[0],
		append(parts[1:], tmpPath)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("running editor %q: %w", editor, err)
	}

	edited, err := os.ReadFile(filepath.Clean(tmpPath))
	if err != nil {
		return "", fmt.Errorf("reading edited file: %w", err)
	}

	return string(edited), nil
}
