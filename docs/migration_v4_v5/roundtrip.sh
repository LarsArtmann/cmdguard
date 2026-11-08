#!/usr/bin/env bash
# roundtrip.sh — proves the v4→v5 codemod end-to-end (Pareto plan R4.6/R4.7).
#
# Takes the committed _fixture consumer (a v4-shaped CLI), builds it against
# the PUBLISHED v4.2.1, runs the codemod, pins cmdguard/v5 v5.0.0, builds,
# runs the binary, and lints the result with cmdguard-lint (expect zero
# CG002 findings). CI-able; requires network for the module proxy.
#
# Usage: bash docs/migration_v4_v5/roundtrip.sh
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
fixture="$repo_root/docs/migration_v4_v5/roundtrip/_fixture/main.go"
codemod="$repo_root/docs/migration_v4_v5/codemod.sh"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

fail() {
	echo "ROUNDTRIP FAIL: $1"
	exit 1
}

# 0. The committed fixture must still be a v4 consumer (guards against an
#    accidental in-repo codemod run turning the fixture into v5).
grep -q "cmdguard/v4/pkg/cmdguard/v4" "$fixture" ||
	fail "committed fixture no longer imports v4 — restore the v4 shape"

cp "$fixture" "$work/main.go"
cd "$work"
go mod init fixture >/dev/null

# 1. Prove the v4 state actually builds against the published v4.2.1.
go get github.com/larsartmann/cmdguard/v4@v4.2.1 2>/dev/null | tail -0
go mod tidy >/dev/null 2>&1
go build -o fixture-v4 . || fail "fixture does not build against v4.2.1 (fixture is broken, not the codemod)"
./fixture-v4 run --name v4 | grep -q "run v4" || fail "v4 binary did not run"

# 2. Run the codemod (the exact script shipped to consumers).
bash "$codemod" >/dev/null

# 3. Pin v5.0.0, tidy, build, run.
# The codemod rewrites the module PATH mechanically but leaves the old
# version string (v5 v4.2.1 is an invalid pairing) — consumers pin the
# version themselves next; the script does the same.
sed -i 's|github.com/larsartmann/cmdguard/v5 v[0-9.]*|github.com/larsartmann/cmdguard/v5 v5.0.0|' go.mod
go get github.com/larsartmann/cmdguard/v5@v5.0.0 2>/dev/null | tail -0
go mod tidy >/dev/null 2>&1
! grep -q "cmdguard/v4" main.go || fail "codemod left a v4 reference behind"
go build -o fixture-v5 . || fail "codemod output does not compile"
./fixture-v5 run --name v5 | grep -q "run v5" || fail "v5 binary did not run"
./fixture-v5 describe | grep -q '"doctor"' || fail "v5 batteries default missing (no doctor in describe contract)"

# 4. Lint the result: the v5 tree must be CG002-clean.
lint_out="$(go run github.com/larsartmann/cmdguard/lint/cmd/cmdguard-lint@v0.3.0 lint --dir "$work" --output text 2>/dev/null || true)"
if echo "$lint_out" | grep -q "CG002"; then
	fail "cmdguard-lint still reports CG002 on the migrated fixture: $lint_out"
fi

echo "ROUNDTRIP OK: v4.2.1 fixture → codemod → v5.0.0 build + run + lint-clean"
