package cache

import (
	"context"
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

func TestClear_PanickingOnEvictStillClearsEveryShard(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const panicValue = "callback failed"
		var armed atomic.Bool
		var callbacks atomic.Int64
		c, _ := newTestCache(t, Config[int, int]{Shards: 4, OnEvict: func(int, int, EvictionReason) {
			callbacks.Add(1)
			if armed.CompareAndSwap(true, false) {
				panic(panicValue)
			}
		}})
		// Two entries per shard also exercise skipping the current shard's
		// remaining callbacks. Choose keys by their actual shard mapping.
		key := 0
		for _, s := range c.in.shards {
			for n := 0; n < 2; key++ {
				if c.in.shardFor(key) == s {
					c.Set(key, key)
					n++
				}
			}
		}
		for c.in.shardFor(key) != c.in.shards[len(c.in.shards)-1] {
			key++
		}
		load, started, release := blockingLoader(777)
		defer release()
		waiter := loadAsync(context.Background(), c, key, load)
		await(t, "last-shard loader started", started)
		before := c.Stats()

		armed.Store(true)
		func() {
			defer func() {
				if got := recover(); got != panicValue {
					t.Errorf("Clear panic = %v, want %q", got, panicValue)
				}
			}()
			c.Clear()
		}()
		if got := c.Len(); got != 0 {
			t.Errorf("Len after Clear = %d, want 0", got)
		}
		if got := callbacks.Load(); got != 1 {
			t.Errorf("OnEvict calls = %d, want 1", got)
		}
		if got := c.Stats(); got != before {
			t.Errorf("Clear changed Stats from %+v to %+v", before, got)
		}

		release()
		if r := recv(t, waiter); r.err != nil || r.v != 777 {
			t.Fatalf("existing waiter got %v, %v; want 777, nil", r.v, r.err)
		}
		if _, ok := c.Peek(key); ok {
			t.Error("a load that finished after Clear resurrected the last-shard key")
		}
		if got := c.Len(); got != 0 {
			t.Errorf("Len after load completed = %d, want 0", got)
		}
		checkInvariants(t, c)
	})
}

func TestDeleteExpired_PanickingOnEvictStillProcessesEveryShard(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const panicValue = "callback failed"
		var armed atomic.Bool
		var callbacks atomic.Int64
		c, clk := newTestCache(t, Config[int, int]{Shards: 4, OnEvict: func(int, int, EvictionReason) {
			callbacks.Add(1)
			if armed.CompareAndSwap(true, false) {
				panic(panicValue)
			}
		}})
		// Put two expiring entries and one non-expiring entry in every shard.
		retained := make(map[int]int)
		key := 0
		for _, s := range c.in.shards {
			for n := 0; n < 3; key++ {
				if c.in.shardFor(key) != s {
					continue
				}
				if n < 2 {
					c.SetWithTTL(key, key, time.Second)
				} else {
					c.SetWithTTL(key, key, NoExpiration)
					retained[key] = key
				}
				n++
			}
		}
		clk.advance(2 * time.Second)

		armed.Store(true)
		func() {
			defer func() {
				if got := recover(); got != panicValue {
					t.Errorf("DeleteExpired panic = %v, want %q", got, panicValue)
				}
			}()
			c.DeleteExpired()
		}()
		if got := c.Len(); got != len(retained) {
			t.Errorf("Len after DeleteExpired = %d, want %d", got, len(retained))
		}
		if got := maps.Collect(c.All()); !maps.Equal(got, retained) {
			t.Errorf("remaining live entries = %v, want %v", got, retained)
		}
		if got := c.Stats(); got != (Stats{Expirations: 8}) {
			t.Errorf("Stats = %+v, want only 8 expirations", got)
		}
		if got := callbacks.Load(); got != 1 {
			t.Errorf("OnEvict calls = %d, want 1", got)
		}
		checkInvariants(t, c)
	})
}
