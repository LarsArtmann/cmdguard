package v2

import (
	"fmt"
	"net"
	"strings"
)

// HostPort combines a hostname and port for network addresses.
// Use this for config fields that specify network endpoints.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type HostPort struct {
	host string
	port Port
}

// ParseHostPort creates a new HostPort from a string.
// Accepts formats like "localhost:8080", "example.com:443", or ":8080" (any host).
func ParseHostPort(s string) (HostPort, error) {
	if strings.TrimSpace(s) == "" {
		return HostPort{}, fmt.Errorf("%w: host:port cannot be empty", ErrInvalidHostPort)
	}

	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		return HostPort{}, fmt.Errorf("%w: %v", ErrInvalidHostPort, err)
	}

	port, err := ParsePort(portStr)
	if err != nil {
		return HostPort{}, fmt.Errorf("%w: invalid port: %v", ErrInvalidHostPort, err)
	}

	return HostPort{host: host, port: port}, nil
}

// MustParseHostPort creates a HostPort from a string, panicking if invalid.
func MustParseHostPort(s string) HostPort {
	hp, err := ParseHostPort(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseHostPort(%q): %v", s, err))
	}

	return hp
}

// NewHostPort creates a HostPort from host and port strings.
func NewHostPort(host, portStr string) (HostPort, error) {
	port, err := ParsePort(portStr)
	if err != nil {
		return HostPort{}, err
	}

	return HostPort{host: host, port: port}, nil
}

// String returns the host:port as a string.
func (hp HostPort) String() string {
	if hp.host == "" {
		return ":" + hp.port.String()
	}

	return net.JoinHostPort(hp.host, hp.port.String())
}

// Host returns the host component.
func (hp HostPort) Host() string {
	return hp.host
}

// Port returns the port component.
func (hp HostPort) Port() Port {
	return hp.port
}

// IsEmpty returns true if both host and port are empty.
func (hp HostPort) IsEmpty() bool {
	return hp.host == "" && hp.port.port == 0
}

// IsAnyHost returns true if the host is empty (meaning "any" or "all interfaces").
func (hp HostPort) IsAnyHost() bool {
	return hp.host == ""
}

// MarshalText implements encoding.TextMarshaler for HostPort.
func (hp HostPort) MarshalText() ([]byte, error) {
	return []byte(hp.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for HostPort.
func (hp *HostPort) UnmarshalText(text []byte) error {
	parsed, err := ParseHostPort(string(text))
	if err != nil {
		return err
	}

	*hp = parsed

	return nil
}
