package v4

import (
	"fmt"
	"strconv"
	"strings"
)

// Well-known port constants for named port lookups.
const (
	portHTTP  = 80
	portHTTPS = 443
	portSSH   = 22
	portFTP   = 21
	portDNS   = 53
	portSMTP  = 25

	portMaxValue      = 65535
	portWellKnownMax  = 1023
	portRegisteredMin = 1024
	portRegisteredMax = 49151
	portDynamicMin    = 49152
)

// validatePortRange returns an error if port is outside 1-65535.
func validatePortRange(port int) error {
	if port < 1 || port > portMaxValue {
		return fmt.Errorf("%w: port %d is out of range (1-%d)", ErrInvalidPort, port, portMaxValue)
	}

	return nil
}

// Port represents a valid network port number (1-65535).
// Use this for config fields that specify TCP/UDP ports.
//

type Port struct {
	port uint16
}

// ParsePort creates a new Port from a string.
// Accepts numeric strings (e.g., "8080") or named ports (e.g., "http", "https").
// Returns an error if the port is not valid.
func ParsePort(s string) (Port, error) {
	err := requireNonEmpty(s, "port", ErrInvalidPort)
	if err != nil {
		return Port{}, err
	}

	// Check for named ports
	switch strings.ToLower(s) {
	case "http":
		return Port{port: portHTTP}, nil
	case "https":
		return Port{port: portHTTPS}, nil
	case "ssh":
		return Port{port: portSSH}, nil
	case "ftp":
		return Port{port: portFTP}, nil
	case "dns":
		return Port{port: portDNS}, nil
	case "smtp":
		return Port{port: portSMTP}, nil
	}

	// Parse numeric port
	port, err := strconv.Atoi(s)
	if err != nil {
		return Port{}, fmt.Errorf("%w: %q is not a valid port number", ErrInvalidPort, s)
	}

	err = validatePortRange(port)
	if err != nil {
		return Port{}, err
	}

	return Port{port: uint16(port)}, nil
}

// PortFromInt creates a Port from an int.
// Returns an error if the port is out of valid range.
func PortFromInt(port int) (Port, error) {
	err := validatePortRange(port)
	if err != nil {
		return Port{}, fmt.Errorf("port=%d: %w", port, err)
	}

	return Port{port: uint16(port)}, nil
}

// Int returns the port as an int.
func (p Port) Int() int {
	return int(p.port)
}

// String returns the port as a string.
func (p Port) String() string {
	return strconv.Itoa(int(p.port))
}

// IsValid returns true if the port is in the valid range.
func (p Port) IsValid() bool {
	return p.port >= 1
}

// IsEmpty returns true if the port is zero (unset).
func (p Port) IsEmpty() bool {
	return p.port == 0
}

// IsWellKnown returns true if the port is in the well-known range (1-1023).
func (p Port) IsWellKnown() bool {
	return p.port >= 1 && p.port <= portWellKnownMax
}

// IsRegistered returns true if the port is in the registered range (1024-49151).
func (p Port) IsRegistered() bool {
	return p.port >= portRegisteredMin && p.port <= portRegisteredMax
}

// IsDynamic returns true if the port is in the dynamic/private range (49152-65535).
func (p Port) IsDynamic() bool {
	return p.port >= portDynamicMin
}

// MarshalText implements encoding.TextMarshaler for Port.
func (p Port) MarshalText() ([]byte, error) {
	return textMarshal(p, Port.String)
}

func (p *Port) UnmarshalText(text []byte) error {
	return textUnmarshal(p, text, ParsePort)
}
