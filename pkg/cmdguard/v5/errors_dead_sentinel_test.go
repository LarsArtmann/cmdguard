package v5

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// TestNoDeadSentinels fails when a package-level sentinel error is declared
// but never referenced outside its own declaration and the canonical facts
// table. The v5 scaffold mechanically re-added five editor/nil sentinels that
// v3 had already cut as dead weight; they sat in the public API and the
// generated ERROR_REFERENCE for weeks while no code path could ever produce
// them (found via `erraudit tree ./... --max` "unreferenced" badges,
// 2026-11-22). This tripwire keeps the declared sentinel set equal to the
// reachable one: a sentinel must be wired (returned, wrapped, or classified
// by non-facts code) or deleted.
func TestNoDeadSentinels(t *testing.T) {
	t.Parallel()

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("globbing package files: %v", err)
	}

	fset := token.NewFileSet()
	parsed := make(map[string]*ast.File, len(files))
	declPos := make(map[string]token.Pos)

	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		af, err := parser.ParseFile(fset, f, nil, 0)
		if err != nil {
			t.Fatalf("parsing %s: %v", f, err)
		}
		parsed[f] = af
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
					if !isSentinelVarName(name.Name) || !isSentinelInitializer(vs, vs.Values[i]) {
						continue
					}
					declPos[name.Name] = name.Pos()
				}
			}
		}
	}

	if len(declPos) == 0 {
		t.Fatal("no sentinel declarations found; the tripwire's declaration scan is broken")
	}

	alive := make(map[string]bool, len(declPos))
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || filepath.Base(f) == "errors_facts.go" {
			continue
		}
		af := parsed[f]
		ast.Inspect(af, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			if pos, isSentinel := declPos[id.Name]; isSentinel && id.Pos() != pos {
				alive[id.Name] = true
			}

			return true
		})
	}

	for name := range declPos {
		if !alive[name] {
			t.Errorf(
				"sentinel %s is declared but never referenced outside its declaration and the facts table; remove it or wire it",
				name,
			)
		}
	}
}

func isSentinelVarName(name string) bool {
	if strings.HasPrefix(name, "Err") {
		rest := name[len("Err"):]

		return rest != "" && unicode.IsUpper(rune(rest[0]))
	}

	if strings.HasPrefix(name, "err") {
		rest := name[len("err"):]

		return rest != ""
	}

	return false
}

func isSentinelInitializer(vs *ast.ValueSpec, value ast.Expr) bool {
	if id, ok := vs.Type.(*ast.Ident); ok && id.Name == "error" {
		return true
	}

	_, ok := value.(*ast.CallExpr)

	return ok
}
