package v2

import (
	"errors"
	"strings"
	"testing"
)

func TestSuggestFlag(t *testing.T) {
	t.Parallel()
	defaultNames := []string{"verbose", "version", "config", "help", "output"}

	tests := []struct {
		name       string
		validNames []string
		input      string
		want       string
		wantOneOf  []string // Alternative: check result is one of these
	}{
		{"exact match returns same name", defaultNames, "verbose", "verbose", nil},
		{"one character typo returns suggestion", defaultNames, "verbosee", "verbose", nil},
		{"missing character returns suggestion", defaultNames, "verbos", "verbose", nil},
		{"transposed characters returns suggestion", defaultNames, "verbsoe", "verbose", nil},
		{"similar prefix returns closest match", defaultNames, "confing", "config", nil},
		{"too different returns empty", defaultNames, "xyzzy", "", nil},
		{"empty valid names returns empty", []string{}, "test", "", nil},
		{"single valid name match", []string{"help"}, "hlep", "help", nil},
		{
			"selects closest match among multiple",
			[]string{"start", "status", "stop"},
			"stat",
			"",
			[]string{"start", "status", "stop"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := SuggestFlag(tt.validNames, tt.input)
			if tt.wantOneOf != nil {
				if !containsString(tt.wantOneOf, result) {
					t.Errorf("SuggestFlag() = %q, expected to be one of %v", result, tt.wantOneOf)
				}
			} else if result != tt.want {
				t.Errorf("SuggestFlag() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestEditDistance(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"identical strings have distance 0", "hello", "hello", 0},
		{"empty to string", "", "hello", 5},
		{"string to empty", "hello", "", 5},
		{"both empty", "", "", 0},
		{"single insertion", "hell", "hello", 1},
		{"single deletion", "hello", "hell", 1},
		{"single substitution", "hello", "hallo", 1},
		{"transposition counts as 2", "ab", "ba", 2},
		{"multiple edits", "kitten", "sitting", 3},
		{"case sensitive", "Hello", "hello", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := editDistance(tt.a, tt.b); got != tt.want {
				t.Errorf("editDistance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

var (
	errUnknownFlag = errors.New("unknown flag")
	errSomeError   = errors.New("some error")
	errInnerError  = errors.New("inner error")
)

func TestNewFlagErrorWithSuggestion(t *testing.T) {
	t.Parallel()
	t.Run("error includes suggestion", func(t *testing.T) {
		t.Parallel()
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
		t.Parallel()
		err := NewFlagError("test", errSomeError)
		if strings.Contains(err.Error(), "did you mean") {
			t.Errorf("error should not contain 'did you mean', got %q", err.Error())
		}
	})

	t.Run("unwraps to inner error", func(t *testing.T) {
		t.Parallel()
		err := NewFlagErrorWithSuggestion("flag", errInnerError, "suggestion")
		if !errors.Is(err, errInnerError) {
			t.Errorf("expected error to unwrap to errInnerError")
		}
	})
}
