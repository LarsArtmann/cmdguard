#!/usr/bin/env python3
"""Truth gate for the cmdguard website docs.

Checks, in order:
  1. No doc page still teaches an obsolete major (``v4.`` / ``cmdguard/v4``),
     except the changelog and the v4-to-v5 migration page where v4 is the subject.
  2. Every ``v5.X`` / ``cli.X`` / sub-module ``X.Y`` symbol mentioned in the docs
     exists in the source tree.
  3. Every full Go program embedded in the docs compiles against the local v5 API.

Exit code 0 means the docs match the code. Run from anywhere; paths resolve
relative to this file.
"""

from __future__ import annotations

import glob
import os
import re
import shutil
import subprocess
import sys
import tempfile

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
DOCS = os.path.join(ROOT, "website", "src", "content", "docs")
SITE_SRC = os.path.join(ROOT, "website", "src")

STALE_MAJOR = re.compile(r"\bv4\.[A-Za-z_]|cmdguard/v4\b")
ALLOWED_V4 = {
    "website/src/content/docs/changelog.mdx",
    "website/src/content/docs/guides/migrating-v4-to-v5.mdx",
}

# APIs that existed in v4 and were removed or reshaped in v5. They may be named
# only in the migration/changelog pages, where they are the subject.
REMOVED_V4_API = {"SetConfig", "SetVersion", "HealthCheckWithContext"}

SUBMODULES = {
    "schema": os.path.join(ROOT, "schema"),
    "glamour": os.path.join(ROOT, "glamour"),
    "prompts": os.path.join(ROOT, "prompts"),
    "spinner": os.path.join(ROOT, "spinner"),
    "telemetry": os.path.join(ROOT, "telemetry"),
    "flightrecorder": os.path.join(ROOT, "flightrecorder"),
}

failures: list[str] = []


def fail(message: str) -> None:
    failures.append(message)
    print(f"DRIFT: {message}")


def site_files() -> list[str]:
    patterns = ("*.mdx", "*.astro", "*.ts", "*.mjs")
    files: list[str] = []
    for dirpath, _, filenames in os.walk(SITE_SRC):
        for name in filenames:
            if name.endswith(patterns):
                files.append(os.path.join(dirpath, name))

    return files


def rel(path: str) -> str:
    return os.path.relpath(path, ROOT)


def check_stale_major() -> None:
    for path in site_files():
        text = open(path, encoding="utf-8").read()
        if rel(path) in ALLOWED_V4:
            continue
        for match in STALE_MAJOR.finditer(text):
            line = text.count("\n", 0, match.start()) + 1
            fail(f"{rel(path)}:{line} teaches an obsolete major: {match.group(0)!r}")


def check_hero_single_source() -> None:
    hero = os.path.join(SITE_SRC, "components", "HeroSection.astro")
    data = os.path.join(SITE_SRC, "data", "hero-code.ts")

    hero_text = open(hero, encoding="utf-8").read()
    data_text = open(data, encoding="utf-8").read()

    if "highlightGo(heroCode)" not in hero_text:
        fail("HeroSection.astro does not render heroCode (snippet must have one source)")

    if "v5.NewCLI" not in data_text:
        fail("hero-code.ts does not teach v5.NewCLI")

    if "cmdguard/v4" in hero_text or "cmdguard/v4" in data_text:
        fail("hero still points at cmdguard/v4")


def check_quick_start() -> None:
    quick_start = os.path.join(DOCS, "getting-started", "quick-start.mdx")
    installation = os.path.join(DOCS, "getting-started", "installation.mdx")

    for path, needle, label in [
        (quick_start, "cmdguard/v5", "quick-start.mdx does not import cmdguard/v5"),
        (quick_start, "NewBuilder", "quick-start.mdx does not teach the v5 Builder"),
        (quick_start, "ON by default", "quick-start.mdx does not mention the batteries default"),
        (installation, "go get github.com/larsartmann/cmdguard/v5", "installation.mdx does not install cmdguard/v5"),
    ]:
        if needle not in open(path, encoding="utf-8").read():
            fail(label)


def go_sources(directory: str, export: bool) -> set[str]:
    names: set[str] = set()
    for path in glob.glob(os.path.join(directory, "*.go")):
        if path.endswith("_test.go"):
            continue
        text = open(path, encoding="utf-8").read()
        if export:
            names |= set(re.findall(r"^func ([A-Z]\w*)", text, re.M))
            names |= set(re.findall(r"^type ([A-Z]\w*)", text, re.M))
            names |= set(re.findall(r"^var ([A-Z]\w*)", text, re.M))
            names |= set(re.findall(r"^const ([A-Z]\w*)", text, re.M))
            names |= set(re.findall(r"^\t([A-Z]\w*)\s+[A-Za-z_][\w.\[\]*]*\s*=", text, re.M))
            names |= set(re.findall(r"^\t([A-Z]\w*)\s*=", text, re.M))
        else:
            names |= set(re.findall(r"^func \([^)]*\) ([A-Z]\w*)", text, re.M))
    return names


def check_symbols() -> None:
    core_exports = go_sources(os.path.join(ROOT, "pkg", "cmdguard", "v5"), export=True)
    cli_methods = go_sources(os.path.join(ROOT, "pkg", "cmdguard", "v5"), export=False)
    module_exports = {
        name: go_sources(path, export=True) for name, path in SUBMODULES.items()
    }

    known: dict[str, set[str]] = {
        "v5": core_exports,
        "cli": cli_methods,
        **module_exports,
    }

    for path in glob.glob(os.path.join(DOCS, "**", "*.mdx"), recursive=True):
        for line_number, line in enumerate(open(path, encoding="utf-8"), start=1):
            for qualifier, symbol in re.findall(
                r"\b(v5|cli|schema|glamour|prompts|spinner|telemetry|flightrecorder)\.([A-Z]\w*)",
                line,
            ):
                if qualifier == "cli" and symbol == "X":
                    continue
                if rel(path) in ALLOWED_V4 and symbol in REMOVED_V4_API:
                    continue
                if symbol not in known[qualifier]:
                    fail(
                        f"{rel(path)}:{line_number} references unknown symbol "
                        f"{qualifier}.{symbol}"
                    )


def go_blocks() -> list[tuple[str, int, str]]:
    blocks: list[tuple[str, int, str]] = []
    for path in glob.glob(os.path.join(DOCS, "**", "*.mdx"), recursive=True):
        text = open(path, encoding="utf-8").read()
        for match in re.finditer(r"```go\n(.*?)```", text, re.S):
            code = match.group(1)
            if not code.lstrip().startswith("package main"):
                continue
            line = text.count("\n", 0, match.start()) + 1
            blocks.append((path, line, code))
    return blocks


def compile_samples() -> None:
    blocks = go_blocks()
    if not blocks:
        return

    go_bin = shutil.which("go")
    if go_bin is None:
        print("SKIP: go not on PATH, cannot compile doc samples")
        return

    work = tempfile.mkdtemp(prefix="cmdguard-docs-")
    try:
        module = os.path.join(work, "go.mod")
        requires = [
            "github.com/larsartmann/cmdguard/v5 v5.0.0",
            *[f"github.com/larsartmann/cmdguard/{name} v0.0.0" for name in SUBMODULES],
            "github.com/samber/do/v2 v2.1.0",
            "github.com/larsartmann/go-output v0.38.7",
            "github.com/larsartmann/samber-do-auditlog v0.12.1",
        ]
        replaces = ["github.com/larsartmann/cmdguard/v5 => " + ROOT] + [
            f"github.com/larsartmann/cmdguard/{name} => {path}"
            for name, path in SUBMODULES.items()
        ]

        with open(module, "w", encoding="utf-8") as handle:
            handle.write("module docscheck\n\ngo 1.27\n\nrequire (\n")
            for requirement in requires:
                handle.write(f"\t{requirement}\n")
            handle.write(")\n\nreplace (\n")
            for replacement in replaces:
                handle.write(f"\t{replacement}\n")
            handle.write(")\n")

        env = os.environ | {"GOWORK": "off", "GOEXPERIMENT": "jsonv2", "GOFLAGS": "-mod=mod"}
        tidy = subprocess.run(
            [go_bin, "mod", "tidy"], cwd=work, env=env, capture_output=True, text=True
        )
        if tidy.returncode != 0:
            fail(f"go mod tidy for doc samples failed:\n{tidy.stderr.strip()[:2000]}")
            return

        for index, (path, line, code) in enumerate(blocks):
            sample_dir = os.path.join(work, f"sample{index}")
            os.makedirs(sample_dir, exist_ok=True)
            with open(os.path.join(sample_dir, "main.go"), "w", encoding="utf-8") as handle:
                handle.write(code)

            build = subprocess.run(
                [go_bin, "build", "-o", os.devnull, "./sample" + str(index)],
                cwd=work,
                env=env,
                capture_output=True,
                text=True,
            )
            if build.returncode != 0:
                fail(
                    f"{rel(path)}:{line} Go sample does not compile:\n"
                    f"{build.stderr.strip()[:2000]}"
                )
    finally:
        shutil.rmtree(work, ignore_errors=True)


def main() -> int:
    check_stale_major()
    check_hero_single_source()
    check_quick_start()
    check_symbols()
    compile_samples()

    if failures:
        print()
        print(f"check-guide-sync FAILED with {len(failures)} finding(s).")
        return 1

    print("check-guide-sync OK: docs, API symbols, and doc samples match the code")
    return 0


if __name__ == "__main__":
    sys.exit(main())
