package v2_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/samber/do/v2"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// newSimpleCmd is a local helper to create simple commands for examples.
func newSimpleCmd[C any](name, message, short string) (v2.Command[C, v2.NoFlags], error) {
	return v2.NewCommand[C, v2.NoFlags](
		name,
		func(_ context.Context, _ *C, _ v2.NoFlags) error {
			fmt.Println(message)

			return nil
		},
		v2.WithShort[C, v2.NoFlags](short),
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

	_, err := v2.NewCommand[config, *greetFlags](
		"greet",
		func(ctx context.Context, cfg *config, flags *greetFlags) error {
			fmt.Printf("Hello, %s!", flags.Name)

			return nil
		},
		v2.WithShort[config, *greetFlags]("Greet someone"),
		v2.WithFlags[config, *greetFlags](&greetFlags{}),
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

	parent, err := v2.NewParentCommand[config, v2.NoFlags](
		"items",
		"Item management commands",
		[]v2.Command[config, v2.NoFlags]{listCmd, createCmd},
		v2.WithShort[config, v2.NoFlags]("Item management"),
	)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("parent created:", parent.Use())
	// Output: parent created: items
}

// ExampleMustNewCommand demonstrates the panic variant for compile-time-known config.
func ExampleMustNewCommand() {
	type config struct{}

	cmd := v2.MustNewCommand[config, v2.NoFlags](
		"version",
		func(ctx context.Context, cfg *config, flags v2.NoFlags) error {
			fmt.Println("v1.0.0")

			return nil
		},
		v2.WithShort[config, v2.NoFlags]("Print version"),
	)

	fmt.Println("command:", cmd.Use())
	// Output: command: version
}

// ExampleNewCLI demonstrates creating a CLI with options.
func ExampleNewCLI() {
	type config struct {
		Verbose bool `flag:"verbose" short:"v" default:"false" help:"Enable verbose output"`
	}

	cli, err := v2.NewCLI[config](
		"myapp", "My application", config{},
		v2.WithCLIVersion[config]("1.0.0"),
		v2.WithEnvPrefix[config]("MYAPP_"),
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

	cli, _ := v2.NewCLI[config]("myapp", "My application", config{})
	scope := cli.Scope()

	// Register a service.
	_ = v2.Provide(scope, func(i do.Injector) (*Database, error) {
		return &Database{DSN: "postgres://localhost"}, nil
	})

	// Invoke the service.
	db, err := v2.Invoke[*Database](scope)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	fmt.Println("dsn:", db.DSN)
	// Output: dsn: postgres://localhost
}

// ExampleOutputTable demonstrates rendering data in multiple formats.
func ExampleOutputTable() {
	headers := []string{"Name", "Age"}
	rows := [][]string{
		{"Alice", "30"},
		{"Bob", "25"},
	}

	// Output as JSON.
	_ = v2.OutputTable(v2.FormatJSON, headers, rows)

	fmt.Println("output rendered")
	// Output:
	// {"Headers":["Name","Age"],"Rows":[["Alice","30"],["Bob","25"]],"Footer":null}
	// output rendered
}

// ExampleTimingMiddleware demonstrates adding timing middleware to a CLI.
func ExampleTimingMiddleware() {
	type config struct{}

	cli, _ := v2.NewCLI[config](
		"myapp", "My application", config{},
		v2.WithMiddleware[config](v2.TimingMiddleware[config](func(name string, d time.Duration) {
			fmt.Printf("%s took %v\n", name, d)
		})),
	)

	fmt.Println("middleware registered:", cli != nil)
	// Output: middleware registered: true
}

// ExampleNewExitError demonstrates creating an error with a custom exit code.
func ExampleNewExitError() {
	innerErr := errors.New("configuration invalid")
	exitErr, err := v2.NewExitError(2, innerErr)
	if err != nil {
		fmt.Println("error:", err)

		return
	}

	if exitCoder, ok := errors.AsType[v2.ExitCoder](exitErr); ok {
		fmt.Println("exit code:", exitCoder.ExitCode())
	}
	// Output: exit code: 2
}
