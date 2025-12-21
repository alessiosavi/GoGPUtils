package mathutil

import (
	"errors"
	"slices"
	"testing"
)

// ============================================================================
// Aggregation Tests
// ============================================================================

func TestSum(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{"positive", []int{1, 2, 3, 4, 5}, 15},
		{"with negative", []int{-1, 2, -3, 4}, 2},
		{"single", []int{42}, 42},
		{"empty", []int{}, 0},
		{"nil", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sum(tt.input); got != tt.want {
				t.Errorf("Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumFloat(t *testing.T) {
	got := Sum([]float64{1.5, 2.5, 3.0})

	want := 7.0

	if got != want {
		t.Errorf("Sum() = %v, want %v", got, want)
	}
}

func TestProduct(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
	}{
		{"positive", []int{1, 2, 3, 4}, 24},
		{"with zero", []int{1, 2, 0, 4}, 0},
		{"single", []int{5}, 5},
		{"empty", []int{}, 1}, // Product of empty set is 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Product(tt.input); got != tt.want {
				t.Errorf("Product() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  float64
	}{
		{"normal", []float64{1, 2, 3, 4, 5}, 3.0},
		{"single", []float64{42}, 42.0},
		{"empty", []float64{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Average(tt.input)
			if got != tt.want {
				t.Errorf("Average() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Extreme Value Tests
// ============================================================================

func TestMin(t *testing.T) {
	got, err := Min([]int{3, 1, 4, 1, 5})
	if err != nil || got != 1 {
		t.Errorf("Min() = %v, %v; want 1, nil", got, err)
	}

	_, err = Min([]int{})
	if !errors.Is(err, ErrEmptySlice) {
		t.Errorf("Min(empty) error = %v, want ErrEmptySlice", err)
	}
}

func TestMax(t *testing.T) {
	got, err := Max([]int{3, 1, 4, 1, 5})
	if err != nil || got != 5 {
		t.Errorf("Max() = %v, %v; want 5, nil", got, err)
	}

	_, err = Max([]int{})
	if !errors.Is(err, ErrEmptySlice) {
		t.Errorf("Max(empty) error = %v, want ErrEmptySlice", err)
	}
}

func TestMinMax(t *testing.T) {
	min, max, err := MinMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
	if err != nil || min != 1 || max != 9 {
		t.Errorf("MinMax() = %v, %v, %v; want 1, 9, nil", min, max, err)
	}
}

func TestMinIndex(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{3, 1, 4, 1, 5}, 1},
		{[]int{5, 4, 3, 2, 1}, 4},
		{[]int{42}, 0},
		{[]int{}, -1},
	}

	for _, tt := range tests {
		if got := MinIndex(tt.input); got != tt.want {
			t.Errorf("MinIndex(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMaxIndex(t *testing.T) {
	tests := []struct {
		input []int
		want  int
	}{
		{[]int{3, 1, 4, 1, 5}, 4},
		{[]int{5, 4, 3, 2, 1}, 0},
		{[]int{42}, 0},
		{[]int{}, -1},
	}

	for _, tt := range tests {
		if got := MaxIndex(tt.input); got != tt.want {
			t.Errorf("MaxIndex(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestRange(t *testing.T) {
	got, err := Range([]int{3, 1, 4, 1, 5, 9, 2, 6})
	if err != nil || got != 8 {
		t.Errorf("Range() = %v, %v; want 8, nil", got, err)
	}
}

// ============================================================================
// Statistical Tests
// ============================================================================

func TestMedian(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  float64
	}{
		{"odd length", []int{1, 2, 3, 4, 5}, 3.0},
		{"even length", []int{1, 2, 3, 4}, 2.5},
		{"single", []int{42}, 42.0},
		{"unsorted", []int{5, 1, 3, 2, 4}, 3.0},
		{"empty", []int{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Median(tt.input)
			if got != tt.want {
				t.Errorf("Median() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMode(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"single mode", []int{1, 2, 2, 3}, []int{2}},
		{"multiple modes", []int{1, 1, 2, 2}, []int{1, 2}},
		{"all same", []int{5, 5, 5}, []int{5}},
		{"all different", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty", []int{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Mode(tt.input)
			if !slices.Equal(got, tt.want) {
				t.Errorf("Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariance(t *testing.T) {
	// Variance of [1, 2, 3, 4, 5] with mean 3
	// Sum of squared differences: 4+1+0+1+4 = 10
	// Population variance: 10/5 = 2
	data := []float64{1, 2, 3, 4, 5}
	got := Variance(data)

	want := 2.0

	if got != want {
		t.Errorf("Variance() = %v, want %v", got, want)
	}
}

func TestSampleVariance(t *testing.T) {
	// Sample variance: 10/4 = 2.5
	data := []float64{1, 2, 3, 4, 5}
	got := SampleVariance(data)

	want := 2.5

	if got != want {
		t.Errorf("SampleVariance() = %v, want %v", got, want)
	}
}

func TestStdDev(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	got := StdDev(data)
	// sqrt(2) ≈ 1.414
	if got < 1.41 || got > 1.42 {
		t.Errorf("StdDev() = %v, want ~1.414", got)
	}
}

func TestPercentile(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	tests := []struct {
		p    float64
		want float64
	}{
		{0, 1.0},
		{50, 5.5},
		{100, 10.0},
	}

	for _, tt := range tests {
		got := Percentile(data, tt.p)
		if got != tt.want {
			t.Errorf("Percentile(%d) = %v, want %v", int(tt.p), got, tt.want)
		}
	}
}

func TestQuartiles(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	q1, q2, q3 := Quartiles(data)

	// Using linear interpolation with (n-1)*p indexing
	// Q1: (10-1)*0.25 = 2.25, interpolate between index 2 and 3: 3 + 0.25*(4-3) = 3.25
	// Q2: (10-1)*0.50 = 4.5, interpolate between index 4 and 5: 5 + 0.5*(6-5) = 5.5
	// Q3: (10-1)*0.75 = 6.75, interpolate between index 6 and 7: 7 + 0.75*(8-7) = 7.75
	if q1 != 3.25 {
		t.Errorf("Q1 = %v, want 3.25", q1)
	}

	if q2 != 5.5 {
		t.Errorf("Q2 = %v, want 5.5", q2)
	}

	if q3 != 7.75 {
		t.Errorf("Q3 = %v, want 7.75", q3)
	}
}

func TestIQR(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := IQR(data)

	want := 4.5 // 7.75 - 3.25

	if got != want {
		t.Errorf("IQR() = %v, want %v", got, want)
	}
}

// ============================================================================
// Number Theory Tests
// ============================================================================

func TestIsPrime(t *testing.T) {
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}
	notPrimes := []int{0, 1, 4, 6, 8, 9, 10, 12, 15, 100}

	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("IsPrime(%d) = false, want true", n)
		}
	}

	for _, n := range notPrimes {
		if IsPrime(n) {
			t.Errorf("IsPrime(%d) = true, want false", n)
		}
	}
}

func TestPrimes(t *testing.T) {
	got := Primes(20)

	want := []int{2, 3, 5, 7, 11, 13, 17, 19}

	if !slices.Equal(got, want) {
		t.Errorf("Primes(20) = %v, want %v", got, want)
	}

	if Primes(1) != nil {
		t.Error("Primes(1) should return nil")
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{12, 18, 6},
		{17, 13, 1},
		{100, 25, 25},
		{0, 5, 5},
		{-12, 18, 6},
	}

	for _, tt := range tests {
		got := GCD(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{4, 6, 12},
		{3, 5, 15},
		{12, 18, 36},
		{0, 5, 0},
	}

	for _, tt := range tests {
		got := LCM(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("LCM(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		n    int
		want int64
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{10, 3628800},
		{20, 2432902008176640000},
		{21, 0}, // Overflow
	}

	for _, tt := range tests {
		got := Factorial(tt.n)
		if got != tt.want {
			t.Errorf("Factorial(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    int
		want int64
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
		{20, 6765},
		{-1, 0},
	}

	for _, tt := range tests {
		got := Fibonacci(tt.n)
		if got != tt.want {
			t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// ============================================================================
// Basic Math Tests
// ============================================================================

func TestAbs(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		if got := Abs(tt.input); got != tt.want {
			t.Errorf("Abs(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestSign(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{5, 1},
		{-5, -1},
		{0, 0},
	}

	for _, tt := range tests {
		if got := Sign(tt.input); got != tt.want {
			t.Errorf("Sign(%d) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		x, min, max int
		want        int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
		{5, 5, 5, 5},
	}

	for _, tt := range tests {
		got := Clamp(tt.x, tt.min, tt.max)
		if got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.x, tt.min, tt.max, got, tt.want)
		}
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{4, 2},
		{9, 3},
		{2, 1.41421356},
		{0, 0},
		{1, 1},
		{-1, 0}, // Negative returns 0
	}

	for _, tt := range tests {
		got := Sqrt(tt.input)
		if absFloat(got-tt.want) > 0.00001 {
			t.Errorf("Sqrt(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		x, n int
		want int
	}{
		{2, 3, 8},
		{2, 0, 1},
		{5, 1, 5},
		{3, 4, 81},
		{2, -1, 0}, // Negative exponent returns 0 for integers
	}

	for _, tt := range tests {
		got := Pow(tt.x, tt.n)
		if got != tt.want {
			t.Errorf("Pow(%d, %d) = %d, want %d", tt.x, tt.n, got, tt.want)
		}
	}
}

func TestPowFloat(t *testing.T) {
	tests := []struct {
		x    float64
		n    int
		want float64
	}{
		{2, 3, 8},
		{2, -1, 0.5},
		{2, 0, 1},
	}

	for _, tt := range tests {
		got := PowFloat(tt.x, tt.n)
		if absFloat(got-tt.want) > 0.00001 {
			t.Errorf("PowFloat(%v, %d) = %v, want %v", tt.x, tt.n, got, tt.want)
		}
	}
}

// ============================================================================
// Vector Tests
// ============================================================================

func TestDotProduct(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	got := DotProduct(a, b)

	want := 32 // 1*4 + 2*5 + 3*6

	if got != want {
		t.Errorf("DotProduct() = %d, want %d", got, want)
	}

	// Different lengths
	if DotProduct([]int{1, 2}, []int{1}) != 0 {
		t.Error("DotProduct with different lengths should return 0")
	}
}

func TestMagnitude(t *testing.T) {
	v := []float64{3, 4}
	got := Magnitude(v)

	want := 5.0

	if got != want {
		t.Errorf("Magnitude() = %v, want %v", got, want)
	}
}

func TestNormalize(t *testing.T) {
	v := []float64{3, 4}
	got := Normalize(v)
	// Should be [0.6, 0.8]
	if absFloat(got[0]-0.6) > 0.00001 || absFloat(got[1]-0.8) > 0.00001 {
		t.Errorf("Normalize() = %v, want [0.6, 0.8]", got)
	}

	// Zero vector
	if Normalize([]float64{0, 0}) != nil {
		t.Error("Normalize zero vector should return nil")
	}
}

func TestCosineSimilarity(t *testing.T) {
	a := []float64{1, 0}
	b := []float64{0, 1}

	got := CosineSimilarity(a, b)
	if got != 0 {
		t.Errorf("Cosine similarity of orthogonal vectors = %v, want 0", got)
	}

	// Same direction
	got = CosineSimilarity([]float64{1, 2, 3}, []float64{2, 4, 6})
	if got < 0.99999 {
		t.Errorf("Cosine similarity of parallel vectors = %v, want 1", got)
	}
}

func TestEuclideanDistance(t *testing.T) {
	a := []float64{0, 0}
	b := []float64{3, 4}
	got := EuclideanDistance(a, b)

	want := 5.0

	if got != want {
		t.Errorf("EuclideanDistance() = %v, want %v", got, want)
	}
}

func TestManhattanDistance(t *testing.T) {
	a := []int{0, 0}
	b := []int{3, 4}
	got := ManhattanDistance(a, b)

	want := 7 // |3-0| + |4-0|

	if got != want {
		t.Errorf("ManhattanDistance() = %d, want %d", got, want)
	}
}

// ============================================================================
// Matrix Tests
// ============================================================================

func TestMatrixMultiply(t *testing.T) {
	a := Matrix[int]{{1, 2}, {3, 4}}
	b := Matrix[int]{{5, 6}, {7, 8}}

	got, err := MatrixMultiply(a, b)
	if err != nil {
		t.Fatalf("MatrixMultiply() error = %v", err)
	}

	want := Matrix[int]{{19, 22}, {43, 50}}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("MatrixMultiply() = %v, want %v", got, want)

			break
		}
	}
}

func TestMatrixMultiplyInvalidDimensions(t *testing.T) {
	a := Matrix[int]{{1, 2, 3}}
	b := Matrix[int]{{1, 2}} // 1x2, needs 3 rows

	_, err := MatrixMultiply(a, b)
	if !errors.Is(err, ErrInvalidDimensions) {
		t.Errorf("Expected ErrInvalidDimensions, got %v", err)
	}
}

func TestMatrixTranspose(t *testing.T) {
	m := Matrix[int]{{1, 2, 3}, {4, 5, 6}}
	got := MatrixTranspose(m)
	want := Matrix[int]{{1, 4}, {2, 5}, {3, 6}}

	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("MatrixTranspose() = %v, want %v", got, want)

			break
		}
	}
}

func TestMatrixAdd(t *testing.T) {
	a := Matrix[int]{{1, 2}, {3, 4}}
	b := Matrix[int]{{5, 6}, {7, 8}}

	got, err := MatrixAdd(a, b)
	if err != nil {
		t.Fatalf("MatrixAdd() error = %v", err)
	}

	want := Matrix[int]{{6, 8}, {10, 12}}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("MatrixAdd() = %v, want %v", got, want)

			break
		}
	}
}

func TestMatrixScalar(t *testing.T) {
	m := Matrix[int]{{1, 2}, {3, 4}}
	got := MatrixScalar(m, 2)
	want := Matrix[int]{{2, 4}, {6, 8}}

	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("MatrixScalar() = %v, want %v", got, want)

			break
		}
	}
}

// ============================================================================
// Utility Function Tests
// ============================================================================

func TestCumsum(t *testing.T) {
	got := Cumsum([]int{1, 2, 3, 4})

	want := []int{1, 3, 6, 10}

	if !slices.Equal(got, want) {
		t.Errorf("Cumsum() = %v, want %v", got, want)
	}
}

func TestDiff(t *testing.T) {
	got := Diff([]int{1, 3, 6, 10})

	want := []int{2, 3, 4}

	if !slices.Equal(got, want) {
		t.Errorf("Diff() = %v, want %v", got, want)
	}

	// Too short
	if Diff([]int{1}) != nil {
		t.Error("Diff of single element should return nil")
	}
}

func TestScale(t *testing.T) {
	got := Scale([]int{1, 2, 3}, 2)

	want := []int{2, 4, 6}

	if !slices.Equal(got, want) {
		t.Errorf("Scale() = %v, want %v", got, want)
	}
}

func TestAdd(t *testing.T) {
	got := Add([]int{1, 2, 3}, []int{4, 5, 6})

	want := []int{5, 7, 9}

	if !slices.Equal(got, want) {
		t.Errorf("Add() = %v, want %v", got, want)
	}
}

func TestSubtract(t *testing.T) {
	got := Subtract([]int{5, 7, 9}, []int{1, 2, 3})

	want := []int{4, 5, 6}

	if !slices.Equal(got, want) {
		t.Errorf("Subtract() = %v, want %v", got, want)
	}
}

func TestMultiply(t *testing.T) {
	got := Multiply([]int{1, 2, 3}, []int{4, 5, 6})

	want := []int{4, 10, 18}

	if !slices.Equal(got, want) {
		t.Errorf("Multiply() = %v, want %v", got, want)
	}
}

func TestLinSpace(t *testing.T) {
	got := LinSpace(0, 10, 5)

	want := []float64{0, 2.5, 5, 7.5, 10}

	if len(got) != len(want) {
		t.Errorf("LinSpace() = %v, want %v", got, want)
	}

	for i := range got {
		if absFloat(got[i]-want[i]) > 0.00001 {
			t.Errorf("LinSpace()[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestArange(t *testing.T) {
	got := Arange(0, 10, 2)

	want := []int{0, 2, 4, 6, 8}

	if !slices.Equal(got, want) {
		t.Errorf("Arange() = %v, want %v", got, want)
	}

	// Zero step
	if Arange(0, 10, 0) != nil {
		t.Error("Arange with zero step should return nil")
	}
}

func TestHistogram(t *testing.T) {
	data := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	bins := []int{1, 2, 3, 4, 5}
	got := Histogram(data, bins)

	want := []int{1, 2, 3, 4}

	if !slices.Equal(got, want) {
		t.Errorf("Histogram() = %v, want %v", got, want)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkSum(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	for b.Loop() {
		Sum(data)
	}
}

func BenchmarkAverage(b *testing.B) {
	data := make([]float64, 10000)
	for i := range data {
		data[i] = float64(i)
	}

	for b.Loop() {
		Average(data)
	}
}

func BenchmarkMedian(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	for b.Loop() {
		Median(data)
	}
}

func BenchmarkIsPrime(b *testing.B) {
	for b.Loop() {
		IsPrime(104729) // 10000th prime
	}
}

func BenchmarkPrimes(b *testing.B) {
	for b.Loop() {
		Primes(10000)
	}
}

func BenchmarkMatrixMultiply(b *testing.B) {
	size := 100
	a := make(Matrix[int], size)

	bMat := make(Matrix[int], size)

	for i := range size {
		a[i] = make([]int, size)

		bMat[i] = make([]int, size)

		for j := range size {
			a[i][j] = i + j
			bMat[i][j] = i - j
		}
	}

	for b.Loop() {
		MatrixMultiply(a, bMat)
	}
}

func BenchmarkCosineSimilarity(b *testing.B) {
	a := make([]float64, 1000)
	bVec := make([]float64, 1000)

	for i := range a {
		a[i] = float64(i)
		bVec[i] = float64(i * 2)
	}

	for b.Loop() {
		CosineSimilarity(a, bVec)
	}
}
