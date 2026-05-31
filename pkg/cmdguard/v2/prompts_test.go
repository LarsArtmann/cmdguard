//nolint:paralleltest // all tests modify package-level defaultPromptRunner
package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

type fakePromptRunner struct {
	stringResults  map[string]string
	selectResults  map[string]string
	confirmResults map[string]bool
	stringErrors   map[string]error
	selectErrors   map[string]error
	confirmErrors  map[string]error
}

func (f *fakePromptRunner) PromptString(title, defaultValue string) (string, error) {
	if err, ok := f.stringErrors[title]; ok {
		return "", err
	}

	if v, ok := f.stringResults[title]; ok {
		return v, nil
	}

	return defaultValue, nil
}

func (f *fakePromptRunner) PromptSelect(title string, options []string) (string, error) {
	if err, ok := f.selectErrors[title]; ok {
		return "", err
	}

	if v, ok := f.selectResults[title]; ok {
		return v, nil
	}

	if len(options) > 0 {
		return options[0], nil
	}

	return "", nil
}

func (f *fakePromptRunner) PromptConfirm(title string) (bool, error) {
	if err, ok := f.confirmErrors[title]; ok {
		return false, err
	}

	if v, ok := f.confirmResults[title]; ok {
		return v, nil
	}

	return false, nil
}

func setFakePromptRunner(t *testing.T, runner PromptRunner) {
	t.Helper()

	old := defaultPromptRunner
	defaultPromptRunner = runner
	t.Cleanup(func() { defaultPromptRunner = old })
}

func TestPromptMissingCommandFlags_StringPrompt(t *testing.T) {
	type Flags struct {
		Name string `flag:"name" prompt:"What is your name?" default:"World"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	setFakePromptRunner(t, &fakePromptRunner{
		stringResults: map[string]string{
			"What is your name?": "Alice",
		},
	})

	if err := promptMissingCommandFlags(cmd, registry); err != nil {
		t.Fatalf("prompting: %v", err)
	}

	flag := cmd.Flags().Lookup("name")
	if flag == nil {
		t.Fatal("flag 'name' not found")
	}

	if flag.Value.String() != "Alice" {
		t.Errorf("flag value = %q, want %q", flag.Value.String(), "Alice")
	}
}

func TestPromptMissingCommandFlags_SelectPrompt(t *testing.T) {
	type Flags struct {
		Color string `flag:"color" prompt:"Pick a color" values:"red,green,blue" default:"red"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	setFakePromptRunner(t, &fakePromptRunner{
		selectResults: map[string]string{
			"Pick a color": "blue",
		},
	})

	if err := promptMissingCommandFlags(cmd, registry); err != nil {
		t.Fatalf("prompting: %v", err)
	}

	flag := cmd.Flags().Lookup("color")
	if flag == nil {
		t.Fatal("flag 'color' not found")
	}

	if flag.Value.String() != "blue" {
		t.Errorf("flag value = %q, want %q", flag.Value.String(), "blue")
	}
}

func TestPromptMissingCommandFlags_ConfirmPrompt(t *testing.T) {
	type Flags struct {
		Verbose bool `flag:"verbose" prompt:"Enable verbose mode?" default:"false"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	setFakePromptRunner(t, &fakePromptRunner{
		confirmResults: map[string]bool{
			"Enable verbose mode?": true,
		},
	})

	if err := promptMissingCommandFlags(cmd, registry); err != nil {
		t.Fatalf("prompting: %v", err)
	}

	flag := cmd.Flags().Lookup("verbose")
	if flag == nil {
		t.Fatal("flag 'verbose' not found")
	}

	if flag.Value.String() != "true" {
		t.Errorf("flag value = %q, want %q", flag.Value.String(), "true")
	}
}

func TestPromptMissingCommandFlags_SkipsWhenFlagChanged(t *testing.T) {
	type Flags struct {
		Name string `flag:"name" prompt:"What is your name?" default:"World"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	if err := cmd.Flags().Set("name", "Bob"); err != nil {
		t.Fatalf("setting flag: %v", err)
	}

	setFakePromptRunner(t, &fakePromptRunner{
		stringResults: map[string]string{
			"What is your name?": "Alice",
		},
	})

	if err := promptMissingCommandFlags(cmd, registry); err != nil {
		t.Fatalf("prompting: %v", err)
	}

	flag := cmd.Flags().Lookup("name")
	if flag.Value.String() != "Bob" {
		t.Errorf("flag value = %q, want %q", flag.Value.String(), "Bob")
	}
}

func TestPromptMissingCommandFlags_SkipsNoPromptTag(t *testing.T) {
	type Flags struct {
		Name    string `flag:"name"    prompt:"What is your name?"`
		Verbose bool   `flag:"verbose"                             default:"false"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	setFakePromptRunner(t, &fakePromptRunner{
		stringResults: map[string]string{
			"What is your name?": "Alice",
		},
	})

	if err := promptMissingCommandFlags(cmd, registry); err != nil {
		t.Fatalf("prompting: %v", err)
	}

	nameFlag := cmd.Flags().Lookup("name")
	if nameFlag.Value.String() != "Alice" {
		t.Errorf("name flag value = %q, want %q", nameFlag.Value.String(), "Alice")
	}

	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag.Value.String() != "false" {
		t.Errorf("verbose flag value = %q, want %q", verboseFlag.Value.String(), "false")
	}
}

func TestPromptMissingCommandFlags_ReturnsErrorOnPromptFailure(t *testing.T) {
	type Flags struct {
		Name string `flag:"name" prompt:"What is your name?"`
	}

	cmd := &cobra.Command{Use: "test"}
	registry, err := NewFlagRegistry(&Flags{})
	if err != nil {
		t.Fatalf("creating registry: %v", err)
	}

	if err := registry.RegisterFlags(cmd); err != nil {
		t.Fatalf("registering flags: %v", err)
	}

	promptErr := errors.New("user cancelled")
	setFakePromptRunner(t, &fakePromptRunner{
		stringErrors: map[string]error{
			"What is your name?": promptErr,
		},
	})

	err = promptMissingCommandFlags(cmd, registry)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, promptErr) {
		t.Errorf("error does not wrap prompt error: %v", err)
	}
}

func TestPromptMissingCommandFlags_NilRegistry(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}

	if err := promptMissingCommandFlags(cmd, nil); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

//nolint:paralleltest // modifies package-level defaultPromptRunner
func TestWithPromptOnMissing_Integration(t *testing.T) {
	t.Run("prompts for missing flag during command execution", func(t *testing.T) {
		type Flags struct {
			Name string `flag:"name" prompt:"What is your name?"`
		}

		setFakePromptRunner(t, &fakePromptRunner{
			stringResults: map[string]string{
				"What is your name?": "Charlie",
			},
		})

		var gotName string

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("creating CLI: %v", err)
		}

		cmd, err := NewCommand[testAppConfig, *Flags](
			"greet",
			func(_ context.Context, _ *testAppConfig, flags *Flags) error {
				gotName = flags.Name

				return nil
			},
			WithShort[testAppConfig, *Flags]("Greet someone"),
			WithFlags[testAppConfig, *Flags](&Flags{}),
			WithPromptOnMissing[testAppConfig, *Flags](),
		)
		if err != nil {
			t.Fatalf("creating command: %v", err)
		}

		if err := AddCommand(cli, cmd); err != nil {
			t.Fatalf("adding command: %v", err)
		}

		if err := cli.ExecuteWithArgs(t.Context(), []string{"greet"}); err != nil {
			t.Fatalf("executing: %v", err)
		}

		if gotName != "Charlie" {
			t.Errorf("name = %q, want %q", gotName, "Charlie")
		}
	})

	t.Run("skips prompt when flag explicitly provided", func(t *testing.T) {
		type Flags struct {
			Name string `flag:"name" prompt:"What is your name?"`
		}

		setFakePromptRunner(t, &fakePromptRunner{
			stringResults: map[string]string{
				"What is your name?": "ShouldNotBeUsed",
			},
		})

		var gotName string

		cli, err := NewCLI[testAppConfig]("myapp", "My CLI", testAppConfig{})
		if err != nil {
			t.Fatalf("creating CLI: %v", err)
		}

		cmd, err := NewCommand[testAppConfig, *Flags](
			"greet",
			func(_ context.Context, _ *testAppConfig, flags *Flags) error {
				gotName = flags.Name

				return nil
			},
			WithShort[testAppConfig, *Flags]("Greet someone"),
			WithFlags[testAppConfig, *Flags](&Flags{}),
			WithPromptOnMissing[testAppConfig, *Flags](),
		)
		if err != nil {
			t.Fatalf("creating command: %v", err)
		}

		if err := AddCommand(cli, cmd); err != nil {
			t.Fatalf("adding command: %v", err)
		}

		if err := cli.ExecuteWithArgs(t.Context(), []string{"greet", "--name", "David"}); err != nil {
			t.Fatalf("executing: %v", err)
		}

		if gotName != "David" {
			t.Errorf("name = %q, want %q", gotName, "David")
		}
	})
}

func TestPromptTag_Parse(t *testing.T) {
	type Config struct {
		Name string `flag:"name" prompt:"What is your name?"`
	}

	tags, err := ParseFlagTags(&Config{})
	if err != nil {
		t.Fatalf("parsing tags: %v", err)
	}

	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}

	if tags[0].Prompt != "What is your name?" {
		t.Errorf("prompt = %q, want %q", tags[0].Prompt, "What is your name?")
	}
}
