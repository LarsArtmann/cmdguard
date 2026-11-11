#!/usr/bin/env bash
# fleet-bump-v5.1.0.sh — bump every clean cmdguard consumer to v5.1.0.
# Skips: cmdguard itself, D6-gated staged worktrees, dirty trees, repos
# already on v5.1.0. Verifies each bump with go build + go vet.
set -uo pipefail
export GOEXPERIMENT=jsonv2

SKIP="cmdguard crm-stage-t13b cv-stage-crmfix BuildFlow-pr44"
DONE=0; FAILED=0; SKIPPED=0

for gomod in /home/lars/projects/*/go.mod; do
  repo=$(basename "$(dirname "$gomod")")
  case " $SKIP " in *" $repo "*) continue;; esac
  grep -q "larsartmann/cmdguard/v5 v" "$gomod" || continue
  grep -q "cmdguard/v5 v5.1.0" "$gomod" && { echo "SKIP $repo (already v5.1.0)"; SKIPPED=$((SKIPPED+1)); continue; }

  cd "$(dirname "$gomod")" || continue
  if [ -n "$(git status --porcelain 2>/dev/null | head -1)" ]; then
    echo "SKIP $repo (dirty tree — do not touch WIP)"; SKIPPED=$((SKIPPED+1)); continue
  fi

  branch=$(git branch --show-current)
  if go get github.com/larsartmann/cmdguard/v5@v5.1.0 >/dev/null 2>&1 && go mod tidy >/dev/null 2>&1; then
    if go build ./... >/dev/null 2>&1 && go vet ./... >/dev/null 2>&1; then
      if git diff --quiet go.mod go.sum 2>/dev/null; then
        echo "NOOP $repo (go get produced no change)"
      else
        git add go.mod go.sum
        git commit -q -m "chore(deps): bump cmdguard/v5 to v5.1.0

Latest release: strict positional-args options (WithStrictPositionalArgs/
WithArgsPassthrough) and a concurrency-safe Builder. Build + vet green.
Ref: LarsArtmann/cmdguard v5.1.0 (2026-11-11)." >/dev/null 2>&1 \
          && echo "OK $repo (branch $branch)" || echo "PARTIAL $repo (bumped, commit failed)"
        DONE=$((DONE+1))
      fi
    else
      echo "FAIL $repo (build/vet red after bump — reverting)"
      git checkout -- go.mod go.sum 2>/dev/null
      FAILED=$((FAILED+1))
    fi
  else
    echo "FAIL $repo (go get/tidy failed)"
    git checkout -- go.mod go.sum 2>/dev/null
    FAILED=$((FAILED+1))
  fi
done

echo ""
echo "bumped=$DONE failed=$FAILED skipped=$SKIPPED"
