//nolint:nlreturn // test file with many inline handler returns
package v3

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// writeJSONConfigFile writes body as a JSON config file inside dir and returns the
// absolute path. Fails the test on write error. Centralizes the temp-dir + write
// pattern shared by config_file subtests that load JSON from disk.
func writeJSONConfigFile(t *testing.T, dir, body string) string {
	t.Helper()

	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestExpandConfigPath(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get user home dir:", err)
	}

	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{"tilde expansion", "~/.config/app/config.json", home},
		{"env expansion", "$HOME/.config/app/config.json", home},
		{"no expansion", "/etc/app/config.json", "/etc/app/config.json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := expandConfigPath(tt.input)
			if tt.contains != "" && !contains(result, tt.contains) {
				t.Errorf("expandConfigPath(%q) = %q, want containing %q", tt.input, result, tt.contains)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(substr) <= len(s) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestCachedHomeDir_MatchesOS(t *testing.T) {
	t.Parallel()

	osHome, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get user home dir:", err)
	}

	cached := cachedHomeDir()
	if cached != osHome {
		t.Errorf("cachedHomeDir() = %q, want %q", cached, osHome)
	}

	again := cachedHomeDir()
	if again != cached {
		t.Errorf("cachedHomeDir() not stable: first=%q, second=%q", cached, again)
	}
}

func TestLoadConfigFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("loads flat config", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name  string `flag:"name"  default:"world"`
			Count int    `flag:"count" default:"1"`
		}

		cfg := Config{}
		data := []byte(`{"name": "config", "count": 5}`)
		setFields, err := loadConfigFromJSON(data, &cfg)
		if err != nil {
			t.Fatalf("Load error: %v", err)
		}

		if cfg.Name != "config" {
			t.Errorf("Name = %q, want %q", cfg.Name, "config")
		}
		if cfg.Count != 5 {
			t.Errorf("Count = %d, want %d", cfg.Count, 5)
		}
		if len(setFields) != 2 {
			t.Errorf("setFields = %v, want 2 fields", setFields)
		}
	})

	t.Run("returns only present keys", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name  string `flag:"name"  default:"world"`
			Count int    `flag:"count" default:"1"`
		}

		cfg := Config{}
		data := []byte(`{"name": "only-name"}`)
		setFields, err := loadConfigFromJSON(data, &cfg)
		if err != nil {
			t.Fatalf("Load error: %v", err)
		}

		if len(setFields) != 1 || setFields[0] != "Name" {
			t.Errorf("setFields = %v, want [Name]", setFields)
		}
		if cfg.Count != 0 {
			t.Errorf("Count = %d, want 0 (not set)", cfg.Count)
		}
	})

	t.Run("invalid json returns error", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name string `flag:"name"`
		}

		cfg := Config{}
		data := []byte(`{invalid`)
		_, err := loadConfigFromJSON(data, &cfg)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
		if !errors.Is(err, ErrConfigFileParse) {
			t.Errorf("error = %v, want ErrConfigFileParse", err)
		}
	})
}

type testJSONFileLoader struct{}

func (l *testJSONFileLoader) Load(data []byte, cfg any) ([]string, error) {
	return loadConfigFromJSON(data, cfg)
}

func TestLoadConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("loads existing file", func(t *testing.T) {
		t.Parallel()

		path := writeJSONConfigFile(t, t.TempDir(), `{"name": "file"}`)

		type Config struct {
			Name string `flag:"name"`
		}
		cfg := Config{}
		loader := &testJSONFileLoader{}

		setFields, err := loadConfigFile([]string{path}, loader, &cfg)
		if err != nil {
			t.Fatalf("loadConfigFile error: %v", err)
		}
		if cfg.Name != "file" {
			t.Errorf("Name = %q, want %q", cfg.Name, "file")
		}
		if len(setFields) != 1 {
			t.Errorf("setFields = %v, want 1 field", setFields)
		}
	})

	t.Run("skips missing files and tries next", func(t *testing.T) {
		t.Parallel()

		path := writeJSONConfigFile(t, t.TempDir(), `{"name": "found"}`)

		type Config struct {
			Name string `flag:"name"`
		}
		cfg := Config{}
		loader := &testJSONFileLoader{}

		setFields, err := loadConfigFile([]string{"/does/not/exist.json", path}, loader, &cfg)
		if err != nil {
			t.Fatalf("loadConfigFile error: %v", err)
		}
		if cfg.Name != "found" {
			t.Errorf("Name = %q, want %q", cfg.Name, "found")
		}
		if len(setFields) != 1 {
			t.Errorf("setFields = %v, want 1 field", setFields)
		}
	})

	t.Run("returns ErrConfigFileNotFound when none exist", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name string `flag:"name"`
		}
		cfg := Config{}
		loader := &testJSONFileLoader{}

		_, err := loadConfigFile([]string{"/does/not/exist.json"}, loader, &cfg)
		if !errors.Is(err, ErrConfigFileNotFound) {
			t.Errorf("error = %v, want ErrConfigFileNotFound", err)
		}
	})
}

func TestUpdateTagDefaultsFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("updates defaults for config-file-set fields", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name  string `flag:"name"  default:"world"`
			Count int    `flag:"count" default:"1"`
		}

		registry, err := NewFlagRegistry(&Config{})
		if err != nil {
			t.Fatal(err)
		}

		cfg := Config{Name: "config", Count: 42}
		registry.updateTagDefaultsFromConfig(&cfg, []string{"Name", "Count"})

		for _, tag := range registry.Tags() {
			switch tag.Name {
			case "name":
				if tag.Default != "config" {
					t.Errorf("name default = %q, want %q", tag.Default, "config")
				}
			case "count":
				if tag.Default != "42" {
					t.Errorf("count default = %q, want %q", tag.Default, "42")
				}
			}
		}
	})

	t.Run("ignores unknown fields", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name string `flag:"name" default:"world"`
		}

		registry, err := NewFlagRegistry(&Config{})
		if err != nil {
			t.Fatal(err)
		}

		cfg := Config{Name: "config"}
		registry.updateTagDefaultsFromConfig(&cfg, []string{"Name", "Unknown"})

		// Should not panic; Name should be updated
		for _, tag := range registry.Tags() {
			if tag.Name == "name" && tag.Default != "config" {
				t.Errorf("name default = %q, want %q", tag.Default, "config")
			}
		}
	})
}

func TestResolveConfigFlag(t *testing.T) {
	t.Parallel()

	t.Run("finds --config flag", func(t *testing.T) {
		t.Parallel()

		result := resolveConfigFlag([]string{"--config", "/custom/config.json", "subcmd"})
		if result != "/custom/config.json" {
			t.Errorf("resolveConfigFlag = %q, want %q", result, "/custom/config.json")
		}
	})

	t.Run("finds --config= flag", func(t *testing.T) {
		t.Parallel()

		result := resolveConfigFlag([]string{"--config=/custom/config.json"})
		if result != "/custom/config.json" {
			t.Errorf("resolveConfigFlag = %q, want %q", result, "/custom/config.json")
		}
	})

	t.Run("returns empty when not present", func(t *testing.T) {
		t.Parallel()

		result := resolveConfigFlag([]string{"subcmd"})
		if result != "" {
			t.Errorf("resolveConfigFlag = %q, want empty", result)
		}
	})
}

func TestLoadConfigFromJSON_RecursiveKeyCollection(t *testing.T) {
	t.Parallel()

	t.Run("flat keys match flag names", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Name string `flag:"name" json:"name"`
			Port string `flag:"port" json:"port"`
		}

		cfg := config{}

		fields, err := loadConfigFromJSON([]byte(`{"name":"test","port":"8080"}`), &cfg)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if cfg.Name != "test" {
			t.Errorf("Name = %q, want %q", cfg.Name, "test")
		}

		if cfg.Port != "8080" {
			t.Errorf("Port = %q, want %q", cfg.Port, "8080")
		}

		if !sliceContains(fields, "Name") {
			t.Errorf("expected fields to contain 'Name', got: %v", fields)
		}
	})

	t.Run("nested object keys are recursively collected", func(t *testing.T) {
		t.Parallel()

		type config struct {
			Host string `flag:"host" json:"host"`
		}

		cfg := config{}

		fields, err := loadConfigFromJSON([]byte(`{"db":{"host":"localhost"}}`), &cfg)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if !sliceContains(fields, "Host") {
			t.Errorf("expected recursive key collection to find 'Host' from nested object, got: %v", fields)
		}
	})

	t.Run("invalid JSON returns parse error", func(t *testing.T) {
		t.Parallel()

		cfg := struct{}{}

		_, err := loadConfigFromJSON([]byte(`{invalid`), &cfg)
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
	})
}

func sliceContains(slice []string, s string) bool {
	return slices.Contains(slice, s)
}
