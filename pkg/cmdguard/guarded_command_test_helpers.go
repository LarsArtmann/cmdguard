package cmdguard

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/testutil"
)

func newTestCommand(use string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		RunE:  testutil.NoOpCobraRunE(),
	}
}

func assertErrorContains(t *testing.T, err error, substrings ...string) {
	t.Helper()
	testutil.AssertErrorContains(t, err, substrings...)
}
