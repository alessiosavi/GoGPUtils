package cache

import (
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

const modelDefaultTTL = 2 * time.Second

var (
	modelTTLs     = []time.Duration{0, NoExpiration, time.Nanosecond, time.Second, 2 * time.Second, time.Duration(math.MaxInt64)}
	modelAdvances = []time.Duration{0, time.Nanosecond, 500 * time.Millisecond, time.Second, 3 * time.Second}
)

type modelEntry struct {
	value     int
	expiresAt int64 // 0 = never
}

type modelOp struct {
	kind uint8
	key  int
	val  int
	ttl  time.Duration
	adv  time.Duration
}

// modelExpiry is an independent transcription of spec §3.2/§4.3.
func modelExpiry(now int64, ttl time.Duration) int64 {
	if ttl == 0 {
		ttl = modelDefaultTTL
	}
	if ttl < 0 {
		return 0
	}
	exp := now + int64(ttl)
	if exp < now {
		return math.MaxInt64
	}
	if exp == 0 {
		return 1
	}

	return exp
}

// runModel applies ops to a single-shard cache, a map model (values and
// deadlines), and the reference SIEVE (exact victims), checking return values,
// OnEvict reasons, counters, and invariants after every op.
func runModel(t *testing.T, capacity int, ops []modelOp) {
	t.Helper()

	rec := &recorder[int, int]{}
	c, clk := newTestCache(t, Config[int, int]{MaxEntries: capacity, Shards: 1, TTL: modelDefaultTTL, OnEvict: rec.record})
	model := map[int]modelEntry{}
	ref := newRefSieve(capacity) // exact SIEVE victim oracle (spec §6)
	var hits, misses, evictions, expirations uint64
	isExpired := func(e modelEntry) bool { return e.expiresAt != 0 && clk.now() >= e.expiresAt }

	for i, op := range ops {
		switch op.kind % 10 {
		case 0: // Get
			v, ok := c.Get(op.key)
			if ok {
				hits++
			} else {
				misses++
			}
			got := rec.take()
			me, present := model[op.key]
			switch {
			case !present:
				if ok || len(got) != 0 {
					t.Fatalf("op %d Get(%d) missing: %v, %v, removals %v", i, op.key, v, ok, got)
				}
			case isExpired(me):
				if ok {
					t.Fatalf("op %d Get(%d) returned an expired entry", i, op.key)
				}
				expectRemovals(t, got, []removal[int, int]{{op.key, me.value, ReasonExpired}})
				delete(model, op.key)
				ref.del(op.key)
				expirations++
			default:
				if !ok || v != me.value || len(got) != 0 {
					t.Fatalf("op %d Get(%d) = %v, %v, removals %v; want %v", i, op.key, v, ok, got, me.value)
				}
				ref.get(op.key)
			}
		case 1: // Peek
			v, ok := c.Peek(op.key)
			me, present := model[op.key]
			live := present && !isExpired(me)
			if ok != live || (live && v != me.value) {
				t.Fatalf("op %d Peek(%d) = %v, %v; model %+v present=%v", i, op.key, v, ok, me, present)
			}
			expectRemovals(t, rec.take(), nil)
		case 2: // TTL
			d, ok := c.TTL(op.key)
			me, present := model[op.key]
			wantD, wantOK := time.Duration(0), false
			switch {
			case !present || isExpired(me):
			case me.expiresAt == 0:
				wantD, wantOK = NoExpiration, true
			default:
				wantD, wantOK = time.Duration(me.expiresAt-clk.now()), true
			}
			if d != wantD || ok != wantOK {
				t.Fatalf("op %d TTL(%d) = %v, %v; want %v, %v", i, op.key, d, ok, wantD, wantOK)
			}
			expectRemovals(t, rec.take(), nil)
		case 3, 4: // Set or SetWithTTL
			if op.kind%10 == 3 {
				op.ttl = 0
				c.Set(op.key, op.val)
			} else {
				c.SetWithTTL(op.key, op.val, op.ttl)
			}
			got := rec.take()
			refVictim, refEvicted := ref.setWith(op.key, func(k int) bool { return isExpired(model[k]) })
			if me, present := model[op.key]; present {
				reason := ReasonReplaced
				if isExpired(me) {
					reason = ReasonExpired
					expirations++
				}
				expectRemovals(t, got, []removal[int, int]{{op.key, me.value, reason}})
			} else if capacity > 0 && len(model) >= capacity {
				if len(got) != 1 {
					t.Fatalf("op %d Set(%d) on a full cache: removals %v, want exactly one", i, op.key, got)
				}
				victim := got[0]
				if !refEvicted || victim.key != refVictim {
					t.Fatalf("op %d: SIEVE victim %d, reference victim %d (evicted=%v)", i, victim.key, refVictim, refEvicted)
				}
				me, ok := model[victim.key]
				if !ok || me.value != victim.value {
					t.Fatalf("op %d: victim %v is not in the model", i, victim)
				}
				switch {
				case isExpired(me) && victim.reason == ReasonExpired:
					expirations++
				case !isExpired(me) && victim.reason == ReasonEvicted:
					evictions++
				default:
					t.Fatalf("op %d: victim %v has the wrong reason for %+v", i, victim, me)
				}
				delete(model, victim.key)
			} else {
				expectRemovals(t, got, nil)
			}
			model[op.key] = modelEntry{value: op.val, expiresAt: modelExpiry(clk.now(), op.ttl)}
		case 5: // Delete
			ok := c.Delete(op.key)
			got := rec.take()
			me, present := model[op.key]
			switch {
			case !present:
				if ok {
					t.Fatalf("op %d Delete(%d) of a missing key returned true", i, op.key)
				}
				expectRemovals(t, got, nil)
			case isExpired(me):
				if ok {
					t.Fatalf("op %d Delete(%d) of an expired key returned true", i, op.key)
				}
				expectRemovals(t, got, []removal[int, int]{{op.key, me.value, ReasonExpired}})
				expirations++
			default:
				if !ok {
					t.Fatalf("op %d Delete(%d) of a live key returned false", i, op.key)
				}
				expectRemovals(t, got, []removal[int, int]{{op.key, me.value, ReasonDeleted}})
			}
			if present {
				ref.del(op.key)
			}
			delete(model, op.key)
		case 6: // DeleteExpired
			n := c.DeleteExpired()
			got := rec.take()
			want := map[int]int{}
			for k, me := range model {
				if isExpired(me) {
					want[k] = me.value
				}
			}
			if n != len(want) || len(got) != len(want) {
				t.Fatalf("op %d DeleteExpired = %d, removals %v; want %d", i, n, got, len(want))
			}
			for _, r := range got {
				if v, ok := want[r.key]; !ok || v != r.value || r.reason != ReasonExpired {
					t.Fatalf("op %d DeleteExpired removed %v unexpectedly", i, r)
				}
				delete(want, r.key)
				delete(model, r.key)
				ref.del(r.key) // same order as the cache, so hand fix-ups match
			}
			expirations += uint64(n)
		case 7: // Clear
			c.Clear()
			got := rec.take()
			if len(got) != len(model) {
				t.Fatalf("op %d Clear removals %v, model has %d entries", i, got, len(model))
			}
			for _, r := range got {
				if me, ok := model[r.key]; !ok || me.value != r.value || r.reason != ReasonDeleted {
					t.Fatalf("op %d Clear removed %v unexpectedly", i, r)
				}
				delete(model, r.key)
			}
			ref = newRefSieve(capacity)
		case 8: // advance the clock
			clk.advance(op.adv)
		case 9: // All
			got := maps.Collect(c.All())
			want := map[int]int{}
			for k, me := range model {
				if !isExpired(me) {
					want[k] = me.value
				}
			}
			if !maps.Equal(got, want) {
				t.Fatalf("op %d All = %v, want %v", i, got, want)
			}
		}

		checkInvariants(t, c)
		if c.Len() != len(model) || len(ref.keys) != len(model) {
			t.Fatalf("op %d: Len = %d, reference %d, model %d", i, c.Len(), len(ref.keys), len(model))
		}
		if st := c.Stats(); st != (Stats{Hits: hits, Misses: misses, Evictions: evictions, Expirations: expirations}) {
			t.Fatalf("op %d: Stats %+v, want %d evictions and %d expirations", i, st, evictions, expirations)
		}
	}
}

func randomOps(rng *rand.Rand, n int) []modelOp {
	ops := make([]modelOp, n)
	for i := range ops {
		ops[i] = modelOp{
			kind: uint8(rng.IntN(10)),
			key:  rng.IntN(8),
			val:  rng.IntN(1000),
			ttl:  modelTTLs[rng.IntN(len(modelTTLs))],
			adv:  modelAdvances[rng.IntN(len(modelAdvances))],
		}
	}

	return ops
}

func TestModel_RandomOperations(t *testing.T) {
	for _, capacity := range []int{0, 1, 2, 5} {
		for seed := range uint64(30) {
			t.Run(fmt.Sprintf("cap=%d/seed=%d", capacity, seed), func(t *testing.T) {
				rng := rand.New(rand.NewPCG(seed, uint64(capacity)))
				runModel(t, capacity, randomOps(rng, 400))
			})
		}
	}
}

func FuzzCache(f *testing.F) {
	f.Add([]byte{1, 4, 1, 2, 0, 8, 0, 0, 1, 0, 1, 0, 0}) // capacity 1; now == deadline
	f.Add([]byte{1, 4, 1, 5, 0, 2, 1, 0, 0})             // MaxInt64 TTL; saturation
	f.Add([]byte{3, 3, 0, 1, 3, 3, 1, 2, 3, 8, 0, 0, 3, 0, 0, 0, 0})
	f.Add([]byte{1, 3, 1, 3, 1, 3, 2, 3, 2, 8, 3, 3, 4, 0, 1, 2, 6, 0, 0, 0, 7, 0, 0, 0, 9, 0, 0, 0})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		rounding := [...]struct{ entries, requested, want int }{
			{0, 0, 16}, {2047, 0, 1}, {2048, 0, 2},
			{4096, 3, 4}, {1, 64, 1}, {MaxShards, MaxShards, MaxShards},
		}
		tc := rounding[int(data[0])%len(rounding)]
		if got := shardCount(tc.entries, tc.requested, 4); got != tc.want {
			t.Fatalf("shardCount(%d, %d, 4) = %d, want %d", tc.entries, tc.requested, got, tc.want)
		}
		capacity := int(data[0] % 6) // 0 = unbounded
		data = data[1:]
		const maxOps = 256
		ops := make([]modelOp, 0, min(len(data)/4, maxOps))
		for len(data) >= 4 && len(ops) < maxOps {
			ops = append(ops, modelOp{
				kind: data[0],
				key:  int(data[1] % 8),
				val:  int(data[2]),
				ttl:  modelTTLs[int(data[2])%len(modelTTLs)],
				adv:  modelAdvances[int(data[3])%len(modelAdvances)],
			})
			data = data[4:]
		}
		runModel(t, capacity, ops)
	})
}
