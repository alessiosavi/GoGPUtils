package mathutil

import (
	"encoding/binary"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// BASE a8b7b93b997b8986200995fdf2c65a5efbf2053b. Body copied verbatim.
//
//nolint:intrange // Keep the frozen BASE loop verbatim.
func baseHistogram[T Number](data []T, bins []T) []int {
	if len(bins) < 2 {
		return nil
	}

	// Sort bins
	sortedBins := make([]T, len(bins))
	copy(sortedBins, bins)
	slices.Sort(sortedBins)

	counts := make([]int, len(sortedBins)-1)

	for _, v := range data {
		for i := 0; i < len(sortedBins)-1; i++ {
			if v >= sortedBins[i] && v < sortedBins[i+1] {
				counts[i]++

				break
			}
			// Include values equal to the last bin edge
			if i == len(sortedBins)-2 && v == sortedBins[i+1] {
				counts[i]++
			}
		}
	}

	return counts
}

func oracleCheckHistogram[T Number](t *testing.T, data, edges []T) {
	t.Helper()
	dataBefore, edgesBefore := slices.Clone(data), slices.Clone(edges)
	want := baseHistogram(dataBefore, edgesBefore)
	got := Histogram(data, edges)
	if (got == nil) != (want == nil) || !slices.Equal(got, want) {
		t.Fatalf("Histogram: got %v (nil=%v), want %v (nil=%v); data bits=%#x; edges bits=%#x",
			got, got == nil, want, want == nil, oracleSliceBits(dataBefore), oracleSliceBits(edgesBefore))
	}
	if !oracleSliceEqual(data, dataBefore) || !oracleSliceEqual(edges, edgesBefore) {
		t.Fatalf("Histogram mutated input: data %s; edges %s",
			oracleSliceDifference(data, dataBefore), oracleSliceDifference(edges, edgesBefore))
	}
	// Changing a result must neither mutate the inputs nor affect a later call.
	for i := range got {
		got[i] = -1
	}
	if !oracleSliceEqual(data, dataBefore) || !oracleSliceEqual(edges, edgesBefore) {
		t.Fatal("Histogram result aliases an input")
	}
	again := Histogram(data, edges)
	if (again == nil) != (want == nil) || !slices.Equal(again, want) {
		t.Fatal("changing Histogram result affected a later call")
	}
	for _, v := range got {
		if v != -1 {
			t.Fatal("later Histogram call changed a previous result")
		}
	}
}

func oracleHistogramDomain[T Number](t *testing.T, alphabet, probes []T, randomValue func(*rand.Rand) T) {
	t.Helper()
	cases := 0
	check := func(data, edges []T) {
		oracleCheckHistogram(t, data, edges)
		cases++
	}
	// Exhaust ordered edge tuples, probing one value at a time so that incorrect
	// counts cannot compensate for each other in an aggregate result.
	var walk func([]T, int)
	walk = func(edges []T, left int) {
		if left == 0 {
			for _, v := range probes {
				check([]T{v}, edges)
			}
			return
		}
		for _, v := range alphabet {
			walk(append(edges, v), left-1)
		}
	}
	for n := range 6 {
		walk(make([]T, 0, n), n)
	}
	check(nil, nil)
	check([]T{}, []T{})
	r := rand.New(rand.NewPCG(20261004, uint64(len(alphabet))))
	for n := range 301 {
		// Include subnormals and near-extreme values as edges, and probe them
		// individually at every size, including m=K-1/K/K+1 for each allowed K.
		cyclic := make([]T, n)
		for i := range cyclic {
			cyclic[i] = probes[(n-i)%len(probes)]
		}
		for _, v := range probes {
			check([]T{v}, cyclic)
		}
		// All-equal edges, including all-NaN, and one/two different edges among
		// them exercise empty bins and the last-edge rule at every bin count.
		for _, v := range alphabet {
			edges := make([]T, n)
			for i := range edges {
				edges[i] = v
			}
			check(probes, edges)
			if n == 0 {
				continue
			}
			for _, w := range alphabet {
				edges[0] = w
				check(probes, edges)
				edges[n-1] = w
				check(probes, edges)
				edges[n-1] = v
			}
		}
		for range 4 {
			edges := make([]T, n)
			for i := range edges {
				edges[i] = alphabet[r.IntN(len(alphabet))]
			}
			check(probes, edges)
			check(nil, edges)
			check([]T{}, edges)
			check(edges, edges)
			if n > 1 {
				check(edges[1:], edges[:n-1]) // Partially overlapping inputs.
			}
		}
	}
	for range 2000 {
		edges := make([]T, r.IntN(301))
		data := make([]T, r.IntN(25))
		for i := range edges {
			edges[i] = randomValue(r)
		}
		for i := range data {
			data[i] = randomValue(r)
		}
		check(data, edges)
	}
	t.Logf("%d differential cases; exhaustive edge tuples 0..5, systematic edge counts 0..300", cases)
}

func oracleHistogramSigned[T Integer](t *testing.T, smallest, largest T) {
	t.Helper()
	alphabet := []T{smallest, ^T(0), 0, 1, largest}
	probes := append(slices.Clone(alphabet), smallest+1, largest-1)
	oracleHistogramDomain(t, alphabet, probes, func(r *rand.Rand) T { return T(r.Uint64()) })
}

func oracleHistogramUnsigned[T Integer](t *testing.T) {
	t.Helper()
	alphabet := []T{0, 1, 2, ^T(0) - 1, ^T(0)}
	oracleHistogramDomain(t, alphabet, alphabet, func(r *rand.Rand) T { return T(r.Uint64()) })
}

func oracleHistogramFloat32[T ~float32](t *testing.T) {
	t.Helper()
	alphabet := []T{T(math.Inf(-1)), T(math.Float32frombits(0x80000000)), 0, 1,
		T(math.Inf(1)), T(math.Float32frombits(0x7fc12345)), T(math.Float32frombits(0xff800001))}
	probes := append(slices.Clone(alphabet), -math.MaxFloat32, math.MaxFloat32,
		-math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, -1)
	oracleHistogramDomain(t, alphabet, probes, func(r *rand.Rand) T { return T(oracleFloat32(r)) })
}

func oracleHistogramFloat64[T ~float64](t *testing.T) {
	t.Helper()
	alphabet := []T{T(math.Inf(-1)), T(math.Copysign(0, -1)), 0, 1,
		T(math.Inf(1)), T(math.Float64frombits(0x7ff8123456789abc)), T(math.Float64frombits(0xfff0000000000001))}
	probes := append(slices.Clone(alphabet), -math.MaxFloat64, math.MaxFloat64,
		-math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, -1)
	oracleHistogramDomain(t, alphabet, probes, func(r *rand.Rand) T { return T(oracleFloat64(r)) })
}

func TestZZHistogramDifferential(t *testing.T) {
	type namedUint uint64
	type namedFloat32 float32
	t.Run("int", func(t *testing.T) { oracleHistogramSigned(t, math.MinInt, math.MaxInt) })
	t.Run("int8", func(t *testing.T) { oracleHistogramSigned[int8](t, math.MinInt8, math.MaxInt8) })
	t.Run("int16", func(t *testing.T) { oracleHistogramSigned[int16](t, math.MinInt16, math.MaxInt16) })
	t.Run("int32", func(t *testing.T) { oracleHistogramSigned[int32](t, math.MinInt32, math.MaxInt32) })
	t.Run("int64", func(t *testing.T) { oracleHistogramSigned[int64](t, math.MinInt64, math.MaxInt64) })
	t.Run("uint", oracleHistogramUnsigned[uint])
	t.Run("uint8", oracleHistogramUnsigned[uint8])
	t.Run("uint16", oracleHistogramUnsigned[uint16])
	t.Run("uint32", oracleHistogramUnsigned[uint32])
	t.Run("uint64", oracleHistogramUnsigned[uint64])
	t.Run("uintptr", oracleHistogramUnsigned[uintptr])
	t.Run("float32", oracleHistogramFloat32[float32])
	t.Run("float64", oracleHistogramFloat64[float64])
	t.Run("named-int64", func(t *testing.T) { oracleHistogramSigned[myInt](t, math.MinInt64, math.MaxInt64) })
	t.Run("named-uint64", oracleHistogramUnsigned[namedUint])
	t.Run("named-float32", oracleHistogramFloat32[namedFloat32])
	t.Run("named-float64", oracleHistogramFloat64[myFloat])
}

func TestZZHistogramSemantics(t *testing.T) {
	nan := math.Float64frombits(0xfff0000000000001)
	for _, tc := range []struct {
		name        string
		data, edges []float64
		want        []int
	}{
		{"NaN,x", []float64{nan, -1, 0, 1, 2, math.Inf(1)}, []float64{nan, 1}, []int{1}},
		{"all-NaN", []float64{nan, 0, 1}, []float64{nan, nan, nan}, []int{0, 0}},
		{"NaN-prefix", []float64{nan, -1, 0, 1, 2}, []float64{1, nan, 0}, []int{0, 2}},
		{"duplicates", []float64{0, 1, 2}, []float64{2, 1, 0, 2, 1}, []int{1, 0, 1, 1}},
		{"signed-zero", []float64{0, math.Copysign(0, -1)}, []float64{0, math.Copysign(0, -1), 0}, []int{0, 2}},
		{"infinities", []float64{math.Inf(-1), -1, 0, math.Inf(1), nan}, []float64{math.Inf(1), 0, math.Inf(-1)}, []int{2, 2}},
		{"last-edge", []float64{2, 2, 3}, []float64{0, 2, 2}, []int{0, 2}},
		{"nil-bins", []float64{1}, nil, nil},
		{"empty-bins", []float64{1}, []float64{}, nil},
		{"one-edge", []float64{1}, []float64{1}, nil},
		{"nil-data", nil, []float64{0, 1, 2}, []int{0, 0}},
		{"empty-data", []float64{}, []float64{0, 1, 2}, []int{0, 0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oracleCheckHistogram(t, tc.data, tc.edges)
			got := Histogram(tc.data, tc.edges)
			if (got == nil) != (tc.want == nil) || !slices.Equal(got, tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func histogramFuzzBytes(values []float64) []byte {
	encoded := make([]byte, 8*len(values))
	for i, v := range values {
		binary.LittleEndian.PutUint64(encoded[8*i:], math.Float64bits(v))
	}
	return encoded
}

func histogramFuzzValues[T Number](encoded []byte, limit int, decode func(uint64) T) []T {
	if encoded == nil {
		return nil
	}
	values := make([]T, min((len(encoded)+7)/8, limit))
	for i := range values {
		var word [8]byte
		copy(word[:], encoded[8*i:])
		values[i] = decode(binary.LittleEndian.Uint64(word[:]))
	}
	return values
}

func FuzzZZHistogramDifferential(f *testing.F) {
	nan := math.Float64frombits(0xfff0000000000001)
	data := histogramFuzzBytes([]float64{nan, math.Inf(-1), -1, math.Copysign(0, -1), 0, 1,
		math.SmallestNonzeroFloat64, 2, math.Inf(1)})
	for _, edges := range [][]float64{nil, {}, {1}, {nan, 1}, {nan, nan, nan}, {2, 1, 0, 2, 1}} {
		f.Add(data, histogramFuzzBytes(edges))
	}
	f.Add([]byte{}, []byte{})
	for _, k := range []int{8, 16, 32, 64, 128, 256} {
		for _, m := range []int{k - 1, k, k + 1} {
			edges := make([]float64, m+1)
			for i := range edges {
				edges[i] = float64((m - i) / 2)
			}
			edges[m/2] = nan
			f.Add(data, histogramFuzzBytes(edges))
		}
	}
	f.Fuzz(func(t *testing.T, dataBytes, edgeBytes []byte) {
		oracleCheckHistogram(t, histogramFuzzValues(dataBytes, 64, math.Float64frombits),
			histogramFuzzValues(edgeBytes, 300, math.Float64frombits))
		float32Value := func(v uint64) float32 { return math.Float32frombits(uint32(v)) }
		oracleCheckHistogram(t, histogramFuzzValues(dataBytes, 64, float32Value),
			histogramFuzzValues(edgeBytes, 300, float32Value))
		intValue := func(v uint64) int64 { return int64(v) }
		oracleCheckHistogram(t, histogramFuzzValues(dataBytes, 64, intValue),
			histogramFuzzValues(edgeBytes, 300, intValue))
		uintValue := func(v uint64) uint64 { return v }
		oracleCheckHistogram(t, histogramFuzzValues(dataBytes, 64, uintValue),
			histogramFuzzValues(edgeBytes, 300, uintValue))
	})
}
