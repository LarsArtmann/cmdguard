package v2

import (
	"context"
	"testing"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

type envTestConfig struct {
	Debug bool `default:"false" flag:"debug" help:"Debug mode"`
}

func TestEnvTag_Integration(t *testing.T) {
	tests := []struct {
		name       string
		envKey     string
		envValue   string
		args       []string
		wantHost   string
		cliOptions []CLIOption[envTestConfig]
	}{
		{
			name:     "env var provides default when flag not set",
			envKey:   "DB_HOST",
			envValue: "db.example.com",
			args:     []string{"connect"},
			wantHost: "db.example.com",
		},
		{
			name:     "explicit flag overrides env var",
			envKey:   "DB_HOST",
			envValue: "db.example.com",
			args:     []string{"connect", "--host", "explicit.example.com"},
			wantHost: "explicit.example.com",
		},
		{
			name:     "default used when env var not set",
			args:     []string{"connect"},
			wantHost: "localhost",
		},
		{
			name:     "env prefix is applied to command flags",
			envKey:   "MYAPP_DB_HOST",
			envValue: "prefixed.example.com",
			args:     []string{"connect"},
			wantHost: "prefixed.example.com",
			cliOptions: []CLIOption[envTestConfig]{
				WithEnvPrefix[envTestConfig]("MYAPP_"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type flags struct {
				Host string `flag:"host" env:"DB_HOST" default:"localhost" help:"Database host"`
			}

			if tt.envKey != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			var result string

			cli, err := NewCLI[envTestConfig]("app", "test", envTestConfig{}, tt.cliOptions...)
			testutil.AssertNoError(t, err)

			cmd, err := NewCommand(
				"connect",
				&flags{},
				func(_ context.Context, _ *envTestConfig, f *flags) error {
					result = f.Host

					return nil
				},
				WithShort("Connect"),
			)
			testutil.AssertNoError(t, err)
			testutil.AssertNoError(t, AddCommand(cli, cmd))

			err = cli.ExecuteWithArgs(t.Context(), tt.args)
			testutil.AssertNoError(t, err)
			if result != tt.wantHost {
				t.Errorf("host = %q, want %q", result, tt.wantHost)
			}
		})
	}
}
