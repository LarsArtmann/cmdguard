package v2

import (
	"fmt"
	"net/mail"
	"strings"
)

// emailPartsCount is the expected number of parts when splitting an email at "@".
const emailPartsCount = 2

// Email wraps a validated email address.
// Use this for config fields that must be valid email addresses.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type Email struct {
	address string
}

// ParseEmail creates a new Email from a string.
// Returns an error if the string is not a valid email address.
func ParseEmail(s string) (Email, error) {
	if strings.TrimSpace(s) == "" {
		return Email{}, fmt.Errorf("%w: email cannot be empty", ErrInvalidEmail)
	}

	addr, err := mail.ParseAddress(s)
	if err != nil {
		return Email{}, fmt.Errorf("%w: %w", ErrInvalidEmail, err)
	}

	return Email{address: addr.Address}, nil
}

// MustParseEmail creates an Email from a string, panicking if invalid.
// Use only when you know the email is valid (e.g., for constants).
func MustParseEmail(s string) Email {
	return MustParse("MustParseEmail", s, ParseEmail)
}

// String returns the email address as a string.
func (e Email) String() string {
	return e.address
}

// Address returns the email address.
func (e Email) Address() string {
	return e.address
}

// IsEmpty returns true if the email has no value.
func (e Email) IsEmpty() bool {
	return e.address == ""
}

// Local returns the local part of the email (before @).
func (e Email) Local() string {
	return e.emailPart(0)
}

// Domain returns the domain part of the email (after @).
func (e Email) Domain() string {
	return e.emailPart(1)
}

// emailPart returns the email part at the given index (0=local, 1=domain).
func (e Email) emailPart(index int) string {
	parts := strings.Split(e.address, "@")
	if len(parts) != emailPartsCount {
		return ""
	}

	return parts[index]
}

// MarshalText implements encoding.TextMarshaler for Email.
func (e Email) MarshalText() ([]byte, error) {
	return []byte(e.address), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Email.
func (e *Email) UnmarshalText(text []byte) error {
	parsed, err := ParseEmail(string(text))
	if err != nil {
		return err
	}

	*e = parsed

	return nil
}
