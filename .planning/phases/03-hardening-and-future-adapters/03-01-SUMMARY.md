---
phase: 03-hardening-and-future-adapters
plan: 01
subsystem: testing
tags: [textnorm, benchmarks, fuzzing, docs, hardening]
requires:
  - phase: 02-token-pipelines-and-presets
    provides: token pipelines, presets, and Unicode-safe normalization stages
provides:
  - benchmark coverage for core pipelines and presets
  - fuzz coverage for malformed input and idempotence-sensitive behavior
  - package docs that keep streaming adapters explicitly deferred
affects: [future optimization, future adapters]
tech-stack:
  added: [go test -bench, go fuzz]
  patterns: [benchmark-first hardening, fuzz-driven stability checks, explicit adapter boundaries]
key-files:
  created: [textnorm/bench_test.go, textnorm/fuzz_test.go]
  modified: [textnorm/doc.go, textnorm/presets.go]
key-decisions:
  - "Measure the current implementation before any optimization work."
  - "Use fuzzing to enforce panic-free, idempotent behavior on exported APIs."
  - "Keep streaming adapters explicitly deferred in the package docs."
patterns-established:
  - "Pattern 1: benchmarks live beside the package they measure"
  - "Pattern 2: fuzzers target public constructors and presets only"
  - "Pattern 3: docs describe the future adapter boundary plainly"
requirements-completed: [future-only]
duration: 50m
completed: 2026-04-02
---

# Phase 3: Hardening and Future Adapters Summary

**Benchmarks and fuzzing now measure and stabilize `textnorm` before any optimization work**

## Performance
- **Duration:** 50m
- **Tasks:** 3
- **Files modified:** 4

## Accomplishments
- Added benchmark coverage for the core pipeline and the main presets using representative corpora.
- Added fuzz coverage for malformed input, repeated normalization, and preset stability.
- Updated the package docs to explain benchmarks, fuzzing, and the deferred adapter boundary.

## Task Commits
1. **Task 1: Add benchmark coverage for core and preset pipelines** - `18169e6`
2. **Task 2: Add fuzz coverage for malformed input and stability** - `9bccfc2`
3. **Task 3: Refresh package docs for hardening scope and adapter boundary** - `e297038`

## Decisions Made
- The canonical and search presets now sanitize malformed input before normalization to keep fuzzed inputs stable.
- Benchmark and fuzz entry points stay inside `textnorm` so the package can harden itself without external scaffolding.
- Streaming adapters remain deferred until there is real usage evidence.

## Deviations from Plan

### Auto-fixed Issues

**1. Malformed-input drift in `CanonicalPreset()` and `SearchPreset()`**
- **Found during:** Task 2 (fuzz coverage)
- **Issue:** Fuzzing found that malformed UTF-8 caused repeated runs of `CanonicalPreset()` to drift.
- **Fix:** Sanitized UTF-8 at the start of `SearchPreset()` and `CanonicalPreset()` so malformed input stabilizes before normalization.
- **Files modified:** `textnorm/presets.go`
- **Verification:** Explicit fuzz smoke runs for all four fuzz targets passed after the fix.
- **Committed in:** `9bccfc2`

**Total deviations:** 1 auto-fixed (1 stability bug)
**Impact on plan:** Correctness fix aligned with hardening goals; no scope creep.

## Issues Encountered

- The combined fuzz selector in `go test` does not accept multiple fuzz targets at once, so the fuzz smoke was run as explicit per-target commands.

## Next Phase Readiness
- Phase 3 is complete.
- The package is now benchmarked, fuzz-hardened, and documented with a clear adapter boundary.

---
*Phase: 03-hardening-and-future-adapters*
*Completed: 2026-04-02*
