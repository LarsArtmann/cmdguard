# Frontend Design Review — cmdguard Website

**Date:** 2026-07-18
**Scope:** 14 `.astro` components, `global.css`, `starlight.css`, `LandingLayout.astro`, `index.astro`
**Mode:** Critique (per `frontend-design` skill — NOT a rebuild)
**Skill:** `frontend-design`

---

## Overall verdict

**Score: 6/10.** The site is competently built: responsive, accessible (focus-visible, skip-link, reduced-motion), dark-mode-first, and the code-in-hero is appropriate for a CLI library. But it reads as **the 2024-2026 developer-tool default**: dark stone background, indigo accent, Space Grotesk, centered everything, gradient-text H1, traffic-light code preview. A screenshot of this site is indistinguishable from hundreds of other dev-tool landing pages generated in the last two years.

The subject is specific — a Go CLI framework that catches errors at construction. The design is not.

---

## The 3 AI-default looks check

Per the `frontend-design` skill, three looks cluster around AI-generated design right now:

1. **Warm cream + serif + terracotta?** No. Dark mode.
2. **Near-black + acid-green/vermilion?** Close in structure (near-black `#0c0a09`), but uses indigo instead of acid-green. The near-black + single-accent pattern is the same default, just with a different hue.
3. **Broadsheet hairline + zero radius?** No. Uses rounded corners (`rounded-xl`, `rounded-lg`).

**Verdict:** Not literally any of the three, but it occupies the equally common fourth default: **dark + stone neutrals + indigo accent + Space Grotesk** — the Vercel/shadcn aesthetic that has become the generic developer-tool look.

---

## Finding 1: The entire palette is out-of-the-box Tailwind defaults

**Current:**

| Token          | Hex                  | Tailwind name       |
| -------------- | -------------------- | ------------------- |
| bg-primary     | `#0c0a09`            | `stone-950`         |
| bg-card-solid  | `#1c1917`            | `stone-900`         |
| text-primary   | `#e7e5e4`            | `stone-200`         |
| text-secondary | `#a8a29e`            | `stone-400`         |
| text-muted     | `#9ca3af`            | `gray-400` (mixed!) |
| accent         | `#818cf8`            | `indigo-400`        |
| accent-hover   | `#a5b4fc`            | `indigo-300`        |
| accent-light   | `#c7d2fe`            | `indigo-200`        |
| border         | `rgba(41,37,36,...)` | `stone-800`         |

Every single color is a Tailwind default token. Zero custom mixing. The accent `#818cf8` is literally the most-used indigo in the Tailwind ecosystem.

**Severity:** Medium. The colors work — they're readable, the contrast is fine, dark/light mode is well-handled. But they carry zero identity. A Go library about "guarding" commands could draw from:

- Go's own brand cyan (`#00ADD8`) as a nod to the ecosystem
- A warm amber/gold (the existing `#fbbf24` is already in the palette as a secondary — promoting it would reference the "shield/guard" metaphor)
- A deep teal that reads as "safety/verified" without being generic

**Recommendation for L1.07:** Shift the accent from `indigo-400` to a custom indigo-teal blend (e.g., `#5ebibf` or a warm-leaning `#7c6fef`) — enough to break the "exact Tailwind default" pattern without a full palette overhaul.

---

## Finding 2: Typography uses the most common dev-tool pairing

**Current:** Space Grotesk (sans/display/body) + JetBrains Mono (code).

Space Grotesk is the 2024-2026 "Inter replacement" for developer tools. It appears on hundreds of CLI/SDK landing pages. JetBrains Mono is equally ubiquitous.

There is **no display face** — the H1 uses `font-bold tracking-tighter` on Space Grotesk, which is competent but generic. No serif, no wide grotesk, no contrast between the headline face and the body face.

**Severity:** Medium. The type is readable and the scale is well-set (`clamp()` for fluid sizing). But the pairing carries no personality.

**Recommendation for L1.07 (if budget allows):** Add a second face for display only — something like **Instrument Serif** or **Fragment Mono** for the H1, keeping Space Grotesk for body. If that's too risky, at minimum add a distinctive type treatment (variable-width tracking, an italic accent word) to the H1.

---

## Finding 3: No signature element — nothing is memorable in 2 seconds

Per the skill: "Signature: the single unique element this page will be remembered by."

**What exists:**

- Hero H1 with gradient-text on one word ("shutdown") — the single most copied hero pattern of 2024-2026
- Code preview with traffic-light dots + filename — the second most copied pattern
- Shield logo with "CG" monogram in a rounded square — generic
- Feature cards with outline icons — standard
- Comparison table — well-executed but standard

If you screenshot this site and put it next to 10 other dev-tool landing pages, there is **no element that identifies it as cmdguard**. The code in the preview is the most specific thing (it shows Go with `v3.NewCLI`), but the code is wrapped in the most generic container possible.

**Severity:** High. This is the core issue. Everything else is polish.

**Recommendation:** The signature element should embody the library's core value: **catching errors at construction, not at runtime**. Ideas (pick one):

- A **before/after split** in the hero: left side shows a raw cobra command with a runtime panic; right side shows the same command in cmdguard with a construction-time error. The visual contrast IS the value prop.
- A **terminal-style hero** (not a code editor card) — a real `$ myapp --invalid-flag` prompt showing a cmdguard validation error with a specific, actionable message. This grounds the site in the terminal, which is where the user actually interacts with the library.
- An **animated construction sequence**: show the CLI being assembled, with each `With*()` option slotting into place like building blocks. The "construction" metaphor is the library's name.

These are L1.07 scope or beyond; for now, the critique documents the gap.

---

## Finding 4: Every section is centered — no structural variety

**Current:** HeroSection, SectionHeader, FeatureGrid, ComparisonSection, UseCasesSection, CTASection — all use `text-center` and `mx-auto` max-width containers. The entire page is a vertical stack of centered blocks.

**Issue:** Centered sections are the default. Asymmetry, editorial columns, or left-aligned section breaks would add visual rhythm. The comparison table especially would benefit from full-width treatment (it currently sits in a `max-w-3xl` narrow container, making it horizontally cramped on wide screens).

**Severity:** Low-Medium. The centering is consistent (not random), which is better than inconsistent alignment. But it contributes to the "templated" feeling.

---

## Finding 5: The hero H1 uses gradient-text on one word — the most overused pattern

**Current:**

```html
<span>From flags to </span
><span class="bg-gradient-to-r from-accent to-accent-light bg-clip-text text-transparent"
  >shutdown</span
><span>, type-safe.</span>
```

Gradient-clip-text on a single word in the H1 is the most copied hero typography pattern of the last three years. It appears on thousands of landing pages.

**Severity:** Low. It's functional and readable. But it's the design equivalent of using `Lorem ipsum` — it signals "no one made a specific choice here."

---

## Finding 6: The logo is generic and the shield is nearly invisible

**Current:** `Logo.astro` draws a rounded-rect background in accent color, a shield outline at `rgba(255,255,255,0.4)` opacity, and "CG" in monospace text.

The shield at 40% white opacity is barely visible — the "CG" monogram dominates, but "CG" is a generic two-letter mark. The shield (which IS the right metaphor for "guard") is treated as decoration rather than the primary element.

**Severity:** Low. Small surface area, but it's the first thing in the nav bar.

---

## What works well (do not change)

1. **Code-in-hero is the right choice.** A CLI library's hero should show code. The syntax highlighting, the import path color-coding, and the copy button are all well-executed.
2. **The comparison table is honest and clear.** Dot + label per cell, with the cmdguard column highlighted via `bg-accent-dim`. Good information density.
3. **Dark-mode-first is correct** for a developer tool.
4. **Accessibility floor is solid:** focus-visible outlines, skip-to-content, aria-labels, reduced-motion respect, semantic HTML.
5. **The `go get` command in a monospace box** as a secondary CTA is a nice pattern — lets the impatient reader copy-paste immediately.
6. **Fluid type scaling** via `clamp()` on H1 and H2 — no breakpoint jankiness.
7. **The newsletter is understated** — border-bottom separator, not a loud CTA band. Good restraint.

---

## Prioritized actions for L1.07 (cap: 6 changes)

| #   | Finding                           | Change                                                                                                                                | Effort | Risk                  |
| --- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | ------ | --------------------- |
| 1   | Palette is Tailwind defaults      | Shift accent hue from `#818cf8` to a custom blend (e.g., `#6d8cef` — slightly warmer/more cyan) across `global.css` + `starlight.css` | S      | Low — single hue swap |
| 2   | Gradient-text H1 is generic       | Replace the gradient-text span with an italic + accent-color treatment on "shutdown"                                                  | XS     | Low                   |
| 3   | Logo shield invisible             | Increase shield stroke opacity from 0.4 to 0.85; shrink "CG" text to give the shield more presence                                    | XS     | Low                   |
| 4   | Code preview is generic container | Replace traffic-light dots with a terminal prompt symbol (`$` or `❯`) in accent color; rename "main.go" to a more descriptive label   | S      | Low                   |
| 5   | Centered everything               | Make the comparison section `width="wide"` instead of `narrow` so the table breathes                                                  | XS     | Low                   |
| 6   | No display face                   | Add `font-style: italic` + tighter tracking on the H1 accent word as a lighter alternative to introducing a new face                  | XS     | Low                   |

**If only 3 changes:** #1 (accent hue), #2 (kill gradient text), #3 (shield visibility). These are the highest-impact, lowest-risk fixes.

---

## What is explicitly NOT recommended (Verschlimmbesserung guards)

- **Do not** introduce a serif face without a full type-scale audit — a badly paired serif is worse than a well-set Space Grotesk
- **Do not** redesign the layout grid — the centered stack works; variety should come from the signature element, not from random asymmetry
- **Do not** add animation beyond the existing fade-in — the skill warns "extra animation contributes to the feeling that the design is AI-generated"
- **Do not** replace the code-in-hero — it is the right choice for a CLI library; it needs a better container, not removal

---

## Verdict

The site is a solid B-. It works, it's accessible, it's responsive. But it doesn't answer the question: **"What makes this cmdguard and not any other dev tool?"** The answer exists in the product (construction-time validation, zero-panic contract, DI lifecycle) — it just hasn't been translated into a visual signature. The 6 changes above are polish; the real work is finding the signature element (Finding 3), which belongs in a future design sprint.

---

_Generated by `frontend-design` skill. This is a critique, not a rebuild. Apply the bounded changes in L1.07; defer the signature-element work to a future session._
