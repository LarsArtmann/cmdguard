package v2

import "errors"

// Config-related sentinel errors.
var (
	// ErrConfigValidation indicates config validation failed.
	ErrConfigValidation = errors.New("config validation failed")

	// ErrConfigNil indicates a nil config was passed where a config struct is required.
	ErrConfigNil = errors.New("config must not be nil")

	// ErrConfigNotPointer indicates config is not a pointer to struct.
	ErrConfigNotPointer = errors.New("config must be a pointer to struct")

	// ErrConfigFileRead indicates a config file could not be read from disk.
	ErrConfigFileRead = errors.New("failed to read config file")

	// ErrConfigFileParse indicates a config file could not be parsed.
	ErrConfigFileParse = errors.New("failed to parse config file")

	// ErrConfigFileLoad indicates config file loading failed.
	ErrConfigFileLoad = errors.New("failed to load config file")

	// ErrConfigFileNotFound indicates no config file was found in any search path.
	ErrConfigFileNotFound = errors.New("config file not found")
)
