//nolint:nlreturn,paralleltest // test file with inline handler returns; top-level parallel blocked by t.Setenv subtests
package v2_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
)

type precedenceConfig struct {
	Name string `flag:"name" default:"default" env:"TEST_NAME"`
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
	cli, err := v2.NewCLI[precedenceConfig](
		"app", "My app", precedenceConfig{},
		v2.WithConfigFile[precedenceConfig](configPath),
	)
	if err != nil {
		t.Fatal(err)
	}

	cmd, err := v2.NewCommand(
		"test",
		v2.NoFlags{},
		func(_ context.Context, cfg *precedenceConfig, _ v2.NoFlags) error {
			if cfg.Name != wantName {
				t.Errorf("Name = %q, want %q", cfg.Name, wantName)
			}
			called = true
			return nil
		},
		v2.WithShort("Test"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := v2.AddCommand(cli, cmd); err != nil {
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
