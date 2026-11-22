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

// FindDeadSentinels scans every Go package under dir (skipping _test.go
// files and vendor/testdata/hidden trees) and reports package-level Err*/
// err* variables that are initialized by a call — or carry an explicit
// error type — but are never referenced outside their declaration and
// allowFile (a base file name whose whole purpose is sentinel bookkeeping,
// e.g. the core error-facts table; pass "" for none). Reference checks run
// per package directory: a sentinel referenced only from another package is
// still dead. Test files never wire a sentinel.
//
// The scan mirrors the 2026-11-22 dead-sentinel audit that removed five
// editor/nil sentinels the v5 scaffold had mechanically re-added: a
// sentinel must be wired (returned, wrapped, or classified) or deleted.
func FindDeadSentinels(dir, allowFile string) (SentinelScan, error) {
	scan := SentinelScan{}
	declFileByDir := map[string]map[string]string{}

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return skipSentinelTree(d, path)
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		dirPath := filepath.Dir(path)
		fset := token.NewFileSet()
		af, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}

		if declFileByDir[dirPath] == nil {
			declFileByDir[dirPath] = map[string]string{}
		}
		collectSentinelDecls(af, path, declFileByDir[dirPath])
		return nil
	})
	if err != nil {
		return SentinelScan{}, err
	}

	for dirPath, declFile := range declFileByDir {
		scan.Declarations += len(declFile)
		alive, err := wiredSentinels(dirPath, declFile, allowFile)
		if err != nil {
			return SentinelScan{}, err
		}
		for name, file := range declFile {
			if !alive[name] {
				scan.Findings = append(scan.Findings, SentinelFinding{
					Name: name,
					File: file,
					Dir:  dirPath,
				})
			}
		}
	}

	return scan, nil
}

// skipSentinelTree skips subtree walks that cannot own sentinel wiring.
func skipSentinelTree(d fs.DirEntry, path string) error {
	name := d.Name()
	if name == "vendor" || name == "testdata" || strings.HasPrefix(name, ".") {
		return filepath.SkipDir
	}
	return nil
}

// collectSentinelDecls records the declaring file of every
// sentinel-candidate var declared at package level in af.
func collectSentinelDecls(af *ast.File, file string, declFile map[string]string) {
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
					declFile[name.Name] = file
				}
			}
		}
	}
}

// wiredSentinels re-parses the package files in dirPath and marks every
// sentinel referenced outside its declaration and the allow-listed
// bookkeeping file as alive.
func wiredSentinels(dirPath string, declFile map[string]string, allowFile string) (map[string]bool, error) {
	alive := make(map[string]bool, len(declFile))
	files, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || filepath.Base(f) == allowFile {
			continue
		}
		fset := token.NewFileSet()
		af, parseErr := parser.ParseFile(fset, f, nil, 0)
		if parseErr != nil {
			return nil, parseErr
		}
		ast.Inspect(af, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if _, isSentinel := declFile[id.Name]; isSentinel {
				alive[id.Name] = true
			}
			return true
		})
	}

	return alive, nil
}

// isSentinelVarName reports whether name follows the sentinel conventions:
// an exported Err + exported-word name, or any unexported err name.
func isSentinelVarName(name string) bool {
	if strings.HasPrefix(name, "Err") {
		rest := name[len("Err"):]
		return rest != "" && unicode.IsUpper(rune(rest[0]))
	}
	if strings.HasPrefix(name, "err") {
		return len(name) > len("err")
	}
	return false
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
