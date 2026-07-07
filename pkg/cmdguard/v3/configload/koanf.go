package configload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	cmdguard "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// KoanfLoader implements ConfigFileLoader using koanf for robust config parsing.
// Supports nested structs, multiple file formats, and automatic format detection.
//
// Unlike the basic loaders (YAML/TOML/JSON), KoanfLoader handles nested config
// structures by flattening dot-separated keys to match flag tag names.
//
// Usage:
//
//	cli, _ := cmdguard.NewCLI[Config]("app", "My app", Config{},
//	    cmdguard.WithConfigFileLoader[Config](
//	        configload.NewKoanfLoader(configload.KoanfWithPaths(
//	            "$HOME/.config/app/config.yaml",
//	            "/etc/app/config.json",
//	        )),
//	    ),
//	)
type KoanfLoader struct {
	paths   []string
	options koanfOptions
}

type koanfOptions struct {
	delimiter string
}

// KoanfOption configures a KoanfLoader.
type KoanfOption func(*koanfOptions)

// KoanfWithDelimiter sets the key delimiter (default ".").
func KoanfWithDelimiter(d string) KoanfOption {
	return func(o *koanfOptions) { o.delimiter = d }
}

// NewKoanfLoader creates a koanf-based config loader.
// Paths are tried in order; the first existing file is loaded.
// File format is detected from the extension (.yaml/.yml, .json, .toml).
func NewKoanfLoader(paths ...string) *KoanfLoader {
	return &KoanfLoader{
		paths: paths,
		options: koanfOptions{
			delimiter: ".",
		},
	}
}

// Load reads config from the first existing file and populates cfg.
// It uses koanf to parse the file and flatten nested keys into flag tag names.
func (l *KoanfLoader) Load(_ []byte, cfg any) ([]string, error) {
	k := koanf.New(l.options.delimiter)

	loaded := false

	for _, p := range l.paths {
		expanded := expandKoanfPath(p)

		parser := parserForPath(expanded)
		if parser == nil {
			continue
		}

		err := k.Load(file.Provider(expanded), parser)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, fmt.Errorf("%w: loading %q with koanf: %w", cmdguard.ErrConfigFileParse, expanded, err)
		}

		loaded = true

		break
	}

	if !loaded {
		return nil, fmt.Errorf("%w: none of %v found", cmdguard.ErrConfigFileNotFound, l.paths)
	}

	if err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{
		Tag:       "flag",
		FlatPaths: true,
	}); err != nil {
		return nil, fmt.Errorf("%w: unmarshaling koanf config: %w", cmdguard.ErrConfigFileParse, err)
	}

	tags, err := cmdguard.ParseFlagTags(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing flag tags: %w", cmdguard.ErrConfigFileParse, err)
	}

	keys := k.Keys()

	present := make(map[string]bool, len(keys))
	for _, k := range keys {
		present[k] = true
	}

	return cmdguard.FilterSetFields(tags, present), nil
}

// parserForPath returns the appropriate koanf parser based on file extension.
func parserForPath(path string) koanf.Parser {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser()
	case ".json":
		return json.Parser()
	default:
		return nil
	}
}

// expandKoanfPath expands environment variables and ~ in a path.
func expandKoanfPath(path string) string {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, "$HOME", "~")

	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return os.ExpandEnv(path)
}
