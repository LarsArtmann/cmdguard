package v2

import (
	"testing"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestSuggestCommand(t *testing.T) {
	t.Parallel()

	validCommands := []string{"start", "stop", "status", "restart"}

	tests := []struct {
		name      string
		commands  []string
		input     string
		wantMatch bool
		wantOneOf []string
	}{
		{
			name:      "exact match returns true",
			commands:  validCommands,
			input:     "start",
			wantMatch: true,
			wantOneOf: []string{"start"},
		},
		{
			name:      "one char typo returns match",
			commands:  validCommands,
			input:     "starz",
			wantMatch: true,
			wantOneOf: []string{"start"},
		},
		{
			name:      "too different returns false",
			commands:  validCommands,
			input:     "xyzzy",
			wantMatch: false,
		},
		{
			name:      "empty commands returns false",
			commands:  []string{},
			input:     "start",
			wantMatch: false,
		},
		{
			name:      "nil commands returns false",
			commands:  nil,
			input:     "start",
			wantMatch: false,
		},
		{
			name:      "close typo returns match",
			commands:  validCommands,
			input:     "stats",
			wantMatch: true,
			wantOneOf: []string{"status", "start"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			match, ok := SuggestCommand(tt.commands, tt.input)
			if ok != tt.wantMatch {
				t.Errorf(
					"SuggestCommand(%q) = (%q, %v), want match=%v",
					tt.input,
					match,
					ok,
					tt.wantMatch,
				)
			}
			if ok && tt.wantOneOf != nil {
				found := false
				for _, w := range tt.wantOneOf {
					if match == w {
						found = true

						break
					}
				}
				if !found {
					t.Errorf(
						"SuggestCommand(%q) = %q, want one of %v",
						tt.input,
						match,
						tt.wantOneOf,
					)
				}
			}
		})
	}
}

func TestSuggestCommand_DelegatesToSuggestFlag(t *testing.T) {
	t.Parallel()

	t.Run("returns same result as SuggestFlag but with bool", func(t *testing.T) {
		t.Parallel()

		commands := []string{"build", "run", "test"}
		input := "buidl"

		flagResult, flagOk := SuggestFlag(commands, input)
		cmdResult, ok := SuggestCommand(commands, input)

		if !flagOk {
			testutil.AssertBoolFalse(t, ok, "should be false when SuggestFlag returns false")
		} else {
			testutil.AssertBoolTrue(t, ok, "should be true when SuggestFlag returns true")
			if cmdResult != flagResult {
				t.Errorf(
					"SuggestCommand = %q, SuggestFlag = %q, should match",
					cmdResult,
					flagResult,
				)
			}
		}
	})
}
