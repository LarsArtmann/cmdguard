package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func TestURL(t *testing.T) {
	t.Parallel()

	t.Run("ParseURL valid", func(t *testing.T) {
		t.Parallel()

		u, err := v2.ParseURL("https://example.com:8080/path")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, u, "https://example.com:8080/path")

		if u.Scheme() != "https" {
			t.Errorf("Scheme() = %q, want %q", u.Scheme(), "https")
		}

		if u.Host() != "example.com:8080" {
			t.Errorf("Host() = %q, want %q", u.Host(), "example.com:8080")
		}

		if u.Hostname() != "example.com" {
			t.Errorf("Hostname() = %q, want %q", u.Hostname(), "example.com")
		}

		if u.Port() != "8080" {
			t.Errorf("Port() = %q, want %q", u.Port(), "8080")
		}

		if u.Path() != "/path" {
			t.Errorf("Path() = %q, want %q", u.Path(), "/path")
		}
	})

	t.Run("ParseURL error cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input string
		}{
			{"empty", ""},
			{"missing scheme", "example.com/path"},
			{"missing host", "http:///path"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				testParseError(t, func() (v2.URL, error) { return v2.ParseURL(tt.input) }, "URL")
			})
		}
	})

	t.Run("MustParseURL valid", func(t *testing.T) {
		t.Parallel()

		u := v2.MustParseURL("https://example.com")
		if u.String() != "https://example.com" {
			t.Errorf("String() = %q, want %q", u.String(), "https://example.com")
		}
	})

	t.Run("MustParseURL panic", func(t *testing.T) {
		t.Parallel()
		testMustParsePanics(t, v2.MustParseURL, "URL")
	})

	t.Run("URL IsEmpty", func(t *testing.T) {
		t.Parallel()

		u, _ := v2.ParseURL("https://example.com")
		if u.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("URL returns copy", func(t *testing.T) {
		t.Parallel()

		u, _ := v2.ParseURL("https://example.com")
		url1 := u.URL()

		url2 := u.URL()
		if url1 == url2 {
			t.Error("URL() should return different pointers")
		}
	})

	t.Run("URL MarshalText", func(t *testing.T) {
		t.Parallel()

		u, _ := v2.ParseURL("https://example.com/path")

		data, err := u.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(data) != "https://example.com/path" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "https://example.com/path")
		}
	})

	t.Run("URL UnmarshalText", func(t *testing.T) {
		t.Parallel()

		var u v2.URL

		err := u.UnmarshalText([]byte("https://example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, u, "https://example.com")
	})
}
