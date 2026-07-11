//nolint:nlreturn,paralleltest,tagliatelle // test file: inline returns, t.Setenv subtests, JSON tags intentionally match flag names
package v3_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

type precedenceConfig struct {
	Name string `flag:"name" default:"default" env:"TEST_NAME"`
}

type dbConfig struct {
	Host string `flag:"db-host" json:"db-host" default:"localhost"`
	Port int    `flag:"db-port" json:"db-port" default:"5432"`
}

type nestedConfigRoot struct {
	Database dbConfig `json:"Database"`
}

func writeNestedConfigFile(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	data := `{"Database": {"db-host": "db.example.com", "db-port": 6543}}`
	if err := os.WriteFile(configPath, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}

	return configPath
}

func TestConfigFileNestedStructs(t *testing.T) {
	t.Parallel()

	configPath := writeNestedConfigFile(t)

	cli, err := v3.NewCLI(
		"app", "My app", nestedConfigRoot{},
		v3.WithConfigFile(configPath),
		v3.WithFang(false),
	)
	if err != nil {
		t.Fatal(err)
	}

	var gotHost string
	var gotPort int

	cmd, err := v3.NewCommand(
		"test",
		v3.NoFlags{},
		func(_ context.Context, cfg *nestedConfigRoot, _ v3.NoFlags) error {
			gotHost = cfg.Database.Host
			gotPort = cfg.Database.Port
			return nil
		},
		v3.WithShort("Test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := v3.AddCommand(cli, cmd); err != nil {
		t.Fatal(err)
	}

	if err := cli.ExecuteWithArgs(context.Background(), []string{"test"}); err != nil {
		t.Fatal(err)
	}

	if gotHost != "db.example.com" {
		t.Errorf("Database.Host = %q, want %q", gotHost, "db.example.com")
	}

	if gotPort != 6543 {
		t.Errorf("Database.Port = %d, want %d", gotPort, 6543)
	}
}

func writeTestConfigFile(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	if err := os.WriteFile(configPath, []byte(`{"name": "file"}`), 0o600); err != nil {
		t.Fatal(err)
	}

	return configPath
}

func runPrecedenceTest(t *testing.T, configPath, envValue string, args []string, wantName string) {
	t.Helper()

	if envValue != "" {
		t.Setenv("TEST_NAME", envValue)
	}

	called := false
	cli, err := v3.NewCLI(
		"app", "My app", precedenceConfig{},
		v3.WithConfigFile(configPath),
	)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err := v3.NewCommand(
		"test",
		v3.NoFlags{},
		func(_ context.Context, cfg *precedenceConfig, _ v3.NoFlags) error {
			if cfg.Name != wantName {
				t.Errorf("Name = %q, want %q", cfg.Name, wantName)
			}
			called = true
			return nil
		},
		v3.WithShort("Test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := v3.AddCommand(cli, cmd); err != nil {
		t.Fatal(err)
	}

	if err := cli.ExecuteWithArgs(context.Background(), args); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("handler was not called")
	}
}

func TestConfigFilePrecedence(t *testing.T) {
	// Not parallel because some subtests use t.Setenv

	tests := []struct {
		name       string
		configPath string
		envValue   string
		args       []string
		wantName   string
		parallel   bool
	}{
		{
			name:       "config file overrides tag default",
			configPath: writeTestConfigFile(t),
			args:       []string{"test"},
			wantName:   "file",
			parallel:   true,
		},
		{
			name:       "flag overrides config file",
			configPath: writeTestConfigFile(t),
			args:       []string{"test", "--name", "flag"},
			wantName:   "flag",
			parallel:   true,
		},
		{
			name:       "env overrides config file",
			configPath: writeTestConfigFile(t),
			envValue:   "env",
			args:       []string{"test"},
			wantName:   "env",
			parallel:   false,
		},
		{
			name:       "flag overrides env and config file",
			configPath: writeTestConfigFile(t),
			envValue:   "env",
			args:       []string{"test", "--name", "flag"},
			wantName:   "flag",
			parallel:   false,
		},
		{
			name:       "missing config file is not an error",
			configPath: "/does/not/exist.json",
			args:       []string{"test"},
			wantName:   "default",
			parallel:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.parallel {
				t.Parallel()
			}
			runPrecedenceTest(t, tt.configPath, tt.envValue, tt.args, tt.wantName)
		})
	}
}
