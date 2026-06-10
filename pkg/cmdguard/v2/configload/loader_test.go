package configload_test

import (
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2/configload"
	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type config struct {
	Name    string `flag:"name"    default:""      help:"Name"`
	Port    int    `flag:"port"    default:"8080"  help:"Port"`
	Verbose bool   `flag:"verbose" default:"false" help:"Verbose"`
}

func TestYAMLLoader(t *testing.T) {
	t.Parallel()

	t.Run("loads flat config", func(t *testing.T) {
		t.Parallel()

		data := []byte("name: test\nport: 9090\n")
		cfg := config{}

		setFields, err := configload.YAML().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Name != "test" {
			t.Errorf("expected name 'test', got %q", cfg.Name)
		}
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
		if len(setFields) != 2 {
			t.Errorf("expected 2 set fields, got %d: %v", len(setFields), setFields)
		}
	})

	t.Run("returns only present keys", func(t *testing.T) {
		t.Parallel()

		data := []byte("name: only-name\n")
		cfg := config{}

		setFields, err := configload.YAML().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if len(setFields) != 1 {
			t.Errorf("expected 1 set field, got %d: %v", len(setFields), setFields)
		}
		if cfg.Name != "only-name" {
			t.Errorf("expected name 'only-name', got %q", cfg.Name)
		}
	})

	t.Run("invalid yaml returns error", func(t *testing.T) {
		t.Parallel()

		data := []byte("{{invalid yaml")
		cfg := config{}

		_, err := configload.YAML().Load(data, &cfg)
		if err == nil {
			t.Fatal("expected error for invalid YAML")
		}
	})

	t.Run("empty document returns empty setFields", func(t *testing.T) {
		t.Parallel()

		data := []byte("---\n")
		cfg := config{}

		setFields, err := configload.YAML().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if len(setFields) != 0 {
			t.Errorf("expected 0 set fields, got %d: %v", len(setFields), setFields)
		}
	})
}

func TestTOMLLoader(t *testing.T) {
	t.Parallel()

	t.Run("loads flat config", func(t *testing.T) {
		t.Parallel()

		data := []byte("name = \"test\"\nport = 9090\n")
		cfg := config{}

		setFields, err := configload.TOML().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Name != "test" {
			t.Errorf("expected name 'test', got %q", cfg.Name)
		}
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
		if len(setFields) != 2 {
			t.Errorf("expected 2 set fields, got %d: %v", len(setFields), setFields)
		}
	})

	t.Run("returns only present keys", func(t *testing.T) {
		t.Parallel()

		data := []byte("verbose = true\n")
		cfg := config{}

		setFields, err := configload.TOML().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if len(setFields) != 1 {
			t.Errorf("expected 1 set field, got %d: %v", len(setFields), setFields)
		}
		if !cfg.Verbose {
			t.Error("expected verbose=true")
		}
	})

	t.Run("invalid toml returns error", func(t *testing.T) {
		t.Parallel()

		data := []byte("{{invalid toml")
		cfg := config{}

		_, err := configload.TOML().Load(data, &cfg)
		if err == nil {
			t.Fatal("expected error for invalid TOML")
		}
	})
}

func TestJSONLoader(t *testing.T) {
	t.Parallel()

	t.Run("loads flat config", func(t *testing.T) {
		t.Parallel()

		data := []byte(`{"name":"test","port":9090}`)
		cfg := config{}

		setFields, err := configload.JSON().Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Name != "test" {
			t.Errorf("expected name 'test', got %q", cfg.Name)
		}
		if cfg.Port != 9090 {
			t.Errorf("expected port 9090, got %d", cfg.Port)
		}
		if len(setFields) != 2 {
			t.Errorf("expected 2 set fields, got %d: %v", len(setFields), setFields)
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		t.Parallel()

		data := []byte("{{invalid json")
		cfg := config{}

		_, err := configload.JSON().Load(data, &cfg)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
	})
}

func TestAutoLoader(t *testing.T) {
	t.Parallel()

	t.Run("detects format automatically", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			data   string
			expect string
		}{
			{"yaml", "name: auto-yaml\n", "auto-yaml"},
			{"json", `{"name":"auto-json"}`, "auto-json"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				cfg := config{}
				setFields, err := configload.Auto().Load([]byte(tt.data), &cfg)
				testutil.AssertNoError(t, err)

				if cfg.Name != tt.expect {
					t.Errorf("expected name %q, got %q", tt.expect, cfg.Name)
				}
				if len(setFields) != 1 {
					t.Errorf("expected 1 set field, got %d: %v", len(setFields), setFields)
				}
			})
		}
	})

	t.Run("invalid data returns error", func(t *testing.T) {
		t.Parallel()

		data := []byte("not valid yaml toml or json [[[")
		cfg := config{}

		_, err := configload.Auto().Load(data, &cfg)
		testutil.AssertExpectedError(t, err)
	})
}

func TestLoaderForPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
	}{
		{"config.yaml"},
		{"config.yml"},
		{"config.toml"},
		{"config.json"},
		{"config.txt"},
		{"config"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			_ = configload.LoaderForPath(tt.path)
		})
	}
}

func TestLoaderForPathLoadsCorrectly(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		path   string
		data   string
		expect string
	}{
		{"yaml via LoaderForPath", "config.yaml", "name: path-test\n", "path-test"},
		{"toml via LoaderForPath", "config.toml", "name = \"path-test\"\n", "path-test"},
		{"json via LoaderForPath", "config.json", `{"name":"path-test"}`, "path-test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			loader := configload.LoaderForPath(tt.path)
			data := []byte(tt.data)

			cfg := config{}
			setFields, err := loader.Load(data, &cfg)
			testutil.AssertNoError(t, err)

			if cfg.Name != tt.expect {
				t.Errorf("expected name %q, got %q", tt.expect, cfg.Name)
			}
			if len(setFields) != 1 {
				t.Errorf("expected 1 set field, got %d", len(setFields))
			}
		})
	}

	t.Run("unknown extension falls back to JSON", func(t *testing.T) {
		t.Parallel()

		loader := configload.LoaderForPath("config.txt")
		data := []byte(`{"name":"fallback-test"}`)

		cfg := config{}
		_, err := loader.Load(data, &cfg)
		testutil.AssertNoError(t, err)

		if cfg.Name != "fallback-test" {
			t.Errorf("expected name 'fallback-test', got %q", cfg.Name)
		}
	})
}
