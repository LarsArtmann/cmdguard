package lint

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// crossFileIndex carries the package-level identifier facts of one directory
// group. Go scoping makes only package-level names visible across files, so
// the index deliberately tracks top-level declarations only: function-local
// variables can never be the cross-file gap, and leaking their names would
// blur distinct variables into false positives.
//
// It closes the file-scoped tracing gap of CG003/CG004/CG006: constructor in
// one file + misuse in another (e.g. `var cli, err = NewCLI(...)` in
// setup.go, `cli.SetVersion(...)` in main.go).
type crossFileIndex struct {
	// cliVars are package-level variable names assigned from a cmdguard
	// NewCLI call anywhere in the group.
	cliVars map[string]bool
	// execErrs are package-level variable names assigned from an Execute
	// call on a CLI variable.
	execErrs map[string]bool
	// constructorErrs are package-level variable names assigned from a
	// cmdguard constructor error result (NewCLI/NewCommand/NewParentCommand/
	// AddCommand).
	constructorErrs map[string]bool
}

// dirGroupOf keys a file into its name-tracing group: same directory AND
// same testness. Go compiles non-test and _test.go files of a directory into
// separate packages, so merging their namespaces would conflate variables
// that never coexist.
func dirGroupOf(file *sourceFile) string {
	dir := filepath.Dir(file.relPath)

	if strings.HasSuffix(file.relPath, "_test.go") {
		return dir + "\x00test"
	}

	return dir + "\x00prod"
}

// buildCrossFileIndex walks every file once and unions its package-level
// facts per directory group. Called once from analyze so all rules share the
// result without per-check recomputation or lazy-init races.
func buildCrossFileIndex(files []sourceFile) map[string]*crossFileIndex {
	index := map[string]*crossFileIndex{}

	groupFor := func(file *sourceFile) *crossFileIndex {
		key := dirGroupOf(file)

		group, ok := index[key]
		if !ok {
			group = &crossFileIndex{
				cliVars:         map[string]bool{},
				execErrs:        map[string]bool{},
				constructorErrs: map[string]bool{},
			}
			index[key] = group
		}

		return group
	}

	for i := range files {
		file := &files[i]
		group := groupFor(file)

		cmdguardPath := firstCmdguardPath(file.imports)

		file.topLevelVarCalls(func(sel *ast.SelectorExpr) bool {
			if cmdguardPath == "" || sel.Sel == nil {
				return false
			}

			if sel.Sel.Name == "NewCLI" || sel.Sel.Name == "Execute" || cmdguardConstructors()[sel.Sel.Name] {
				return file.selectorFrom(sel, cmdguardPath)
			}

			return false
		}, func(names []*ast.Ident, call *ast.CallExpr) {
			sel, _ := call.Fun.(*ast.SelectorExpr)
			if sel == nil || sel.Sel == nil {
				return
			}

			switch sel.Sel.Name {
			case "NewCLI":
				// var cli, err = NewCLI(...): cli is the value slot, err the error.
				addFirstIdent(group.cliVars, names)
				addLastIdent(group.constructorErrs, names)
			case "Execute":
				addLastIdent(group.execErrs, names)
			default:
				// NewCommand / NewParentCommand / AddCommand: the error is
				// the (last) result.
				addLastIdent(group.constructorErrs, names)
			}
		})
	}

	return index
}

// cliVarsFor returns the CLI variable names visible in file: its own
// file-scoped NewCLI assignments plus the group's package-level ones.
func (p *project) cliVarsFor(file *sourceFile) map[string]bool {
	vars := file.cliVars()

	if group := p.crossFile[dirGroupOf(file)]; group != nil {
		for name := range group.cliVars {
			vars[name] = true
		}
	}

	return vars
}

// execErrIdentsFor returns the error identifiers of Execute calls visible in
// file: its own file-scoped ones plus the group's package-level ones.
func (p *project) execErrIdentsFor(file *sourceFile, cliVars map[string]bool) map[string]bool {
	errs := collectExecuteErrIdents(file.file, cliVars)

	if group := p.crossFile[dirGroupOf(file)]; group != nil {
		for name := range group.execErrs {
			errs[name] = true
		}
	}

	return errs
}

// constructorErrIdentsFor returns the constructor-error identifiers visible
// in file from the group's package-level declarations (function-body locals
// stay with the body-local analysis of the check itself).
func (p *project) constructorErrIdentsFor(file *sourceFile) map[string]bool {
	if group := p.crossFile[dirGroupOf(file)]; group != nil {
		return group.constructorErrs
	}

	return nil
}

// cmdguardConstructors returns the constructor symbol set shared by the
// panic check and the cross-file index (fresh map per call; called once per
// analyze, so the allocation is irrelevant).
func cmdguardConstructors() map[string]bool {
	return map[string]bool{
		"NewCLI":           true,
		"NewCommand":       true,
		"NewParentCommand": true,
		"AddCommand":       true,
	}
}

// topLevelVarCalls visits package-level var declarations whose value is a
// selector call the match function accepts. For `var a, b = f()` (one call,
// many names) visit receives all names; for `var a = f(), b = g()` it
// receives each call with its aligned name. Declarations that are not calls
// are skipped.
func (f *sourceFile) topLevelVarCalls(
	match func(*ast.SelectorExpr) bool,
	visit func(names []*ast.Ident, call *ast.CallExpr),
) {
	for _, decl := range f.file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.VAR {
			continue
		}

		for _, spec := range gen.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, value := range valueSpec.Values {
				call, ok := value.(*ast.CallExpr)
				if !ok {
					continue
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok || !match(sel) {
					continue
				}

				visit(alignedNames(valueSpec, i), call)
			}
		}
	}
}

// alignedNames resolves which declared identifiers a spec's i-th value
// feeds: index-aligned when counts match, all names for a single
// multi-return call, and none when the shape is unclear.
func alignedNames(spec *ast.ValueSpec, i int) []*ast.Ident {
	switch {
	case len(spec.Values) == len(spec.Names):
		if i < len(spec.Names) {
			return spec.Names[i : i+1]
		}
	case len(spec.Values) == 1 && len(spec.Names) > 1:
		return spec.Names
	}

	return nil
}

// addFirstIdent records the leading identifier of names (the value slot of
// a (value, error) call) into the set, skipping blank identifiers.
func addFirstIdent(set map[string]bool, names []*ast.Ident) {
	if len(names) == 0 {
		return
	}

	first := names[0]
	if first != nil && first.Name != "_" {
		set[first.Name] = true
	}
}

// addLastIdent records the final identifier of names (the error slot of a
// (value, error) call) into the set, skipping blank identifiers.
func addLastIdent(set map[string]bool, names []*ast.Ident) {
	if len(names) == 0 {
		return
	}

	last := names[len(names)-1]
	if last != nil && last.Name != "_" {
		set[last.Name] = true
	}
}
