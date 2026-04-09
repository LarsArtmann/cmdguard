# Migration to Nix Flakes — Detailed Proposal

**Project:** cmdguard v2.1.0
**Date:** 2026-04-09
**Status:** Proposal — Not Yet Implemented

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current State Analysis](#2-current-state-analysis)
3. [Why Nix Flakes](#3-why-nix-flakes)
4. [Proposed Architecture](#4-proposed-architecture)
5. [File-by-File Specification](#5-file-by-file-specification)
6. [Implementation Plan](#6-implementation-plan)
7. [Migration Strategy](#7-migration-strategy)
8. [Risk Assessment](#8-risk-assessment)
9. [Open Questions](#9-open-questions)
10. [References](#10-references)

---

## 1. Executive Summary

This proposal migrates cmdguard's development environment and CI pipeline to **Nix Flakes**, providing:

- **Reproducible dev environments** — one `nix develop` gives every contributor identical Go, golangci-lint, just, and tool versions across macOS and Linux.
- **Pinned toolchain** — Go 1.26, golangci-lint 2.x, and all tools are locked in `flake.lock`, eliminating "works on my machine".
- **Hermetic CI** — `nix flake check` replaces ad-hoc `actions/setup-go` + `golangci-lint-action` with a single declarative pipeline.
- **Declarative formatting** — `nix fmt` formats Nix files; Go formatting still delegates to `golangci-lint fmt`.
- **Zero breaking changes** — the justfile, go.mod, and all existing workflows remain fully functional.

### What This Is NOT

- Not a replacement for `go.mod` dependency management.
- Not building/publishing cmdguard as a Nix package (it's a Go library, not a binary).
- Not replacing the justfile — just becomes a devShell tool, not a CI dependency.

---

## 2. Current State Analysis

### 2.1 Toolchain Inventory

The following tools are required to build, test, lint, and develop cmdguard:

| Tool          | Version            | Source                           | Used By           |
| ------------- | ------------------ | -------------------------------- | ----------------- |
| Go            | 1.26 (go.mod)      | Manual / CI setup-go             | Build, test, run  |
| golangci-lint | 2.x (latest in CI) | Manual / CI golangci-lint-action | Lint, format      |
| just          | Latest (CI)        | Manual / CI setup-just           | Task runner       |
| golines       | Latest (via gci)   | Implicit via golangci-lint       | Code formatting   |
| gci           | Latest             | Implicit via golangci-lint       | Import ordering   |
| gofumpt       | Latest             | Implicit via golangci-lint       | Code formatting   |
| goimports     | Latest             | Implicit via golangci-lint       | Import management |

### 2.2 Build Targets (from justfile)

```
build          → go build ./...
build-version  → go build -ldflags "..." ./...
test           → go test ./...
test-v         → go test -v ./...
test-cover     → go test -cover ./...
test-race      → go test -race ./...
lint           → golangci-lint run --enable=errcheck ./...
fmt            → go fmt ./...
tidy           → go mod tidy
verify         → build + test + lint
run-basic      → go run ./examples/basic/main.go hello
run-advanced   → go run ./examples/advanced/main.go db migrate
run-guarded    → go run ./examples/guarded/main.go validate
clean          → go clean + rm artifacts
deps           → go mod graph
deps-list      → go list -m all
update         → go get -u + go mod tidy
dogfood        → Self-validation checks
```

### 2.3 CI Pipeline (from `.github/workflows/ci.yml`)

| Job        | Go Version(s)    | Steps                                                                   |
| ---------- | ---------------- | ----------------------------------------------------------------------- |
| `test`     | 1.24, 1.25, 1.26 | checkout → setup-go → cache → download → build → test → race → coverage |
| `lint`     | 1.26             | checkout → setup-go → golangci-lint-action (latest)                     |
| `verify`   | 1.26             | checkout → setup-go → setup-just → just verify                          |
| `examples` | 1.26             | checkout → setup-go → go run basic → go run typed                       |

### 2.4 Key Observations

1. **No existing Nix infrastructure** — zero `.envrc`, `shell.nix`, or `flake.nix` files exist.
2. **CI uses `latest` for golangci-lint** — `golangci-lint-action@v6` with `version: latest` means non-reproducible CI.
3. **Multi-version Go matrix** — CI tests against Go 1.24, 1.25, and 1.26.
4. **No shell scripts** — all task orchestration is in the justfile.
5. **Library, not binary** — `buildGoModule` would only be needed for examples; the library itself has no `main` package.
6. **Go 1.26 is cutting-edge** — may not be in stable nixpkgs yet; may require `nixpkgs-unstable` or an overlay.
7. **golangci-lint v2** — uses config version `"2"` format with `formatters` section.
8. **Build tags** — `.golangci.yml` uses experimental build tags: `goexperiment.goroutineleakprofile`, `goexperiment.jsonv2`, `goexperiment.simd`.

---

## 3. Why Nix Flakes

### 3.1 Problems Solved

| Problem                            | Current State                                      | With Nix Flakes                        |
| ---------------------------------- | -------------------------------------------------- | -------------------------------------- |
| Version drift between developers   | "Install Go 1.26, golangci-lint 2.x, just"         | `nix develop` — exact versions, locked |
| CI uses `latest` for golangci-lint | Builds can break when upstream releases changes    | Pinned in flake.lock                   |
| Onboarding friction                | Multiple manual tool installations                 | Install Nix → `nix develop` → done     |
| macOS vs Linux inconsistencies     | Different package managers, different versions     | Same derivation on both platforms      |
| No unified `check` command         | Separate CI jobs with different setups             | `nix flake check` — all checks in one  |
| Formatter not enforced             | `go fmt` manual; golangci-lint formatters separate | `nix fmt` for Nix; flake check for Go  |

### 3.2 What We Keep

- **justfile** — remains the primary task runner; available inside devShell
- **go.mod / go.sum** — remains the Go dependency manager
- **.golangci.yml** — remains the lint config
- **GitHub Actions CI** — enhanced with Nix, not replaced (see §7)

### 3.3 Trade-offs

| Benefit                          | Cost                                      |
| -------------------------------- | ----------------------------------------- |
| Reproducible environments        | Nix learning curve for contributors       |
| Pinned tooling                   | `flake.lock` needs periodic updates       |
| `nix flake check` for everything | Longer initial `nix develop` (downloads)  |
| Cross-platform consistency       | Go 1.26 may need nixpkgs-unstable channel |

---

## 4. Proposed Architecture

### 4.1 File Structure

```
cmdguard/
├── flake.nix                    # Main flake definition
├── flake.lock                   # Pinned dependency versions (generated)
├── .envrc                       # direnv auto-shell activation
├── nix/
│   └── checks.nix               # Custom check derivations (build, test, lint)
├── justfile                     # Unchanged — available in devShell
├── go.mod                       # Unchanged
├── go.sum                       # Unchanged
├── .golangci.yml                # Unchanged
└── .github/workflows/ci.yml     # Enhanced (see §7.2)
```

### 4.2 Flake Outputs

```
flake.nix
├── inputs
│   ├── nixpkgs (nixos-unstable for Go 1.26)
│   ├── flake-utils (multi-system support)
│   └── treefmt-nix (optional: Nix formatting)
│
└── outputs
    ├── devShells.${system}.default    # Go 1.26 + golangci-lint + just + gopls
    ├── checks.${system}               # build, test, lint, gofmt
    ├── formatter.${system}            # nixfmt-rfc-style for .nix files
    └── packages.${system}             # (optional: example binaries)
```

### 4.3 Design Decisions

| Decision                 | Choice                | Rationale                                            |
| ------------------------ | --------------------- | ---------------------------------------------------- |
| Flake framework          | Raw flake-utils       | Simple; no need for flake-parts complexity           |
| Go version source        | nixpkgs-unstable      | Go 1.26 may not be in stable nixpkgs                 |
| Library packaging        | Skip                  | cmdguard is a Go library, not a distributable binary |
| Example binary packaging | Optional (Phase 2)    | Nice-to-have; not critical for dev workflow          |
| Vendor hash management   | N/A (library)         | Only needed if building binaries with buildGoModule  |
| Nix formatter            | nixfmt-rfc-style      | Official Nix formatting standard                     |
| Go formatting via Nix    | Via golangci-lint fmt | Respect existing .golangci.yml formatter config      |
| direnv integration       | Yes                   | Seamless shell activation on cd                      |
| Replace CI or augment    | Augment               | Keep GitHub Actions; add Nix-based CI option         |

---

## 5. File-by-File Specification

### 5.1 `flake.nix`

```nix
{
  description = "cmdguard — CLI Guard Library for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix.url = "github:numtide/treefmt-nix";
  };

  outputs =
    { self, nixpkgs, flake-utils, treefmt-nix }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        goVersion = 26; # Go 1.26

        # Use Go 1.26 from nixpkgs
        goPkg = pkgs."go_1_${toString goVersion}" or pkgs.go;

        # Build inputs shared across devShell and checks
        buildInputs = with pkgs; [ ];

        nativeBuildInputs = with pkgs; [
          goPkg
          golangci-lint
          just
          gopls
          gotools
        ];

        # treefmt for Nix file formatting
        treefmtEval = treefmt-nix.lib.evalModule pkgs {
          projectRootFile = "flake.nix";
          programs.nixfmt.enable = true; # nixfmt-rfc-style
        };

      in
      {
        # ── Development Shell ──────────────────────────────────
        devShells.default = pkgs.mkShell {
          name = "cmdguard-dev";

          inherit nativeBuildInputs buildInputs;

          shellHook = ''
            echo "cmdguard dev shell"
            echo "Go:    $(go version)"
            echo "Lint:  $(golangci-lint version)"
            echo "Just:  $(just --version)"
          '';
        };

        # ── Checks (nix flake check) ──────────────────────────
        checks = {
          # Verify the project compiles
          build = pkgs.runCommand "cmdguard-build" {
            nativeBuildInputs = [ goPkg ];
            src = self;
          } ''
            cp -r $src/* .
            chmod -R u+w .
            go build ./...
            touch $out
          '';

          # Run all tests with race detector
          test = pkgs.runCommand "cmdguard-test" {
            nativeBuildInputs = [ goPkg ];
            src = self;
          } ''
            cp -r $src/* .
            chmod -R u+w .
            go test -race -count=1 -timeout 120s ./...
            touch $out
          '';

          # Run golangci-lint
          lint = pkgs.runCommand "cmdguard-lint" {
            nativeBuildInputs = [ goPkg golangci-lint ];
            src = self;
          } ''
            cp -r $src/* .
            chmod -R u+w .
            golangci-lint run --enable=errcheck ./...
            touch $out
          '';

          # Verify Go formatting (gofumpt via golangci-lint fmt)
          gofmt = pkgs.runCommand "cmdguard-gofmt" {
            nativeBuildInputs = [ goPkg golangci-lint ];
            src = self;
          } ''
            cp -r $src/* .
            chmod -R u+w .
            export HOME=$(mktemp -d)
            golangci-lint fmt ./...
            if ! git diff --quiet; then
              echo "ERROR: Go files are not formatted. Run 'golangci-lint fmt ./...'"
              git diff
              exit 1
            fi
            touch $out
          '';

          # treefmt check for Nix files
          formatting = treefmtEval.config.build.check self;
        };

        # ── Formatter (nix fmt) ────────────────────────────────
        formatter = treefmtEval.config.build.wrapper;
      }
    );
}
```

#### Key Design Notes

- **`runCommand` pattern** — Each check copies the source, makes it writable, runs the tool, and creates `$out` on success. This is the standard Nixpkgs pattern for check derivations.
- **No `buildGoModule`** — Since cmdguard is a library with no `main` package, we use `go build ./...` directly in check derivations. `buildGoModule` would require a `vendorHash` and `subPackages`, which don't apply here.
- **Multi-system via `eachDefaultSystem`** — Supports `x86_64-linux`, `aarch64-linux`, `x86_64-darwin` (Intel Mac), `aarch64-darwin` (Apple Silicon).
- **`goPkg` fallback** — Uses `pkgs.go_1_26` if available, falls back to `pkgs.go`. This handles the case where Go 1.26 hasn't landed in nixpkgs yet.
- **golangci-lint formatting** — The `gofmt` check uses `golangci-lint fmt` which respects `.golangci.yml`'s `formatters` section (gofumpt, gci, goimports, golines), not just plain `go fmt`.
- **`HOME` in gofmt check** — golangci-lint needs a writable home directory for cache; we provide a temp one.

### 5.2 `.envrc`

```bash
# Enable nix integrations for automatic shell activation
# Requires: direnv (https://direnv.net/) and nix-direnv
use flake
```

#### Setup Requirements

```bash
# Install direnv (macOS)
brew install direnv

# Add to shell rc (~/.zshrc or ~/.bashrc)
eval "$(direnv hook zsh)"  # or bash

# Install nix-direnv for persistent GC roots
nix profile install nixpkgs#nix-direnv

# Then in project root
direnv allow
```

### 5.3 `nix/checks.nix` (Optional Extract)

For projects that grow beyond a single `flake.nix`, extract checks into a separate file:

```nix
# nix/checks.nix
# This file is OPTIONAL — only extract if flake.nix becomes too large.
# All checks can remain inline in flake.nix initially.

{ pkgs, self }:
{
  build = pkgs.runCommand "cmdguard-build" {
    nativeBuildInputs = [ pkgs.go ];
    src = self;
  } ''
    cp -r $src/* .
    chmod -R u+w .
    go build ./...
    touch $out
  '';

  test = pkgs.runCommand "cmdguard-test" {
    nativeBuildInputs = [ pkgs.go ];
    src = self;
  } ''
    cp -r $src/* .
    chmod -R u+w .
    go test -race -count=1 -timeout 120s ./...
    touch $out
  '';
}
```

Usage in `flake.nix`:

```nix
checks = import ./nix/checks.nix { inherit pkgs self; };
```

---

## 6. Implementation Plan

### Phase 1: Foundation (Minimal Viable Flake)

**Goal:** Get `nix develop` working with the correct toolchain.

| Step | Action                                                                | Verification                                |
| ---- | --------------------------------------------------------------------- | ------------------------------------------- |
| 1.1  | Create `flake.nix` with devShell only (no checks yet)                 | `nix develop` drops into shell with Go 1.26 |
| 1.2  | Create `.envrc` with `use flake`                                      | `direnv allow` activates shell on `cd`      |
| 1.3  | Add `flake.lock` to git                                               | `git add flake.nix flake.lock .envrc`       |
| 1.4  | Verify `go build ./...` works inside devShell                         | `nix develop -c go build ./...`             |
| 1.5  | Verify `just verify` works inside devShell                            | `nix develop -c just verify`                |
| 1.6  | Add `.direnv/` and `.envrc` to `.gitignore` if `.direnv/` not already | Check .gitignore                            |

**Estimated time:** 1–2 hours

**Potential blocker:** Go 1.26 availability in nixpkgs. Mitigation: use `nixpkgs-unstable` channel or a Go overlay.

### Phase 2: Checks

**Goal:** `nix flake check` runs build + test + lint + formatting.

| Step | Action                                        | Verification                                 |
| ---- | --------------------------------------------- | -------------------------------------------- |
| 2.1  | Add `build` check to `flake.nix`              | `nix flake check -L` passes build            |
| 2.2  | Add `test` check (with race detector)         | `nix flake check -L` passes tests            |
| 2.3  | Add `lint` check (golangci-lint)              | `nix flake check -L` passes lint             |
| 2.4  | Add `gofmt` check (golangci-lint fmt dry-run) | `nix flake check -L` passes formatting check |
| 2.5  | Add `formatting` check (treefmt for .nix)     | `nix flake check -L` passes Nix formatting   |
| 2.6  | Verify all checks pass together               | `nix flake check -L` — all green             |

**Estimated time:** 1–2 hours

**Potential blocker:** golangci-lint version mismatch between nixpkgs and CI. Mitigation: check nixpkgs version, override if needed.

### Phase 3: Formatter

**Goal:** `nix fmt` formats Nix files; `golangci-lint fmt` remains for Go files.

| Step | Action                                         | Verification                               |
| ---- | ---------------------------------------------- | ------------------------------------------ |
| 3.1  | Configure treefmt with `nixfmt` in `flake.nix` | `nix fmt` formats `.nix` files             |
| 3.2  | Run `nix fmt` and verify output                | `nix fmt && git diff —-stat` shows changes |
| 3.3  | Commit formatted Nix files                     | `git commit`                               |

**Estimated time:** 30 minutes

### Phase 4: CI Integration

**Goal:** Add Nix-based CI job alongside existing GitHub Actions.

| Step | Action                                                                       | Verification                 |
| ---- | ---------------------------------------------------------------------------- | ---------------------------- |
| 4.1  | Add `nix-flake-check` job to `.github/workflows/ci.yml`                      | CI passes with new job       |
| 4.2  | Use `DeterminantSystems/nix-installer-action` or `cachix/install-nix-action` | Nix installs in CI           |
| 4.3  | Enable Cachix (optional) for CI caching                                      | Subsequent CI runs faster    |
| 4.4  | Keep existing jobs (test matrix, lint, verify, examples) unchanged           | All existing jobs still pass |

**Estimated time:** 2–3 hours

### Phase 5: Polish

| Step | Action                                                          | Verification                            |
| ---- | --------------------------------------------------------------- | --------------------------------------- |
| 5.1  | Update `.gitignore` with `.direnv/`                             | Nix artifacts not tracked               |
| 5.2  | Update `AGENTS.md` with Nix commands                            | Documentation reflects new workflow     |
| 5.3  | Update `CONTRIBUTING.md` with Nix setup instructions            | Contributors know how to use Nix        |
| 5.4  | Add `nix develop` as recommended onboarding path in `README.md` | New contributors see Nix as first-class |
| 5.5  | Verify `nix flake check` passes on both macOS and Linux         | Cross-platform validation               |

**Estimated time:** 1–2 hours

---

## 7. Migration Strategy

### 7.1 Non-Breaking Approach

The migration is **additive** — no existing files are modified (except `.gitignore` and docs):

```
ADDED:   flake.nix, flake.lock, .envrc
MODIFIED: .gitignore (add .direnv/), AGENTS.md, CONTRIBUTING.md
KEPT:    justfile, go.mod, go.sum, .golangci.yml, all .go files
```

Contributors who don't use Nix continue to:

- Install Go, golangci-lint, and just manually
- Use `just verify`, `just test`, `just lint`
- Run CI via existing GitHub Actions

### 7.2 CI Strategy

Two options, in order of preference:

#### Option A: Nix CI alongside existing CI (Recommended)

```yaml
# Add to existing ci.yml as a new job
nix-check:
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: DeterminantSystems/nix-installer-action@main
    - uses: DeterminantSystems/magic-nix-cache-action@main
    - run: nix flake check -L
```

This adds Nix validation without removing existing matrix jobs. Both must pass.

#### Option B: Replace existing CI with Nix

Replace `test`, `lint`, and `verify` jobs with a single `nix flake check`. Keep `examples` job as-is. This is more aggressive but simpler.

**Recommendation:** Start with Option A. Migrate to Option B once Nix CI is proven stable.

### 7.3 Go Version Handling

The CI matrix currently tests Go 1.24, 1.25, and 1.26. With Nix:

- The devShell and `nix flake check` use **Go 1.26** (primary).
- The multi-version matrix stays in GitHub Actions (non-Nix job).
- If multi-version Nix testing is desired later, add a flake input or overlay for each Go version.

### 7.4 Vendor Hash Considerations

Since cmdguard is a **library** (no `main` package), `buildGoModule` is not needed. The checks use `go build ./...` and `go test ./...` directly, which download modules via the Go module proxy (respecting `GOPROXY` and `GONOSUMCHECK`).

If example binaries are packaged in the future:

```nix
packages.default = pkgs.buildGoModule {
  pname = "cmdguard-typed-example";
  version = "0.0.0";
  src = ./.;
  vendorHash = "sha256-XXXXXXXXXXXX"; # Must update on dependency changes
  subPackages = [ "examples/typed" ];
};
```

The `vendorHash` would need updating whenever `go.sum` changes. Tools like `nix-update` or `gomod2nix` can automate this.

---

## 8. Risk Assessment

### 8.1 Technical Risks

| Risk                                               | Likelihood | Impact | Mitigation                                             |
| -------------------------------------------------- | ---------- | ------ | ------------------------------------------------------ |
| Go 1.26 not available in nixpkgs                   | Medium     | High   | Use `nixpkgs-unstable`; create Go overlay if needed    |
| golangci-lint v2 not in nixpkgs                    | Low        | High   | Override with `buildGoModule` or use nixpkgs-unstable  |
| `go build ./...` in Nix sandbox has network issues | Low        | Medium | Set `GOPROXY=off` and pre-fetch with `go mod download` |
| Build tags not working in Nix environment          | Low        | Low    | Pass `GOFLAGS=-tags=...` in check derivation           |
| Race detector not available in Nix Go package      | Very Low   | Medium | Verify with `go test -race` in devShell first          |
| `golangci-lint fmt` requires `$HOME`               | Medium     | Low    | Set `HOME=$(mktemp -d)` in check derivation            |

### 8.2 Adoption Risks

| Risk                                 | Likelihood | Impact | Mitigation                                  |
| ------------------------------------ | ---------- | ------ | ------------------------------------------- |
| Contributors unfamiliar with Nix     | High       | Low    | Keep non-Nix workflow; Nix is optional      |
| Nix daemon overhead on macOS         | Low        | Low    | Document `nix-daemon` setup                 |
| `flake.lock` conflicts in PRs        | Medium     | Low    | Document update process: `nix flake update` |
| Longer CI times (Nix build overhead) | Medium     | Medium | Use `magic-nix-cache-action` or Cachix      |

### 8.3 Go 1.26 + Nix Specific Concerns

1. **Go 1.26 release timing** — If Go 1.26 was released very recently, it may only be in `nixpkgs-unstable`. This is fine for a dev environment but may cause CI instability if the channel updates unexpectedly.

2. **Overlay fallback** — If Go 1.26 is not available, an overlay can build it:

   ```nix
   goPkg = pkgs.go_1_26 or (pkgs.callPackage (
     { buildGo123Module, fetchFromGitHub }:
     buildGo123Module rec {
       pname = "go";
       version = "1.26.0";
       src = fetchFromGitHub {
         owner = "golang";
         repo = "go";
         rev = "go${version}";
         hash = "sha256-XXXX";
       };
       # ... (full bootstrap build)
     }
   ) {});
   ```

3. **Experimental build tags** — The `.golangci.yml` uses experimental tags (`goexperiment.goroutineleakprofile`, `goexperiment.jsonv2`, `goexperiment.simd`). These may require specific Go build configurations that aren't the default in nixpkgs' Go package.

---

## 9. Open Questions

| #   | Question                                                            | Owner       | Decision Needed Before |
| --- | ------------------------------------------------------------------- | ----------- | ---------------------- |
| 1   | Is Go 1.26 available in nixpkgs-unstable today?                     | Implementor | Phase 1 start          |
| 2   | Is golangci-lint v2 (major version 2) available in nixpkgs?         | Implementor | Phase 2 start          |
| 3   | Should we use `flake-parts` or raw `flake-utils`?                   | Maintainer  | Phase 1 start          |
| 4   | Should we enable Cachix for CI caching?                             | Maintainer  | Phase 4 start          |
| 5   | Should `nix flake check` eventually replace the existing CI matrix? | Maintainer  | After Phase 4          |
| 6   | Should example binaries be packaged as Nix packages?                | Maintainer  | Phase 5 or later       |
| 7   | Should we add `gomod2nix` for Go dependency management in Nix?      | Maintainer  | Phase 5 or later       |
| 8   | Minimum Nix version required? (Flakes stable since Nix 2.4+)        | Implementor | Phase 1 start          |
| 9   | Should `.envrc` be committed or added to `.gitignore`?              | Maintainer  | Phase 1                |
| 10  | Should we use `nixfmt` or `alejandra` for Nix formatting?           | Implementor | Phase 3                |
| 11  | Do the experimental Go build tags affect Nix-built Go?              | Implementor | Phase 2                |

---

## 10. References

### Nix Documentation

- [Nix Flakes (nix.dev)](https://nix.dev/concepts/flakes)
- [buildGoModule (Nixpkgs manual)](https://nixos.org/manual/nixpkgs/unstable/#sec-language-go)
- [Nix Flakes Wiki](https://wiki.nixos.org/wiki/Flakes)
- [treefmt-nix](https://github.com/numtide/treefmt-nix)
- [flake-utils](https://github.com/numtide/flake-utils)
- [nixfmt (RFC 166)](https://github.com/NixOS/nixfmt)

### CI Integration

- [DeterminantSystems/nix-installer-action](https://github.com/DeterminantSystems/nix-installer-action)
- [DeterminantSystems/magic-nix-cache-action](https://github.com/DeterminantSystems/magic-nix-cache-action)
- [cachix/install-nix-action](https://github.com/cachix/install-nix-action)

### direnv

- [direnv documentation](https://direnv.net/)
- [nix-direnv](https://github.com/nix-community/nix-direnv)

### Similar Projects

- [hoopsnake/flake.nix](https://github.com/boinkor-net/hoopsnake) — Go project with flake-parts, devshell, golangci-lint
- [kiln/flake.nix](https://github.com/Thunderbottom/kiln) — Go project with buildGoModule, checks, devShells
- [scip/flake.nix](https://github.com/scip-code/scip) — Go project with checks, treefmt, multi-language

---

## Appendix A: Quick Start for Reviewers

```bash
# 1. Install Nix (if not already installed)
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install

# 2. Enable flakes (add to ~/.config/nix/nix.conf or /etc/nix/nix.conf)
experimental-features = nix-command flakes

# 3. Enter dev shell
nix develop

# 4. Run all checks
nix flake check -L

# 5. Format Nix files
nix fmt
```

## Appendix B: Go 1.26 Availability Check

Before starting implementation, verify Go 1.26 availability:

```bash
# Check if go_1_26 exists in nixpkgs-unstable
nix eval nixpkgs#go_1_26.meta.description

# If not available, check what Go versions exist
nix search nixpkgs go_1_

# Fallback: use nixpkgs-unstable
nix eval github:NixOS/nixpkgs/nixpkgs-unstable#go.meta.version
```

## Appendix C: Commands Summary

| Command                                 | What It Does                                 |
| --------------------------------------- | -------------------------------------------- |
| `nix develop`                           | Enter dev shell with Go, golangci-lint, just |
| `nix develop -c just verify`            | Run just verify inside Nix shell             |
| `nix flake check -L`                    | Run all checks (build, test, lint, format)   |
| `nix flake check -L .#test`             | Run only the test check                      |
| `nix fmt`                               | Format Nix files                             |
| `nix flake update`                      | Update all flake inputs (nixpkgs, etc.)      |
| `nix flake lock --update-input nixpkgs` | Update only nixpkgs                          |
| `direnv allow`                          | Enable automatic shell activation on cd      |
