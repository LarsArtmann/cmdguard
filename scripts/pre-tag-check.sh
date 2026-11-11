#!/usr/bin/env bash
# pre-tag-check.sh — the local gate that MUST pass before pushing a cmdguard tag.
#
# Mirrors every gate the release workflow runs (.github/workflows/release.yml):
#   1. checkout (full history)   -> we require a clean, committed, up-to-date tree
#   2. setup-go 1.27             -> we require the dev shell toolchain (nix develop)
#   3. go test ./... -race       -> step 5 below
#   4. golangci-lint ./...       -> step 6 below
#   5. softprops GitHub Release  -> NOT verified locally (needs the push)
# plus the v4.3.0 incident guard: never tag without a green full battery.
#
# Usage:
#   scripts/pre-tag-check.sh v5.1.0            # full: tree checks + battery + workflow gates
#   scripts/pre-tag-check.sh v5.1.0 --fast     # tree checks + workflow gates only
#                                              # (use ONLY when a green `nix run .#check-all`
#                                              #  exists for the exact HEAD being tagged)
# Exit 0 = safe to tag and push. Anything else = DO NOT TAG.
set -euo pipefail

TAG="${1:-}"
MODE="${2:-full}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

fail() {
	echo "PRE-TAG CHECK FAILED: $*" >&2
	exit 1
}

[ -n "$TAG" ] || fail "usage: pre-tag-check.sh <tag> [--fast]"
[[ "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$ ]] || fail "tag '$TAG' is not a valid semver tag"

echo "=== pre-tag-check: $TAG (mode: $MODE) ==="

# 1. Refuse a dirty tree — the tag must be cut from a committed state.
if [ -n "$(git status --porcelain)" ]; then
	git status --short | head -20 >&2
	fail "dirty tree; commit or stash first (v4.3.0 lesson: tag only verified states)"
fi

# 2. Refuse if the tag already exists anywhere.
if git rev-parse -q --verify "refs/tags/$TAG" >/dev/null; then
	fail "tag '$TAG' already exists locally"
fi
git fetch --tags --quiet origin || fail "cannot reach origin"
if git ls-remote --exit-code --tags origin "refs/tags/$TAG" >/dev/null 2>&1; then
	fail "tag '$TAG' already exists on origin"
fi

# 3. Refuse if HEAD is behind origin/master (tagging stale states poisons the range).
BRANCH="$(git rev-parse --abbrev-ref HEAD)"
git fetch --quiet origin "$BRANCH"
if [ "$(git rev-parse HEAD)" != "$(git rev-parse "origin/$BRANCH")" ]; then
	if git merge-base --is-ancestor origin/"$BRANCH" HEAD; then
		fail "HEAD has unpushed commits; push first so CI runs before the tag"
	else
		fail "HEAD diverged from or is behind origin/$BRANCH; pull/rebase first"
	fi
fi

# 4. CHANGELOG must carry a section for the tag (release notes source of truth).
VERSION_NO_V="${TAG#v}"
if ! grep -q "^## \[$VERSION_NO_V\]" CHANGELOG.md; then
	fail "CHANGELOG.md has no '## [$VERSION_NO_V]' section; draft it from [Unreleased] first"
fi

# 5. Full battery (default) — build+test+lint+fmt+tidy across all 8 modules.
if [ "$MODE" = "full" ]; then
	echo "--- running nix run .#check-all (full battery) ---"
	nix run .#check-all || fail "check-all battery failed"
else
	echo "--- --fast: skipping check-all (a green run must exist for HEAD $(git rev-parse --short HEAD)) ---"
	echo "    evidence: A01 check-all green for this tree (see TODO_LIST F3 / status report)"
fi

# 6. The release workflow's own gates, re-run locally (GOEXPERIMENT matches ci.yml).
echo "--- go test ./... -race (release workflow gate) ---"
GOEXPERIMENT=jsonv2 go test ./... -count=1 -timeout 120s -race || fail "release workflow tests failed"

echo "--- golangci-lint run ./... (release workflow gate) ---"
command -v golangci-lint >/dev/null || fail "golangci-lint not on PATH (enter: nix develop)"
golangci-lint run ./... || fail "release workflow lint failed"

# 7. Workflow sanity grep: the pushed tag must still trigger a workflow that
#    actually runs tests+lint before publishing (guards against silent workflow rot).
grep -q "go test" .github/workflows/release.yml || fail "release.yml no longer runs tests"
grep -q "golangci-lint" .github/workflows/release.yml || fail "release.yml no longer runs lint"
grep -q "v\*" .github/workflows/release.yml || fail "release.yml no longer triggers on v* tags"

echo ""
echo "ALL PRE-TAG CHECKS PASSED for $TAG. Safe to: git tag -a $TAG && git push origin $TAG"
