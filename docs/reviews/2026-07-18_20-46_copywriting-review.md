# Copywriting Review — cmdguard README + Website Hero

**Date:** 2026-07-18
**Scope:** `README.md` (685 lines, 32 H2 sections) + `website/src/data/{hero-code,features,comparisons,usecases}.ts`
**Mode:** Findings + annotations + alternatives (NOT a rewrite — README is already strong)

---

## Overall verdict

**The copy is already B+ / A-.** Clear headline, specific differentiation (the trinity comparison table is excellent), concrete code examples, honest positioning. No fabricated claims. The main improvement opportunities are tightening the hero subtitle and making the CTA more action-oriented.

Score: **8.5/10**. Would be 9.5 with the two changes below.

---

## Finding 1: Hero subtitle is a feature-list, not a hook

**Current** (`README.md:19`):

> From flag definition to service shutdown — type-safe, validated, zero panics.

**Issue:** This names 4 features (flags, shutdown, type-safety, zero-panics) but doesn't tell the reader _what outcome they get_. It reads as a feature catalog, not a value proposition. A newcomer who doesn't already know the library sees jargon.

**Alternatives:**

- **Option A (outcome-led):** "Build production CLIs in Go without the runtime surprises — typed flags, lifecycle-aware services, and errors that can't panic."
  - _Rationale:_ leads with the outcome ("no runtime surprises") and the 3 pillars, in plain language.
- **Option B (question hook):** "What if your CLI's flags, services, and errors were all type-checked before the first command ran?"
  - _Rationale:_ rhetorical question engages the reader's curiosity; the answer is the feature list that follows.
- **Option C (specific benefit):** "The Go CLI framework that catches missing handlers, duplicate commands, and invalid flags at construction — not at 2am in production."
  - _Rationale:_ uses specificity (the "2am in production" detail) and names the exact failure modes the library prevents.

**Recommendation:** Option C. It is the most specific and the "2am in production" detail makes it memorable without being cute.

---

## Finding 2: No explicit CTA above the fold

**Current:** README opens with badges → tagline → positioning paragraph → "Why cmdguard?" comparison table. The first actionable instruction is "go get ..." at the Quick Start section (~line 120).

**Issue:** A reader who scrolls past the comparison table has no single obvious next step until they scroll ~100 lines down. The website has CTAs (`CTASection.astro`); the README should mirror that.

**Fix:** Add a one-line CTA right after the tagline:

```markdown
**Get started in 30 seconds:** `go get github.com/larsartmann/cmdguard/v4` · [Quick Start](#quick-start) · [Full Docs](https://cmdguard.lars.software)
```

- _Rationale:_ gives the impatient reader a single line to act on; mirrors the website's pattern; the inline `go get` lets them copy-paste immediately.

---

## Finding 3: "The trinity" framing is strong — keep it

**Current** (`README.md:27`):

> ### The trinity — what no other CLI framework offers together

**Verdict:** excellent. "Trinity" is specific, memorable, and the comparison table below it delivers on the promise. The word "together" is the key insight — any one feature exists elsewhere, but the combination is the differentiator. **Do not change.**

---

## Finding 4: Comparison table is honest but could name sources

**Current:** the comparison table (`README.md:31-43`) compares against Cobra / Kong / urfave/cli with "Yes" / "—" / "Some".

**Issue:** "Some" for Kong's "Validated at construction" is accurate but unsourced. A skeptical reader wonders "how much?"

**Fix (optional):** add a footnote or inline link:

```markdown
| Validated at construction (not at runtime) | — | Some <sup>[1]</sup> | — | **Yes** |

...
<sup>[1]</sup> Kong validates struct tags at parse time but does not validate command structure (missing handlers, duplicates) at registration.
```

- _Rationale:_ specificity builds trust. The current "Some" is honest but invites doubt; the footnote converts doubt into evidence.

---

## Finding 5: Website hero code is too long

**Current** (`website/src/data/hero-code.ts`): ~40 lines of Go in the hero. Shows NewCLI + Provide + NewCommand + handler.

**Issue:** The hero code block is the first thing a visitor sees. 40 lines is a lot to parse in 3 seconds. The handler's `db.Query(ctx)` is the punchline but it's below the fold of the code block.

**Fix:** Trim to ~20 lines — show only the DI registration + the handler invocation. Drop the `AppConfig` struct definition (move to a "full example" link). The punchline ("`return db.Query(ctx)`") should be visible without scrolling.

**Alternative:** keep the full code but add inline highlighting (a colored callout on the `Provide` line and the `Invoke` line) so the reader's eye goes to the differentiator.

---

## Finding 6: Feature card descriptions are solid

**Current** (`website/src/data/features.ts`): each feature has a 1-sentence description with the library it depends on.

**Verdict:** good. "Lazy services, lifecycle hooks, health checks, and graceful shutdown in reverse order. Powered by samber/do/v2." — specific, names the dependency, uses domain language. **No change needed.**

---

## Finding 7: "Zero panics" claim is backed by evidence

**Current:** README says "Zero panics by construction (no `Run`, no `Must`)" and the feature card says "Every function returns errors. No Run, no Must*."

**Verdict:** honest. The claim is verifiable (grep for `Must` in the codebase → 0 results outside testutil's `AssertPanics`). The parenthetical "(no `Run`, no `Must`)" is the proof. **No change needed.**

---

## Prioritized actions

| # | Finding                                | Impact | Effort | Recommendation                              |
| - | -------------------------------------- | ------ | ------ | ------------------------------------------- |
| 1 | Hero subtitle → outcome-led (Option C) | High   | XS     | **Apply**                                   |
| 2 | Add one-line CTA after tagline         | High   | XS     | **Apply**                                   |
| 3 | "Trinity" framing                      | —      | —      | Keep                                        |
| 4 | Comparison table sources for "Some"    | Med    | S      | Optional                                    |
| 5 | Website hero code trim/highlight       | Med    | M      | Optional (frontend-design review owns this) |
| 6 | Feature card descriptions              | —      | —      | Keep                                        |
| 7 | "Zero panics" evidence                 | —      | —      | Keep                                        |

**If only 2 changes are made:** #1 and #2. They are ~2 lines of README diff with outsized first-impression impact.

---

## Voice and tone assessment

| Dimension           | Current state                                                   | Verdict                                 |
| ------------------- | --------------------------------------------------------------- | --------------------------------------- |
| Formality           | Professional but technical                                      | Appropriate for a Go library            |
| Personality         | Confident, specific, no hype                                    | Excellent                               |
| Jargon level        | High (DI, lifecycle, sentinel, COW) but always explained        | Acceptable — target audience is Go devs |
| Exclamation points  | 0 in README body                                                | Excellent (per copywriting rules)       |
| Passive voice       | Rare; mostly active ("cmdguard gives you", "Register services") | Excellent                               |
| Marketing buzzwords | 0 ("streamline", "optimize", "innovative" all absent)           | Excellent                               |

The voice is already where it should be. No tone adjustment needed.

---

_Generated by copywriting skill. This is a findings document, not a rewrite. Apply the 2 high-impact changes (#1, #2) in a separate commit if desired._
