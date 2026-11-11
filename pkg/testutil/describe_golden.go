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

// updateDescribeGolden mirrors the -update-golden convention used by the
// core drift tests: run
//
//	go test ./... -run TestDescribeContract -update-golden
//
// to regenerate the committed contract snapshots after an intentional
// contract change. The describe JSON carries a schema_version, so a
// consumer bumping majors can compare the old and new contracts
// deliberately instead of discovering flag drift at runtime.
var updateDescribeGolden = flag.Bool("update-golden", false, "rewrite the describe-contract goldens") //nolint:gochecknoglobals // flag registration is package-level by design (Go testing convention)

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

	if *updateDescribeGolden {
		const goldenDirPerm = 0o755

		if mkErr := os.MkdirAll(filepath.Dir(goldenPath), goldenDirPerm); mkErr != nil {
			t.Fatalf("creating golden dir: %v", mkErr)
		}
		if writeErr := os.WriteFile(goldenPath, []byte(got), 0o600); writeErr != nil { //nolint:gosec // committed test fixture, not a secret
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
