package v2

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestParseEnum(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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

				if !slices.Equal(e.Allowed(), tt.allowed) {
					t.Errorf("Allowed() = %v, want %v", e.Allowed(), tt.allowed)
				}
			}
		})
	}
}

func TestParseEnum_ErrorCases(t *testing.T) {
	t.Parallel()
	t.Run("returns error on invalid", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()

	e, err := ParseEnum("test", []string{"a", "test", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		if e.String() != "test" {
			t.Errorf("String() = %q, want %q", e.String(), "test")
		}
	})

	t.Run("Value", func(t *testing.T) {
		t.Parallel()

		if e.Value() != "test" {
			t.Errorf("Value() = %q, want %q", e.Value(), "test")
		}
	})

	t.Run("Allowed", func(t *testing.T) {
		t.Parallel()

		if !slices.Equal(e.Allowed(), []string{"a", "test", "b"}) {
			t.Errorf("Allowed() = %v, want %v", e.Allowed(), []string{"a", "test", "b"})
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()

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
	t.Parallel()

	type config struct {
		Level Enum `json:"level"`
	}

	validLevel, err := ParseEnum("info", []string{"debug", "info", "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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
		t.Parallel()

		var c config

		c.Level = Enum{allowed: []string{"debug", "info"}}

		err := json.Unmarshal([]byte(`{"level":"invalid"}`), &c)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("unmarshal with no allowed", func(t *testing.T) {
		t.Parallel()

		var c config

		err := json.Unmarshal([]byte(`{"level":"any"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Level.String() != "any" {
			t.Errorf("unmarshaled Level = %q, want %q", c.Level.String(), "any")
		}

		if !slices.Equal(c.Level.Allowed(), []string{"any"}) {
			t.Errorf("Allowed() = %v, want %v", c.Level.Allowed(), []string{"any"})
		}
	})
}
