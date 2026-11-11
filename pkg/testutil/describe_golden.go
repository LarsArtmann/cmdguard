package testutil

import (
	"bytes"
	"flag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// RegisterUpdateGoldensFlag wires the -update-golden convention flag into a
// test binary. It is a NO-OP when any package in the binary already
// registered the flag (the core v5 drift tests own their own registration);
// a double registration panics at init with "flag redefined", which is how
// the 2026-11-11 CI red on a docs-only commit happened. Call it from an
// init in the package whose goldens the harness guards:
//
//	func init() { testutil.RegisterUpdateGoldensFlag() }
//
// AssertDescribeContract always honors the flag whether this binary
// registered it here or somewhere else.
func RegisterUpdateGoldensFlag() {
	if flag.Lookup("update-golden") != nil {
		return
	}

	flag.Bool("update-golden", false, "rewrite the describe-contract goldens")
}

// updateGoldensRequested reports whether the binary's -update-golden flag is
// set, regardless of which package registered it.
func updateGoldensRequested() bool {
	f := flag.Lookup("update-golden")

	return f != nil && f.Value.String() == "true"
}

// DescribeDescriber is the slice of *v5.CLI[T] the harness needs. Structural
// on purpose: consumers keep their concrete CLI[T] type at the call site.
type DescribeDescriber interface {
	WriteDescription(w io.Writer) error
}

// AssertDescribeContract fails the test when the CLI's machine contract
// (the `describe` JSON) drifts from the committed golden file. With
// -update-golden it rewrites the golden instead.
//
// The golden captures the FULL command tree: every command's flags, env
// names, defaults, required/secret markers, and the schema version. A
// consumer CI runs one line per CLI:
//
//	testutil.AssertDescribeContract(t, cli, "testdata/describe.json")
func AssertDescribeContract(t *testing.T, cli DescribeDescriber, goldenPath string) {
	t.Helper()

	var buf bytes.Buffer
	if err := cli.WriteDescription(&buf); err != nil {
		t.Fatalf("WriteDescription: %v", err)
	}

	got := buf.String()

	if updateGoldensRequested() {
		const goldenDirPerm = 0o755

		if mkErr := os.MkdirAll(filepath.Dir(goldenPath), goldenDirPerm); mkErr != nil {
			t.Fatalf("creating golden dir: %v", mkErr)
		}

		const goldenFilePerm = 0o600

		if writeErr := os.WriteFile(goldenPath, []byte(got), goldenFilePerm); writeErr != nil {
			t.Fatalf("writing golden: %v", writeErr)
		}

		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("golden %s missing; run with -update-golden to create it", goldenPath)
		}

		t.Fatalf("reading golden: %v", err)
	}

	if got != string(want) {
		t.Fatalf(
			"describe contract drifted from %s\nregenerate with: go test ./... -run TestDescribeContract -update-golden",
			goldenPath,
		)
	}
}

// DescribeContractFromBinary runs a built CLI binary with its hidden
// `describe` command and returns the JSON contract — the CI-side capture
// for consumers whose CLI is exercised as a process, not a library.
func DescribeContractFromBinary(t *testing.T, binary string) string {
	t.Helper()

	ctx := t.Context()

	out, err := exec.CommandContext(ctx, binary, "describe").Output()
	if err != nil {
		t.Fatalf("%s describe: %v", binary, err)
	}

	return string(out)
}
