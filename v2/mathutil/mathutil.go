package mathutil

import (
	"errors"
	"slices"
	"sort"
)

// Number is a constraint for all numeric types.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// Integer is a constraint for all integer types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float is a constraint for floating-point types.
type Float interface {
	~float32 | ~float64
}

// Errors returned by math operations.
var (
	ErrEmptySlice        = errors.New("empty slice")
	ErrDivisionByZero    = errors.New("division by zero")
	ErrInvalidDimensions = errors.New("invalid matrix dimensions")
	ErrNegativeInput     = errors.New("negative input not allowed")
)

// ============================================================================
// Aggregation Functions
// ============================================================================

// Sum returns the sum of all elements in the slice.
// Returns 0 for empty slices.
//
// Example:
//
//	Sum([]int{1, 2, 3, 4, 5})  // 15
func Sum[T Number](s []T) T {
	var sum T
	for _, v := range s {
		sum += v
	}
	return sum
}

// Product returns the product of all elements in the slice.
// Returns 1 for empty slices.
//
// Example:
//
//	Product([]int{1, 2, 3, 4})  // 24
func Product[T Number](s []T) T {
	if len(s) == 0 {
		return 1
	}
	product := s[0]
	for _, v := range s[1:] {
		product *= v
	}
	return product
}

// Average returns the arithmetic mean of the slice.
// Returns 0 for empty slices.
//
// Example:
//
//	Average([]float64{1, 2, 3, 4, 5})  // 3.0
func Average[T Number](s []T) float64 {
	if len(s) == 0 {
		return 0
	}
	return float64(Sum(s)) / float64(len(s))
}

// ============================================================================
// Extreme Value Functions
// ============================================================================

// Min returns the minimum value in the slice.
// Returns zero value and ErrEmptySlice for empty slices.
func Min[T Number](s []T) (T, error) {
	if len(s) == 0 {
		var zero T
		return zero, ErrEmptySlice
	}
	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}
	return min, nil
}

// Max returns the maximum value in the slice.
// Returns zero value and ErrEmptySlice for empty slices.
func Max[T Number](s []T) (T, error) {
	if len(s) == 0 {
		var zero T
		return zero, ErrEmptySlice
	}
	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}
	return max, nil
}

// MinMax returns both minimum and maximum values in a single pass.
// Returns zero values and ErrEmptySlice for empty slices.
func MinMax[T Number](s []T) (min, max T, err error) {
	if len(s) == 0 {
		return min, max, ErrEmptySlice
	}
	min, max = s[0], s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max, nil
}

// MinIndex returns the index of the minimum value.
// Returns -1 for empty slices.
func MinIndex[T Number](s []T) int {
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

// MaxIndex returns the index of the maximum value.
// Returns -1 for empty slices.
func MaxIndex[T Number](s []T) int {
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

// Range returns the difference between max and min values.
// Returns 0 and ErrEmptySlice for empty slices.
func Range[T Number](s []T) (T, error) {
	min, max, err := MinMax(s)
	if err != nil {
		var zero T
		return zero, err
	}
	return max - min, nil
}

// ============================================================================
// Statistical Functions
// ============================================================================

// Median returns the middle value of a sorted slice.
// For even-length slices, returns the average of the two middle values.
// Returns 0 for empty slices.
//
// Example:
//
//	Median([]int{1, 2, 3, 4, 5})  // 3.0
//	Median([]int{1, 2, 3, 4})    // 2.5
func Median[T Number](s []T) float64 {
	if len(s) == 0 {
		return 0
	}

	// Sort a copy
	sorted := make([]T, len(s))
	copy(sorted, s)
	slices.Sort(sorted)

	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return (float64(sorted[mid-1]) + float64(sorted[mid])) / 2
	}
	return float64(sorted[mid])
}

// Mode returns the most frequently occurring values.
// Returns nil for empty slices.
// May return multiple values if there are ties.
//
// Example:
//
//	Mode([]int{1, 2, 2, 3, 3, 3})  // [3]
//	Mode([]int{1, 1, 2, 2})       // [1, 2] (tie)
func Mode[T Number](s []T) []T {
	if len(s) == 0 {
		return nil
	}

	counts := make(map[T]int)
	maxCount := 0
	for _, v := range s {
		counts[v]++
		if counts[v] > maxCount {
			maxCount = counts[v]
		}
	}

	var modes []T
	for v, count := range counts {
		if count == maxCount {
			modes = append(modes, v)
		}
	}

	// Sort for deterministic output
	slices.Sort(modes)
	return modes
}

// Variance returns the population variance of the slice.
// Returns 0 for empty slices.
//
// For sample variance (N-1 denominator), use SampleVariance.
func Variance[T Number](s []T) float64 {
	if len(s) == 0 {
		return 0
	}

	avg := Average(s)
	var sumSquares float64
	for _, v := range s {
		diff := float64(v) - avg
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(s))
}

// SampleVariance returns the sample variance (using N-1 denominator).
// Returns 0 for slices with fewer than 2 elements.
func SampleVariance[T Number](s []T) float64 {
	if len(s) < 2 {
		return 0
	}

	avg := Average(s)
	var sumSquares float64
	for _, v := range s {
		diff := float64(v) - avg
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(s)-1)
}

// StdDev returns the population standard deviation.
// Returns 0 for empty slices.
func StdDev[T Number](s []T) float64 {
	return Sqrt(Variance(s))
}

// SampleStdDev returns the sample standard deviation.
// Returns 0 for slices with fewer than 2 elements.
func SampleStdDev[T Number](s []T) float64 {
	return Sqrt(SampleVariance(s))
}

// Percentile returns the value at the given percentile (0-100).
// Uses linear interpolation between data points.
// Returns 0 for empty slices.
//
// Example:
//
//	Percentile([]int{1, 2, 3, 4, 5}, 50)  // 3.0 (median)
func Percentile[T Number](s []T, p float64) float64 {
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

// Quartiles returns Q1, Q2 (median), and Q3.
func Quartiles[T Number](s []T) (q1, q2, q3 float64) {
	q1 = Percentile(s, 25)
	q2 = Percentile(s, 50)
	q3 = Percentile(s, 75)
	return
}

// IQR returns the interquartile range (Q3 - Q1).
func IQR[T Number](s []T) float64 {
	q1, _, q3 := Quartiles(s)
	return q3 - q1
}

// ============================================================================
// Number Theory Functions
// ============================================================================

// IsPrime reports whether n is a prime number.
// Returns false for n < 2.
//
// Example:
//
//	IsPrime(17)  // true
//	IsPrime(15)  // false
func IsPrime[T Integer](n T) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Check odd divisors up to sqrt(n)
	i := T(3)
	for i*i <= n {
		if n%i == 0 {
			return false
		}
		i += 2
	}
	return true
}

// Primes returns all prime numbers up to and including n.
// Uses the Sieve of Eratosthenes.
//
// Example:
//
//	Primes(20)  // [2, 3, 5, 7, 11, 13, 17, 19]
func Primes(n int) []int {
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

// GCD returns the greatest common divisor of a and b using Euclidean algorithm.
//
// Example:
//
//	GCD(12, 18)  // 6
func GCD[T Integer](a, b T) T {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM returns the least common multiple of a and b.
//
// Example:
//
//	LCM(4, 6)  // 12
func LCM[T Integer](a, b T) T {
	if a == 0 || b == 0 {
		return 0
	}
	return Abs(a*b) / GCD(a, b)
}

// Factorial returns n! (n factorial).
// Returns 1 for n <= 1.
// Returns 0 if result would overflow int64.
func Factorial(n int) int64 {
	if n <= 1 {
		return 1
	}
	if n > 20 { // 21! overflows int64
		return 0
	}
	result := int64(1)
	for i := 2; i <= n; i++ {
		result *= int64(i)
	}
	return result
}

// Fibonacci returns the nth Fibonacci number (0-indexed).
// F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2)
//
// Example:
//
//	Fibonacci(10)  // 55
func Fibonacci(n int) int64 {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return int64(n)
	}

	a, b := int64(0), int64(1)
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// ============================================================================
// Basic Math Functions
// ============================================================================

// Abs returns the absolute value of x.
func Abs[T Number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

// Sign returns -1 for negative, 0 for zero, 1 for positive.
func Sign[T Number](x T) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

// Clamp restricts x to the range [min, max].
//
// Example:
//
//	Clamp(15, 0, 10)  // 10
//	Clamp(-5, 0, 10)  // 0
//	Clamp(5, 0, 10)   // 5
func Clamp[T Number](x, min, max T) T {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// Sqrt returns the square root using Newton's method.
// Returns 0 for negative numbers.
func Sqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x == 0 || x == 1 {
		return x
	}

	guess := x / 2
	for i := 0; i < 100; i++ { // Sufficient iterations for precision
		newGuess := (guess + x/guess) / 2
		if absFloat(newGuess-guess) < 1e-15 {
			break
		}
		guess = newGuess
	}
	return guess
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Pow returns x raised to the power of n.
// For negative exponents with integer types, returns 0.
func Pow[T Number](x T, n int) T {
	if n == 0 {
		return 1
	}
	if n < 0 {
		return 0 // Integer types can't represent fractions
	}

	result := x
	for i := 1; i < n; i++ {
		result *= x
	}
	return result
}

// PowFloat returns x raised to the power of n for float types.
func PowFloat(x float64, n int) float64 {
	if n == 0 {
		return 1
	}
	negative := n < 0
	if negative {
		n = -n
	}

	result := 1.0
	for i := 0; i < n; i++ {
		result *= x
	}

	if negative {
		return 1 / result
	}
	return result
}

// ============================================================================
// Vector Operations
// ============================================================================

// DotProduct computes the dot product of two vectors.
// Returns 0 if vectors have different lengths.
func DotProduct[T Number](a, b []T) T {
	if len(a) != len(b) {
		var zero T
		return zero
	}

	var sum T
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Magnitude returns the Euclidean length of a vector.
func Magnitude[T Number](v []T) float64 {
	var sumSquares float64
	for _, x := range v {
		sumSquares += float64(x) * float64(x)
	}
	return Sqrt(sumSquares)
}

// Normalize returns a unit vector in the same direction.
// Returns nil if the input is a zero vector.
func Normalize[T Float](v []T) []T {
	mag := Magnitude(v)
	if mag == 0 {
		return nil
	}

	result := make([]T, len(v))
	for i, x := range v {
		result[i] = T(float64(x) / mag)
	}
	return result
}

// CosineSimilarity returns the cosine similarity between two vectors.
// Returns 0 if vectors have different lengths or either is zero vector.
func CosineSimilarity[T Number](a, b []T) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	dot := float64(DotProduct(a, b))
	magA := Magnitude(a)
	magB := Magnitude(b)

	if magA == 0 || magB == 0 {
		return 0
	}

	return dot / (magA * magB)
}

// EuclideanDistance returns the Euclidean distance between two points.
func EuclideanDistance[T Number](a, b []T) float64 {
	if len(a) != len(b) {
		return 0
	}

	var sumSquares float64
	for i := range a {
		diff := float64(a[i]) - float64(b[i])
		sumSquares += diff * diff
	}
	return Sqrt(sumSquares)
}

// ManhattanDistance returns the Manhattan (L1) distance between two points.
func ManhattanDistance[T Number](a, b []T) T {
	if len(a) != len(b) {
		var zero T
		return zero
	}

	var sum T
	for i := range a {
		sum += Abs(a[i] - b[i])
	}
	return sum
}

// ============================================================================
// Matrix Operations
// ============================================================================

// Matrix represents a 2D matrix.
type Matrix[T Number] [][]T

// MatrixMultiply multiplies two matrices.
// Returns nil and ErrInvalidDimensions if dimensions don't match.
//
// Matrix A (m×n) × Matrix B (n×p) = Matrix C (m×p)
func MatrixMultiply[T Number](a, b Matrix[T]) (Matrix[T], error) {
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
		for j := 0; j < p; j++ {
			var sum T
			for k := 0; k < n; k++ {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result, nil
}

// MatrixTranspose returns the transpose of a matrix.
func MatrixTranspose[T Number](m Matrix[T]) Matrix[T] {
	if len(m) == 0 {
		return nil
	}

	rows := len(m)
	cols := len(m[0])

	result := make(Matrix[T], cols)
	for i := range result {
		result[i] = make([]T, rows)
		for j := 0; j < rows; j++ {
			result[i][j] = m[j][i]
		}
	}
	return result
}

// MatrixAdd adds two matrices element-wise.
func MatrixAdd[T Number](a, b Matrix[T]) (Matrix[T], error) {
	if len(a) != len(b) || len(a) == 0 {
		return nil, ErrInvalidDimensions
	}
	if len(a[0]) != len(b[0]) {
		return nil, ErrInvalidDimensions
	}

	result := make(Matrix[T], len(a))
	for i := range result {
		result[i] = make([]T, len(a[i]))
		for j := range result[i] {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result, nil
}

// MatrixScalar multiplies a matrix by a scalar.
func MatrixScalar[T Number](m Matrix[T], scalar T) Matrix[T] {
	if len(m) == 0 {
		return nil
	}

	result := make(Matrix[T], len(m))
	for i := range result {
		result[i] = make([]T, len(m[i]))
		for j := range result[i] {
			result[i][j] = m[i][j] * scalar
		}
	}
	return result
}

// ============================================================================
// Utility Functions
// ============================================================================

// Cumsum returns the cumulative sum of the slice.
//
// Example:
//
//	Cumsum([]int{1, 2, 3, 4})  // [1, 3, 6, 10]
func Cumsum[T Number](s []T) []T {
	if s == nil {
		return nil
	}

	result := make([]T, len(s))
	var sum T
	for i, v := range s {
		sum += v
		result[i] = sum
	}
	return result
}

// Diff returns the differences between consecutive elements.
//
// Example:
//
//	Diff([]int{1, 3, 6, 10})  // [2, 3, 4]
func Diff[T Number](s []T) []T {
	if len(s) < 2 {
		return nil
	}

	result := make([]T, len(s)-1)
	for i := 0; i < len(s)-1; i++ {
		result[i] = s[i+1] - s[i]
	}
	return result
}

// Scale multiplies all elements by a scalar and returns a new slice.
func Scale[T Number](s []T, scalar T) []T {
	if s == nil {
		return nil
	}

	result := make([]T, len(s))
	for i, v := range s {
		result[i] = v * scalar
	}
	return result
}

// Add performs element-wise addition of two slices.
// The shorter slice determines the result length.
func Add[T Number](a, b []T) []T {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]T, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] + b[i]
	}
	return result
}

// Subtract performs element-wise subtraction (a - b).
func Subtract[T Number](a, b []T) []T {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]T, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] - b[i]
	}
	return result
}

// Multiply performs element-wise multiplication.
func Multiply[T Number](a, b []T) []T {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}

	result := make([]T, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] * b[i]
	}
	return result
}

// LinSpace returns n evenly spaced values from start to end (inclusive).
func LinSpace(start, end float64, n int) []float64 {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []float64{start}
	}

	result := make([]float64, n)
	step := (end - start) / float64(n-1)
	for i := range result {
		result[i] = start + float64(i)*step
	}
	return result
}

// Arange returns values from start to stop (exclusive) with given step.
func Arange[T Number](start, stop, step T) []T {
	if step == 0 {
		return nil
	}
	if (step > 0 && start >= stop) || (step < 0 && start <= stop) {
		return nil
	}

	var result []T
	for v := start; (step > 0 && v < stop) || (step < 0 && v > stop); v += step {
		result = append(result, v)
	}
	return result
}

// Histogram returns the frequency count of values in bins.
// bins specifies the bin edges (n+1 edges for n bins).
func Histogram[T Number](data []T, bins []T) []int {
	if len(bins) < 2 {
		return nil
	}

	// Sort bins
	sortedBins := make([]T, len(bins))
	copy(sortedBins, bins)
	sort.Slice(sortedBins, func(i, j int) bool { return sortedBins[i] < sortedBins[j] })

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
