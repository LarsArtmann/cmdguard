// Package benchmarks provides performance benchmarks for cmdguard v2.
package benchmarks

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

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

// noOpRunE is a shared no-op RunE function to reduce duplication.
func noOpRunE[T any, F any](_ context.Context, _ *T, _ F) error {
	return nil
}

// newBenchCommand creates a command with standard benchmark configuration.
func newBenchCommand(use, short string) v2.Command[BenchConfig, v2.NoFlags] {
	return v2.Command[BenchConfig, v2.NoFlags]{
		Use:   use,
		Short: short,
		RunE:  noOpRunE[BenchConfig, v2.NoFlags],
	}
}

// BenchmarkNew measures CLI creation performance.
func BenchmarkNew(b *testing.B) {
	defaults := BenchConfig{}

	for b.Loop() {
		cli, err := v2.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
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
		cli, err := v2.NewCLI[BenchConfig](
			"myapp",
			"short",
			defaults,
			v2.WithCLILong[BenchConfig]("long description"),
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
		testCli, err := v2.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
		if err != nil {
			b.Fatal(err)
		}

		err = v2.AddCommand(testCli, newBenchCommand("greet", "Greet someone"))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExecute measures command execution.
func BenchmarkExecute(b *testing.B) {
	defaults := BenchConfig{}

	cli, err := v2.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
	if err != nil {
		b.Fatal(err)
	}

	if err := v2.AddCommand(cli, newBenchCommand("hello", "Say hello")); err != nil {
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
		cmd, err := v2.NewCommand(
			"greet",
			v2.WithShort[BenchConfig, v2.NoFlags]("Greet someone"),
			v2.WithRunE(noOpRunE[BenchConfig, v2.NoFlags]),
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
	cmd := newBenchCommand("test", "Test command")

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

// BenchmarkParseURL measures URL parsing.
func BenchmarkParseURL(b *testing.B) {
	urls := []string{
		"https://example.com",
		"https://example.com:8080/path",
		"https://user:pass@example.com/api/v1",
		"http://localhost:3000/health",
		"https://api.example.com/v2/users?page=1",
	}

	for b.Loop() {
		for _, u := range urls {
			url, err := v2.ParseURL(u)
			if err != nil {
				b.Fatal(err)
			}

			_ = url
		}
	}
}

// BenchmarkParseEmail measures email parsing.
func BenchmarkParseEmail(b *testing.B) {
	emails := []string{
		"user@example.com",
		"test.user@company.org",
		"admin+alias@subdomain.example.co.uk",
		"user123@mail.example.com",
		"contact@support.example.net",
	}

	for b.Loop() {
		for _, e := range emails {
			email, err := v2.ParseEmail(e)
			if err != nil {
				b.Fatal(err)
			}

			_ = email
		}
	}
}

// benchmarkParsePorts is a helper for parsing port benchmarks.
func benchmarkParsePorts(b *testing.B, ports []string) {
	for b.Loop() {
		for _, p := range ports {
			port, err := v2.ParsePort(p)
			if err != nil {
				b.Fatal(err)
			}

			_ = port
		}
	}
}

// BenchmarkParsePort measures port parsing (numeric).
func BenchmarkParsePort(b *testing.B) {
	ports := []string{
		"8080",
		"3000",
		"443",
		"22",
		"9000",
	}

	benchmarkParsePorts(b, ports)
}

// BenchmarkParsePortNamed measures port parsing (named ports).
func BenchmarkParsePortNamed(b *testing.B) {
	ports := []string{"http", "https", "ssh", "ftp", "smtp"}

	benchmarkParsePorts(b, ports)
}

// BenchmarkParseFilePath measures file path parsing.
func BenchmarkParseFilePath(b *testing.B) {
	paths := []string{
		"/tmp/test.txt",
		"/home/user/projects/cmdguard/main.go",
		"/var/log/app.log",
		"./relative/path/file.yaml",
		"/etc/config/app.toml",
	}

	for b.Loop() {
		for _, p := range paths {
			fp, err := v2.ParseFilePath(p, false)
			if err != nil {
				b.Fatal(err)
			}

			_ = fp
		}
	}
}

// BenchmarkParseHostPort measures host:port parsing.
func BenchmarkParseHostPort(b *testing.B) {
	hostports := []string{
		"localhost:8080",
		"example.com:443",
		"127.0.0.1:3000",
		":8080",
		"api.example.com:9000",
	}

	for b.Loop() {
		for _, hp := range hostports {
			hostPort, err := v2.ParseHostPort(hp)
			if err != nil {
				b.Fatal(err)
			}

			_ = hostPort
		}
	}
}

// benchService is a shared service type for scope benchmarks.
type benchService struct {
	Name string
}

// provideBenchService is a helper that provides a benchService to a scope.
func provideBenchService(scope *v2.Scope) error {
	return v2.Provide[benchService](scope, func(i do.Injector) (benchService, error) {
		return benchService{Name: "test"}, nil
	})
}

// BenchmarkScopeProvide measures DI service registration.
func BenchmarkScopeProvide(b *testing.B) {
	scope := v2.NewScope("bench")

	for b.Loop() {
		err := provideBenchService(scope)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScopeInvoke measures DI service retrieval.
func BenchmarkScopeInvoke(b *testing.B) {
	scope := v2.NewScope("bench")

	err := provideBenchService(scope)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		svc, err := v2.Invoke[benchService](scope)
		if err != nil {
			b.Fatal(err)
		}

		_ = svc
	}
}
