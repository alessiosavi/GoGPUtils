# Roadmap: Text Normalization Pipeline

**Core Value:** Turn messy text into consistent, search-friendly output through a fluent, explicit pipeline that is easy to reuse and hard to misuse.

## Phase 1: Core Pipeline and Unicode Safety

**Goal:** Establish the immutable fluent pipeline and the Unicode-safe cleanup foundation.

**Requirements:** PIPE-01, PIPE-02, PIPE-03, UNIC-01, UNIC-02

**Plans:** 1 plan

**Plan List:**
- [x] `01-core-pipeline-and-unicode-safety-01-PLAN.md` — build the immutable `textnorm` pipeline and Unicode-safe cleanup stages

**Deliverables:**
- New package scaffold and public docs
- Immutable pipeline type and ordered stage execution
- Unicode normalization and accent removal
- Whitespace trimming/collapsing and UTF-8/NUL sanitization
- Core unit tests for determinism and Unicode edge cases

**Success Criteria:**
1. A caller can chain stages and get a new reusable pipeline value.
2. The same input and pipeline always produce the same output.
3. Unicode fixtures with combining marks and invalid UTF-8 pass.
4. Whitespace and database-safety cleanup behave predictably.

## Phase 2: Token Pipelines and Presets

**Goal:** Add token-level composition and the search/DB-oriented presets that make the library useful in practice.

**Requirements:** UNIC-03, UNIC-04, TOKN-01, TOKN-02, TOKN-03, TOKN-04, PRES-01, PRES-02, PRES-03, PRES-04

**Plans:** 2 plans

**Plan List:**
- [x] `02-token-pipelines-and-presets-01-PLAN.md` — add token-level stages, case/rune transforms, and width folding
- [x] `02-token-pipelines-and-presets-02-PLAN.md` — add search, canonical, and DB-safe presets plus docs/tests

**Deliverables:**
- Split/map/filter/join token stages
- Case folding and rune/class filtering stages
- Search, canonical, and DB-safe presets
- Optional width folding for mixed-width text
- Tests for token boundaries, collision risks, and preset parity

**Success Criteria:**
1. Callers can normalize text at the token level without rewriting the core pipeline.
2. Presets stay thin and reuse the same underlying stages.
3. Search-key output is stable and explicit about destructive transforms.
4. Width folding is opt-in and does not affect default behavior.

## Phase 3: Hardening and Future Adapters

**Goal:** Harden the package and leave room for future streaming or tracing features without bloating the core.

**Requirements:** none from v1; future adapters only

**Plans:** 1 plan

**Plan List:**
- [ ] `03-hardening-and-future-adapters-01-PLAN.md` — benchmark, fuzz, and document the hardened `textnorm` surface without introducing framework creep

**Deliverables:**
- Benchmarks on representative text corpora
- Fuzz tests for malformed input and idempotence
- Optional streaming adapters if they prove worth the extra surface area
- Public examples and package docs that show real-world usage

**Success Criteria:**
1. Hot paths are measured before they are optimized.
2. Edge cases are covered beyond ASCII-only fixtures.
3. No framework creep leaks into the core API.
4. The package remains small, deterministic, and easy to reason about.

## Phase Ordering Rationale

- Phase 1 comes first because every later feature depends on a trustworthy execution model.
- Phase 2 follows because token pipelines only make sense once whole-string cleanup is correct.
- Phase 3 stays last because adapters and performance work should wrap a proven core, not define it.

## Research Flags

- **Phase 1:** standard pattern, no deeper research needed.
- **Phase 2:** revisit if preset semantics or collision tolerance need more precision.
- **Phase 3:** revisit if streaming adapters become a first-class need.

---

## Phase 1 Details

**Name:** Core Pipeline and Unicode Safety

**Goal:** Establish the immutable fluent pipeline and the Unicode-safe cleanup foundation.

**Requirements:** PIPE-01, PIPE-02, PIPE-03, UNIC-01, UNIC-02

**Plans:** 1 plan

**Plan List:**
- [ ] `01-core-pipeline-and-unicode-safety-01-PLAN.md` — implement the new top-level `textnorm` package, immutable stage chaining, and Unicode-safe cleanup

**Success Criteria:**
1. Fluent chaining returns a reusable pipeline value.
2. Ordered execution is deterministic and test-covered.
3. Unicode normalization, accent removal, and cleanup behave correctly on real text.
4. The core package has no global state and no hidden stage reordering.

## Phase 2 Details

**Name:** Token Pipelines and Presets

**Goal:** Add token-level composition and search/DB-oriented presets.

**Requirements:** UNIC-03, UNIC-04, TOKN-01, TOKN-02, TOKN-03, TOKN-04, PRES-01, PRES-02, PRES-03, PRES-04

**Success Criteria:**
1. Token split/map/filter/join works as a fluent extension of the core pipeline.
2. Presets cover search, canonical, and DB-safe use cases without hiding behavior.
3. Case folding, rune filtering, and width folding are explicit and testable.
4. The package can produce stable search keys without accidental hidden transformations.

## Phase 3 Details

**Name:** Hardening and Future Adapters

**Goal:** Improve confidence and leave room for later streaming or tracing support.

**Requirements:** future-only

**Plans:** 1 plan

**Plan List:**
- [ ] `03-hardening-and-future-adapters-01-PLAN.md` — add benchmarks, fuzzing, and docs that keep streaming adapters deferred

**Success Criteria:**
1. Benchmarks and fuzz tests cover the core pipeline.
2. Optional adapters are isolated from the core API.
3. The library still feels like a small utility package, not a framework.
