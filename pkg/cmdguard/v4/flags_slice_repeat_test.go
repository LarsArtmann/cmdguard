package v4

import (
	"context"
	"slices"
	"testing"
)

// TestSliceFlagRepeatedOccurrencesParsesElements pins the stringly flag-parse
// path for slice flags: a changed pflag slice flag is read through
// flag.Value.String(), which renders the bracketed form "[a,b]". Before the
// splitSliceValue fix, that bracketed string was split on commas and the
// handler received "[100" / "101]" instead of "100" / "101" (seen live with
// a consumer's repeatable --document ID flag).
func TestSliceFlagRepeatedOccurrencesParsesElements(t *testing.T) {
	t.Parallel()

	type repeatFlags struct {
		Documents []string `flag:"document" default:"" help:"repeatable document ID"`
	}

	type repeatConfig struct{}

	var got []string

	cli, err := NewCLI[repeatConfig]("slicetest", "slice flag test", repeatConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cmd, err := NewCommand(
		"collect",
		repeatFlags{},
		func(_ context.Context, _ *repeatConfig, flags repeatFlags) error {
			got = flags.Documents

			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := AddCommand(cli, cmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = cli.ExecuteWithArgs(t.Context(), []string{
		"collect", "--document", "100", "--document", "101",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !slices.Equal(got, []string{"100", "101"}) {
		t.Errorf("documents = %q, want [100 101]", got)
	}
}

// TestSplitSliceValue covers the renderer-form variants splitSliceValue must
// normalize.
func TestSplitSliceValue(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value string
		want  []string
	}{
		{name: "single element bracketed", value: "[101]", want: []string{"101"}},
		{name: "two elements bracketed", value: "[100,101]", want: []string{"100", "101"}},
		{name: "empty bracketed", value: "[]", want: []string{}},
		{name: "bare single", value: "101", want: []string{"101"}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := splitSliceValue(testCase.value)
			if !slices.Equal(got, testCase.want) {
				t.Errorf("splitSliceValue(%q) = %q, want %q", testCase.value, got, testCase.want)
			}
		})
	}
}
