package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
)

// ConfigFileLoader loads configuration from a file.
// Implementations read raw bytes and populate a config struct,
// returning the list of struct field names that were explicitly set.
type ConfigFileLoader interface {
	Load(data []byte, cfg any) (setFields []string, err error)
}

// jsonLoader loads configuration from JSON files.
// Supports flat key-value objects where keys match flag tag names.
type jsonLoader struct{}

// Load unmarshals JSON data into cfg and returns the list of fields that were set.
func (l *jsonLoader) Load(data []byte, cfg any) (setFields []string, err error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileParse, err)
	}

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", ErrConfigFileParse, err)
	}

	setFields = make([]string, 0, len(raw))

	for _, tag := range tags {
		if _, ok := raw[tag.Name]; ok {
			setFields = append(setFields, tag.Field)
		}
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileParse, err)
	}

	return setFields, nil
}

// loadConfigFile tries to load a config file from the given paths.
// Paths are expanded via expandConfigPath. Missing files are skipped.
// Returns ErrConfigFileNotFound if none of the paths exist.
func loadConfigFile(paths []string, loader ConfigFileLoader, cfg any) (setFields []string, err error) {
	for _, path := range paths {
		path = expandConfigPath(path)

		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, fmt.Errorf("%w: reading %q: %w", ErrConfigFileRead, path, err)
		}

		setFields, err := loader.Load(data, cfg)
		if err != nil {
			return nil, fmt.Errorf("%w: loading %q: %w", ErrConfigFileLoad, path, err)
		}

		return setFields, nil
	}

	return nil, ErrConfigFileNotFound
}

// expandConfigPath expands environment variables and the leading ~ in a path.
func expandConfigPath(path string) string {
	path = os.ExpandEnv(path)

	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return path
}

// resolveConfigFlag checks os.Args for a --config flag override.
// Returns the flag value if found, otherwise an empty string.
func resolveConfigFlag(flagName string) string {
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.String(flagName, "", "")

	_ = fs.Parse(os.Args[1:])

	if path, _ := fs.GetString(flagName); path != "" {
		return path
	}

	return ""
}

// updateTagDefaultsFromConfig updates tag defaults for fields set by a config file.
// This makes config file values act as defaults that flags and env vars can override.
func (r *FlagRegistry) updateTagDefaultsFromConfig(cfg any, setFields []string) {
	if len(setFields) == 0 {
		return
	}

	setSet := make(map[string]struct{}, len(setFields))
	for _, f := range setFields {
		setSet[f] = struct{}{}
	}

	v, err := derefPointerToStruct(cfg)
	if err != nil {
		return
	}

	for i := range r.tags {
		if _, ok := setSet[r.tags[i].Field]; !ok {
			continue
		}

		field := v.FieldByName(r.tags[i].Field)
		if !field.IsValid() {
			continue
		}

		if s, ok := getFieldValue(field); ok {
			r.tags[i].Default = s
		}
	}
}

// WithConfigFile adds config file loading with the given search paths.
// Paths are tried in order; the first existing file wins.
// Environment variables and ~ are expanded in paths.
// Only JSON files are supported in the core package;
// use WithConfigFileLoader for YAML/TOML support.
func WithConfigFile[T any](paths ...string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.configFilePaths = paths
		cli.configFileLoader = &jsonLoader{}
	}
}

// WithConfigFileLoader adds config file loading with a custom loader.
// Paths are tried in order; the first existing file wins.
func WithConfigFileLoader[T any](loader ConfigFileLoader, paths ...string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.configFilePaths = paths
		cli.configFileLoader = loader
	}
}

// loadConfigFileOrSkip attempts to load a config file, returning nil on "not found".
// This is the helper used during CLI initialization.
func (cli *CLI[T]) loadConfigFileOrSkip() ([]string, error) {
	if cli.configFileLoader == nil || len(cli.configFilePaths) == 0 {
		return nil, nil
	}

	paths := cli.configFilePaths

	// Check for --config flag override.
	if override := resolveConfigFlag("config"); override != "" {
		paths = []string{override}
	}

	setFields, err := loadConfigFile(paths, cli.configFileLoader, cli.config)
	if err != nil && errors.Is(err, ErrConfigFileNotFound) {
		return nil, nil
	}

	return setFields, err
}
