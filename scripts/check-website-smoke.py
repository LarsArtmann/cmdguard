#!/usr/bin/env python3
"""Live-domain smoke test for the deployed cmdguard website.

Run AFTER a deploy (the website workflow's deploy job, or a local
`firebase deploy`):

    python3 scripts/check-website-smoke.py
    python3 scripts/check-website-smoke.py --expect-commit <sha>

Asserts, with bounded retries per route (transient 5xx/429/timeouts only):
  - HTTP 200 on every route the site serves (landing, getting started,
    all guides, reference pages) and on sitemap.xml/robots.txt/llms.txt
  - a real 404 status for a path that must not exist
  - when --expect-commit is given: the served /deploy-meta.json carries
    THIS deploy's commit (post-deploy freshness — the served release must
    be the just-created one, not the previous one still riding out the
    CDN propagation window). Bounded freshness window, default 90s.

Exit code 0 means the deployment is verifiably alive and fresh.
"""

from __future__ import annotations

import argparse
import json
import sys
import time
import urllib.error
import urllib.request

BASE = "https://cmdguard.lars.software"

ROUTE_TIMEOUT_S = 10
ROUTE_ATTEMPTS = 3
RETRY_BACKOFF_S = 2
FRESHNESS_WINDOW_S = 90
FRESHNESS_POLL_S = 10

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
    "/deploy-meta.json",
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

TRANSIENT_STATUS = {429, 500, 502, 503, 504}


def fetch_once(path: str) -> tuple[int | None, str]:
    req = urllib.request.Request(
        BASE + path,
        method="GET",
        headers={"User-Agent": "cmdguard-website-smoke/1"},
    )
    try:
        with urllib.request.urlopen(req, timeout=ROUTE_TIMEOUT_S) as resp:
            return resp.status, resp.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as err:
        return err.code, ""
    except (urllib.error.URLError, TimeoutError, ConnectionError, OSError) as err:
        print(f"retry {path}: transient error: {err}")
        return None, ""


def fetch(path: str) -> tuple[int | None, str]:
    """GET with bounded retries; only transient failures are retried."""
    code: int | None = None
    body = ""
    for attempt in range(ROUTE_ATTEMPTS):
        code, body = fetch_once(path)
        if code is not None and code not in TRANSIENT_STATUS:
            return code, body
        if attempt < ROUTE_ATTEMPTS - 1:
            time.sleep(RETRY_BACKOFF_S * (attempt + 1))
    return code, body


def check_routes(failures: list[str]) -> int:
    ok_200 = 0

    for path in MUST_BE_200:
        code, _ = fetch(path)
        mark = "ok " if code == 200 else "FAIL"
        print(f"{mark} {code} {path}")
        if code == 200:
            ok_200 += 1
        else:
            failures.append(f"{path}: expected 200, got {code}")

    for path in MUST_BE_404:
        code, _ = fetch(path)
        mark = "ok " if code == 404 else "FAIL"
        print(f"{mark} {code} {path}")
        if code != 404:
            failures.append(f"{path}: expected 404, got {code}")

    return ok_200


def check_freshness(expect_commit: str, failures: list[str]) -> None:
    """Assert the CDN serves THIS deploy's commit within a bounded window."""
    deadline = time.monotonic() + FRESHNESS_WINDOW_S

    while True:
        code, body = fetch("/deploy-meta.json")
        if code == 200 and body:
            try:
                meta = json.loads(body)
                served = meta.get("commit", "")
            except json.JSONDecodeError:
                served = ""
            if served == expect_commit:
                print(f"ok  freshness: served commit matches {expect_commit}")
                return
            print(f"wait freshness: serving {served or 'unknown'}, expecting {expect_commit}")
        else:
            print(f"wait freshness: deploy-meta.json -> {code}")

        if time.monotonic() >= deadline:
            failures.append(
                f"freshness: after {FRESHNESS_WINDOW_S}s the CDN still does not "
                f"serve commit {expect_commit} (deploy-meta.json -> {code})"
            )
            return

        time.sleep(FRESHNESS_POLL_S)


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--expect-commit",
        help="commit SHA the just-finished deploy must serve (freshness check)",
    )
    args = parser.parse_args()

    failures: list[str] = []

    ok_200 = check_routes(failures)

    if args.expect_commit:
        check_freshness(args.expect_commit, failures)

    if failures:
        print()
        print(f"SMOKE FAILED — {len(failures)} failure(s), {ok_200} route(s) ok:")
        for line in failures:
            print(f"  - {line}")
        print(f"Debug: curl -I {BASE}/deploy-meta.json and the failing route above.")
        return 1

    print()
    print(f"website smoke OK — {len(MUST_BE_200)} routes 200, 404 path correct")
    return 0


if __name__ == "__main__":
    sys.exit(main())
