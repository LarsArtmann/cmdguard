package v2

import (
	"fmt"
	"net/url"
)

// URL wraps url.URL with parsing validation and text marshaling.
// Use this for config fields that must be valid URLs.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type URL struct {
	url *url.URL
}

// ParseURL creates a new URL from a string.
// Returns an error if the string is not a valid URL.
func ParseURL(s string) (URL, error) {
	if err := requireNonEmpty(s, "URL", ErrInvalidURL); err != nil {
		return URL{}, err
	}

	u, err := url.Parse(s)
	if err != nil {
		return URL{}, fmt.Errorf("%w: %w", ErrInvalidURL, err)
	}

	if u.Scheme == "" {
		return URL{}, fmt.Errorf("%w: URL must have a scheme (e.g., http, https)", ErrInvalidURL)
	}

	if u.Host == "" {
		return URL{}, fmt.Errorf("%w: URL must have a host", ErrInvalidURL)
	}

	return URL{url: u}, nil
}

// URL returns the underlying *url.URL.
// Returns nil if the URL was not properly initialized.
func (u URL) URL() *url.URL {
	if u.url == nil {
		return nil
	}

	// Return a copy to prevent mutation
	cpy := *u.url

	return &cpy
}

// String returns the URL as a string.
func (u URL) String() string {
	if u.url == nil {
		return ""
	}

	return u.url.String()
}

// IsEmpty returns true if the URL has no value.
func (u URL) IsEmpty() bool {
	return u.url == nil || u.url.String() == ""
}

// Scheme returns the URL scheme (e.g., "https").
func (u URL) Scheme() string {
	if u.url == nil {
		return ""
	}

	return u.url.Scheme
}

// Host returns the URL host (e.g., "example.com:8080").
func (u URL) Host() string {
	if u.url == nil {
		return ""
	}

	return u.url.Host
}

// Hostname returns the URL hostname without port (e.g., "example.com").
func (u URL) Hostname() string {
	if u.url == nil {
		return ""
	}

	return u.url.Hostname()
}

// Port returns the URL port (e.g., "8080").
func (u URL) Port() string {
	if u.url == nil {
		return ""
	}

	return u.url.Port()
}

// Path returns the URL path.
func (u URL) Path() string {
	if u.url == nil {
		return ""
	}

	return u.url.Path
}

// MarshalText implements encoding.TextMarshaler for URL.
func (u URL) MarshalText() ([]byte, error) {
	return textMarshal(u, URL.String)
}

func (u *URL) UnmarshalText(text []byte) error {
	return textUnmarshal(u, text, ParseURL)
}
