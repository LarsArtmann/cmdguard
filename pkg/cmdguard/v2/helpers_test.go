package v2

import (
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestLogLevel(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		if LogLevelDebug.String() != "debug" {
			t.Errorf("LogLevelDebug.String() = %q, want %q", LogLevelDebug.String(), "debug")
		}

		if LogLevelInfo.String() != "info" {
			t.Errorf("LogLevelInfo.String() = %q, want %q", LogLevelInfo.String(), "info")
		}

		if LogLevelWarn.String() != "warn" {
			t.Errorf("LogLevelWarn.String() = %q, want %q", LogLevelWarn.String(), "warn")
		}

		if LogLevelError.String() != "error" {
			t.Errorf("LogLevelError.String() = %q, want %q", LogLevelError.String(), "error")
		}
	})

	t.Run("ParseLogLevel valid", func(t *testing.T) {
		tests := []string{"debug", "info", "warn", "error"}
		for _, v := range tests {
			l, err := ParseLogLevel(v)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", v, err)
			}

			if l.String() != v {
				t.Errorf("ParseLogLevel(%q).String() = %q, want %q", v, l.String(), v)
			}
		}
	})

	t.Run("ParseLogLevel invalid", func(t *testing.T) {
		_, err := ParseLogLevel("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("SlogLevel conversion", func(t *testing.T) {
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
	t.Run("constants", func(t *testing.T) {
		if LogFormatText.String() != "text" {
			t.Errorf("LogFormatText.String() = %q, want %q", LogFormatText.String(), "text")
		}

		if LogFormatJSON.String() != "json" {
			t.Errorf("LogFormatJSON.String() = %q, want %q", LogFormatJSON.String(), "json")
		}
	})

	t.Run("ParseLogFormat valid", func(t *testing.T) {
		tests := []string{"text", "json"}
		for _, v := range tests {
			f, err := ParseLogFormat(v)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", v, err)
			}

			if f.String() != v {
				t.Errorf("ParseLogFormat(%q).String() = %q, want %q", v, f.String(), v)
			}
		}
	})

	t.Run("ParseLogFormat invalid", func(t *testing.T) {
		_, err := ParseLogFormat("xml")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := 42

		p := Ptr(v)
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != 42 {
			t.Errorf("*p = %d, want %d", *p, 42)
		}
	})

	t.Run("string", func(t *testing.T) {
		v := "hello"

		p := Ptr(v)
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if *p != "hello" {
			t.Errorf("*p = %q, want %q", *p, "hello")
		}
	})

	t.Run("struct", func(t *testing.T) {
		type s struct{ Name string }

		v := s{Name: "test"}

		p := Ptr(v)
		if p == nil {
			t.Fatal("expected non-nil pointer")
		}

		if p.Name != "test" {
			t.Errorf("p.Name = %q, want %q", p.Name, "test")
		}
	})
}

func TestValueOrDefault(t *testing.T) {
	t.Run("nil pointer returns default", func(t *testing.T) {
		var p *int

		result := ValueOrDefault(p, 10)
		if result != 10 {
			t.Errorf("ValueOrDefault(nil, 10) = %d, want %d", result, 10)
		}
	})

	t.Run("non-nil pointer returns value", func(t *testing.T) {
		v := 42
		p := &v

		result := ValueOrDefault(p, 10)
		if result != 42 {
			t.Errorf("ValueOrDefault(&42, 10) = %d, want %d", result, 42)
		}
	})

	t.Run("empty string default", func(t *testing.T) {
		var p *string

		result := ValueOrDefault(p, "default")
		if result != "default" {
			t.Errorf("ValueOrDefault(nil, \"default\") = %q, want %q", result, "default")
		}
	})
}

func TestEnsureValid(t *testing.T) {
	t.Run("nil returns error", func(t *testing.T) {
		var p *int

		err := EnsureValid(p, "myField")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "myField") {
			t.Errorf("error should contain 'myField', got %q", err.Error())
		}

		if !strings.Contains(err.Error(), "must not be nil") {
			t.Errorf("error should contain 'must not be nil', got %q", err.Error())
		}
	})

	t.Run("non-nil returns nil", func(t *testing.T) {
		v := 42
		p := &v

		err := EnsureValid(p, "myField")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestLogLevel_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Value LogLevel `json:"value"`
	}

	validLevel := LogLevelInfo

	t.Run("marshal", func(t *testing.T) {
		c := config{Value: validLevel}

		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `{"value":"info"}` {
			t.Errorf("json.Marshal() = %q, want %q", string(data), `{"value":"info"}`)
		}
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"value":"info"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Value.String() != "info" {
			t.Errorf("unmarshaled Value = %q, want %q", c.Value.String(), "info")
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"value":"trace"}`), &c)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestLogFormat_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Value LogFormat `json:"value"`
	}

	validFormat := LogFormatJSON

	t.Run("marshal", func(t *testing.T) {
		c := config{Value: validFormat}

		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `{"value":"json"}` {
			t.Errorf("json.Marshal() = %q, want %q", string(data), `{"value":"json"}`)
		}
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"value":"json"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Value.String() != "json" {
			t.Errorf("unmarshaled Value = %q, want %q", c.Value.String(), "json")
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"value":"xml"}`), &c)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
