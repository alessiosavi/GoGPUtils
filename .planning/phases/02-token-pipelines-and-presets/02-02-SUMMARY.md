---
phase: 02-token-pipelines-and-presets
plan: 02
subsystem: library
tags: [textnorm, presets, search, canonical, db-safe, x-text]
requires:
  - phase: 02-token-pipelines-and-presets
    provides: token and normalization stage API from wave 1
provides:
  - search, canonical, and DB-safe preset constructors
  - preset-focused package documentation
  - regression tests for preset behavior and width folding opt-in
affects: [Phase 3 hardening and future adapters]
tech-stack:
  added: [none]
  patterns: [thin presets, stage composition reuse, opt-in width folding, explicit preset docs]
key-files:
  created: [textnorm/presets.go, textnorm/presets_test.go]
  modified: [textnorm/doc.go]
key-decisions:
  - "Presets are thin wrappers over the existing stage API."
  - "Search preset filters to letters, numbers, and spaces before token rejoin."
  - "Width folding is toggled only through `WithWidthFold()`."
patterns-established:
  - "Pattern 1: presets are composed, not hand-rolled"
  - "Pattern 2: documentation mirrors exported preset behavior"
  - "Pattern 3: DB-safe cleanup starts with explicit UTF-8/NUL sanitization"
requirements-completed: [PRES-01, PRES-02, PRES-03]
duration: 35m
completed: 2026-04-02
---

# Phase 2: Token Pipelines and Presets - Plan 02 Summary

**Thin reusable presets for search, canonical, and DB-safe normalization**

## Performance
- **Duration:** 35m
- **Tasks:** 3
- **Files modified:** 3

## Accomplishments
- Added preset constructors for search, canonical, and DB-safe normalization.
- Documented the package surface and opt-in width-fold behavior.
- Added regression tests for preset output, sanitization, and width-fold opt-in.

## Task Commits
1. **Task 1: Add preset constructors and width-fold option wiring** - `efa626a`
2. **Task 2: Document the preset surface and usage examples** - `f49676a`
3. **Task 3: Add preset regression coverage** - `fd3c35b`

## Decisions Made
- Search preset uses tokenization and rejoin so punctuation stripping stays explicit.
- Canonical preset preserves punctuation and only normalizes case/space.
- DB-safe preset sanitizes invalid UTF-8 and NUL bytes before canonical cleanup.

## Deviations from Plan

None - plan executed as specified.

## Next Phase Readiness
- Phase 2 is complete.
- Phase 3 can now harden the package with benchmarks, fuzzing, and optional adapters.
