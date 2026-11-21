#!/usr/bin/env bash
# fleet-bump.sh <version> — bump every clean cmdguard consumer to <version>.
# Skips: cmdguard itself, staged worktrees, dirty trees, repos already on
# <version>. Verifies each bump with go build + go vet.
set -uo pipefail
export GOEXPERIMENT=jsonv2

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
	echo "usage: fleet-bump.sh <version> (e.g. fleet-bump.sh v5.2.0)" >&2
	exit 2
fi
if ! printf '%s' "$VERSION" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$'; then
	echo "fleet-bump: '$VERSION' is not a vMAJOR.MINOR.PATCH tag" >&2
	exit 2
fi

SKIP="cmdguard crm-stage-t13b cv-stage-crmfix BuildFlow-pr44"
DONE=0
FAILED=0
SKIPPED=0

for gomod in /home/lars/projects/*/go.mod; do
	repo=$(basename "$(dirname "$gomod")")
	case " $SKIP " in *" $repo "*) continue ;; esac
	grep -q "larsartmann/cmdguard/v5 v" "$gomod" || continue
	grep -q "cmdguard/v5 $VERSION" "$gomod" && {
		echo "SKIP $repo (already $VERSION)"
		SKIPPED=$((SKIPPED + 1))
		continue
	}

	cd "$(dirname "$gomod")" || continue
	if [ -n "$(git status --porcelain 2>/dev/null | head -1)" ]; then
		echo "SKIP $repo (dirty tree — do not touch WIP)"
		SKIPPED=$((SKIPPED + 1))
		continue
	fi

	branch=$(git branch --show-current)
	if go get "github.com/larsartmann/cmdguard/v5@$VERSION" >/dev/null 2>&1 && go mod tidy >/dev/null 2>&1; then
		if go build ./... >/dev/null 2>&1 && go vet ./... >/dev/null 2>&1; then
			if git diff --quiet go.mod go.sum 2>/dev/null; then
				echo "NOOP $repo (go get produced no change)"
			else
				git add go.mod go.sum
				git commit -q -m "chore(deps): bump cmdguard/v5 to $VERSION

Build + vet green. Ref: LarsArtmann/cmdguard $VERSION." >/dev/null 2>&1 &&
					echo "OK $repo (branch $branch)" || echo "PARTIAL $repo (bumped, commit failed)"
				DONE=$((DONE + 1))
			fi
		else
			echo "FAIL $repo (build/vet red after bump — reverting)"
			git checkout -- go.mod go.sum 2>/dev/null
			FAILED=$((FAILED + 1))
		fi
	else
		echo "FAIL $repo (go get/tidy failed)"
		git checkout -- go.mod go.sum 2>/dev/null
		FAILED=$((FAILED + 1))
	fi
done

echo ""
echo "bumped=$DONE failed=$FAILED skipped=$SKIPPED"
