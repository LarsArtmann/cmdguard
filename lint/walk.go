package lint

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CurrentMajor is the cmdguard major version this linter tracks. Imports of
// older majors are stale (frozen, unfixed). When cmdguard ships a new major,
// bump this constant and re-verify CG002's message.
const CurrentMajor = "v4"

// cmdguardImportPrefix matches every cmdguard import path ever published:
// the pre-module root ("github.com/larsartmann/cmdguard") and each major
// ("github.com/larsartmann/cmdguard/v3/pkg/cmdguard/v3", sub-modules such as
// "github.com/larsartmann/cmdguard/spinner", and so on).
const cmdguardImportPrefix = "github.com/larsartmann/cmdguard"

// fangImportPath is charm.land/fang — the cobra styling runner whose direct
// use bypasses the cmdguard lifecycle layer (finding CG001).
const fangImportPath = "charm.land/fang/v2"

// knownPackageNames maps import paths this linter reasons about to their real
// Go package names, which cannot always be derived from the path (e.g. fang/v2
// is package "fang", cmdguard/v4's package is "v4"). Aliased imports override
// the guess at resolution time.
var knownPackageNames = map[string]string{ //nolint:gochecknoglobals // static table of well-known import paths
	fangImportPath: "fang",
}

// skippedDirs are directory names never descended into: vendored code is not
// the consumer's own, testdata is fixtures, hidden dirs are tooling state.
var skippedDirs = map[string]bool{ //nolint:gochecknoglobals // static table
	"vendor":    true,
	"testdata":  true,
	"node_modules": true,
}

// sourceFile is one parsed Go file plus its resolved import aliases.
type sourceFile struct {
	// relPath is the file path relative to the scan root.
	relPath string
	// fset positions this file's nodes.
	fset *token.FileSet
	// file is the parsed AST (with comments).
	file *ast.File
	// imports maps local package name -> import path for every import
	// (including "." for dot imports).
	imports map[string]string
	// cmdguardMajor is "", "v1", "v2", "v3", "v4", ... — the major of the
	// cmdguard import this file carries, or "" when it imports none.
	cmdguardMajor string
	// lines caches the split source lines for suppression lookups.
	lines []string
}

// project is a fully analyzed directory tree.
type project struct {
	// root is the absolute scan root.
	root string
	// files are the parsed Go files, in walk order.
	files []sourceFile
	// importsCmdguard reports whether any file imports any cmdguard path.
	importsCmdguard bool
	// skipped records files that failed to parse (kept for diagnostics).
	skipped []string
}

// analyzeCache memoizes analyze results per absolute directory so the six
// rules can share one parse pass when driven through the per-rule SDK
// interface. Valid for the lifetime of the process: a long-lived consumer
// re-invoking with a mutated tree should use [analyze] via a fresh process or
// call [ClearCache].
var (
	analyzeCacheMu sync.Mutex
	analyzeCache   = map[string]*project{} //nolint:gochecknoglobals // process-lifetime memo, see ClearCache
)

// ClearCache drops memoized analysis results. Long-lived processes that
// re-lint a mutated tree must call this between runs.
func ClearCache() {
	analyzeCacheMu.Lock()
	defer analyzeCacheMu.Unlock()

	analyzeCache = map[string]*project{}
}

// analyze walks dir, parses every relevant Go file, and resolves import
// aliases. Results are memoized per absolute directory (see [ClearCache]).
func analyze(_ context.Context, dir string) (*project, error) {
	root, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	analyzeCacheMu.Lock()
	cached, ok := analyzeCache[root]
	analyzeCacheMu.Unlock()

	if ok {
		return cached, nil
	}

	proj := &project{root: root}

	fset := token.NewFileSet()

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() {
			if path != root && (skippedDirs[d.Name()] || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") {
			return nil
		}

		src, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		if isGenerated(src) {
			return nil
		}

		parsed, parseErr := parser.ParseFile(fset, path, src, parser.ParseComments|parser.SkipObjectResolution)
		if parseErr != nil {
			// A file that does not parse is someone else's error to report
			// (the compiler); skipping keeps the linter resilient on dirty trees.
			rel, _ := filepath.Rel(root, path)
			proj.skipped = append(proj.skipped, rel)

			return nil
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}

		sf := sourceFile{
			relPath:       rel,
			fset:          fset,
			file:          parsed,
			imports:       resolveImports(parsed),
			cmdguardMajor: "",
			lines:         strings.Split(string(src), "\n"),
		}

		if major, found := cmdguardImport(sf.imports); found {
			sf.cmdguardMajor = major
			proj.importsCmdguard = true
		}

		proj.files = append(proj.files, sf)

		return nil
	})
	if err != nil {
		return nil, err
	}

	analyzeCacheMu.Lock()
	analyzeCache[root] = proj
	analyzeCacheMu.Unlock()

	return proj, nil
}

// cmdguardImport reports whether any resolved import is a cmdguard path, and
// which core major it belongs to. Sub-module imports (major "") count as
// cmdguard usage without pinning a core major.
func cmdguardImport(imports map[string]string) (major string, found bool) {
	for _, importPath := range imports {
		if strings.HasPrefix(importPath, cmdguardImportPrefix) {
			return majorOf(importPath), true
		}
	}

	return "", false
}

// resolveImports maps each import's local package name to its path. Aliases
// win; otherwise the package name is guessed from the known table or the last
// path segment; dot imports map to "."; blank imports are dropped.
func resolveImports(file *ast.File) map[string]string {
	imports := map[string]string{}

	for _, spec := range file.Imports {
		path := strings.Trim(spec.Path.Value, `"`)

		var local string

		switch {
		case spec.Name == nil:
			if known, ok := knownPackageNames[path]; ok {
				local = known
			} else {
				local = guessPackageName(path)
			}
		case spec.Name.Name == "_":
			continue
		case spec.Name.Name == ".":
			local = "."
		default:
			local = spec.Name.Name
		}

		imports[local] = path
	}

	return imports
}

// guessPackageName derives a package name from an import path: the last
// segment, with a trailing major-version segment ("/v2") skipped so
// "example.com/foo/v2" guesses "foo".
func guessPackageName(path string) string {
	segments := strings.Split(path, "/")
	last := segments[len(segments)-1]

	if len(segments) > 1 && len(last) >= 2 && last[0] == 'v' && last[1] >= '1' && last[1] <= '9' {
		return segments[len(segments)-2]
	}

	return last
}

// majorOf extracts the major-version segment ("", "v1", "v2", ...) from a
// cmdguard import path. Sub-module paths ("cmdguard/spinner") inherit no
// major from the core module — they report "" which callers treat as
// "current-line module, not the frozen core".
func majorOf(importPath string) string {
	rest := strings.TrimPrefix(importPath, cmdguardImportPrefix)

	if rest == "" {
		return "v1"
	}

	if !strings.HasPrefix(rest, "/") {
		return ""
	}

	segment := strings.SplitN(strings.TrimPrefix(rest, "/"), "/", 2)[0]

	if len(segment) >= 2 && segment[0] == 'v' && segment[1] >= '1' && segment[1] <= '9' {
		return segment
	}

	return ""
}

// isGenerated reports whether a file's header carries the standard
// "Code generated ... DO NOT EDIT." marker.
func isGenerated(src []byte) bool {
	for i, line := range strings.Split(string(src), "\n") {
		if i >= 5 {
			break
		}

		if strings.Contains(line, "Code generated") && strings.Contains(line, "DO NOT EDIT") {
			return true
		}
	}

	return false
}
