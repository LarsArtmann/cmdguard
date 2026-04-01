package v2

import (
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	if strings.TrimSpace(s) == "" {
		return URL{}, fmt.Errorf("%w: URL cannot be empty", ErrInvalidURL)
	}

	u, err := url.Parse(s)
	if err != nil {
		return URL{}, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if u.Scheme == "" {
		return URL{}, fmt.Errorf("%w: URL must have a scheme (e.g., http, https)", ErrInvalidURL)
	}

	if u.Host == "" {
		return URL{}, fmt.Errorf("%w: URL must have a host", ErrInvalidURL)
	}

	return URL{url: u}, nil
}

// MustParseURL creates a URL from a string, panicking if invalid.
// Use only when you know the URL is valid (e.g., for constants).
func MustParseURL(s string) URL {
	u, err := ParseURL(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseURL(%q): %v", s, err))
	}

	return u
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
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for URL.
func (u *URL) UnmarshalText(text []byte) error {
	parsed, err := ParseURL(string(text))
	if err != nil {
		return err
	}

	*u = parsed

	return nil
}

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
		return Email{}, fmt.Errorf("%w: %v", ErrInvalidEmail, err)
	}

	return Email{address: addr.Address}, nil
}

// MustParseEmail creates an Email from a string, panicking if invalid.
// Use only when you know the email is valid (e.g., for constants).
func MustParseEmail(s string) Email {
	e, err := ParseEmail(s)
	if err != nil {
		panic(fmt.Sprintf("MustParseEmail(%q): %v", s, err))
	}

	return e
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
	parts := strings.Split(e.address, "@")
	if len(parts) != 2 {
		return ""
	}

	return parts[0]
}

// Domain returns the domain part of the email (after @).
func (e Email) Domain() string {
	parts := strings.Split(e.address, "@")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
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

// FromInt creates a Port from an int.
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

// FilePath represents a validated file system path.
// Use this for config fields that specify file or directory paths.
// Provides validation, existence checks, and path cleaning.
//
//nolint:recvcheck // MarshalText/UnmarshalText require different receivers per Go convention
type FilePath struct {
	path     string
	absolute string
	exists   bool
}

// ParseFilePath creates a new FilePath from a string.
// The path is cleaned (removes .. and . components) and converted to absolute.
// Set checkExists to true to verify the path exists on the filesystem.
func ParseFilePath(s string, checkExists bool) (FilePath, error) {
	if strings.TrimSpace(s) == "" {
		return FilePath{}, fmt.Errorf("%w: path cannot be empty", ErrInvalidFilePath)
	}

	// Clean the path
	cleanPath := filepath.Clean(s)

	// Convert to absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return FilePath{}, fmt.Errorf("%w: failed to resolve absolute path: %v", ErrInvalidFilePath, err)
	}

	fp := FilePath{
		path:     cleanPath,
		absolute: absPath,
		exists:   false,
	}

	// Check existence if requested
	if checkExists {
		_, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				return FilePath{}, fmt.Errorf("%w: path does not exist: %s", ErrInvalidFilePath, absPath)
			}

			return FilePath{}, fmt.Errorf("%w: cannot access path: %v", ErrInvalidFilePath, err)
		}

		fp.exists = true
	}

	return fp, nil
}

// MustParseFilePath creates a FilePath from a string, panicking if invalid.
// Use only when you know the path is valid (e.g., for constants).
func MustParseFilePath(s string, checkExists bool) FilePath {
	fp, err := ParseFilePath(s, checkExists)
	if err != nil {
		panic(fmt.Sprintf("MustParseFilePath(%q): %v", s, err))
	}

	return fp
}

// String returns the original (cleaned) path.
func (fp FilePath) String() string {
	return fp.path
}

// Absolute returns the absolute path.
func (fp FilePath) Absolute() string {
	return fp.absolute
}

// Exists returns true if the path was verified to exist.
func (fp FilePath) Exists() bool {
	return fp.exists
}

// IsEmpty returns true if the path is empty.
func (fp FilePath) IsEmpty() bool {
	return fp.path == ""
}

// IsDir returns true if the path is a directory (requires Exists() to be true).
func (fp FilePath) IsDir() bool {
	if !fp.exists {
		return false
	}

	info, err := os.Stat(fp.absolute)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// IsFile returns true if the path is a regular file (requires Exists() to be true).
func (fp FilePath) IsFile() bool {
	if !fp.exists {
		return false
	}

	info, err := os.Stat(fp.absolute)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// Dir returns the directory component of the path.
func (fp FilePath) Dir() string {
	return filepath.Dir(fp.absolute)
}

// Base returns the last component of the path.
func (fp FilePath) Base() string {
	return filepath.Base(fp.absolute)
}

// Ext returns the file extension.
func (fp FilePath) Ext() string {
	return filepath.Ext(fp.absolute)
}

// Join joins the path with additional components.
func (fp FilePath) Join(elem ...string) FilePath {
	newPath := filepath.Join(append([]string{fp.absolute}, elem...)...)

	return FilePath{
		path:     newPath,
		absolute: newPath,
		exists:   false,
	}
}

// MarshalText implements encoding.TextMarshaler for FilePath.
func (fp FilePath) MarshalText() ([]byte, error) {
	return []byte(fp.path), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for FilePath.
// Note: This does NOT check if the path exists.
func (fp *FilePath) UnmarshalText(text []byte) error {
	parsed, err := ParseFilePath(string(text), false)
	if err != nil {
		return err
	}

	*fp = parsed

	return nil
}

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
