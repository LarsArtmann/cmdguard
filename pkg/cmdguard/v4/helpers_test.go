package v4

import (
	"encoding/json/v2"
	"errors"
	"log/slog"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

func TestLogLevel(t *testing.T) {
	t.Parallel()
	t.Run("constants", func(t *testing.T) {
		t.Parallel()

		testutil.AssertFieldEqString(t, LogLevelDebug.String(), "debug", "LogLevelDebug.String()")

		if LogLevelInfo.String() != "info" {
			testutil.AssertFieldEqString(t, LogLevelInfo.String(), "info", "LogLevelInfo.String()")
		}

		if LogLevelWarn.String() != "warn" {
			testutil.AssertFieldEqString(t, LogLevelWarn.String(), "warn", "LogLevelWarn.String()")
		}

		if LogLevelError.String() != "error" {
			testutil.AssertFieldEqString(
				t,
				LogLevelError.String(),
				"error",
				"LogLevelError.String()",
			)
		}
	})

	t.Run("ParseLogLevel valid", func(t *testing.T) {
		t.Parallel()

		tests := []string{"debug", "info", "warn", "error"}
		for _, v := range tests {
			l, err := ParseLogLevel(v)
			testutil.AssertNoError(t, err)

			testutil.AssertFieldEqString(t, l.String(), v, "ParseLogLevel().String()")
		}
	})

	t.Run("ParseLogLevel invalid", func(t *testing.T) {
		t.Parallel()

		_, err := ParseLogLevel("invalid")
		testutil.AssertExpectedError(t, err)

		if !errors.Is(err, ErrLogLevel) {
			t.Error("expected error to match ErrLogLevel via errors.Is")
		}

		if !errors.Is(err, ErrInvalidEnum) {
			t.Error("expected error to match ErrInvalidEnum via errors.Is")
		}
	})

	t.Run("LogLevel IsEmpty", func(t *testing.T) {
		t.Parallel()

		if LogLevelInfo.IsEmpty() {
			t.Error("LogLevelInfo.IsEmpty() = true, want false")
		}

		var empty LogLevel
		if !empty.IsEmpty() {
			t.Error("empty LogLevel IsEmpty() = false, want true")
		}
	})

	t.Run("SlogLevel conversion", func(t *testing.T) {
		t.Parallel()

		if LogLevelDebug.SlogLevel() != slog.LevelDebug {
			t.Errorf(
				"LogLevelDebug.SlogLevel() = %v, want %v",
				LogLevelDebug.SlogLevel(),
				slog.LevelDebug,
			)
		}

		if LogLevelInfo.SlogLevel() != slog.LevelInfo {
			t.Errorf(
				"LogLevelInfo.SlogLevel() = %v, want %v",
				LogLevelInfo.SlogLevel(),
				slog.LevelInfo,
			)
		}

		if LogLevelWarn.SlogLevel() != slog.LevelWarn {
			t.Errorf(
				"LogLevelWarn.SlogLevel() = %v, want %v",
				LogLevelWarn.SlogLevel(),
				slog.LevelWarn,
			)
		}

		if LogLevelError.SlogLevel() != slog.LevelError {
			t.Errorf(
				"LogLevelError.SlogLevel() = %v, want %v",
				LogLevelError.SlogLevel(),
				slog.LevelError,
			)
		}
	})
}

func TestLogFormat(t *testing.T) {
	t.Parallel()
	t.Run("constants", func(t *testing.T) {
		t.Parallel()

		if LogFormatText.String() != "text" {
			testutil.AssertFieldEqString(
				t,
				LogFormatText.String(),
				"text",
				"LogFormatText.String()",
			)
		}

		if LogFormatJSON.String() != "json" {
			testutil.AssertFieldEqString(
				t,
				LogFormatJSON.String(),
				"json",
				"LogFormatJSON.String()",
			)
		}
	})

	t.Run("ParseLogFormat valid", func(t *testing.T) {
		t.Parallel()

		tests := []string{"text", "json"}
		for _, v := range tests {
			f, err := ParseLogFormat(v)
			testutil.AssertNoError(t, err)

			testutil.AssertFieldEqString(t, f.String(), v, "ParseLogFormat().String()")
		}
	})

	t.Run("LogFormat IsEmpty", func(t *testing.T) {
		t.Parallel()

		if LogFormatText.IsEmpty() {
			t.Error("LogFormatText.IsEmpty() = true, want false")
		}

		var empty LogFormat
		if !empty.IsEmpty() {
			t.Error("empty LogFormat IsEmpty() = false, want true")
		}
	})

	t.Run("ParseLogFormat invalid", func(t *testing.T) {
		t.Parallel()

		_, err := ParseLogFormat("xml")
		testutil.AssertExpectedError(t, err)

		if !errors.Is(err, ErrLogFormat) {
			t.Error("expected error to match ErrLogFormat via errors.Is")
		}

		if !errors.Is(err, ErrInvalidEnum) {
			t.Error("expected error to match ErrInvalidEnum via errors.Is")
		}
	})
}

// testPtrGeneric tests that Ptr returns a valid pointer to the given value.
func testPtrGeneric[T comparable](t *testing.T, v T) {
	t.Helper()

	p := new(v)

	if *p != v {
		t.Errorf("*p = %v, want %v", *p, v)
	}
}

// runPtrTest runs a single pointer test case.
func runPtrTest[T comparable](t *testing.T, name string, value T) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()
		testPtrGeneric(t, value)
	})
}

func TestPtr(t *testing.T) {
	t.Parallel()

	runPtrTest(t, "int", 42)
	runPtrTest(t, "string", "hello")

	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		type s struct{ Name string }

		v := s{Name: "test"}

		p := new(v)
		testutil.AssertNotNil(t, p)

		if p.Name != "test" {
			testutil.AssertFieldEqString(t, p.Name, "test", "p.Name")
		}
	})
}

func TestValueOrDefault(t *testing.T) {
	t.Parallel()
	t.Run("nil pointer returns default", func(t *testing.T) {
		t.Parallel()

		var p *int

		result := ValueOrDefault(p, 10)
		testutil.AssertFieldEq(t, result, 10, "result")
	})

	t.Run("non-nil pointer returns value", func(t *testing.T) {
		t.Parallel()

		v := 42
		p := &v

		result := ValueOrDefault(p, 10)
		testutil.AssertFieldEq(t, result, 42, "result")
	})

	t.Run("empty string default", func(t *testing.T) {
		t.Parallel()

		var p *string

		result := ValueOrDefault(p, "default")
		testutil.AssertFieldEqString(t, result, "default", "result")
	})
}

func TestEnsureValid(t *testing.T) {
	t.Parallel()
	t.Run("nil returns error", func(t *testing.T) {
		t.Parallel()

		var p *int

		err := EnsureValid(p, "myField")
		testutil.AssertExpectedError(t, err)

		assertErrorContains(t, err, "myField", "must not be nil")
	})

	t.Run("non-nil returns nil", func(t *testing.T) {
		t.Parallel()

		v := 42
		p := &v

		err := EnsureValid(p, "myField")
		testutil.AssertNoError(t, err)
	})
}

func TestLogLevel_MarshalUnmarshal(t *testing.T) {
	t.Parallel()

	type config struct {
		Value LogLevel `json:"value"`
	}

	validLevel := LogLevelInfo

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()

		c := config{Value: validLevel}

		data, err := json.Marshal(c)
		testutil.AssertNoError(t, err)

		testutil.AssertJSONMarshal(t, data, `{"value":"info"}`)
	})

	// unmarshal valid
	var c1 config
	expectUnmarshalValidString(
		t,
		&c1,
		`{"value":"info"}`,
		"info",
		func() string { return c1.Value.String() },
	)

	runUnmarshalErrorTest[config](t, "unmarshal invalid", `{"value":"trace"}`)
}

func TestLogFormat_MarshalUnmarshal(t *testing.T) {
	t.Parallel()

	type config struct {
		Value LogFormat `json:"value"`
	}

	validFormat := LogFormatJSON

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()

		c := config{Value: validFormat}

		data, err := json.Marshal(c)
		testutil.AssertNoError(t, err)

		testutil.AssertJSONMarshal(t, data, `{"value":"json"}`)
	})

	// unmarshal valid
	var c2 config
	expectUnmarshalValidString(
		t,
		&c2,
		`{"value":"json"}`,
		"json",
		func() string { return c2.Value.String() },
	)

	runUnmarshalErrorTest[config](t, "unmarshal invalid", `{"value":"xml"}`)
}

// expectUnmarshalError tests that unmarshaling invalid JSON returns an error.
func expectUnmarshalError(t *testing.T, target any, jsonStr string) {
	t.Helper()

	err := json.Unmarshal([]byte(jsonStr), target)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// runUnmarshalErrorTest runs an unmarshal error test with the given name and JSON.
func runUnmarshalErrorTest[T any](t *testing.T, name, jsonStr string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		var c T
		expectUnmarshalError(t, &c, jsonStr)
	})
}

// expectUnmarshalValidString tests that unmarshaling valid JSON succeeds and the extracted string matches expected.
func expectUnmarshalValidString(
	t *testing.T,
	target any,
	jsonStr, expected string,
	extractStr func() string,
) {
	t.Helper()

	err := json.Unmarshal([]byte(jsonStr), target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := extractStr(); got != expected {
		t.Errorf("unmarshaled Value = %q, want %q", got, expected)
	}
}
