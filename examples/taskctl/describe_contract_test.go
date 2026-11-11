package main

import (
	"path/filepath"
	"testing"

	"github.com/larsartmann/cmdguard/v5/pkg/testutil"
)

func init() { //nolint:gochecknoinits // the convention flag must be registered before testing's flag.Parse
	testutil.RegisterUpdateGoldensFlag()
}

// TestDescribeContract pins the machine contract (the hidden `describe`
// command's JSON) against a committed golden. The contract is what agents
// and tooling parse to discover commands, flags, env names, defaults, and
// required/secret markers — a change here is a BREAKING consumer change
// even when every Go signature is stable. Regenerate deliberately:
//
//	go test ./... -run TestDescribeContract -update-golden
//
// The auto-registered `describe` command is only wired at Execute time
// (see describe.go), so the golden is built from WriteDescription on the
// fully-constructed CLI — the same source the runtime command renders.
func TestDescribeContract(t *testing.T) {
	t.Parallel()

	cli := newEmptyTestCLI(t)

	testutil.AssertDescribeContract(t, cli, filepath.Join("testdata", "describe.golden.json"))
}
