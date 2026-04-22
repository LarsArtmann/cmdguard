package v2_test

import (
	"context"
	"fmt"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// newSimpleCmd is a local helper to create simple commands for examples.
func newSimpleCmd[C any](name, message, short string) (v2.Command[C, v2.NoFlags], error) {
	return v2.NewCommand[C, v2.NoFlags](name,
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

	_, err := v2.NewCommand[config, *greetFlags]("greet",
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

	parent, err := v2.NewParentCommand[config, v2.NoFlags]("items",
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

	cmd := v2.MustNewCommand[config, v2.NoFlags]("version",
		func(ctx context.Context, cfg *config, flags v2.NoFlags) error {
			fmt.Println("v1.0.0")

			return nil
		},
		v2.WithShort[config, v2.NoFlags]("Print version"),
	)

	fmt.Println("command:", cmd.Use())
	// Output: command: version
}
