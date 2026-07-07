package v3_test

import (
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestPort(t *testing.T) {
	t.Parallel()

	t.Run("ParsePort numeric", func(t *testing.T) {
		t.Parallel()

		p, err := v3.ParsePort("8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}

		testutil.AssertStringerEq(t, p, "8080")
	})

	t.Run("ParsePort named", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			input    string
			expected int
		}{
			{"http", "http", 80},
			{"https", "https", 443},
			{"ssh", "ssh", 22},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				p, err := v3.ParsePort(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if p.Int() != tt.expected {
					t.Errorf("Int() = %d, want %d", p.Int(), tt.expected)
				}
			})
		}
	})

	t.Run("ParsePort error cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
		}{
			{"empty", ""},
			{"out of range", "70000"},
			{"zero", "0"},
			{"negative", "-1"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				testParseError(t, func() (v3.Port, error) { return v3.ParsePort(tt.input) }, "port")
			})
		}
	})

	t.Run("PortFromInt valid", func(t *testing.T) {
		t.Parallel()

		p, err := v3.PortFromInt(8080)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}
	})

	t.Run("PortFromInt invalid", func(t *testing.T) {
		t.Parallel()
		testParseError(t, func() (v3.Port, error) { return v3.PortFromInt(70000) }, "port")
	})

	t.Run("Port IsEmpty", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("8080")
		if p.IsEmpty() {
			t.Error("IsEmpty() = true, want false for parsed port")
		}

		var zero v3.Port
		if !zero.IsEmpty() {
			t.Error("zero Port IsEmpty() = false, want true")
		}
	})

	t.Run("Port IsValid", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("8080")
		if !p.IsValid() {
			t.Error("IsValid() = false, want true")
		}
	})

	t.Run("Port IsWellKnown", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("80")
		if !p.IsWellKnown() {
			t.Error("IsWellKnown() = false, want true for port 80")
		}

		p2, _ := v3.ParsePort("8080")
		if p2.IsWellKnown() {
			t.Error("IsWellKnown() = true, want false for port 8080")
		}
	})

	t.Run("Port IsRegistered", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("8080")
		if !p.IsRegistered() {
			t.Error("IsRegistered() = false, want true for port 8080")
		}
	})

	t.Run("Port IsDynamic", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("50000")
		if !p.IsDynamic() {
			t.Error("IsDynamic() = false, want true for port 50000")
		}
	})

	t.Run("ParsePort valid", func(t *testing.T) {
		t.Parallel()

		p, err := v3.ParsePort("443")
		if err != nil {
			t.Fatal(err)
		}
		if p.Int() != 443 {
			t.Errorf("got %d, want %d", p.Int(), 443)
		}
	})

	t.Run("Port MarshalText", func(t *testing.T) {
		t.Parallel()

		p, _ := v3.ParsePort("8080")

		data, err := p.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != "8080" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "8080")
		}
	})

	t.Run("Port UnmarshalText", func(t *testing.T) {
		t.Parallel()

		var p v3.Port

		err := p.UnmarshalText([]byte("9090"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if p.Int() != 9090 {
			t.Errorf("Int() = %d, want %d", p.Int(), 9090)
		}
	})
}
