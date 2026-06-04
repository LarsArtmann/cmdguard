package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestHostPort(t *testing.T) {
	t.Parallel()

	t.Run("ParseHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp, err := v2.ParseHostPort("localhost:8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, hp, "localhost:8080")

		if hp.Host() != "localhost" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "localhost")
		}

		testHostPortPortInt(t, hp, 8080)
	})

	t.Run("ParseHostPort any host", func(t *testing.T) {
		t.Parallel()

		hp, err := v2.ParseHostPort(":8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if hp.Host() != "" {
			t.Errorf("Host() = %q, want empty", hp.Host())
		}

		if !hp.IsAnyHost() {
			t.Error("IsAnyHost() = false, want true")
		}
	})

	t.Run("ParseHostPort error cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
		}{
			{"empty", ""},
			{"invalid", "not-valid"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				testParseError(
					t,
					func() (v2.HostPort, error) { return v2.ParseHostPort(tt.input) },
					"host:port",
				)
			})
		}
	})

	t.Run("NewHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp, err := v2.NewHostPort("example.com", "443")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if hp.Host() != "example.com" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "example.com")
		}

		testHostPortPortInt(t, hp, 443)
	})

	t.Run("MustParseHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp := v2.MustParseHostPort("127.0.0.1:3000")
		if hp.Host() != "127.0.0.1" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "127.0.0.1")
		}

		testHostPortPortInt(t, hp, 3000)
	})

	t.Run("MustParseHostPort panic", func(t *testing.T) {
		t.Parallel()
		testMustParsePanics(t, v2.MustParseHostPort, "host:port")
	})

	t.Run("HostPort IsEmpty", func(t *testing.T) {
		t.Parallel()

		hp, _ := v2.ParseHostPort("localhost:8080")
		if hp.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("HostPort MarshalText", func(t *testing.T) {
		t.Parallel()

		hp, _ := v2.ParseHostPort("example.com:443")

		data, err := hp.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != "example.com:443" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "example.com:443")
		}
	})

	t.Run("HostPort UnmarshalText", func(t *testing.T) {
		t.Parallel()

		var hp v2.HostPort

		err := hp.UnmarshalText([]byte("localhost:9090"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if hp.Host() != "localhost" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "localhost")
		}

		testHostPortPortInt(t, hp, 9090)
	})
}
