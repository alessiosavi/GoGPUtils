package mathutil

import (
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"runtime"
	"slices"
	"testing"
	"weak"
)

// BASE 4c05787953a944b57478180f4c27100869db12f6. Bodies are copied verbatim,
// except baseQuartiles calls basePercentile to keep the reference independent.
func basePercentile[T Number](s []T, p float64) float64 {
	if len(s) == 0 || p < 0 || p > 100 {
		return 0
	}

	sorted := make([]T, len(s))
	copy(sorted, s)
	slices.Sort(sorted)

	if p == 0 {
		return float64(sorted[0])
	}

	if p == 100 {
		return float64(sorted[len(sorted)-1])
	}

	rank := (p / 100) * float64(len(sorted)-1)
	lower := int(rank)

	upper := lower + 1
	if upper >= len(sorted) {
		return float64(sorted[len(sorted)-1])
	}

	fraction := rank - float64(lower)

	return float64(sorted[lower]) + fraction*(float64(sorted[upper])-float64(sorted[lower]))
}

func baseQuartiles[T Number](s []T) (q1, q2, q3 float64) {
	q1 = basePercentile(s, 25)
	q2 = basePercentile(s, 50)
	q3 = basePercentile(s, 75)

	return
}

func baseMatrixMultiply[T Number](a, b Matrix[T]) (Matrix[T], error) {
	if len(a) == 0 || len(b) == 0 {
		return nil, ErrInvalidDimensions
	}

	m := len(a)

	n := len(a[0])
	if len(b) != n {
		return nil, ErrInvalidDimensions
	}

	p := len(b[0])

	// Verify all rows have consistent length
	for _, row := range a {
		if len(row) != n {
			return nil, ErrInvalidDimensions
		}
	}

	for _, row := range b {
		if len(row) != p {
			return nil, ErrInvalidDimensions
		}
	}

	result := make(Matrix[T], m)
	for i := range result {
		result[i] = make([]T, p)

		for j := range p {
			var sum T
			for k := range n {
				sum += a[i][k] * b[k][j]
			}

			result[i][j] = sum
		}
	}

	return result, nil
}

func baseMinIndex[T Number](s []T) int {
	if len(s) == 0 {
		return -1
	}

	minIdx := 0
	for i, v := range s {
		if v < s[minIdx] {
			minIdx = i
		}
	}

	return minIdx
}

func baseMaxIndex[T Number](s []T) int {
	if len(s) == 0 {
		return -1
	}

	maxIdx := 0
	for i, v := range s {
		if v > s[maxIdx] {
			maxIdx = i
		}
	}

	return maxIdx
}

func basePrimes(n int) []int {
	if n < 2 {
		return nil
	}

	// Sieve of Eratosthenes
	sieve := make([]bool, n+1)
	for i := range sieve {
		sieve[i] = true
	}

	sieve[0], sieve[1] = false, false

	for i := 2; i*i <= n; i++ {
		if sieve[i] {
			for j := i * i; j <= n; j += i {
				sieve[j] = false
			}
		}
	}

	var primes []int

	for i, isPrime := range sieve {
		if isPrime {
			primes = append(primes, i)
		}
	}

	return primes
}

type myFloat float64
type myInt int64

func oracleBits[T Number](v T) uint64 {
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Float64:
		return math.Float64bits(float64(v))
	case reflect.Float32:
		return uint64(math.Float32bits(float32(v)))
	default:
		return uint64(v)
	}
}

func oracleSliceEqual[T Number](a, b []T) bool {
	if (a == nil) != (b == nil) || len(a) != len(b) {
		return false
	}
	for i := range a {
		if oracleBits(a[i]) != oracleBits(b[i]) {
			return false
		}
	}
	return true
}

// Diagnostic helpers are called only on failure, keeping formatting out of the random loops.
func oracleSliceBits[T Number](s []T) []uint64 {
	if s == nil {
		return nil
	}
	bits := make([]uint64, len(s))
	for i, v := range s {
		bits[i] = oracleBits(v)
	}
	return bits
}

func oracleSliceDifference[T Number](got, want []T) string {
	for i := range min(len(got), len(want)) {
		if oracleBits(got[i]) != oracleBits(want[i]) {
			return fmt.Sprintf("index %d: want bits=0x%016x, got bits=0x%016x", i, oracleBits(want[i]), oracleBits(got[i]))
		}
	}
	if len(got) != len(want) {
		return fmt.Sprintf("first missing index %d: want length=%d, got length=%d", min(len(got), len(want)), len(want), len(got))
	}
	if (got == nil) != (want == nil) {
		return fmt.Sprintf("want nil=%v, got nil=%v", want == nil, got == nil)
	}
	return ""
}

func oracleFloat64(r *rand.Rand) float64 {
	switch r.IntN(12) {
	case 0:
		return 0
	case 1:
		return math.Copysign(0, -1)
	case 2:
		return math.Inf(1)
	case 3:
		return math.Inf(-1)
	case 4:
		return math.Float64frombits((r.Uint64() & (1<<63 | (1 << 52) - 1)) | 0x7ff0000000000001)
	case 5:
		return math.Float64frombits((r.Uint64() & (1<<63 | (1 << 52) - 1)) | 1)
	case 6:
		return math.MaxFloat64
	case 7:
		return -math.MaxFloat64
	default:
		return math.Float64frombits(r.Uint64())
	}
}

func oracleFloat32(r *rand.Rand) float32 {
	switch r.IntN(12) {
	case 0:
		return 0
	case 1:
		return math.Float32frombits(1 << 31)
	case 2:
		return float32(math.Inf(1))
	case 3:
		return float32(math.Inf(-1))
	case 4:
		return math.Float32frombits((r.Uint32() & (1<<31 | (1 << 23) - 1)) | 0x7f800001)
	case 5:
		return math.Float32frombits((r.Uint32() & (1<<31 | (1 << 23) - 1)) | 1)
	case 6:
		return math.MaxFloat32
	case 7:
		return -math.MaxFloat32
	default:
		return math.Float32frombits(r.Uint32())
	}
}

// Start in the panicked state so even panic(nil) cannot look like a normal return.
func outcome[R any](f func() R) (r R, panicked bool, msg string) {
	panicked = true
	defer func() {
		if v := recover(); panicked {
			msg = fmt.Sprintf("%T: %v", v, v)
		}
	}()
	r = f()
	panicked = false
	return
}

func oracleCompare[R any](t *testing.T, label string, gotFn, wantFn func() R, equal func(R, R) bool, diagnostic func(R, R) string) (R, bool) {
	t.Helper()
	got, gp, gm := outcome(gotFn)
	want, wp, wm := outcome(wantFn)
	if gp != wp || (gp && gm != wm) || (!gp && !equal(got, want)) {
		t.Fatalf("%s: %s; got %#v panic=%v %q; want %#v panic=%v %q", label, diagnostic(got, want), got, gp, gm, want, wp, wm)
	}
	return got, gp
}

func oracleFloatEqual(a, b float64) bool { return math.Float64bits(a) == math.Float64bits(b) }

func oracleQuartileEqual(a, b [3]float64) bool {
	return oracleFloatEqual(a[0], b[0]) && oracleFloatEqual(a[1], b[1]) && oracleFloatEqual(a[2], b[2])
}

func oracleCheckPercentile[T Number](t *testing.T, s []T, p float64, caseIndex int) {
	t.Helper()
	oracleCompare(t, "Percentile", func() float64 { return Percentile(s, p) },
		func() float64 { return basePercentile(s, p) }, oracleFloatEqual, func(got, want float64) string {
			return fmt.Sprintf("case %d p=%g (bits=0x%016x): want bits=0x%016x, got bits=0x%016x; input bits=%#x",
				caseIndex, p, math.Float64bits(p), math.Float64bits(want), math.Float64bits(got), oracleSliceBits(s))
		})
}

func oracleQuartiles[T Number](t *testing.T, next func(*rand.Rand) T) {
	t.Helper()
	r := rand.New(rand.NewPCG(0x123, 0x456))
	for c := range 20000 {
		n := r.IntN(41)
		if c%200 == 0 {
			n = 41 + r.IntN(1960)
		}
		var s []T
		if c%2 == 0 || n > 0 {
			s = make([]T, n)
		}
		for i := range s {
			s[i] = next(r)
		}
		before := slices.Clone(s)
		oracleCompare(t, "Quartiles", func() [3]float64 {
			a, b, c := Quartiles(s)
			return [3]float64{a, b, c}
		}, func() [3]float64 {
			a, b, c := baseQuartiles(s)
			return [3]float64{a, b, c}
		}, oracleQuartileEqual, func(got, want [3]float64) string {
			return fmt.Sprintf("case %d: %s; input bits=%#x", c, oracleSliceDifference(got[:], want[:]), oracleSliceBits(s))
		})
		oracleCompare(t, "IQR", func() float64 { return IQR(s) }, func() float64 {
			a, _, c := baseQuartiles(s)
			return c - a
		}, oracleFloatEqual, func(got, want float64) string {
			return fmt.Sprintf("case %d: want bits=0x%016x, got bits=0x%016x; input bits=%#x",
				c, math.Float64bits(want), math.Float64bits(got), oracleSliceBits(s))
		})
		// C-07 intentionally changes NaN p; TestPercentileNaN pins its zero result.
		// Keep BASE equivalence for all other p, including slices containing NaNs.
		for _, p := range []float64{-1, 0, 0.5, 25, 50, 75, 99.9, 100, 101, math.Inf(-1), math.Inf(1)} {
			oracleCheckPercentile(t, s, p, c)
		}
		if !oracleSliceEqual(s, before) {
			t.Fatalf("case %d mutated input: %s", c, oracleSliceDifference(s, before))
		}
	}
}

func TestOracleQuartilesAndPercentile(t *testing.T) {
	t.Run("float64", func(t *testing.T) { oracleQuartiles(t, oracleFloat64) })
	t.Run("float32", func(t *testing.T) { oracleQuartiles(t, oracleFloat32) })
	t.Run("int", func(t *testing.T) { oracleQuartiles(t, func(r *rand.Rand) int { return int(r.Uint64()) }) })
	t.Run("int8", func(t *testing.T) { oracleQuartiles(t, func(r *rand.Rand) int8 { return int8(r.Uint32()) }) })
	t.Run("myFloat", func(t *testing.T) {
		oracleQuartiles(t, func(r *rand.Rand) myFloat { return myFloat(oracleFloat64(r)) })
	})
	t.Run("myInt", func(t *testing.T) {
		oracleQuartiles(t, func(r *rand.Rand) myInt { return myInt(r.Uint64()) })
	})
}

func TestPercentileValidationAllocs(t *testing.T) {
	big := make([]float64, 2048)
	for i := range big {
		big[i] = float64(len(big) - i)
	}
	for _, c := range []struct {
		name  string
		input []float64
	}{
		{"big", big}, {"nil", nil}, {"empty", []float64{}},
	} {
		for _, p := range []float64{-1, 101, math.Inf(-1), math.Inf(1), math.NaN()} {
			t.Run(fmt.Sprintf("%s/p=%v", c.name, p), func(t *testing.T) {
				if got := testing.AllocsPerRun(50, func() { Percentile(c.input, p) }); got != 0 {
					t.Fatalf("allocations=%g want 0", got)
				}
			})
		}
		if len(c.input) == 0 {
			// Empty input must also short-circuit otherwise-valid p.
			for _, p := range []float64{0, 50, 100} {
				t.Run(fmt.Sprintf("%s/p=%v", c.name, p), func(t *testing.T) {
					if got := testing.AllocsPerRun(50, func() { Percentile(c.input, p) }); got != 0 {
						t.Fatalf("allocations=%g want 0", got)
					}
				})
			}
			t.Run(c.name+"/Quartiles", func(t *testing.T) {
				if got := testing.AllocsPerRun(50, func() { Quartiles(c.input) }); got != 0 {
					t.Fatalf("allocations=%g want 0", got)
				}
			})
		}
	}
}

func oracleCloneMatrix[T Number](m Matrix[T]) Matrix[T] {
	if m == nil {
		return nil
	}
	out := make(Matrix[T], len(m))
	for i := range m {
		out[i] = slices.Clone(m[i])
	}
	return out
}

func oracleMatrixEqual[T Number](a, b Matrix[T]) bool {
	if (a == nil) != (b == nil) || len(a) != len(b) {
		return false
	}
	for i := range a {
		if !oracleSliceEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// Owner decision 2026-10-04, findings §2 item 5: only MatrixMultiply results may
// differ in NaN payload/sign; NaN-ness and every non-NaN bit must still match BASE.
func oracleMatrixResultCellEqual[T Number](a, b T) bool {
	return oracleBits(a) == oracleBits(b) || (math.IsNaN(float64(a)) && math.IsNaN(float64(b)))
}

func oracleMatrixResultEqual[T Number](a, b Matrix[T]) bool {
	if (a == nil) != (b == nil) || len(a) != len(b) {
		return false
	}
	for i := range a {
		if (a[i] == nil) != (b[i] == nil) || len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if !oracleMatrixResultCellEqual(a[i][j], b[i][j]) {
				return false
			}
		}
	}
	return true
}

func oracleMatrixBits[T Number](m Matrix[T]) [][]uint64 {
	if m == nil {
		return nil
	}
	bits := make([][]uint64, len(m))
	for i, row := range m {
		bits[i] = oracleSliceBits(row)
	}
	return bits
}

func oracleMatrixDifference[T Number](got, want Matrix[T], equal func(T, T) bool) string {
	for i := range min(len(got), len(want)) {
		for j := range min(len(got[i]), len(want[i])) {
			if !equal(got[i][j], want[i][j]) {
				return fmt.Sprintf("cell [%d][%d]: want bits=0x%016x, got bits=0x%016x", i, j, oracleBits(want[i][j]), oracleBits(got[i][j]))
			}
		}
		if len(got[i]) != len(want[i]) {
			return fmt.Sprintf("row %d first missing index %d: want length=%d, got length=%d", i, min(len(got[i]), len(want[i])), len(want[i]), len(got[i]))
		}
		if (got[i] == nil) != (want[i] == nil) {
			return fmt.Sprintf("row %d: want nil=%v, got nil=%v", i, want[i] == nil, got[i] == nil)
		}
	}
	if len(got) != len(want) {
		return fmt.Sprintf("first missing row %d: want rows=%d, got rows=%d", min(len(got), len(want)), len(want), len(got))
	}
	if (got == nil) != (want == nil) {
		return fmt.Sprintf("want nil=%v, got nil=%v", want == nil, got == nil)
	}
	return ""
}

func oracleCheckMatrix[T Number](t *testing.T, a, b Matrix[T], caseIndex int) {
	t.Helper()
	ac, bc := oracleCloneMatrix(a), oracleCloneMatrix(b)
	type result struct {
		matrix Matrix[T]
		err    error
	}
	out, panicked := oracleCompare(t, "MatrixMultiply", func() result {
		m, err := MatrixMultiply(a, b)
		return result{m, err}
	}, func() result {
		m, err := baseMatrixMultiply(a, b)
		return result{m, err}
	}, func(a, b result) bool { return a.err == b.err && oracleMatrixResultEqual(a.matrix, b.matrix) }, //nolint:errorlint // identity comparison against BASE is intended
		func(got, want result) string {
			return fmt.Sprintf("case %d: %s; a bits=%#x; b bits=%#x", caseIndex,
				oracleMatrixDifference(got.matrix, want.matrix, oracleMatrixResultCellEqual[T]), oracleMatrixBits(ac), oracleMatrixBits(bc))
		})
	got := out.matrix
	bitEqual := func(x, y T) bool { return oracleBits(x) == oracleBits(y) }
	if !oracleMatrixEqual(a, ac) {
		t.Fatalf("case %d matrix a input mutated: %s", caseIndex, oracleMatrixDifference(a, ac, bitEqual))
	}
	if !oracleMatrixEqual(b, bc) {
		t.Fatalf("case %d matrix b input mutated: %s", caseIndex, oracleMatrixDifference(b, bc, bitEqual))
	}
	if panicked {
		return
	}
	// Also check overlapping rows; the weak-pointer test checks allocation independence.
	for i := range got {
		for j := i + 1; j < len(got); j++ {
			for x := range got[i][:cap(got[i])] {
				for y := range got[j][:cap(got[j])] {
					if &got[i][:cap(got[i])][x] == &got[j][:cap(got[j])][y] {
						t.Fatalf("case %d result rows overlap at [%d][%d] and [%d][%d]", caseIndex, i, x, j, y)
					}
				}
			}
		}
	}
	if len(got) > 1 {
		before := slices.Clone(got[1])
		got[0] = append(got[0], T(123))
		if !oracleSliceEqual(got[1], before) {
			t.Fatalf("case %d append row 0 changed row 1: %s", caseIndex, oracleSliceDifference(got[1], before))
		}
	}
}

func oracleMatrices[T Number](t *testing.T, next func(*rand.Rand) T) {
	t.Helper()
	r := rand.New(rand.NewPCG(0x789, 0xabc))
	for c := range 4000 {
		m, n, p := 1+r.IntN(9), 1+r.IntN(9), 1+r.IntN(9)
		a, b := make(Matrix[T], m), make(Matrix[T], n)
		for i := range a {
			a[i] = make([]T, n)
			for j := range a[i] {
				a[i][j] = next(r)
			}
		}
		for i := range b {
			b[i] = make([]T, p)
			for j := range b[i] {
				b[i][j] = next(r)
			}
		}
		oracleCheckMatrix(t, a, b, c)
	}
	for i, c := range []struct{ a, b Matrix[T] }{
		{nil, nil}, {Matrix[T]{}, Matrix[T]{{1}}}, {Matrix[T]{{1}}, nil},
		{Matrix[T]{{}}, Matrix[T]{{1}}}, {Matrix[T]{{1, 2}}, Matrix[T]{{1}}},
		{Matrix[T]{{1}, {1, 2}}, Matrix[T]{{1}}},
		{Matrix[T]{{1, 2}}, Matrix[T]{{1}, {1, 2}}},
		{Matrix[T]{{1}, {2}}, Matrix[T]{{}}},
	} {
		oracleCheckMatrix(t, c.a, c.b, i)
	}
}

func TestOracleMatrixMultiply(t *testing.T) {
	t.Run("float64", func(t *testing.T) { oracleMatrices(t, oracleFloat64) })
	t.Run("float32", func(t *testing.T) { oracleMatrices(t, oracleFloat32) })
	t.Run("int32", func(t *testing.T) { oracleMatrices(t, func(r *rand.Rand) int32 { return int32(r.Uint32()) }) })
	t.Run("int64", func(t *testing.T) { oracleMatrices(t, func(r *rand.Rand) int64 { return int64(r.Uint64()) }) })
	t.Run("myFloat", func(t *testing.T) { oracleMatrices(t, func(r *rand.Rand) myFloat { return myFloat(oracleFloat64(r)) }) })
	t.Run("myInt", func(t *testing.T) { oracleMatrices(t, func(r *rand.Rand) myInt { return myInt(r.Uint64()) }) })
	// Deterministic overflowing products and sums, in addition to random values.
	t.Run("overflow-int32", func(t *testing.T) { oracleMatrixOverflow[int32](t, math.MinInt32, math.MaxInt32) })
	t.Run("overflow-int64", func(t *testing.T) { oracleMatrixOverflow[int64](t, math.MinInt64, math.MaxInt64) })
	t.Run("overflow-myInt", func(t *testing.T) { oracleMatrixOverflow[myInt](t, math.MinInt64, math.MaxInt64) })
}

func TestOracleMatrixMultiplyNaNPropagation(t *testing.T) {
	// Construct distinct payloads and opposite signs at each precision, including
	// signaling NaNs. These cases assert NaN-ness parity with BASE, not NaN bits.
	t.Run("float64/quiet", func(t *testing.T) {
		oracleMatrixNaNPropagation(t, math.Float64frombits(0x7ff8000000000041), math.Float64frombits(0xfff8000000000082))
	})
	t.Run("float64/signaling", func(t *testing.T) {
		oracleMatrixNaNPropagation(t, math.Float64frombits(0x7ff0000000000041), math.Float64frombits(0xfff0000000000082))
	})
	t.Run("float32/quiet", func(t *testing.T) {
		oracleMatrixNaNPropagation(t, math.Float32frombits(0x7fc00041), math.Float32frombits(0xffc00082))
	})
	t.Run("float32/signaling", func(t *testing.T) {
		oracleMatrixNaNPropagation(t, math.Float32frombits(0x7f800041), math.Float32frombits(0xff800082))
	})
}

func oracleMatrixNaNPropagation[T Float](t *testing.T, nanA, nanB T) {
	t.Helper()
	inf := T(math.Inf(1))
	for i, c := range []struct {
		name string
		a, b Matrix[T]
	}{
		{
			name: "left-payload-a-then-b",
			a:    Matrix[T]{{nanA, nanB}},
			b:    Matrix[T]{{1}, {1}},
		},
		{
			name: "left-payload-b-then-a",
			a:    Matrix[T]{{nanB, nanA}},
			b:    Matrix[T]{{1}, {1}},
		},
		{
			name: "right-payload-a-then-b",
			a:    Matrix[T]{{1, 1}},
			b:    Matrix[T]{{nanA}, {nanB}},
		},
		{
			name: "right-payload-b-then-a",
			a:    Matrix[T]{{1, 1}},
			b:    Matrix[T]{{nanB}, {nanA}},
		},
		{
			name: "inf-times-zero-then-input",
			a:    Matrix[T]{{inf, nanA}, {inf, nanB}},
			b:    Matrix[T]{{0}, {1}},
		},
		{
			name: "input-then-inf-times-zero",
			a:    Matrix[T]{{nanA, inf}, {nanB, inf}},
			b:    Matrix[T]{{1}, {0}},
		},
		{
			name: "inf-minus-inf-then-input",
			a:    Matrix[T]{{inf, -inf, nanA}, {inf, -inf, nanB}},
			b:    Matrix[T]{{1}, {1}, {1}},
		},
		{
			name: "both-operands-a-times-b",
			a:    Matrix[T]{{nanA}},
			b:    Matrix[T]{{nanB}},
		},
		{
			name: "both-operands-b-times-a",
			a:    Matrix[T]{{nanB}},
			b:    Matrix[T]{{nanA}},
		},
	} {
		t.Run(c.name, func(t *testing.T) { oracleCheckMatrix(t, c.a, c.b, i) })
	}
}

func TestOracleMatrixResultEqual(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		oracleCheckMatrixResultEqual(t, math.Float64frombits(0x7ff8000000000041), math.Float64frombits(0xfff0000000000082))
	})
	t.Run("float32", func(t *testing.T) {
		oracleCheckMatrixResultEqual(t, math.Float32frombits(0x7fc00041), math.Float32frombits(0xff800082))
	})
	t.Run("myFloat", func(t *testing.T) {
		oracleCheckMatrixResultEqual(t, myFloat(math.Float64frombits(0x7ff8000000000041)), myFloat(math.Float64frombits(0xfff0000000000082)))
	})
	t.Run("int64", func(t *testing.T) {
		if oracleMatrixResultEqual(Matrix[int64]{{1 << 53}}, Matrix[int64]{{1<<53 + 1}}) {
			t.Fatal("distinct integer bits compared equal")
		}
	})
}

func oracleCheckMatrixResultEqual[T Float](t *testing.T, nanA, nanB T) {
	t.Helper()
	nz, inf := T(math.Copysign(0, -1)), T(math.Inf(1))
	for _, c := range []struct {
		name      string
		got, want Matrix[T]
		equal     bool
	}{
		{"equal-bits", Matrix[T]{{0, nz, inf, -inf, 1}}, Matrix[T]{{0, nz, inf, -inf, 1}}, true},
		{"both-nan", Matrix[T]{{nanA}}, Matrix[T]{{nanB}}, true},
		{"nan-vs-number", Matrix[T]{{nanA}}, Matrix[T]{{1}}, false},
		{"number-vs-nan", Matrix[T]{{1}}, Matrix[T]{{nanB}}, false},
		{"positive-vs-negative-zero", Matrix[T]{{0}}, Matrix[T]{{nz}}, false},
		{"negative-vs-positive-zero", Matrix[T]{{nz}}, Matrix[T]{{0}}, false},
		{"positive-vs-negative-inf", Matrix[T]{{inf}}, Matrix[T]{{-inf}}, false},
		{"different-finite", Matrix[T]{{1}}, Matrix[T]{{2}}, false},
		{"nan-before-different-bits", Matrix[T]{{nanA, 0}}, Matrix[T]{{nanB, nz}}, false},
		{"nil-vs-empty", nil, Matrix[T]{}, false},
		{"nil-vs-empty-row", Matrix[T]{nil}, Matrix[T]{{}}, false},
		{"different-rows", Matrix[T]{{1}}, Matrix[T]{{1}, {1}}, false},
		{"different-columns", Matrix[T]{{1}}, Matrix[T]{{1, 1}}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := oracleMatrixResultEqual(c.got, c.want); got != c.equal {
				t.Fatalf("equal=%v, want %v; got bits=%#x, want bits=%#x", got, c.equal, oracleMatrixBits(c.got), oracleMatrixBits(c.want))
			}
		})
	}
	// Input immutability and row independence must still reject changes to NaN bits.
	if oracleMatrixEqual(Matrix[T]{{nanA}}, Matrix[T]{{nanB}}) || oracleSliceEqual([]T{nanA}, []T{nanB}) {
		t.Fatal("strict input/row comparator accepted different NaN bits")
	}
}

func TestOracleDifferenceDiagnostics(t *testing.T) {
	nanA, nanB := math.Float64frombits(0x7ff8000000000041), math.Float64frombits(0xfff8000000000082)
	got, want := Matrix[float64]{{nanA, 1}, {0, 2}}, Matrix[float64]{{nanB, 1}, {math.Copysign(0, -1), 3}}
	if diff := oracleMatrixDifference(got, want, oracleMatrixResultCellEqual[float64]); diff != "cell [1][0]: want bits=0x8000000000000000, got bits=0x0000000000000000" {
		t.Fatalf("unexpected first matrix difference: %s", diff)
	}
	if diff := oracleSliceDifference(got[0], want[0]); diff != "index 0: want bits=0xfff8000000000082, got bits=0x7ff8000000000041" {
		t.Fatalf("unexpected first strict slice difference: %s", diff)
	}
	oracleCompare(t, "equal", func() float64 { return 1 }, func() float64 { return 1 }, oracleFloatEqual, func(_, _ float64) string {
		t.Fatal("formatted diagnostic without a mismatch")
		return ""
	})
}

func oracleMatrixOverflow[T Integer](t *testing.T, lo, hi T) {
	t.Helper()
	oracleCheckMatrix(t, Matrix[T]{{hi, hi}, {lo, hi}}, Matrix[T]{{2, hi}, {2, 1}}, 0)
}

//go:noinline
func oracleRetainedRow(fn func(Matrix[float64], Matrix[float64]) (Matrix[float64], error)) ([]float64, weak.Pointer[float64]) {
	// 8x512 times 512x512 produces eight distinct 4 KB result rows.
	a, b := make(Matrix[float64], 8), make(Matrix[float64], 512)
	for i := range a {
		a[i] = make([]float64, 512)
	}
	for i := range b {
		b[i] = make([]float64, 512)
	}
	r, err := fn(a, b)
	if err != nil {
		panic(err)
	}
	// Returning from a noinline helper drops r before collection.
	return r[0], weak.Make(&r[1][0])
}

func TestMatrixMultiplyRowsIndependentAllocations(t *testing.T) {
	for _, c := range []struct {
		name string
		fn   func(Matrix[float64], Matrix[float64]) (Matrix[float64], error)
	}{{"BASE", baseMatrixMultiply[float64]}, {"candidate", MatrixMultiply[float64]}} {
		t.Run(c.name, func(t *testing.T) {
			row, other := oracleRetainedRow(c.fn)
			runtime.GC()
			runtime.GC()
			retained := other.Value() != nil
			runtime.KeepAlive(row)
			if retained {
				t.Fatal("keeping row 0 retains row 1's allocation")
			}
		})
	}
}

func oracleIndexes[T Number](t *testing.T, next func(*rand.Rand) T, cases [][]T) {
	t.Helper()
	r := rand.New(rand.NewPCG(0xdef, 0x123))
	cases = append(cases, nil, []T{}, []T{1, 1, 1}, []T{2, 1, 1, 2})
	for range 20000 {
		s := make([]T, r.IntN(41))
		for i := range s {
			s[i] = next(r)
		}
		cases = append(cases, s)
	}
	for c, s := range cases {
		before := slices.Clone(s)
		equal := func(a, b int) bool { return a == b }
		diagnostic := func(got, want int) string {
			return fmt.Sprintf("case %d: want index=%d (bits=0x%016x), got index=%d (bits=0x%016x); input bits=%#x",
				c, want, oracleBits(want), got, oracleBits(got), oracleSliceBits(s))
		}
		oracleCompare(t, "MinIndex", func() int { return MinIndex(s) }, func() int { return baseMinIndex(s) }, equal, diagnostic)
		oracleCompare(t, "MaxIndex", func() int { return MaxIndex(s) }, func() int { return baseMaxIndex(s) }, equal, diagnostic)
		if !oracleSliceEqual(s, before) {
			t.Fatalf("case %d extrema input mutated: %s", c, oracleSliceDifference(s, before))
		}
	}
}

func TestOracleMinMaxIndex(t *testing.T) {
	t.Run("float64", func(t *testing.T) { oracleFloatIndexes(t, oracleFloat64, math.MaxFloat64, math.SmallestNonzeroFloat64) })
	t.Run("float32", func(t *testing.T) { oracleFloatIndexes(t, oracleFloat32, math.MaxFloat32, math.SmallestNonzeroFloat32) })
	t.Run("myFloat", func(t *testing.T) {
		oracleFloatIndexes(t, func(r *rand.Rand) myFloat { return myFloat(oracleFloat64(r)) }, myFloat(math.MaxFloat64), myFloat(math.SmallestNonzeroFloat64))
	})
	t.Run("int", func(t *testing.T) { oracleIntegerIndexes[int](t) })
	t.Run("int8", func(t *testing.T) { oracleIntegerIndexes[int8](t) })
	t.Run("int16", func(t *testing.T) { oracleIntegerIndexes[int16](t) })
	t.Run("int32", func(t *testing.T) { oracleIntegerIndexes[int32](t) })
	t.Run("int64", func(t *testing.T) { oracleIntegerIndexes[int64](t) })
	t.Run("uint", func(t *testing.T) { oracleIntegerIndexes[uint](t) })
	t.Run("uint8", func(t *testing.T) { oracleIntegerIndexes[uint8](t) })
	t.Run("uint16", func(t *testing.T) { oracleIntegerIndexes[uint16](t) })
	t.Run("uint32", func(t *testing.T) { oracleIntegerIndexes[uint32](t) })
	t.Run("uint64", func(t *testing.T) { oracleIntegerIndexes[uint64](t) })
	t.Run("uintptr", func(t *testing.T) { oracleIntegerIndexes[uintptr](t) })
	t.Run("myInt", func(t *testing.T) { oracleIntegerIndexes[myInt](t) })
}

func oracleFloatIndexes[T Float](t *testing.T, next func(*rand.Rand) T, largest, smallest T) {
	t.Helper()
	nan, nz, inf := T(math.NaN()), T(math.Copysign(0, -1)), T(math.Inf(1))
	oracleIndexes(t, next, [][]T{
		{nan, 1, -1}, {1, nan, -1}, {0, nz}, {nz, 0}, {2, 2, 1, 1},
		{largest, -largest}, {-largest, largest}, {inf, -inf}, {-inf, inf},
		{smallest, -smallest}, {-smallest, smallest},
	})
}

func oracleIntegerIndexes[T Integer](t *testing.T) {
	t.Helper()
	typ := reflect.TypeFor[T]()
	bits := typ.Bits()
	signed := typ.Kind() >= reflect.Int && typ.Kind() <= reflect.Int64
	var lo, hi T
	if signed {
		hi = T((uint64(1) << (bits - 1)) - 1)
		lo = ^hi
	} else {
		hi = ^T(0)
	}
	cases := [][]T{
		{hi, lo}, {lo, hi}, {hi, hi - 1}, {hi - 1, hi},
		{lo + 1, lo}, {lo, lo + 1}, {lo, lo, hi, hi},
	}
	if bits >= 64 {
		edge := uint64(1) << 53
		cases = append(cases, []T{T(edge + 1), T(edge)}, []T{T(edge), T(edge + 1)})
		if signed {
			cases = append(cases, []T{-T(edge + 1), -T(edge)}, []T{-T(edge), -T(edge + 1)})
		}
	}
	oracleIndexes(t, func(r *rand.Rand) T { return T(r.Uint64()) }, cases)
}

func TestOraclePrimes(t *testing.T) {
	bounds := make([]int, 0, 5007)
	for n := -5; n <= 5000; n++ {
		bounds = append(bounds, n)
	}
	bounds = append(bounds, 65536)
	for c, n := range bounds {
		oracleCompare(t, "Primes", func() []int { return Primes(n) },
			func() []int { return basePrimes(n) }, oracleSliceEqual[int], func(got, want []int) string {
				return fmt.Sprintf("case %d n=%d: %s", c, n, oracleSliceDifference(got, want))
			})
	}
}
