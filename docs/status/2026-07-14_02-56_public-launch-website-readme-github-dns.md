# Status Report — 2026-07-14 02:56 — Public Launch: Website, README, GitHub, DNS

> Session goal: Make cmdguard superb as a public repo across README, website, GitHub metadata, and DNS/Firebase hosting.

---

## A) FULLY DONE

### 1. Website Infrastructure (Astro + Starlight + Firebase)

- **Full website scaffold** at `website/` — 55+ files
- Astro 7 + Starlight + Tailwind v4 + Nix flake (identical stack to gogenfilter/go-atomic-write)
- `astro.config.mjs` with CSP, sitemap, prefetch, font optimization (Space Grotesk + JetBrains Mono)
- `firebase.json` with security headers (HSTS, X-Frame-Options DENY, CORP, COOP, Permissions-Policy), cache rules, cleanUrls, `/docs/*` redirect
- `.firebaserc` mapping `cmdguard` target to `lars-software` project
- `flake.nix` with `dev`/`build`/`preview`/`deploy` apps + treefmt
- `scripts/fix-csp.mjs` post-build CSP hash patching (copied from gogenfilter)
- `tsconfig.json`, `.htmlvalidate.json`, `.node-version`, `.gitignore`

### 2. Landing Page

- `HeroSection.astro` — animated code preview with syntax highlighting, GitHub stars badge, `go get` command, copy-to-clipboard
- `FeatureGrid.astro` — 8 feature cards (Zero Panics, Type-Safe Flags, DI, Constructor Validation, Config Files, 16 Output Formats, Rich Flag Tags, Cobra Escape Hatch)
- `ComparisonSection.astro` — Cobra-vs-cmdguard table (6 rows: usage-on-error, error output, exit codes, panics, flag lookups, missing handler)
- `UseCasesSection.astro` — 3 cards (Production CLIs, DevTools, Microservice CLIs)
- `CTASection.astro` — Read the Docs + View on GitHub buttons
- `LandingLayout.astro` — full SEO (OG/Twitter meta, JSON-LD SoftwareApplication schema, canonical, manifest)
- Indigo/violet accent theme (#818cf8) — distinct from gogenfilter's cyan (#22d3ee)

### 3. Starlight Documentation (16 pages)

- **Getting Started:** Installation, Quick Start
- **Guides (8):** Type-Safe Flags, Dependency Injection, Error Handling, Config Files, Rich Output, Middleware, Sub-Modules, Migrating from Cobra
- **API Reference:** Overview (constructors, options, error types, built-in commands)
- **Community:** Changelog, Contributing, Related Tools
- All content written from scratch — adapted from README/AGENTS.md source material

### 4. README.md Polish

- Added website/docs/pkg.go.dev navigation bar under badges
- Added coverage badge (87.6%)
- Replaced documentation section with links to hosted docs site + retained local docs references

### 5. GitHub Metadata

- Description updated (unchanged — already good)
- Homepage URL set to `https://cmdguard.lars.software`
- Added 6 new topics: `cobra-wrapper`, `zero-panic`, `samber-do`, `go-cli`, `functional-options`, `exit-codes`
- Total topics: 20 (was 14)

### 6. Firebase Hosting

- Created hosting site `cmdguard` in project `lars-software`
- **Deployed successfully** — `https://cmdguard.web.app` returns HTTP 200
- 71 files uploaded, CSP patched on all 16 HTML pages

### 7. DNS Record

- Added `cmdguard` CNAME → `cmdguard.web.app.` to `lars.software.tf` in domains repo
- Terraform formatting verified
- Record is staged, ready for `terraform apply`

---

## B) PARTIALLY DONE

### 1. DNS Propagation

- CNAME record added to Terraform but **NOT applied** (`terraform apply` failed — Namecheap API IP not whitelisted in this environment)
- `cmdguard.lars.software` does NOT resolve yet
- `cmdguard.web.app` IS live and serving content

### 2. README "Documentation" Links

- Links point to `cmdguard.lars.software` which doesn't resolve yet (will work after DNS apply)
- The `cmdguard.web.app` URL would work today but is less polished

### 3. Website CI/CD Pipeline

- No GitHub Actions workflow for auto-deploying the website on push
- gogenfilter has a CI workflow for this; cmdguard website deployment is currently manual (`nix run .#deploy`)

---

## C) NOT STARTED

### 1. Website Favicon/Logo Quality

- Favicon is a basic SVG placeholder (shield shape with "CG" text) — functional but not designed
- No OG image (uses favicon.svg as fallback in og:image meta tag)
- No app icons beyond the SVG favicon (no PNG icons for social media cards)

### 2. GitHub Social Preview Image

- No social preview image uploaded to GitHub repo settings
- The repo will show a generic preview when shared on social media

### 3. Website GitHub Actions CI

- gogenfilter and go-atomic-write both have website CI workflows
- cmdguard website has no CI — builds and deploys are manual

### 4. Firebase Custom Domain Connection

- Even after DNS CNAME is applied, Firebase needs to be told to connect the custom domain `cmdguard.lars.software` to the hosting site
- This is done via `firebase hosting:sites:create` or the Firebase console
- Not started — the Firebase site only responds to `cmdguard.web.app` currently

### 5. `dependents.astro`-style Page

- gogenfilter has a "Who Uses gogenfilter" page at `/dependents`
- cmdguard has no equivalent page listing projects that depend on it

### 6. Website `.gitignore` for `package-lock.json`

- `package-lock.json` is untracked (shown as `??` in git status)
- gogenfilter tracks `package-lock.json`; go-atomic-write does not
- Decision needed: track or gitignore

### 7. Go Workspace Integration

- The website directory was `git add`-ed to make it visible to Nix
- `website/node_modules/` and `website/dist/` are properly gitignored
- `website/flake.lock` was auto-generated and is tracked

---

## D) TOTALLY FUCKED UP

### 1. `hero-code.ts` Left Dead Import

- Initial version had `import { siteConfig }` + `const _ = siteConfig` + redundant `export { heroCode }` causing a duplicate export error
- Fixed during build, but this was sloppy — should have been clean from the start
- **Impact:** Wasted one build cycle. Fixed and working now.

### 2. No Staging/Preview Before Deploy

- Deployed directly to production Firebase hosting without previewing the built site locally first
- The site works, but this violates best practices — should have run `astro preview` to verify before `firebase deploy`

### 3. Changelog.mdx Referenced Non-Existent Component

- Initially wrote `import { Changelog } from '../../components/Changelog'` in changelog.mdx — the component didn't exist
- Fixed by removing the import, but this was a copy-paste error that should have been caught

---

## E) WHAT WE SHOULD IMPROVE

### Website Quality

1. **Custom-designed logo** — The SVG favicon is a placeholder. A real logo would significantly improve perceived quality.
2. **OG image generation** — gogenfilter uses `astro-og-canvas` for dynamic OG images. cmdguard has none.
3. **Lighthouse audit** — Never ran Lighthouse against the deployed site. Should verify performance/accessibility/SEO scores.
4. **Comparison page depth** — The comparison table is 6 rows. Could expand with Kong, urfave/cli, go-flags comparisons (docs/COMPARISON.md already has this content).
5. **Interactive examples** — No runnable examples or playground. The taskctl example is impressive but not showcased.
6. **Search functionality** — Starlight's built-in Pagefind search is configured, but never verified it works correctly.
7. **Analytics** — No analytics (gogenfilter removed Plausible, so maybe intentional, but worth deciding).

### README Quality

8. **Website link prominence** — The website link could be more prominent (badge-style or hero banner).
9. **Animated demo** — No GIF/video showing cmdguard in action in the README.
10. **Contributors section** — No contributors list or acknowledgment.
11. **"As seen on" / trust signals** — No badges for Go Report Card score, specific coverage %, etc.

### Operational

12. **Website CI workflow** — Need a GitHub Actions workflow that builds and deploys the website on push to `website/`.
13. **Firebase custom domain** — Must connect `cmdguard.lars.software` in Firebase console after DNS propagates.
14. **DNS apply** — User must run `terraform apply` from a whitelisted IP.
15. **Staging environment** — No preview channel in Firebase for staging before production deploy.

### Content Depth

16. **More guides** — Could add guides for: BranchingFlowContext, Audit Log export, Plugin system, Doctor/Version commands, Shell completion, Typo suggestions, Signal handling patterns, WithCleanup lifecycle.
17. **API reference depth** — Current API reference is an overview table. Could have per-function documentation pages.
18. **Migration guide** — The "Migrating from Cobra" guide is good but could include a full before/after code diff of a real app.
19. **Changelog** — Only shows v3.0.0 highlights. Could include the full CHANGELOG.md content or more versions.

---

## F) Up to 50 Things to Get Done Next

### Immediate (Blocks `cmdguard.lars.software`)

1. **Apply Terraform DNS** — Run `terraform apply` in domains repo from whitelisted IP to create the CNAME record
2. **Connect custom domain in Firebase** — Run `firebase hosting:channels:create` or console: add `cmdguard.lars.software` to the `cmdguard` hosting site
3. **Verify DNS propagation** — `dig cmdguard.lars.software` after apply, verify it resolves to Firebase
4. **Verify SSL** — Firebase auto-provisions SSL; verify `https://cmdguard.lars.software` works after DNS + custom domain connection

### Website CI/CD

5. **Create website deploy workflow** — GitHub Actions: build on push to `website/**`, deploy to Firebase hosting
6. **Add Firebase token secret** — `FIREBASE_TOKEN` to GitHub Actions secrets for deploy
7. **Add build status badge** — Link CI workflow badge in README and website
8. **Add Lighthouse CI** — Performance/accessibility regression checks on the website

### Website Polish

9. **Design a proper logo/favicon** — Commission or create a polished SVG logo for cmdguard
10. **Generate OG images** — Add `astro-og-canvas` integration for dynamic social share images
11. **Upload GitHub social preview** — Create and upload a social preview image to repo settings
12. **Add `package-lock.json` to git** — Decide: track it (like gogenfilter) or add to .gitignore (like go-atomic-write)
13. **Run Lighthouse audit** — Verify performance, accessibility, SEO, best practices scores
14. **Verify Pagefind search** — Confirm Starlight search works on the deployed site
15. **Add error page styling** — Verify the 404 page looks good (it was built but never checked)
16. **Mobile responsive audit** — Test landing page on mobile viewport sizes
17. **Add a "Why cmdguard?" section to the website** — The README has a great comparison table; the website landing page should echo it
18. **Fix og:image fallback** — Currently uses favicon.svg; should use a proper OG image

### README Improvements

19. **Add website link badge** — Add a "Website" badge in the badge row for visibility
20. **Add animated demo** — GIF or asciinema showing cmdguard in action
21. **Add "Who uses cmdguard?" section** — Once there are users
22. **Add Go Report Card dynamic badge** — Replace static coverage badge with a live one
23. **Tighten Quick Start** — The Quick Start in README is good but long; could link to website for the full version
24. **Add badges row for sub-module status** — Show the 5 sub-modules are independently versioned

### Documentation Content

25. **Expand guides** — Add guides for BranchingFlowContext, Audit Log, Plugin system, Doctor command, Shell completion, WithCleanup
26. **Add a full tutorial** — Port docs/TUTORIAL.md to the website as a multi-page guide
27. **Add per-function API docs** — Expand the API reference into per-function pages
28. **Port the framework comparison** — docs/COMPARISON.md content to a website page comparing Kong, urfave/cli, go-flags
29. **Add benchmarks page** — Port docs/PERFORMANCE.md to the website
30. **Add CLI design principles** — Port docs/CLI_DESIGN_PRINCIPLES.md to the website
31. **Add full changelog** — Include more version history beyond v3.0.0
32. **Add migration guide** — Port docs/MIGRATION_v2_v3.md to the website
33. **Add Cobra footguns doc** — Port docs/COBRA_FOOTGUNS.md to the website

### Code & Repo

34. **Add `website/` to `.buildflow.yml`** — If BuildFlow is configured, add website build/lint to the pipeline
35. **Update AGENTS.md** — Document the website directory, build commands, and deploy process
36. **Update root `flake.nix`** — Consider adding a `website` app to the root flake for convenience
37. **Add `.github/workflows/website.yml`** — Separate CI for website builds
38. **Remove stale docs** — `docs/status/` has the jsonv2 migration status doc; clean up or archive
39. **Add `website/README.md`** — Document the website build/deploy process within the website directory itself

### Marketing & Discoverability

40. **Write a launch blog post** — Content for Reddit r/golang, Hacker News, etc.
41. **Add to Awesome Go lists** — Submit to cobra-awesome, awesome-go-cli
42. **Create pkg.go.dev landing** — Ensure pkg.go.dev renders well (check module path, readme rendering)
43. **Add to Go project wiki** — Submit to github.com/golang/go/wiki/Projects
44. **Social media announcement** — Prepare Twitter/X/LinkedIn content
45. **Add "Sponsor" button** — `.github/FUNDING.yml` for GitHub Sponsors

### Sub-Modules

46. **Sub-module docs pages** — Add dedicated docs pages for each sub-module (glamour, manpage, prompts, spinner, telemetry)
47. **Sub-module version badges** — Show current version of each sub-module in the docs
48. **Sub-module examples** — Add standalone example apps for each sub-module

### Firebase & Infrastructure

49. **Set up Firebase preview channels** — For staging deployments on PRs
50. **Configure Firebase cache optimization** — Review cache headers for Starlight assets
51. **Add `_acme-challenge.cmdguard` TXT record** — May be needed for Firebase SSL cert verification (Firebase will provide the value after custom domain connection)

---

## G) Top 2 Questions I Cannot Answer Myself

### 1. Should the website live at `cmdguard.lars.software` or a custom domain like `cmdguard.dev`?

The current plan uses `cmdguard.lars.software` (consistent with gogenfilter, go-atomic-write). However, for a public open-source library aimed at wide adoption, a dedicated domain like `cmdguard.dev` might be more memorable and professional. **Do you want to register/use a separate domain, or stick with the `lars.software` subdomain pattern?**

### 2. Should `package-lock.json` be committed?

gogenfilter tracks it; go-atomic-write does not. Both approaches work. The npm docs recommend committing it for applications (which this is). **What's your preference?**

---

## Appendix: Post-Report Fix — `animations.js` (02:58)

### Bug: Entire page below hero section was pure black

**Root cause:** `public/js/animations.js` was written with the nav-toggle code (a duplicate of `header.js`) instead of the `IntersectionObserver` scroll-fade handler. The CSS rule `[data-animate]:not(.animate-fade-in) { opacity: 0 }` hides all animated sections by default. Without the observer adding the `animate-fade-in` class as elements scroll into view, every section after the hero (FeatureGrid, ComparisonSection, UseCasesSection, CTASection, Footer) stayed at `opacity: 0` — invisible against the black background.

**Fix:** Replaced `animations.js` with the correct `IntersectionObserver` implementation (matching gogenfilter's working version):

```js
if (!window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add("animate-fade-in");
          observer.unobserve(entry.target);
        }
      });
    },
    { threshold: 0.1 },
  );
  document.querySelectorAll("[data-animate]").forEach((el) => observer.observe(el));
} else {
  document.querySelectorAll("[data-animate]").forEach((el) => el.classList.add("animate-fade-in"));
}
```

**Redeployed:** `npm run build && firebase deploy --only hosting:cmdguard` — 3 files changed, live at `cmdguard.web.app`.
