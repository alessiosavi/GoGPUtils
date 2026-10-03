package cache

import (
	"container/list"
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"testing"
)

const benchKeys = 1 << 16

func benchCache(b *testing.B, shards int) *Cache[int, int] {
	b.Helper()

	c, err := New(Config[int, int]{MaxEntries: benchKeys, Shards: shards})
	if err != nil {
		b.Fatal(err)
	}
	for i := range benchKeys {
		c.Set(i, i)
	}

	return c
}

func BenchmarkGetHitParallel(b *testing.B) {
	for _, shards := range []int{1, 0} {
		b.Run(fmt.Sprintf("shards=%d", shards), func(b *testing.B) {
			c := benchCache(b, shards)
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					c.Get(i & (benchKeys - 1))
					i++
				}
			})
		})
	}
}

func BenchmarkSet(b *testing.B) {
	c := benchCache(b, 0)
	i := 0
	for b.Loop() {
		c.Set(i&(2*benchKeys-1), i)
		i++
	}
}

func BenchmarkMixedParallel(b *testing.B) {
	for _, readPct := range []int{90, 50} {
		b.Run(fmt.Sprintf("read=%d%%", readPct), func(b *testing.B) {
			c := benchCache(b, 0)
			b.RunParallel(func(pb *testing.PB) {
				rng := rand.New(rand.NewPCG(rand.Uint64(), 1))
				for pb.Next() {
					k := rng.IntN(2 * benchKeys)
					if rng.IntN(100) < readPct {
						c.Get(k)
					} else {
						c.Set(k, k)
					}
				}
			})
		})
	}
}

func BenchmarkGetOrLoadHit(b *testing.B) {
	c := benchCache(b, 0)
	ctx := context.Background()
	load := func(context.Context) (int, error) { return 0, nil }
	i := 0
	for b.Loop() {
		_, _ = c.GetOrLoad(ctx, i&(benchKeys-1), load)
		i++
	}
}

func BenchmarkBaselineRWMutexMapParallel(b *testing.B) {
	var mu sync.RWMutex
	m := make(map[int]int, benchKeys)
	for i := range benchKeys {
		m[i] = i
	}
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mu.RLock()
			_ = m[i&(benchKeys-1)]
			mu.RUnlock()
			i++
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

// BenchmarkHitRatioZipf reports SIEVE and LRU hit ratios on the same Zipf
// trace (informational, not asserted).
func BenchmarkHitRatioZipf(b *testing.B) {
	const capacity, accesses = 1000, 200_000
	zipf := rand.NewZipf(rand.New(rand.NewPCG(1, 2)), 1.07, 1, 100_000)
	trace := make([]int, accesses)
	for i := range trace {
		trace[i] = int(zipf.Uint64())
	}
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
