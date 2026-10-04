package stringutil

import (
	"math"
	"strconv"
	"testing"
)

// Added with Task 4: these production helpers do not exist on BASE.
func TestPR3CosineGuard(t *testing.T) {
	const limit = 94_906_265
	for _, n := range []int{0, limit - 1, limit, limit + 1, math.MaxInt} {
		for _, pair := range [][2]int{{n, 0}, {0, n}, {n, limit}, {limit, n}} {
			want := n <= limit
			if got := cosineExactLengths(pair[0], pair[1]); got != want {
				t.Errorf("lengths=%v: got=%v want=%v", pair, got, want)
			}
		}
	}
	if cosineExactLengths(limit+1, limit+1) {
		t.Error("both oversized lengths accepted")
	}
}

// Artificial 64-bit counts exceed the public guarded domain to isolate
// conversion ordering without multi-GB strings. 50,000 is reachable on 32-bit.
func TestPR3CosineProducts(t *testing.T) {
	counts := []int64{50_000}
	if strconv.IntSize == 64 {
		counts = append(counts, 4_000_000_000)
	}
	for _, c := range counts {
		m1, m2 := map[string]int{"a": int(c)}, map[string]int{"a": int(c)}
		v := float64(c)
		want := (v * v) / (math.Sqrt(v*v) * math.Sqrt(v*v))
		if got := cosineFast(m1, m2); math.Float64bits(got) != math.Float64bits(want) {
			t.Errorf("count=%d: got=%g want=%g", c, got, want)
		}
	}
}
