# `cache` package — development handoff (2026-10-03)

Read this first, then the spec and the plan. It captures the state, decisions, workflow,
and gotchas from the session that produced them, so development can start in a fresh or
compacted context without re-deriving anything.

## Where things stand

| Item | State |
|---|---|
| Audit fixes | Merged to `master` (squash): #37 hygiene/CI/lint/deps, #36 fileutil integrity, #38 core High fixes, #39 AWS correctness. `master` = `65f6a41`, CI green. |
| Spec | `docs/superpowers/specs/2026-10-03-cache-design.md` — **READY** (Codex spec review rounds 1–4; §9 = review log). Do not reopen settled decisions. |
| Plan | `docs/superpowers/plans/2026-10-03-cache.md` — 9 tasks, full code, TDD steps. **READY** after Codex executable plan reviews r1–r2: the committed plan was replayed task-by-task in a scratch copy — every red step failed as expected, every gate passed (race ×3, Go 1.26.4, 60 s fuzz ≈11M execs, staticcheck, golangci-lint 0 issues, coverage 99.8 %, all exported funcs 100 %); Claude independently re-ran race/vet/gofmt/coverage on the replay and confirmed the locked-hit mutant is caught. |
| Branch / worktree | `feat/cache` at `/opt/SP/Workspace/Go/GoGPUtils/.worktrees/feat-cache` (rebased on `65f6a41`). Contains only the spec, plan, and this file. **No implementation yet.** Not pushed. |
| Primary checkout | `/opt/SP/Workspace/Go/GoGPUtils` is still on the old `master` (`5bf64f2`) with the owner's uncommitted `go.mod`/`go.sum` (identical to what #37 merged). Owner discards them before pulling. Never touch them. |

## Owner decisions (fixed)

General-purpose cache (bound + TTL + loader); SIEVE eviction; entry-count capacity;
lazy expiry + `DeleteExpired` + opt-in janitor with `Close`; per-call deduplicated
`GetOrLoad` on a detached context; extras: `Stats`, `OnEvict`, `All()` (`iter.Seq2`),
`Peek`, `TTL(key)`; sharded SIEVE; `Config[K, V]` struct (not functional options);
squash-merge PRs in order after bot approval and green CI.

## Workflow (owner's global CLAUDE.md — strict)

- **Roles:** Claude orchestrates and reviews; **Codex implements**. Each plan task is one
  Codex handoff (implementation) followed by Claude's review of the actual diff and an
  independent re-run of the task's checks, before the next task starts.
- **Codex invocation (implementation):**
  `codex exec -p worker -C /opt/SP/Workspace/Go/GoGPUtils/.worktrees/feat-cache -o <out.md> - < <prompt.md>`
  (finite stdin; `-p worker` = workspace-write + network, writes limited to `-C`, /tmp, Go caches).
- **Codex invocation (review-only):** run with `-C <scratchpad> --skip-git-repo-check`;
  **never** pass `--add-dir <repo>` to a review (it makes the repo writable).
- **Git:** Codex's sandbox cannot write `.git`; Claude does every `git add/commit/push`.
  Conventional Commits; **no AI attribution** (no Co-Authored-By, no "Generated with").
- **Delivery:** push `feat/cache` → PR to `master` → wait for bot `git-code-reviewer-v1`
  (posts "Starting… / Review complete!" comments, then a review) and CI → verify every bot
  finding against the code (fix, or reject with evidence in a thread reply) → **ask the owner
  before merging** → squash-merge → remove the worktree from the primary checkout
  (`cd /opt/SP/Workspace/Go/GoGPUtils && git worktree remove .worktrees/feat-cache && git worktree prune
  && git branch -D feat/cache && git push origin --delete feat/cache`).
- `master` is not branch-protected; still never push to it directly.

## Environment gotchas (learned the hard way)

- `rtk` rewrites shell commands and truncates long `grep`/`ls` output; use `rtk proxy <cmd>`
  for complete output.
- Inside Codex sandboxes, staticcheck/golangci caches are not writable: use
  `STATICCHECK_CACHE=/tmp/cache-staticcheck` and `GOLANGCI_LINT_CACHE=/tmp/gl`.
- `testing/synctest`: create and `Close` caches with janitors **inside** the bubble (their
  channels are bubbled); do not call `t.Run` inside a bubble; `waitFor` in the plan calls
  `synctest.Wait()` and is valid only inside a bubble.
- CI lint is still non-blocking (`continue-on-error`) and CI does not enforce coverage, so the
  PR must quote the Task 9 gate outputs explicitly.
- A monitor/watch that greps `gh` output can miss events; confirm with
  `gh pr view <n> --json statusCheckRollup,reviews` before concluding.

## Acceptance gates (plan Task 9 / spec §7)

`gofmt -l .` empty; `go build ./... && go vet ./...`; `go test -race -count=1 ./...`;
`go test -run '^$' -fuzz FuzzCache -fuzztime 60s ./cache/`; staticcheck + golangci-lint clean
on `./cache/...`; coverage ≥ 90 % and no exported function < 80 %; stdlib-only dependency
check empty; docs updated; PR → bot → green CI.

## Audit backlog (not part of the cache work)

Medium findings from the 2026-10-03 audit still open (candidates for follow-up PRs):
mathutil overflow (IsPrime narrow ints, Average/Sum narrow ints, LCM intermediate overflow,
unsigned ManhattanDistance, Pow float negative exponent, PowFloat MinInt), randutil
IntRange/Int64Range overflow panics, sliceutil.Shuffle fixed seed 1 + global mutex + 32-bit
negative index (keep SeedShuffle working), collection Stack/Queue pointer retention, Set zero
value, BST NaN ordering/unbalanced, stringutil Levenshtein empty-path byte counts,
CommonPrefix/Suffix rune boundaries, Wrap/Pad edge cases, textnorm JoinTokens double execution
and stopword accent guidance, cryptoutil DeriveKey (add PBKDF2, keep legacy), dead hex
helpers, SecureString throughput for long strings (per-byte `crypto/rand.Int`), performance
items (Intersect, quantiles sort once, Histogram binary search, Queue ring buffer), and making
lint blocking after a cleanup sweep. Deferred by decision: AWS submodule, BST balancing,
typed Zip, `iter.Seq` APIs, s3 transfermanager migration.
