package cache

import (
	"math"
	"sync"
	"testing"
	"time"
)

func TestTTL_DefaultOverrideAndNoExpiration(t *testing.T) {
	c, clk := newTestCache(t, Config[string, int]{TTL: 10 * time.Second})
	c.Set("default", 1)
	c.SetWithTTL("short", 2, time.Second)
	c.SetWithTTL("never", 3, NoExpiration)
	c.SetWithTTL("negative", 4, -5*time.Second)

	assertTTL := func(key string, want time.Duration, wantOK bool) {
		t.Helper()
		if got, ok := c.TTL(key); got != want || ok != wantOK {
			t.Fatalf("TTL(%q) = %v, %v; want %v, %v", key, got, ok, want, wantOK)
		}
	}
	assertTTL("default", 10*time.Second, true)
	assertTTL("short", time.Second, true)
	assertTTL("never", NoExpiration, true)
	assertTTL("negative", NoExpiration, true)
	assertTTL("missing", 0, false)

	clk.advance(time.Second) // now == expiresAt counts as expired
	assertTTL("short", 0, false)
	if _, ok := c.Get("short"); ok {
		t.Fatal("Get returned an entry at its deadline")
	}
	assertTTL("default", 9*time.Second, true)
}

func TestTTL_ZeroDefaultMeansNever(t *testing.T) {
	c, clk := newTestCache(t, Config[string, int]{})
	c.Set("a", 1)
	clk.advance(1000 * time.Hour)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get = %v, %v; want 1, true", v, ok)
	}
	if d, ok := c.TTL("a"); d != NoExpiration || !ok {
		t.Fatalf("TTL = %v, %v; want NoExpiration, true", d, ok)
	}
}

func TestTTL_SaturatesAtMaxDuration(t *testing.T) {
	c, clk := newTestCache(t, Config[string, int]{Shards: 1})
	clk.advance(time.Hour)
	c.SetWithTTL("a", 1, time.Duration(math.MaxInt64))
	if got := c.in.shards[0].items["a"].expiresAt; got != math.MaxInt64 {
		t.Fatalf("expiresAt = %d, want math.MaxInt64", got)
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("entry with a saturated deadline expired")
	}
	if d, ok := c.TTL("a"); !ok || d <= 0 {
		t.Fatalf("TTL = %v, %v; want positive, true", d, ok)
	}
}

func TestExpiredEntry_PerMethodSemantics(t *testing.T) {
	setup := func(t *testing.T) (*Cache[string, int], *recorder[string, int]) {
		rec := &recorder[string, int]{}
		c, clk := newTestCache(t, Config[string, int]{Shards: 1, OnEvict: rec.record})
		c.SetWithTTL("k", 1, time.Second)
		clk.advance(2 * time.Second)

		return c, rec
	}
	expired := []removal[string, int]{{"k", 1, ReasonExpired}}

	t.Run("Peek keeps it", func(t *testing.T) {
		c, rec := setup(t)
		if _, ok := c.Peek("k"); ok || c.Len() != 1 {
			t.Fatal("Peek returned or removed an expired entry")
		}
		expectRemovals(t, rec.take(), nil)
	})
	t.Run("TTL keeps it", func(t *testing.T) {
		c, rec := setup(t)
		if d, ok := c.TTL("k"); d != 0 || ok || c.Len() != 1 {
			t.Fatal("TTL returned or removed an expired entry")
		}
		expectRemovals(t, rec.take(), nil)
	})
	t.Run("Get removes it", func(t *testing.T) {
		c, rec := setup(t)
		if _, ok := c.Get("k"); ok || c.Len() != 0 {
			t.Fatal("Get did not remove an expired entry")
		}
		expectRemovals(t, rec.take(), expired)
	})
	t.Run("Set replaces it as expired", func(t *testing.T) {
		c, rec := setup(t)
		c.Set("k", 2)
		expectRemovals(t, rec.take(), expired)
		if v, ok := c.Peek("k"); !ok || v != 2 {
			t.Fatalf("Peek = %v, %v; want 2, true", v, ok)
		}
	})
	t.Run("Delete reports false", func(t *testing.T) {
		c, rec := setup(t)
		if c.Delete("k") || c.Len() != 0 {
			t.Fatal("Delete of an expired entry returned true or kept it")
		}
		expectRemovals(t, rec.take(), expired)
	})
	t.Run("DeleteExpired removes it", func(t *testing.T) {
		c, rec := setup(t)
		if n := c.DeleteExpired(); n != 1 {
			t.Fatalf("DeleteExpired = %d, want 1", n)
		}
		expectRemovals(t, rec.take(), expired)
	})
	t.Run("Clear reports deleted", func(t *testing.T) {
		c, rec := setup(t)
		c.Clear()
		expectRemovals(t, rec.take(), []removal[string, int]{{"k", 1, ReasonDeleted}})
	})
}

func TestEviction_ExpiredVisitedEntryIsVictim(t *testing.T) {
	rec := &recorder[string, int]{}
	c, clk := newTestCache(t, Config[string, int]{MaxEntries: 2, Shards: 1, OnEvict: rec.record})
	c.SetWithTTL("a", 1, time.Second)
	c.Set("b", 2)
	c.Get("a")
	c.Get("b") // both visited
	clk.advance(2 * time.Second)
	c.Set("c", 3)
	expectRemovals(t, rec.take(), []removal[string, int]{{"a", 1, ReasonExpired}})
	if _, ok := c.Peek("b"); !ok {
		t.Fatal("live entry evicted instead of the expired one")
	}
	checkInvariants(t, c)
}

func TestDeleteExpired_CountsAcrossShards(t *testing.T) {
	c, clk := newTestCache(t, Config[int, int]{Shards: 4})
	for i := range 100 {
		if i%2 == 0 {
			c.SetWithTTL(i, i, time.Second)
		} else {
			c.Set(i, i)
		}
	}
	clk.advance(time.Second)
	if n := c.DeleteExpired(); n != 50 {
		t.Fatalf("DeleteExpired = %d, want 50", n)
	}
	if c.Len() != 50 {
		t.Fatalf("Len = %d, want 50", c.Len())
	}
	checkInvariants(t, c)
}

func TestGet_ExpiredEntryReplacedDuringLockUpgrade(t *testing.T) {
	var c *Cache[string, int]
	var once sync.Once
	hooks := testHooks{afterReadUnlock: func() { once.Do(func() { c.Set("k", 2) }) }}
	c, clk := newTestCacheWithHooks(t, Config[string, int]{Shards: 1}, hooks)
	c.SetWithTTL("k", 1, time.Second)
	clk.advance(2 * time.Second)
	if v, ok := c.Get("k"); !ok || v != 2 {
		t.Fatalf("Get = %v, %v; want the live replacement 2, true", v, ok)
	}
}

func TestTTL_NonReflexiveKey(t *testing.T) {
	c, _ := newTestCache(t, Config[float64, int]{})
	nan := math.NaN()
	c.SetWithTTL(nan, 1, time.Second)
	if d, ok := c.TTL(nan); d != 0 || ok || c.Len() != 0 {
		t.Fatal("NaN key stored or reported by TTL")
	}
}
