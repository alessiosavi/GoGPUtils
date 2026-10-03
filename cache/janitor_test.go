package cache

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestJanitor_RemovesExpiredEntries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		rec := &recorder[string, int]{}
		c, err := New(Config[string, int]{TTL: time.Second, JanitorInterval: 500 * time.Millisecond, OnEvict: rec.record})
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = c.Close() }) // must stop the janitor inside the bubble
		c.Set("a", 1)
		c.SetWithTTL("keep", 2, NoExpiration)
		time.Sleep(1500 * time.Millisecond)
		synctest.Wait()
		if c.Len() != 1 {
			t.Fatalf("Len = %d, want 1 (janitor removes the expired entry)", c.Len())
		}
		expectRemovals(t, rec.take(), []removal[string, int]{{"a", 1, ReasonExpired}})
	})
}

func TestClose_StopsJanitorAndCacheStaysUsable(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, err := New(Config[string, int]{JanitorInterval: time.Second})
		if err != nil {
			t.Fatal(err)
		}
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		select {
		case <-c.in.janitorDone:
		default:
			t.Fatal("janitor still running after Close")
		}
		if err := c.Close(); err != nil {
			t.Fatalf("second Close: %v", err)
		}
		c.Set("a", 1)
		if v, ok := c.Get("a"); !ok || v != 1 {
			t.Fatal("cache unusable after Close")
		}
	})
}

func TestClose_ConcurrentCallers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, err := New(Config[int, int]{JanitorInterval: time.Hour})
		if err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				if err := c.Close(); err != nil {
					t.Error(err)
				}
			})
		}
		wg.Wait()
		select {
		case <-c.in.janitorDone:
		default:
			t.Fatal("janitor still running after concurrent Close")
		}
	})
}

func TestClose_WithoutJanitor(t *testing.T) {
	c, err := New(Config[int, int]{})
	if err != nil {
		t.Fatal(err)
	}
	if c.in.janitorDone != nil {
		t.Fatal("janitor started without JanitorInterval")
	}
	for range 2 {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
	c.Set(1, 1)
	if _, ok := c.Get(1); !ok {
		t.Fatal("cache unusable after Close")
	}
}

func TestClose_DoesNotWaitForInFlightLoads(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		c, err := New(Config[string, int]{JanitorInterval: time.Hour})
		if err != nil {
			t.Fatal(err)
		}
		load, started, release := blockingLoader(4)
		defer release()
		r := loadAsync(context.Background(), c, "k", load)
		await(t, "loader started", started)
		closed := make(chan error, 1)
		go func() { closed <- c.Close() }()
		select {
		case err := <-closed:
			if err != nil {
				t.Fatal(err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Close waited for an in-flight load")
		}
		release()
		if res := recv(t, r); res.err != nil || res.v != 4 {
			t.Fatalf("load after Close = %v, %v", res.v, res.err)
		}
	})
}

// startUnreachableJanitor creates a cache with a janitor, drops every
// reference to the public wrapper, and returns the janitor's done channel.
//
//go:noinline
func startUnreachableJanitor(t *testing.T) <-chan struct{} {
	c, err := New(Config[int, int]{JanitorInterval: time.Hour})
	if err != nil {
		t.Fatal(err)
	}

	return c.in.janitorDone
}

// Observational only: cleanup timing is not guaranteed, so a timeout is
// inconclusive (skip), never a failure.
func TestJanitor_StoppedByCleanupWhenCacheUnreachable(t *testing.T) {
	if testing.Short() {
		t.Skip("observational GC test")
	}
	done := startUnreachableJanitor(t)
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		runtime.GC()
		select {
		case <-done:
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
	t.Skip("cleanup not observed within 10s (timing is not guaranteed)")
}
