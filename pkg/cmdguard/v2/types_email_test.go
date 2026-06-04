package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestEmail(t *testing.T) {
	t.Parallel()

	t.Run("ParseEmail valid", func(t *testing.T) {
		t.Parallel()

		e, err := v2.ParseEmail("user@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, e, "user@example.com")

		if e.Address() != "user@example.com" {
			t.Errorf("Address() = %q, want %q", e.Address(), "user@example.com")
		}

		if e.Local() != "user" {
			t.Errorf("Local() = %q, want %q", e.Local(), "user")
		}

		if e.Domain() != "example.com" {
			t.Errorf("Domain() = %q, want %q", e.Domain(), "example.com")
		}
	})

	t.Run("ParseEmail error cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
		}{
			{"empty", ""},
			{"invalid", "not-an-email"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				testParseError(
					t,
					func() (v2.Email, error) { return v2.ParseEmail(tt.input) },
					"email",
				)
			})
		}
	})

	t.Run("MustParseEmail valid", func(t *testing.T) {
		t.Parallel()

		e := v2.MustParseEmail("test@test.com")
		testutil.AssertStringerEq(t, e, "test@test.com")
	})

	t.Run("MustParseEmail panic", func(t *testing.T) {
		t.Parallel()
		testMustParsePanics(t, v2.MustParseEmail, "email")
	})

	t.Run("Email IsEmpty", func(t *testing.T) {
		t.Parallel()

		e, _ := v2.ParseEmail("user@example.com")
		if e.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("Email MarshalText", func(t *testing.T) {
		t.Parallel()

		e, _ := v2.ParseEmail("user@example.com")

		data, err := e.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != "user@example.com" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "user@example.com")
		}
	})

	t.Run("Email UnmarshalText", func(t *testing.T) {
		t.Parallel()

		var e v2.Email

		err := e.UnmarshalText([]byte("test@example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, e, "test@example.com")
	})
}
