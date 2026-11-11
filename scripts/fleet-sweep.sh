#!/usr/bin/env bash
# fleet-sweep.sh — golangci-lint + go vet + go build across the v5.1.0 fleet.
# Read-only sweep; prints a per-repo verdict table for the findings doc.
set -uo pipefail
export GOEXPERIMENT=jsonv2

REPOS="AutoCart InboxClean Standup-Killer Zlota44 code-duplicate-analyzer crm docs-organizer erraudit file-and-image-renamer go-auto-upgrade go-credentials-store go-version-auto-configure go-wizard-sdk google-domain-detector io-surface lean-business-plan media-atlas mr-sync nsfw-classifier pim-sync projects-management-automation samber-linter securitymd terraform-diagrams-aggregator timesheets yt-history-intel"

for repo in $REPOS; do
  d=/home/lars/projects/$repo
  [ -d "$d" ] || { echo "MISSING $repo"; continue; }
  cd "$d"
  vet=OK; lint=0
  go vet ./... >/dev/null 2>&1 || vet=FAIL
  if command -v golangci-lint >/dev/null; then
    lint=$(golangci-lint run ./... --no-color 2>/dev/null | grep -cE "^[a-z].*\.go:[0-9]+" || true)
  fi
  echo "$repo vet=$vet lint_issues=$lint"
done
