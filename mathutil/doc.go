// Package mathutil provides mathematical and statistical operations on numeric slices.
//
// The package uses Go generics to work with any numeric type (int, float64, etc.)
// through the Number constraint.
//
// # Aggregations
//
// Basic mathematical operations on slices:
//
//	sum := mathutil.Sum([]int{1, 2, 3, 4, 5})      // 15
//	product := mathutil.Product([]int{1, 2, 3})   // 6
//	avg := mathutil.Average([]float64{1, 2, 3})   // 2.0
//
// # Statistics
//
// Statistical functions:
//
//	median := mathutil.Median([]int{1, 2, 3, 4, 5})      // 3.0
//	mode := mathutil.Mode([]int{1, 2, 2, 3})             // [2]
//	stddev := mathutil.StdDev([]float64{1, 2, 3, 4, 5})  // ~1.41
//	variance := mathutil.Variance(data)
//
// # Extremes
//
// Finding minimum and maximum values:
//
//	min, max := mathutil.MinMax([]int{3, 1, 4, 1, 5})
//	minIdx := mathutil.MinIndex(data)
//	maxIdx := mathutil.MaxIndex(data)
//
// # Number Theory
//
// Prime numbers and related functions:
//
//	if mathutil.IsPrime(17) { ... }
//	primes := mathutil.Primes(100)
//	gcd := mathutil.GCD(12, 18)
//	lcm := mathutil.LCM(4, 6)
//
// # Matrix Operations
//
// Basic matrix operations:
//
//	result := mathutil.MatrixMultiply(a, b)
//	transposed := mathutil.MatrixTranspose(m)
//
// All functions are designed to handle edge cases gracefully,
// returning zero values or errors where appropriate.
package mathutil
