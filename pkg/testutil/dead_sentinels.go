package testutil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"unicode"
)

// SentinelFinding describes a package-level sentinel error variable that is
// declared but never referenced outside its own declaration (and any
// allow-listed bookkeeping file).
type SentinelFinding struct {
	Name string // sentinel variable name
	File string // declaring file path as discovered under dir
	Dir  string // package directory that owns the declaration
}

// SentinelScan is the outcome of one FindDeadSentinels pass over a tree.
type SentinelScan struct {
	// Declarations counts sentinel-candidate vars found across the tree.
	// A count of zero for a tree that should own sentinels means the scan
	// (not the tree) is broken — assert it in the caller.
	Declarations int

	// Findings lists declared-but-unwired sentinels.
	Findings []SentinelFinding
}

// sentinelDecl is one sentinel-candidate var: where it was declared and at
// which position (valid within its own file's parse).
type sentinelDecl struct {
	file string
	pos  token.Pos
}

// sentinelScanner carries one FindDeadSentinels pass's accumulated state.
type sentinelScanner struct {
	allowFile    string
	declsByDir   map[string]map[string]sentinelDecl
	exportedRefs map[string]bool
}

// FindDeadSentinels scans every Go package under dir (skipping _test.go
// files and vendor/testdata/hidden trees) and reports package-level Err*/
// err* variables that are initialized by a call — or carry an explicit
// error type — but are never referenced outside their declaration and
// allowFile (a base file name whose whole purpose is sentinel bookkeeping,
// e.g. the core error-facts table; pass "" for none). Reference checks are
// package-scoped: unexported sentinels referenced only from another package
// are dead, and a sentinel referenced only from a test file never counts as
// wired. Exported Err* sentinels are additionally kept alive by any
// qualified identifier (pkg.ErrX) anywhere in the tree — another package in
// the module or an external consumer may be the only wire — so deleting
// one still requires an rg-level check of real consumers.
//
// The scan mirrors the 2026-11-22 dead-sentinel audit that removed five
// editor/nil sentinels the v5 scaffold had mechanically re-added: a
// sentinel must be wired (returned, wrapped, or classified) or deleted.
func FindDeadSentinels(dir, allowFile string) (SentinelScan, error) {
	scanner := &sentinelScanner{
		allowFile:    allowFile,
		declsByDir:   map[string]map[string]sentinelDecl{},
		exportedRefs: map[string]bool{},
	}

	if err := scanner.walk(dir); err != nil {
		return SentinelScan{Declarations: 0, Findings: nil}, fmt.Errorf("scanning %s: %w", dir, err)
	}

	findings, err := scanner.resolve()
	if err != nil {
		return SentinelScan{Declarations: 0, Findings: nil}, fmt.Errorf("resolving sentinels in %s: %w", dir, err)
	}

	return SentinelScan{
		Declarations: scanner.declarationCount(),
		Findings:     findings,
	}, nil
}

// walk parses every non-test Go file under dir, collecting sentinel
// declarations per package directory and exported selector names tree-wide.
func (s *sentinelScanner) walk(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walking %s: %w", path, walkErr)
		}

		if d.IsDir() {
			return s.visitDir(dir, path, d)
		}

		return s.visitFile(path)
	})
	if err != nil {
		return fmt.Errorf("walking %s: %w", dir, err)
	}

	return nil
}

// visitDir skips subtree walks that cannot own sentinel wiring. The walk
// root itself is never skipped.
func (s *sentinelScanner) visitDir(root, path string, d fs.DirEntry) error {
	if path == root || !skipSentinelTree(d) {
		return nil
	}

	return filepath.SkipDir
}

// visitFile parses one non-test Go file into the scanner state.
func (s *sentinelScanner) visitFile(path string) error {
	if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
		return nil
	}

	af, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", path, err)
	}

	dirPath := filepath.Dir(path)

	if s.declsByDir[dirPath] == nil {
		s.declsByDir[dirPath] = map[string]sentinelDecl{}
	}

	collectSentinelDecls(af, path, s.declsByDir[dirPath])
	collectExportedSelectorNames(af, s.exportedRefs)

	return nil
}

// declarationCount totals the sentinel candidates found across all packages.
func (s *sentinelScanner) declarationCount() int {
	total := 0

	for _, decls := range s.declsByDir {
		total += len(decls)
	}

	return total
}

// resolve walks the collected packages and returns the dead sentinels.
func (s *sentinelScanner) resolve() ([]SentinelFinding, error) {
	var findings []SentinelFinding

	for dirPath, decls := range s.declsByDir {
		dirFindings, err := s.resolveDir(dirPath, decls)
		if err != nil {
			return nil, fmt.Errorf("resolving sentinels in %s: %w", dirPath, err)
		}

		findings = append(findings, dirFindings...)
	}

	return findings, nil
}

// resolveDir reports one package directory's declared-but-unwired sentinels.
func (s *sentinelScanner) resolveDir(dirPath string, decls map[string]sentinelDecl) ([]SentinelFinding, error) {
	alive, err := wiredSentinels(dirPath, decls, s.allowFile)
	if err != nil {
		return nil, fmt.Errorf("scanning references in %s: %w", dirPath, err)
	}

	var findings []SentinelFinding

	for name, decl := range decls {
		if s.isWired(alive, name) {
			continue
		}

		findings = append(findings, SentinelFinding{
			Name: name,
			File: decl.file,
			Dir:  dirPath,
		})
	}

	return findings, nil
}

// isWired reports whether a sentinel is referenced in its own package (or
// via a qualified identifier anywhere in the tree, for exported sentinels).
func (s *sentinelScanner) isWired(alive map[string]bool, name string) bool {
	if alive[name] {
		return true
	}

	return isExportedSentinelName(name) && s.exportedRefs[name]
}

// skipSentinelTree reports whether a subtree cannot own sentinel wiring:
// vendored code, test fixtures, and hidden trees are all skipped. The walk
// root itself is exempt (checked by the caller).
func skipSentinelTree(d fs.DirEntry) bool {
	name := d.Name()

	return name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".")
}

// collectSentinelDecls records every sentinel-candidate var declared at
// package level in af (its declaring file and declaration position).
func collectSentinelDecls(af *ast.File, file string, decls map[string]sentinelDecl) {
	for _, d := range af.Decls {
		gd, ok := d.(*ast.GenDecl)

		if !ok || gd.Tok != token.VAR {
			continue
		}

		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)

			if !ok || len(vs.Names) != len(vs.Values) {
				continue
			}

			for i, name := range vs.Names {
				if isSentinelVarName(name.Name) && isSentinelInitializer(vs, vs.Values[i]) {
					decls[name.Name] = sentinelDecl{file: file, pos: name.Pos()}
				}
			}
		}
	}
}

// collectExportedSelectorNames records every exported ErrX identifier used
// in selector position (pkg.ErrX) anywhere in af — cross-package wiring for
// exported sentinels.
func collectExportedSelectorNames(af *ast.File, refs map[string]bool) {
	ast.Inspect(af, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)

		if !ok {
			return true
		}

		if isExportedSentinelName(sel.Sel.Name) {
			refs[sel.Sel.Name] = true
		}

		return true
	})
}

// wiredSentinels re-parses the package files in dirPath and marks every
// sentinel referenced outside its declaration and the allow-listed
// bookkeeping file as alive.
func wiredSentinels(dirPath string, decls map[string]sentinelDecl, allowFile string) (map[string]bool, error) {
	alive := make(map[string]bool, len(decls))

	files, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("globbing %s: %w", dirPath, err)
	}

	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || filepath.Base(f) == allowFile {
			continue
		}

		af, parseErr := parser.ParseFile(token.NewFileSet(), f, nil, 0)
		if parseErr != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, parseErr)
		}

		scanWiring(af, decls, f, alive)
	}

	return alive, nil
}

// scanWiring marks every tracked sentinel identifier in af — except the
// declaration itself — as alive.
func scanWiring(af *ast.File, decls map[string]sentinelDecl, file string, alive map[string]bool) {
	selfPos := selfDeclPositions(af, decls, file)

	ast.Inspect(af, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)

		if !ok {
			return true
		}

		if _, isSentinel := decls[id.Name]; !isSentinel {
			return true
		}

		if pos, isSelf := selfPos[id.Name]; isSelf && id.Pos() == pos {
			return true
		}

		alive[id.Name] = true

		return true
	})
}

// selfDeclPositions maps the sentinel names DECLARED in af (restricted to
// the tracked set owned by file) to their declaration positions, so the
// reference scan can exclude the declaration itself while still counting a
// genuine same-file use.
func selfDeclPositions(af *ast.File, decls map[string]sentinelDecl, file string) map[string]token.Pos {
	self := map[string]token.Pos{}

	for _, d := range af.Decls {
		gd, ok := d.(*ast.GenDecl)

		if !ok || gd.Tok != token.VAR {
			continue
		}

		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)

			if !ok {
				continue
			}

			for _, name := range vs.Names {
				if decl, tracked := decls[name.Name]; tracked && decl.file == file {
					self[name.Name] = name.Pos()
				}
			}
		}
	}

	return self
}

// isSentinelVarName reports whether name follows the sentinel conventions:
// an exported Err + exported-word name, or any unexported err name.
func isSentinelVarName(name string) bool {
	return isExportedSentinelName(name) || isUnexportedSentinelName(name)
}

// isExportedSentinelName reports whether name is an exported sentinel:
// Err followed by another exported word.
func isExportedSentinelName(name string) bool {
	if !strings.HasPrefix(name, "Err") {
		return false
	}

	rest := name[len("Err"):]

	return rest != "" && unicode.IsUpper(rune(rest[0]))
}

// isUnexportedSentinelName reports whether name is an unexported sentinel:
// err followed by anything.
func isUnexportedSentinelName(name string) bool {
	return strings.HasPrefix(name, "err") && len(name) > len("err")
}

// isSentinelInitializer reports whether the value spec is shaped like the
// sentinel declarations this tripwire guards: an explicit error type or a
// call-expression initializer.
func isSentinelInitializer(vs *ast.ValueSpec, value ast.Expr) bool {
	if id, ok := vs.Type.(*ast.Ident); ok && id.Name == "error" {
		return true
	}

	_, ok := value.(*ast.CallExpr)

	return ok
}
