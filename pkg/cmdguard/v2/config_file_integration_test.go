//nolint:nlreturn,tparallel // test file with inline handler returns; top-level parallel blocked by t.Setenv subtests
package v2_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestConfigFilePrecedence(t *testing.T) {
	// Not parallel because some subtests use t.Setenv

	t.Run("config file overrides tag default", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(configPath, []byte(`{"name": "file"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name" default:"default"`
		}

		called := false
		cli, err := v2.NewCLI[Config]("app", "My app", Config{},
			v2.WithConfigFile[Config](configPath),
		)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err := v2.NewCommand[Config, v2.NoFlags]("test",
			func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
				if cfg.Name != "file" {
					t.Errorf("Name = %q, want %q", cfg.Name, "file")
				}
				called = true
				return nil
			},
			v2.WithShort[Config, v2.NoFlags]("Test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatal(err)
		}

		if err := cli.ExecuteWithArgs(context.Background(), []string{"test"}); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})

	t.Run("flag overrides config file", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(configPath, []byte(`{"name": "file"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name" default:"default"`
		}

		called := false
		cli, err := v2.NewCLI[Config]("app", "My app", Config{},
			v2.WithConfigFile[Config](configPath),
		)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err := v2.NewCommand[Config, v2.NoFlags]("test",
			func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
				if cfg.Name != "flag" {
					t.Errorf("Name = %q, want %q", cfg.Name, "flag")
				}
				called = true
				return nil
			},
			v2.WithShort[Config, v2.NoFlags]("Test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatal(err)
		}

		if err := cli.ExecuteWithArgs(context.Background(), []string{"test", "--name", "flag"}); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})

	t.Run("env overrides config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(configPath, []byte(`{"name": "file"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name" default:"default" env:"TEST_NAME"`
		}

		t.Setenv("TEST_NAME", "env")

		called := false
		cli, err := v2.NewCLI[Config]("app", "My app", Config{},
			v2.WithConfigFile[Config](configPath),
		)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err := v2.NewCommand[Config, v2.NoFlags]("test",
			func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
				if cfg.Name != "env" {
					t.Errorf("Name = %q, want %q", cfg.Name, "env")
				}
				called = true
				return nil
			},
			v2.WithShort[Config, v2.NoFlags]("Test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatal(err)
		}

		if err := cli.ExecuteWithArgs(context.Background(), []string{"test"}); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})

	t.Run("flag overrides env and config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.json")
		if err := os.WriteFile(configPath, []byte(`{"name": "file"}`), 0o600); err != nil {
			t.Fatal(err)
		}

		type Config struct {
			Name string `flag:"name" default:"default" env:"TEST_NAME2"`
		}

		t.Setenv("TEST_NAME2", "env")

		called := false
		cli, err := v2.NewCLI[Config]("app", "My app", Config{},
			v2.WithConfigFile[Config](configPath),
		)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err := v2.NewCommand[Config, v2.NoFlags]("test",
			func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
				if cfg.Name != "flag" {
					t.Errorf("Name = %q, want %q", cfg.Name, "flag")
				}
				called = true
				return nil
			},
			v2.WithShort[Config, v2.NoFlags]("Test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatal(err)
		}

		if err := cli.ExecuteWithArgs(context.Background(), []string{"test", "--name", "flag"}); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})

	t.Run("missing config file is not an error", func(t *testing.T) {
		t.Parallel()

		type Config struct {
			Name string `flag:"name" default:"default"`
		}

		called := false
		cli, err := v2.NewCLI[Config]("app", "My app", Config{},
			v2.WithConfigFile[Config]("/does/not/exist.json"),
		)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err := v2.NewCommand[Config, v2.NoFlags]("test",
			func(_ context.Context, cfg *Config, _ v2.NoFlags) error {
				if cfg.Name != "default" {
					t.Errorf("Name = %q, want %q", cfg.Name, "default")
				}
				called = true
				return nil
			},
			v2.WithShort[Config, v2.NoFlags]("Test"),
		)
		if err != nil {
			t.Fatal(err)
		}

		if err := v2.AddCommand(cli, cmd); err != nil {
			t.Fatal(err)
		}

		if err := cli.ExecuteWithArgs(context.Background(), []string{"test"}); err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("handler was not called")
		}
	})
}
