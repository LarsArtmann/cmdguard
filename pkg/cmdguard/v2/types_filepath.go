package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
		return FilePath{}, fmt.Errorf(
			"%w: failed to resolve absolute path: %w",
			ErrInvalidFilePath,
			err,
		)
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
				return FilePath{}, fmt.Errorf(
					"%w: path does not exist: %s",
					ErrInvalidFilePath,
					absPath,
				)
			}

			return FilePath{}, fmt.Errorf("%w: cannot access path: %w", ErrInvalidFilePath, err)
		}

		fp.exists = true
	}

	return fp, nil
}

// MustParseFilePath creates a FilePath from a string, panicking if invalid.
// Use only when you know the path is valid (e.g., for constants).
func MustParseFilePath(s string, checkExists bool) FilePath {
	return MustParse("MustParseFilePath", s, func(v string) (FilePath, error) {
		return ParseFilePath(v, checkExists)
	})
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
	return textMarshal(fp, func(fp FilePath) string { return fp.path })
}

func (fp *FilePath) UnmarshalText(text []byte) error {
	return textUnmarshal(
		fp,
		text,
		func(s string) (FilePath, error) { return ParseFilePath(s, false) },
	)
}
