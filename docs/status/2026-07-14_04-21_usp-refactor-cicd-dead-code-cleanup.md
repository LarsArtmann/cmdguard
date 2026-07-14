# Status Report — 2026-07-14 04:21 — USP Refactor, CI/CD, Dead Code Cleanup

> Session covered: website creation, USP pivot to "trinity", CI/CD workflow, type model cleanup, README polish.

---

## A) FULLY DONE

### 1. Website — Full Astro + Starlight + Firebase Deployment

- **16-page Starlight documentation site** (installation, quick-start, 8 guides, API reference, changelog, contributing, related tools)
- **Custom landing page** with hero (animated code preview, GitHub stars badge), feature grid (8 cards), comparison table, use cases, CTAs
- **Indigo/violet theme** distinct from sibling projects (cyan gogenfilter, emerald go-atomic-write)
- **SEO**: OG/Twitter meta, JSON-LD SoftwareApplication schema, sitemap, robots.txt, PWA manifest
- **Security**: CSP headers (HSTS, X-Frame-Options DENY, CORP, COOP, Permissions-Policy), cleanUrls, `/docs/*` redirect
- **Firebase hosting site `cmdguard`** created in `lars-software` project, deployed (HTTP 200 live)
- **Nix flake** with `dev`/`build`/`preview`/`deploy` apps + treefmt
- **Build pipeline**: `astro build && node scripts/fix-csp.mjs` patches inline script hashes

### 2. USP Pivot — "The Trinity"

- Rewrote README pitch from "type-safe flags" (table stakes, Kong does this) to **"From flags to shutdown, type-safe"** — the only framework unifying type-safe flags + DI lifecycle + zero-panic contracts
- Rewrote website hero code to showcase DI (`v3.Provide`, `v3.Invoke`, `WithGracefulShutdown`)
- Expanded comparison from 2-column (Cobra vs cmdguard) to **5-column matrix** (Cobra/Kong/urfave/cli/cmdguard)
- Updated GitHub description, all website metadata (config.ts, astro.config.mjs, OG descriptions)
- Reordered feature grid: DI first, Zero Panics second, Validated at Construction third

### 3. CI/CD Pipeline

- Created `.github/workflows/website.yml` — triggers on `website/**` and `website.yml` changes
- Build job: npm install, typecheck (`astro check`), build, HTML validation, import-path verification
- Deploy job: runs on push to master only, deploys to Firebase `cmdguard` hosting target
- Added website CI badge to README

### 4. Code Quality — Dead Type Removal

- Removed unused `FrameworkColumn` type (exported, never imported)
- Removed `cmdguardBest: boolean` field from `ComparisonRow` (always `true`, never read by component)
- Removed redundant Cobra footguns table from README (trinity matrix above it covers the same 4 capabilities)

### 5. GitHub Metadata

- Homepage URL set to `https://cmdguard.lars.software`
- 20 topics (added: `cobra-wrapper`, `zero-panic`, `samber-do`, `go-cli`, `functional-options`, `exit-codes`)
- Description updated to trinity messaging

### 6. DNS Record

- `cmdguard` CNAME → `cmdguard.web.app.` added to `lars.software.tf` in domains repo

---

## B) PARTIALLY DONE

### 1. DNS Propagation — BLOCKED

- CNAME added to Terraform but NOT applied (Namecheap API IP not whitelisted in this environment)
- `cmdguard.lars.software` does NOT resolve yet
- `cmdguard.web.app` IS live and serving content

### 2. CI/CD Secret

- Workflow file exists but **will fail on deploy** without `FIREBASE_SERVICE_ACCOUNT` secret
- Build job will work without the secret (no deploy on PRs)
- User must add the secret to GitHub repo settings

### 3. Website CI/CD — First Run Not Verified

- Workflow has never run (created after the last push)
- Need to confirm the build job passes on GitHub Actions (npm install, astro check, html-validate)
- Deploy job blocked by missing secret

### 4. README Links Point to Unresolvable Domain

- All website links in README point to `cmdguard.lars.software` which doesn't resolve yet
- They'll work once DNS is applied, but until then they 404

---

## C) NOT STARTED

### 1. Firebase Custom Domain Connection

- Firebase console needs `cmdguard.lars.soference` added as a custom domain on the `cmdguard` hosting site
- This generates the `_acme-challenge.cmdguard` TXT record value for SSL provisioning
- Not started — blocked by DNS propagation

### 2. Lighthouse CI

- gogenfilter has a `lighthouse.yml` workflow + `lighthouserc.json` config
- cmdguard has neither — no performance/accessibility/SEO regression checks

### 3. OG Image Generation

- gogenfilter uses `astro-og-canvas` for dynamic social share images
- cmdguard uses `favicon.svg` as a fallback — no proper OG image

### 4. GitHub Social Preview Image

- No social preview uploaded to GitHub repo settings

### 5. `package-lock.json` Decision

- Currently tracked in git (needed for CI npm cache)
- No explicit decision made — it just got committed alongside everything else

### 6. Website `.gitignore` for `dist/` Verification

- `dist/` is gitignored and confirmed not tracked

### 7. Update AGENTS.md with Website Section

- AGENTS.md project structure doesn't mention the `website/` directory
- No documentation of build/deploy commands for the website

### 8. Stale Status Reports Cleanup

- `docs/status/` now has 3 status reports from this session
- The oldest one (02:14 jsonv2 migration) may be stale

---

## D) TOTALLY FUCKED UP

### 1. `animations.js` Copy-Paste Error (FIXED)

- Initially wrote `header.js` content into `animations.js` — nav toggle code instead of IntersectionObserver
- Caused the entire page below the hero to be pure black (opacity: 0 never removed)
- Fixed in the same session, redeployed

### 2. `hero-code.ts` Double Export (FIXED)

- Initial version had `export const heroCode` at top AND `export { heroCode }` at bottom
- Caused build failure: "Duplicated export 'heroCode'"
- Fixed by removing the redundant export

### 3. Changelog.mdx Referenced Non-Existent Component (FIXED)

- Initially wrote `import { Changelog } from '../../components/Changelog'` — component never existed
- Fixed by removing the import

### 4. No Staging/Preview Before Deploy

- Deployed directly to production Firebase without `astro preview` first
- Worked, but violates best practice

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Type Model

1. **Website data model is anemic** — `ComparisonRow` could use a `Support` enum (`full` | `partial` | `none`) instead of free-text strings, enabling icon/styling per cell
2. **`sections.ts` mixes concerns** — comparison data and use-case data in one file; should be `comparisons.ts` + `usecases.ts`
3. **No shared `Framework` type** — framework names ("Cobra", "Kong", etc.) are repeated as column headers; a `Framework` type with display name + URL would enable linking to each framework
4. **Hero code is hardcoded HTML string** — could use a syntax highlighter library (Shiki, Starlight already bundles it) instead of manual `<span>` wrapping

### Process

5. **Always `astro preview` before `firebase deploy`** — should be wired into the flake deploy app
6. **Never copy-paste JS files** — the animations.js disaster was preventable by reading the source before writing
7. **Lighthouse baseline** — should establish performance scores before iterating on design
8. **Firebase preview channel** — deploy to a staging channel before production

---

## F) Up to 50 Things to Get Done Next

### Immediate Blockers (User Actions)

1. **Apply Terraform DNS** — `cd domains && terraform apply` from whitelisted IP
2. **Add `FIREBASE_SERVICE_ACCOUNT` secret** — GCP service account JSON with Firebase Hosting admin on `lars-software`
3. **Connect custom domain in Firebase** — add `cmdguard.lars.software` to `cmdguard` hosting site
4. **Verify SSL** — confirm `https://cmdguard.lars.software` works after DNS + custom domain
5. **Verify website CI passes** — push to trigger `website.yml` build job

### Type Model & Architecture (This Repo)

6. **Add `Support` enum to `ComparisonRow`** — `full | partial | none | diy` instead of free-text
7. **Split `sections.ts`** into `comparisons.ts` + `usecases.ts`
8. **Add `Framework` type** with display name + URL for linking from comparison table
9. **Use Shiki for hero code highlighting** — Starlight already bundles it, no manual `<span>` needed
10. **Add `LandingPageData` aggregate type** — bundle config + features + comparisons + usecases into one import

### Website Content Depth

11. **Add BranchingFlowContext guide** — track command path + share values across hierarchy
12. **Add Audit Log guide** — `WithAuditLog`, `ExportAuditLog`, 11 formats
13. **Add Plugin System guide** — `Plugin` interface, `RegisterPlugin`, `WithPlugin`
14. **Add Doctor/Version commands guide** — `DoctorCommand`, health checks, `VersionCommand`
15. **Add Shell completion guide** — `WithCompletion`, dynamic completion
16. **Add Signal handling guide** — `WithSignalHandling` vs `WithGracefulShutdown`
17. **Add WithCleanup guide** — cleanup hooks that fire even when RunE errors
18. **Port docs/TUTORIAL.md** to website as multi-page tutorial
19. **Port docs/COMPARISON.md** to website (Kong, urfave/cli, go-flags, sflags)
20. **Port docs/PERFORMANCE.md** to website benchmarks page
21. **Port docs/CLI_DESIGN_PRINCIPLES.md** to website
22. **Port docs/COBRA_FOOTGUNS.md** to website
23. **Full changelog on website** — currently only v3.0.0 highlights
24. **Port docs/MIGRATION_v2_v3.md** to website

### Website Polish

25. **Design proper logo/favicon** — current SVG is a placeholder
26. **Add OG image generation** — `astro-og-canvas` integration
27. **Upload GitHub social preview** image
28. **Run Lighthouse audit** — verify performance/accessibility/SEO scores
29. **Add Lighthouse CI workflow** — `lighthouse.yml` + `lighthouserc.json`
30. **Mobile responsive audit** — verify on phone-width viewports
31. **Verify Pagefind search** — confirm Starlight search works on deployed site
32. **Verify 404 page styling** — built but never checked
33. **Add `dedup.sh` script** — gogenfilter has jscpd code duplication check for website
34. **Add `validate-docs` app** — gogenfilter uses `md-go-validator` to validate code blocks in docs
35. **Add website README** — document build/deploy process inside `website/`

### CI/CD & Operations

36. **Verify website CI first run** — confirm build job passes on GitHub Actions
37. **Add Firebase preview channels** — staging before production
38. **Add `website/` to root AGENTS.md** — document build commands, deploy process, flake apps
39. **Update root `flake.nix`** — consider convenience app for website deploy
40. **Add sub-module docs pages** — dedicated pages for glamour, manpage, prompts, spinner, telemetry
41. **Website CI: stale reference check** — gogenfilter checks for references to deleted files
42. **Website CI: CHANGELOG sync** — gogenfilter verifies changelog.mdx matches CHANGELOG.md

### Marketing & Discoverability

43. **Write launch announcement** — content for r/golang, HN
44. **Submit to Awesome Go** lists
45. **Add `.github/FUNDING.yml`** — GitHub Sponsors
46. **Add animated demo** to README — asciinema or GIF
47. **Add "Who uses cmdguard?"** section once there are users
48. **Add sub-module version badges** — show each is independently versioned

### Code & Repo Hygiene

49. **Clean up stale status reports** — `docs/status/` has 3 from today
50. **Remove `docs/status/2026-07-14_02-14_jsonv2-migration-embraced.md`** if stale — it's pre-website
51. **Update FEATURES.md** — add website as a feature entry

---

## G) Top 2 Questions I Cannot Answer Myself

### 1. Should the website comparison table link to each framework?

The 5-column comparison table mentions Cobra, Kong, and urfave/cli by name. Linking to their GitHub repos would be good etiquette and improve credibility. But it also means maintaining external links. **Should I add hyperlinks to each framework in the comparison table headers?**

### 2. Should we add `astro-og-canvas` now or defer?

OG images significantly improve social media sharing appearance. gogenfilter uses `astro-og-canvas` for dynamic generation. But adding it now means another dependency, another build step, and design decisions about the image layout. **Should I add OG image generation now, or defer until we have a proper logo?**
