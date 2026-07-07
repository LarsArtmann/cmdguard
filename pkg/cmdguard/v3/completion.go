package v3

import (
	"github.com/spf13/cobra"
)

// CompletionFunc is a function that returns completion candidates for a command.
// It follows cobra's ValidArgsFunction signature for compatibility.
type CompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)
