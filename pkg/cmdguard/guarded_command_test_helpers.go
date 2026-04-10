package cmdguard

import (
	"github.com/spf13/cobra"
)

func newTestCommand(use string) *cobra.Command {
	return &cobra.Command{
		Use: use,
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}
}
