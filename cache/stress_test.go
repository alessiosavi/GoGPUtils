package cache

import (
	"context"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestConcurrentStress(t *testing.T) {
	var removed atomic.Int64
	c, err := New(Config[int, int]{
		MaxEntries:      512,
		Shards:          4,
		TTL:             5 * time.Millisecond,
		JanitorInterval: time.Millisecond,
		OnEvict:         func(int, int, EvictionReason) { removed.Add(1) },
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })

	ops := 3000
	if testing.Short() {
		ops = 500
	}

	var wg sync.WaitGroup
	for g := range 8 {
		wg.Go(func() {
			rng := rand.New(rand.NewPCG(uint64(g), 7))
			ctx := context.Background()
			for range ops {
				k := rng.IntN(2048)
				switch rng.IntN(20) {
				case 0:
					c.Clear()
				case 1:
					for range c.All() {
					}
				case 2:
					c.DeleteExpired()
				case 3, 4:
					c.Delete(k)
				case 5, 6, 7:
					_, _ = c.GetOrLoad(ctx, k, func(context.Context) (int, error) { return k, nil })
				case 8, 9, 10, 11:
					c.SetWithTTL(k, k, time.Duration(rng.IntN(10))*time.Millisecond)
				default:
					c.Get(k)
				}
			}
		})
	}
	wg.Wait()

	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
	checkInvariants(t, c)
	if c.Len() > 512 {
		t.Fatalf("Len = %d exceeds MaxEntries", c.Len())
	}
}
