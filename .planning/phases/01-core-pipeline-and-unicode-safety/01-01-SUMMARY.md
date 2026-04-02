---
phase: 01-core-pipeline-and-unicode-safety
plan: 01
subsystem: library
tags: [textnorm, unicode, pipeline, normalization, x-text]
requires:
  - phase: none
    provides: project initialization context
provides:
  - immutable fluent `textnorm` pipeline core
  - Unicode normalization and accent-removal stages
  - whitespace trimming, collapse, and UTF-8/NUL sanitization stages
  - tests covering immutability, ordering, and normalization behavior
affects: [Phase 2 token pipelines, Phase 3 hardening/adapters]
tech-stack:
  added: [golang.org/x/text]
  patterns: [immutable fluent pipeline, ordered stage execution, Unicode-first normalization, explicit sanitization stages]
key-files:
  created: [textnorm/doc.go, textnorm/pipeline.go, textnorm/stages.go, textnorm/unicode.go, textnorm/whitespace.go, textnorm/pipeline_test.go, textnorm/unicode_test.go, textnorm/whitespace_test.go]
  modified: [none]
key-decisions:
  - "Create a new `textnorm` package instead of extending `stringutil`."
  - "Use `golang.org/x/text` primitives for Unicode normalization and diacritic removal."
  - "Keep cleanup stages explicit and deterministic; no tokenization or HTML parsing in Phase 1."
patterns-established:
  - "Pattern 1: immutable pipeline values return a new chain on every append"
  - "Pattern 2: stages execute strictly left-to-right with no hidden reordering"
  - "Pattern 3: Unicode and whitespace cleanup live behind explicit pipeline stages"
requirements-completed: [PIPE-01, PIPE-02, PIPE-03, UNIC-01, UNIC-02]
duration: 1h 10m
completed: 2026-04-02
---

# Phase 1: Core Pipeline and Unicode Safety Summary

**Immutable `textnorm` pipeline core with Unicode-safe cleanup and deterministic stage execution**

## Performance
- **Duration:** 1h 10m
- **Started:** 2026-04-02T17:00:00Z
- **Completed:** 2026-04-02T18:10:19Z
- **Tasks:** 2
- **Files modified:** 8

## Accomplishments
- Added a new top-level `textnorm` package with an immutable fluent `Pipeline` API.
- Implemented Unicode normalization, accent removal, whitespace cleanup, and UTF-8/NUL sanitization stages.
- Added tests that verify chaining immutability, stage ordering, and normalization output.
- Verified the full repository test suite passes with the new package in place.

## Task Commits

1. **Task 1: Scaffold the immutable pipeline core** - `9360aa8` (`feat`)
2. **Task 2: Add Unicode and whitespace cleanup stages** - `6afa81b` (`feat`)

## Files Created/Modified
- `textnorm/doc.go` - package documentation for the new normalization pipeline.
- `textnorm/pipeline.go` - immutable pipeline type and ordered execution.
- `textnorm/stages.go` - exported stage contract.
- `textnorm/unicode.go` - Unicode normalization and accent removal stages.
- `textnorm/whitespace.go` - whitespace and UTF-8 sanitization stages.
- `textnorm/pipeline_test.go` - immutability and ordering tests.
- `textnorm/unicode_test.go` - Unicode normalization tests.
- `textnorm/whitespace_test.go` - whitespace and sanitization tests.

## Decisions Made
- The new package is canonical for normalization work; `stringutil` remains a separate helper package.
- Unicode work is built on `golang.org/x/text` instead of custom rune hacks.
- Phase 1 stays intentionally small: no token pipeline, transliteration, or HTML parsing.

## Deviations from Plan

None - plan executed as specified.

## Issues Encountered

None.

## User Setup Required

None - no external service configuration required.

## Next Phase Readiness
- Phase 1 is complete and stable.
- Phase 2 can now add token pipelines and presets on top of the new core.

---
*Phase: 01-core-pipeline-and-unicode-safety*
*Completed: 2026-04-02*
