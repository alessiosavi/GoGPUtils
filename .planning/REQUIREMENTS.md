# Requirements: Text Normalization Pipeline

**Defined:** 2026-04-02
**Core Value:** Turn messy text into consistent, search-friendly output through a fluent, explicit pipeline that is easy to reuse and hard to misuse.

## v1 Requirements

Requirements for the initial release. Each maps to a roadmap phase.

### Pipeline Core

- [ ] **PIPE-01**: User can build an immutable normalization pipeline by chaining stages.
- [ ] **PIPE-02**: User can run a pipeline on input text and get deterministic output or an error.
- [ ] **PIPE-03**: Stage execution preserves declaration order and never reorders steps implicitly.

### Unicode Cleanup

- [ ] **UNIC-01**: User can normalize Unicode text and remove combining marks for accent-insensitive matching.
- [ ] **UNIC-02**: User can trim and collapse whitespace, and sanitize invalid UTF-8 and NUL bytes.
- [ ] **UNIC-03**: User can apply case folding or lowercasing as part of a normalization pipeline.
- [ ] **UNIC-04**: User can filter or map characters or runes by class when building canonical output.

### Token Pipelines

- [ ] **TOKN-01**: User can split normalized text into tokens.
- [ ] **TOKN-02**: User can map and filter tokens before joining them back together.
- [ ] **TOKN-03**: User can join tokens with a caller-chosen separator.
- [ ] **TOKN-04**: User can define a stable token-boundary strategy for whitespace-first normalization.

### Presets

- [ ] **PRES-01**: User can use a search-oriented preset for database and search keys.
- [ ] **PRES-02**: User can use a canonicalization preset for general text normalization.
- [ ] **PRES-03**: User can use a DB-safe preset that keeps output valid for persistence and indexing.
- [ ] **PRES-04**: User can opt into width folding for mixed-width scraped text when needed.

## v2 Requirements

Deferred to a later release.

### Adapters

- **ADPT-01**: User can adapt a pipeline to streaming interfaces for large inputs.
- **ADPT-02**: User can trace or explain which stages changed the text.
- **ADPT-03**: User can transliterate to ASCII when explicitly requested.
- **ADPT-04**: User can apply locale-specific rules explicitly.

## Out of Scope

Explicitly excluded for now.

| Feature | Reason |
|---------|--------|
| Full NLP, stemming, or lemmatization | This library normalizes text, it does not analyze language. |
| HTML parsing or DOM extraction | Markup cleanup belongs upstream or in a separate adapter. |
| Search indexing, ranking, or query execution | The library prepares text; it does not store or score it. |
| Hidden global config or mutable shared state | Breaks determinism and makes pipelines hard to reuse safely. |
| Default ASCII transliteration | Too destructive to make implicit. |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| PIPE-01 | Phase 1 | Pending |
| PIPE-02 | Phase 1 | Pending |
| PIPE-03 | Phase 1 | Pending |
| UNIC-01 | Phase 1 | Pending |
| UNIC-02 | Phase 1 | Pending |
| UNIC-03 | Phase 2 | Pending |
| UNIC-04 | Phase 2 | Pending |
| TOKN-01 | Phase 2 | Pending |
| TOKN-02 | Phase 2 | Pending |
| TOKN-03 | Phase 2 | Pending |
| TOKN-04 | Phase 2 | Pending |
| PRES-01 | Phase 2 | Pending |
| PRES-02 | Phase 2 | Pending |
| PRES-03 | Phase 2 | Pending |
| PRES-04 | Phase 2 | Pending |

**Coverage:**
- v1 requirements: 15 total
- Mapped to phases: 15
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-02*
*Last updated: 2026-04-02 after initial definition*
