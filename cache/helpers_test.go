package cache

import (
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// fakeClock is a concurrency-safe manual clock in nanoseconds.
type fakeClock struct{ t atomic.Int64 }

func newFakeClock() *fakeClock {
	clk := &fakeClock{}
	clk.t.Store(1) // start after the epoch so deadlines are never 0

	return clk
}

func (f *fakeClock) now() int64              { return f.t.Load() }
func (f *fakeClock) advance(d time.Duration) { f.t.Add(int64(d)) }

// recorder collects OnEvict calls; safe for concurrent use.
type recorder[K comparable, V any] struct {
	mu  sync.Mutex
	got []removal[K, V]
}

func (r *recorder[K, V]) record(key K, value V, reason EvictionReason) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.got = append(r.got, removal[K, V]{key: key, value: value, reason: reason})
}

// take returns and clears the recorded removals.
func (r *recorder[K, V]) take() []removal[K, V] {
	r.mu.Lock()
	defer r.mu.Unlock()

	got := r.got
	r.got = nil

	return got
}

func newTestCache[K comparable, V any](t testing.TB, cfg Config[K, V]) (*Cache[K, V], *fakeClock) {
	t.Helper()

	return newTestCacheWithHooks(t, cfg, testHooks{})
}

func newTestCacheWithHooks[K comparable, V any](t testing.TB, cfg Config[K, V], hooks testHooks) (*Cache[K, V], *fakeClock) {
	t.Helper()

	clk := newFakeClock()
	c, err := newCache(cfg, clk.now, hooks)
	if err != nil {
		t.Fatalf("newCache: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })

	return c, clk
}

// checkInvariants asserts that every shard's map and list agree, the hand is
// nil or a linked entry, and no shard exceeds its capacity.
func checkInvariants[K comparable, V any](t testing.TB, c *Cache[K, V]) {
	t.Helper()

	for i, s := range c.in.shards {
		s.mu.RLock()
		linked := make(map[*entry[K, V]]bool, len(s.items))
		var prev *entry[K, V]
		for e := s.tail; e != nil; e = e.next {
			if e.prev != prev {
				s.mu.RUnlock()
				t.Fatalf("shard %d: broken prev link at key %v", i, e.key)
			}
			if s.items[e.key] != e {
				s.mu.RUnlock()
				t.Fatalf("shard %d: list entry %v is not in the map", i, e.key)
			}
			linked[e] = true
			prev = e
		}
		ok := prev == s.head && len(linked) == len(s.items) &&
			(s.hand == nil || linked[s.hand]) &&
			(s.capacity == 0 || len(s.items) <= s.capacity)
		n, capacity := len(s.items), s.capacity
		s.mu.RUnlock()
		if !ok {
			t.Fatalf("shard %d: invariant broken (list %d, map %d, capacity %d, hand linked=%v)",
				i, len(linked), n, capacity, s.hand == nil || linked[s.hand])
		}
	}
}

// expectRemovals fails unless got equals want, in order.
func expectRemovals[K comparable, V comparable](t testing.TB, got, want []removal[K, V]) {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Fatalf("removals = %v, want %v", got, want)
	}
}
