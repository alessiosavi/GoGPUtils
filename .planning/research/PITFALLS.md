# Domain Pitfalls

**Domain:** Go text-normalization pipeline library
**Researched:** 2026-04-02

## Critical Pitfalls

### 1. Turning a normalizer into a mini-framework
**What goes wrong:** The API grows into registries, plugins, reflection, middleware, and config-driven magic instead of a small ordered pipeline.
**Why it happens:** Text cleanup has lots of edge cases, so every new request feels like it needs another abstraction layer.
**Consequences:** Hard-to-predict outputs, brittle ordering, more surface area than behavior, and a library that is harder to trust than ad hoc helpers.
**Prevention:** Keep the core as a tiny chain of pure transform functions with explicit order; no global registry, no reflection, no hidden execution graph.
**Detection:** New stages require framework concepts just to compose; users ask how to "register" transforms instead of just chaining them.
**Phase mapping:** Phase 1 should lock the minimal fluent API and refuse framework features.

### 2. Blurring normalization with NLP / linguistic analysis
**What goes wrong:** Stemming, lemmatization, synonym expansion, locale-specific "smart" rewrites, or other language semantics creep into a library meant for deterministic cleanup.
**Why it happens:** Search and scraping users want better matches, so the scope drifts from normalization into language processing.
**Consequences:** Non-deterministic or locale-dependent behavior, surprising data loss, and a package that cannot stay small or well-tested.
**Prevention:** Normalize form, not meaning. Keep stemming, translation, and language-specific analysis out of scope; if needed, make them separate packages.
**Detection:** Roadmap items mention stemming, token classification, synonyms, or "intelligent" cleanup.
**Phase mapping:** Phase 1 should state the boundary; later phases should not expand it.

### 3. Over-normalizing and creating collisions too early
**What goes wrong:** Lowercasing, accent stripping, ASCII folding, whitespace collapse, and punctuation removal happen by default, so distinct inputs map to the same output.
**Why it happens:** Search quality pushes teams toward aggressive canonicalization.
**Consequences:** False dedupe, broken identifiers, impossible round-trips, and search keys that no longer preserve enough information.
**Prevention:** Separate "search key" pipelines from user-facing text; keep destructive transforms opt-in and document collision risk loudly.
**Detection:** Different source values normalize to the same key and callers start using normalized output as an identity field.
**Phase mapping:** Phase 1 should define safe defaults; Phase 2 can add explicit presets for search/DB use.

### 4. Getting Unicode wrong by assuming ASCII rules
**What goes wrong:** Byte slicing, ASCII regexes, naive case conversion, and incomplete diacritic handling break on combining marks, emoji, RTL text, and non-Latin scripts.
**Why it happens:** Early tests are usually English-only, so Unicode edge cases stay invisible until production.
**Consequences:** Corrupted output, broken truncation, flaky matching, and "works on my machine" bugs for real-world text.
**Prevention:** Use rune-aware operations and `golang.org/x/text` for normalization; add fixtures for combining marks, emoji, CJK, and invalid UTF-8; fuzz idempotence.
**Detection:** Tests only cover ASCII; bugs appear around accents, emoji, or multi-byte truncation.
**Phase mapping:** Phase 1 must get correctness right; Phase 2 should expand test coverage and fuzzing.

### 5. Folding HTML/scraping cleanup into the normalization core
**What goes wrong:** The pipeline starts doing tag stripping, entity decoding, script/style handling, or malformed-markup recovery.
**Why it happens:** Scraped input is messy, and it feels convenient to "just clean it here".
**Consequences:** Regex-driven HTML bugs, parser responsibility creep, security assumptions, and a core library that no longer operates on plain text.
**Prevention:** Keep parsing/extraction upstream; if HTML cleanup is needed, make it a separate adapter that depends on a real parser.
**Detection:** Requests to support nested tags, broken markup, or script/style removal inside the core pipeline.
**Phase mapping:** Phase 2 or later only, and only as an adapter outside the core text-normalization path.

### 6. Making behavior depend on hidden state or option order
**What goes wrong:** Global configuration, mutable builders, map iteration, or undocumented option precedence make outputs vary with runtime details.
**Why it happens:** Options are easy to add, but easy to make ambiguous.
**Consequences:** Flaky tests, irreproducible search keys, and pipelines that are hard to reason about.
**Prevention:** Keep transforms pure, freeze config at build time, and make execution order explicit and tested.
**Detection:** Reordering options changes output unexpectedly or parallel tests start to fail intermittently.
**Phase mapping:** Phase 1 must define deterministic ordering; Phase 2 should add golden tests for full pipelines.

### 7. Duplicating semantics between the new pipeline package and existing helpers
**What goes wrong:** The new package reimplements `stringutil` behavior with slightly different rules, so there are two competing definitions of "normalize".
**Why it happens:** Existing helpers feel close enough to reuse, but pipeline semantics slowly drift.
**Consequences:** Callers get inconsistent output, migration becomes confusing, and bugs show up as "why did these two code paths disagree?"
**Prevention:** Make the new package the canonical normalization path; keep old helpers as thin wrappers or explicitly deprecated adapters with parity tests where needed.
**Detection:** Two functions with the same intent produce different normalized forms or docs cannot explain which one to use.
**Phase mapping:** Phase 1 should decide the canonical API; Phase 2 should clean up overlap.

## Moderate Pitfalls

### 1. Paying the allocation tax too early
**What goes wrong:** Every stage converts to `[]rune`, re-scans the entire string, or materializes intermediate buffers even for no-op pipelines.
**Prevention:** Add benchmarks on representative corpora, keep fast paths for empty/no-op input, and combine compatible transforms only after the semantics are locked.
**Phase mapping:** Phase 2 should introduce benchmarks and optimize only after correctness is stable.

### 2. Shipping unproven presets instead of proven composition
**What goes wrong:** The package publishes lots of convenience presets (`search`, `slug`, `db`, `scrape`) before the primitive stages are validated.
**Prevention:** Ship a small core first, then add presets only when backed by real usage and tests.
**Phase mapping:** Phase 2 can add named presets; Phase 1 should stay primitive-only.

## Minor Pitfalls

### 1. Under-testing idempotence and edge cases
**What goes wrong:** Tests cover happy-path English text but not "run it twice" behavior, invalid UTF-8, combining marks, or odd whitespace.
**Prevention:** Add idempotence checks for normalization stages, table-driven Unicode fixtures, and fuzz tests for malformed input.

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Phase 1: Minimal core pipeline | Framework creep, hidden state, Unicode mistakes | Keep only pure ordered transforms, explicit config, and Unicode-first tests |
| Phase 2: Presets and extra stages | Over-normalization, HTML cleanup creep, duplicate semantics | Keep presets explicit, separate adapters from core, and reconcile behavior with existing helpers |
| Phase 3: Hardening and scale | Allocation blowups, missing edge-case coverage | Benchmark real corpora, fuzz malformed input, and lock in golden outputs |

## Sources

- `/opt/SP/Workspace/Go/GoGPUtils/.planning/PROJECT.md`
- `/opt/SP/Workspace/Go/GoGPUtils/.planning/codebase/CONCERNS.md`
- `/opt/SP/Workspace/Go/GoGPUtils/stringutil/clean.go`
- `/opt/SP/Workspace/Go/GoGPUtils/stringutil/stringutil.go`
