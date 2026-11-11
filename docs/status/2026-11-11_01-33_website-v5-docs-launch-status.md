# Website v5 docs rewrite and launch: status report

- **Date:** 2026-11-11 01:33 CET
- **Scope:** this session only: cmdguard website docs port (v4 to v5), docs truth gates, Firebase deploy, custom domain attach.
- **Verdict:** the website is live at https://cmdguard.lars.software, teaches v5 everywhere, and builds clean. It is **not yet a complete launch**: no CI, no demo video, no versioned v4 docs, and no visual QA.

## Evidence collected

| Check                                      | Result                                                                        |
| ------------------------------------------ | ----------------------------------------------------------------------------- |
| `pnpm build`                               | 28 pages, sitemap + pagefind index, CSP patch 28/28                           |
| `pnpm exec astro check`                    | 0 errors, 0 warnings, 0 hints (30 files)                                      |
| `pnpm exec html-validate "dist/**/*.html"` | clean                                                                         |
| `python3 scripts/check-website-links.py`   | 28 pages, 0 broken internal targets                                           |
| `bash scripts/check-guide-sync.sh`         | 0 stale majors, 0 unknown symbols, all full doc programs compile              |
| `firebase deploy --only hosting:cmdguard`  | release complete                                                              |
| `https://cmdguard.lars.software`           | 200 on every sampled route, cert `CN=cmdguard.lars.software`                  |
| Domain record state                        | `HOST_ACTIVE`, `OWNERSHIP_ACTIVE`, `cert: CERT_PROPAGATING`, type `TEMPORARY` |

Before this session: 14 of 16 doc pages used `v4.` (120 occurrences), the visible hero snippet was hand-written v4 while its Copy button served v5, every internal link and sitemap URL 301'd, the domain returned 404 "Site Not Found" with a `firebaseapp.com` certificate, and the docs guard passed green over all of it.

## a) Fully done

1. **Foundation fixes**
   - `website/src/components/HeroSection.astro`: server-side Go highlighter; the rendered snippet comes from `heroCode` (`src/data/hero-code.ts`), so visible text and Copy cannot drift. Star count now singular/plural correct.
   - `website/src/data/config.ts`: pkg.go.dev link moved to v5, keys sorted, em dash removed.
   - `astro.config.mjs`: v5 pkg.go.dev sidebar link, new navigation groups (Core Concepts, Guides, Sub-Modules & Tooling, Migration), `lastUpdated: true`, `editLink`, `trailingSlash: "always"`.
   - `website/firebase.json`: `trailingSlash: true`, which removes the 301 on every internal navigation and every sitemap URL.
   - `package.json`: TypeScript pinned to `^6.0.3` so `astro check` is usable again (TS 7 is unsupported by Astro).
2. **Guides ported to v5 (14 pages)**: type-safe-flags, custom-types, config-files, dependency-injection, error-handling, lifecycle, middleware, plugins, rich-output, shell-completion, sub-modules, flow-context, audit-log, migrating-from-cobra.
3. **New pages (6)**: production-batteries, machine-contract (describe JSON), doctor, schema, lint (CG001-CG008 + baseline), migrating-v4-to-v5.
4. **Rewritten pages**: api-reference (real constructor and method signatures, all option groups, sealed `AddableCommand` not misrepresented), changelog (v5 entry added, links moved from `main` to `master`), contributing (Go 1.27, seven sub-modules, `v5` test package), related-tools (real `glamour.WithHelp()`).
5. **Landing data**: features and comparison rows now describe v5 (batteries, machine contract, hints/codes, 16 exact formats); "all three" copy replaced; `di` card wording corrected to dependency-order shutdown.
6. **Docs truth gate**: `scripts/check-guide-sync.sh` rewritten on top of `scripts/check-docs-api.py`; it scans all `website/src` (mdx, astro, ts, mjs), audits every `v5.X` / `cli.X` / sub-module symbol against source, and compiles every complete Go program in the docs. Wired into `nix run .#check-all` with Python and Go on PATH.
7. **Link gate**: `scripts/check-website-links.py` verifies internal targets and same-page anchors across `dist`.
8. **Godoc corrections found on the way**: `glamour.WithHelp[T]()` to `WithHelp()`, `WithCLIVersion[Config]` to `WithCLIVersion`, the non-existent `schema.ValidateFile` example, ADR-001 option signatures.
9. **Memory corrections**: `WithBatteries()` was documented as a no-op in `AGENTS.md` and the cmdguard skill; source shows it re-enables the full set. Fixed in both plus the two site pages that repeated it.
10. **Deploy and domain**: deployed to `cmdguard` site (`lars-software` project), attached `cmdguard.lars.software` via the Firebase REST API, added the ACME TXT through the domains Terraform (and merged its duplicate `required_providers` blocks, which blocked `terraform init`), verified valid TLS and 200s on all sampled routes.

## b) Partially done

1. **Certificate**: the domain serves a valid cert for `cmdguard.lars.software`, but the API still reports `CERT_PROPAGATING` with type `TEMPORARY`. The managed/project-grouped cert swap is not confirmed.
2. **Visual QA**: no headless browser exists on this host. Verified HTTP 200, CSS token presence (`color-accent`), rendered HTML structure, and html-validate. No screenshot, no mobile layout check, no dark/light toggle check.
3. **Search**: Pagefind indexed 28 pages, but the search UI was never exercised in a browser.
4. **Docs gate coverage**: only complete `package main` programs are compiled. Fragments (the majority of samples) are symbol-audited, not type-checked. The broken `v5.Add(b, listCmd)` example in api-reference is proof: it was a fragment and slipped through until manual review.
5. **Live verification**: routes were checked with `--resolve` because the local resolver was refusing connections; general DNS health on the workstation was not fixed.
6. **Commit history**: everything landed through the auto-commit daemon's heuristic messages. The changes are correct but not reviewable as intentional commits.
7. **Domains repo**: the ACME TXT and the provider-block merge were committed by its own auto-commit daemon; the provider merge is a drive-by fix with no written note in that repo.
8. **README and site cross-links**: the repo README already points at v5 and the live domain, but there is no demo video link, no "watch the tour" badge, and no link from the site back to the repo's docs artifacts.
9. **`pnpm typecheck`**: `astro check` works now, but it is still not wired into a script or CI job.

## c) Not started

1. CI/CD for the website (GitHub Actions: frozen-lockfile install, build, checks, deploy).
2. Demo video (HyperFrames) as the launch centerpiece, which the website-launch skill treats as the default.
3. Versioned v4 documentation (the v4 pages are gone; only the migration guide covers the old API).
4. Per-page OG images and social cards (only `public/og/home.png` exists).
5. External link checking and link monitoring.
6. `llms.txt` / AI-discovery file and a machine-readable site map for agents.
7. Generated error-reference page on the site (currently links to GitHub).
8. Google Search Console / sitemap submission.
9. Accessibility audit (axe or pa11y) and Lighthouse budget.
10. Smoke-test script that hits the live domain after deploy.
11. Rollback plan and release pinning for hosting (`firebase hosting:rollback` is unverified).
12. Site version banner or "you are reading the v5 docs" affordance.
13. Deploy documentation in `website/README.md` (none exists).
14. Anchor coverage for headings in the link checker (only same-page `#` is checked today).
15. Verification that the temporary cert has been replaced after propagation finishes.
16. Verification that `firebase.json` redirect `/docs/:path*` still matches the new URL scheme.
17. Checking that pagefind excludes the 404 page and includes the new guides.
18. Testing `prefers-color-scheme` and the theme toggle in a real browser.
19. Re-running `terraform` plan for the whole domains repo after the provider merge (targeted apply only).
20. Removing the leftover `/tmp/cgdoccheck` and `/tmp` node scripts used for ad-hoc verification.

## d) Totally fucked up

1. **The first version of the api-reference Builder example was wrong**: `v5.Add(b, listCmd)` does not exist; `Add` constructs the command inline (`Add(b, use, flags, runE, opts...)`). The docs gate could not catch it because it only compiles full programs. Fixed, but only through manual review.
2. **The first docs-gate implementation was buggy**: `go build ./sampleN` collided with the sample directory name and top-level `const` declarations were not indexed, producing a burst of false findings. Fixed within the session.
3. **A subagent's temp build directory was committed to git history** (`tmp_docs_check` in `20636b9`, removed in `3bde98b`). The working tree is clean, but the history carries two junk commits.
4. **The site was deployed before I had a rollback check**: no previous release ID was recorded and no `hosting:rollback` drill was performed.
5. **I changed another repo without a note**: `~/projects/domains` got a provider-block merge and an ACME TXT. Correct and necessary, but undisclosed at the time and auto-committed.

## e) What we should improve

1. **Make the gate stronger than the reviewer.** A docs gate that only compiles full programs gives false confidence. Either mark every snippet as compilable and generate runnable wrappers, or add a type-check pass over snippets with a small harness.
2. **One source of truth per fact.** Error reference, API tables, and flag tables should be generated from the same data the library uses, not hand-mirrored. The `describe` contract and `ConfigKeys` already make this possible.
3. **Publish only through a gate.** A `scripts/deploy-website.sh` with preflight, build, all checks, deploy, and live smoke test would have made this session reproducible and would catch the next regression.
4. **No agent writes temp dirs inside the repo.** The git-history pollution was avoidable; subagents should use `/tmp`.
5. **Never claim a version is "published" until the custom domain serves 200 with a valid cert.** This session had the site deployed and the domain broken at the same time.
6. **Treat docs as code with tests.** The v4-to-v5 rename would have been caught months earlier by the symbol audit now in place.
7. **Keep a docs changelog.** The site changelog jumped from v4.0.0 to v5.0.0 with nothing in between; readers cannot see when `describe`, hints, or schema appeared.
8. **Verify cross-repo Terraform changes in the same session that makes them.** The domains repo is currently only targeted-applied.
9. **Distinguish "temporary cert" from "done"** in every launch report.
10. **Add a browser to the flake** (chromium headless) so visual QA is executable instead of a caveat.

## f) Up to 50 things we should get done next

**Deploy and operations**

1. Write `website/scripts/deploy.sh`: frozen lockfile install, build, all gates, deploy, live smoke test.
2. Add a GitHub Actions workflow for PR checks (build, astro check, html-validate, link check, docs gate).
3. Add a deploy job on `master` using a Firebase service account secret.
4. Record the current hosting release ID and drill `firebase hosting:rollback`.
5. Smoke-test script: assert 200 for `/`, all sidebar routes, `/sitemap-index.xml`, `/robots.txt`, `/404.html`.
6. Confirm `cert` reaches `CERT_ACTIVE`; then remove the temporary-cert caveat from reports.
7. Run a full `terraform plan` on `~/projects/domains` and confirm the provider merge changed nothing else.
8. Add an uptime check (or a scheduled smoke test) for `cmdguard.lars.software`.
9. Cache-bust verification: assert hashed `_astro/` assets are immutable and HTML is not cached.
10. Verify the legacy `/docs/:path*` redirect against a few old URLs.

**Docs correctness**

11. Add a snippet type-check harness so fragments are compiled too.
12. Convert the most important fragments into full programs marked for compilation.
13. Generate the error-reference page on the site from `sentinelFacts()`.
14. Generate the flag-tag table and the option tables from source.
15. Add a docs test that every internal link target has a title and description.
16. Add a test that every guide has a "Where to go next" section.
17. Verify the claim "62 sentinels" against `GenerateErrorReference()` in CI.
18. Add a "docs build" check that fails on duplicate headings or empty pages.
19. Audit every page for stale counts ("16 formats", "seven modules") against source constants.
20. Add a v4 docs snapshot under `/v4/` or a separate site, with a banner.

**Content and UX**

21. Add "Who is this for?" and "When NOT to use cmdguard" to the landing page or quick start.
22. Produce the HyperFrames demo video and embed it above the fold.
23. Add a "Show me the machine contract" live example with real `describe` output.
24. Add a cookbook page: config + env + flags precedence in one runnable program.
25. Add a "production checklist" page (version, doctor, signals, hints, lint).
26. Add an FAQ (module path, v4 vs v5, Cobra compatibility, DI scope lifetime).
27. Add a comparison page with dated, sourced claims instead of only the table.
28. Add per-page OG images (title + one-line summary).
29. Add a site version badge and a link to the migration guide in the header for v4 visitors.
30. Add a "docs feedback" link per page (GitHub issue prefilled).
31. Add a changelog feed (RSS or Atom) for releases.
32. Add copy-to-clipboard verification for every code block (the hero now works).

**SEO and discovery**

33. Submit the sitemap to Google Search Console and Bing.
34. Add `llms.txt` with links and a one-paragraph product summary.
35. Validate JSON-LD (SoftwareApplication) with a structured-data test.
36. Add canonical checks to the link gate (every page's canonical must be the live domain).
37. Add `hreflang`/locale only if a second language is planned; otherwise document English-only.
38. Add a `humans.txt` or security.txt only if there is an actual reporting address.
39. Check that 404 returns the custom page and a 404 status on the custom domain.
40. Add external-link checking on a schedule (not per PR).

**Quality gates**

41. Wire `astro check` and html-validate into `nix run .#check-website`.
42. Add `pnpm install --frozen-lockfile` to CI to prevent lockfile drift.
43. Add a markdown/MDX linter (frontmatter required, heading order, code fence language).
44. Add an accessibility gate (axe or pa11y) for the landing page and one guide.
45. Add a Lighthouse budget for LCP/CLS on the landing page.
46. Add a test for the hero highlighter (keywords, strings, comments, HTML escaping).
47. Add unit tests for `check-docs-api.py` and `check-website-links.py`.
48. Add a visual regression snapshot once a headless browser is in the flake.
49. Add a "docs-only" path filter so PRs that touch only `website/**` run the website jobs.
50. Review and squash the auto-commit history around the website work before any public release.

## g) Three questions I cannot answer myself

1. **Merge and release intent:** should this work stay on `v5-scaffold`, or should it be merged to `master` and tied to a release (v5.0.1 / v5.1.0) with a tag and changelog entry?
2. **CI and credentials:** do you want a GitHub Actions deploy pipeline for the website, which requires creating a Firebase service-account key and storing it as a repository secret?
3. **Launch assets:** do you want the HyperFrames demo video produced now as the landing-page centerpiece, and if so, should it be narrated or a silent motion piece, and roughly how long?

## Session-level self-assessment

**What I forgot in the moment**

- I treated "deployed" as "published" while the domain was still broken.
- I did not record the hosting release before deploying, so there was no rollback reference.
- I did not check the v4 docs' fate (they are simply gone).
- I did not verify page search, mobile layout, or the theme toggle.
- I did not wire any of the new checks into CI.
- I did not squash or annotate the auto-commit history.

**What I would do differently**

- Run the domain attach and DNS steps first, then deploy, then verify from the public URL.
- Build the snippet type-check harness before rewriting 20 pages, not after.
- Keep all temporary work outside the repo so the auto-commit daemon cannot capture it.
- Produce a single deploy script and use it, so the next person does not reconstruct the sequence.

**Biggest residual risk**

The docs are accurate today because a gate now compares them to source. That gate only knows about complete programs and exported identifiers; semantic claims (counts, default behaviour, ordering) are still hand-checked. The next drift will be a sentence, not a symbol.

**Confidence**

- Deployed and reachable: high.
- Docs API accuracy: high for symbols and full programs, medium for prose claims.
- Launch completeness: medium; the product story is good, the launch machinery (video, CI, versioned docs) is missing.
