package v4

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestErrorCode(t *testing.T) {
	t.Parallel()

	t.Run("nil error has no code", func(t *testing.T) {
		t.Parallel()

		testutil.AssertFieldEqString(t, ErrorCode(nil), "", "nil code")
	})

	t.Run("unknown error has no code", func(t *testing.T) {
		t.Parallel()

		testutil.AssertFieldEqString(t, ErrorCode(errors.New("unexpected")), "", "unknown code")
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

		testutil.AssertFieldEqString(t, ErrorCode(wrapped), "required_flag", "wrapped code")
		testutil.AssertFieldEqString(t, ErrorCode(ErrRequiredFlag), "required_flag", "sentinel code")
		testutil.AssertFieldEqString(t, ErrorCode(ErrFlagParseFailed), "flag_parse_failed", "wrapper code")
	})

	t.Run("every fact carries a stable code and hint", func(t *testing.T) {
		t.Parallel()

		for _, fact := range sentinelFacts() {
			if fact.code == "" {
				t.Errorf("sentinel %q has no machine error code", fact.match)
			}

			if fact.hint == "" {
				t.Errorf("sentinel %q (code %q) has no actionable hint", fact.match, fact.code)
			}
		}
	})

	t.Run("codes are unique", func(t *testing.T) {
		t.Parallel()

		seen := make(map[string]error, len(sentinelFacts()))

		for _, fact := range sentinelFacts() {
			if prev, dup := seen[fact.code]; dup {
				t.Errorf("code %q is claimed by both %v and %v", fact.code, prev, fact.match)
			}

			seen[fact.code] = fact.match
		}
	})

	t.Run("every classified sentinel has a code", func(t *testing.T) {
		t.Parallel()

		for _, entry := range errorTypeNames() {
			if ErrorCode(entry.match) == "" {
				t.Errorf("sentinel %q (class %q) has no machine error code", entry.match, entry.name)
			}
		}
	})
}

func TestWriteJSONError_IncludesErrorCode(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := writeJSONError(&buf, fmt.Errorf("%w: --name is required", ErrRequiredFlag))
	testutil.AssertNoError(t, err)

	var envelope jsonErrorEnvelope
	testutil.AssertNoError(t, json.Unmarshal(buf.Bytes(), &envelope))

	testutil.AssertFieldEqString(t, envelope.Error.ErrorCode, "required_flag", "error_code")
	testutil.AssertFieldEqString(t, envelope.Error.Hint, ErrorHint(ErrRequiredFlag), "hint")
}
