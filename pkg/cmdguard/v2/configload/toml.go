package configload

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// TOML returns a ConfigFileLoader for TOML files.
// Supports flat key-value objects where keys match flag tag names.
func TOML() v2.ConfigFileLoader {
	return &tomlLoader{}
}

type tomlLoader struct{}

func (l *tomlLoader) Load(data []byte, cfg any) ([]string, error) {
	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
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

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", v2.ErrConfigFileParse, err)
	}

	return setFields, nil
}
