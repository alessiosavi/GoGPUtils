// Package cache provides a generic, concurrency-safe in-memory cache.
//
// A Cache combines three workloads in one type: a bounded key/value store
// (entry-count limit with SIEVE eviction), a TTL store (per-entry expiration
// with lazy removal and an optional janitor), and a read-through memoizer
// (GetOrLoad with per-key load deduplication).
//
// Example:
//
//	c, err := cache.New(cache.Config[string, User]{
//	    MaxEntries: 10_000,
//	    TTL:        5 * time.Minute,
//	})
//	if err != nil {
//	    return err
//	}
//	c.Set("alice", alice)
//	u, ok := c.Get("alice")
//
// # Eviction
//
// When a shard is full, SIEVE (NSDI'24) picks the victim: a hand walks from
// the oldest entry toward the newest, giving recently used entries a second
// chance. Expired entries the hand reaches are removed first. Hits take only
// the shard's read lock and set an atomic bit. The cache is split into
// power-of-two shards; MaxEntries is divided exactly across them, so Len never
// exceeds it when MaxEntries > 0 (0 means unbounded). Capacity is enforced per
// shard, so evictions can begin slightly before Len reaches MaxEntries. The
// automatic shard count gives a cache bounded below 2048 entries a single
// shard (exact global SIEVE order), where parallel readers contend on one
// lock; set Config.Shards for read-heavy parallel workloads on small caches.
//
// # Expiration
//
// Config.TTL is the default lifetime (0 = never); SetWithTTL overrides it per
// entry (NoExpiration = never). Get, GetOrLoad, Delete, replacement, eviction,
// DeleteExpired, and the janitor remove or replace expired entries. Peek and
// TTL never remove them. The janitor is started with
// Config.JanitorInterval. Expiry uses the monotonic clock. Call Close to stop
// the janitor; a cleanup stops it on a best-effort basis if the Cache becomes
// unreachable, but Close is the only deterministic shutdown. Cleanup timing
// is not guaranteed; cleanup cannot run if OnEvict or stored keys or values
// keep the public Cache wrapper reachable.
//
// # Loading
//
//	u, err := c.GetOrLoad(ctx, id, func(ctx context.Context) (User, error) {
//	    return repo.Find(ctx, id)
//	})
//
// Concurrent misses for a key share one load, which runs on a context
// detached from caller cancellation; give loaders their own timeouts. Errors
// are never cached; panics become *PanicError. MaxEntries bounds stored
// entries, not concurrent loads, and Close does not cancel loads. A loader
// must not call GetOrLoad for its own key (or form cycles with other loads).
//
// # Callbacks and keys
//
// Config.OnEvict runs outside internal locks in the goroutine that caused the
// removal; it may call any method except Close. Callbacks from different
// goroutines may run concurrently; each operation reports removals in removal
// order. Callback panics are not recovered and skip the rest of the batch (Clear
// and DeleteExpired still finish every shard, without callbacks, before the panic
// propagates); a panic in a janitor or load goroutine terminates the process.
// Values are stored as-is (shallow), so
// mutable values need their own synchronization. Keys must equal themselves:
// NaN keys (or keys containing NaN) are never stored.
//
// # Zero value
//
// The zero Cache is not usable; create caches with New and do not copy them.
package cache
