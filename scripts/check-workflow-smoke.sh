#!/usr/bin/env bash
# check-workflow-smoke.sh — bounded smoke for the GitHub Actions surface
# (the SU46 promotion of the ad-hoc pre-push workflow smoke).
#
# Two checks, no network writes:
#   1. go-run pins: every `go run <module>@<version>` in .github/workflows
#      must resolve, build, and answer with a usage banner (docs-gates and
#      brandid-lint style CLIs; dependabot cannot see these pins, so a
#      rotten pin would otherwise surface only as a red CI job).
#   2. permissions checklist: any job whose steps can push (git push,
#      release upload, firebase/docker deploy) must run under an explicit
#      workflow- or job-level `permissions:` block — the 2026-11-21
#      ci.yml incident class (Contents: read only silently blocking tags).
#
# Usage: scripts/check-workflow-smoke.sh   (exit 0 = smoke clean)
set -uo pipefail
cd "$(dirname "$0")/.."

PIN_TIMEOUT_S=120
failures=0

echo "=== workflow smoke: go-run pins ==="
pins=$(grep -rhoE 'go run [^ ]+@[^ ]+' .github/workflows/ | awk '{print $3}' | sort -u)
if [ -z "$pins" ]; then
	echo "no go run pins found (nothing to probe)"
fi

for pin in $pins; do
	out=$(GOFLAGS=-mod=mod timeout "$PIN_TIMEOUT_S" go run "$pin" -h 2>&1)
	code=$?
	if [ $code -eq 0 ] || printf '%s' "$out" | grep -q "Usage"; then
		echo "ok   $pin (resolves, answers usage)"
	else
		echo "FAIL $pin (exit $code, no usage banner; rotten pin?)"
		failures=$((failures + 1))
	fi
done

echo ""
echo "=== workflow smoke: permissions checklist ==="
for wf in .github/workflows/*.yml .github/workflows/*.yaml; do
	[ -f "$wf" ] || continue
	push_jobs=$(awk '
		/^  [A-Za-z0-9_-]+:$/ && $0 !~ /^   / { job = substr($1, 1, length($1)-1) }
		/git push|gh release|firebase deploy|docker push/ { print job }
	' "$wf" | sort -u)
	[ -z "$push_jobs" ] && continue
	for job in $push_jobs; do
		if grep -q "^permissions:" "$wf" || awk -v j="$job:" '
			index($0, j) == 7 { injob = 1; next }
			/^  [^ ]/ && injob { injob = 0 }
			injob && /^    permissions:/ { found = 1 }
			END { exit !found }
		' "$wf"; then
			echo "ok   $wf: $job runs under an explicit permissions block"
		else
			echo "FAIL $wf: $job can push but has no explicit permissions block"
			failures=$((failures + 1))
		fi
	done
done

echo ""
if [ "$failures" -gt 0 ]; then
	echo "workflow smoke FAILED — $failures finding(s)"
	exit 1
fi
echo "workflow smoke OK"
