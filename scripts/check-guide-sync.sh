#!/usr/bin/env bash
# check-guide-sync.sh — docs-vs-code truth gate for the website (2026-11-11).
#
# The drift class this kills: the website teaching an obsolete major after the
# library moved on. The 2026-11-08 version of this guard only grepped
# content/docs/**/*.mdx and one data file, so 14 guide pages stayed on v4 while
# the hero's visible snippet was hand-written v4 and its copy button served v5.
#
# The real checks live in scripts/check-docs-api.py:
#   1. no stale major anywhere in website/src (mdx, astro, ts, mjs)
#   2. every v5.X / cli.X / sub-module symbol in the docs exists in source
#   3. every complete Go program embedded in the docs compiles
#   4. the hero snippet has a single source of truth
#
# Run it directly or through `nix run .#check-all`.
set -euo pipefail
cd "$(dirname "$0")/.."

exec python3 scripts/check-docs-api.py
