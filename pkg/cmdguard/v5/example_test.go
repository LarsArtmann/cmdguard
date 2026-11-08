package v5_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	output "github.com/larsartmann/go-output"
	"github.com/samber/do/v2"

	v5 "github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5"
)

// newSimpleCmd is a local helper to create simple commands for examples.
func newSimpleCmd[C any](name, message, short string) (v5.Command[C, v5.NoFlags], error) {
	return v5.NewCommand(
		name,
		v5.NoFlags{},
		func(_ context.Context, _ *C, _ v5.NoFlags) error {
			fmt.Println(message)

			return nil
		},
		v5.WithShort(short),
	)
}

// ExampleNewCommand demonstrates creating a leaf command with the constructor API.
func ExampleNewCommand() {
	type config struct {
		Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	}

	cmd, err := newSimpleCmd[config]("hello", "Hello, World!", "Say hello")
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("command created:", cmd.Use())
	// Output: command created: hello
}

// ExampleNewCommand_withFlags demonstrates creating a command with typed flags.
func ExampleNewCommand_withFlags() {
	type config struct{}
	type greetFlags struct {
		Name string `flag:"name" short:"n" default:"World" help:"Name to greet"`
	}

	_, err := v5.NewCommand(
		"greet",
		&greetFlags{},
		func(ctx context.Context, cfg *config, flags *greetFlags) error {
			fmt.Printf("Hello, %s!", flags.Name)

			return nil
		},
		v5.WithShort("Greet someone"),
	)
	if err != nil {
		fmt.Println("error:", err)
	}

	// Output:
}

// ExampleNewParentCommand demonstrates creating a parent command with subcommands.
func ExampleNewParentCommand() {
	type config struct{}

	listCmd, _ := newSimpleCmd[config]("list", "listing items...", "List items")
	createCmd, _ := newSimpleCmd[config]("create", "creating item...", "Create item")

	parent, err := v5.NewParentCommand[config](
		"items",
		"Item management commands",
		v5.NoFlags{},
		v5.WithSubcommands(listCmd, createCmd),
		v5.WithShort("Item management"),
	)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("parent created:", parent.Use())
	// Output: parent created: items
}

// ExampleNewCommand_minimal demonstrates minimal command creation with error handling.
func ExampleNewCommand_minimal() {
	type config struct{}

	cmd, err := v5.NewCommand(
		"version",
		v5.NoFlags{},
		func(ctx context.Context, cfg *config, flags v5.NoFlags) error {
			fmt.Println("v1.0.0")

			return nil
		},
		v5.WithShort("Print version"),
	)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("command:", cmd.Use())
	// Output: command: version
}

// ExampleNewCLI demonstrates creating a CLI with options.
func ExampleNewCLI() {
	type config struct {
		Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	}

	cli, err := v5.NewCLI(
		"myapp", "My application", config{},
		v5.WithCLIVersion("1.0.0"),
		v5.WithEnvPrefix("MYAPP_"),
	)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("cli created:", cli.RootCommand().Use)
	// Output: cli created: myapp
}

// ExampleProvide demonstrates dependency injection with Provide and Invoke.
func ExampleProvide() {
	type config struct{}

	type Database struct {
		DSN string
	}

	cli, _ := v5.NewCLI("myapp", "My application", config{})
	scope := cli.Scope()

	// Register a service.
	_ = v5.Provide(scope, func(i do.Injector) (*Database, error) {
		return &Database{DSN: "postgres://localhost"}, nil
	})

	// Invoke the service.
	db, err := v5.Invoke[*Database](scope)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("dsn:", db.DSN)
	// Output: dsn: postgres://localhost
}

// ExampleOutputResult demonstrates rendering arbitrary data.
func ExampleOutputResult() {
	type Person struct {
		Name string `json:"name"`
	}

	cfg := v5.OutputConfig{Format: output.FormatJSON}
	_ = v5.OutputResult(cfg, Person{Name: "Alice"})

	// Output:
	// {
	//   "name": "Alice"
	// }
}

// ExampleTimingMiddleware demonstrates adding timing middleware to a CLI.
func ExampleTimingMiddleware() {
	type config struct{}

	cli, _ := v5.NewCLI(
		"myapp", "My application", config{},
		v5.WithMiddleware(v5.TimingMiddleware[config](func(name string, d time.Duration, err error) {
			fmt.Printf("%s took %v (err=%v)\n", name, d, err)
		})),
	)

	fmt.Println("middleware registered:", cli != nil)
	// Output: middleware registered: true
}

// ExampleNewExitError demonstrates creating an error with a custom exit code.
func ExampleNewExitError() {
	innerErr := errors.New("configuration invalid")
	exitErr, err := v5.NewExitError(2, innerErr)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	if exitCoder, ok := errors.AsType[v5.ExitCoder](exitErr); ok {
		fmt.Println("exit code:", exitCoder.ExitCode())
	}
	// Output: exit code: 2
}

// ExampleCLI_ExplainConfig demonstrates reporting where each configuration
// value came from: flag, environment, config file, or declared default.
func ExampleCLI_ExplainConfig() {
	type config struct {
		LogLevel string `flag:"log-level" env:"APP_LOG_LEVEL" default:"info" help:"Log level"`
		Token    string `flag:"token"     env:"APP_TOKEN"                    help:"API token" secret:"true"`
	}

	cli, err := v5.NewCLI[config]("myapp", "Explain configuration", config{})
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	explanation := cli.ExplainConfig()
	for _, flag := range explanation.Flags {
		value := flag.Value
		if flag.Redacted {
			value = "***"
		}

		fmt.Printf("%s: source=%s value=%s\n", flag.Name, flag.Source, value)
	}
	// Output:
	// log-level: source=default value=info
	// token: source=unset value=***
}
