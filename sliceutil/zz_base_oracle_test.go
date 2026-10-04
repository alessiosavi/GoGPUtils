package sliceutil

import (
	"fmt"
	"math"
	"math/bits"
	"math/rand/v2"
	"reflect"
	"slices"
	"sync"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// Private copies from BASE 4c05787, with only identifiers renamed.
var baseFastrandMu sync.Mutex
var baseFastrandState uint64 = 1

func baseUnion[T comparable](a, b []T) []T {
	if a == nil && b == nil {
		return nil
	}

	combined := make([]T, 0, len(a)+len(b))
	combined = append(combined, a...)
	combined = append(combined, b...)

	return baseUnique(combined)
}

func baseIntersect[T comparable](a, b []T) []T {
	if a == nil || b == nil {
		return nil
	}

	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)

	for _, v := range a {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}

	return baseUnique(result)
}

func baseUnique[T comparable](s []T) []T {
	if s == nil {
		return nil
	}

	if len(s) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{}, len(s))

	result := make([]T, 0, len(s))

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}

			result = append(result, v)
		}
	}

	return result
}

func baseShuffle[T any](s []T) []T {
	if s == nil {
		return nil
	}

	result := slices.Clone(s)
	baseShuffleInPlace(result)

	return result
}

func baseShuffleInPlace[T any](s []T) {
	// Fisher-Yates shuffle using simple PRNG seeded from time
	// For tests/determinism, see ShuffleWithSeed
	for i := len(s) - 1; i > 0; i-- {
		// Simple LCG for shuffling - not crypto secure but fast
		j := int(baseFastrand()) % (i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

func baseFastrand() uint32 {
	// xorshift64 — read-modify-write cannot be made atomic without a lock
	baseFastrandMu.Lock()
	baseFastrandState ^= baseFastrandState << 13
	baseFastrandState ^= baseFastrandState >> 7
	baseFastrandState ^= baseFastrandState << 17
	state := uint32(baseFastrandState)
	baseFastrandMu.Unlock()

	return state
}

func outcome[R any](f func() R) (r R, panicked bool, msg string) {
	completed := false
	defer func() {
		if !completed {
			panicked = true
			msg = fmt.Sprint(recover())
		}
	}()
	r = f()
	completed = true
	return
}

// Recursive bit equality also handles NaNs inside comparable interface values.
func oracleEqualValue(a, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	}
	if a.Type() != b.Type() {
		return false
	}
	switch a.Kind() {
	case reflect.Func:
		return a.Pointer() == b.Pointer()
	case reflect.Float64:
		return math.Float64bits(a.Float()) == math.Float64bits(b.Float())
	case reflect.Float32:
		return math.Float32bits(float32(a.Float())) == math.Float32bits(float32(b.Float()))
	case reflect.Interface:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return oracleEqualValue(a.Elem(), b.Elem())
	case reflect.Slice:
		if a.IsNil() != b.IsNil() || a.Len() != b.Len() {
			return false
		}
		for i := range a.Len() {
			if !oracleEqualValue(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Array, reflect.Struct:
		if a.Kind() == reflect.Array {
			for i := range a.Len() {
				if !oracleEqualValue(a.Index(i), b.Index(i)) {
					return false
				}
			}
		} else {
			for i := range a.NumField() {
				if !oracleEqualValue(a.Field(i), b.Field(i)) {
					return false
				}
			}
		}
		return true
	default:
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
}
func oracleEqual[T any](a, b T) bool { return oracleEqualValue(reflect.ValueOf(a), reflect.ValueOf(b)) }

func checkSets[T comparable](t *testing.T, cases [][]T, random func() T) {
	t.Helper()
	check := func(a, b []T) {
		t.Helper()
		for _, f := range []struct {
			name      string
			got, want func([]T, []T) []T
		}{{"Union", Union[T], baseUnion[T]}, {"Intersect", Intersect[T], baseIntersect[T]}} {
			aa, bb := slices.Clone(a), slices.Clone(b)
			got, gp, gm := outcome(func() []T { return f.got(a, b) })
			want, wp, wm := outcome(func() []T { return f.want(a, b) })
			if gp != wp || gm != wm || !oracleEqual(got, want) {
				t.Fatalf("%s a=%v b=%v: got %v panic=%t %q; want %v panic=%t %q", f.name, a, b, got, gp, gm, want, wp, wm)
			}
			if !oracleEqual(a, aa) || !oracleEqual(b, bb) {
				t.Fatalf("%s mutated input", f.name)
			}
			// Every returned element must have storage independent of either input.
			for i := range got {
				for j := range a {
					if &got[i] == &a[j] {
						t.Fatalf("%s aliases a", f.name)
					}
				}
				for j := range b {
					if &got[i] == &b[j] {
						t.Fatalf("%s aliases b", f.name)
					}
				}
			}
		}
	}
	for _, a := range cases {
		for _, b := range cases {
			check(a, b)
		}
	}
	r := rand.New(rand.NewPCG(2026, 8))
	for iter := range 2000 {
		n, m := r.IntN(41), r.IntN(41)
		if iter%200 == 0 {
			n, m = 1000, 1000
		}
		a, b := make([]T, n), make([]T, m)
		for i := range a {
			a[i] = random()
		}
		for i := range b {
			b[i] = random()
		}
		check(a, b)
	}
}

func TestOracleSets(t *testing.T) {
	r := rand.New(rand.NewPCG(8, 1))
	t.Run("int", func(t *testing.T) {
		cases := [][]int{nil, {}, {2, 1, 2, 3}, {1, 2, 1, 4}, {3, 2, 1}, {7}}
		// Exhaust all sequences of length <=3 over a two-value domain, pairwise.
		for n := 1; n <= 3; n++ {
			for mask := range 1 << n {
				s := make([]int, n)
				for i := range s {
					s[i] = (mask >> i) & 1
				}
				cases = append(cases, s)
			}
		}
		checkSets(t, cases, func() int { return r.IntN(16) - 8 })
	})
	t.Run("string", func(t *testing.T) {
		checkSets(t, [][]string{nil, {}, {"b", "a", "b"}, {"a", "b"}, {"\xff", "é", ""}}, func() string { return []string{"", "a", "b", "é", "\xff"}[r.IntN(5)] })
	})
	t.Run("float64", func(t *testing.T) {
		pool := []float64{0, math.Copysign(0, -1), math.Float64frombits(0x7ff8000000000001), math.Float64frombits(0xfff8000000000002), math.Inf(1), math.Inf(-1), math.SmallestNonzeroFloat64, 1, 2}
		checkSets(t, [][]float64{nil, {}, pool, {pool[2], pool[2], 0, pool[1]}, {pool[1], 0, 2, 1}}, func() float64 {
			if r.IntN(4) == 0 {
				return math.Float64frombits(r.Uint64())
			}
			return pool[r.IntN(len(pool))]
		})
	})
	t.Run("any", func(t *testing.T) {
		pool := []any{nil, 1, "a", 0.0, math.Copysign(0, -1), math.Float64frombits(0x7ff8000000000001), [2]int{1, 2}, []int{1}, map[int]int{1: 1}, func() {}}
		cases := [][]any{nil, {}, {1, "a", 1}, {"a", 1}, pool}
		for _, v := range pool {
			cases = append(cases, []any{v})
		}
		checkSets(t, cases, func() any { return pool[r.IntN(len(pool))] })
	})
}

func setBaseSeed(seed uint64) {
	if seed == 0 {
		seed = 1
	}
	baseFastrandMu.Lock()
	baseFastrandState = seed
	baseFastrandMu.Unlock()
}
func states() (uint64, uint64) {
	fastrandMu.Lock()
	s := fastrandState
	fastrandMu.Unlock()
	baseFastrandMu.Lock()
	b := baseFastrandState
	baseFastrandMu.Unlock()
	return s, b
}
func TestOracleShuffleSequence(t *testing.T) {
	fastrandMu.Lock()
	saved := fastrandState
	fastrandMu.Unlock()
	t.Cleanup(func() { SeedShuffle(saved) })
	r := rand.New(rand.NewPCG(8, 5))
	seeds := []uint64{0, 1, 42}
	for range 32 {
		seeds = append(seeds, r.Uint64())
	}
	// A length-n shuffle consumes n-1 draws: 17 and 33 end full batches;
	// 18 and 34 add a partial batch. Retain the former 64-draw boundaries too.
	sizes := []int{
		0, 1, 2, 15, 16, 17, 18, 31, 32, 33, 34,
		63, 64, 65, 127, 128, 129, 1000,
	}
	for _, seed := range seeds {
		for _, copying := range []bool{false, true} {
			SeedShuffle(seed)
			setBaseSeed(seed)
			for repeat := range 3 {
				for _, index := range r.Perm(len(sizes)) {
					n := sizes[index]
					a, b := make([]int, n), make([]int, n)
					for i := range a {
						a[i] = i
						b[i] = i
					}
					if n == 0 && repeat%2 == 0 {
						a, b = nil, nil
					}
					var got, want []int
					var gp, wp bool
					var gm, wm string
					if copying {
						got, gp, gm = outcome(func() []int { return Shuffle(a) })
						want, wp, wm = outcome(func() []int { return baseShuffle(b) })
					} else {
						got, gp, gm = outcome(func() []int { ShuffleInPlace(a); return a })
						want, wp, wm = outcome(func() []int { baseShuffleInPlace(b); return b })
					}
					gs, ws := states()
					if gp != wp || gm != wm || !oracleEqual(got, want) || !oracleEqual(a, b) || gs != ws {
						t.Fatalf("seed=%d n=%d copy=%t: panic %t/%t %q/%q; state %x/%x; output equal=%t partial equal=%t", seed, n, copying, gp, wp, gm, wm, gs, ws, oracleEqual(got, want), oracleEqual(a, b))
					}
					if copying && !gp {
						for i := range got {
							for j := range a {
								if &got[i] == &a[j] {
									t.Fatal("Shuffle aliases input")
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestShuffleConcurrentSeeding(t *testing.T) {
	if bits.UintSize == 32 {
		t.Skip("the separate sequence oracle covers existing 32-bit panics")
	}
	var wg sync.WaitGroup
	for worker := range 8 {
		wg.Go(func() {
			s := make([]int, 129+worker)
			for k := range 300 {
				if worker%3 == 0 {
					SeedShuffle(uint64(k))
				} else {
					ShuffleInPlace(s)
				}
			}
		})
	}
	wg.Wait()
}

// panicDetails preserves the recovered dynamic type as well as its message.
func panicDetails(f func()) (panicked bool, typ reflect.Type, msg string) {
	completed := false
	defer func() {
		if !completed {
			r := recover()
			panicked, typ, msg = true, reflect.TypeOf(r), fmt.Sprint(r)
		}
	}()
	f()
	completed = true
	return
}

func TestIntersectConsumedKeyPanic(t *testing.T) {
	a, b := []any{1, []int{2}}, []any{1}
	gp, gt, gm := panicDetails(func() { Intersect(a, b) })
	wp, wt, wm := panicDetails(func() { baseIntersect(a, b) })
	t.Logf("candidate panic=%t type=%v message=%q; BASE panic=%t type=%v message=%q", gp, gt, gm, wp, wt, wm)
	if !wp || gp != wp || gt != wt || gm != wm {
		t.Fatal("consuming the last key changed the subsequent unhashable-key panic")
	}
}

func TestShuffleConcurrentDrawAccounting(t *testing.T) {
	if bits.UintSize == 32 {
		t.Skip("32-bit BASE panics are covered separately")
	}
	saved, _ := states()
	t.Cleanup(func() { SeedShuffle(saved) })
	const workers, calls, size = 8, 200, 129
	const draws = workers * calls * (size - 1)
	SeedShuffle(42)
	setBaseSeed(42)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			s := make([]int, size)
			<-start
			for range calls {
				ShuffleInPlace(s)
			}
		})
	}
	close(start)
	wg.Wait()
	for range draws {
		baseFastrand()
	}
	got, want := states()
	if got != want {
		t.Fatalf("lost/duplicated draws: state=%x want=%x after %d draws", got, want, draws)
	}
	t.Logf("state=%x matches BASE after %d draws", got, draws)
}

func TestShuffle32PanicState(t *testing.T) {
	if bits.UintSize != 32 {
		t.Skip("requires 32-bit runtime")
	}
	saved, _ := states()
	t.Cleanup(func() { SeedShuffle(saved) })
	setBaseSeed(42)
	baseFastrand()
	_, oneDraw := states()
	SeedShuffle(42)
	setBaseSeed(42)
	a, b := make([]int, 16), make([]int, 16)
	for i := range a {
		a[i], b[i] = i, i
	}
	gp, gt, gm := panicDetails(func() { ShuffleInPlace(a) })
	if !fastrandMu.TryLock() {
		t.Fatal("global mutex remains locked after panic")
	}
	fastrandMu.Unlock()
	wp, wt, wm := panicDetails(func() { baseShuffleInPlace(b) })
	if !baseFastrandMu.TryLock() {
		t.Fatal("BASE mutex remains locked after panic")
	}
	baseFastrandMu.Unlock()
	got, want := states()
	if !wp || gp != wp || gt != wt || gm != wm || got != want || want != oneDraw || !oracleEqual(a, b) {
		t.Fatalf("panic=%t/%t type=%v/%v message=%q/%q state=%x/%x oneDraw=%x partialEqual=%t",
			gp, wp, gt, wt, gm, wm, got, want, oneDraw, oracleEqual(a, b))
	}
}

func allocationCheck[T comparable](t *testing.T, src []T) {
	var sink []T
	a, b := src[:16], src[8:]
	base := testing.AllocsPerRun(1000, func() { sink = baseIntersect(a, b) })
	got := testing.AllocsPerRun(1000, func() { sink = Intersect(a, b) })
	if len(sink) == 0 {
		t.Fatal("fixture unexpectedly empty")
	}
	t.Logf("%T n=16 BASE=%g candidate=%g allocs/run", src, base, got)
	if got > base {
		t.Fatalf("allocation gate %g > %g", got, base)
	}
}
func TestIntersectSmallAllocations(t *testing.T) {
	t.Run("int", func(t *testing.T) { allocationCheck(t, benchkit.Ints(24)) })
	t.Run("string", func(t *testing.T) { allocationCheck(t, benchkit.Strings(24, 8, "abcdefgh")) })
}
