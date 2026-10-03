# `cache` package — design spec

- **Date:** 2026-10-03
- **Status:** v2 — approved design (brainstorming §1–§3) + Codex spec review
  rounds 1–2 (all findings accepted, see §9); pending owner review
- **Package:** `github.com/alessiosavi/GoGPUtils/cache`
- **Go:** module minimum (`go 1.26.4`); standard library only

## 1. Goal

A general-purpose, concurrency-safe, generic in-memory cache that covers three
workloads with one type:

1. **Memoizing expensive calls** — read-through `GetOrLoad` with per-key load
   deduplication.
2. **Bounded key/value store** — explicit `Set`/`Get`/`Delete` with an
   entry-count bound and SIEVE eviction.
3. **TTL store** — per-entry expiration, lazy removal, optional background
   janitor.

### Non-goals (v1)

- Cost/byte-weighted capacity (entry count only).
- Pluggable eviction policies; W-TinyLFU admission.
- Negative caching, stale-while-revalidate, refresh-ahead.
- Persistence, serialization, distribution.
- Per-call TTL for values stored by `GetOrLoad` (they use the default TTL).
- Load cancellation, loader worker pools, cycle detection between loaders.

## 2. Repository conventions to follow

- Package doc in `cache/doc.go` with `// Example:` blocks, like other packages.
- Sentinel errors in a `var` block, prefixed `cache:`.
- `(V, bool)` for lookups; errors only for invalid input or load failures.
- No panics on misuse that can be reported as an error; no package-level mutable
  state; no imports outside the standard library.
- Zero value of `Cache` is **not** usable; use `New` (same rule as
  `collection.Set`). A `*Cache` must not be dereferenced and copied.
- Lessons from the 2026-10-03 audit: zero removed slots/pointers so evicted
  values are collectable; no unseeded/global randomness; no unbounded internal
  structures; every documented behavior claim is backed by a test.

## 3. Public API

```go
package cache

// NoExpiration, passed to SetWithTTL, stores an entry that never expires.
// TTL returns it for entries without an expiration.
const NoExpiration time.Duration = -1

// MaxShards is the largest accepted Config.Shards value.
const MaxShards = 1 << 16

type Config[K comparable, V any] struct {
	// MaxEntries bounds the number of stored entries. 0 means unbounded.
	MaxEntries int
	// TTL is the default time-to-live. 0 means entries never expire.
	TTL time.Duration
	// Shards is the number of independent shards. 0 selects automatically.
	// Other values are rounded up to a power of two (see §4.1).
	Shards int
	// JanitorInterval > 0 starts a goroutine that removes expired entries at
	// that interval. Call Close to stop it. 0 disables the janitor.
	JanitorInterval time.Duration
	// OnEvict, if set, is called for every entry removed from the cache,
	// outside any internal lock, in the goroutine that caused the removal.
	// It must not call Close (see §4.5).
	OnEvict func(key K, value V, reason EvictionReason)
}

type EvictionReason int

const (
	ReasonEvicted  EvictionReason = iota + 1 // removed to make room (capacity)
	ReasonExpired                            // TTL elapsed
	ReasonDeleted                            // Delete or Clear
	ReasonReplaced                           // Set/SetWithTTL overwrote the key
)

func (r EvictionReason) String() string // "evicted", "expired", "deleted", "replaced"

var (
	ErrInvalidConfig = errors.New("cache: invalid config")
	ErrNilLoader     = errors.New("cache: nil loader")
	ErrNilContext    = errors.New("cache: nil context")
	ErrLoaderGoexit  = errors.New("cache: loader called runtime.Goexit")
)

// PanicError is returned by GetOrLoad when the loader panics.
type PanicError struct {
	Value any    // the value passed to panic
	Stack []byte // stack trace captured at recovery
}

func (e *PanicError) Error() string // "cache: loader panicked: <Value>"
func (e *PanicError) Unwrap() error  // Value if it is an error, else nil

type Stats struct {
	Hits, Misses, Loads, LoadErrors, Evictions, Expirations uint64
}

func (s Stats) HitRatio() float64 // Hits/(Hits+Misses); 0 when both are 0

func New[K comparable, V any](cfg Config[K, V]) (*Cache[K, V], error)

func (c *Cache[K, V]) Get(key K) (V, bool)
func (c *Cache[K, V]) Peek(key K) (V, bool)
func (c *Cache[K, V]) TTL(key K) (time.Duration, bool)
func (c *Cache[K, V]) Set(key K, value V)
func (c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration)
func (c *Cache[K, V]) Delete(key K) bool
func (c *Cache[K, V]) GetOrLoad(ctx context.Context, key K,
	load func(context.Context) (V, error)) (V, error)
func (c *Cache[K, V]) DeleteExpired() int
func (c *Cache[K, V]) Clear()
func (c *Cache[K, V]) Len() int
func (c *Cache[K, V]) All() iter.Seq2[K, V]
func (c *Cache[K, V]) Stats() Stats
func (c *Cache[K, V]) Close() error
```

### 3.1 Keys (non-reflexive keys)

Keys must be reflexive (`k == k`). Keys that are not — NaN floats, or arrays,
structs, or interfaces containing NaN — cannot be found or deleted in a Go map,
so the cache **does not store them**: `Set`/`SetWithTTL` are no-ops, `Get`/
`Peek`/`TTL` report a miss, `Delete` returns `false`, and `GetOrLoad` performs
an ordinary load (§4.4: same validation, detached execution, panic/`Goexit`
handling, cancellation, error semantics, and `Misses`/`Loads`/`LoadErrors`
counters) through an **unregistered** call — only deduplication and storage are
bypassed.

The guard is `key != key`, evaluated before hashing or locking (one equality
operation; proportional to the size of composite keys). If that self-comparison
panics (an interface key whose dynamic value is not comparable), the panic
propagates; the cache does not promise exact parity with Go map panics (e.g. for
`any([2]any{math.NaN(), []int{1}})` the comparison stops at the NaN and the key is
treated as non-reflexive).

### 3.2 Method contracts

All methods are safe for concurrent use. "Expired" means `now >= expiresAt`,
evaluated with one clock sample per decision.

| Method | Contract |
|---|---|
| `New` | Returns `ErrInvalidConfig` (wrapped via `fmt.Errorf("%w: <field> ...")`) if `MaxEntries < 0`, `TTL < 0`, `Shards < 0`, `Shards > MaxShards`, or `JanitorInterval < 0`. Starts no goroutine unless `JanitorInterval > 0`. |
| `Get` | See the state table (§3.3). |
| `Peek` | Never marks visited, never updates stats, never removes. |
| `TTL` | Live entry with expiry: `(remaining > 0, true)`. Live entry without expiry: `(NoExpiration, true)`. Missing or expired: `(0, false)`. Never removes, no stats, no visited mark. |
| `Set` | `SetWithTTL(key, value, 0)`. |
| `SetWithTTL` | `ttl == 0` → default TTL; `ttl < 0` → never expires; `ttl > 0` → expires after `ttl` (saturating, §4.3). Invalidates any in-flight load for the key (§4.4). |
| `Delete` | Returns `true` iff a live entry was removed. Invalidates any in-flight load for the key. |
| `GetOrLoad` | §4.4. |
| `DeleteExpired` | Removes every expired entry, returns the count. Processes one shard at a time (not atomic across shards). |
| `Clear` | Removes every entry with `ReasonDeleted` (including expired ones) and invalidates every in-flight load, one shard at a time (not atomic across shards). Stats are kept. |
| `Len` | Sum of per-shard stored-entry counts (may include expired entries not yet removed). Each shard count is `≤` its capacity, so the sum is `≤ MaxEntries` when bounded. Not an atomic snapshot across shards. |
| `All` | Each invocation of the returned iterator takes a fresh pass. Per shard: copy entries live at copy time under the read lock, release, yield. Entries may expire between copy and yield. Order unspecified; no visited marks or stats; the cache may be mutated during iteration; stops when `yield` returns false. |
| `Stats` | Sum of per-shard atomic counters; not an atomic snapshot. |
| `Close` | Stops the janitor (if any), waits for it to exit, and stops the cleanup registration (§4.6). Idempotent and safe to call concurrently; returns `nil`. Does not cancel in-flight loads. The cache remains fully usable afterwards. |

### 3.3 Method × entry-state table

`R` = reason reported to `OnEvict`; counters are incremented only as listed.

| Method | missing | live | expired |
|---|---|---|---|
| `Get` | miss; `Misses++` | hit; mark visited; `Hits++` | remove `R=Expired`, `Expirations++`; miss; `Misses++` (if, after re-locking, a live replacement is found → treated as a hit) |
| `Peek` | `(zero,false)` | `(v,true)` | `(zero,false)`, entry kept |
| `TTL` | `(0,false)` | `(rem,true)` / `(NoExpiration,true)` | `(0,false)`, entry kept |
| `SetWithTTL` | insert unvisited (may evict one entry, `R=Evicted`/`Expired`) | replace in place, mark visited, `R=Replaced` (old value) | replace in place, mark visited, `R=Expired` (old value), `Expirations++` |
| `Delete` | `false` | remove `R=Deleted`, `true` | remove `R=Expired`, `Expirations++`, `false` |
| `DeleteExpired` | — | kept | remove `R=Expired`, `Expirations++` |
| `Clear` | — | remove `R=Deleted` | remove `R=Deleted` (no `Expirations++`) |
| eviction (§4.2) | — | unvisited → remove `R=Evicted`, `Evictions++` | remove `R=Expired`, `Expirations++` |

Counter definitions: `Evictions` counts only `ReasonEvicted`; `Expirations`
counts only `ReasonExpired`; `ReasonDeleted`/`ReasonReplaced` increment neither.
`Hits`/`Misses` come from `Get` and `GetOrLoad` (§4.4); `Loads` counts loader
executions started; `LoadErrors` counts loads that returned an error, panicked,
or called `runtime.Goexit` (caller cancellation is not a load error).

## 4. Internals

### 4.1 Sharding

- `n` shards, `n` a power of two. Shard index =
  `maphash.Comparable(seed, key) & (n-1)` with a per-cache `maphash.MakeSeed()`.
  When `n == 1`, hashing is skipped.
- Auto (`Shards == 0`):
  - `base = nextPow2(4 * runtime.GOMAXPROCS(0))`, capped at `MaxShards`.
  - Unbounded: `n = base`.
  - Bounded: `n = min(base, prevPow2(max(1, MaxEntries/1024)))` — every shard
    holds ≥ 1024 entries, so `MaxEntries < 2048` gives a single shard (exact
    global SIEVE).
- Explicit `1 ≤ Shards ≤ MaxShards`: `n = nextPow2(Shards)` (≤ `MaxShards`, no
  overflow); if bounded and `n > MaxEntries`, `n = prevPow2(MaxEntries)`.
- Capacity split: `MaxEntries / n` per shard, the first `MaxEntries % n` shards
  get one more. Sum equals `MaxEntries` exactly.

### 4.2 Shard and SIEVE

```go
type entry[K comparable, V any] struct {
	key        K
	value      V
	expiresAt  int64 // monotonic ns since cache epoch; 0 = never
	visited    atomic.Bool
	prev, next *entry[K, V] // prev = toward tail (older), next = toward head (newer)
}

type shard[K comparable, V any] struct {
	mu       sync.RWMutex
	items    map[K]*entry[K, V]
	head     *entry[K, V] // newest
	tail     *entry[K, V] // oldest
	hand     *entry[K, V] // SIEVE hand; nil = start at tail
	capacity int          // 0 = unbounded
	inflight map[K]*call[V]
	stats    shardStats   // atomic counters
}
```

- **Lookup hit path:** `RLock`, map lookup, expiry check, `visited.Store(true)`,
  copy value, `RUnlock`. No write lock on hits.
- **Expired on lookup:** release the read lock, take the write lock, **re-fetch**
  the key: missing → miss; live (replaced meanwhile) → hit; still expired →
  remove (`ReasonExpired`) → miss.
- **Insert:** at `head`, `visited = false`.
- **Evict (one victim, write lock held):** `e := hand`, or `tail` if `hand` is
  nil. Loop: if `e` is expired → victim (`ReasonExpired`); else if `e.visited`
  → clear it, `e = e.next` (toward head), and if nil wrap to `tail`; else →
  victim (`ReasonEvicted`). For the victim: save `p := victim.next` **before**
  unlinking, remove it, set `hand = p` (nil → next eviction starts at `tail`).
  At most `m+1` inspections for `m` entries.
- **Removal (any path):** if `hand == e`, set `hand = e.next` (saved before
  unlinking); unlink; delete from the map; zero `value`, `prev`, `next` so the
  value is collectable. Invariant (checked in tests): map and list contain the
  same entries; `hand` is nil or a linked entry.

### 4.3 Time

- `now()` = `int64(time.Since(epoch))` where `epoch` is captured in `New`;
  monotonic, so wall-clock jumps never expire or revive entries.
- `expiresAt = saturatingAdd(now, ttl)` (clamps at `math.MaxInt64`); a computed
  value of `0` is bumped to `1` so `0` keeps meaning "never".
- The clock is an unexported, concurrency-safe field installed **before** any
  goroutine starts (internal constructor used by tests), so tests drive time
  deterministically.
- Entries with `expiresAt == 0` skip the clock read on lookup.

### 4.4 `GetOrLoad`

Validation, in order, before any lookup effect: `ctx == nil` → `ErrNilContext`;
`load == nil` → `ErrNilLoader`; `ctx.Err() != nil` → `ctx.Err()`.
Non-reflexive key → §3.1.

1. Fast path (read lock): live hit → mark visited, `Hits++`, return.
2. Write lock; re-fetch: live → mark visited, `Hits++`, unlock, return.
   Expired → remove (`ReasonExpired`, `Expirations++`). Then `Misses++`.
   If `inflight[key]` exists, join it; otherwise create
   `call{done: make(chan struct{})}`, set `inflight[key] = call`, `Loads++`,
   and start `go c.runLoad(context.WithoutCancel(ctx), key, call, load)`.
   Unlock. Dispatch any collected `OnEvict` callbacks, then wait.
3. Wait: `select { case <-call.done: return call.val, call.err; case <-ctx.Done(): return zero, ctx.Err() }`.
   If both are ready, either result may be returned. A caller giving up never
   cancels the load or affects other waiters.
4. **Invalidation:** `SetWithTTL`, `Delete`, and `Clear` (under the shard lock)
   set `call.invalidated = true` **and remove the call from `inflight`**. Existing
   waiters still receive the old call's result; callers arriving afterwards start
   a new load (new generation).
5. **`runLoad`** guarantees exactly one completion via a deferred `complete`:
   - The loader runs inside a nested function `invoke() (v V, err error)` with
     its own deferred `recover` that converts a panic into `*PanicError`
     (captures `debug.Stack()`), so `runLoad` continues normally.
   - If the loader calls `runtime.Goexit`, `invoke` never returns normally;
     `runLoad`'s deferred `complete` detects this (a `returned` flag) and
     completes with `ErrLoaderGoexit`.
   - `complete` (write lock): if `inflight[key] == call`, delete it (an
     invalidated older call never removes a newer one); on success and
     `!call.invalidated`, store the value with the default TTL via the `Set`
     insert path (capacity eviction allowed; the key cannot hold a live value
     because any write would have invalidated the call); on error
     `LoadErrors++` and the value is discarded (`call.val` = zero). Unlock;
     publish `call.val`/`call.err`; `close(call.done)`; then dispatch collected
     `OnEvict` callbacks (in the load goroutine).
6. Errors are never cached.

### 4.5 Callbacks

- Removals collect `(key, value, reason)` while holding a shard lock; `OnEvict`
  runs after unlocking, in the goroutine that triggered the removal (caller,
  load goroutine, or janitor). No shard lock is held during callbacks.
- Callbacks from different goroutines may run concurrently; there is no global
  ordering. Within one operation they run in removal order.
- Callbacks may call any cache method **except `Close`** (calling `Close` from a
  callback running on the janitor would wait for itself). Documented.
- A panicking callback panics its goroutine and skips the rest of that batch;
  in the janitor or a load goroutine this terminates the process (documented,
  not recovered).

### 4.6 Janitor and lifecycle

- Public `Cache` is a small wrapper holding `*inner`; all state lives in a
  separately allocated `inner`. The janitor goroutine references only `inner`.
- `JanitorInterval > 0`: start a goroutine with a `time.Ticker` calling
  `inner.deleteExpired`; it exits when `inner.stop` is closed.
- `inner.requestStop()` closes `stop` exactly once (`stopOnce sync.Once`); it
  only signals and never waits. It is the only code that closes `stop`.
- `Close` (every call, concurrently safe): `inner.requestStop()`; if a janitor
  was started, wait on its `done` channel (already-closed → returns at once);
  `cleanup.Stop()` (only if a cleanup was registered); `runtime.KeepAlive(c)` so
  the wrapper stays reachable until `Stop` returns. Without a janitor, `Close`
  does nothing beyond `requestStop` and returns `nil`.
- **Best-effort safety net:** registered only when a janitor is started:
  `runtime.AddCleanup(outer, func(in *inner) { in.requestStop() }, inner)`
  (the function must not capture `outer`); it signals the janitor without
  waiting if the wrapper becomes unreachable. Cleanup timing is not guaranteed, and it never
  runs if `OnEvict` or stored keys/values reference the wrapper. `Close` is the
  only deterministic shutdown; docs say so.

### 4.7 Documentation notes (doc.go)

- Values are stored as-is (shallow); mutable contents need caller synchronization.
- `MaxEntries` bounds stored entries, not concurrent loads; loaders run on a
  detached context, so give them their own timeouts. `Close` does not cancel
  loads.
- A loader must not call `GetOrLoad` for its own key (or form cycles with other
  loads) — it would wait for itself until the caller's context ends.
- Do not copy `Cache` values; use the `*Cache` returned by `New`.

## 5. File layout

```
cache/
  doc.go          package documentation and examples
  cache.go        Config, New, Cache wrapper, errors, EvictionReason, PanicError
  shard.go        shard, entry, SIEVE list operations, eviction, removal
  load.go         GetOrLoad, call, runLoad, invalidation
  janitor.go      janitor goroutine, Close, cleanup registration
  stats.go        Stats, shardStats
  *_test.go       see §6
docs/packages/cache.md, README.md (package table + section), docs/architecture.md
```

## 6. Testing

TTL tests use the injected clock; timer/goroutine tests use `testing/synctest`
where applicable (runtime cleanups run outside synctest bubbles, so GC-cleanup
tests stay separate). Unexported test hooks may pause execution at named points
(e.g. between read-unlock and write-lock, inside `runLoad`) to force
interleavings deterministically.

- **SIEVE fidelity:** single-shard cache vs a reference SIEVE (test-only, from
  the NSDI'24 paper) on scripted and seeded pseudo-random traces (`math/rand/v2`
  PCG); eviction order must match exactly. Cases: all entries visited; arbitrary
  starting hand; wrap-around; hand pointing at an entry removed via `Delete`,
  `Clear`, expiry, and eviction; replacement does not move the entry;
  expired-and-visited entry is the victim; singleton/head/tail/middle/empty.
  After every operation assert map/list agreement and hand validity.
- **Model-based:** random op sequences (Get/Peek/TTL/Set/SetWithTTL/Delete/
  DeleteExpired/Clear/advance-clock) on a single-shard bounded cache, with the
  reference SIEVE as the eviction oracle and a map model for TTL. Invariants:
  `Len() ≤ MaxEntries`; no expired value returned; every removal reported once
  with the reason/counters of §3.3; `All` yields exactly the live set. Also as
  `FuzzCache` with capped op counts/allocations, including capacity 1,
  deadline equality (`now == expiresAt`), extreme durations
  (`math.MaxInt64` TTL), and shard-rounding boundaries.
- **Keys:** NaN float32/float64, complex with NaN part, struct/array/interface
  containing NaN → not stored, per §3.1, and `GetOrLoad` on such a key behaves as
  an ordinary non-deduplicated load (cancellation, panic, error semantics,
  counters); interface key with non-comparable dynamic value panics; mixed key
  `any([2]any{math.NaN(), []int{1}})` is treated as non-reflexive.
- **TTL:** default vs override vs `NoExpiration`; `TTL(key)` values; saturating
  expiry; `DeleteExpired` count; eviction treats an expired entry as the victim
  even if visited (no global "expired first" claim).
- **GetOrLoad:** hold the loader until 100 callers have missed, then assert
  `Misses=100`, `Loads=1`, all get the value; cancelling one waiter does not
  cancel the load or others and does not count a `LoadError`; all waiters
  cancelled → load still completes and stores; errors not cached (next call
  reloads); panic → `*PanicError` with `Unwrap`; `runtime.Goexit` →
  `ErrLoaderGoexit`; `Set`/`Delete`/`Clear` during a load → no stale write-back,
  waiters get the old value, a caller arriving after invalidation runs a new
  loader, and the old completion does not remove the new call; validation order
  (nil ctx, nil loader, done ctx) with no lookup effects; both hit paths mark
  visited; expired entry replaced during lock re-acquisition → hit.
- **Callbacks & stats:** each reason and counter per §3.3; re-entrant callbacks
  (Get/Set/GetOrLoad from capacity-eviction and expiration callbacks) do not
  deadlock; callback completion synchronized explicitly in tests (`done` closes
  before callbacks run); concurrent callbacks recorded with a synchronized
  recorder; `Peek`/`TTL`/`All` change neither stats nor SIEVE order; `HitRatio`
  zero case.
- **Config:** each invalid field → `ErrInvalidConfig`; `Shards = MaxShards` and
  `MaxShards+1`; auto/explicit rounding; capacity split sums to `MaxEntries`.
- **Lifecycle:** janitor removes expired entries; `Close` stops the janitor
  (verified via its own done signal, not global goroutine counts), is
  idempotent and safe concurrently, works without a janitor, and the cache stays
  usable after it. Correctness is gated on deterministic tests of
  `requestStop`/`Close`/janitor exit. The actual GC-cleanup path is an
  observational test: drop the last reference, poll `runtime.GC` up to 10 s,
  and if the janitor has not exited, `t.Skip("cleanup not observed (timing not
  guaranteed)")` — inconclusive, never a failure; skipped under `-short`.
- **Race stress:** concurrent Get/Set/Delete/GetOrLoad/Clear/All/DeleteExpired
  with janitor, multiple shards, under `-race`.
- **Examples:** `example_test.go` with runnable `Example*` functions.

### Benchmarks (not gates)

Parallel Get-hit; Set; 90/10 and 50/50 read/write; GetOrLoad hit; shards 1 vs
auto; baseline `map + sync.RWMutex`; SIEVE vs simple LRU hit ratio on a Zipf
trace (reported, not asserted).

## 7. Acceptance criteria

Each item is recorded with its actual command output in the PR (green CI alone
does not prove 3–5, because CI lint is non-blocking and coverage is not
enforced there).

1. `gofmt -l .` empty; `go build ./...`, `go vet ./...` clean.
2. `go test -race -count=1 ./...` passes; `go test -run '^$' -fuzz FuzzCache -fuzztime 60s ./cache/` finds nothing.
3. `staticcheck ./cache/...` clean; `golangci-lint run ./cache/...` (repo config) reports no findings.
4. `go test -cover ./cache/` ≥ 90%, and no exported function below 80% (`go tool cover -func`).
5. `go list -deps -f '{{if and .DepOnly (not .Standard)}}{{.ImportPath}}{{end}}' ./cache/` prints nothing.
6. Docs: `cache/doc.go` (incl. §4.7 notes), runnable examples,
   `docs/packages/cache.md`, README and architecture entries.
7. Delivered via feature branch `feat/cache` → PR → bot review → green CI.

## 8. Open questions

None. Implementation-time questions go back to the owner before deviating from
this spec.

## 9. Review log

- **Codex round 1 (2026-10-03): changes needed, 14 findings, all accepted.**
  F1 `Close` from callback → prohibited (§4.5). F2 non-reflexive keys → not
  stored (§3.1). F3 cleanup → best-effort, separate allocation, `cleanup.Stop`
  in `Close` (§4.6). F4 post-invalidation joiners → invalidation detaches the
  call; completion deletes only its own registry entry (§4.4). F5 recover scope
  and `Goexit` → nested `invoke`, deferred `complete`, `ErrLoaderGoexit`
  (§4.4). F6 validation order, visited on both hit paths, re-fetch after
  re-lock, either-result race (§3.3, §4.4). F7 per-shard aggregate semantics
  (§3.2). F8 saturating expiry, `MaxShards` (§3, §4.1, §4.3). F9 hand saved
  before unlink, removal fix-up (§4.2). F10 counters by reason, non-removing
  `TTL`/`Peek`, state table (§3.3). F11 callback dispatch in `GetOrLoad`,
  concurrency, panic semantics (§4.4, §4.5). F12 clock before goroutines, janitor
  exit via its own signal, synctest (§4.3, §6). F13 controlled interleavings and
  SIEVE oracle in the model (§6). F14 `DepOnly` dependency check, explicit
  lint/coverage evidence (§7).
- **Codex round 2 (2026-10-03): 11/14 resolved; F2/F3/F12 partial + N1–N3,
  all accepted as prescribed.** N1 non-reflexive `GetOrLoad` uses the ordinary
  load machinery via an unregistered call (§3.1). N2 guard before hashing,
  comparison panics propagate, no exact map-parity claim (§3.1, §6). N3
  signal-only `requestStop` shared by `Close` and cleanup, `KeepAlive` after
  `cleanup.Stop` (§4.6). F12 residual: GC-cleanup observation is inconclusive on
  timeout, never a failure (§6).
