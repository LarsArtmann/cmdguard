package v4

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	koanfjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
)

// KoanfLoader implements [ConfigFileLoader] using koanf for multi-format config
// parsing. It supports YAML, TOML, and JSON files with automatic format detection
// based on file extension (.yaml/.yml, .toml, .json).
//
// KoanfLoader reads its own files from configured paths — the [ConfigFileLoader.Load]
// data parameter is ignored. Internally, koanf parses the file into its native map
// representation, marshals it to JSON, then reuses the same JSON processing path
// as the rest of the config system: [collectKeysRecursive] for key detection,
// [FilterSetFields] for set-field tracking, and json.Unmarshal with
// MatchCaseInsensitiveNames for case-insensitive struct population.
//
// Usage:
//
//	cli, _ := cmdguard.NewCLI[Config]("app", "My app", Config{},
//	    cmdguard.WithConfigFile(
//	        "$HOME/.config/app/config.yaml",
//	        "/etc/app/config.json",
//	    ),
//	)
//
// WithConfigFile creates a KoanfLoader automatically. For advanced use cases,
// construct one explicitly and pass it to [WithConfigFileLoader].
type KoanfLoader struct {
	paths []string
}

// NewKoanfLoader creates a koanf-based config loader.
// Paths are tried in order; the first existing file wins.
// File format is detected from the extension (.yaml/.yml, .json, .toml).
func NewKoanfLoader(paths ...string) *KoanfLoader {
	return &KoanfLoader{paths: paths}
}

// SetPaths replaces the search paths. Used by [WithConfigFile] to apply
// a --config flag override at initialization time.
func (l *KoanfLoader) SetPaths(paths ...string) {
	l.paths = paths
}

// Load reads config from the first existing file and populates cfg.
// The data parameter is ignored — KoanfLoader reads its own files.
func (l *KoanfLoader) Load(_ []byte, cfg any) ([]string, error) {
	for _, p := range l.paths {
		expanded := expandConfigPath(p)

		parser := koanfParserForPath(expanded)
		if parser == nil {
			continue
		}

		data, err := os.ReadFile(expanded) //nolint:gosec // intentional config file loading with user-provided paths
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, fmt.Errorf("%w: reading %q: %w", ErrConfigFileRead, expanded, err)
		}

		k := koanf.New(".")
		if err := k.Load(&bytesProvider{data: data}, parser); err != nil {
			return nil, fmt.Errorf("%w: parsing %q: %w", ErrConfigFileParse, expanded, err)
		}

		jsonBytes, err := k.Marshal(koanfjson.Parser())
		if err != nil {
			return nil, fmt.Errorf("%w: converting %q to JSON: %w", ErrConfigFileParse, expanded, err)
		}

		return loadConfigFromJSON(jsonBytes, cfg)
	}

	return nil, fmt.Errorf("%w: none of %v found", ErrConfigFileNotFound, l.paths)
}

// bytesProvider implements [koanf.Provider] for in-memory byte slices,
// avoiding a dependency on koanf/providers/rawbytes. When a Parser is provided
// to koanf.Load, only ReadBytes is called.
type bytesProvider struct{ data []byte }

func (b *bytesProvider) ReadBytes() ([]byte, error) { return b.data, nil }

func (b *bytesProvider) Read() (map[string]any, error) {
	return nil, errKoanfReadNotImplemented
}

// errKoanfReadNotImplemented is returned by bytesProvider.Read, which is never
// called by koanf when a Parser is provided to Load.
var errKoanfReadNotImplemented = errors.New("bytesProvider.Read not implemented; koanf should use ReadBytes")

// koanfParserForPath returns the appropriate koanf parser based on file extension.
func koanfParserForPath(path string) koanf.Parser {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser()
	case ".json":
		return koanfjson.Parser()
	case ".toml":
		return toml.Parser()
	default:
		return nil
	}
}

// loadConfigFromJSON parses JSON bytes into cfg using case-insensitive field
// matching. It is the shared processing path used by KoanfLoader (after
// converting YAML/TOML to JSON via koanf) and by any future loader that
// produces JSON bytes.
func loadConfigFromJSON(data []byte, cfg any) ([]string, error) {
	var raw map[string]jsontext.Value

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

	err = json.Unmarshal(data, cfg, json.MatchCaseInsensitiveNames(true))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfigFileParse, err)
	}

	return setFields, nil
}
