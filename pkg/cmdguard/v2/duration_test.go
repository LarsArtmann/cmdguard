package v2

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestParseDuration(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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
	t.Parallel()
	t.Run("returns error on invalid", func(t *testing.T) {
		t.Parallel()

		_, err := ParseDuration("invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestFromDuration(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	d, err := ParseDuration("2h30m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("Duration", func(t *testing.T) {
		t.Parallel()

		if d.Duration() != 2*time.Hour+30*time.Minute {
			t.Errorf("Duration() = %v, want %v", d.Duration(), 2*time.Hour+30*time.Minute)
		}
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		testutil.AssertStringerEq(t, d, "2h30m0s")
	})

	t.Run("IsZero", func(t *testing.T) {
		t.Parallel()

		if d.IsZero() {
			t.Error("IsZero() = true, want false")
		}

		var zero Duration
		if !zero.IsZero() {
			t.Error("zero.IsZero() = false, want true")
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()

		if d.IsEmpty() {
			t.Error("IsEmpty() = true, want false for non-zero duration")
		}

		var zero Duration
		if !zero.IsEmpty() {
			t.Error("zero.IsEmpty() = false, want true")
		}
	})

	t.Run("Milliseconds", func(t *testing.T) {
		t.Parallel()

		if d.Milliseconds() != 9000000 {
			t.Errorf("Milliseconds() = %d, want %d", d.Milliseconds(), 9000000)
		}
	})

	t.Run("Seconds", func(t *testing.T) {
		t.Parallel()

		got := d.Seconds()
		if got < 8999.999 || got > 9000.001 {
			t.Errorf("Seconds() = %f, want approximately 9000", got)
		}
	})
}

func TestDuration_MarshalUnmarshal(t *testing.T) {
	t.Parallel()

	type config struct {
		Timeout Duration `json:"timeout"`
	}

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

		var c config

		err := json.Unmarshal([]byte(`{"timeout":"1h"}`), &c)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if c.Timeout.Duration() != time.Hour {
			t.Errorf("unmarshaled Timeout = %v, want %v", c.Timeout.Duration(), time.Hour)
		}
	})

	runUnmarshalErrorTest[config](t, "unmarshal invalid", `{"timeout":"invalid"}`)
}
