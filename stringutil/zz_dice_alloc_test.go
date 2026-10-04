//go:build !race

package stringutil

import "testing"

// Race instrumentation can change escape analysis. These lowercase inputs
// avoid ToLower allocations and keep both bigram maps small enough for the stack.
func TestDiceCoefficientShortAllocs(t *testing.T) {
	var got float64
	allocs := testing.AllocsPerRun(100, func() {
		got = DiceCoefficient("abcdefgh", "abc#efg#")
	})
	if want := 4.0 / 7.0; got != want {
		t.Fatalf("DiceCoefficient = %g, want %g", got, want)
	}
	t.Logf("DiceCoefficient short lowercase ASCII: %g allocs/op", allocs)
	if allocs != 0 {
		t.Fatalf("DiceCoefficient allocated %g times, want 0", allocs)
	}
}
