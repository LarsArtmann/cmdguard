package testutil

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestAssertionFailurePaths exercises the error/fatal branches of every assertion
// helper by re-executing the test binary in a subprocess. Each case sets an env
// var that selects which assertion to trigger; the parent verifies the
// subprocess failed with the expected message.
func TestAssertionFailurePaths(t *testing.T) {
	t.Parallel()

	if label := os.Getenv("CMDGUARD_FAIL_TEST"); label != "" {
		runFailureCase(t, label)

		return
	}

	cases := []struct {
		label  string
		expect string
	}{
		// Errorf paths
		{"equal_mismatch", "got 1, want 2"},
		{"equalf_mismatch", "got 1, want 2"},
		{"not_equal", "expected not to equal"},
		{"nil_notnil", "expected nil"},
		{"notnil_nil", "expected non-nil"},
		{"error_is", "error does not match"},
		{"error_isf", "error does not match"},
		{"panics_no_panic", "expected panic"},
		{"does_not_panic", "expected no panic"},
		{"bool_field", "expected enabled"},
		{"bool_true", "expected flag to be true"},
		{"bool_false", "expected flag to be false"},
		{"string_contains", "greeting should contain"},
		{"output_contains", "output should contain"},
		{"no_error", "expected no error"},
		{"pointer_eq", "got"},
		{"json_marshal", "json.Marshal()"},
		{"stringer_eq", "String()"},
		{"field_eq_int", "expected count"},
		{"field_len", "len(items)"},
		// Fatal paths
		{"error_contains_nil", "expected error, got nil"},
		{"expected_error_nil", "expected error, got nil"},
		{"string_slices_len", "test:"},
		// Multiple-branch paths
		{"error_contains_missing", "should contain"},
		{"stderr_missing", "stderr should contain"},
		{"string_slices_mismatch", "test[1]:"},
		{"flag_not_registered", "expected"},
		{"flag_registered_unexpected", "expected"},
		{"contains_string_fail", "should contain"},
		{"not_contains_string_fail", "should not contain"},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			t.Parallel()

			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestAssertionFailurePaths$", "-test.v")
			cmd.Env = append(os.Environ(), "CMDGUARD_FAIL_TEST="+tc.label)
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			if err == nil {
				t.Fatalf("expected failure for %s, but subprocess passed", tc.label)
			}
			if !strings.Contains(out.String(), tc.expect) {
				t.Errorf("case %s: expected %q in output, got:\n%s", tc.label, tc.expect, out.String())
			}
		})
	}
}

func runFailureCase(t *testing.T, label string) {
	switch label {
	// Errorf paths
	case "equal_mismatch":
		AssertEqual(t, 1, 2)
	case "equalf_mismatch":
		AssertEqualf(t, 1, 2, "test ctx")
	case "not_equal":
		AssertNotEqual(t, 1, 1)
	case "nil_notnil":
		v := 42
		AssertNil(t, &v)
	case "notnil_nil":
		var p *int
		AssertNotNil(t, p)
	case "error_is":
		AssertErrorIs(t, errors.New("a"), errors.New("b"))
	case "error_isf":
		AssertErrorIsf(t, errors.New("a"), errors.New("b"), "ctx")
	case "panics_no_panic":
		AssertPanics(t, func() {})
	case "does_not_panic":
		AssertDoesNotPanic(t, func() { panic("boom") })
	case "bool_field":
		AssertBoolField(t, true, false, "enabled")
	case "bool_true":
		AssertBoolTrue(t, false, "flag")
	case "bool_false":
		AssertBoolFalse(t, true, "flag")
	case "string_contains":
		AssertStringFieldContains(t, "hello", "world", "greeting")
	case "output_contains":
		AssertOutputContains(t, "hello", "world")
	case "no_error":
		AssertNoError(t, errors.New("unexpected"))
	case "pointer_eq":
		a, b := 1, 2
		AssertPointerEq(t, &a, &b)
	case "json_marshal":
		AssertJSONMarshal(t, []byte(`{"a":1}`), `{"a":2}`)
	case "stringer_eq":
		AssertStringerEq(t, testStringer("a"), "b")
	case "field_eq_int":
		AssertFieldEq(t, 1, 2, "count")
	case "field_len":
		AssertFieldLen(t, []int{1}, 2, "items")
	// Fatal paths
	case "error_contains_nil":
		AssertErrorContains(t, nil, "test")
	case "expected_error_nil":
		AssertExpectedError(t, nil)
	case "string_slices_len":
		AssertStringSlicesEqual(t, []string{"a"}, []string{"a", "b"}, "test")
	// Multiple-branch paths
	case "error_contains_missing":
		AssertErrorContains(t, errors.New("file not found"), "missing")
	case "stderr_missing":
		AssertStderrContains(t, "all good", "error")
	case "string_slices_mismatch":
		AssertStringSlicesEqual(t, []string{"a", "x"}, []string{"a", "b"}, "test")
	case "flag_not_registered":
		cmd := &cobra.Command{Use: "test"}
		AssertFlagRegistered(t, cmd, "nonexistent")
	case "flag_registered_unexpected":
		cmd := &cobra.Command{Use: "test"}
		cmd.Flags().String("name", "", "name")
		AssertFlagNotRegistered(t, cmd, "name")
	case "contains_string_fail":
		AssertContainsString(t, []string{"a", "b"}, "c")
	case "not_contains_string_fail":
		AssertNotContainsString(t, []string{"a", "b"}, "a")
	default:
		t.Logf("unknown label: %s", label)
	}
}
