// Package lint provides the cmdguard usage linter: static-analysis rules that
// detect anti-patterns and missed opportunities in codebases consuming cmdguard.
//
// The rules are derived from real consumer audits (timesheets deep dive,
// 2026-10-04) and encode the documented cmdguard execution contract:
//
//   - CG001 execute-bypass: handing cli.RootCommand() to fang.Execute bypasses
//     cli.Execute, silently disabling signal handling, graceful DI shutdown,
//     cleanup hooks, and the single-error-display contract.
//   - CG002 stale-major-import: importing a frozen cmdguard major (anything
//     below the current one) blocks all fixes and sub-module features.
//   - CG003 panic-on-constructor-error: panicking on NewCLI/NewCommand/
//     NewParentCommand/AddCommand errors re-introduces the panic cmdguard
//     removed by design.
//   - CG004 set-version-runtime: calling cli.SetVersion after construction
//     instead of passing WithCLIVersion to NewCLI.
//   - CG005 duplicate-version-options: combining WithCLIVersion with
//     WithFangOptions(fang.WithVersion(...)) produces duplicate fang options.
//   - CG006 execute-error-reprint: re-printing the error returned by
//     cli.Execute double-reports it; it is for exit-code mapping only.
//
// Suppressions use reason-bearing in-source directives:
//
//	//cmdguard-lint:ignore CG001 <reason>
//
// placed on the offending line or the line directly above it.
//
// The package is dependency-lean: it depends only on go-finding (finding
// model) and go-linter-sdk (rule/registry scaffolding). The BuildFlow
// integration lives in the sibling provider package; the standalone CLI in
// cmd/cmdguard-lint.
package lint
