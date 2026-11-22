package testutil

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// fixturePreamble is the shared package header for temp-dir probe files.
const fixturePreamble = `package probe

import "errors"

func newErr(msg string) error { return errors.New(msg) }
`

// writeFixtureTree materializes files (slash-relative path -> content) under
// a fresh temp dir and returns the dir. The dead-sentinel tripwire's
// detection logic is regression-tested on throwaway mini-packages instead
// of probing real source files (the 2026-11-22 append+sed incident).
func writeFixtureTree(t *testing.T, files map[string]string) string {
	t.Helper()

	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("creating dir for %s: %v", name, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
	}

	return dir
}

// assertScan runs FindDeadSentinels over dir and checks the declaration
// count and the exact set of dead sentinel names (order-insensitive).
func assertScan(t *testing.T, dir, allowFile string, wantDecls int, wantDead []string) {
	t.Helper()

	scan, err := FindDeadSentinels(dir, allowFile)
	if err != nil {
		t.Fatalf("FindDeadSentinels: %v", err)
	}
	if scan.Declarations != wantDecls {
		t.Errorf("Declarations = %d, want %d", scan.Declarations, wantDecls)
	}

	got := make([]string, 0, len(scan.Findings))
	for _, f := range scan.Findings {
		got = append(got, f.Name)
	}
	slices.Sort(got)
	want := slices.Clone(wantDead)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("dead sentinels = %v, want %v", got, want)
	}
}

// TestFindDeadSentinelsDetectsDead pins the detection half of the tripwire:
// call-initialized and typed-error sentinels referenced only from the facts
// file, a test file, or another package are all dead.
func TestFindDeadSentinelsDetectsDead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		files     map[string]string
		allowFile string
		wantDecls int
		wantDead  []string
	}{
		{
			name: "call-initialized sentinel never referenced",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var ErrDead = newErr("dead")

var ErrWired = newErr("wired")

func consume() error { return ErrWired }
`,
			},
			allowFile: "",
			wantDecls: 2,
			wantDead:  []string{"ErrDead"},
		},
		{
			name: "typed error sentinel never referenced",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var errTyped error = newErr("typed")
`,
			},
			allowFile: "",
			wantDecls: 1,
			wantDead:  []string{"errTyped"},
		},
		{
			name: "reference only inside the bookkeeping file",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var ErrBook = newErr("book")
`,
				"facts.go": `package probe

var _ = ErrBook
`,
			},
			allowFile: "facts.go",
			wantDecls: 1,
			wantDead:  []string{"ErrBook"},
		},
		{
			name: "reference only in a test file",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var ErrTestOnly = newErr("test only")
`,
				"sentinel_test.go": `package probe

import "testing"

func TestReference(t *testing.T) {
	t.Parallel()
	_ = ErrTestOnly
}
`,
			},
			allowFile: "",
			wantDecls: 1,
			wantDead:  []string{"ErrTestOnly"},
		},
		{
			name: "unexported sentinel referenced from another package stays dead",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var errLonely = newErr("lonely")
`,
				"other/other.go": `package other

import probe "probe"

func useLonely() error { return probe.errLonely }
`,
			},
			allowFile: "",
			wantDecls: 1,
			wantDead:  []string{"errLonely"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertScan(t, writeFixtureTree(t, tt.files), tt.allowFile, tt.wantDecls, tt.wantDead)
		})
	}
}

// TestFindDeadSentinelsAcceptsWired pins the acceptance half of the
// tripwire: plain returns, %w wraps, and grouped var blocks all keep a
// sentinel alive, and the bookkeeping file keeps wiring when not excluded.
func TestFindDeadSentinelsAcceptsWired(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		files     map[string]string
		allowFile string
		wantDecls int
		wantDead  []string
	}{
		{
			name: "plain return and %w wrap wire sentinels",
			files: map[string]string{
				"sentinel.go": `package probe

import (
	"errors"
	"fmt"
)

func newErr(msg string) error { return errors.New(msg) }

var ErrPlain = newErr("plain")

var ErrWrapped = newErr("wrapped")

func consumePlain() error { return ErrPlain }

func consumeWrapped() error { return fmt.Errorf("ctx: %w", ErrWrapped) }
`,
			},
			allowFile: "",
			wantDecls: 2,
			wantDead:  nil,
		},
		{
			name: "grouped var block with both wired",
			files: map[string]string{
				"sentinel.go": `package probe

import "errors"

var (
	// ErrGroupA is wired below.
	ErrGroupA = errors.New("a")
	// ErrGroupB is wired below.
	ErrGroupB = errors.New("b")
)

func consumeA() error { return ErrGroupA }

func consumeB() error { return ErrGroupB }
`,
			},
			allowFile: "",
			wantDecls: 2,
			wantDead:  nil,
		},
		{
			name: "exported sentinel wired only from another package",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var ErrRemote = newErr("remote")
`,
				"other/other.go": `package other

import probe "probe"

func useRemote() error { return probe.ErrRemote }
`,
			},
			allowFile: "",
			wantDecls: 1,
			wantDead:  nil,
		},
		{
			name: "bookkeeping file wires when not excluded",
			files: map[string]string{
				"sentinel.go": fixturePreamble + `
var ErrBook = newErr("book")
`,
				"facts.go": `package probe

var _ = ErrBook
`,
			},
			allowFile: "",
			wantDecls: 1,
			wantDead:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertScan(t, writeFixtureTree(t, tt.files), tt.allowFile, tt.wantDecls, tt.wantDead)
		})
	}
}

// TestFindDeadSentinelsIgnoresNonSentinelNames keeps the Err/err name
// filter from over-matching: bare "Err" and non-error var names are not
// tripwire candidates.
func TestFindDeadSentinelsIgnoresNonSentinelNames(t *testing.T) {
	t.Parallel()

	dir := writeFixtureTree(t, map[string]string{
		"sentinel.go": `package probe

import "errors"

var Err = errors.New("bare Err has no suffix")

var warnValue = errors.New("does not start with err")
`,
	})
	assertScan(t, dir, "", 0, nil)
}

// TestFindDeadSentinelsRejectsBrokenSource surfaces parse failures instead
// of silently scanning nothing.
func TestFindDeadSentinelsRejectsBrokenSource(t *testing.T) {
	t.Parallel()

	dir := writeFixtureTree(t, map[string]string{
		"bad.go": "package probe\n\nfunc broken( {\n",
	})
	if _, err := FindDeadSentinels(dir, ""); err == nil {
		t.Error("FindDeadSentinels on broken source = nil error, want parse error")
	}
}
