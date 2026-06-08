// Package configload provides optional config file loaders for YAML and TOML.
//
// Use these with v2.WithConfigFileLoader[T]():
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithConfigFileLoader[Config](configload.YAML(), "$HOME/.config/app/config.yaml"),
//	)
package configload

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-faster/yaml"
	"github.com/pelletier/go-toml/v2"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

// unmarshalFunc parses bytes into a target value.
type unmarshalFunc func(data []byte, v any) error

// genericLoader is a ConfigFileLoader that uses a generic unmarshal function.
type genericLoader struct {
	unmarshal unmarshalFunc
}

func (l *genericLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]any
	if err := l.unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	tags, err := v2.ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", v2.ErrConfigFileParse, err)
	}

	present := make(map[string]bool, len(raw))
	for k := range raw {
		present[k] = true
	}

	setFields := v2.FilterSetFields(tags, present)

	if err := l.unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	return setFields, nil
}

// YAML returns a ConfigFileLoader for YAML files.
// Supports flat key-value objects where keys match flag tag names.
func YAML() v2.ConfigFileLoader {
	return &genericLoader{unmarshal: yaml.Unmarshal}
}

// TOML returns a ConfigFileLoader for TOML files.
// Supports flat key-value objects where keys match flag tag names.
func TOML() v2.ConfigFileLoader {
	return &genericLoader{unmarshal: toml.Unmarshal}
}

// JSON returns a ConfigFileLoader for JSON files.
// This is identical to the core package's built-in JSON loader.
func JSON() v2.ConfigFileLoader {
	return &jsonLoader{}
}

type jsonLoader struct{}

func (l *jsonLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	tags, err := v2.ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", v2.ErrConfigFileParse, err)
	}

	present := make(map[string]bool, len(raw))
	for k := range raw {
		present[k] = true
	}

	setFields := v2.FilterSetFields(tags, present)

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	return setFields, nil
}

// Auto returns a ConfigFileLoader that selects the appropriate loader
// based on the file extension of the first existing path.
// Falls back to JSON if the extension is unknown.
//
// This is useful when you don't know the file format at compile time:
//
//	cli, _ := v2.NewCLI[Config]("app", "My app", Config{},
//	    v2.WithConfigFileLoader[Config](configload.Auto(),
//	        "$HOME/.config/app/config.yaml",
//	        "$HOME/.config/app/config.json",
//	    ),
//	)
func Auto() v2.ConfigFileLoader {
	return &autoLoader{}
}

type autoLoader struct{}

func (l *autoLoader) Load(data []byte, cfg any) ([]string, error) {
	setFields, err := JSON().Load(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("auto-detect fallback to JSON: %w", err)
	}

	return setFields, nil
}

// LoaderForPath returns the appropriate loader for a file path based on its extension.
// Returns JSON loader for unknown extensions.
func LoaderForPath(path string) v2.ConfigFileLoader {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		return YAML()
	case ".toml":
		return TOML()
	case ".json":
		return JSON()
	default:
		return JSON()
	}
}
