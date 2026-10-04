package lint

import (
	"go/ast"
	"go/token"
)

// isSelectorCall reports whether call's function expression is a selector
// `pkgName.Sel` where pkgName resolves to importPath in this file (a dot
// import makes every selector match).
func (f *sourceFile) isSelectorCall(call *ast.CallExpr, importPath, symbol string) bool {
	if call == nil || call.Fun == nil {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != symbol {
		return false
	}

	return f.selectorFrom(sel, importPath)
}

// selectorFrom reports whether sel's receiver expression resolves to the
// given import path in this file.
func (f *sourceFile) selectorFrom(sel *ast.SelectorExpr, importPath string) bool {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if resolved, ok := f.imports[ident.Name]; ok {
		return resolved == importPath
	}

	// A dot import places the package's symbols in file scope directly, so
	// any bare selector could be from it.
	if _, ok := f.imports["."]; ok {
		return true
	}

	return false
}

// isMethodCallOn reports whether call is `recvIdent.Sel(...)` where recvIdent
// has the given name. This is deliberately name-based (dataflow-lite): the
// caller guarantees recvIdent was traced to the relevant construction site.
func isMethodCallOn(call *ast.CallExpr, recvIdent, method string) bool {
	if call == nil || call.Fun == nil {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil || sel.Sel.Name != method {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)

	return ok && ident.Name == recvIdent
}

// callArgIsRootCommand reports whether any argument of call is a
// `x.RootCommand()` selector call or an identifier previously assigned from
// one (tracked in rootCommandIdents).
func (f *sourceFile) callArgIsRootCommand(call *ast.CallExpr, rootCommandIdents map[string]bool) bool {
	for _, arg := range call.Args {
		switch expr := arg.(type) {
		case *ast.CallExpr:
			if sel, ok := expr.Fun.(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "RootCommand" {
				return true
			}
		case *ast.Ident:
			if rootCommandIdents[expr.Name] {
				return true
			}
		}
	}

	return false
}

// collectRootCommandIdents returns the identifiers in the file assigned from
// `x.RootCommand()` calls (e.g. `root := cli.RootCommand()`).
func collectRootCommandIdents(file *ast.File) map[string]bool {
	id := map[string]bool{}

	ast.Inspect(file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}

		for i, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok || call.Fun == nil {
				continue
			}

			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel != nil && sel.Sel.Name == "RootCommand" {
				if ident, ok := assign.Lhs[i].(*ast.Ident); ok {
					id[ident.Name] = true
				}
			}
		}

		return true
	})

	return id
}

// cliVars returns identifiers assigned from `<cmdguard>.NewCLI(...)` calls in
// the file (the first LHS identifier of each assignment). These are treated
// as cmdguard CLI values for the dataflow-lite rules CG004 and CG006.
func (f *sourceFile) cliVars() map[string]bool {
	vars := map[string]bool{}

	ast.Inspect(f.file, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}

		for i, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}

			if !f.isSelectorCall(call, firstCmdguardPath(f.imports), "NewCLI") {
				continue
			}

			if ident, ok := assign.Lhs[i].(*ast.Ident); ok {
				vars[ident.Name] = true
			}
		}

		return true
	})

	return vars
}

// firstCmdguardPath returns one import path in the file that belongs to the
// cmdguard core module (major-versioned), preferring the current major. It
// exists so symbol matching (NewCLI, NewCommand, ...) resolves against the
// core package whichever major the file imports.
func firstCmdguardPath(imports map[string]string) string {
	fallback := ""

	for _, importPath := range imports {
		if majorOf(importPath) == CurrentMajor {
			return importPath
		}

		if fallback == "" && majorOf(importPath) != "" {
			fallback = importPath
		}
	}

	return fallback
}

// forFuncs invokes fn for every function body in the file (declarations and
// literals). Rules that reason about intra-function dataflow use this to keep
// their tracking function-scoped.
func forFuncs(file *ast.File, fn func(body *ast.BlockStmt)) {
	ast.Inspect(file, func(n ast.Node) bool {
		switch fnNode := n.(type) {
		case *ast.FuncDecl:
			if fnNode.Body != nil {
				fn(fnNode.Body)
			}
		case *ast.FuncLit:
			if fnNode.Body != nil {
				fn(fnNode.Body)
			}
		}

		return true
	})
}

// cmdguardConstructorErrIdents walks a function body and returns the names of
// identifiers that receive the error return of a cmdguard constructor call
// (`x, err := v4.NewCommand(...)`, `err := v4.AddCommand(...)`, and the
// if-init form `if err := v4.NewCLI(...); err != nil`).
func (f *sourceFile) cmdguardConstructorErrIdents(body *ast.BlockStmt) map[string]bool {
	constructorNames := map[string]bool{"NewCLI": true, "NewCommand": true, "NewParentCommand": true, "AddCommand": true}
	cmdguardPath := firstCmdguardPath(f.imports)
	errs := map[string]bool{}

	if cmdguardPath == "" {
		return errs
	}

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}

		for _, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel == nil || !constructorNames[sel.Sel.Name] {
				continue
			}

			if !f.selectorFrom(sel, cmdguardPath) {
				continue
			}

			// Error result is the last LHS identifier (NewCLI/NewCommand
			// return (value, error); AddCommand returns just error).
			if ident, ok := assign.Lhs[len(assign.Lhs)-1].(*ast.Ident); ok {
				errs[ident.Name] = true
			}
		}

		return true
	})

	return errs
}

// condIsErrNotNil reports whether cond has the shape `IDENT != nil` or
// `nil != IDENT`, returning the identifier name when it does.
func condIsErrNotNil(cond ast.Expr) (string, bool) {
	bin, ok := cond.(*ast.BinaryExpr)
	if !ok || bin.Op != token.NEQ {
		return "", false
	}

	left, lok := bin.X.(*ast.Ident)
	rightNil, rok := bin.Y.(*ast.Ident)

	if lok && rok && rightNil.Name == "nil" {
		return left.Name, true
	}

	leftNil, lok2 := bin.X.(*ast.Ident)
	right, rok2 := bin.Y.(*ast.Ident)

	if lok2 && rok2 && leftNil.Name == "nil" {
		return right.Name, true
	}

	return "", false
}
