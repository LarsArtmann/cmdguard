package testutil

import (
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
	scan := SentinelScan{}
	declsByDir := map[string]map[string]sentinelDecl{}
	exportedRefs := map[string]bool{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if path != dir && skipSentinelTree(d) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		dirPath := filepath.Dir(path)
		af, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if parseErr != nil {
			return parseErr
		}

		if declsByDir[dirPath] == nil {
			declsByDir[dirPath] = map[string]sentinelDecl{}
		}
		collectSentinelDecls(af, path, declsByDir[dirPath])
		collectExportedSelectorNames(af, exportedRefs)
		return nil
	})
	if err != nil {
		return SentinelScan{}, err
	}

	for dirPath, decls := range declsByDir {
		scan.Declarations += len(decls)
		alive, err := wiredSentinels(dirPath, decls, allowFile)
		if err != nil {
			return SentinelScan{}, err
		}
		for name, decl := range decls {
			if alive[name] {
				continue
			}
			if isExportedSentinelName(name) && exportedRefs[name] {
				continue
			}
			scan.Findings = append(scan.Findings, SentinelFinding{
				Name: name,
				File: decl.file,
				Dir:  dirPath,
			})
		}
	}

	return scan, nil
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

// wiredSentinels re-parses the package files in dirPath and marks every
// sentinel referenced outside its declaration and the allow-listed
// bookkeeping file as alive.
func wiredSentinels(dirPath string, decls map[string]sentinelDecl, allowFile string) (map[string]bool, error) {
	alive := make(map[string]bool, len(decls))
	files, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || filepath.Base(f) == allowFile {
			continue
		}
		af, parseErr := parser.ParseFile(token.NewFileSet(), f, nil, 0)
		if parseErr != nil {
			return nil, parseErr
		}

		selfPos := selfDeclPositions(af, decls, f)
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

	return alive, nil
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
