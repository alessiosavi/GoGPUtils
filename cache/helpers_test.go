package cache

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
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

// waitFor checks a condition after bubble goroutines are durably blocked.
// Call it only from inside synctest.Test.
func waitFor(t testing.TB, what string, cond func() bool) {
	t.Helper()
	synctest.Wait()
	if !cond() {
		t.Fatalf("condition not reached: %s", what)
	}
}

func await[T any](t testing.TB, what string, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
		var zero T
		return zero
	}
}

type result[V any] struct {
	v   V
	err error
}

// loadAsync runs GetOrLoad in a goroutine and delivers its result.
func loadAsync[K comparable, V any](ctx context.Context, c *Cache[K, V], key K, load func(context.Context) (V, error)) <-chan result[V] {
	ch := make(chan result[V], 1)
	go func() {
		v, err := c.GetOrLoad(ctx, key, load)
		ch <- result[V]{v, err}
	}()

	return ch
}

func recv[V any](t testing.TB, ch <-chan result[V]) result[V] {
	t.Helper()

	select {
	case r := <-ch:
		return r
	case <-time.After(5 * time.Second):
		t.Fatal("GetOrLoad did not return")

		return result[V]{}
	}
}

// blockingLoader returns a loader that closes started when it begins and then
// waits for release before returning value.
func blockingLoader[V any](value V) (load func(context.Context) (V, error), started <-chan struct{}, release func()) {
	s, r := make(chan struct{}), make(chan struct{})
	var once sync.Once
	load = func(context.Context) (V, error) {
		once.Do(func() { close(s) })
		<-r

		return value, nil
	}

	return load, s, sync.OnceFunc(func() { close(r) })
}

func mustNotLoad[V any](t testing.TB) func(context.Context) (V, error) {
	return func(context.Context) (V, error) {
		t.Error("loader must not run")

		var zero V

		return zero, nil
	}
}
