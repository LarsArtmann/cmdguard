package v2

import (
	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

// CLIOption is a functional option for configuring a CLI.
type CLIOption[T any] func(*CLI[T])

// WithCLIVersion sets the version string.
func WithCLIVersion[T any](version string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.version = version
		cli.rootCmd.Version = version
	}
}

// WithCLILong sets the long description.
func WithCLILong[T any](long string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.long = long
		cli.rootCmd.Long = long
	}
}

// WithCLIScope sets a custom DI scope.
func WithCLIScope[T any](scope *Scope) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.scope = scope
	}
}

// WithSilenceErrors suppresses automatic error printing from cobra.
func WithSilenceErrors[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.SilenceErrors = true
	}
}

// WithSilenceUsage suppresses automatic usage printing on error.
func WithSilenceUsage[T any]() CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.SilenceUsage = true
	}
}

// WithColor enables or disables colored output from fang.
// When disabled, falls back to cobra's default plain text output.
func WithColor[T any](enabled bool) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.useFang = enabled
	}
}

// WithFangOptions sets fang options for the CLI's Execute method.
func WithFangOptions[T any](opts ...fang.Option) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.fangOpts = append(cli.fangOpts, opts...)
	}
}

// WithMiddleware adds middleware that wraps every command handler.
// Middleware are applied in order: first wraps the second, etc.
func WithMiddleware[T any](mw ...Middleware[T]) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.middleware = append(cli.middleware, mw...)
	}
}

// WithGroup registers a command group on the root command.
// Groups organize commands in help output under titled sections.
// Use the Group field on Command to assign a command to a registered group.
func WithGroup[T any](id, title string) CLIOption[T] {
	return func(cli *CLI[T]) {
		cli.rootCmd.AddGroup(&cobra.Group{ID: id, Title: title})
	}
}
