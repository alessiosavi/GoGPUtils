# STATE

## Project Reference

See: `.planning/PROJECT.md` (updated 2026-04-02)

**Core value:** Turn messy text into consistent, search-friendly output through a fluent, explicit pipeline that is easy to reuse and hard to misuse.
**Current focus:** Phase 1 planning

## Current Snapshot

- New package direction: dedicated text-normalization pipeline, not an extension of `stringutil`
- Primary use cases: search/DB canonicalization and scraping cleanup
- Stack baseline: Go 1.25.0 + `golang.org/x/text` v0.35.0
- Workflow posture: YOLO, parallel execution, research enabled, plan check enabled, verifier enabled

## Open Questions

- Final package name (`textnorm` is the current working name)
- Whether width folding should be a core stage or only a preset option
- Whether streaming adapters should remain deferred until usage proves the need
