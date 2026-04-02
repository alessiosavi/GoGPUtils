---
phase: 02-token-pipelines-and-presets
plan: 01
subsystem: library
tags: [textnorm, tokens, unicode, cases, width, x-text]
requires:
  - phase: 01-core-pipeline-and-unicode-safety
    provides: immutable pipeline core and Unicode-safe cleanup stages
provides:
  - token pipeline contracts and whitespace-first splitting
  - Unicode case-folding, lowercasing, rune transforms, and width folding stages
  - regression tests for token and stage behavior
affects: [Phase 2 presets, future search/canonical pipelines]
tech-stack:
  added: [golang.org/x/text/cases, golang.org/x/text/runes, golang.org/x/text/width]
  patterns: [immutable token pipeline, whitespace-first tokenization, opt-in width folding, explicit rune transforms]
key-files:
  created: [textnorm/tokens.go, textnorm/case.go, textnorm/width.go, textnorm/tokens_test.go, textnorm/case_test.go, textnorm/width_test.go, textnorm/regression_test.go]
  modified: [textnorm/unicode.go]
key-decisions:
  - "Token pipelines are split from the string pipeline via whitespace boundaries only."
  - "Width folding is explicit and opt-in; default pipeline behavior stays unchanged."
  - "Case and rune transforms are exposed as plain pipeline stages, not presets."
patterns-established:
  - "Pattern 1: token pipelines mirror the immutable string-pipeline style"
  - "Pattern 2: `strings.Fields` defines token boundaries"
  - "Pattern 3: Unicode stage combinations remain explicit and deterministic"
requirements-completed: [UNIC-03, UNIC-04, TOKN-01, TOKN-02, TOKN-03, TOKN-04, PRES-04]
duration: 55m
completed: 2026-04-02
---

# Phase 2: Token Pipelines and Presets - Plan 01 Summary

**Immutable token pipelines plus Unicode-aware case, rune, and width stages**

## Performance
- **Duration:** 55m
- **Tasks:** 3
- **Files modified:** 7

## Accomplishments
- Added a token pipeline API that splits on whitespace and can map, filter, and join tokens immutably.
- Added case-folding, lowercasing, rune transform, and explicit width-folding stages.
- Added regression tests proving token boundaries, case behavior, and width opt-in behavior.

## Task Commits
1. **Task 1: Add token pipeline contracts and whitespace-first splitting** - `13b9799`
2. **Task 2: Add unicode case and width stages** - `2fa005b`
3. **Task 3: Add mixed normalization regression test** - `99e3194`

## Decisions Made
- `SplitTokens()` uses `strings.Fields` on the source pipeline output.
- `FoldWidth()` is opt-in only and does not change default normalization.
- `NormalizeUnicode()` now avoids width folding so presets can control it explicitly.

## Deviations from Plan

None - plan executed as specified.

## Next Phase Readiness
- Wave 2 presets can now compose the token and normalization primitives.
- Default pipeline behavior remains stable while width folding is isolated behind an option.
