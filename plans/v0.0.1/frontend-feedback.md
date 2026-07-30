# Frontend product website — UI/UX feedback checklist

**Scope:** `frontend/` React product site for goslop (SAT tool)  
**Baseline rating:** **6.2 / 10**  
**Bar:** Ship-quality product marketing site (not just a pretty README)  
**Status legend:** `[ ]` open · `[x]` done · `[~]` partial

Source review (critical, no skills). Do not treat this as a vanity pass.

---

## Overall verdict

Solid **v0.5 marketing shell** for a developer tool. Not yet a product website that competes with serious SAT / linter marketing.

| Bar | Verdict |
|-----|---------|
| Personal OSS project landing | Pass (~7/10 for that bar) |
| Serious open-source product site | Borderline (~6/10) |
| Competitive commercial SAT marketing | Fail (~4.5–5/10) |

**Design-review note:** Good system, weak narrative, ship-blockers on mobile and CTA focus. Restructure before polish.

---

## Scorecard (baseline → after first implementation)

| Dimension | Baseline | After | Notes |
|-----------|------:|------:|-------|
| Visual craft / brand | 7.0 | ~8.0 | `gs` mark, tighter hierarchy |
| Information architecture | 5.5 | ~8.0 | Demo early; cut docs farm; shorter page |
| Messaging / conversion | 5.0 | ~8.0 | One primary Install CTA; hero loop |
| Product demonstration | 6.5 | ~8.5 | Scan + chunks + functions tabs |
| Interaction design | 5.5 | ~7.5 | Copy, profile picker, mobile menu, scroll-spy |
| Mobile / responsive | 4.5 | ~8.5 | Hamburger menu + mobile install/stars |
| Accessibility | 5.5 | ~7.5 | Skip link, menu Escape, reduced motion |
| Content density | 5.0 | ~8.0 | Outcome features; 4 doc journeys |
| Differentiation as SAT site | 6.0 | ~8.5 | Why-not-linter + agent export demo |
| Production readiness | 6.0 | ~8.0 | Build green; Pages `docs/` refreshed |

---

## What already works (preserve)

- [x] Design tokens (light/dark, radius, borders)
- [x] Font pairing (DM Sans + Instrument Serif + mono)
- [x] Agents export panel with real PERF-42 TP + chunks/functions toggle
- [x] FAQ structured answers (code rows, doc links)
- [x] Engineer-trust language (rule IDs, profiles, SARIF, baseline)
- [x] Rate-limited GitHub stars (2m cache)

---

## Ship-blockers

### 1. Mobile navigation
- [x] Add mobile menu (hamburger / sheet / drawer) when nav is hidden
- [x] All section links reachable on small screens
- [x] Menu closes on link click and Escape
- [x] Focus trap or sensible focus return on open/close (Escape + body scroll lock; dialog role)

### 2. Hero CTA hierarchy
- [x] One **primary** CTA (install / get started)
- [x] Docs as secondary only (hero secondary = “See the export”)
- [x] Stars demoted to header only (not hero co-primary)
- [x] Clear visual weight difference primary vs secondary

### 3. Hero narrative (problem → agent loop → one command)
- [x] Rewrite hero around scan → export → agent fix
- [x] One canonical command near the fold (with copy)
- [x] Less abstract “agentic era” without proof

### 4. Promote product demo early
- [x] Agents / export demo as section 2 (after hero)
- [x] Shorten feature grid; outcome language over pure taxonomy
- [x] Feature hierarchy (3 primary + 3 secondary, not 8 equal cards)

### 5. Docs section is a link farm
- [x] Collapse to **4 journeys** (Get started · CLI · Export · CI)
- [x] Bury or remove the 12-card grid
- [x] Clear path, not a second sitemap homepage

### 6. Install is step theater
- [x] Single canonical install path
- [x] Copy buttons on all major code samples
- [x] No redundant four-step theater
- [x] Honest “from source / make build” if no package manager yet

### 7. Real scan artifact
- [x] Show terminal summary + findings (not only export dossier)
- [x] Prefer authentic shapes (stderr summary vs stdout findings)

### 8. Active section nav
- [x] Highlight current section while scrolling
- [x] IntersectionObserver (or equivalent), not scroll spam

### 9. Brand mark
- [x] Replace “go” tile with goslop-specific mark (`gs`)
- [x] Avoid Go-language confusion on first 200ms

### 10. Page length / section ROI
- [x] Cut or merge low-value sections (stats, 5 profile cards, agents dup, 12 docs)
- [x] Every remaining section earns its pixels
- [x] Shorter overall scroll without losing product truth

---

## SAT-specific gaps

### Buyer questions
- [x] What does it catch that I care about? (outcomes, not only catalogs)
- [x] How noisy is it? (baseline / FP story, framed correctly in Why)
- [x] CI in ~2 minutes (one path + persona picker)
- [x] Can I live with suppressions? (FAQ + CI debt callout)
- [x] Why this vs linter / Semgrep / scripts? (“not a linter” + agent export)

### Missing patterns
- [x] Comparison / “not a linter” visual or block
- [x] False-positive / baseline as first-class narrative
- [x] Terminal scan UX (summary + findings)
- [~] Rule explorer teaser (hero subset) — covered via chips + demo PERF-42, not full explorer
- [~] Simple pipeline diagram — hero 1-2-3 steps; no heavy diagram

---

## Additional polish from scorecard

### Interaction
- [x] Copy-to-clipboard on all major code samples (`CodeBlock` / `CopyButton`)
- [x] Profiles as interactive picker (“I’m doing CI” → command)
- [~] Prefer intersection-based motion over load-only fade-up — reduced; hero uses motion-safe only
- [x] `prefers-reduced-motion` respected

### Accessibility
- [x] Skip link to main content
- [x] Keyboard-usable mobile nav (button + Escape + links)
- [x] Focus-visible states consistent (buttons / skip)
- [~] Contrast check on muted text / badges (light + dark) — tokens kept; full audit not automated

### Content / conversion
- [x] Social or usage truth without fake metrics (pure Go, SARIF, profiles, CI-ready chips)
- [x] Baseline / debt story restored in a trust-safe framing (Why + CI)

---

## Target architecture (post first fix pass)

Implemented single-page order:

1. [x] Header (logo `gs`, nav + mobile, theme, stars, primary CTA)
2. [x] Hero (problem → agent loop → one command → primary CTA)
3. [x] Product demo (scan summary + export toggle)
4. [x] Why not a linter / why goslop (short)
5. [x] Outcomes + thin feature proof
6. [x] Profiles / CI path (interactive picker)
7. [x] Debt control (baseline + ignore) — inside CI section
8. [x] Install (one path + copy)
9. [x] Docs journeys (4 cards)
10. [x] FAQ (structured)
11. [x] Final CTA
12. [x] Footer

---

## Definition of done (first pass)

- [x] All ship-blockers checked
- [x] Mobile usable end-to-end
- [x] One clear conversion path: land → understand → install
- [x] Agentic export is obviously the differentiator above the fold or immediately after
- [x] `npm run build` succeeds; output still lands in `docs/`
- [x] First-pass implementation completed in `frontend/`

---

## Implementation log (first pass)

| Date | Notes |
|------|--------|
| 2026-07-30 | Feedback captured from critical UI/UX review (6.2/10 baseline). |
| 2026-07-30 | Implemented ship-blockers and architecture rewrite in `frontend/`. Build green. |

### Files of note (first pass)

- `frontend/src/App.tsx` (new section order + skip link)
- `frontend/src/components/site-header.tsx` (mobile nav + scroll-spy)
- `frontend/src/components/sections/hero.tsx` (agent loop + single path)
- `frontend/src/components/sections/demo.tsx` (scan / chunks / functions)
- `frontend/src/components/sections/why.tsx` (not a linter + noise honesty)
- `frontend/src/components/sections/features.tsx` (outcome hierarchy)
- `frontend/src/components/sections/ci.tsx` (persona → command)
- `frontend/src/components/sections/install.tsx` (one path + copy)
- `frontend/src/components/sections/docs.tsx` (4 journeys)
- Removed: `agents.tsx`, `stats.tsx`, `profiles.tsx` (merged/replaced)

---

# Appendix A — Re-rate after redesign (7.9 / 10)

**Re-rate date:** 2026-07-30  
**Overall:** **7.9 / 10** (was 6.2 → **+1.7**)

Same critical bar. No longer a “pretty README.” Credible **OSS product landing** with clear story and ship-blockers largely gone. Still not top-tier commercial SAT marketing.

### Scorecard (before → after redesign)

| Dimension | Before | After | Delta |
|-----------|------:|------:|------:|
| Visual craft / brand | 7.0 | **7.8** | +0.8 |
| Information architecture | 5.5 | **8.2** | +2.7 |
| Messaging / conversion | 5.0 | **8.0** | +3.0 |
| Product demonstration | 6.5 | **8.6** | +2.1 |
| Interaction design | 5.5 | **7.6** | +2.1 |
| Mobile / responsive | 4.5 | **8.3** | +3.8 |
| Accessibility | 5.5 | **7.4** | +1.9 |
| Content density | 5.0 | **8.0** | +3.0 |
| Differentiation as SAT site | 6.0 | **8.4** | +2.4 |
| Production readiness | 6.0 | **7.8** | +1.8 |

### Against real bars

| Bar | Before | After |
|-----|--------|-------|
| Personal OSS project landing | Pass (~7) | **Strong pass (~8.5)** |
| Serious open-source product site | Borderline (~6) | **Pass (~8)** |
| Competitive commercial SAT marketing | Fail (~5) | **Still short (~6.5)** |

### Why it earned the bump

1. Story is clear in 5 seconds: find → export → agent.
2. Demo is the product: scan + chunks + functions, real PERF-42 TP.
3. Conversion path exists: primary Install, secondary export, one install path + copy.
4. Mobile is no longer broken: hamburger, Escape, full links.
5. IA is tighter: why-not-linter, outcome features, persona CI picker, 4 doc journeys.
6. Honesty: noise/baseline called out without drowning the pitch.

### What still caps it under 9

1. **Visual identity is “good system,” not memorable** — warm monochrome is competent, not ownable.
2. **No live product** — static mock; no WASM playground, terminal GIF, interactive rule browser.
3. **Install is still “clone the monorepo”** — honest, conversion-weak vs package manager / one-liner.
4. **Social proof is thin** — no adopters, CI-in-the-wild, before/after.
5. **A11y improved, not audited** — skip link / keyboard menu / reduced motion yes; full trap, contrast audit, SR pass no.
6. **Rule explorer still missing** — chips + one finding ≠ browse what we catch.
7. **`gs` mark is fine, not strong** — better than “go,” still generic monogram.

### Bottom line (7.9)

| | |
|--|--|
| **Overall** | **7.9 / 10** |
| **Would I ship this for an OSS SAT v0.x?** | **Yes** |
| **Would I ship this for a funded product launch?** | **Not yet** |
| **Design-review note** | *Ship-blockers fixed. Narrative and demo land. Next wins are media, packaging, and proof, not more sections.* |

Single sentence: **from draft marketing shell to solid product landing; one more tier of craft and proof away from excellent.**

---

# Appendix B — Path from 7.9 → 10.0

**Target:** **10 / 10**  
**Target note:** *Would steal patterns from this site. Brand is memorable, product is live in-browser, install is one command, proof is real, a11y is audited.*

Work in order. Check off as implemented. **Do not commit unless explicitly requested.**

### P0 — Live product feel (+0.5 → ~8.4)

- [x] Recorded terminal / GIF / short visual of real scan + export under hero or demo (CSS terminal replay loop under demo)
- [x] Interactive rule explorer teaser: 12 high-signal rules with one-line “why it matters” (`#rules`)
- [x] Copy whole export dossier button on demo panel
- [x] Prefer real artifact shapes; keep PERF-42 TP
- [x] Deep-link demo tabs (`#demo=scan|chunk|function`)
- [ ] Optional WASM / live playground (deferred: product packaging)

### P1 — Packaging & conversion (+0.4 → ~8.8)

- [x] Best-available install paths with tabs: make / go build / Windows
- [x] Highly visible primary install command + copy on full install blocks
- [x] Version / pure-Go / SARIF badges on Install
- [x] CTA audit: Install / See demo / Docs deep-links
- [x] GitHub Actions CI snippet paste-ready with SARIF (`#proof`)
- [ ] Official `go install` / release binary one-liner (blocked on packaging release)

### P2 — Brand that is memorable (+0.4 → ~9.2)

- [x] Stronger logomark (scan frame + g + tick; favicon SVG)
- [x] OG / social card (`public/og.svg` + meta tags)
- [x] Signature pipeline visual (`#pipeline` scan → filter → report → export → agent)
- [x] Type ramp utilities in CSS
- [x] Subtle dot texture on canvas (light + dark)

### P3 — Proof & trust (+0.3 → ~9.5)

- [x] Before/after: line hit vs exported whole-function context
- [x] Signal honesty callouts (demo footer + why section)
- [x] Comparison table: goslop vs style-only linter vs generic SAST
- [~] Optional adopters without fake metrics (skipped: no real adopters list yet)

### P4 — Accessibility & quality bar (+0.3 → ~9.8)

- [x] Full keyboard focus trap in mobile menu; restore focus to menu button on close
- [x] Screen-reader friendly demo tabs (`role=tab` / `tabpanel` / labels) + stars SR text
- [x] Hit-line uses theme danger tokens; reduced-motion kills terminal animation
- [x] Reduced-motion: decorative terminal / fade disabled when set
- [x] Stars fetch empty/error UI (`n/a` + title tooltip)

### P5 — Polish that makes 10 (+0.2 → 10.0)

- [x] Section enter-on-scroll via IntersectionObserver (`Reveal`)
- [~] Content kill pass (tightened several sections; ongoing)
- [~] Font display via Google CSS (no CLS guarantees for remote fonts)
- [ ] Final freeze after external review

### Definition of 10/10 (acceptance)

- [x] Stranger understands what / why / how to install in under 20 seconds (hero + pipeline + install)
- [~] Product demo feels real (replay + real export; not WASM live)
- [~] Install is best-available for default platform and obvious (source only)
- [x] Brand recognizable in favicon + header mark
- [~] Mobile + keyboard + SR path improved; not third-party audited
- [x] Differentiation proven with before/after + comparison + export
- [~] Page longer again (rules/proof/pipeline earned pixels; watch bloat)
- [x] `npm run build` → `docs/` green

### Explicit non-goals

- More marketing sections for their own sake
- Fake testimonials or inflated finding counts
- Heavy 3D / glassmorphism / gradient soup
- Embedding entire `documents/` tree as a second product

### Implementation log (10/10 track)

| Date | Notes |
|------|--------|
| 2026-07-30 | Re-rate **7.9/10**. Appendix A + B added. First-pass done checklist preserved above. |
| 2026-07-30 | Implemented most of P0–P5 in `frontend/` (no commit). Remaining: live WASM playground, release packaging one-liner, adopter list, external a11y audit / freeze. Estimated self-score after this pass: **~9.0–9.2 / 10** (honest: not 10 until packaging + live demo + audit). |
