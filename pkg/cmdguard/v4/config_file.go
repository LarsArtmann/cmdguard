package v4

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
)

// ConfigFileLoader loads configuration from a file.
// Implementations populate a config struct and return the list of struct field
// names that were explicitly set.
//
// The data parameter may be nil for loaders that handle their own file reading
// (e.g. KoanfLoader, which reads from configured paths via SetPaths). Custom
// byte-based loaders receive the raw file contents in data.
type ConfigFileLoader interface {
	Load(data []byte, cfg any) (setFields []string, err error)
}

// StrictConfigKeyLoader is an optional [ConfigFileLoader] capability: the
// loader can reject config-file keys the config struct does not define.
// [KoanfLoader] implements it, and [WithStrictConfigKeys] turns it on. Custom
// loaders that do not implement it keep accepting unknown keys, because
// cmdguard cannot inspect their parse step.
type StrictConfigKeyLoader interface {
	SetStrictKeys(enabled bool)
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

// FilterSetFields returns the struct field names whose flag name or config-file
// key is present. Both spellings are accepted because a config file may address
// a field either way; the loader normalizes keys before this runs, so the
// canonical key is the one that normally matches.
func FilterSetFields(tags []FlagTag, present map[string]bool) []string {
	setFields := make([]string, 0, len(present))

	for _, tag := range tags {
		if present[tag.Name] || (tag.Key != "" && present[tag.Key]) {
			setFields = append(setFields, tag.Field)
		}
	}

	return setFields
}

// collectKeysRecursive walks a JSON raw-message map, recording every key at
// every nesting level. This lets FilterSetFields detect leaf-level flag names
// that appear inside nested config-file objects (e.g. {"db":{"host":"x"}} → "host").
func collectKeysRecursive(raw map[string]jsontext.Value, keys map[string]bool) {
	for k, v := range raw {
		keys[k] = true

		var nested map[string]jsontext.Value
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

// loadConfigFileOrSkip attempts to load a config file, returning nil on "not found".
// KoanfLoader manages its own file reading; custom loaders receive bytes via
// loadConfigFile.
//
// It also returns the path of the file that was loaded (empty when none was),
// so the CLI can report where a value came from instead of guessing.
func (cli *CLI[T]) loadConfigFileOrSkip() ([]string, string, error) {
	if cli.spec.configLoader == nil || len(cli.spec.configFilePaths) == 0 {
		return nil, "", nil
	}

	paths := cli.configSearchPaths()

	if strict, ok := cli.spec.configLoader.(StrictConfigKeyLoader); ok {
		strict.SetStrictKeys(cli.spec.strictConfigKeys)
	}

	if kl, ok := cli.spec.configLoader.(*KoanfLoader); ok {
		kl.SetPaths(paths...)

		setFields, err := skipNotFound(kl.Load(nil, cli.config))

		return setFields, kl.LoadedPath(), err
	}

	setFields, err := skipNotFound(loadConfigFile(paths, cli.spec.configLoader, cli.config))

	return setFields, firstExistingConfigPath(paths), err
}

// configSearchPaths returns the config file search order, honouring a --config
// override found in the process arguments.
func (cli *CLI[T]) configSearchPaths() []string {
	paths := slices.Clone(cli.spec.configFilePaths)

	if override := resolveConfigFlag(os.Args[1:]); override != "" {
		paths = []string{override}
	}

	return paths
}

// firstExistingConfigPath returns the first expanded path that exists, or an
// empty string. It exists to report which file a custom loader read without
// duplicating the loader's own search loop.
func firstExistingConfigPath(paths []string) string {
	for _, path := range paths {
		expanded := expandConfigPath(path)

		info, err := os.Stat(expanded)
		if err == nil && !info.IsDir() {
			return expanded
		}
	}

	return ""
}

func skipNotFound(setFields []string, err error) ([]string, error) {
	if err != nil && errors.Is(err, ErrConfigFileNotFound) {
		return nil, nil
	}

	return setFields, err
}
