// Package benchmarks provides performance benchmarks for cmdguard v2.
package benchmarks

import (
	"context"
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

type BenchConfig struct {
	Verbose bool   `flag:"verbose" short:"v" default:"false"`
	Output  string `flag:"output" short:"o" default:"text"`
}

type BenchFlags struct {
	Name  string `flag:"name" short:"n" default:"World"`
	Count int    `flag:"count" short:"c" default:"1"`
}

// BenchmarkNew measures CLI creation performance.
func BenchmarkNew(b *testing.B) {
	defaults := BenchConfig{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli, err := v2.NewWithLong[BenchConfig, v2.NoFlags]("myapp", "short", "long description", defaults)
		if err != nil {
			b.Fatal(err)
		}
		_ = cli
	}
}

// BenchmarkAddCommand measures adding a command to CLI.
func BenchmarkAddCommand(b *testing.B) {
	defaults := BenchConfig{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Need fresh CLI for each iteration since we can't add same command twice
		testCli, err := v2.New[BenchConfig, v2.NoFlags]("myapp", "My CLI", defaults)
		if err != nil {
			b.Fatal(err)
		}

		cmd := v2.Command[BenchConfig, v2.NoFlags]{
			Use:   "greet",
			Short: "Greet someone",
			RunE: func(ctx context.Context, cfg *BenchConfig, flags v2.NoFlags) error {
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
		RunE: func(ctx context.Context, cfg *BenchConfig, flags v2.NoFlags) error {
			return nil
		},
	}

	if err := cli.AddCommand(cmd); err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Execute with help to avoid actual command running
		err := cli.ExecuteWithArgs(ctx, []string{"--help"})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCommandCreation measures creating command definitions.
func BenchmarkCommandCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := v2.Command[BenchConfig, *BenchFlags]{
			Use:   "greet",
			Short: "Greet someone",
			Flags: &BenchFlags{},
			RunE: func(ctx context.Context, cfg *BenchConfig, flags *BenchFlags) error {
				return nil
			},
		}
		_ = cmd
	}
}

// BenchmarkNewCommand measures the NewCommand constructor.
func BenchmarkNewCommand(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd, err := v2.NewCommand[BenchConfig, v2.NoFlags]("greet",
			v2.WithShort[BenchConfig, v2.NoFlags]("Greet someone"),
			v2.WithRunE[BenchConfig, v2.NoFlags](func(ctx context.Context, cfg *BenchConfig, flags v2.NoFlags) error {
				return nil
			}),
		)
		if err != nil {
			b.Fatal(err)
		}
		_ = cmd
	}
}
