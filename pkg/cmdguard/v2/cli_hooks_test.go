package v2

import (
	"context"
	"errors"
	"slices"
	"testing"
)

var errTest = errors.New("test error")

func TestCLI_PreRunE_PostRunE(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		hookName string
		setupCmd func(order *[]string) Command[testAppConfig, NoFlags]
		want     []string
	}{
		{
			name:     "calls PreRunE before RunE",
			hookName: "pre",
			setupCmd: func(order *[]string) Command[testAppConfig, NoFlags] {
				return Command[testAppConfig, NoFlags]{
					spec: commandSpec{
						use:     "test",
						preRunE: &typedHook[testAppConfig, NoFlags]{fn: makeHookRunE(order, "pre")},
					},
					runE: makeHookRunE(order, "run"),
				}
			},
			want: []string{"pre", "run"},
		},
		{
			name:     "calls PostRunE after RunE",
			hookName: "post",
			setupCmd: func(order *[]string) Command[testAppConfig, NoFlags] {
				return Command[testAppConfig, NoFlags]{
					spec: commandSpec{
						use:      "test",
						postRunE: &typedHook[testAppConfig, NoFlags]{fn: makeHookRunE(order, "post")},
					},
					runE: makeHookRunE(order, "run"),
				}
			},
			want: []string{"run", "post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var order []string

			cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cmd := tt.setupCmd(&order)
			if err := AddCommand(cli, cmd); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			err = cli.ExecuteWithArgs(t.Context(), []string{"test"})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !slices.Equal(order, tt.want) {
				t.Errorf("order = %v, want %v", order, tt.want)
			}
		})
	}

	t.Run("PreRunE error stops execution", func(t *testing.T) {
		t.Parallel()

		called := false

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		cmd := Command[testAppConfig, NoFlags]{
			spec: commandSpec{
				use: "test",
				preRunE: &typedHook[testAppConfig, NoFlags]{
					fn: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
						return errTest
					},
				},
			},
			runE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
				called = true

				return nil
			},
		}
		addCommand(t, cli, cmd)

		err = cli.ExecuteWithArgs(t.Context(), []string{"test"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if called {
			t.Error("RunE should not have been called when PreRunE fails")
		}
	})
}

func TestCLI_CommandOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		use            string
		hidden         bool
		deprecated     string
		aliases        []string
		version        string
		wantHidden     bool
		wantDeprecated string
		wantAliases    []string
		wantVersion    string
	}{
		{
			name:       "hidden command",
			use:        "secret",
			hidden:     true,
			wantHidden: true,
		},
		{
			name:           "deprecated command",
			use:            "old",
			deprecated:     "use new-cmd instead",
			wantDeprecated: "use new-cmd instead",
		},
		{
			name:        "command with aliases",
			use:         "list",
			aliases:     []string{"ls", "l"},
			wantAliases: []string{"ls", "l"},
		},
		{
			name:        "command with version",
			use:         "versioned",
			version:     "v1.2.3",
			wantVersion: "v1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cmd := Command[testAppConfig, NoFlags]{
				spec: commandSpec{
					use:        tt.use,
					hidden:     tt.hidden,
					deprecated: tt.deprecated,
					aliases:    tt.aliases,
					version:    tt.version,
				},
				runE: noOpHandlerForTestAppConfig(),
			}
			if err := AddCommand(cli, cmd); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			cobraCmd := cli.RootCommand().Commands()[0]
			if cobraCmd.Hidden != tt.wantHidden {
				t.Errorf("Hidden = %v, want %v", cobraCmd.Hidden, tt.wantHidden)
			}

			if cobraCmd.Deprecated != tt.wantDeprecated {
				t.Errorf("Deprecated = %q, want %q", cobraCmd.Deprecated, tt.wantDeprecated)
			}

			if !slices.Equal(cobraCmd.Aliases, tt.wantAliases) {
				t.Errorf("Aliases = %v, want %v", cobraCmd.Aliases, tt.wantAliases)
			}

			if cobraCmd.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", cobraCmd.Version, tt.wantVersion)
			}
		})
	}
}
