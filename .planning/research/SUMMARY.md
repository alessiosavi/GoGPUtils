# Text Normalization Pipeline — Research Summary

## Executive Summary

This project should be a new, dedicated Go package for composing deterministic text-normalization pipelines, not an extension of `stringutil`. The core use cases are search/DB canonicalization and scraping cleanup: take messy Unicode text, clean it in a predictable order, and return stable normalized output for indexing, dedupe, and matching.

The recommended build approach is a small immutable fluent pipeline backed by `golang.org/x/text` (`transform`, `unicode/norm`, `runes`, optional `width`) on the current Go 1.25.0 toolchain. Start with explicit whole-string stages, then add token-level split/map/filter/join stages, then publish thin presets. Keep the package library-only, Unicode-first, and strict about ordering and state.

Main risks are scope creep and correctness drift: do not let this become a mini-framework, an NLP stack, or an HTML parser. Also avoid aggressive default normalization that causes collisions; destructive transforms must be explicit and searchable presets should be separate from general-purpose cleanup.

## Key Findings

### STACK.md

- **Go 1.25.0** — keep the package aligned with the repo’s existing toolchain and stdlib-first style.
- **`golang.org/x/text` v0.35.0** — the right foundation for Unicode normalization and composable transforms.
- **`transform` / `norm` / `runes`** — use these instead of hand-rolled rune logic for NFC/NFKC/NFD/NFKD and diacritic removal.
- **`width` (optional)** — only for explicit fullwidth/halfwidth folding when scraped or East Asian text needs it.
- **Stdlib only for non-Unicode cleanup** — `strings`, `unicode`, `utf8`, `html`, `testing` are enough for whitespace, validation, and entity decoding.

### FEATURES.md

- **Table stakes:** Unicode normalization, accent removal, case folding, whitespace cleanup, invalid UTF-8/NUL sanitization, rune/class filtering, and deterministic stage order.
- **Core fluent feature:** split → map/filter → join token pipelines.
- **Differentiators:** immutable pipelines, opinionated presets (`Search`, `Canonical`, `DBSafe`), per-token processing, and optional width folding.
- **Defer:** ASCII transliteration, streaming adapters, explain/tracing mode, and locale-specific rules.
- **Hard boundary:** no NLP, no HTML parsing, no ranking/indexing, no global config, no hidden precedence rules.

### ARCHITECTURE.md

- **Shape:** one new top-level package (e.g. `textnorm/`) with an immutable `Pipeline` at the center.
- **Execution model:** stages run left-to-right, exactly as declared; no hidden reordering.
- **Boundaries:** whole-string cleanup first, token stages only after boundaries are established.
- **Implementation rule:** adapter to `x/text`, don’t duplicate Unicode tables or semantics.
- **Build order:** pipeline contract → core cleanup → Unicode layer → token pipeline → presets/adapters.

### PITFALLS.md

- **Framework creep** — resist registries, reflection, plugins, or config-driven magic.
- **Normalization drift into NLP** — keep stemming, synonyms, and language analysis out of scope.
- **Over-normalization collisions** — make destructive search-key behavior opt-in and explicit.
- **ASCII-first bugs** — add Unicode fixtures, invalid UTF-8 tests, and fuzz idempotence early.
- **HTML/scraping leakage** — parsing/extraction stays upstream; core normalization only handles text.
- **Duplicate semantics with `stringutil`** — this new package should become the canonical normalization API.

## Implications for Roadmap

### Phase 1 — Core pipeline + Unicode correctness

**Deliver:** immutable fluent `Pipeline`, explicit stage interface, ordered execution, trimming/whitespace cleanup, invalid UTF-8/NUL handling, Unicode normalization, accent removal, and core tests/fuzzing.

**Why first:** everything else depends on a trustworthy execution model and Unicode-safe primitives.

**Includes:** table-stakes normalization, deterministic ordering, canonical text cleanup.

**Avoid:** framework features, token presets, transliteration, HTML cleanup, and any hidden defaults.

**Research flag:** likely **skip** deeper research; the patterns are well established and the x/text path is clear.

### Phase 2 — Token pipeline + search/DB presets

**Deliver:** split/map/filter/join stages, case folding, punctuation/class filtering, width folding where needed, and thin presets like `Search`, `Canonical`, and `DBSafe`.

**Why second:** token behavior depends on stable core cleanup and carries the highest collision risk.

**Includes:** the main differentiator of the library: reusable fluent pipelines for search-friendly canonicalization.

**Avoid:** ASCII transliteration by default, locale magic, and any preset that silently changes too much text.

**Research flag:** **needs /gsd-research-phase** if preset semantics or collision rules are unclear.

### Phase 3 — Adapters, scale, and optional enhancements

**Deliver:** optional `transform`/`Reader`/`Writer` adapters, explain/tracing mode if still desired, and performance hardening from real corpora.

**Why last:** these are useful, but they should wrap a proven core rather than define it.

**Includes:** integration conveniences and optimization work.

**Avoid:** turning the package into a streaming framework or adding heavy transliteration databases.

**Research flag:** **needs /gsd-research-phase** for streaming/performance tradeoffs; safe to defer until after MVP usage.

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | High | Strong alignment with current repo toolchain and established `x/text` primitives. |
| Features | Medium-High | MVP boundaries are clear; preset semantics and future extras still need validation. |
| Architecture | High | Immutable fluent pipeline + ordered stages is the right shape for this domain. |
| Pitfalls | High | The main failure modes are well understood and directly shape roadmap ordering. |

### Gaps to Address

- Exact public API naming for the fluent chain.
- Final preset semantics for search/DB-safe output and collision tolerance.
- Whether width folding belongs in core or only in presets.
- Whether streaming adapters are worth the complexity before real usage proves it.

## Sources

- `.planning/PROJECT.md`
- `.planning/research/STACK.md`
- `.planning/research/FEATURES.md`
- `.planning/research/ARCHITECTURE.md`
- `.planning/research/PITFALLS.md`
