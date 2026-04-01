package v2

import (
	"fmt"
	"strconv"
	"strings"
)

// Port represents a valid network port number (1-65535).
// Use this for config fields that specify TCP/UDP ports.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type Port struct {
	port int
}

// ParsePort creates a new Port from a string.
// Accepts numeric strings (e.g., "8080") or named ports (e.g., "http", "https").
// Returns an error if the port is not valid.
func ParsePort(s string) (Port, error) {
	if strings.TrimSpace(s) == "" {
		return Port{}, fmt.Errorf("%w: port cannot be empty", ErrInvalidPort)
	}

	// Check for named ports
	switch strings.ToLower(s) {
	case "http":
		return Port{port: 80}, nil
	case "https":
		return Port{port: 443}, nil
	case "ssh":
		return Port{port: 22}, nil
	case "ftp":
		return Port{port: 21}, nil
	case "dns":
		return Port{port: 53}, nil
	case "smtp":
		return Port{port: 25}, nil
	}

	// Parse numeric port
	port, err := strconv.Atoi(s)
	if err != nil {
		return Port{}, fmt.Errorf("%w: %q is not a valid port number", ErrInvalidPort, s)
	}

	if port < 1 || port > 65535 {
		return Port{}, fmt.Errorf("%w: port %d is out of range (1-65535)", ErrInvalidPort, port)
	}

	return Port{port: port}, nil
}

// MustParsePort creates a Port from a string, panicking if invalid.
// Use only when you know the port is valid (e.g., for constants).
func MustParsePort(s string) Port {
	p, err := ParsePort(s)
	if err != nil {
		panic(fmt.Sprintf("MustParsePort(%q): %v", s, err))
	}

	return p
}

// PortFromInt creates a Port from an int.
// Returns an error if the port is out of valid range.
func PortFromInt(port int) (Port, error) {
	if port < 1 || port > 65535 {
		return Port{}, fmt.Errorf("%w: port %d is out of range (1-65535)", ErrInvalidPort, port)
	}

	return Port{port: port}, nil
}

// Int returns the port as an int.
func (p Port) Int() int {
	return p.port
}

// String returns the port as a string.
func (p Port) String() string {
	return strconv.Itoa(p.port)
}

// IsValid returns true if the port is in the valid range.
func (p Port) IsValid() bool {
	return p.port >= 1 && p.port <= 65535
}

// IsWellKnown returns true if the port is in the well-known range (1-1023).
func (p Port) IsWellKnown() bool {
	return p.port >= 1 && p.port <= 1023
}

// IsRegistered returns true if the port is in the registered range (1024-49151).
func (p Port) IsRegistered() bool {
	return p.port >= 1024 && p.port <= 49151
}

// IsDynamic returns true if the port is in the dynamic/private range (49152-65535).
func (p Port) IsDynamic() bool {
	return p.port >= 49152 && p.port <= 65535
}

// MarshalText implements encoding.TextMarshaler for Port.
func (p Port) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for Port.
func (p *Port) UnmarshalText(text []byte) error {
	parsed, err := ParsePort(string(text))
	if err != nil {
		return err
	}

	*p = parsed

	return nil
}
