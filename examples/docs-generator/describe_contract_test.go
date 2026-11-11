package main

import (
	"path/filepath"
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

// TestDescribeContract pins the machine contract (see taskctl's
// TestDescribeContract for the rationale) for the docs-generator example —
// the second pilot for the consumer describe-snapshot harness.
func TestDescribeContract(t *testing.T) {
	t.Parallel()

	// main() is the constructor here; mirror its construction without
	// executing (Execute would run the root handler).
	cli, err := newContractCLI()
	if err != nil {
		t.Fatalf("constructing CLI: %v", err)
	}

	testutil.AssertDescribeContract(t, cli, filepath.Join("testdata", "describe.golden.json"))
}
