// Package benchmarks provides performance benchmarks for cmdguard v2.
package benchmarks

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type BenchConfig struct {
	Verbose bool   `default:"false" flag:"verbose" short:"v"`
	Output  string `default:"text"  flag:"output"  short:"o"`
}

type BenchFlags struct {
	Name  string `default:"World" flag:"name"  short:"n"`
	Count int    `default:"1"     flag:"count" short:"c"`
}

// BenchmarkNew measures CLI creation performance.
func BenchmarkNew(b *testing.B) {
	defaults := BenchConfig{}

	for b.Loop() {
		cli, err := v2.New[BenchConfig, v2.NoFlags]("myapp", "My CLI", defaults)
		if err != nil {
			b.Fatal(err)
		}

		_ = cli
	}
}

// BenchmarkNewWithLong measures CLI creation with long description.
func BenchmarkNewWithLong(b *testing.B) {
	defaults := BenchConfig{}

	for b.Loop() {
		cli, err := v2.NewWithLong[BenchConfig, v2.NoFlags](
			"myapp",
			"short",
			"long description",
			defaults,
		)
		if err != nil {
			b.Fatal(err)
		}

		_ = cli
	}
}

// BenchmarkAddCommand measures adding a command to CLI.
func BenchmarkAddCommand(b *testing.B) {
	defaults := BenchConfig{}

	for b.Loop() {
		// Need fresh CLI for each iteration since we can't add same command twice
		testCli, err := v2.New[BenchConfig, v2.NoFlags]("myapp", "My CLI", defaults)
		if err != nil {
			b.Fatal(err)
		}

		cmd := v2.Command[BenchConfig, v2.NoFlags]{
			Use:   "greet",
			Short: "Greet someone",
			RunE: func(_ context.Context, _ *BenchConfig, _ v2.NoFlags) error {
				return nil
			},
		}

		err = testCli.AddCommand(cmd)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExecute measures command execution.
func BenchmarkExecute(b *testing.B) {
	defaults := BenchConfig{}

	cli, err := v2.New[BenchConfig, v2.NoFlags]("myapp", "My CLI", defaults)
	if err != nil {
		b.Fatal(err)
	}

	cmd := v2.Command[BenchConfig, v2.NoFlags]{
		Use:   "hello",
		Short: "Say hello",
		RunE: func(_ context.Context, _ *BenchConfig, _ v2.NoFlags) error {
			return nil
		},
	}

	if err := cli.AddCommand(cmd); err != nil {
		b.Fatal(err)
	}

	ctx := b.Context()

	for b.Loop() {
		// Execute with help to avoid actual command running
		err := cli.ExecuteWithArgs(ctx, []string{"--help"})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCommandCreation measures creating command definitions.
func BenchmarkCommandCreation(b *testing.B) {
	for b.Loop() {
		cmd := v2.Command[BenchConfig, *BenchFlags]{
			Use:   "greet",
			Short: "Greet someone",
			Flags: &BenchFlags{},
			RunE: func(_ context.Context, _ *BenchConfig, _ *BenchFlags) error {
				return nil
			},
		}
		_ = cmd
	}
}

// BenchmarkNewCommand measures the NewCommand constructor.
func BenchmarkNewCommand(b *testing.B) {
	for b.Loop() {
		cmd, err := v2.NewCommand[BenchConfig, v2.NoFlags](
			"greet",
			v2.WithShort[BenchConfig, v2.NoFlags]("Greet someone"),
			v2.WithRunE[BenchConfig, v2.NoFlags](
				func(_ context.Context, _ *BenchConfig, _ v2.NoFlags) error {
					return nil
				},
			),
		)
		if err != nil {
			b.Fatal(err)
		}

		_ = cmd
	}
}

// BenchmarkScopeCreation measures DI scope creation.
func BenchmarkScopeCreation(b *testing.B) {
	for b.Loop() {
		scope := v2.NewScope("benchmark")
		_ = scope
	}
}

// BenchmarkParseFlagTags measures flag tag parsing.
func BenchmarkParseFlagTags(b *testing.B) {
	type TestConfig struct {
		Name    string `default:"test"  flag:"name"    help:"Name"    short:"n"`
		Verbose bool   `default:"false" flag:"verbose" help:"Verbose" short:"v"`
		Count   int    `default:"1"     flag:"count"   help:"Count"   short:"c"`
		Timeout string `default:"30s"   flag:"timeout" help:"Timeout"`
	}

	cfg := &TestConfig{}

	for b.Loop() {
		tags, err := v2.ParseFlagTags(cfg)
		if err != nil {
			b.Fatal(err)
		}

		_ = tags
	}
}

// BenchmarkFlagRegistryCreation measures FlagRegistry creation.
func BenchmarkFlagRegistryCreation(b *testing.B) {
	type TestConfig struct {
		Name    string `default:"test"  flag:"name"    help:"Name"`
		Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
	}

	cfg := &TestConfig{}

	for b.Loop() {
		registry, err := v2.NewFlagRegistry(cfg)
		if err != nil {
			b.Fatal(err)
		}

		_ = registry
	}
}

// BenchmarkCommandValidate measures command validation.
func BenchmarkCommandValidate(b *testing.B) {
	cmd := v2.Command[BenchConfig, v2.NoFlags]{
		Use:   "test",
		Short: "Test command",
		RunE: func(_ context.Context, _ *BenchConfig, _ v2.NoFlags) error {
			return nil
		},
	}

	for b.Loop() {
		err := cmd.Validate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseDuration measures duration parsing.
func BenchmarkParseDuration(b *testing.B) {
	durations := []string{
		"30s",
		"5m",
		"1h30m",
		"24h",
		"100ms",
	}

	for b.Loop() {
		for _, d := range durations {
			duration, err := v2.ParseDuration(d)
			if err != nil {
				b.Fatal(err)
			}

			_ = duration
		}
	}
}

// BenchmarkParseLogLevel measures log level parsing.
func BenchmarkParseLogLevel(b *testing.B) {
	levels := []string{"debug", "info", "warn", "error"}

	for b.Loop() {
		for _, level := range levels {
			ll, err := v2.ParseLogLevel(level)
			if err != nil {
				b.Fatal(err)
			}

			_ = ll
		}
	}
}
