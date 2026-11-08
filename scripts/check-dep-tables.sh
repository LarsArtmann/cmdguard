#!/usr/bin/env bash
# check-dep-tables.sh — kills the version-drift class where AGENTS.md's
# dependency tables (Key Dependencies versions, the per-module dependency
# budget) silently diverge from the go.mod files (Pareto plan D42.1).
#
# Derived FROM go.mod (the truth), diffed AGAINST AGENTS.md (the claim).
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0

viol() {
	echo "DEP-DRIFT: $1"
	fail=1
}

# 1. Key Dependencies versions in AGENTS.md vs root go.mod.
for lib in "github.com/spf13/cobra" "github.com/samber/do/v2" "github.com/spf13/pflag" "github.com/larsartmann/go-output" "github.com/larsartmann/samber-do-auditlog"; do
	mod_ver="$(awk -v lib="$lib" '$1 == lib { print $2 }' go.mod | head -1)"
	if [ -z "$mod_ver" ]; then
		viol "$lib not found in root go.mod"
		continue
	fi
	if ! grep -qE "\`$lib\`.*\| *$mod_ver" AGENTS.md; then
		viol "AGENTS.md Key Dependencies does not list $lib at $mod_ver (go.mod truth)"
	fi
done

# 2. Dependency budget line: "core N, glamour N, ..." vs computed direct
#    requires per module (excludes indirect; includes sibling requires).
declare -A budget
budget_line="$(grep -o 'Dependency budget.*' AGENTS.md | head -1 || true)"
parse_budget() {
	local key="$1"
	local val
	val="$(echo "$budget_line" | grep -oE "$key [0-9]+" | grep -oE "[0-9]+" | head -1)"
	echo "$val"
}
budget[core]="$(parse_budget 'core')"
budget[glamour]="$(parse_budget 'glamour')"
budget[prompts]="$(parse_budget 'prompts')"
budget[spinner]="$(parse_budget 'spinner')"
budget[telemetry]="$(parse_budget 'telemetry')"
budget[flightrecorder]="$(parse_budget 'flightrecorder')"
budget[lint]="$(parse_budget 'lint')"
budget[schema]="$(parse_budget 'schema')"

if [ -z "${budget[core]:-}" ]; then
	viol "AGENTS.md dependency-budget line not found or not parseable"
fi

compute_direct() {
	local mod="$1"
	awk '/^require [(]/,/[)]$/' "$mod/go.mod" 2>/dev/null |
		grep -v '// indirect' | grep -v 'require (' | grep -cv '^)' || true
}

for mod in core glamour prompts spinner telemetry flightrecorder lint schema; do
	[ -n "${budget[$mod]:-}" ] || continue
	dir="."
	[ "$mod" != "core" ] && dir="$mod"
	actual="$(compute_direct "$dir")"
	if [ "$actual" != "${budget[$mod]}" ]; then
		viol "$mod: AGENTS.md budget says ${budget[$mod]} direct requires, go.mod has $actual"
	fi
done

if [ "$fail" -eq 1 ]; then
	echo
	echo "check-dep-tables FAILED — update AGENTS.md from the go.mod truth (or fix the go.mod)."
	exit 1
fi

echo "check-dep-tables OK — AGENTS.md dependency tables match go.mod"
