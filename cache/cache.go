package cache

import (
	"errors"
	"fmt"
	"hash/maphash"
	"math/bits"
	"runtime"
	"sync"
	"time"
)

// NoExpiration tells SetWithTTL to store an entry that never expires.
// TTL returns it for entries without an expiration.
const NoExpiration time.Duration = -1

// MaxShards is the largest accepted Config.Shards value.
const MaxShards = 1 << 16

// Errors returned by the cache.
var (
	// ErrInvalidConfig is returned by New for an invalid Config.
	ErrInvalidConfig = errors.New("cache: invalid config")
	// ErrNilLoader is returned by GetOrLoad when the loader is nil.
	ErrNilLoader = errors.New("cache: nil loader")
	// ErrNilContext is returned by GetOrLoad when the context is nil.
	ErrNilContext = errors.New("cache: nil context")
	// ErrLoaderGoexit is returned by GetOrLoad when the loader calls runtime.Goexit.
	ErrLoaderGoexit = errors.New("cache: loader called runtime.Goexit")
)

// Config configures a Cache. The zero Config is valid: an unbounded cache
// whose entries never expire.
type Config[K comparable, V any] struct {
	// MaxEntries bounds the number of stored entries. 0 means unbounded.
	MaxEntries int
	// TTL is the default time-to-live. 0 means entries never expire.
	TTL time.Duration
	// Shards is the number of independent shards. 0 selects automatically.
	// Other values are rounded up to a power of two and clamped so that no
	// shard has zero capacity.
	Shards int
	// JanitorInterval > 0 starts a goroutine that removes expired entries at
	// that interval. Call Close to stop it. 0 disables the janitor.
	JanitorInterval time.Duration
	// OnEvict, if set, is called for every entry removed from the cache,
	// outside any internal lock, in the goroutine that caused the removal.
	// It may call any Cache method except Close.
	OnEvict func(key K, value V, reason EvictionReason)
}

func (cfg Config[K, V]) validate() error {
	switch {
	case cfg.MaxEntries < 0:
		return fmt.Errorf("%w: MaxEntries must be >= 0, got %d", ErrInvalidConfig, cfg.MaxEntries)
	case cfg.TTL < 0:
		return fmt.Errorf("%w: TTL must be >= 0, got %v", ErrInvalidConfig, cfg.TTL)
	case cfg.Shards < 0 || cfg.Shards > MaxShards:
		return fmt.Errorf("%w: Shards must be in [0, %d], got %d", ErrInvalidConfig, MaxShards, cfg.Shards)
	case cfg.JanitorInterval < 0:
		return fmt.Errorf("%w: JanitorInterval must be >= 0, got %v", ErrInvalidConfig, cfg.JanitorInterval)
	}

	return nil
}

// EvictionReason tells OnEvict why an entry was removed.
type EvictionReason int

// Eviction reasons.
const (
	// ReasonEvicted: removed to make room for a new entry.
	ReasonEvicted EvictionReason = iota + 1
	// ReasonExpired: its TTL elapsed.
	ReasonExpired
	// ReasonDeleted: removed by Delete or Clear.
	ReasonDeleted
	// ReasonReplaced: overwritten by Set or SetWithTTL.
	ReasonReplaced
)

// String returns "evicted", "expired", "deleted", or "replaced".
func (r EvictionReason) String() string {
	switch r {
	case ReasonEvicted:
		return "evicted"
	case ReasonExpired:
		return "expired"
	case ReasonDeleted:
		return "deleted"
	case ReasonReplaced:
		return "replaced"
	default:
		return fmt.Sprintf("EvictionReason(%d)", int(r))
	}
}

// PanicError is returned by GetOrLoad when the loader panics.
type PanicError struct {
	// Value is the value passed to panic.
	Value any
	// Stack is the stack trace captured when the panic was recovered.
	Stack []byte
}

// Error implements error.
func (e *PanicError) Error() string {
	return fmt.Sprintf("cache: loader panicked: %v", e.Value)
}

// Unwrap returns the panic value if it is an error, else nil.
func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error)

	return err
}

// Cache is a concurrency-safe, generic in-memory cache with SIEVE eviction,
// TTL expiration, and deduplicated loading. The zero value is not usable;
// create one with New. Do not copy a Cache; use the pointer returned by New.
type Cache[K comparable, V any] struct {
	in          *inner[K, V]
	cleanup     runtime.Cleanup
	hasCleanup  bool
	cleanupOnce sync.Once
}

// inner holds all cache state. It is allocated separately from Cache so the
// janitor can run without keeping the public wrapper reachable.
type inner[K comparable, V any] struct {
	shards  []*shard[K, V]
	mask    uint64
	seed    maphash.Seed
	ttl     time.Duration
	now     func() int64
	onEvict func(K, V, EvictionReason)
	hooks   testHooks

	stop        chan struct{}
	stopOnce    sync.Once
	janitorDone chan struct{} // nil when no janitor was started
}

// testHooks lets in-package tests force interleavings. All fields are nil in
// production.
type testHooks struct {
	// afterReadUnlock runs between releasing the read lock and acquiring the
	// write lock on the Get-expired and GetOrLoad-miss paths.
	afterReadUnlock func()
}

// New creates a cache from cfg.
//
// Example:
//
//	c, err := cache.New(cache.Config[string, User]{MaxEntries: 10_000, TTL: 5 * time.Minute})
//	if err != nil {
//	    return err
//	}
//	c.Set("alice", alice)
func New[K comparable, V any](cfg Config[K, V]) (*Cache[K, V], error) {
	return newCache(cfg, monotonicClock(), testHooks{})
}

// monotonicClock returns nanoseconds elapsed since its creation. It uses the
// monotonic clock, so wall-clock changes never expire or revive entries.
func monotonicClock() func() int64 {
	epoch := time.Now()

	return func() int64 { return int64(time.Since(epoch)) }
}

func newCache[K comparable, V any](cfg Config[K, V], now func() int64, hooks testHooks) (*Cache[K, V], error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	n := shardCount(cfg.MaxEntries, cfg.Shards, runtime.GOMAXPROCS(0))
	in := &inner[K, V]{
		shards:  make([]*shard[K, V], n),
		mask:    uint64(n - 1), //nolint:gosec // shardCount guarantees 1 <= n <= MaxShards.
		seed:    maphash.MakeSeed(),
		ttl:     cfg.TTL,
		now:     now,
		onEvict: cfg.OnEvict,
		hooks:   hooks,
		stop:    make(chan struct{}),
	}
	for i, capacity := range splitCapacity(cfg.MaxEntries, n) {
		in.shards[i] = newShard[K, V](capacity)
	}

	return &Cache[K, V]{in: in}, nil
}

// shardCount returns the number of shards (a power of two) for the config.
// Explicit counts are rounded up and clamped to maxEntries when bounded; the
// automatic count keeps at least 1024 entries per shard when bounded.
func shardCount(maxEntries, shards, procs int) int {
	if shards > 0 {
		n := nextPow2(shards)
		if maxEntries > 0 && n > maxEntries {
			n = prevPow2(maxEntries)
		}

		return n
	}

	base := nextPow2(min(max(4*procs, 1), MaxShards))
	if maxEntries == 0 {
		return base
	}

	return min(base, prevPow2(max(1, maxEntries/1024)))
}

// nextPow2 returns the smallest power of two >= n (1 for n <= 1).
func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}

	return 1 << bits.Len(uint(n-1))
}

// prevPow2 returns the largest power of two <= n. n must be >= 1.
func prevPow2(n int) int {
	return 1 << (bits.Len(uint(n)) - 1)
}

// splitCapacity divides maxEntries across n shards; the first maxEntries%n
// shards get one extra slot, so the capacities sum to maxEntries.
func splitCapacity(maxEntries, n int) []int {
	caps := make([]int, n)
	if maxEntries == 0 {
		return caps
	}

	base, extra := maxEntries/n, maxEntries%n
	for i := range caps {
		caps[i] = base
		if i < extra {
			caps[i]++
		}
	}

	return caps
}

func (c *inner[K, V]) shardFor(key K) *shard[K, V] {
	if len(c.shards) == 1 {
		return c.shards[0]
	}

	return c.shards[maphash.Comparable(c.seed, key)&c.mask]
}

// reflexive reports whether key equals itself. It is false for NaN floats and
// for arrays, structs, and interfaces containing them; such keys can never be
// found in a Go map, so the cache does not store them. Comparing an interface
// key whose dynamic value is not comparable panics.
func reflexive[K comparable](key K) bool {
	other := key

	return key == other
}
