# STATE

## Project Reference

See: `.planning/PROJECT.md` (updated 2026-04-02)

**Core value:** Turn messy text into consistent, search-friendly output through a fluent, explicit pipeline that is easy to reuse and hard to misuse.
**Current focus:** v1 complete; maintain the package and defer adapters until there is usage evidence

## Current Snapshot

- New package direction: dedicated text-normalization pipeline, not an extension of `stringutil`
- Primary use cases: search/DB canonicalization and scraping cleanup
- Stack baseline: Go 1.25.0 + `golang.org/x/text` v0.35.0
- Workflow posture: YOLO, parallel execution, research enabled, plan check enabled, verifier enabled
- Phase 1 plan completed: `.planning/phases/01-core-pipeline-and-unicode-safety/01-01-SUMMARY.md`
- Phase 1 commits: `9360aa8`, `6afa81b`, `dfdfb04`
- Phase 2 plan summaries: `.planning/phases/02-token-pipelines-and-presets/02-01-SUMMARY.md`, `.planning/phases/02-token-pipelines-and-presets/02-02-SUMMARY.md`
- Phase 2 commits: `13b9799`, `2fa005b`, `99e3194`, `efa626a`, `f49676a`, `fd3c35b`
- Phase 3 plan summary: `.planning/phases/03-hardening-and-future-adapters/03-01-SUMMARY.md`
- Phase 3 commits: `18169e6`, `9bccfc2`, `e297038`
- v1 requirements are validated in `.planning/REQUIREMENTS.md`

## Open Questions

- Final package name (`textnorm` is the current working name)
- Whether width folding should be a core stage or only a preset option
- Whether streaming adapters should remain deferred until usage proves the need
