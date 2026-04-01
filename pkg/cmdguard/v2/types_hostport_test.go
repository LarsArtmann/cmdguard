package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestHostPort(t *testing.T) {
	t.Parallel()

	t.Run("ParseHostPort valid", func(t *testing.T) {
		t.Parallel()
		hp, err := v2.ParseHostPort("localhost:8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hp.String() != "localhost:8080" {
			t.Errorf("String() = %q, want %q", hp.String(), "localhost:8080")
		}
		if hp.Host() != "localhost" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "localhost")
		}
		if hp.Port().Int() != 8080 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 8080)
		}
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

	t.Run("ParseHostPort empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseHostPort("")
		if err == nil {
			t.Fatal("expected error for empty host:port")
		}
	})

	t.Run("ParseHostPort invalid", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseHostPort("not-valid")
		if err == nil {
			t.Fatal("expected error for invalid host:port")
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
		if hp.Port().Int() != 443 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 443)
		}
	})

	t.Run("MustParseHostPort valid", func(t *testing.T) {
		t.Parallel()
		hp := v2.MustParseHostPort("127.0.0.1:3000")
		if hp.Host() != "127.0.0.1" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "127.0.0.1")
		}
		if hp.Port().Int() != 3000 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 3000)
		}
	})

	t.Run("MustParseHostPort panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid host:port")
			}
		}()
		v2.MustParseHostPort("invalid")
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
		if hp.Port().Int() != 9090 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 9090)
		}
	})
}
