package flightrecorder

import (
	"testing"
	"unicode/utf8"
)

func FuzzSanitizeFilename(f *testing.F) {
	f.Add("deploy")
	f.Add("deploy prod")
	f.Add("")
	f.Add("app/deploy:v2.0")
	f.Add("café")

	f.Fuzz(func(t *testing.T, input string) {
		result := sanitizeFilename(input)

		// Invariant: output must only contain safe characters.
		for _, char := range result {
			if !isSafeFilenameChar(char) {
				t.Errorf("sanitizeFilename(%q) = %q: unsafe char %q", input, result, char)
			}
		}

		// Invariant: output rune count must equal input rune count.
		// sanitizeFilename replaces unsafe runes with '-', it never adds or removes.
		if utf8.RuneCountInString(result) != utf8.RuneCountInString(input) {
			t.Errorf("sanitizeFilename(%q) = %q: rune count changed from %d to %d",
				input, result, utf8.RuneCountInString(input), utf8.RuneCountInString(result))
		}
	})
}

// isSafeFilenameChar returns true for characters that sanitizeFilename preserves.
func isSafeFilenameChar(char rune) bool {
	switch {
	case char >= 'a' && char <= 'z',
		char >= 'A' && char <= 'Z',
		char >= '0' && char <= '9',
		char == '-', char == '_', char == '.':
		return true

	default:
		return false
	}
}
