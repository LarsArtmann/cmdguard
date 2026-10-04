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

// visitCallAssignments walks every assignment in root whose right side is a
// call matching match, invoking visit with the call and the identifiers
// receiving its results. It handles both the pairwise form (`x := f()`,
// one LHS per RHS) and the multi-value form (`v, err := f()`, two LHS for
// one RHS): in the multi-value form all LHS identifiers belong to the call.
func visitCallAssignments(
	root ast.Node,
	match func(*ast.CallExpr) bool,
	visit func(call *ast.CallExpr, resultIdents []string),
) {
	ast.Inspect(root, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, rhs := range assign.Rhs {
			call, ok := rhs.(*ast.CallExpr)
			if !ok || !match(call) {
				continue
			}

			visit(call, lhsIdentsFor(assign, i))
		}

		return true
	})
}

// lhsIdentsFor returns the LHS identifiers that receive the result of the
// call at assign.Rhs[i].
func lhsIdentsFor(assign *ast.AssignStmt, rhsIndex int) []string {
	var idents []string

	if len(assign.Lhs) > len(assign.Rhs) {
		// Multi-value form: every LHS identifier receives one result of the
		// single call (e.g. `v, err := NewCLI(...)`).
		for _, lhs := range assign.Lhs {
			if ident, ok := lhs.(*ast.Ident); ok {
				idents = append(idents, ident.Name)
			}
		}

		return idents
	}

	if rhsIndex < len(assign.Lhs) {
		if ident, ok := assign.Lhs[rhsIndex].(*ast.Ident); ok {
			idents = append(idents, ident.Name)
		}
	}

	return idents
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

	visitCallAssignments(file, func(call *ast.CallExpr) bool {
		sel, ok := call.Fun.(*ast.SelectorExpr)

		return ok && sel.Sel != nil && sel.Sel.Name == "RootCommand"
	}, func(_ *ast.CallExpr, resultIdents []string) {
		for _, name := range resultIdents {
			id[name] = true
		}
	})

	return id
}

// cliVars returns identifiers assigned from `<cmdguard>.NewCLI(...)` calls in
// the file (the value identifier of each assignment). These are treated as
// cmdguard CLI values for the dataflow-lite rules CG004 and CG006.
func (f *sourceFile) cliVars() map[string]bool {
	vars := map[string]bool{}

	visitCallAssignments(f.file, func(call *ast.CallExpr) bool {
		return f.isSelectorCall(call, firstCmdguardPath(f.imports), "NewCLI")
	}, func(_ *ast.CallExpr, resultIdents []string) {
		if len(resultIdents) > 0 {
			vars[resultIdents[0]] = true
		}
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
// if-init form `if err := v4.AddCommand(...); err != nil`).
func (f *sourceFile) cmdguardConstructorErrIdents(body *ast.BlockStmt) map[string]bool {
	constructorNames := map[string]bool{
		"NewCLI":           true,
		"NewCommand":       true,
		"NewParentCommand": true,
		"AddCommand":       true,
	}
	cmdguardPath := firstCmdguardPath(f.imports)
	errs := map[string]bool{}

	if cmdguardPath == "" {
		return errs
	}

	visitCallAssignments(body, func(call *ast.CallExpr) bool {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || !constructorNames[sel.Sel.Name] {
			return false
		}

		return f.selectorFrom(sel, cmdguardPath)
	}, func(_ *ast.CallExpr, resultIdents []string) {
		// Error result is the last identifier (NewCLI/NewCommand return
		// (value, error); AddCommand returns just error, which is also last).
		if len(resultIdents) > 0 {
			errs[resultIdents[len(resultIdents)-1]] = true
		}
	})

	return errs
}

// collectExecuteErrIdents returns identifiers assigned from
// `<cliVar>.Execute(...)` calls for the given CLI variables.
func collectExecuteErrIdents(file *ast.File, cliVars map[string]bool) map[string]bool {
	errs := map[string]bool{}

	visitCallAssignments(file, func(call *ast.CallExpr) bool {
		for name := range cliVars {
			if isMethodCallOn(call, name, "Execute") {
				return true
			}
		}

		return false
	}, func(_ *ast.CallExpr, resultIdents []string) {
		for _, name := range resultIdents {
			errs[name] = true
		}
	})

	return errs
}

// isFmtPrintCall reports whether call is one of the fmt display functions
// (Print, Println, Printf, Fprint, Fprintln, Fprintf).
func isFmtPrintCall(file *sourceFile, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel == nil {
		return false
	}

	switch sel.Sel.Name {
	case "Print", "Println", "Printf", "Fprint", "Fprintln", "Fprintf":
		return file.selectorFrom(sel, "fmt")
	default:
		return false
	}
}

// condIsErrNotNil reports whether cond has the shape `IDENT != nil` or
// `nil != IDENT`, returning the identifier name when it does.
func condIsErrNotNil(cond ast.Expr) (string, bool) {
	bin, ok := cond.(*ast.BinaryExpr)
	if !ok || bin.Op != token.NEQ {
		return "", false
	}

	left, lok := bin.X.(*ast.Ident)
	right, rok := bin.Y.(*ast.Ident)

	if lok && rok {
		if right.Name == "nil" && left.Name != "nil" {
			return left.Name, true
		}

		if left.Name == "nil" && right.Name != "nil" {
			return right.Name, true
		}
	}

	return "", false
}
