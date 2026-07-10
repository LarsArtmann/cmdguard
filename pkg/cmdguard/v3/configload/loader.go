// Package configload provides optional config file loaders for YAML and TOML.
//
// Use these with cmdguard.WithConfigFileLoader[T]():
//
//	cli, _ := cmdguard.NewCLI[Config]("app", "My app", Config{},
//	    cmdguard.WithConfigFileLoader[Config](configload.YAML(), "$HOME/.config/app/config.yaml"),
//	)
package configload

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-faster/yaml"
	toml "github.com/pelletier/go-toml/v2"

	cmdguard "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// unmarshalFunc parses bytes into a target value.
type unmarshalFunc func(data []byte, v any) error

// genericLoader is a ConfigFileLoader that uses a generic unmarshal function.
type genericLoader struct {
	unmarshal unmarshalFunc
}

func (l *genericLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]any

	err := l.unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cmdguard.ErrConfigFileParse, err)
	}

	tags, err := cmdguard.ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", cmdguard.ErrConfigFileParse, err)
	}

	present := make(map[string]bool, len(raw))
	for k := range raw {
		present[k] = true
	}

	setFields := cmdguard.FilterSetFields(tags, present)

	err = l.unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cmdguard.ErrConfigFileParse, err)
	}

	return setFields, nil
}

// YAML returns a ConfigFileLoader for YAML files.
// Supports flat key-value objects where keys match flag tag names.
func YAML() cmdguard.ConfigFileLoader {
	return &genericLoader{unmarshal: yaml.Unmarshal}
}

// TOML returns a ConfigFileLoader for TOML files.
// Supports flat key-value objects where keys match flag tag names.
func TOML() cmdguard.ConfigFileLoader {
	return &genericLoader{unmarshal: toml.Unmarshal}
}

// JSON returns a ConfigFileLoader for JSON files.
// This is identical to the core package's built-in JSON loader.
func JSON() cmdguard.ConfigFileLoader {
	return &jsonLoader{}
}

type jsonLoader struct{}

func (l *jsonLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]json.RawMessage

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cmdguard.ErrConfigFileParse, err)
	}

	tags, err := cmdguard.ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", cmdguard.ErrConfigFileParse, err)
	}

	present := make(map[string]bool, len(raw))
	for k := range raw {
		present[k] = true
	}

	setFields := cmdguard.FilterSetFields(tags, present)

	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cmdguard.ErrConfigFileParse, err)
	}

	return setFields, nil
}

// Auto returns a ConfigFileLoader that tries to parse data as YAML, then TOML,
// then JSON, returning the first successful parse result.
// Since JSON is valid YAML, JSON data is handled by the YAML parser first.
// Use LoaderForPath when the file extension is known for precise format selection.
func Auto() cmdguard.ConfigFileLoader {
	return &autoLoader{}
}

type autoLoader struct{}

func (l *autoLoader) Load(data []byte, cfg any) ([]string, error) {
	for _, attempt := range []cmdguard.ConfigFileLoader{YAML(), TOML(), JSON()} {
		setFields, err := attempt.Load(data, cfg)
		if err == nil {
			return setFields, nil
		}
	}

	return nil, fmt.Errorf("%w: auto-detect failed: tried YAML, TOML, JSON", cmdguard.ErrConfigFileParse)
}

// LoaderForPath returns the appropriate loader for a file path based on its extension.
// Returns JSON loader for unknown extensions.
func LoaderForPath(path string) cmdguard.ConfigFileLoader {
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
