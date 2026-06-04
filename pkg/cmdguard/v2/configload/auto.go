package configload

import (
	"fmt"
	"path/filepath"
	"strings"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

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
	// Auto-detect cannot determine format from raw bytes alone for the Load call.
	// The extension check happens at the CLI level before Load is called.
	// For safety, fall back to JSON.
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
