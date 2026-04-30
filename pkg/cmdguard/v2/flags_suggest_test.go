package v2

import (
	"errors"
	"slices"
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

			result, ok := SuggestFlag(tt.validNames, tt.input)
			if tt.wantOneOf != nil {
				if !ok || !slices.Contains(tt.wantOneOf, result) {
					t.Errorf(
						"SuggestFlag() = (%q, %v), expected to be one of %v",
						result,
						ok,
						tt.wantOneOf,
					)
				}
			} else if result != tt.want || ok != (tt.want != "") {
				t.Errorf(
					"SuggestFlag() = (%q, %v), want (%q, %v)",
					result,
					ok,
					tt.want,
					tt.want != "",
				)
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
		assertStringContains(t, err.Error(), "verboose", "unknown flag", "did you mean --verbose")
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

func TestGenerateHelp(t *testing.T) {
	t.Parallel()

	t.Run("generates help text for all flags", func(t *testing.T) {
		t.Parallel()

		type helpConfig struct {
			Name    string `flag:"name"    short:"n" help:"The name"   default:"world"`
			Count   int    `flag:"count"             help:"The count"  default:"10"`
			Verbose bool   `flag:"verbose" short:"v" help:"Be verbose" default:"false"`
		}

		registry, err := NewFlagRegistry(helpConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		help := registry.GenerateHelp()
		assertStringContains(
			t,
			help,
			"--name",
			"--count",
			"--verbose",
			"The name",
			"default: world",
		)
	})

	t.Run("empty config returns empty help", func(t *testing.T) {
		t.Parallel()

		registry, err := NewFlagRegistry(struct{}{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		help := registry.GenerateHelp()
		if help != "" {
			t.Errorf("expected empty help, got %q", help)
		}
	})
}

func TestFlagNames(t *testing.T) {
	t.Parallel()

	t.Run("returns all flag names", func(t *testing.T) {
		t.Parallel()

		type namesConfig struct {
			Host string `flag:"host" help:"host" default:"localhost"`
			Port int    `flag:"port" help:"port" default:"8080"`
		}

		registry, err := NewFlagRegistry(namesConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		names := registry.FlagNames()
		if len(names) != 2 {
			t.Fatalf("expected 2 names, got %d", len(names))
		}

		if names[0] != "host" || names[1] != "port" {
			t.Errorf("names = %v, want [host port]", names)
		}
	})
}
