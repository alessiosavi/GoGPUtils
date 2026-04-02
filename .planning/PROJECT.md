# Text Normalization Pipeline

## What This Is

A dedicated Go package for building composable text-normalization pipelines. It is meant for noisy input from scraping, data cleaning, and search or database indexing, where callers need deterministic, reusable normalization rather than one-off string hacks.

## Core Value

Turn messy text into consistent, search-friendly output through a fluent, explicit pipeline that is easy to reuse and hard to misuse.

## Requirements

### Validated

- [x] Build a fluent normalization pipeline API for composing reusable text transforms.
- [x] Support state-of-the-art Unicode normalization and diacritic removal for search-friendly text.
- [x] Support token-level cleanup steps such as stripping, filtering, mapping, and joining.
- [x] Make the package reusable for scraping cleanup, database search, and general text canonicalization.

### Active

- (None)

### Out of Scope

- [ ] Full natural-language processing or stemming — this library is for normalization, not linguistic analysis.
- [ ] HTML parsing and DOM extraction — callers should clean raw markup before feeding text into the pipeline.
- [ ] Search indexing or query execution — this library prepares text, it does not store or rank it.

## Context

The repository already has a `stringutil` package with basic normalization helpers, including Unicode cleanup and diacritic removal. This project should be a dedicated new package with a more composable, fluent pipeline API instead of more ad hoc helpers.

The target implementation should fit the existing repo style: package-oriented Go utilities, explicit errors over panics, zero global state, and colocated tests. The repo already uses Go 1.25.0 and `golang.org/x/text`, which is a good fit for Unicode normalization.

Phase 1 through Phase 3 are complete: `textnorm` now has the fluent core, Unicode-safe cleanup, token pipelines, presets, benchmarks, fuzz coverage, and docs that keep streaming adapters deferred.

## Constraints

- **Tech stack**: Go 1.25.0 and the existing `golang.org/x/text` dependency — keep the package lightweight and library-only.
- **Design**: Fluent, composable API with reusable pipeline stages — avoid global configuration and hidden state.
- **Compatibility**: Keep behavior deterministic and testable for search and cleaning workflows.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| New package instead of extending `stringutil` | The desired API is a dedicated fluent pipeline, not another collection of free functions | — Pending |
| Use a fluent chain model | Matches the desired ergonomics for reusable normalization pipelines | — Pending |
| Focus v1 on search/DB normalization | The strongest user value is canonicalizing text for indexing, matching, and dedupe | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-04-02 after milestone completion*
