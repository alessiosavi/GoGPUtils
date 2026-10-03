package cache

import (
	"maps"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestStats_CountersFollowRemovalReasons(t *testing.T) {
	c, clk := newTestCache(t, Config[string, int]{MaxEntries: 2, Shards: 1})
	c.Get("missing") // Misses=1
	c.Set("a", 1)
	c.Get("a")    // Hits=1
	c.Set("a", 2) // replaced: no counter
	c.Set("b", 2)
	c.Set("c", 3) // one capacity eviction
	c.SetWithTTL("c", 3, time.Second)
	clk.advance(2 * time.Second)
	c.Get("c") // expired: Expirations=1, Misses=2
	c.Clear()  // deleted: no counter
	want := Stats{Hits: 1, Misses: 2, Evictions: 1, Expirations: 1}
	if got := c.Stats(); got != want {
		t.Fatalf("Stats = %+v, want %+v", got, want)
	}
}

func TestStats_SumsAcrossShards(t *testing.T) {
	c, _ := newTestCache(t, Config[int, int]{Shards: 8})
	for i := range 1000 {
		c.Set(i, i)
	}
	for i := range 1500 {
		c.Get(i)
	}
	st := c.Stats()
	if st.Hits != 1000 || st.Misses != 500 {
		t.Fatalf("Stats = %+v, want 1000 hits and 500 misses", st)
	}
	if got := st.HitRatio(); got != 1000.0/1500.0 {
		t.Fatalf("HitRatio = %v", got)
	}
}

func TestPeekTTLAndAllDoNotTouchStatsOrSieve(t *testing.T) {
	rec := &recorder[string, int]{}
	c, _ := newTestCache(t, Config[string, int]{MaxEntries: 2, Shards: 1, OnEvict: rec.record})
	c.Set("a", 1)
	c.Set("b", 2)
	c.Peek("a")
	c.TTL("a")
	for range c.All() {
	}
	c.Set("c", 3) // "a" was never marked visited, so it is the victim
	expectRemovals(t, rec.take(), []removal[string, int]{{"a", 1, ReasonEvicted}})
	if got := c.Stats(); got != (Stats{Evictions: 1}) {
		t.Fatalf("Stats = %+v, want only one eviction", got)
	}
}

func TestAll(t *testing.T) {
	c, clk := newTestCache(t, Config[int, int]{Shards: 4})
	for i := range 10 {
		c.Set(i, i*10)
	}
	c.SetWithTTL(100, 1, time.Second)
	clk.advance(time.Second)

	got := maps.Collect(c.All())
	want := map[int]int{}
	for i := range 10 {
		want[i] = i * 10
	}
	if !maps.Equal(got, want) {
		t.Fatalf("All = %v, want %v (expired entry excluded)", got, want)
	}

	n := 0
	for range c.All() {
		if n++; n == 3 {
			break
		}
	}
	if n != 3 {
		t.Fatalf("early break yielded %d entries", n)
	}

	for k := range c.All() { // mutating during iteration must not deadlock
		c.Delete(k)
		c.Set(k+1000, k)
	}

	seq := c.All()
	count := func() int {
		total := 0
		for range seq {
			total++
		}

		return total
	}
	if first, second := count(), count(); first != second {
		t.Fatalf("re-invoking the iterator gave %d then %d entries", first, second)
	}
	checkInvariants(t, c)
}

func TestOnEvict_RunsWithoutLocksAndMayReenter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c *Cache[int, int]
		var nested atomic.Int64
		cfg := Config[int, int]{MaxEntries: 1, Shards: 1, OnEvict: func(k, _ int, r EvictionReason) {
			// Re-entering the same shard deadlocks if a lock is still held.
			c.Get(k)
			c.Peek(k)
			c.Len()
			for range c.All() {
			}
			if r == ReasonEvicted && nested.Add(1) == 1 {
				c.Set(1000, 1000) // triggers one nested eviction
			}
		}}
		c, _ = newTestCache(t, cfg)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Set(1, 1)
			c.Set(2, 2)
			c.Delete(1000)
			c.Set(3, 3)
			c.Clear()
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("re-entrant OnEvict deadlocked")
		}
		checkInvariants(t, c)
	})
}

func TestOnEvict_PanicLeavesCacheUsable(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		cfg := Config[int, int]{MaxEntries: 1, Shards: 1, OnEvict: func(int, int, EvictionReason) {
			panic("callback failed")
		}}
		c, _ := newTestCache(t, cfg)
		c.Set(1, 1)
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expected the callback panic to propagate to the caller")
				}
			}()
			c.Set(2, 2) // evicts 1; the callback panics after the lock is released
		}()
		if v, ok := c.Get(2); !ok || v != 2 {
			t.Fatal("cache unusable after a callback panic")
		}
		checkInvariants(t, c)
	})
}
