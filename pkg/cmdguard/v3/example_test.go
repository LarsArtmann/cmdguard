package v3_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	output "github.com/larsartmann/go-output"
	"github.com/samber/do/v2"

	v3 "github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3"
)

// newSimpleCmd is a local helper to create simple commands for examples.
func newSimpleCmd[C any](name, message, short string) (v3.Command[C, v3.NoFlags], error) {
	return v3.NewCommand(
		name,
		v3.NoFlags{},
		func(_ context.Context, _ *C, _ v3.NoFlags) error {
			fmt.Println(message)

			return nil
		},
		v3.WithShort(short),
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

	_, err := v3.NewCommand(
		"greet",
		&greetFlags{},
		func(ctx context.Context, cfg *config, flags *greetFlags) error {
			fmt.Printf("Hello, %s!", flags.Name)

			return nil
		},
		v3.WithShort("Greet someone"),
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

	parent, err := v3.NewParentCommand[config](
		"items",
		"Item management commands",
		v3.NoFlags{},
		v3.WithSubcommands(listCmd, createCmd),
		v3.WithShort("Item management"),
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

	cmd, err := v3.NewCommand(
		"version",
		v3.NoFlags{},
		func(ctx context.Context, cfg *config, flags v3.NoFlags) error {
			fmt.Println("v1.0.0")

			return nil
		},
		v3.WithShort("Print version"),
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

	cli, err := v3.NewCLI[config](
		"myapp", "My application", config{},
		v3.WithCLIVersion("1.0.0"),
		v3.WithEnvPrefix("MYAPP_"),
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

	cli, _ := v3.NewCLI[config]("myapp", "My application", config{})
	scope := cli.Scope()

	// Register a service.
	_ = v3.Provide(scope, func(i do.Injector) (*Database, error) {
		return &Database{DSN: "postgres://localhost"}, nil
	})

	// Invoke the service.
	db, err := v3.Invoke[*Database](scope)
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

	cfg := v3.OutputConfig{Format: output.FormatJSON}
	_ = v3.OutputResult(cfg, Person{Name: "Alice"})

	// Output:
	// {
	//   "name": "Alice"
	// }
}

// ExampleTimingMiddleware demonstrates adding timing middleware to a CLI.
func ExampleTimingMiddleware() {
	type config struct{}

	cli, _ := v3.NewCLI[config](
		"myapp", "My application", config{},
		v3.WithMiddleware(v3.TimingMiddleware[config](func(name string, d time.Duration, err error) {
			fmt.Printf("%s took %v (err=%v)\n", name, d, err)
		})),
	)

	fmt.Println("middleware registered:", cli != nil)
	// Output: middleware registered: true
}

// ExampleNewExitError demonstrates creating an error with a custom exit code.
func ExampleNewExitError() {
	innerErr := errors.New("configuration invalid")
	exitErr, err := v3.NewExitError(2, innerErr)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	if exitCoder, ok := errors.AsType[v3.ExitCoder](exitErr); ok {
		fmt.Println("exit code:", exitCoder.ExitCode())
	}
	// Output: exit code: 2
}
