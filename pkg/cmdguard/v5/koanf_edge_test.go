package v5

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// writeEdgeConfig writes content to a temp file with the given extension and
// returns the absolute path.
func writeEdgeConfig(t *testing.T, dir, name, content string) string {
	t.Helper()

	path := filepath.Join(dir, name)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing %s: %v", name, err)
	}

	return path
}

// TestKoanfLoader_TOMLDatetimeRoundTrip pins the TOML-native datetime path:
// koanf parses TOML datetimes into time.Time, marshals them to JSON as
// RFC 3339 text, and a string config field receives that text.
func TestKoanfLoader_TOMLDatetimeRoundTrip(t *testing.T) {
	t.Parallel()

	type Config struct {
		Started string `flag:"started"`
	}

	dir := t.TempDir()
	path := writeEdgeConfig(t, dir, "config.toml", "started = 2024-06-01T12:30:00Z\n")

	cfg := Config{}
	loader := NewKoanfLoader(path)

	setFields, err := loader.Load(nil, &cfg)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Started != "2024-06-01T12:30:00Z" {
		t.Errorf("Started = %q, want RFC 3339 text %q", cfg.Started, "2024-06-01T12:30:00Z")
	}

	if len(setFields) != 1 {
		t.Errorf("setFields = %v, want [started]", setFields)
	}
}

// TestKoanfLoader_YAMLAnchorsAndMergeKeys covers YAML anchors with the merge
// key. yaml.v3 resolves anchors during map decoding, so a TOP-LEVEL merge
// contributes ordinary root keys; koanf flattens NESTED maps into dotted
// paths, which never address flat config fields.
func TestKoanfLoader_YAMLAnchorsAndMergeKeys(t *testing.T) {
	t.Parallel()

	type Config struct {
		Host    string `flag:"host"`
		Port    int    `flag:"port"`
		Timeout string `flag:"timeout"`
	}

	t.Run("top-level merge lands in flat fields", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := writeEdgeConfig(t, dir, "config.yaml", `
defaults: &defaults
  host: localhost
  timeout: 30s

<<: *defaults
`)

		cfg := Config{}
		loader := NewKoanfLoader(path)

		if _, err := loader.Load(nil, &cfg); err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Host != "localhost" {
			t.Errorf("Host = %q, want %q (merged from anchor)", cfg.Host, "localhost")
		}

		if cfg.Timeout != "30s" {
			t.Errorf("Timeout = %q, want %q (merged from anchor)", cfg.Timeout, "30s")
		}
	})

	t.Run("nested keys do not leak into flat fields", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := writeEdgeConfig(t, dir, "config.yaml", `
defaults: &defaults
  host: localhost

server:
  <<: *defaults
  port: 8080
`)

		cfg := Config{}
		loader := NewKoanfLoader(path)

		if _, err := loader.Load(nil, &cfg); err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Host != "" || cfg.Port != 0 {
			t.Errorf("nested server.* keys leaked into flat fields: host=%q port=%d", cfg.Host, cfg.Port)
		}
	})
}

// TestKoanfLoader_IntVersusFloat64 covers numeric fidelity across formats:
// integral YAML/TOML numbers must reach int fields as ints (not 8080.0
// float64 text) and fractional values reach float64 fields.
func TestKoanfLoader_IntVersusFloat64(t *testing.T) {
	t.Parallel()

	type Config struct {
		Port  int     `flag:"port"`
		Ratio float64 `flag:"ratio"`
		Limit int64   `flag:"limit"`
	}

	t.Run("yaml", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := writeEdgeConfig(t, dir, "config.yaml", "port: 8080\nratio: 1.5\nlimit: 9000000000\n")

		cfg := Config{}
		loader := NewKoanfLoader(path)

		if _, err := loader.Load(nil, &cfg); err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Port != 8080 {
			t.Errorf("Port = %d, want 8080", cfg.Port)
		}

		if cfg.Ratio != 1.5 {
			t.Errorf("Ratio = %v, want 1.5", cfg.Ratio)
		}

		if cfg.Limit != 9000000000 {
			t.Errorf("Limit = %d, want 9000000000 (int64 survives the JSON round trip)", cfg.Limit)
		}
	})

	t.Run("toml", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := writeEdgeConfig(t, dir, "config.toml", "port = 9090\nratio = 0.25\nlimit = 7777\n")

		cfg := Config{}
		loader := NewKoanfLoader(path)

		if _, err := loader.Load(nil, &cfg); err != nil {
			t.Fatalf("Load: %v", err)
		}

		if cfg.Port != 9090 {
			t.Errorf("Port = %d, want 9090", cfg.Port)
		}

		if cfg.Ratio != 0.25 {
			t.Errorf("Ratio = %v, want 0.25", cfg.Ratio)
		}
	})
}

// TestKoanfLoader_EmptyFile pins the 0-byte config file behavior: every
// format treats an empty file like a missing one (skipped, no fields, no
// error), so an empty config cannot fail startup.
func TestKoanfLoader_EmptyFile(t *testing.T) {
	t.Parallel()

	type Config struct {
		Name string `flag:"name"`
	}

	for _, name := range []string{"config.yaml", "config.toml", "config.json"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			path := writeEdgeConfig(t, dir, name, "")

			cfg := Config{}
			loader := NewKoanfLoader(path)

			setFields, err := loader.Load(nil, &cfg)
			if !errors.Is(err, ErrConfigFileNotFound) {
				t.Fatalf("Load(%s) = %v, want ErrConfigFileNotFound (empty file is treated as absent)", name, err)
			}

			if len(setFields) != 0 {
				t.Errorf("setFields = %v, want none", setFields)
			}

			if cfg.Name != "" {
				t.Errorf("Name = %q, want empty", cfg.Name)
			}
		})
	}
}
