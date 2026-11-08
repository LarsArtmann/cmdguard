#!/usr/bin/env bash
# check-guide-sync.sh — docs-vs-website consistency guard (2026-11-08).
#
# The drift class this kills: docs/ teach the current major while website/
# still teach the previous one (the v5 release found website/quick-start.mdx
# and installation.mdx on v4 with pointer-flag code samples, 24 hours after
# docs/QUICKSTART.md had been rewritten to v5).
#
# Scope note: docs/QUICKSTART.md (full tutorial) and the website quick-start
# page (1-minute minimal walk) are DIFFERENT documents by design, so this is
# NOT a heading-diff — it asserts the version-consistency facts both must
# carry. Extend the mapped pairs below as more guides get ported (TUTORIAL,
# COMPARISON, ... per the Pareto plan R3.4).
set -euo pipefail
cd "$(dirname "$0")/.."

website="website/src/content/docs"
fail=0

viol() {
	echo "DRIFT: $1"
	fail=1
}

# 1. No stale module path in any website doc, except changelog/migration
#    history where the old path is the subject being described.
while IFS= read -r f; do
	case "$f" in
	*/changelog.mdx | */migrating-from-cobra.mdx | */api-reference.mdx) continue ;;
	esac
	if grep -q "cmdguard/v4" "$f"; then
		viol "$f references cmdguard/v4 (stale major; /v5 is current)"
	fi
done < <(find "$website" -name '*.mdx' 2>/dev/null)

# 2. The website quick-start must teach the v5 batteries default, so it
#    cannot silently regress to the opt-in story.
qs="$website/getting-started/quick-start.mdx"
[ -f "$qs" ] || viol "missing $qs"
grep -q "cmdguard/v5" "$qs" || viol "quick-start.mdx does not import cmdguard/v5"
grep -q "ON by default" "$qs" || viol "quick-start.mdx does not mention the v5 batteries default (must say batteries are ON by default)"
grep -q "NewBuilder" "$qs" || viol "quick-start.mdx does not teach the v5 Builder"

# 3. Website installation page must point at the current major.
inst="$website/getting-started/installation.mdx"
[ -f "$inst" ] || viol "missing $inst"
grep -q "go get github.com/larsartmann/cmdguard/v5" "$inst" || viol "installation.mdx does not install cmdguard/v5"
grep -qE "Go 1\.2[0-9]+ or later" "$inst" || viol "installation.mdx has no Go version floor"

# 4. docs/QUICKSTART.md must agree with the same v5 story.
grep -q "cmdguard/v5" docs/QUICKSTART.md || viol "docs/QUICKSTART.md does not import cmdguard/v5"
grep -q "NewBuilder" docs/QUICKSTART.md || viol "docs/QUICKSTART.md does not teach the v5 Builder"

# 5. The website hero snippet (hero-code.ts) must teach the current major —
#    it is the FIRST code a visitor sees (H4 guard).
hero="website/src/data/hero-code.ts"
[ -f "$hero" ] || viol "missing $hero"
if grep -q "cmdguard/v4\|v4\." "$hero"; then
	viol "$hero still teaches cmdguard v4"
fi
grep -q "v5.NewCLI" "$hero" || viol "$hero does not use v5.NewCLI"

if [ "$fail" -eq 1 ]; then
	echo
	echo "check-guide-sync FAILED — fix the DRIFT lines above."
	echo "Website teaches the previous major; port the page, then re-run."
	exit 1
fi

echo "check-guide-sync OK — docs and website agree on the current major"
