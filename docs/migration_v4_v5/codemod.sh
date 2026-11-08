#!/usr/bin/env bash
# cmdguard v4 → v5 codemod: rewrites import paths and applies the v5 renames
# to a consumer repository. Run from the consumer repo root:
#
#   bash /path/to/cmdguard/docs/migration_v4_v5/codemod.sh
#
# What is rewritten automatically:
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
# After the run: go mod tidy && go build ./... && go test ./...

set -euo pipefail

set_default() {
	local src dst
	for f in $(find . -type f \( -name '*.go' -o -name 'go.mod' \) ! -path './.git/*' ! -path './vendor/*' 2>/dev/null); do
		sed -i 's|github.com/larsartmann/cmdguard/v4/pkg/cmdguard/v4|github.com/larsartmann/cmdguard/v5/pkg/cmdguard/v5|g' "$f"
		sed -i 's|github.com/larsartmann/cmdguard/v4|github.com/larsartmann/cmdguard/v5|g' "$f"
	done

	for f in $(find . -type f -name '*.go' ! -path './.git/*' ! -path './vendor/*' 2>/dev/null); do
		sed -i 's/\bv4 "github.com\/larsartmann\/cmdguard/v5 "github.com\/larsartmann\/cmdguard/g' "$f"
		sed -i 's/\bv4\.\([A-Z]\)/v5.\1/g' "$f"
		sed -i 's/\bCommandInfo\b/CommandMetadata/g' "$f"
		sed -i 's/\bTypeHandlerFunc\b/TypeCodecFunc/g' "$f"
		sed -i 's/\bRegisterTypeHandler\b/RegisterTypeCodec/g' "$f"
		sed -i 's/\bTypeHandler\b/TypeCodec/g' "$f"
		sed -i 's/\bHuhRunner\b/HuhPrompter/g' "$f"
	done
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
	echo "==> cmdguard v4 → v5 codemod"

	set_default
	flag_manual

	echo "==> done. Next steps:"
	echo "    go mod edit -require github.com/larsartmann/cmdguard/v5@latest"
	echo "    go mod tidy && go build ./... && go test ./..."
}

main "$@"
