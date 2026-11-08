package v5

import (
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

// defaultsFixtureConfig is the config type for the default-flip matrix tests.
type defaultsFixtureConfig struct {
	Name string `flag:"name" default:"world" help:"Name to greet"`
}

// TestV5Defaults_OnWithoutOptions pins the v5 battery defaults: a bare
// NewCLI ships strict validation, signal handling with graceful DI shutdown,
// strict config keys, error hints, fang styling, and the batteries bundle.
func TestV5Defaults_OnWithoutOptions(t *testing.T) {
	t.Parallel()

	cli, err := NewCLI("app", "App", defaultsFixtureConfig{})
	testutil.AssertNoError(t, err)

	if cli.spec.validationMode != Strict {
		t.Errorf("default validationMode = %v, want Strict", cli.spec.validationMode)
	}

	if !cli.spec.signalHandling {
		t.Error("default signalHandling = false, want true")
	}

	if !cli.spec.gracefulShutdown {
		t.Error("default gracefulShutdown = false, want true")
	}

	if !cli.spec.strictConfigKeys {
		t.Error("default strictConfigKeys = false, want true")
	}

	if !cli.spec.errorHints {
		t.Error("default errorHints = false, want true")
	}

	if !cli.spec.batteries {
		t.Error("default batteries = false, want true")
	}
}

// TestV5DefaultOptOuts is the B3 matrix: each Without* opt-out flips exactly
// its own flags and leaves the rest of the bundle untouched.
func TestV5DefaultOptOuts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		opts   []CLIOption
		checks func(*testing.T, *CLI[defaultsFixtureConfig])
	}{
		{
			name: "WithoutStrictValidation reverts to lenient only",
			opts: []CLIOption{WithoutStrictValidation()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if cli.spec.validationMode != Lenient {
					t.Errorf("validationMode = %v, want Lenient", cli.spec.validationMode)
				}

				if !cli.spec.signalHandling || !cli.spec.gracefulShutdown || !cli.spec.strictConfigKeys {
					t.Error("other battery defaults must stay on")
				}
			},
		},
		{
			name: "WithoutSignalHandling disables signals and shutdown",
			opts: []CLIOption{WithoutSignalHandling()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if cli.spec.signalHandling {
					t.Error("signalHandling = true, want false")
				}

				if cli.spec.gracefulShutdown {
					t.Error("gracefulShutdown = true, want false")
				}

				if cli.spec.validationMode != Strict || !cli.spec.strictConfigKeys {
					t.Error("other battery defaults must stay on")
				}
			},
		},
		{
			name: "WithoutGracefulShutdown keeps signal handling",
			opts: []CLIOption{WithoutGracefulShutdown()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if !cli.spec.signalHandling {
					t.Error("signalHandling = false, want true")
				}

				if cli.spec.gracefulShutdown {
					t.Error("gracefulShutdown = true, want false")
				}
			},
		},
		{
			name: "WithoutStrictConfigKeys reverts keys only",
			opts: []CLIOption{WithoutStrictConfigKeys()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if cli.spec.strictConfigKeys {
					t.Error("strictConfigKeys = true, want false")
				}

				if cli.spec.validationMode != Strict || !cli.spec.signalHandling {
					t.Error("other battery defaults must stay on")
				}
			},
		},
		{
			name: "WithoutErrorHints reverts hints only",
			opts: []CLIOption{WithoutErrorHints()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if cli.spec.errorHints {
					t.Error("errorHints = true, want false")
				}

				if cli.spec.validationMode != Strict || !cli.spec.signalHandling {
					t.Error("other battery defaults must stay on")
				}
			},
		},
		{
			name: "WithoutBatteries reverts the whole bundle",
			opts: []CLIOption{WithoutBatteries()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if cli.spec.batteries || cli.spec.signalHandling ||
					cli.spec.gracefulShutdown || cli.spec.strictConfigKeys {
					t.Error("WithoutBatteries must revert the bundle")
				}

				if cli.spec.validationMode != Lenient {
					t.Errorf("validationMode = %v, want Lenient", cli.spec.validationMode)
				}

				if !cli.spec.errorHints {
					t.Error("errorHints must stay on (own opt-out is WithoutErrorHints)")
				}
			},
		},
		{
			name: "WithBatteries restores the bundle after opt-outs",
			opts: []CLIOption{WithoutBatteries(), WithBatteries()},
			checks: func(t *testing.T, cli *CLI[defaultsFixtureConfig]) {
				t.Helper()

				if !cli.spec.batteries || !cli.spec.signalHandling ||
					!cli.spec.gracefulShutdown || !cli.spec.strictConfigKeys ||
					cli.spec.validationMode != Strict {
					t.Error("WithBatteries must restore the full bundle")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, err := NewCLI("app", "App", defaultsFixtureConfig{}, tt.opts...)
			testutil.AssertNoError(t, err)

			tt.checks(t, cli)
		})
	}
}
