package cache

import (
	"context"
	"errors"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestGetOrLoad_HitAndMiss(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		calls := 0
		load := func(context.Context) (int, error) { calls++; return 7, nil }
		for range 2 {
			if v, err := c.GetOrLoad(context.Background(), "k", load); err != nil || v != 7 {
				t.Fatalf("GetOrLoad = %v, %v; want 7, nil", v, err)
			}
		}
		if calls != 1 {
			t.Fatalf("loader ran %d times, want 1", calls)
		}
		if st := c.Stats(); st.Hits != 1 || st.Misses != 1 || st.Loads != 1 {
			t.Fatalf("Stats = %+v", st)
		}
	})
}

func TestGetOrLoad_DeduplicatesConcurrentMisses(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		var loads atomic.Int64
		release := make(chan struct{})
		releaseLoad := sync.OnceFunc(func() { close(release) })
		defer releaseLoad()
		load := func(context.Context) (int, error) { //nolint:unparam // GetOrLoad requires an error result even for a success-only loader.
			loads.Add(1)
			<-release

			return 42, nil
		}
		const callers = 100
		results := make(chan result[int], callers)
		for range callers {
			go func() {
				v, err := c.GetOrLoad(context.Background(), "k", load)
				results <- result[int]{v, err}
			}()
		}
		waitFor(t, "all callers to miss", func() bool { return c.Stats().Misses == callers })
		releaseLoad()
		for range callers {
			if r := recv[int](t, results); r.err != nil || r.v != 42 {
				t.Fatalf("caller got %v, %v", r.v, r.err)
			}
		}
		if loads.Load() != 1 {
			t.Fatalf("loader ran %d times, want 1", loads.Load())
		}
		if st := c.Stats(); st.Loads != 1 || st.Misses != callers {
			t.Fatalf("Stats = %+v", st)
		}
	})
}

func TestGetOrLoad_CancelledWaiterDoesNotCancelLoad(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		load, started, release := blockingLoader(5)
		defer release()
		ctx1, cancel1 := context.WithCancel(context.Background())
		r1 := loadAsync(ctx1, c, "k", load)
		await(t, "loader started", started)
		r2 := loadAsync(context.Background(), c, "k", load)
		waitFor(t, "second caller to join", func() bool { return c.Stats().Misses == 2 })
		cancel1()
		if r := recv(t, r1); !errors.Is(r.err, context.Canceled) {
			t.Fatalf("cancelled caller got %v, %v", r.v, r.err)
		}
		release()
		if r := recv(t, r2); r.err != nil || r.v != 5 {
			t.Fatalf("remaining caller got %v, %v", r.v, r.err)
		}
		if st := c.Stats(); st.Loads != 1 || st.LoadErrors != 0 {
			t.Fatalf("Stats = %+v", st)
		}
	})
}

func TestGetOrLoad_LoadCompletesAfterAllWaitersCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		load, started, release := blockingLoader(9)
		defer release()
		ctx, cancel := context.WithCancel(context.Background())
		r := loadAsync(ctx, c, "k", load)
		await(t, "loader started", started)
		cancel()
		if res := recv(t, r); !errors.Is(res.err, context.Canceled) {
			t.Fatalf("got %v, %v", res.v, res.err)
		}
		release()
		waitFor(t, "the load to store its value", func() bool {
			v, ok := c.Peek("k")

			return ok && v == 9
		})
	})
}

func TestGetOrLoad_LoaderGetsDetachedContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		type ctxKey struct{}
		c, _ := newTestCache(t, Config[string, string]{Shards: 1})
		ctx, cancel := context.WithCancel(context.WithValue(context.Background(), ctxKey{}, "v"))
		defer cancel()
		observed := make(chan error, 1)
		_, _ = c.GetOrLoad(ctx, "k", func(lctx context.Context) (string, error) {
			cancel() // the caller's cancellation must not reach the loader
			switch {
			case lctx.Err() != nil:
				observed <- lctx.Err()
			case lctx.Value(ctxKey{}) != "v":
				observed <- errors.New("context value lost")
			default:
				observed <- nil
			}

			return "ok", nil
		})
		if err := await(t, "detached context observation", observed); err != nil {
			t.Fatalf("loader context: %v", err)
		}
	})
}

func TestGetOrLoad_ErrorsAreNotCached(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		boom := errors.New("boom")
		v, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { return 99, boom })
		if !errors.Is(err, boom) || v != 0 {
			t.Fatalf("GetOrLoad = %v, %v; want 0, boom", v, err)
		}
		if _, ok := c.Peek("k"); ok {
			t.Fatal("error result was cached")
		}
		if v, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { return 1, nil }); err != nil || v != 1 {
			t.Fatalf("retry = %v, %v", v, err)
		}
		if st := c.Stats(); st.Loads != 2 || st.LoadErrors != 1 {
			t.Fatalf("Stats = %+v", st)
		}
	})
}

func TestGetOrLoad_PanicBecomesPanicError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		boom := errors.New("boom")
		_, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { panic(boom) })
		var pe *PanicError
		if !errors.As(err, &pe) || !errors.Is(err, boom) || len(pe.Stack) == 0 {
			t.Fatalf("err = %v, want *PanicError wrapping boom with a stack", err)
		}
		if st := c.Stats(); st.LoadErrors != 1 {
			t.Fatalf("Stats = %+v", st)
		}
	})
}

func TestGetOrLoad_GoexitBecomesError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		_, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) {
			runtime.Goexit()

			return 0, nil
		})
		if !errors.Is(err, ErrLoaderGoexit) {
			t.Fatalf("err = %v, want ErrLoaderGoexit", err)
		}
		if v, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { return 3, nil }); err != nil || v != 3 {
			t.Fatalf("key stuck in flight after Goexit: %v, %v", v, err)
		}
	})
}

func TestGetOrLoad_Validation(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		load := mustNotLoad[int](t)
		var nilCtx context.Context
		done, cancel := context.WithCancel(context.Background())
		cancel()

		tests := []struct {
			name string
			ctx  context.Context
			load func(context.Context) (int, error)
			want error
		}{
			{"nil context", nilCtx, load, ErrNilContext},
			{"nil context before nil loader", nilCtx, nil, ErrNilContext},
			{"nil loader", context.Background(), nil, ErrNilLoader},
			{"nil loader before done context", done, nil, ErrNilLoader},
			{"done context", done, load, context.Canceled},
		}
		for _, tt := range tests {
			if _, err := c.GetOrLoad(tt.ctx, "k", tt.load); !errors.Is(err, tt.want) {
				t.Errorf("%s: err = %v, want %v", tt.name, err, tt.want)
			}
		}
		if st := c.Stats(); st != (Stats{}) {
			t.Fatalf("validation failures touched stats: %+v", st)
		}
	})
}

func TestGetOrLoad_SetDuringLoadWins(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		load, started, release := blockingLoader(1)
		defer release()
		r1 := loadAsync(context.Background(), c, "k", load)
		await(t, "loader started", started)
		c.Set("k", 5)
		if v, err := c.GetOrLoad(context.Background(), "k", mustNotLoad[int](t)); err != nil || v != 5 {
			t.Fatalf("caller after Set got %v, %v; want the live value 5", v, err)
		}
		release()
		if r := recv(t, r1); r.err != nil || r.v != 1 {
			t.Fatalf("existing waiter got %v, %v; want the old call's 1", r.v, r.err)
		}
		if v, _ := c.Get("k"); v != 5 {
			t.Fatalf("stale load overwrote the newer Set: %v", v)
		}
	})
}

func TestGetOrLoad_DeleteDuringLoadStartsNewGeneration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		loadA, startedA, releaseA := blockingLoader(1)
		defer releaseA()
		loadB, startedB, releaseB := blockingLoader(2)
		defer releaseB()
		r1 := loadAsync(context.Background(), c, "k", loadA)
		await(t, "loader started", startedA)
		c.Delete("k")
		r2 := loadAsync(context.Background(), c, "k", loadB)
		await(t, "loader B started", startedB) // a new generation started instead of joining A
		releaseA()
		if r := recv(t, r1); r.err != nil || r.v != 1 {
			t.Fatalf("generation A waiter got %v, %v", r.v, r.err)
		}
		r3 := loadAsync(context.Background(), c, "k", mustNotLoad[int](t))
		waitFor(t, "third caller to join generation B", func() bool { return c.Stats().Misses == 3 })
		releaseB()
		for _, ch := range []<-chan result[int]{r2, r3} {
			if r := recv(t, ch); r.err != nil || r.v != 2 {
				t.Fatalf("generation B waiter got %v, %v", r.v, r.err)
			}
		}
		if v, ok := c.Get("k"); !ok || v != 2 {
			t.Fatalf("Get = %v, %v; want generation B's 2", v, ok)
		}
		if st := c.Stats(); st.Loads != 2 {
			t.Fatalf("Loads = %d, want 2", st.Loads)
		}
	})
}

func TestGetOrLoad_ClearDuringLoadPreventsWriteBack(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[string, int]{Shards: 1})
		load, started, release := blockingLoader(1)
		defer release()
		r := loadAsync(context.Background(), c, "k", load)
		await(t, "loader started", started)
		c.Clear()
		release()
		if res := recv(t, r); res.err != nil || res.v != 1 {
			t.Fatalf("waiter got %v, %v", res.v, res.err)
		}
		if _, ok := c.Peek("k"); ok {
			t.Fatal("a load that finished after Clear resurrected the key")
		}
	})
}

func TestGetOrLoad_HitPathsMarkVisited(t *testing.T) {
	t.Run("fast path", func(t *testing.T) {
		rec := &recorder[string, int]{}
		c, _ := newTestCache(t, Config[string, int]{MaxEntries: 2, Shards: 1, OnEvict: rec.record})
		c.Set("a", 1)
		c.Set("b", 2)
		if v, err := c.GetOrLoad(context.Background(), "a", mustNotLoad[int](t)); err != nil || v != 1 {
			t.Fatalf("GetOrLoad = %v, %v", v, err)
		}
		c.Set("c", 3) // "a" is visited, so "b" is evicted
		expectRemovals(t, rec.take(), []removal[string, int]{{"b", 2, ReasonEvicted}})
	})
	t.Run("locked path", func(t *testing.T) {
		rec := &recorder[string, int]{}
		var c *Cache[string, int]
		var once sync.Once
		hooks := testHooks{afterReadUnlock: func() { once.Do(func() { c.Set("a", 1) }) }}
		c, _ = newTestCacheWithHooks(t, Config[string, int]{MaxEntries: 2, Shards: 1, OnEvict: rec.record}, hooks)
		c.Set("b", 2)
		// The fast path misses, the hook inserts "a", the locked re-check hits.
		if v, err := c.GetOrLoad(context.Background(), "a", mustNotLoad[int](t)); err != nil || v != 1 {
			t.Fatalf("GetOrLoad = %v, %v", v, err)
		}
		if st := c.Stats(); st.Hits != 1 || st.Misses != 0 || st.Loads != 0 {
			t.Fatalf("Stats = %+v", st)
		}
		c.in.shards[0].hand = c.in.shards[0].items["a"]
		c.Set("c", 3) // a must be skipped; an unmarked a would be evicted
		expectRemovals(t, rec.take(), []removal[string, int]{{"b", 2, ReasonEvicted}})
	})
}

func TestGetOrLoad_ReentrantFromOnEvict(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var c *Cache[string, int]
		var fired atomic.Bool
		reloaded := make(chan result[int], 1)
		cfg := Config[string, int]{MaxEntries: 1, Shards: 1, OnEvict: func(k string, _ int, r EvictionReason) {
			if r != ReasonEvicted || !fired.CompareAndSwap(false, true) {
				return
			}
			v, err := c.GetOrLoad(context.Background(), k, func(context.Context) (int, error) { return 100, nil })
			reloaded <- result[int]{v, err}
		}}
		c, _ = newTestCache(t, cfg)
		c.Set("a", 1)

		done := make(chan struct{})
		go func() {
			defer close(done)
			c.Set("b", 2) // evicts "a"; the callback reloads "a", which evicts "b"
		}()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("GetOrLoad from OnEvict deadlocked")
		}
		if r := recv(t, reloaded); r.err != nil || r.v != 100 {
			t.Fatalf("reload = %v, %v", r.v, r.err)
		}
		checkInvariants(t, c)
	})
}

func TestGetOrLoad_NonReflexiveKeyLoadsWithoutCaching(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, _ := newTestCache(t, Config[float64, int]{Shards: 1})
		nan := math.NaN()
		var calls atomic.Int64
		load := func(context.Context) (int, error) { return int(calls.Add(1)), nil }
		for want := 1; want <= 2; want++ {
			if v, err := c.GetOrLoad(context.Background(), nan, load); err != nil || v != want {
				t.Fatalf("GetOrLoad = %v, %v; want %d (no deduplication or caching)", v, err, want)
			}
		}
		if c.Len() != 0 {
			t.Fatal("NaN key was stored")
		}
		if st := c.Stats(); st.Loads != 2 || st.Misses != 2 {
			t.Fatalf("Stats = %+v", st)
		}

		boom := errors.New("boom")
		_, err := c.GetOrLoad(context.Background(), nan, func(context.Context) (int, error) { panic(boom) })
		if _, ok := errors.AsType[*PanicError](err); !ok {
			t.Fatalf("panic not converted for a NaN key: %v", err)
		}

		blocking, started, release := blockingLoader(7)

		defer release()
		ctx, cancel := context.WithCancel(context.Background())
		r := loadAsync(ctx, c, nan, blocking)
		await(t, "loader started", started)
		cancel()
		if res := recv(t, r); !errors.Is(res.err, context.Canceled) {
			t.Fatalf("cancelled NaN-key caller got %v, %v", res.v, res.err)
		}
		release()
	})
}
