package mathutil

import (
	"fmt"
	"math"
	"slices"
	"testing"
)

func TestPercentileNaN(t *testing.T) {
	t.Run("int", func(t *testing.T) { checkPercentileNaN[int](t) })
	t.Run("int8", func(t *testing.T) { checkPercentileNaN[int8](t) })
	t.Run("float32", func(t *testing.T) { checkPercentileNaN[float32](t) })
	t.Run("float64", func(t *testing.T) { checkPercentileNaN[float64](t) })
	t.Run("myInt", func(t *testing.T) { checkPercentileNaN[myInt](t) })
	t.Run("myFloat", func(t *testing.T) { checkPercentileNaN[myFloat](t) })
	t.Run("special-data", func(t *testing.T) {
		checkPercentileNaNInput(t, []float64{
			math.Inf(1), math.Float64frombits(0x7ff0000000000041),
			0, math.Copysign(0, -1), math.Inf(-1),
			math.Float64frombits(0xfff8000000000082),
		})
	})
}

func checkPercentileNaN[T Number](t *testing.T) {
	t.Helper()
	big := make([]T, 2048)
	for i := range big {
		big[i] = T((len(big) - i) % 128)
	}
	for _, tc := range []struct {
		name  string
		input []T
	}{
		{"nil", nil},
		{"empty", []T{}},
		{"singleton", []T{7}},
		{"pair", []T{2, 1}},
		{"multiple", []T{3, 1, 2}},
		{"big", big},
	} {
		t.Run(tc.name, func(t *testing.T) { checkPercentileNaNInput(t, tc.input) })
	}
}

func checkPercentileNaNInput[T Number](t *testing.T, input []T) {
	t.Helper()
	// Both signs of quiet and signaling NaNs must take the invalid-p path.
	for _, bits := range []uint64{
		0x7ff8000000000041, 0xfff8000000000082,
		0x7ff0000000000001, 0xfff0000000000042,
	} {
		t.Run(fmt.Sprintf("p=0x%016x", bits), func(t *testing.T) {
			p := math.Float64frombits(bits)
			before := slices.Clone(input)
			got, panicked, msg := outcome(func() float64 { return Percentile(input, p) })
			if !oracleSliceEqual(input, before) {
				t.Errorf("Percentile mutated input: %s", oracleSliceDifference(input, before))
			}
			if panicked {
				t.Fatalf("Percentile(input, NaN) panicked: %s", msg)
			}
			if math.Float64bits(got) != 0 {
				t.Errorf("Percentile(input, NaN) = %v (bits=0x%016x), want +0", got, math.Float64bits(got))
			}
			if allocs := testing.AllocsPerRun(50, func() { Percentile(input, p) }); allocs != 0 {
				t.Errorf("Percentile(input, NaN) allocations=%g, want 0", allocs)
			}
		})
	}
}
