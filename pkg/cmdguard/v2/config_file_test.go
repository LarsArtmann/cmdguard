//nolint:nlreturn // test file with many inline handler returns
package v2

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

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

func TestJSONLoader_Load(t *testing.T) {
	t.Parallel()

	loader := &jsonLoader{}

	t.Run("loads flat config", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name  string `flag:"name"  default:"world"`
			Count int    `flag:"count" default:"1"`
		}

		cfg := Config{}
		data := []byte(`{"name": "config", "count": 5}`)
		setFields, err := loader.Load(data, &cfg)
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
		setFields, err := loader.Load(data, &cfg)
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
		_, err := loader.Load(data, &cfg)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
		if !errors.Is(err, ErrConfigFileParse) {
			t.Errorf("error = %v, want ErrConfigFileParse", err)
		}
	})
}

func TestLoadConfigFile(t *testing.T) {
	t.Parallel()

	t.Run("loads existing file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(path, []byte(`{"name": "file"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name"`
		}
		cfg := Config{}
		loader := &jsonLoader{}

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

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(path, []byte(`{"name": "found"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name"`
		}
		cfg := Config{}
		loader := &jsonLoader{}

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
		loader := &jsonLoader{}

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
