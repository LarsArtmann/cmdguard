package v3

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
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

// NewJSONLoader returns a ConfigFileLoader for JSON files.
// Supports flat key-value objects and nested objects where keys match flag tag names.
func NewJSONLoader() ConfigFileLoader {
	return &jsonLoader{}
}

// Load unmarshals JSON data into cfg and returns the list of fields that were set.
func (l *jsonLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]json.RawMessage

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileParse, err)
	}

	tags, err := ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", ErrConfigFileParse, err)
	}

	present := make(map[string]bool, len(raw))
	collectKeysRecursive(raw, present)

	setFields := FilterSetFields(tags, present)

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileParse, err)
	}

	return setFields, nil
}

// loadConfigFile tries to load a config file from the given paths.
// Paths are expanded via expandConfigPath. Missing files are skipped.
// Returns ErrConfigFileNotFound if none of the paths exist.
func loadConfigFile(paths []string, loader ConfigFileLoader, cfg any) ([]string, error) {
	for _, path := range paths {
		path = expandConfigPath(path)

		data, err := os.ReadFile(path) //nolint:gosec // intentional config file loading with user-provided paths
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

	return nil, fmt.Errorf("%w: none of %v found", ErrConfigFileNotFound, paths)
}

// cachedHomeDir caches the user's home directory for the process lifetime.
// os.UserHomeDir() performs a syscall (getenv) on every call; the result never
// changes during a process lifetime, so caching avoids redundant syscalls when
// multiple config paths use ~/ expansion.
//
//nolint:gochecknoglobals // Process-lifetime cache, idiomatic sync.OnceValue usage
var cachedHomeDir = sync.OnceValue(func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return home
})

// expandConfigPath expands environment variables and the leading ~ in a path.
func expandConfigPath(path string) string {
	path = os.ExpandEnv(path)

	if strings.HasPrefix(path, "~/") {
		if home := cachedHomeDir(); home != "" {
			path = filepath.Join(home, path[2:])
		}
	}

	return path
}

// FilterSetFields returns the field names from tags whose flag name is present.
func FilterSetFields(tags []FlagTag, present map[string]bool) []string {
	setFields := make([]string, 0, len(present))
	for _, tag := range tags {
		if present[tag.Name] {
			setFields = append(setFields, tag.Field)
		}
	}

	return setFields
}

// collectKeysRecursive walks a JSON raw-message map, recording every key at
// every nesting level. This lets FilterSetFields detect leaf-level flag names
// that appear inside nested config-file objects (e.g. {"db":{"host":"x"}} → "host").
func collectKeysRecursive(raw map[string]json.RawMessage, keys map[string]bool) {
	for k, v := range raw {
		keys[k] = true

		var nested map[string]json.RawMessage
		if json.Unmarshal(v, &nested) == nil {
			collectKeysRecursive(nested, keys)
		}
	}
}

// fieldValueToString converts a reflect.Value to its string representation.
// Delegates to the unified formatFieldValue in flag_helpers.go.
func fieldValueToString(field reflect.Value) (string, bool) {
	if !field.IsValid() {
		return "", false
	}

	result := formatFieldValue(field)

	return result, result != ""
}

// resolveConfigFlag scans args for a --config flag override.
// Supports --config <value> and --config=<value> forms.
// Returns the flag value if found, otherwise an empty string.
func resolveConfigFlag(args []string) string {
	for i, arg := range args {
		if arg == "--config" && i+1 < len(args) {
			return args[i+1]
		}

		if val, ok := strings.CutPrefix(arg, "--config="); ok {
			return val
		}
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

		field := fieldByTag(v, r.tags[i])
		if !field.IsValid() {
			continue
		}

		if s, ok := fieldValueToString(field); ok {
			r.tags[i].Default = s
		}
	}
}

// WithConfigFile adds config file loading with the given search paths.
// Paths are tried in order; the first existing file wins.
// Environment variables and ~ are expanded in paths.
// Only JSON files are supported in the core package;
// use WithConfigFileLoader for YAML/TOML support.
// loadConfigFileOrSkip attempts to load a config file, returning nil on "not found".
// This is the helper used during CLI initialization.
func (cli *CLI[T]) loadConfigFileOrSkip() ([]string, error) {
	if cli.spec.configLoader == nil || len(cli.spec.configFilePaths) == 0 {
		return nil, nil
	}

	paths := cli.spec.configFilePaths

	// Check for --config flag override.
	if override := resolveConfigFlag(os.Args[1:]); override != "" {
		paths = []string{override}
	}

	setFields, err := loadConfigFile(paths, cli.spec.configLoader, cli.config)
	if err != nil && errors.Is(err, ErrConfigFileNotFound) {
		return nil, nil
	}

	return setFields, err
}
