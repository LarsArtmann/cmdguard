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

	"github.com/go-faster/yaml"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// YAML returns a ConfigFileLoader for YAML files.
// Supports flat key-value objects where keys match flag tag names.
func YAML() v2.ConfigFileLoader {
	return &yamlLoader{}
}

type yamlLoader struct{}

func (l *yamlLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
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

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	return setFields, nil
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
