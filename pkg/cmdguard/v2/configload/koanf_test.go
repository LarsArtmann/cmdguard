package configload_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2/configload"
	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

// flatConfig matches the config struct used in loader_test.go
// to verify koanf produces identical results for flat files.
type koanfFlatConfig struct {
	Name    string `flag:"name"    default:""      help:"Name"`
	Port    int    `flag:"port"    default:"8080"  help:"Port"`
	Verbose bool   `flag:"verbose" default:"false" help:"Verbose"`
}

// nestedConfig uses dotted flag names to match flattened nested YAML keys.
type koanfNestedConfig struct {
	ServerPort   int    `flag:"server.port"   default:"8080" help:"Server port"`
	ServerHost   string `flag:"server.host"   default:""     help:"Server host"`
	DatabaseName string `flag:"database.name" default:""     help:"Database name"`
	DatabasePort int    `flag:"database.port" default:"5432" help:"Database port"`
	LogLevel     string `flag:"log.level"     default:"info" help:"Log level"`
}

func TestKoanfLoader_YAML(t *testing.T) {
	t.Parallel()

	t.Run("loads flat YAML config", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		err := os.WriteFile(path, []byte("name: test\nport: 9090\nverbose: true\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEqString(t, cfg.Name, "test", "Name")
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
		testutil.AssertBoolTrue(t, cfg.Verbose, "Verbose")
		testutil.AssertFieldLen(t, setFields, 3, "setFields")
	})

	t.Run("returns only present keys", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		err := os.WriteFile(path, []byte("name: only-name\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldLen(t, setFields, 1, "setFields")
		testutil.AssertFieldEqString(t, cfg.Name, "only-name", "Name")
	})

	t.Run("invalid YAML returns ErrConfigFileParse", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		err := os.WriteFile(path, []byte("{{invalid yaml\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, v2.ErrConfigFileParse) {
			t.Errorf("expected ErrConfigFileParse, got: %v", err)
		}
	})
}

func TestKoanfLoader_JSON(t *testing.T) {
	t.Parallel()

	t.Run("loads flat JSON config", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		err := os.WriteFile(path, []byte(`{"name":"test","port":9090,"verbose":true}`), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEqString(t, cfg.Name, "test", "Name")
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
		testutil.AssertBoolTrue(t, cfg.Verbose, "Verbose")
		testutil.AssertFieldLen(t, setFields, 3, "setFields")
	})

	t.Run("invalid JSON returns ErrConfigFileParse", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		err := os.WriteFile(path, []byte("{{invalid json"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, v2.ErrConfigFileParse) {
			t.Errorf("expected ErrConfigFileParse, got: %v", err)
		}
	})
}

func TestKoanfLoader_NestedConfig(t *testing.T) {
	t.Parallel()

	t.Run("flattens nested YAML keys to dotted flag names", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `server:
  port: 3000
  host: localhost
database:
  name: myapp
  port: 5433
log:
  level: debug
`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.ServerPort != 3000 {
			t.Errorf("expected server.port 3000, got %d", cfg.ServerPort)
		}
		testutil.AssertFieldEqString(t, cfg.ServerHost, "localhost", "ServerHost")
		testutil.AssertFieldEqString(t, cfg.DatabaseName, "myapp", "DatabaseName")
		if cfg.DatabasePort != 5433 {
			t.Errorf("expected database.port 5433, got %d", cfg.DatabasePort)
		}
		testutil.AssertFieldEqString(t, cfg.LogLevel, "debug", "LogLevel")
		testutil.AssertFieldLen(t, setFields, 5, "setFields")
	})

	t.Run("returns only deeply-nested present keys", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `server:
  port: 3000
`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldLen(t, setFields, 1, "setFields")
		if cfg.ServerPort != 3000 {
			t.Errorf("expected server.port 3000, got %d", cfg.ServerPort)
		}
	})

	t.Run("flattens nested JSON keys to dotted flag names", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		content := `{"server":{"port":3000,"host":"example.com"},"database":{"name":"testdb"}}`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.ServerPort != 3000 {
			t.Errorf("expected server.port 3000, got %d", cfg.ServerPort)
		}
		testutil.AssertFieldEqString(t, cfg.ServerHost, "example.com", "ServerHost")
		testutil.AssertFieldEqString(t, cfg.DatabaseName, "testdb", "DatabaseName")
		testutil.AssertFieldLen(t, setFields, 3, "setFields")
	})
}

func TestKoanfLoader_MultiplePaths(t *testing.T) {
	t.Parallel()

	t.Run("loads first existing file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		first := filepath.Join(dir, "first.yaml")
		second := filepath.Join(dir, "second.yaml")

		err := os.WriteFile(second, []byte("name: from-second\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(first, second)
		cfg := koanfFlatConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEqString(t, cfg.Name, "from-second", "Name")
		testutil.AssertFieldLen(t, setFields, 1, "setFields")
	})

	t.Run("skips missing files and loads first found", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		missing := filepath.Join(dir, "missing.yaml")
		existing := filepath.Join(dir, "existing.yaml")

		err := os.WriteFile(existing, []byte("name: found-it\nport: 4444\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(missing, existing)
		cfg := koanfFlatConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldEqString(t, cfg.Name, "found-it", "Name")
		if cfg.Port != 4444 {
			t.Errorf("expected port 4444, got %d", cfg.Port)
		}
		testutil.AssertFieldLen(t, setFields, 2, "setFields")
	})

	t.Run("returns ErrConfigFileNotFound when no paths exist", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		loader := configload.NewKoanfLoader(
			filepath.Join(dir, "a.yaml"),
			filepath.Join(dir, "b.json"),
		)
		cfg := koanfFlatConfig{}

		_, err := loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, v2.ErrConfigFileNotFound) {
			t.Errorf("expected ErrConfigFileNotFound, got: %v", err)
		}
	})
}

func TestKoanfLoader_FormatDetection(t *testing.T) {
	t.Parallel()

	t.Run("detects .yml extension", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yml")
		err := os.WriteFile(path, []byte("name: yml-ext\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Name, "yml-ext", "Name")
	})

	t.Run("unknown extension returns ErrConfigFileNotFound", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.toml")
		err := os.WriteFile(path, []byte("name = \"toml\"\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := configload.NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, v2.ErrConfigFileNotFound) {
			t.Errorf("expected ErrConfigFileNotFound for unsupported extension, got: %v", err)
		}
	})
}

func TestKoanfLoader_PathExpansion(t *testing.T) {
	t.Run("expands environment variables in path", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		err := os.WriteFile(path, []byte("name: env-expanded\n"), 0o600)
		testutil.AssertNoError(t, err)

		t.Setenv("MY_TEST_CONFIG_DIR", dir)

		loader := configload.NewKoanfLoader("$MY_TEST_CONFIG_DIR/config.yaml")
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Name, "env-expanded", "Name")
	})
}
