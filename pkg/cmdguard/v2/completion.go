package v2

import (
	"github.com/spf13/cobra"
)

// CompletionFunc is a function that returns completion candidates for a command.
// It follows cobra's ValidArgsFunction signature for compatibility.
type CompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// WithCompletion sets the shell completion function for a command.
// This wires into cobra's ValidArgsFunction for bash/zsh/fish/powershell completion.
//
// Usage:
//
//	cmd, _ := v2.NewCommand[Config, *Flags]("search", searchHandler,
//	    v2.WithCompletion[Config, *Flags](func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
//	        return []string{"apple", "banana", "cherry"}, cobra.ShellCompDirectiveNoFileComp
//	    }),
//	)
func WithCompletion[T, F any](fn CompletionFunc) CommandOption[T, F] {
	return func(cmd *Command[T, F]) {
		cmd.completionFn = fn
	}
}

// WithValidArgs sets static valid arguments for a command.
// This is a simpler alternative to WithCompletion for fixed argument lists.
func WithValidArgs[T, F any](args ...string) CommandOption[T, F] {
	return func(cmd *Command[T, F]) {
		cmd.validArgs = args
	}
}
