package cmdguard

import (
	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func newTestCommand(use string) *cobra.Command {
	return &cobra.Command{
		Use:  use,
		RunE: testutil.NoOpCobraRunE(),
	}
}
