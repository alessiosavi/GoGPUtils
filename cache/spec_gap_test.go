package cache

import (
	"context"
	"errors"
	"maps"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestGetOrLoad_InvalidationGenerationOrders(t *testing.T) {
	for _, how := range []string{"delete", "clear"} {
		for _, oldFirst := range []bool{true, false} {
			t.Run(how+map[bool]string{true: "/old-first", false: "/new-first"}[oldFirst], func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					c, _ := newTestCache(t, Config[string, int]{Shards: 1})
					a, startedA, releaseA := blockingLoader(1)
					defer releaseA()
					b, startedB, releaseB := blockingLoader(2)
					defer releaseB()
					ra := loadAsync(context.Background(), c, "k", a)
					await(t, "generation A", startedA)
					if how == "delete" {
						c.Delete("k")
					} else {
						c.Clear()
					}
					rb := loadAsync(context.Background(), c, "k", b)
					await(t, "generation B", startedB)
					if oldFirst {
						releaseA()
						if got := recv(t, ra); got.err != nil || got.v != 1 {
							t.Fatalf("A = %+v", got)
						}
					}
					rc := loadAsync(context.Background(), c, "k", mustNotLoad[int](t))
					waitFor(t, "B joiner", func() bool { return c.Stats().Misses == 3 })
					releaseB()
					for _, ch := range []<-chan result[int]{rb, rc} {
						if got := recv(t, ch); got.err != nil || got.v != 2 {
							t.Fatalf("B = %+v", got)
						}
					}
					if !oldFirst {
						releaseA()
						if got := recv(t, ra); got.err != nil || got.v != 1 {
							t.Fatalf("A = %+v", got)
						}
					}
					if v, ok := c.Peek("k"); !ok || v != 2 {
						t.Fatalf("stored = %v, %v", v, ok)
					}
					if got := c.Stats(); got.Loads != 2 || got.LoadErrors != 0 {
						t.Fatalf("Stats = %+v", got)
					}
				})
			})
		}
	}
}

func TestGetOrLoad_ExpiredReplacementAndDefaultTTL(t *testing.T) {
	var c *Cache[string, int]
	hooks := testHooks{afterReadUnlock: func() { c.Set("k", 2) }}
	c, clk := newTestCacheWithHooks(t, Config[string, int]{Shards: 1}, hooks)
	c.SetWithTTL("k", 1, time.Second)
	clk.advance(time.Second)
	if v, err := c.GetOrLoad(context.Background(), "k", mustNotLoad[int](t)); err != nil || v != 2 {
		t.Fatalf("got %v, %v", v, err)
	}
	if got := c.Stats(); got != (Stats{Hits: 1, Expirations: 1}) {
		t.Fatalf("Stats = %+v", got)
	}
	d, dclk := newTestCache(t, Config[string, int]{Shards: 1, TTL: time.Second})
	if _, err := d.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { return 3, nil }); err != nil {
		t.Fatal(err)
	}
	if ttl, ok := d.TTL("k"); !ok || ttl != time.Second {
		t.Fatalf("TTL = %v, %v", ttl, ok)
	}
	dclk.advance(time.Second)
	if _, ok := d.Get("k"); ok {
		t.Fatal("load did not use the default TTL")
	}
}

func TestGetOrLoad_ExpirationCallbackMayReenter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c *Cache[string, int]
		var fired atomic.Bool
		callback := make(chan result[int], 1)
		cfg := Config[string, int]{Shards: 1, OnEvict: func(k string, _ int, reason EvictionReason) {
			if reason != ReasonExpired || !fired.CompareAndSwap(false, true) {
				return
			}
			c.Get("missing")
			c.Set("other", 10)
			v, err := c.GetOrLoad(context.Background(), k, mustNotLoad[int](t))
			callback <- result[int]{v, err}
		}}
		var clk *fakeClock
		c, clk = newTestCache(t, cfg)
		c.SetWithTTL("k", 1, time.Second)
		clk.advance(time.Second)
		r := loadAsync(context.Background(), c, "k", func(context.Context) (int, error) { return 2, nil })
		if got := recv(t, r); got.err != nil || got.v != 2 {
			t.Fatalf("load = %+v", got)
		}
		if got := recv(t, callback); got.err != nil || got.v != 2 {
			t.Fatalf("callback = %+v", got)
		}
	})
}

func TestGetOrLoad_ResultPublishedBeforeCompletionCallback(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		started, release := make(chan struct{}), make(chan struct{})
		defer close(release)
		c, _ := newTestCache(t, Config[string, int]{MaxEntries: 1, Shards: 1, OnEvict: func(string, int, EvictionReason) {
			close(started)
			<-release
		}})
		c.Set("a", 1)
		r := loadAsync(context.Background(), c, "b", func(context.Context) (int, error) { return 2, nil })
		await(t, "completion callback", started)
		if got := recv(t, r); got.err != nil || got.v != 2 {
			t.Fatalf("load = %+v", got)
		}
	})
}

func TestGetOrLoad_NonReflexiveErrorAndGoexit(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[float64, int]{Shards: 1})
		nan := math.NaN()
		boom := errors.New("boom")
		if v, err := c.GetOrLoad(context.Background(), nan, func(context.Context) (int, error) { return 99, boom }); v != 0 || !errors.Is(err, boom) {
			t.Fatalf("error result = %v, %v", v, err)
		}
		if v, err := c.GetOrLoad(context.Background(), nan, func(context.Context) (int, error) { runtime.Goexit(); return 99, nil }); v != 0 || !errors.Is(err, ErrLoaderGoexit) {
			t.Fatalf("Goexit result = %v, %v", v, err)
		}
		if v, err := c.GetOrLoad(context.Background(), nan, func(context.Context) (int, error) { panic(boom) }); v != 0 || !errors.Is(err, boom) {
			t.Fatalf("panic result = %v, %v", v, err)
		}
		a, startedA, releaseA := blockingLoader(1)
		defer releaseA()
		b, startedB, releaseB := blockingLoader(2)
		defer releaseB()
		ra := loadAsync(context.Background(), c, nan, a)
		await(t, "NaN load A", startedA)
		rb := loadAsync(context.Background(), c, nan, b)
		await(t, "independent NaN load B", startedB)
		releaseA()
		releaseB()
		if got := recv(t, ra); got.v != 1 || got.err != nil {
			t.Fatalf("A = %+v", got)
		}
		if got := recv(t, rb); got.v != 2 || got.err != nil {
			t.Fatalf("B = %+v", got)
		}
		if c.Len() != 0 {
			t.Fatal("NaN was stored")
		}
		if got := c.Stats(); got != (Stats{Misses: 5, Loads: 5, LoadErrors: 3}) {
			t.Fatalf("Stats = %+v", got)
		}
	})
}

func TestRequestStop_IsIdempotentAndStopsJanitor(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, err := New(Config[string, int]{JanitorInterval: time.Hour})
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = c.Close() }()
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(c.in.requestStop)
		}
		wg.Wait()
		await(t, "requestStop janitor exit", c.in.janitorDone)
	})
}

func TestAll_SecondInvocationSeesMutation(t *testing.T) {
	c, _ := newTestCache(t, Config[string, int]{Shards: 1})
	c.Set("a", 1)
	seq := c.All()
	if got := maps.Collect(seq); !maps.Equal(got, map[string]int{"a": 1}) {
		t.Fatalf("first = %v", got)
	}
	c.Clear()
	c.Set("b", 2)
	if got := maps.Collect(seq); !maps.Equal(got, map[string]int{"b": 2}) {
		t.Fatalf("second = %v", got)
	}
}

func checkNonReflexiveKey[K comparable](t *testing.T, key K) {
	t.Helper()
	c, _ := newTestCache(t, Config[K, int]{Shards: 1})
	c.Set(key, 1)
	c.SetWithTTL(key, 2, time.Second)
	if _, ok := c.Get(key); ok {
		t.Fatal("Get hit")
	}
	if _, ok := c.Peek(key); ok {
		t.Fatal("Peek hit")
	}
	if ttl, ok := c.TTL(key); ok || ttl != 0 {
		t.Fatal("TTL hit")
	}
	if c.Delete(key) || c.Len() != 0 {
		t.Fatal("non-reflexive key stored")
	}
}

func TestNonReflexiveKeyMatrix(t *testing.T) {
	nan := math.NaN()
	checkNonReflexiveKey(t, nan)
	checkNonReflexiveKey(t, float32(nan))
	checkNonReflexiveKey(t, complex(nan, 0))
	checkNonReflexiveKey(t, complex(0, nan))
	checkNonReflexiveKey(t, struct{ F float64 }{nan})
	checkNonReflexiveKey(t, [2]float64{1, nan})
	checkNonReflexiveKey(t, any(nan))
	checkNonReflexiveKey(t, any([2]any{nan, []int{1}}))
}

func TestGetOrLoad_ValidationDoesNotTouchEntries(t *testing.T) {
	done, cancel := context.WithCancel(context.Background())
	cancel()
	for _, expired := range []bool{false, true} {
		rec := &recorder[string, int]{}
		c, clk := newTestCache(t, Config[string, int]{Shards: 1, OnEvict: rec.record})
		c.SetWithTTL("k", 1, time.Second)
		if expired {
			clk.advance(time.Second)
		}
		var nilCtx context.Context
		tests := []struct {
			ctx  context.Context
			load func(context.Context) (int, error)
			want error
		}{
			{nilCtx, nil, ErrNilContext},
			{done, nil, ErrNilLoader},
			{done, mustNotLoad[int](t), context.Canceled},
		}
		for _, tt := range tests {
			if _, err := c.GetOrLoad(tt.ctx, "k", tt.load); !errors.Is(err, tt.want) {
				t.Fatalf("validation = %v, want %v", err, tt.want)
			}
		}
		if c.Len() != 1 || c.Stats() != (Stats{}) || c.in.shards[0].items["k"].visited.Load() {
			t.Fatal("validation affected the entry or statistics")
		}
		expectRemovals(t, rec.take(), nil)
	}
}
