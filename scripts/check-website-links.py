#!/usr/bin/env python3
"""Internal link and asset checker for the built cmdguard website.

Run after `pnpm build`, from the website directory or the repo root:

    python3 scripts/check-website-links.py

It verifies that every root-relative or relative href/src in dist/**.html
resolves to a file in dist, and that same-page fragment links resolve to an
existing id. External URLs are not fetched.
"""

from __future__ import annotations

import glob
import html
import os
import re
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DIST = os.path.join(ROOT, "website", "dist")

LINK = re.compile(r'(?:href|src)="([^"]+)"')
ID = re.compile(r'\bid="([^"]+)"')
EXTERNAL = ("http://", "https://", "//", "mailto:", "tel:", "data:", "javascript:")

failures: list[str] = []


def target_for(base: str, url: str) -> str | None:
    path = url.split("#", 1)[0].split("?", 1)[0]
    if not path:
        return None
    if path.startswith("/"):
        candidate = os.path.join(DIST, path.lstrip("/"))
    else:
        candidate = os.path.normpath(os.path.join(os.path.dirname(base), path))

    if candidate.endswith("/"):
        candidate = os.path.join(candidate, "index.html")
    if os.path.isdir(candidate):
        candidate = os.path.join(candidate, "index.html")
    if not os.path.splitext(candidate)[1]:
        for extension in ("/index.html", ".html"):
            if os.path.isfile(candidate + extension):
                return candidate + extension

    return candidate


def ids_of(path: str) -> set[str]:
    return set(ID.findall(open(path, encoding="utf-8").read()))


def main() -> int:
    pages = glob.glob(os.path.join(DIST, "**", "*.html"), recursive=True)
    if not pages:
        print(f"no built HTML under {DIST}; run pnpm build first")
        return 1

    for page in pages:
        text = open(page, encoding="utf-8").read()
        for raw in LINK.findall(text):
            url = html.unescape(raw)
            if url.startswith(EXTERNAL) or url == "#":
                continue

            target = target_for(page, url)
            if target is None:
                continue

            if not os.path.isfile(target):
                failures.append(f"{os.path.relpath(page, DIST)} -> missing {url}")
                continue

            fragment = url.split("#", 1)[1] if "#" in url else ""
            if fragment and target is not page and target.endswith(".html"):
                if fragment not in ids_of(target):
                    failures.append(
                        f"{os.path.relpath(page, DIST)} -> missing anchor #{fragment} in {url}"
                    )

    if failures:
        print(f"check-website-links FAILED with {len(failures)} finding(s):")
        for failure in failures[:60]:
            print(f"  {failure}")
        return 1

    print(f"check-website-links OK: {len(pages)} pages, no broken internal targets")
    return 0


if __name__ == "__main__":
    sys.exit(main())
