// Package benchmarks provides performance benchmarks for cmdguard v4.
package benchmarks

import (
	"os"
	"reflect"
	"testing"

	"github.com/samber/do/v2"

	v4 "github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4"
	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type BenchConfig struct {
	Verbose bool   `default:"false" flag:"verbose" short:"v"`
	Output  string `default:"text"  flag:"output"  short:"o"`
}

type BenchFlags struct {
	Name  string `default:"World" flag:"name"  short:"n"`
	Count int    `default:"1"     flag:"count" short:"c"`
}

// newBenchCommand creates a command with standard benchmark configuration.
func newBenchCommand(b *testing.B, use, short string) v4.Command[BenchConfig, v4.NoFlags] {
	b.Helper()

	cmd, err := v4.NewCommand(
		use,
		v4.NoFlags{},
		testutil.NoOpRunE[BenchConfig, v4.NoFlags],
		v4.WithShort(short),
	)
	if err != nil {
		b.Fatal(err)
	}

	return cmd
}

// BenchmarkNew measures CLI creation performance.
func BenchmarkNew(b *testing.B) {
	defaults := BenchConfig{}

	for b.Loop() {
		cli, err := v4.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
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
		cli, err := v4.NewCLI[BenchConfig](
			"myapp",
			"short",
			defaults,
			v4.WithCLILong("long description"),
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
		testCli, err := v4.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
		if err != nil {
			b.Fatal(err)
		}

		err = v4.AddCommand(testCli, newBenchCommand(b, "greet", "Greet someone"))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExecute measures command execution.
func BenchmarkExecute(b *testing.B) {
	defaults := BenchConfig{}

	cli, err := v4.NewCLI[BenchConfig]("myapp", "My CLI", defaults)
	if err != nil {
		b.Fatal(err)
	}

	if err := v4.AddCommand(cli, newBenchCommand(b, "hello", "Say hello")); err != nil {
		b.Fatal(err)
	}

	ctx := b.Context()

	// Suppress stdout to prevent help-text rendering from causing I/O
	// contention that inflates co-running benchmarks.
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		b.Fatal(err)
	}
	defer devNull.Close()

	origStdout := os.Stdout
	os.Stdout = devNull
	defer func() { os.Stdout = origStdout }()

	for b.Loop() {
		// Execute with help to avoid actual command running
		err := cli.ExecuteWithArgs(ctx, []string{"--help"})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCommandCreation measures creating command definitions.

// BenchmarkNewCommand measures the NewCommand constructor.
func BenchmarkNewCommand(b *testing.B) {
	for b.Loop() {
		cmd, err := v4.NewCommand(
			"greet",
			v4.NoFlags{},
			testutil.NoOpRunE[BenchConfig, v4.NoFlags],
			v4.WithShort("Greet someone"),
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
		scope := v4.NewScope("benchmark")
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
		tags, err := v4.ParseFlagTags(cfg)
		if err != nil {
			b.Fatal(err)
		}

		_ = tags
	}
}

// runFlagRegistryBench is the shared body for the FlagRegistry creation
// benchmarks. Kept as a function (not a sub-benchmark) so each variant shows
// up as a distinct row in `go test -bench` output for comparison.
func runFlagRegistryBench(b *testing.B) {
	type TestConfig struct {
		Name    string `default:"test"  flag:"name"    help:"Name"`
		Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
	}

	cfg := &TestConfig{}

	for b.Loop() {
		registry, err := v4.NewFlagRegistry(cfg)
		if err != nil {
			b.Fatal(err)
		}

		_ = registry
	}
}

// BenchmarkFlagRegistryCreation measures FlagRegistry creation.
func BenchmarkFlagRegistryCreation(b *testing.B) {
	runFlagRegistryBench(b)
}

// BenchmarkFlagRegistryCOW measures FlagRegistry creation with copy-on-write
// (no per-instance customization — should be cheaper than eager clone).
func BenchmarkFlagRegistryCOW(b *testing.B) {
	runFlagRegistryBench(b)
}

// BenchmarkCommandValidate measures command validation.
func BenchmarkCommandValidate(b *testing.B) {
	cmd := newBenchCommand(b, "test", "Test command")

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
			duration, err := v4.ParseDuration(d)
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
			ll, err := v4.ParseLogLevel(level)
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
			url, err := v4.ParseURL(u)
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
			email, err := v4.ParseEmail(e)
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
			port, err := v4.ParsePort(p)
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
			fp, err := v4.ParseFilePath(p, false)
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
			hostPort, err := v4.ParseHostPort(hp)
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

// provideBenchService is a helper that provides a benchService with the given name to a scope.
func provideBenchService(scope *v4.Scope, name string) error {
	return v4.Provide[benchService](scope, func(i do.Injector) (benchService, error) {
		return benchService{Name: name}, nil
	})
}

// BenchmarkScopeProvide measures DI service registration.
func BenchmarkScopeProvide(b *testing.B) {
	for b.Loop() {
		scope := v4.NewScope("bench")
		err := provideBenchService(scope, "test")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScopeInvoke measures DI service retrieval.
func BenchmarkScopeInvoke(b *testing.B) {
	scope := v4.NewScope("bench")

	err := provideBenchService(scope, "test")
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		svc, err := v4.Invoke[benchService](scope)
		if err != nil {
			b.Fatal(err)
		}

		_ = svc
	}
}

// BenchmarkScopeCreationWithOpts measures DI scope creation with InjectorOpts.
func BenchmarkScopeCreationWithOpts(b *testing.B) {
	for b.Loop() {
		scope := v4.NewScopeWithOpts("bench", nil)
		_ = scope
	}
}

// BenchmarkCloneScope measures scope cloning for test isolation.
func BenchmarkCloneScope(b *testing.B) {
	scope := v4.NewScope("bench")
	err := provideBenchService(scope, "test")
	if err != nil {
		b.Fatal(err)
	}

	_, err = v4.Invoke[benchService](scope)
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		cloned := v4.CloneScope(scope)
		_, err := v4.Invoke[benchService](cloned)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScopeProvideInvokeCycle measures full register-then-retrieve cycle.
func BenchmarkScopeProvideInvokeCycle(b *testing.B) {
	for b.Loop() {
		scope := v4.NewScope("bench")
		err := provideBenchService(scope, "cycle")
		if err != nil {
			b.Fatal(err)
		}

		_, err = v4.Invoke[benchService](scope)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFlagRegistryCOWWithWrite measures FlagRegistry creation followed by
// a per-instance write (triggers lazy clone).
func BenchmarkFlagRegistryCOWWithWrite(b *testing.B) {
	type TestConfig struct {
		Name string `default:"test" flag:"name" help:"Name"`
	}

	cfg := &TestConfig{}

	for b.Loop() {
		registry, err := v4.NewFlagRegistry(cfg)
		if err != nil {
			b.Fatal(err)
		}

		registry.RegisterTypeHandler(reflect.TypeFor[string](), v4.TypeHandlerFunc{
			ParseFunc:   func(value string, _ v4.FlagTag) (any, error) { return value, nil },
			DefaultFunc: func(_ v4.FlagTag) any { return "" },
		})
	}
}

// BenchmarkTagsSeq measures iterator-based tag traversal (zero alloc).
func BenchmarkTagsSeq(b *testing.B) {
	type TestConfig struct {
		Name    string `default:"test"  flag:"name"    help:"Name"`
		Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
	}

	registry, err := v4.NewFlagRegistry(&TestConfig{})
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		for range registry.TagsSeq() {
		}
	}
}

// BenchmarkTagsSlice measures slice-based tag traversal (allocates).
func BenchmarkTagsSlice(b *testing.B) {
	type TestConfig struct {
		Name    string `default:"test"  flag:"name"    help:"Name"`
		Verbose bool   `default:"false" flag:"verbose" help:"Verbose"`
	}

	registry, err := v4.NewFlagRegistry(&TestConfig{})
	if err != nil {
		b.Fatal(err)
	}

	for b.Loop() {
		tags := registry.Tags()
		for range tags {
		}
	}
}
