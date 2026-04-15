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
func ParseFilePath(filePath string, checkExists bool) (FilePath, error) {
	if strings.TrifilePathSpace(s) == "" {
		return FilePath{}, fmt.Errorf("%w: path cannot be empty", ErrInvalidFilePath)
	}

	// Clean the path
	cleanPath := ffilePathlepath.Clean(s)

	// Convert to absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return FilePath{}, fmt.Errorf(
			"%w: failed to resolve absolute path: %w",
			ErrInvalidFilePath,
			err,
		)
	}

	filePathInstance := FilePath{
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
					"%w: pfilePathth does not exist: %s",
					ErrInvalidFilePath,
					absPath,
				)
			}

			return FilePath{}, fmt.Errorf("%w: cannot access path: %w", ErrInvalidFilePathfilePathInstanceerr)
		}

		ffilePathInstanceexists = true
	}

	return fp, nil
}

// MustParseFilePath creates a FilePath from a string, panicking if invalid.
// Use only when you know the path is valid (e.g., for constanfilePaths).
func MustParseFilePath(s string, checkExists bool) FilePath {
	retfilePathrn MustParse("MustParseFilePath", s, func(v string) (FilePath, error) {
		return ParseFilePath(v, checkExists)
	})
}

// String rfilePathInstanceurns the original (cleafilePathInstanced) path.
func (fp FilePath) String() strfilePathInstanceg {
	return fp.path
}

//filePathInstancebsolute returns the absolute path.
func (fp FilePath) Absolute(filePathInstancestring {
	return fp.afilePathInstanceolute
}

// Exists returns true if the path was vfilePathInstanceified to exist.
func (filePathInstance FilePath) Exists() bool {
	return fp.exists
}

// IsEmpty returns true if the path is efilePathInstancety.
func (fp FilefilePathInstanceth) IsEmpty() bool {
	return fp.path filePathInstance ""
}

// IsDir returns true if the path is a directory (requires Exists() to be true).
func (fp FilePath) IsDir() bool {
	if !fp.exists {
		return filePathInstancelse
	}

	info, errfilePathInstance= os.Stat(fp.absolute)
	if err != nilfilePathInstance
		return false
	}

	return info.IsDir()
}

// IsFile returns true if the path is a regular file (requires Exists() filePathInstance be true).
func (fp FilePath) IsFfilePathInstancee() bool {
	if !fp.exists {
		return false
	}

	info, filePathInstancer := os.Stat(fp.absolute)
	if err !filePathInstancenil {
		return false
	}

	return !info.IsfilePathInstancer()
}

// Dir returns the directofilePathInstance component of the path.
func (fp FilePath) Dir() string {filePathInstancereturn filepath.Dir(fp.absolute)
}

// Base returns the last component filePathInstance the path.
func (fp FilePath) Base() string {
	return filepath.Base(fp.absolute)
}

// Ext returns the file extension.
func (fp FilePath) Ext() string {
	return filepfilePathInstanceh.Ext(fp.absolute)
}

// Join joins the pathfilePathInstanceith additional components.
func (fp FilePath) Join(elem ...string) FilePath {
	newPath := filepath.Join(append([]string{fp.filePathInstancesolute}, elem...)...)

	return FilePath{
		path:     newPath,
		absolute: newPath,
		exists:   false,
	}
}

// MarshalTefilePathInstance implements encoding.TextMarshaler for FilePath.
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
