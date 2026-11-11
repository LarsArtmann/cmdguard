#!/usr/bin/env bash
# cmdguard v4 → v5 codemod: rewrites import paths and applies the v5 renames
# to a consumer repository. Run from the consumer repo root:
#
#   bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh            # apply
#   bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh --dry-run  # preview
#
# What is rewritten automatically (identifier renames scoped to
# cmdguard-importing files — see set_renames below):
#   - import paths:  github.com/larsartmann/cmdguard/v4 → /v5
#   - import alias:  v4 "..." / v4 qualifier → v5
#   - CommandInfo    → CommandMetadata
#   - TypeHandler    → TypeCodec (incl. RegisterTypeHandler, TypeHandlerFunc)
#   - HuhRunner      → HuhPrompter
#
# What is flagged for manual follow-up (printed, not changed):
#   - cli.SetConfig(...)        — REMOVED in v5; use config-file options or
#                                 direct struct fields before NewCLI
#   - RegisterInScope(...)      — signature is generic in v5; provider now
#                                 returns (T, error) instead of (do.Injector)-shaped values
#
# The identifier renames (CommandInfo, TypeHandler, ...) apply ONLY to files
# that reference cmdguard: a consumer's own TypeHandler type, or a v4 alias
# for a DIFFERENT package, must never be renamed (2026-11-11 fleet
# incidents: non-cmdguard files were corrupted by the unscoped renames).
#
# After the run: go mod tidy && go build ./... && go test ./...

set -euo pipefail

DRY_RUN=0
for arg in "$@"; do
	case "$arg" in
	--dry-run) DRY_RUN=1 ;;
	*)
		echo "unknown argument: $arg (supported: --dry-run)" >&2
		exit 2
		;;
	esac
done

# cmdguard_go_files lists the Go files that reference cmdguard — the ONLY
# files the renames may touch. Import-path rewrites run repo-wide (a
# cmdguard/v4 string anywhere IS cmdguard), but identifier renames stay
# scoped: a consumer's own TypeHandler in a file that never imports
# cmdguard is theirs, not ours.
cmdguard_go_files() {
	local f
	while IFS= read -r f; do
		grep -ql 'larsartmann/cmdguard' "$f" && printf '%s\n' "$f"
	done < <(find . -type f -name '*.go' ! -path './.git/*' ! -path './vendor/*' 2>/dev/null)
}

apply() {
	if [ "$DRY_RUN" -eq 1 ]; then
		echo "  [dry-run] would rewrite: $*"
	else
		sed -i "$@" || echo "  sed failed on: $*" >&2
	fi
}

set_default() {
	local f
	while IFS= read -r f; do
		apply "s|github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4|github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5|g" "$f"
		apply "s|github.com/larsartmann/cmdguard/v4|github.com/larsartmann/cmdguard/v5|g" "$f"
	done < <(find . -type f \( -name '*.go' -o -name 'go.mod' \) ! -path './.git/*' ! -path './vendor/*' 2>/dev/null)
}

set_renames() {
	local f
	while IFS= read -r f; do
		apply 's/\bv4 "github.com\/larsartmann\/cmdguard/v5 "github.com\/larsartmann\/cmdguard/g' "$f"
		apply 's/\bv4\.\([A-Z]\)/v5.\1/g' "$f"
		apply 's/\bCommandInfo\b/CommandMetadata/g' "$f"
		apply 's/\bTypeHandlerFunc\b/TypeCodecFunc/g' "$f"
		apply 's/\bRegisterTypeHandler\b/RegisterTypeCodec/g' "$f"
		apply 's/\bTypeHandler\b/TypeCodec/g' "$f"
		apply 's/\bHuhRunner\b/HuhPrompter/g' "$f"
	done < <(cmdguard_go_files)
}

flag_manual() {
	echo "==> Manual follow-up required (these APIs changed shape in v5):"

	local found=0

	for f in $(find . -type f -name '*.go' ! -path './.git/*' ! -path './vendor/*' 2>/dev/null); do
		if grep -nq '\.SetConfig(' "$f"; then
			echo "  SetConfig removed: $f"
			grep -n '\.SetConfig(' "$f" | sed 's/^/    /'
			found=1
		fi

		if grep -nq 'RegisterInScope(' "$f"; then
			echo "  RegisterInScope is generic now (provider returns (T, error)): $f"
			grep -n 'RegisterInScope(' "$f" | sed 's/^/    /'
			found=1
		fi

		if grep -nq "cli\.HealthCheckWithContext(" "$f"; then
			echo "  CLI.HealthCheckWithContext removed (HealthCheck is ctx-aware): $f"
			grep -n "cli\.HealthCheckWithContext(" "$f" | sed "s/^/    /"
			found=1
		fi
	done

	if [ "$found" -eq 0 ]; then
		echo "  none — all rewrites were mechanical"
	fi
}

main() {
	echo "==> cmdguard v4 → v5 codemod${DRY_RUN:+ (dry run — no files changed)}"

	set_default
	set_renames
	flag_manual

	if [ "$DRY_RUN" -eq 1 ]; then
		echo "==> dry run complete. Re-run without --dry-run to apply."
		return
	fi

	echo "==> done. Next steps:"
	echo "    go mod edit -require github.com/larsartmann/cmdguard/v5@latest"
	echo "    go mod tidy && go build ./... && go test ./..."
}

main "$@"
