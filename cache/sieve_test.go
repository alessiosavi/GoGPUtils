package cache

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// refSieve is a direct, slice-based transcription of SIEVE (NSDI'24,
// Algorithm 1) extended with Delete and TTL-aware eviction, used as a test
// oracle (Task 2 fidelity tests and the Task 7 model test).
type refSieve struct {
	capacity int
	keys     []int        // index 0 = tail (oldest), last = head (newest)
	visited  map[int]bool // presence = cached
	hand     int          // index into keys; -1 = start at the tail
}

func newRefSieve(capacity int) *refSieve {
	return &refSieve{capacity: capacity, visited: map[int]bool{}, hand: -1}
}

func (r *refSieve) has(k int) bool {
	_, ok := r.visited[k]

	return ok
}

func (r *refSieve) get(k int) bool {
	if !r.has(k) {
		return false
	}
	r.visited[k] = true

	return true
}

// set inserts or refreshes k and returns the evicted key, if any.
func (r *refSieve) set(k int) (int, bool) {
	return r.setWith(k, func(int) bool { return false })
}

// setWith is set with TTL awareness: the hand evicts a key for which expired
// returns true even if it is marked visited. capacity 0 means unbounded.
func (r *refSieve) setWith(k int, expired func(int) bool) (int, bool) {
	if r.has(k) {
		r.visited[k] = true

		return 0, false
	}

	victim, evicted := 0, false
	if r.capacity > 0 && len(r.keys) == r.capacity {
		i := max(r.hand, 0)
		for !expired(r.keys[i]) && r.visited[r.keys[i]] {
			r.visited[r.keys[i]] = false
			if i++; i == len(r.keys) {
				i = 0
			}
		}
		victim, evicted = r.keys[i], true
		r.removeAt(i)
		// The entry that was toward the head of the victim now sits at i.
		if i < len(r.keys) {
			r.hand = i
		} else {
			r.hand = -1
		}
	}
	r.keys = append(r.keys, k)
	r.visited[k] = false

	return victim, evicted
}

func (r *refSieve) del(k int) bool {
	if !r.has(k) {
		return false
	}
	i := slices.Index(r.keys, k)
	r.removeAt(i)
	switch {
	case r.hand == i && i >= len(r.keys):
		r.hand = -1 // the hand was at the head: wrap to the tail
	case r.hand > i:
		r.hand--
	}

	return true
}

func (r *refSieve) removeAt(i int) {
	delete(r.visited, r.keys[i])
	r.keys = slices.Delete(r.keys, i, i+1)
}

func TestSieve_ScriptedTrace(t *testing.T) {
	rec := &recorder[string, int]{}
	c, _ := newTestCache(t, Config[string, int]{MaxEntries: 3, Shards: 1, OnEvict: rec.record})
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Get("a")    // a visited
	c.Set("d", 4) // hand: a cleared, b evicted; hand -> c
	expectRemovals(t, rec.take(), []removal[string, int]{{"b", 2, ReasonEvicted}})
	c.Delete("c") // hand pointed at c -> moves toward the head, to d
	expectRemovals(t, rec.take(), []removal[string, int]{{"c", 3, ReasonDeleted}})
	c.Set("e", 5) // room available: no eviction
	expectRemovals(t, rec.take(), nil)
	c.Set("f", 6) // hand at d (unvisited) -> d evicted
	expectRemovals(t, rec.take(), []removal[string, int]{{"d", 4, ReasonEvicted}})
	checkInvariants(t, c)
}

func TestSieve_MatchesReference(t *testing.T) {
	for _, capacity := range []int{1, 2, 3, 7, 64} {
		for seed := range uint64(20) {
			t.Run(fmt.Sprintf("cap=%d/seed=%d", capacity, seed), func(t *testing.T) {
				rng := rand.New(rand.NewPCG(seed, uint64(capacity)))
				ref := newRefSieve(capacity)
				rec := &recorder[int, int]{}
				c, _ := newTestCache(t, Config[int, int]{MaxEntries: capacity, Shards: 1, OnEvict: rec.record})

				set := func(k int) {
					existed := ref.has(k)
					c.Set(k, k)
					victim, evicted := ref.set(k)
					var want []removal[int, int]
					switch {
					case existed:
						want = []removal[int, int]{{k, k, ReasonReplaced}}
					case evicted:
						want = []removal[int, int]{{victim, victim, ReasonEvicted}}
					}
					expectRemovals(t, rec.take(), want)
				}

				for step := range 2000 {
					k := rng.IntN(capacity * 3)
					if rng.IntN(2) == 0 {
						k = rng.IntN(max(1, capacity/2)) // hot keys exercise visited bits
					}
					switch op := rng.IntN(10); {
					case op < 6: // access: Get, Set on miss
						_, hit := c.Get(k)
						if want := ref.get(k); hit != want {
							t.Fatalf("step %d: Get(%d) hit=%v, reference %v", step, k, hit, want)
						}
						if !hit {
							set(k)
						}
					case op < 8:
						set(k)
					default:
						got, want := c.Delete(k), ref.del(k)
						if got != want {
							t.Fatalf("step %d: Delete(%d) = %v, reference %v", step, k, got, want)
						}
						var wantRem []removal[int, int]
						if want {
							wantRem = []removal[int, int]{{k, k, ReasonDeleted}}
						}
						expectRemovals(t, rec.take(), wantRem)
					}
					checkInvariants(t, c)
				}
			})
		}
	}
}
