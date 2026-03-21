package v2

import (
	"errors"
	"strings"
	"testing"
)

func TestSuggestFlag(t *testing.T) {
	validNames := []string{"verbose", "version", "config", "help", "output"}

	t.Run("exact match returns same name", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbose")
		if result != "verbose" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "verbose")
		}
	})

	t.Run("one character typo returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbosee")
		if result != "verbose" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "verbose")
		}
	})

	t.Run("missing character returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbos")
		if result != "verbose" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "verbose")
		}
	})

	t.Run("transposed characters returns suggestion", func(t *testing.T) {
		result := SuggestFlag(validNames, "verbsoe")
		if result != "verbose" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "verbose")
		}
	})

	t.Run("similar prefix returns closest match", func(t *testing.T) {
		result := SuggestFlag(validNames, "confing")
		if result != "config" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "config")
		}
	})

	t.Run("too different returns empty", func(t *testing.T) {
		result := SuggestFlag(validNames, "xyzzy")
		if result != "" {
			t.Errorf("SuggestFlag() = %q, want empty", result)
		}
	})

	t.Run("empty valid names returns empty", func(t *testing.T) {
		result := SuggestFlag([]string{}, "test")
		if result != "" {
			t.Errorf("SuggestFlag() = %q, want empty", result)
		}
	})

	t.Run("single valid name match", func(t *testing.T) {
		result := SuggestFlag([]string{"help"}, "hlep")
		if result != "help" {
			t.Errorf("SuggestFlag() = %q, want %q", result, "help")
		}
	})

	t.Run("selects closest match among multiple", func(t *testing.T) {
		names := []string{"start", "status", "stop"}

		result := SuggestFlag(names, "stat")
		if !containsString(names, result) {
			t.Errorf("SuggestFlag() = %q, expected to be one of %v", result, names)
		}
	})
}

func TestEditDistance(t *testing.T) {
	t.Run("identical strings have distance 0", func(t *testing.T) {
		if editDistance("hello", "hello") != 0 {
			t.Errorf(
				"editDistance(\"hello\", \"hello\") = %d, want 0",
				editDistance("hello", "hello"),
			)
		}
	})

	t.Run("empty strings", func(t *testing.T) {
		if editDistance("", "hello") != 5 {
			t.Errorf("editDistance(\"\", \"hello\") = %d, want 5", editDistance("", "hello"))
		}

		if editDistance("hello", "") != 5 {
			t.Errorf("editDistance(\"hello\", \"\") = %d, want 5", editDistance("hello", ""))
		}

		if editDistance("", "") != 0 {
			t.Errorf("editDistance(\"\", \"\") = %d, want 0", editDistance("", ""))
		}
	})

	t.Run("single insertion", func(t *testing.T) {
		if editDistance("hell", "hello") != 1 {
			t.Errorf(
				"editDistance(\"hell\", \"hello\") = %d, want 1",
				editDistance("hell", "hello"),
			)
		}
	})

	t.Run("single deletion", func(t *testing.T) {
		if editDistance("hello", "hell") != 1 {
			t.Errorf(
				"editDistance(\"hello\", \"hell\") = %d, want 1",
				editDistance("hello", "hell"),
			)
		}
	})

	t.Run("single substitution", func(t *testing.T) {
		if editDistance("hello", "hallo") != 1 {
			t.Errorf(
				"editDistance(\"hello\", \"hallo\") = %d, want 1",
				editDistance("hello", "hallo"),
			)
		}
	})

	t.Run("transposition counts as 2", func(t *testing.T) {
		if editDistance("ab", "ba") != 2 {
			t.Errorf("editDistance(\"ab\", \"ba\") = %d, want 2", editDistance("ab", "ba"))
		}
	})

	t.Run("multiple edits", func(t *testing.T) {
		if editDistance("kitten", "sitting") != 3 {
			t.Errorf(
				"editDistance(\"kitten\", \"sitting\") = %d, want 3",
				editDistance("kitten", "sitting"),
			)
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		if editDistance("Hello", "hello") != 1 {
			t.Errorf(
				"editDistance(\"Hello\", \"hello\") = %d, want 1",
				editDistance("Hello", "hello"),
			)
		}
	})
}

var (
	errUnknownFlag = errors.New("unknown flag")
	errSomeError   = errors.New("some error")
	errInnerError  = errors.New("inner error")
)

func TestNewFlagErrorWithSuggestion(t *testing.T) {
	t.Run("error includes suggestion", func(t *testing.T) {
		err := NewFlagErrorWithSuggestion("verboose", errUnknownFlag, "verbose")

		errMsg := err.Error()
		if !strings.Contains(errMsg, "verboose") {
			t.Errorf("error should contain 'verboose', got %q", errMsg)
		}

		if !strings.Contains(errMsg, "unknown flag") {
			t.Errorf("error should contain 'unknown flag', got %q", errMsg)
		}

		if !strings.Contains(errMsg, "did you mean --verbose") {
			t.Errorf("error should contain 'did you mean --verbose', got %q", errMsg)
		}
	})

	t.Run("empty suggestion omits hint", func(t *testing.T) {
		err := NewFlagError("test", errSomeError)
		if strings.Contains(err.Error(), "did you mean") {
			t.Errorf("error should not contain 'did you mean', got %q", err.Error())
		}
	})

	t.Run("unwraps to inner error", func(t *testing.T) {
		err := NewFlagErrorWithSuggestion("flag", errInnerError, "suggestion")
		if !errors.Is(err, errInnerError) {
			t.Errorf("expected error to unwrap to errInnerError")
		}
	})
}
