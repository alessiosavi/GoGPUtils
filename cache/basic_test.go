package cache

import (
	"math"
	"testing"
)

func TestSetGetPeekDelete(t *testing.T) {
	c, _ := newTestCache(t, Config[string, int]{})
	if _, ok := c.Get("a"); ok {
		t.Fatal("Get on an empty cache hit")
	}
	c.Set("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get = %v, %v; want 1, true", v, ok)
	}
	if v, ok := c.Peek("a"); !ok || v != 1 {
		t.Fatalf("Peek = %v, %v; want 1, true", v, ok)
	}
	c.Set("a", 2)
	if v, _ := c.Get("a"); v != 2 {
		t.Fatalf("Get after replace = %v, want 2", v)
	}
	if c.Len() != 1 {
		t.Fatalf("Len = %d, want 1", c.Len())
	}
	if !c.Delete("a") {
		t.Fatal("Delete of a live key returned false")
	}
	if c.Delete("a") {
		t.Fatal("Delete of a missing key returned true")
	}
	if _, ok := c.Peek("a"); ok || c.Len() != 0 {
		t.Fatal("key still present after Delete")
	}
	checkInvariants(t, c)
}

func TestSet_ReplaceKeepsPositionAndMarksVisited(t *testing.T) {
	rec := &recorder[string, int]{}
	c, _ := newTestCache(t, Config[string, int]{MaxEntries: 3, Shards: 1, OnEvict: rec.record})
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("a", 10) // replaced in place at the tail and marked visited
	expectRemovals(t, rec.take(), []removal[string, int]{{"a", 1, ReasonReplaced}})
	c.Set("d", 4) // hand: a visited -> cleared; b unvisited -> evicted
	expectRemovals(t, rec.take(), []removal[string, int]{{"b", 2, ReasonEvicted}})
	if v, ok := c.Peek("a"); !ok || v != 10 {
		t.Fatalf("Peek(a) = %v, %v; want 10, true", v, ok)
	}
	checkInvariants(t, c)
}

func TestClear(t *testing.T) {
	rec := &recorder[int, int]{}
	c, _ := newTestCache(t, Config[int, int]{MaxEntries: 8, Shards: 1, OnEvict: rec.record})
	for i := range 5 {
		c.Set(i, i)
	}
	c.Clear()
	got := rec.take()
	if len(got) != 5 {
		t.Fatalf("Clear reported %d removals, want 5", len(got))
	}
	for _, r := range got {
		if r.reason != ReasonDeleted {
			t.Fatalf("Clear removal %v, want ReasonDeleted", r)
		}
	}
	if c.Len() != 0 {
		t.Fatalf("Len after Clear = %d", c.Len())
	}
	c.Set(9, 9)
	if v, ok := c.Get(9); !ok || v != 9 {
		t.Fatal("cache unusable after Clear")
	}
	checkInvariants(t, c)
}

func TestRemovedEntriesReleaseValues(t *testing.T) {
	c, _ := newTestCache(t, Config[string, *int]{MaxEntries: 1, Shards: 1})
	assertReleased := func(e *entry[string, *int], how string) {
		t.Helper()
		if e.value != nil || e.prev != nil || e.next != nil {
			t.Fatalf("entry removed by %s still references its value or neighbours", how)
		}
	}
	c.Set("a", new(int))
	e := c.in.shards[0].items["a"]
	c.Set("b", new(int)) // evicts "a"
	assertReleased(e, "eviction")
	e = c.in.shards[0].items["b"]
	c.Delete("b")
	assertReleased(e, "Delete")
	c.Set("c", new(int))
	e = c.in.shards[0].items["c"]
	c.Clear()
	assertReleased(e, "Clear")
}

func TestBoundedLenAcrossShards(t *testing.T) {
	c, _ := newTestCache(t, Config[int, int]{MaxEntries: 4096, Shards: 8})
	for i := range 20000 {
		c.Set(i, i)
	}
	if got := c.Len(); got != 4096 {
		t.Fatalf("Len = %d, want 4096", got)
	}
	checkInvariants(t, c)

	tiny, _ := newTestCache(t, Config[int, int]{MaxEntries: 1, Shards: 64})
	for i := range 10 {
		tiny.Set(i, i)
	}
	if tiny.Len() != 1 {
		t.Fatalf("Len = %d, want 1", tiny.Len())
	}
}

func TestNonReflexiveKeys_NotStored(t *testing.T) {
	nan := math.NaN()

	t.Run("float64", func(t *testing.T) {
		c, _ := newTestCache(t, Config[float64, int]{})
		c.Set(nan, 1)
		if c.Len() != 0 {
			t.Fatal("NaN key was stored")
		}
		if _, ok := c.Get(nan); ok {
			t.Fatal("Get(NaN) hit")
		}
		if _, ok := c.Peek(nan); ok {
			t.Fatal("Peek(NaN) hit")
		}
		if c.Delete(nan) {
			t.Fatal("Delete(NaN) returned true")
		}
	})
	t.Run("struct with NaN", func(t *testing.T) {
		type key struct{ f float32 }
		c, _ := newTestCache(t, Config[key, int]{})
		c.Set(key{float32(nan)}, 1)
		if c.Len() != 0 {
			t.Fatal("struct key containing NaN was stored")
		}
	})
	t.Run("array with NaN", func(t *testing.T) {
		c, _ := newTestCache(t, Config[[2]float64, int]{})
		c.Set([2]float64{1, nan}, 1)
		if c.Len() != 0 {
			t.Fatal("array key containing NaN was stored")
		}
	})
	t.Run("complex with NaN", func(t *testing.T) {
		c, _ := newTestCache(t, Config[complex128, int]{})
		c.Set(complex(nan, 0), 1)
		if c.Len() != 0 {
			t.Fatal("complex NaN key was stored")
		}
	})
	t.Run("interface with NaN", func(t *testing.T) {
		c, _ := newTestCache(t, Config[any, int]{})
		c.Set(any(nan), 1)
		c.Set(any([2]any{nan, []int{1}}), 1) // comparison stops at the NaN: no panic
		if c.Len() != 0 {
			t.Fatal("interface key containing NaN was stored")
		}
	})
	t.Run("interface with non-comparable value panics", func(t *testing.T) {
		c, _ := newTestCache(t, Config[any, int]{})
		defer func() {
			if recover() == nil {
				t.Fatal("expected a panic for a non-comparable interface key")
			}
		}()
		c.Set([]int{1}, 1)
	})
}
