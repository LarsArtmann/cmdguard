package v2

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestParseEnum(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		allowed   []string
		wantErr   bool
		wantValue string
	}{
		{
			name:      "valid value",
			value:     "debug",
			allowed:   []string{"debug", "info", "warn", "error"},
			wantErr:   false,
			wantValue: "debug",
		},
		{
			name:    "invalid value",
			value:   "invalid",
			allowed: []string{"debug", "info", "warn", "error"},
			wantErr: true,
		},
		{
			name:      "empty allowed list",
			value:     "anything",
			allowed:   []string{},
			wantErr:   true,
			wantValue: "",
		},
		{
			name:      "single allowed value",
			value:     "only",
			allowed:   []string{"only"},
			wantErr:   false,
			wantValue: "only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := ParseEnum(tt.value, tt.allowed)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !errors.Is(err, ErrInvalidEnum) {
					t.Errorf("expected EnumError (ErrInvalidEnum), got %T", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if e.String() != tt.wantValue {
					t.Errorf("String() = %q, want %q", e.String(), tt.wantValue)
				}

				if !slicesEqual(e.Allowed(), tt.allowed) {
					t.Errorf("Allowed() = %v, want %v", e.Allowed(), tt.allowed)
				}
			}
		})
	}
}

func TestParseEnum_ErrorCases(t *testing.T) {
	t.Run("returns error on invalid", func(t *testing.T) {
		_, err := ParseEnum("invalid", []string{"valid"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "invalid value") {
			t.Errorf("error should contain 'invalid value', got %q", err.Error())
		}
	})
}

func TestEnum_Methods(t *testing.T) {
	e, err := ParseEnum("test", []string{"a", "test", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("String", func(t *testing.T) {
		if e.String() != "test" {
			t.Errorf("String() = %q, want %q", e.String(), "test")
		}
	})

	t.Run("Value", func(t *testing.T) {
		if e.Value() != "test" {
			t.Errorf("Value() = %q, want %q", e.Value(), "test")
		}
	})

	t.Run("Allowed", func(t *testing.T) {
		if !slicesEqual(e.Allowed(), []string{"a", "test", "b"}) {
			t.Errorf("Allowed() = %v, want %v", e.Allowed(), []string{"a", "test", "b"})
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		if e.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}

		var empty Enum
		if !empty.IsEmpty() {
			t.Error("empty.IsEmpty() = false, want true")
		}
	})
}

func TestEnum_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Level Enum `json:"level"`
	}

	validLevel, err := ParseEnum("info", []string{"debug", "info", "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("marshal", func(t *testing.T) {
		c := config{Level: validLevel}

		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `{"level":"info"}` {
			t.Errorf("json.Marshal() = %q, want %q", string(data), `{"level":"info"}`)
		}
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config

		c.Level = Enum{allowed: []string{"debug", "info", "warn"}}

		err := json.Unmarshal([]byte(`{"level":"info"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Level.String() != "info" {
			t.Errorf("unmarshaled Level = %q, want %q", c.Level.String(), "info")
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config

		c.Level = Enum{allowed: []string{"debug", "info"}}

		err := json.Unmarshal([]byte(`{"level":"invalid"}`), &c)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("unmarshal with no allowed", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"level":"any"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Level.String() != "any" {
			t.Errorf("unmarshaled Level = %q, want %q", c.Level.String(), "any")
		}

		if !slicesEqual(c.Level.Allowed(), []string{"any"}) {
			t.Errorf("Allowed() = %v, want %v", c.Level.Allowed(), []string{"any"})
		}
	})
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantMs  int64
	}{
		{"seconds", "30s", false, 30000},
		{"minutes", "5m", false, 300000},
		{"hours", "1h", false, 3600000},
		{"complex", "1h30m", false, 5400000},
		{"milliseconds", "500ms", false, 500},
		{"invalid", "not-a-duration", true, 0},
		{"empty", "", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ParseDuration(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !errors.Is(err, ErrInvalidDuration) {
					t.Errorf("expected DurationError (ErrInvalidDuration), got %T", err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if d.Milliseconds() != tt.wantMs {
					t.Errorf("Milliseconds() = %d, want %d", d.Milliseconds(), tt.wantMs)
				}
			}
		})
	}
}

func TestParseDuration_ErrorCases(t *testing.T) {
	t.Run("returns error on invalid", func(t *testing.T) {
		_, err := ParseDuration("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestFromDuration(t *testing.T) {
	td := 5 * time.Minute

	d := FromDuration(td)
	if d.Duration() != td {
		t.Errorf("Duration() = %v, want %v", d.Duration(), td)
	}

	if d.Milliseconds() != 300000 {
		t.Errorf("Milliseconds() = %d, want %d", d.Milliseconds(), 300000)
	}

	gotSeconds := d.Seconds()
	if gotSeconds < 299.999 || gotSeconds > 300.001 {
		t.Errorf("Seconds() = %f, want approximately 300", gotSeconds)
	}
}

func TestDuration_Methods(t *testing.T) {
	d, err := ParseDuration("2h30m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("Duration", func(t *testing.T) {
		if d.Duration() != 2*time.Hour+30*time.Minute {
			t.Errorf("Duration() = %v, want %v", d.Duration(), 2*time.Hour+30*time.Minute)
		}
	})

	t.Run("String", func(t *testing.T) {
		if d.String() != "2h30m0s" {
			t.Errorf("String() = %q, want %q", d.String(), "2h30m0s")
		}
	})

	t.Run("IsZero", func(t *testing.T) {
		if d.IsZero() {
			t.Error("IsZero() = true, want false")
		}

		var zero Duration
		if !zero.IsZero() {
			t.Error("zero.IsZero() = false, want true")
		}
	})

	t.Run("Milliseconds", func(t *testing.T) {
		if d.Milliseconds() != 9000000 {
			t.Errorf("Milliseconds() = %d, want %d", d.Milliseconds(), 9000000)
		}
	})

	t.Run("Seconds", func(t *testing.T) {
		got := d.Seconds()
		if got < 8999.999 || got > 9000.001 {
			t.Errorf("Seconds() = %f, want approximately 9000", got)
		}
	})
}

func TestDuration_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Timeout Duration `json:"timeout"`
	}

	t.Run("marshal", func(t *testing.T) {
		validDuration, err := ParseDuration("30s")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		c := config{Timeout: validDuration}

		data, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != `{"timeout":"30s"}` {
			t.Errorf("json.Marshal() = %q, want %q", string(data), `{"timeout":"30s"}`)
		}
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"timeout":"1h"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Timeout.Duration() != time.Hour {
			t.Errorf("unmarshaled Timeout = %v, want %v", c.Timeout.Duration(), time.Hour)
		}
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config

		err := json.Unmarshal([]byte(`{"timeout":"invalid"}`), &c)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

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
