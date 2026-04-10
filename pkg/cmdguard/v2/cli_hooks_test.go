package v2

import (
	"context"
	"errors"
	"testing"
)

var errTest = errors.New("test error")

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func makeHookRunE(
	order *[]string,
	msg string,
) func(context.Context, *testAppConfig, NoFlags) error {
	return func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
		*order = append(*order, msg)

		return nil
	}
}

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
					Use:     "test",
					PreRunE: makeHookRunE(order, "pre"),
					RunE:    makeHookRunE(order, "run"),
				}
			},
			want: []string{"pre", "run"},
		},
		{
			name:     "calls PostRunE after RunE",
			hookName: "post",
			setupCmd: func(order *[]string) Command[testAppConfig, NoFlags] {
				return Command[testAppConfig, NoFlags]{
					Use:      "test",
					RunE:     makeHookRunE(order, "run"),
					PostRunE: makeHookRunE(order, "post"),
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

			if !slicesEqual(order, tt.want) {
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
			Use: "test",
			PreRunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
				return errTest
			},
			RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
				called = true

				return nil
			},
		}
		if err := AddCommand(cli, cmd); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
				Use:        tt.use,
				Hidden:     tt.hidden,
				Deprecated: tt.deprecated,
				Aliases:    tt.aliases,
				Version:    tt.version,
				RunE: func(_ context.Context, _ *testAppConfig, _ NoFlags) error {
					return nil
				},
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

			if !slicesEqual(cobraCmd.Aliases, tt.wantAliases) {
				t.Errorf("Aliases = %v, want %v", cobraCmd.Aliases, tt.wantAliases)
			}

			if cobraCmd.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", cobraCmd.Version, tt.wantVersion)
			}
		})
	}
}
