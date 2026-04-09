package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestPort(t *testing.T) {
	t.Parallel()

	t.Run("ParsePort numeric", func(t *testing.T) {
		t.Parallel()
		p, err := v2.ParsePort("8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}
		if p.String() != "8080" {
			t.Errorf("String() = %q, want %q", p.String(), "8080")
		}
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
				p, err := v2.ParsePort(tt.input)
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
				testParseError(t, func() (v2.Port, error) { return v2.ParsePort(tt.input) }, "port")
			})
		}
	})

	t.Run("PortFromInt valid", func(t *testing.T) {
		t.Parallel()
		p, err := v2.PortFromInt(8080)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}
	})

	t.Run("PortFromInt invalid", func(t *testing.T) {
		t.Parallel()
		testParseError(t, func() (v2.Port, error) { return v2.PortFromInt(70000) }, "port")
	})

	t.Run("Port IsValid", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
		if !p.IsValid() {
			t.Error("IsValid() = false, want true")
		}
	})

	t.Run("Port IsWellKnown", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("80")
		if !p.IsWellKnown() {
			t.Error("IsWellKnown() = false, want true for port 80")
		}
		p2, _ := v2.ParsePort("8080")
		if p2.IsWellKnown() {
			t.Error("IsWellKnown() = true, want false for port 8080")
		}
	})

	t.Run("Port IsRegistered", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
		if !p.IsRegistered() {
			t.Error("IsRegistered() = false, want true for port 8080")
		}
	})

	t.Run("Port IsDynamic", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("50000")
		if !p.IsDynamic() {
			t.Error("IsDynamic() = false, want true for port 50000")
		}
	})

	t.Run("MustParsePort valid", func(t *testing.T) {
		t.Parallel()
		p := v2.MustParsePort("443")
		if p.Int() != 443 {
			t.Errorf("Int() = %d, want %d", p.Int(), 443)
		}
	})

	t.Run("MustParsePort panic", func(t *testing.T) {
		t.Parallel()
		testMustParsePanics(t, v2.MustParsePort, "port")
	})

	t.Run("Port MarshalText", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
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
		var p v2.Port
		err := p.UnmarshalText([]byte("9090"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 9090 {
			t.Errorf("Int() = %d, want %d", p.Int(), 9090)
		}
	})
}
