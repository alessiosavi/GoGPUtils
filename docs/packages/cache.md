---
title: cache
parent: Packages
nav_order: 9
---

# cache

Generic, concurrency-safe in-memory cache: SIEVE eviction, TTL expiration, and deduplicated loading.
{: .fs-6 .fw-300 }

## Overview

```go
import "github.com/alessiosavi/GoGPUtils/cache"
```

One `Cache[K, V]` type covers three workloads:

| Workload | API |
|---|---|
| Bounded key/value store | `Config.MaxEntries`, `Set`, `Get`, `Peek`, `Delete` |
| TTL store | `Config.TTL`, `SetWithTTL`, `TTL`, `DeleteExpired`, `Config.JanitorInterval` + `Close` |
| Read-through memoizer | `GetOrLoad` (one load per key for concurrent misses) |

Standard library only. All methods are safe for concurrent use. The zero `Cache` is not usable; create one with `New`.

## Quick start

```go
c, err := cache.New(cache.Config[string, User]{
    MaxEntries: 10_000,
    TTL:        5 * time.Minute,
})
if err != nil {
    return err
}

c.Set("alice", alice)
if u, ok := c.Get("alice"); ok {
    fmt.Println(u.Name)
}

u, err := c.GetOrLoad(ctx, "bob", func(ctx context.Context) (User, error) {
    return repo.Find(ctx, "bob")
})
```

## Configuration

| Field | Meaning | Zero value |
|---|---|---|
| `MaxEntries` | Maximum stored entries | unbounded |
| `TTL` | Default lifetime | never expires |
| `Shards` | Shard count (rounded up to a power of two, ≤ `MaxShards`) | automatic (1 shard for a bound of 1–2047 entries) |
| `JanitorInterval` | Background removal of expired entries; call `Close` | no janitor |
| `OnEvict` | `func(key, value, reason)` for every removal | none |

`New` returns `ErrInvalidConfig` (naming the field) for negative values or `Shards > MaxShards`.

## Eviction (SIEVE)

When a shard is full, a hand walks from the oldest entry toward the newest: recently used entries get a second chance, and the first unused (or expired) entry is removed. Hits take only the shard's read lock and set an atomic bit, so reads on different shards do not contend. `MaxEntries` is split exactly across shards, so `Len()` never exceeds it when `MaxEntries > 0`; zero means unbounded. Capacity is enforced per shard, so evictions can begin slightly before `Len()` reaches `MaxEntries`. A cache bounded below 2048 entries gets one shard by default (exact global SIEVE order), where parallel readers contend on one lock; set `Shards` explicitly for read-heavy parallel workloads on small caches.

## Expiration

- `Set` uses `Config.TTL`; `SetWithTTL(k, v, ttl)` overrides it (`NoExpiration` = never).
- An entry is expired when `now >= deadline`; deadlines use the monotonic clock.
- `Get`, `Delete`, eviction, `DeleteExpired`, and the janitor remove expired entries; `Peek` and `TTL` never remove.
- `Close` stops the janitor deterministically. If a cache with a janitor becomes unreachable without `Close`, a runtime cleanup stops it on a best-effort basis only. Cleanup cannot run if `OnEvict` or stored keys/values keep the public wrapper reachable.

## Loading

`GetOrLoad(ctx, key, load)`:

- Concurrent misses for a key share one call to `load`.
- `load` runs on a context detached from caller cancellation (values preserved): one caller giving up never fails the others, and each caller still returns `ctx.Err()` when its own context ends. Give loaders their own timeouts.
- Results are stored with the default TTL unless `Set`, `Delete`, or `Clear` touched the key meanwhile.
- Errors are never cached. Panics become `*PanicError`; `runtime.Goexit` becomes `ErrLoaderGoexit`.
- Validation order: nil context (`ErrNilContext`), nil loader (`ErrNilLoader`), finished context (`ctx.Err()`).
- A loader must not call `GetOrLoad` for its own key.

## Callbacks, stats, iteration

| Reason | When | Counter |
|---|---|---|
| `ReasonEvicted` | removed to make room | `Stats.Evictions` |
| `ReasonExpired` | deadline passed | `Stats.Expirations` |
| `ReasonDeleted` | `Delete`, `Clear` | — |
| `ReasonReplaced` | `Set` over a live key | — |

`OnEvict` runs outside internal locks in the goroutine that caused the removal and may call any method except `Close`. Callbacks from different goroutines may run concurrently; within an operation they follow removal order. Callback panics are not recovered and skip the rest of the batch (Clear and DeleteExpired still finish every shard, without callbacks, before the panic propagates); a panic in a janitor or load goroutine terminates the process. `Stats()` also reports `Hits`, `Misses`, `Loads`, `LoadErrors`, and `HitRatio()`. `All()` iterates live entries (per-shard snapshot, unspecified order) without affecting stats or eviction order.

## Keys and values

Keys must equal themselves: `NaN` keys (or keys containing `NaN`) are never stored, and `GetOrLoad` loads them without caching. Values are stored as-is; synchronize mutable values yourself.
