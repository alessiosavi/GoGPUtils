package cache

import (
	"container/list"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// benchCache fills every shard to capacity without evicting setup keys. The
// public cache's random hash seed can distribute consecutive keys unevenly,
// so setup accepts a key only while its shard has room. Values come from
// benchkit; keys are returned in insertion order for repeatable access cycles.
func benchCache(b *testing.B, values []int, shards int) (*Cache[int, int], []int) {
	b.Helper()
	c, err := New(Config[int, int]{MaxEntries: len(values), Shards: shards})
	if err != nil {
		b.Fatal(err)
	}

	return c, benchFillCache(b, c, values)
}

// benchFillCache returns keys distributed to fill the supplied cache's shards.
func benchFillCache(b *testing.B, c *Cache[int, int], values []int) []int {
	b.Helper()
	keys := make([]int, 0, len(values))
	filled := make(map[*shard[int, int]]int, len(c.in.shards))
	for key := 0; len(keys) < len(values); key++ {
		s := c.in.shardFor(key)
		if filled[s] == s.capacity {
			continue
		}
		c.Set(key, values[len(keys)])
		keys = append(keys, key)
		filled[s]++
	}
	if got := c.Len(); got != len(values) {
		b.Fatalf("setup entries = %d, want %d", got, len(values))
	}
	for i, key := range keys {
		if v, ok := c.Peek(key); !ok || v != values[i] {
			b.Fatalf("setup key %d: value=%d present=%t", key, v, ok)
		}
	}

	return keys
}

// BenchmarkNew measures empty cache construction. n is configured capacity
// for bounded and janitor cases; unbounded has no size parameter. The janitor
// case includes Close on every iteration so live goroutines stay bounded.
func BenchmarkNew(b *testing.B) {
	b.Run("case=unbounded", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := New(Config[int, int]{}); err != nil {
				b.Fatal(err)
			}
		}
	})
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=bounded", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := New(Config[int, int]{MaxEntries: n}); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run("case=janitor", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				c, err := New(Config[int, int]{MaxEntries: n, JanitorInterval: time.Hour})
				if err != nil {
					b.Fatal(err)
				}
				if err := c.Close(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkCache_Get measures resident hits, absent misses, and expired-entry
// removal. n is capacity and, except for expired, the initial entry count. Hit cases cycle through
// all resident keys with one or automatic shards, serially and in parallel.
// case=expired starts with n-1 resident entries and includes re-SetWithTTL(1ns)
// and a Len check that verifies removal each iteration. Its clock reads the real
// monotonic clock plus an atomic offset; a measured 1us advance after SetWithTTL
// guarantees expiry even when consecutive real clock reads return the same tick.
func BenchmarkCache_Get(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		for _, sh := range []struct {
			name  string
			count int
		}{{"1", 1}, {"auto", 0}} {
			b.Run("case=hit/shards="+sh.name+"/mode=serial", func(b *testing.B) {
				c, keys := benchCache(b, values, sh.count)
				b.ReportAllocs()
				i := 0
				for b.Loop() {
					j := i & (n - 1)
					if v, ok := c.Get(keys[j]); !ok || v != values[j] {
						b.Fatalf("Get hit: value=%d present=%t", v, ok)
					}
					i++
				}
			})
			b.Run("case=hit/shards="+sh.name+"/mode=parallel", func(b *testing.B) {
				c, keys := benchCache(b, values, sh.count)
				var failed atomic.Pointer[error]
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						j := i & (n - 1)
						if v, ok := c.Get(keys[j]); !ok || v != values[j] {
							err := fmt.Errorf("Get hit: value=%d present=%t", v, ok)
							failed.CompareAndSwap(nil, &err)
							b.Fail()

							return
						}
						i++
					}
				})
				if p := failed.Load(); p != nil {
					b.Fatal(*p)
				}
			})
		}
		b.Run("case=miss", func(b *testing.B) {
			c, _ := benchCache(b, values, 0)
			b.ReportAllocs()
			for b.Loop() {
				if _, ok := c.Get(-1); ok {
					b.Fatal("Get returned an absent key")
				}
			}
		})
		b.Run("case=expired", func(b *testing.B) {
			now := monotonicClock()
			var offset atomic.Int64
			c, err := newCache(Config[int, int]{MaxEntries: n}, func() int64 {
				return now() + offset.Load()
			}, testHooks{})
			if err != nil {
				b.Fatal(err)
			}
			keys := benchFillCache(b, c, values)
			if !c.Delete(keys[0]) {
				b.Fatal("expired setup did not remove the target key")
			}
			b.ReportAllocs()
			for b.Loop() {
				c.SetWithTTL(keys[0], values[0], time.Nanosecond)
				offset.Add(int64(time.Microsecond))
				if _, ok := c.Get(keys[0]); ok {
					b.Fatal("Get returned an expired key")
				}
				if c.Len() != n-1 {
					b.Fatal("Get did not remove the expired entry")
				}
			}
		})
	})
}

// BenchmarkCache_Peek measures hits cycling over n resident entries in an
// automatically sharded cache. Each returned value is checked.
func BenchmarkCache_Peek(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=hit", func(b *testing.B) {
			values := benchkit.Ints(n)
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				j := i & (n - 1)
				if v, ok := c.Peek(keys[j]); !ok || v != values[j] {
					b.Fatalf("Peek hit: value=%d present=%t", v, ok)
				}
				i++
			}
		})
	})
}

// BenchmarkCache_Set measures updates to n resident entries, serially and
// in parallel, and new insertions into a full automatically sharded cache.
// n is capacity; case=new uses fresh increasing keys, evicting on every Set.
// Both cases keep live state bounded by n without measured restoration.
func BenchmarkCache_Set(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		b.Run("case=update/mode=serial", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				j := i & (n - 1)
				c.Set(keys[j], values[j])
				i++
			}
		})
		b.Run("case=update/mode=parallel", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					j := i & (n - 1)
					c.Set(keys[j], values[j])
					i++
				}
			})
		})
		b.Run("case=new", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			next := keys[len(keys)-1] + 1
			b.ReportAllocs()
			for b.Loop() {
				c.Set(next, values[next&(n-1)])
				next++
			}
		})
	})
}

// BenchmarkCache_Delete measures successful deletion and re-Set of a resident
// key in a full cache of n entries. The restoring Set is measured.
func BenchmarkCache_Delete(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=delete-and-set", func(b *testing.B) {
			values := benchkit.Ints(n)
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			for b.Loop() {
				if !c.Delete(keys[0]) {
					b.Fatal("Delete did not remove a live entry")
				}
				c.Set(keys[0], values[0])
			}
		})
	})
}

// BenchmarkCache_Len measures Len on a full automatically sharded cache;
// n is resident entries. The returned count is checked.
func BenchmarkCache_Len(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("shards=auto", func(b *testing.B) {
			c, _ := benchCache(b, benchkit.Ints(n), 0)
			b.ReportAllocs()
			for b.Loop() {
				if got := c.Len(); got != n {
					b.Fatalf("Len = %d, want %d", got, n)
				}
			}
		})
	})
}

// BenchmarkCache_Clear measures refilling n entries and clearing all of them.
// The n Set calls and before/after Len checks are included in the measurement.
func BenchmarkCache_Clear(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=refill-and-clear", func(b *testing.B) {
			values := benchkit.Ints(n)
			c, keys := benchCache(b, values, 0)
			c.Clear()
			b.ReportAllocs()
			for b.Loop() {
				for i, key := range keys {
					c.Set(key, values[i])
				}
				if c.Len() != n {
					b.Fatal("Clear requires a full cache")
				}
				c.Clear()
				if c.Len() != 0 {
					b.Fatal("Clear left entries behind")
				}
			}
		})
	})
}

// BenchmarkCache_SetWithTTL measures updating resident keys with a one-hour
// TTL. n is the full cache's entry count; cycling keys keeps state bounded.
func BenchmarkCache_SetWithTTL(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=update", func(b *testing.B) {
			values := benchkit.Ints(n)
			c, keys := benchCache(b, values, 0)
			for i, key := range keys {
				c.SetWithTTL(key, values[i], time.Hour)
			}
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				j := i & (n - 1)
				c.SetWithTTL(keys[j], values[j], time.Hour)
				i++
			}
		})
	})
}

// BenchmarkCache_TTL measures remaining-lifetime hits over n resident entries
// with a one-hour TTL. Each result must be present and strictly positive.
func BenchmarkCache_TTL(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=hit", func(b *testing.B) {
			values := benchkit.Ints(n)
			c, keys := benchCache(b, values, 0)
			for i, key := range keys {
				c.SetWithTTL(key, values[i], time.Hour)
			}
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				if ttl, ok := c.TTL(keys[i&(n-1)]); !ok || ttl <= 0 {
					b.Fatalf("TTL hit: ttl=%v present=%t", ttl, ok)
				}
				i++
			}
		})
	})
}

// BenchmarkCache_DeleteExpired measures refilling n entries with a 1ns TTL
// and deleting all of them. The refill and precondition Len check are measured;
// each iteration requires exactly n expired removals. The clock reads the real
// monotonic clock plus an atomic offset; a measured 1us advance after the refill
// guarantees expiry even when consecutive real clock reads return the same tick.
func BenchmarkCache_DeleteExpired(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=refill-and-expire", func(b *testing.B) {
			values := benchkit.Ints(n)
			now := monotonicClock()
			var offset atomic.Int64
			c, err := newCache(Config[int, int]{MaxEntries: n}, func() int64 {
				return now() + offset.Load()
			}, testHooks{})
			if err != nil {
				b.Fatal(err)
			}
			keys := benchFillCache(b, c, values)
			c.Clear()
			b.ReportAllocs()
			for b.Loop() {
				for i, key := range keys {
					c.SetWithTTL(key, values[i], time.Nanosecond)
				}
				offset.Add(int64(time.Microsecond))
				if c.Len() != n {
					b.Fatal("DeleteExpired requires n stored entries")
				}
				if removed := c.DeleteExpired(); removed != n {
					b.Fatalf("DeleteExpired removed %d, want %d", removed, n)
				}
			}
		})
	})
}

// BenchmarkCache_All measures consuming a fresh iterator over n resident
// entries, including its snapshot allocation and count/value checksum checks.
func BenchmarkCache_All(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		c, _ := benchCache(b, values, 0)
		want := 0
		for _, v := range values {
			want ^= v
		}
		b.ReportAllocs()
		for b.Loop() {
			count, checksum := 0, 0
			for _, v := range c.All() {
				count++
				checksum ^= v
			}
			if count != n || checksum != want {
				b.Fatalf("All: count=%d checksum=%d, want %d and %d", count, checksum, n, want)
			}
		}
	})
}

// BenchmarkCache_Stats measures counter aggregation on a full automatically
// sharded cache of n entries, prepopulated with n hits and n/4 misses.
func BenchmarkCache_Stats(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("shards=auto", func(b *testing.B) {
			c, keys := benchCache(b, benchkit.Ints(n), 0)
			for _, key := range keys {
				if _, ok := c.Get(key); !ok {
					b.Fatal("Stats setup missed a resident key")
				}
			}
			for range n / 4 {
				if _, ok := c.Get(-1); ok {
					b.Fatal("Stats setup hit an absent key")
				}
			}
			b.ReportAllocs()
			for b.Loop() {
				if s := c.Stats(); s.Hits != uint64(n) || s.Misses != uint64(n/4) {
					b.Fatalf("unexpected Stats: %+v", s)
				}
			}
		})
	})
}

// BenchmarkCache_GetOrLoad measures hits on n resident keys serially and in
// parallel, and a miss loading a fresh key and evicting one. n is capacity.
// Hits reject loader execution; misses include the loader goroutine and store.
func BenchmarkCache_GetOrLoad(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		ctx := context.Background()
		unexpectedLoad := errors.New("GetOrLoad hit invoked the loader")
		hitLoad := func(context.Context) (int, error) { return 0, unexpectedLoad }
		load := func(ctx context.Context) (int, error) { return values[0], ctx.Err() }
		b.Run("case=hit/mode=serial", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				j := i & (n - 1)
				v, err := c.GetOrLoad(ctx, keys[j], hitLoad)
				if err != nil {
					b.Fatal(err)
				}
				if v != values[j] {
					b.Fatalf("GetOrLoad hit = %d, want %d", v, values[j])
				}
				i++
			}
		})
		b.Run("case=hit/mode=parallel", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			var failed atomic.Pointer[error]
			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					j := i & (n - 1)
					v, err := c.GetOrLoad(ctx, keys[j], hitLoad)
					if err == nil && v != values[j] {
						err = fmt.Errorf("GetOrLoad hit = %d, want %d", v, values[j])
					}
					if err != nil {
						failed.CompareAndSwap(nil, &err)
						b.Fail()

						return
					}
					i++
				}
			})
			if p := failed.Load(); p != nil {
				b.Fatal(*p)
			}
		})
		b.Run("case=miss", func(b *testing.B) {
			c, keys := benchCache(b, values, 0)
			b.ReportAllocs()
			next := keys[len(keys)-1] + 1
			for b.Loop() {
				v, err := c.GetOrLoad(ctx, next, load)
				if err != nil {
					b.Fatal(err)
				}
				if v != values[0] {
					b.Fatalf("GetOrLoad miss = %d, want %d", v, values[0])
				}
				next++
			}
		})
	})
}

// BenchmarkCache_Close measures New plus Close of an empty janitor cache on
// every iteration; n is configured capacity. Construction is the measured
// restoration, so Close always stops a fresh janitor rather than a closed one.
func BenchmarkCache_Close(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.Run("case=new-and-close", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				c, err := New(Config[int, int]{MaxEntries: n, JanitorInterval: time.Hour})
				if err != nil {
					b.Fatal(err)
				}
				if err := c.Close(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkStats_HitRatio measures the scalar ratio calculation on a populated
// snapshot with 900 hits and 100 misses; there is no size parameter.
func BenchmarkStats_HitRatio(b *testing.B) {
	s := Stats{Hits: 900, Misses: 100}
	b.ReportAllocs()
	for b.Loop() {
		s.HitRatio()
	}
}

// BenchmarkCache_MixedParallel measures automatically sharded contention with
// 90% or 50% reads (in expectation). n is capacity, initially full; uniformly
// sampled keys span [0,2*n), mixing hits, misses, updates and eviction. Each
// worker has a PCG stream seeded from benchkit, with one atomic stream assignment
// at worker startup. Random draws are measured; live entries stay bounded by n.
func BenchmarkCache_MixedParallel(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		seed := benchkit.Bytes(16)
		seed1 := binary.LittleEndian.Uint64(seed[:8])
		seed2 := binary.LittleEndian.Uint64(seed[8:])
		for _, readPct := range []int{90, 50} {
			b.Run(fmt.Sprintf("read=%d/mode=parallel", readPct), func(b *testing.B) {
				c, _ := benchCache(b, values, 0)
				var streams atomic.Uint64
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					rng := rand.New(rand.NewPCG(seed1, seed2+streams.Add(1)))
					checksum := 0
					for pb.Next() {
						k := rng.IntN(2 * n)
						if rng.IntN(100) < readPct {
							v, ok := c.Get(k)
							if ok {
								checksum ^= v
							}
						} else {
							c.Set(k, values[k&(n-1)])
						}
					}
					runtime.KeepAlive(checksum)
				})
			})
		}
	})
}

// BenchmarkBaselineRWMutexMapParallel measures parallel hits on an RWMutex-
// guarded map with n entries, cycling the same key/value workload as Get hits.
func BenchmarkBaselineRWMutexMapParallel(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		values := benchkit.Ints(n)
		var mu sync.RWMutex
		m := make(map[int]int, n)
		for i, v := range values {
			m[i] = v
		}
		var failed atomic.Pointer[error]
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				j := i & (n - 1)
				mu.RLock()
				v, ok := m[j]
				mu.RUnlock()
				if !ok || v != values[j] {
					err := fmt.Errorf("map hit: value=%d present=%t", v, ok)
					failed.CompareAndSwap(nil, &err)
					b.Fail()

					return
				}
				i++
			}
		})
		if p := failed.Load(); p != nil {
			b.Fatal(*p)
		}
	})
}

// lruCache is a minimal LRU used only to report a comparative hit ratio.
type lruCache struct {
	capacity int
	ll       *list.List
	items    map[int]*list.Element
}

func (l *lruCache) access(k int) bool {
	if e, ok := l.items[k]; ok {
		l.ll.MoveToFront(e)

		return true
	}
	if l.ll.Len() == l.capacity {
		oldest := l.ll.Back()
		l.ll.Remove(oldest)
		delete(l.items, oldest.Value.(int))
	}
	l.items[k] = l.ll.PushFront(k)

	return false
}

// BenchmarkHitRatio_SIEVEvsLRU reports informational hit ratios on the same
// 200000-access Zipf(s=1.07,v=1,imax=100000) trace, seeded from benchkit bytes.
// Each iteration includes constructing and populating both initially empty
// caches of fixed capacity 1000, keeping state bounded and the trace identical.
func BenchmarkHitRatio_SIEVEvsLRU(b *testing.B) {
	const capacity, accesses = 1000, 200_000
	seed := benchkit.Bytes(16)
	zipf := rand.NewZipf(rand.New(rand.NewPCG(binary.LittleEndian.Uint64(seed[:8]), binary.LittleEndian.Uint64(seed[8:]))), 1.07, 1, 100_000)
	trace := make([]int, accesses)
	for i := range trace {
		trace[i] = int(zipf.Uint64())
	}
	b.ReportAllocs()
	for b.Loop() {
		c, err := New(Config[int, int]{MaxEntries: capacity, Shards: 1})
		if err != nil {
			b.Fatal(err)
		}
		lru := &lruCache{capacity: capacity, ll: list.New(), items: map[int]*list.Element{}}
		lruHits := 0
		for _, k := range trace {
			if _, ok := c.Get(k); !ok {
				c.Set(k, k)
			}
			if lru.access(k) {
				lruHits++
			}
		}
		b.ReportMetric(100*c.Stats().HitRatio(), "sieve-hit%")
		b.ReportMetric(100*float64(lruHits)/accesses, "lru-hit%")
	}
}
