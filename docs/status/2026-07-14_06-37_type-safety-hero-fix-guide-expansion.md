# Status Report — 2026-07-14 06:37

## Session Goal

"Make the cmdguard public repo superb" — follow-up improvement pass after the initial website + README + CI/CD launch. Focus: type safety, missing content, build verification, correctness bugs.

---

> **Update 2026-07-23:** The type-safety improvements and 6 guide pages shipped. DNS, the Firebase service-account secret, and the OG image remain external blockers. The workspace now has 4 optional sub-modules; `manpage` was removed in `34a0c6e`.

## A) FULLY DONE

### 1. Fixed Hero Code Contradiction (Critical)

The hero code sample on the landing page showed `panic(err)` — directly contradicting the "zero panics" USP that the entire marketing message is built around. Replaced with proper `fmt.Fprintf(os.Stderr, ...) + os.Exit(1)` pattern, matching the README Quick Start.

**Files:** `website/src/data/hero-code.ts`, `website/src/components/HeroSection.astro`

### 2. Type Model: SupportLevel Enum + File Split

Replaced the stringly-typed `ComparisonRow` (all fields were `string`) with a proper typed model:

- Added `SupportLevel` union type: `"full" | "partial" | "diy" | "none" | "native"`
- Added `FrameworkSupport` interface: `{ level: SupportLevel, note: string }`
- Added `FrameworkKey` and `frameworkKeys` const array
- `ComparisonRow` now has `FrameworkSupport` fields instead of raw strings
- `ComparisonSection.astro` renders colored status dots driven by `SupportLevel` (green=full, amber=partial, gray=diy/none, indigo=native)
- Split `sections.ts` (mixed comparisons + useCases) into `comparisons.ts` + `usecases.ts`

**Files:** `website/src/data/types.ts`, `website/src/data/comparisons.ts` (new), `website/src/data/usecases.ts` (new), `website/src/data/sections.ts` (deleted), `website/src/components/ComparisonSection.astro`, `website/src/components/UseCasesSection.astro`

### 3. Six Missing Documentation Guide Pages

Added 6 guide pages that were referenced in docs/status but never written. Website went from 16 to 22 pages.

| Guide                  | Content                                                                          |
| ---------------------- | -------------------------------------------------------------------------------- |
| `custom-types.mdx`     | 8 built-in value types, RegisterTypeHandler, COW registry, custom validators     |
| `audit-log.mdx`        | WithAuditLog wiring, programmatic access, 11-format export, consumer pattern     |
| `lifecycle.mdx`        | Signal handling, graceful shutdown, WithCleanup, Doctor command, Version command |
| `flow-context.mdx`     | BranchingFlowContext, path tracking, hierarchical value propagation              |
| `plugins.mdx`          | Plugin interface, 3 registration scopes, writing a plugin                        |
| `shell-completion.mdx` | Dynamic completion, WithCompletion, static arg validation                        |

**Files:** 6 new `.mdx` files in `website/src/content/docs/guides/`, updated sidebar in `astro.config.mjs`

### 4. CI/CD Import-Path Check Bug Fix (Critical)

The "Verify import paths contain /v3" step in `website.yml` was silently passing because:

- Job-level `defaults.run.working-directory: website` applied to the step
- The step grepped for `README.md` and `website/src/` — which from `website/` resolves to `website/README.md` (doesn't exist) and `website/website/src/` (doesn't exist)
- `grep -rl` found zero files, the loop body never executed, `ERRORS` stayed empty, step passed

Fixed by adding `working-directory: ${{ github.workspace }}` to that specific step.

**File:** `.github/workflows/website.yml`

### 5. .gitignore Blocking audit-log.mdx

The root `.gitignore` has `audit-log.*` (for BuildFlow artifacts). This silently blocked `website/src/content/docs/guides/audit-log.mdx` from being tracked. The file existed on disk but `git status` showed nothing. Added exception: `!website/src/content/docs/guides/audit-log.mdx`.

**File:** `.gitignore`

### 6. README Polish

- Added documentation links for all 6 new guides
- Fixed PostRunE guidance: was telling users to "put `defer` in RunE" — now points to `WithCleanup[T]` guide (the actual cmdguard solution)

**File:** `README.md`

### 7. Full Build Verification

- `astro check`: 0 errors, 0 warnings, 0 hints (29 files)
- `pnpm run build`: 22 pages built successfully
- `html-validate`: 0 errors
- Import path check: passes from repo root
- BuildFlow pre-commit: passed (commit intact, auto-fixes applied)

---

## B) PARTIALLY DONE

### 1. CI/CD Workflow Verification

**Fixed:** The import-path check bug.
**Not done:** The workflow has NEVER actually run on GitHub Actions. We verified locally but the real test is a push trigger. The deploy job will fail without `FIREBASE_SERVICE_ACCOUNT` secret.

### 2. Comparison Table Consistency

**Done:** Website `ComparisonSection` now uses typed `SupportLevel` with colored dots.
**Not done:** The README's trinity comparison table still uses plain text/dashes (`—`, `Yes`, `Some`). The two tables now diverge in both content and format. The README table has 9 rows; the website has 8 (missing "Health checks" row).

### 3. Guide Page Content Quality

**Done:** All 6 pages written with code examples, compile successfully.
**Not done:** Several pages have thin coverage or potential accuracy issues:

- `shell-completion.mdx` — completion code examples are inferred from CHANGELOG entries, not verified against actual API signatures
- `plugins.mdx` — the `Plugin` interface definition is a guess based on FEATURES.md description, not verified against source code
- `custom-types.mdx` — `RegisterValidator` example uses a regex pattern but the actual API signature wasn't verified
- No cross-references between guides (e.g., lifecycle.mdx should link to dependency-injection.mdx for DI services)

---

## C) NOT STARTED

1. **Website deployment** — updated site not deployed to Firebase; `cmdguard.web.app` still runs the old 16-page version
2. **DNS propagation** — `cmdguard.lars.software` still doesn't resolve; 12 README links + all OG metadata point there
3. **Firebase custom domain** — not connected in Firebase console
4. **GitHub secret** — `FIREBASE_SERVICE_ACCOUNT` not added
5. **OG image** — still using `favicon.svg` (32x32) as `og:image`; will look terrible on social media
6. **Lighthouse audit** — never run; don't know performance/accessibility/SEO scores
7. **Mobile audit** — never visually checked on mobile viewport
8. **Social preview** — no Twitter card image, no GitHub social preview
9. **Website logo** — still using the placeholder "CG" monogram SVG
10. **README GitHub topics/description** — mentioned as done in prior session but not verified this session

---

## D) TOTALLY FUCKED UP

### 1. The `panic(err)` Was In The Code From Day One

The hero code sample — the FIRST thing every visitor sees — contained `panic(err)` while the entire value proposition is "zero panics by construction." This was written in the initial website launch commit and survived through multiple "polish" and "USP pivot" passes without anyone catching it. This is a credibility-destroying contradiction for a library whose identity is built on never panicking.

### 2. CI/CD Import-Path Check Was Security Theater

The "Verify import paths contain /v3" step was supposed to catch broken import paths in documentation. It was silently passing because the working directory was wrong — it found ZERO files every single run. If someone had written `github.com/larsartmann/cmdguard"` (missing `/v3`) in any doc, the check would have happily passed. This is the worst kind of CI check: one that looks like it's protecting you but isn't.

### 3. .gitignore Silently Eating Documentation

`audit-log.*` in `.gitignore` silently blocked `audit-log.mdx`. The file existed on disk, built correctly, but would never have been committed. If I hadn't run `git status` with `--untracked-files=all` and noticed only 6 of 7 new files, the audit-log guide would have been silently excluded from the commit. The CI build would have 21 pages instead of 22, and the sidebar link would 404.

### 4. Never Visually Verified ANYTHING

I built the website, ran typecheck, ran HTML validation — all programmatic checks. But I never once looked at the actual rendered pages. I don't know if:

- The comparison table colored dots look good or like a mess
- The hero code formatting is correct after the panic(err) fix
- The new guides render with proper code highlighting
- The mobile nav still works
- Anything looks "superb" at all

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Code Quality

1. **README comparison table should use the same typed data model** — currently two sources of truth that diverge
2. **Extract comparison data into a shared JSON/TS file** that both README table generation and website consume
3. **The `supportLevels` array includes `"native"` but it's used exactly once** — consider whether it's a meaningful distinction from `"full"` or if it muddies the type
4. **Guide pages need internal cross-references** — lifecycle.mdx → dependency-injection.mdx, plugins.mdx → custom-types.mdx, etc.
5. **`api-reference.mdx` should link to all guides** — it's currently a thin page that doesn't reference the new content
6. **The landing page feature grid should link to relevant guides** — currently no path from landing page to the new 6 guides
7. **No "Getting Started" tutorial** — installation + quick-start exist but there's no end-to-end walkthrough

### CI/CD

8. **Deploy step uses `pnpm add -g firebase-tools`** — should use a versioned GitHub Action (`w9jds/firebase-action`) for reproducibility
9. **`pnpm audit` runs with `continue-on-error: true`** — security vulnerabilities are silently ignored
10. **No Lighthouse CI** — should add performance/accessibility/SEO regression detection
11. **No link checker** — should verify all internal links resolve (especially after adding 6 new pages)
12. **No preview deployments** — PRs should deploy previews, not just build-check

### Content

13. **Verify guide code examples compile** — several were written from inference, not verified against the actual API
14. **Add a "Full Tutorial" guide** — docs/TUTORIAL.md exists but isn't on the website
15. **Port docs/COMPARISON.md to a website page** — detailed framework comparison already exists in repo
16. **Port docs/PERFORMANCE.md to a website page** — benchmark data is consumer-relevant
17. **Add "Examples" page** — the taskctl example is the best marketing material and isn't showcased on the website
18. **Changelog page is likely stale** — should be auto-generated or synced from CHANGELOG.md

### Design & UX

19. **OG image generation** — use `astro-og-canvas` or similar for proper social share images
20. **Favicon/logo design** — current "CG" monogram is a placeholder
21. **Mobile audit** — never tested; could be broken
22. **Accessibility audit** — skip link exists but keyboard nav, ARIA, color contrast never verified
23. **Dark/light mode consistency** — Starlight docs pages have their own theme toggle separate from the landing page toggle; could be confusing
24. **Page load performance** — no idea if fonts, JS, or CSS are optimized

### Operations

25. **Firebase deploy pipeline** — manually deploying works but CI auto-deploy doesn't (missing secret)
26. **DNS still not applied** — the most critical blocker; nothing at `cmdguard.lars.software` works
27. **No monitoring/alerting** — if the site goes down, nobody knows
28. **No error tracking** — client-side JS errors (header toggle, theme toggle, copy button) are silent

---

## F) UP TO 50 THINGS TO GET DONE NEXT

### Immediate (blocks everything)

1. ~~Fix hero code panic~~ ✅ DONE
2. Push commit to origin/master
3. Deploy updated website to Firebase (`firebase deploy --only hosting:cmdguard`)
4. Apply DNS terraform (`terraform apply` in domains repo, from whitelisted IP)
5. Add `FIREBASE_SERVICE_ACCOUNT` GitHub secret
6. Connect `cmdguard.lars.software` as Firebase custom domain
7. Verify CI/CD workflow passes on GitHub Actions

### Content (high impact)

8. Verify all 6 new guide code examples against actual cmdguard API source
9. Fix `plugins.mdx` Plugin interface definition — read `plugin.go` and match exactly
10. Fix `shell-completion.mdx` WithCompletion signature — read `completion.go`
11. Fix `custom-types.mdx` RegisterValidator signature — verify against source
12. Add cross-references between all 14 guide pages
13. Sync README trinity comparison table with website `comparisons.ts` data
14. Add "Health checks" row to website comparison table (exists in README, missing from website)
15. Port docs/TUTORIAL.md to a website "Full Tutorial" guide
16. Port docs/COMPARISON.md to a website "Framework Comparison" page
17. Add an "Examples" page showcasing the taskctl CLI
18. Update `api-reference.mdx` with links to all guides
19. Add a `getting-started/full-tutorial.mdx` end-to-end walkthrough
20. Verify changelog.mdx is current with CHANGELOG.md

### CI/CD & Quality

21. Replace `pnpm add -g firebase-tools` with `w9jds/firebase-action@v15`
22. Remove `continue-on-error: true` from `pnpm audit` or document why it stays
23. Add a link-checker step to CI (lychee or similar)
24. Add Lighthouse CI for performance regression detection
25. Add Firebase preview channel deploys on PRs
26. Add a `website-test` job that runs `astro build` on PRs (already partially done)
27. Pin pnpm dependencies to exact versions (currently using `^` ranges)

### Design & UX

28. Generate proper OG image (1200x630) with the trinity messaging
29. Design a proper logo/favicon (not "CG" monogram)
30. Run Lighthouse audit and fix all issues
31. Test on mobile viewport (Chrome DevTools device emulation at minimum)
32. Test keyboard navigation through landing page
33. Verify color contrast meets WCAG AA
34. Audit Starlight docs theme vs landing page theme consistency
35. Add scroll-to-top button on docs pages
36. Add reading time estimates on guide pages

### Type Model & Architecture

37. Extract comparison data into shared format (README + website)
38. Consider generating README comparison table from `comparisons.ts` via script
39. Add `Schema.org` structured data for the new guide pages (currently only landing page has JSON-LD)
40. Add `lastmod` to sitemap entries
41. Add `breadcrumb` structured data to guide pages

### Operations

42. Set up uptime monitoring for `cmdguard.web.app`
43. Set up uptime monitoring for `cmdguard.lars.software` (once DNS works)
44. Add ` SECURITY.md` file
45. Add GitHub issue templates
46. Add GitHub discussion templates
47. Set up Dependabot for website pnpm dependencies
48. Add `CITATION.cff` for academic citations
49. Create a GitHub release for v3.0.0 (if not already done)
50. Add the website to the GitHub repo "About" section with description

---

## G) TOP 2 QUESTIONS

### 1. Should the README comparison table be generated from the website's `comparisons.ts`?

Right now there are two comparison tables (README + website) with different data (9 rows vs 8 rows) and different formats (text dashes vs typed SupportLevel dots). They WILL diverge further over time. Options:

- **A) Generate README from TS** — Write a script that reads `comparisons.ts` and outputs a markdown table. Single source of truth. But README becomes harder to edit directly.
- **B) Keep separate, accept divergence** — README is for GitHub scrollers, website is for visitors. Different audiences, different formats.
- **C) Remove the comparison from README entirely** — Link to the website instead. Less maintenance, but README becomes less compelling for GitHub discovery.

I cannot decide this because it depends on whether you value README self-containedness or single-source-of-truth more.

### 2. Should I deploy the updated website now, or wait until DNS + custom domain is ready?

The website is live at `cmdguard.web.app` but running the old 16-page version. The new version (22 pages, fixed hero code, typed comparison table) is committed but not deployed. Options:

- **A) Deploy now to `.web.app`** — The fixed hero code (no `panic`) goes live immediately. But all OG metadata, canonical URLs, and sitemap entries point to `cmdguard.lars.software` which doesn't resolve yet — SEO/crawlers will see broken canonical links.
- **B) Wait for DNS** — Everything goes live at once with correct URLs. But the `panic(err)` stays visible on the live site until DNS is applied, which could be days.

I cannot decide this because I don't know how long DNS will take or whether you care more about the hero fix being live vs SEO correctness during the transition.

## Resolution (2026-07-23)

- §A items 1–5 are complete.
- External blockers (DNS/secret/OG) remain unchanged.
- `manpage` was removed in `34a0c6e`; current sub-modules are glamour, prompts, spinner, telemetry.
