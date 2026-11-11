#!/usr/bin/env python3
"""Live-domain smoke test for the deployed cmdguard website.

Run AFTER a deploy (the website workflow's deploy job, or a local
`firebase deploy`):

    python3 scripts/check-website-smoke.sh

It asserts HTTP 200 on every route the site serves (landing, getting
started, all guides, reference pages), on sitemap.xml/robots.txt/llms.txt,
and asserts a real 404 status for a path that must not exist. Exit code 0
means the deployment is verifiably alive.
"""

from __future__ import annotations

import sys
import urllib.request
import urllib.error

BASE = "https://cmdguard.lars.software"

MUST_BE_200 = [
    "/",
    "/getting-started/installation/",
    "/getting-started/quick-start/",
    "/api-reference/",
    "/changelog/",
    "/contributing/",
    "/related-tools/",
    "/llms.txt",
    "/robots.txt",
    "/sitemap-index.xml",
] + [
    f"/guides/{slug}/"
    for slug in [
        "audit-log",
        "config-files",
        "custom-types",
        "dependency-injection",
        "doctor",
        "error-handling",
        "flow-context",
        "lifecycle",
        "lint",
        "machine-contract",
        "middleware",
        "migrating-from-cobra",
        "migrating-v4-to-v5",
        "plugins",
        "production-batteries",
        "rich-output",
        "schema",
        "shell-completion",
        "sub-modules",
        "type-safe-flags",
    ]
]

MUST_BE_404 = ["/this-page-must-not-exist/"]


def status_for(path: str) -> int:
    req = urllib.request.Request(
        BASE + path,
        method="GET",
        headers={"User-Agent": "cmdguard-website-smoke/1"},
    )
    try:
        with urllib.request.urlopen(req, timeout=20) as resp:
            return resp.status
    except urllib.error.HTTPError as err:
        return err.code


def main() -> int:
    failures: list[str] = []

    for path in MUST_BE_200:
        code = status_for(path)
        mark = "ok " if code == 200 else "FAIL"
        print(f"{mark} {code} {path}")
        if code != 200:
            failures.append(f"{path}: expected 200, got {code}")

    for path in MUST_BE_404:
        code = status_for(path)
        mark = "ok " if code == 404 else "FAIL"
        print(f"{mark} {code} {path}")
        if code != 404:
            failures.append(f"{path}: expected 404, got {code}")

    if failures:
        print()
        print(f"{len(failures)} smoke failure(s):")
        for line in failures:
            print(f"  - {line}")
        return 1

    print()
    print(f"website smoke OK — {len(MUST_BE_200)} routes 200, 404 path correct")
    return 0


if __name__ == "__main__":
    sys.exit(main())
