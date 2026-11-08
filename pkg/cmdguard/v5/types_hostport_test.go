package v5_test

import (
	"testing"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

func TestHostPort(t *testing.T) {
	t.Parallel()

	t.Run("ParseHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp, err := v5.ParseHostPort("localhost:8080")
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

		hp, err := v5.ParseHostPort(":8080")
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
					func() (v5.HostPort, error) { return v5.ParseHostPort(tt.input) },
					"host:port",
				)
			})
		}
	})

	t.Run("NewHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp, err := v5.NewHostPort("example.com", "443")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if hp.Host() != "example.com" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "example.com")
		}

		testHostPortPortInt(t, hp, 443)
	})

	t.Run("ParseHostPort valid", func(t *testing.T) {
		t.Parallel()

		hp, err := v5.ParseHostPort("127.0.0.1:3000")
		if err != nil {
			t.Fatal(err)
		}
		if hp.Host() != "127.0.0.1" {
			t.Errorf("got %q, want %q", hp.Host(), "127.0.0.1")
		}
		if hp.Port().Int() != 3000 {
			t.Errorf("got %d, want %d", hp.Port().Int(), 3000)
		}
	})

	t.Run("HostPort IsEmpty", func(t *testing.T) {
		t.Parallel()

		hp, _ := v5.ParseHostPort("localhost:8080")
		if hp.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("HostPort MarshalText", func(t *testing.T) {
		t.Parallel()

		hp, _ := v5.ParseHostPort("example.com:443")

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

		var hp v5.HostPort

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
