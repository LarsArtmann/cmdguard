package v4

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type koanfFlatConfig struct {
	Name    string `flag:"name"    default:""      help:"Name"`
	Port    int    `flag:"port"    default:"8080"  help:"Port"`
	Verbose bool   `flag:"verbose" default:"false" help:"Verbose"`
}

type koanfServerConfig struct {
	Port int    `flag:"port" default:"8080" help:"Server port"`
	Host string `flag:"host" default:""     help:"Server host"`
}

type koanfDatabaseConfig struct {
	Name string `flag:"name" default:"" help:"Database name"`
}

type koanfLogConfig struct {
	Level string `flag:"level" default:"info" help:"Log level"`
}

type koanfNestedConfig struct {
	Server   koanfServerConfig   // no flag tag → recurse
	Database koanfDatabaseConfig // no flag tag → recurse
	Log      koanfLogConfig      // no flag tag → recurse
}

func TestKoanfLoader_YAML(t *testing.T) {
	t.Parallel()

	t.Run("loads flat YAML config", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		err := os.WriteFile(path, []byte("name: test\nport: 9090\nverbose: true\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
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

		loader := NewKoanfLoader(path)
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

		loader := NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrConfigFileParse) {
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

		loader := NewKoanfLoader(path)
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

		loader := NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrConfigFileParse) {
			t.Errorf("expected ErrConfigFileParse, got: %v", err)
		}
	})
}

func TestKoanfLoader_TOML(t *testing.T) {
	t.Parallel()

	t.Run("loads flat TOML config", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.toml")
		err := os.WriteFile(path, []byte(`name = "test"`+"\n"+"port = 9090\n"+"verbose = true\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
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
}

func TestKoanfLoader_NestedConfig(t *testing.T) {
	t.Parallel()

	t.Run("loads nested YAML into nested structs", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `server:
  port: 3000
  host: localhost
database:
  name: myapp
log:
  level: debug
`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Server.Port != 3000 {
			t.Errorf("expected Server.Port 3000, got %d", cfg.Server.Port)
		}
		testutil.AssertFieldEqString(t, cfg.Server.Host, "localhost", "Server.Host")
		testutil.AssertFieldEqString(t, cfg.Database.Name, "myapp", "Database.Name")
		testutil.AssertFieldEqString(t, cfg.Log.Level, "debug", "Log.Level")
		testutil.AssertFieldLen(t, setFields, 4, "setFields")
	})

	t.Run("returns only present nested keys", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.yaml")
		content := `server:
  port: 3000
`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		testutil.AssertFieldLen(t, setFields, 1, "setFields")
		if cfg.Server.Port != 3000 {
			t.Errorf("expected Server.Port 3000, got %d", cfg.Server.Port)
		}
	})

	t.Run("loads nested JSON into nested structs", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		content := `{"server":{"port":3000,"host":"example.com"},"database":{"name":"testdb"}}`
		err := os.WriteFile(path, []byte(content), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
		cfg := koanfNestedConfig{}

		setFields, err := loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Server.Port != 3000 {
			t.Errorf("expected Server.Port 3000, got %d", cfg.Server.Port)
		}
		testutil.AssertFieldEqString(t, cfg.Server.Host, "example.com", "Server.Host")
		testutil.AssertFieldEqString(t, cfg.Database.Name, "testdb", "Database.Name")
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

		loader := NewKoanfLoader(first, second)
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

		loader := NewKoanfLoader(missing, existing)
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
		loader := NewKoanfLoader(
			filepath.Join(dir, "a.yaml"),
			filepath.Join(dir, "b.json"),
		)
		cfg := koanfFlatConfig{}

		_, err := loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrConfigFileNotFound) {
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

		loader := NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Name, "yml-ext", "Name")
	})

	t.Run("detects .toml extension", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.toml")
		err := os.WriteFile(path, []byte(`name = "toml-ext"`+"\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertNoError(t, err)
		testutil.AssertFieldEqString(t, cfg.Name, "toml-ext", "Name")
	})

	t.Run("unsupported extension returns ErrConfigFileNotFound", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		path := filepath.Join(dir, "config.ini")
		err := os.WriteFile(path, []byte("name=ini\n"), 0o600)
		testutil.AssertNoError(t, err)

		loader := NewKoanfLoader(path)
		cfg := koanfFlatConfig{}

		_, err = loader.Load(nil, &cfg)
		testutil.AssertExpectedError(t, err)
		if !errors.Is(err, ErrConfigFileNotFound) {
			t.Errorf("expected ErrConfigFileNotFound for unsupported extension, got: %v", err)
		}
	})
}

func TestKoanfLoader_PathExpansion(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(path, []byte("name: env-expanded\n"), 0o600)
	testutil.AssertNoError(t, err)

	t.Setenv("MY_TEST_CONFIG_DIR", dir)

	loader := NewKoanfLoader("$MY_TEST_CONFIG_DIR/config.yaml")
	cfg := koanfFlatConfig{}

	_, err = loader.Load(nil, &cfg)
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, cfg.Name, "env-expanded", "Name")
}

func TestKoanfLoader_SetPaths(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	first := filepath.Join(dir, "first.yaml")
	second := filepath.Join(dir, "second.yaml")

	err := os.WriteFile(second, []byte("name: from-second\n"), 0o600)
	testutil.AssertNoError(t, err)

	loader := NewKoanfLoader(first)
	loader.SetPaths(first, second)
	cfg := koanfFlatConfig{}

	_, err = loader.Load(nil, &cfg)
	testutil.AssertNoError(t, err)
	testutil.AssertFieldEqString(t, cfg.Name, "from-second", "Name")
}
