package v4

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestErrorHint(t *testing.T) {
	t.Parallel()

	t.Run("nil error has no hint", func(t *testing.T) {
		t.Parallel()

		testutil.AssertFieldEqString(t, ErrorHint(nil), "", "nil hint")
	})

	t.Run("unknown error has no hint", func(t *testing.T) {
		t.Parallel()

		testutil.AssertFieldEqString(t, ErrorHint(errors.New("unexpected")), "", "unknown hint")
	})

	t.Run("specific sentinel beats generic wrapper", func(t *testing.T) {
		t.Parallel()

		wrapped := fmt.Errorf(
			"%w: validating required flag %q on command %q: %w",
			ErrFlagParseFailed,
			"name",
			"app",
			fmt.Errorf("%w: name is required", ErrRequiredFlag),
		)

		got := ErrorHint(wrapped)
		if got == "" {
			t.Fatal("expected a hint for a wrapped required-flag error")
		}

		if got == ErrorHint(ErrFlagParseFailed) {
			t.Error("wrapped required-flag error must use the specific hint, not the generic flag-parse hint")
		}

		if got != ErrorHint(ErrRequiredFlag) {
			t.Errorf("hint = %q, want the ErrRequiredFlag hint", got)
		}
	})

	t.Run("every classified sentinel has a hint", func(t *testing.T) {
		t.Parallel()

		for _, entry := range errorTypeNames() {
			if ErrorHint(entry.match) == "" {
				t.Errorf("sentinel %q (class %q) has no actionable hint", entry.match, entry.name)
			}
		}
	})
}

func TestWriteJSONError_IncludesHint(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := writeJSONError(&buf, fmt.Errorf("%w: --name is required", ErrRequiredFlag))
	testutil.AssertNoError(t, err)

	var envelope jsonErrorEnvelope
	testutil.AssertNoError(t, json.Unmarshal(buf.Bytes(), &envelope))

	if envelope.Error.Hint == "" {
		t.Fatal("JSON error envelope must include an actionable hint")
	}

	testutil.AssertFieldEqString(t, envelope.Error.Hint, ErrorHint(ErrRequiredFlag), "hint")
}
