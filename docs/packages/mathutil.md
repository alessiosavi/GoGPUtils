---
title: mathutil
parent: Packages
nav_order: 3
---

# mathutil

Mathematical and statistical operations on numeric slices using Go generics.
{: .fs-6 .fw-300 }

## Overview

The `mathutil` package provides a comprehensive collection of mathematical and statistical functions that work with any numeric type (`int`, `float64`, etc.) through Go generics. It is organized into logical categories: aggregations, extreme values, statistics, number theory, basic math, vector operations, matrix operations, and utility functions.

All functions handle edge cases gracefully, returning zero values or errors where appropriate rather than panicking.

```go
import "github.com/alessiosavi/GoGPUtils/mathutil"
```

---

## Type Constraints

The package defines three constraint interfaces for generic functions:

| Constraint | Types Included                                                                   |
| ---------- | -------------------------------------------------------------------------------- |
| `Number`   | All integer and floating-point types (`int`, `uint`, `float32`, `float64`, etc.) |
| `Integer`  | All integer types only (`int`, `uint`, `int64`, etc.)                            |
| `Float`    | Floating-point types only (`float32`, `float64`)                                 |

```go
type Number interface {
    ~int | ~int8 | ~int16 | ~int32 | ~int64 |
        ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
        ~float32 | ~float64
}
```

---

## Errors

| Error                  | Description                                           |
| ---------------------- | ----------------------------------------------------- |
| `ErrEmptySlice`        | Returned when an operation requires a non-empty slice |
| `ErrDivisionByZero`    | Returned on division by zero                          |
| `ErrInvalidDimensions` | Returned when matrix dimensions are incompatible      |
| `ErrNegativeInput`     | Returned when negative input is not allowed           |

---

## Aggregation Functions

### Sum

Returns the sum of all elements in the slice. Returns `0` for empty slices.

```go
func Sum[T Number](s []T) T
```

**Complexity:** O(n) time, O(1) space

```go
Sum([]int{1, 2, 3, 4, 5})        // 15
Sum([]float64{1.5, 2.5, 3.0})    // 7.0
Sum([]int{})                      // 0
```

---

### Product

Returns the product of all elements in the slice. Returns `1` for empty slices (multiplicative identity).

```go
func Product[T Number](s []T) T
```

**Complexity:** O(n) time, O(1) space

```go
Product([]int{1, 2, 3, 4})       // 24
Product([]int{1, 2, 0, 4})       // 0
Product([]int{})                  // 1
```

---

### Average

Returns the arithmetic mean of the slice. Returns `0` for empty slices.

```go
func Average[T Number](s []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
Average([]float64{1, 2, 3, 4, 5})  // 3.0
Average([]float64{42})             // 42.0
Average([]float64{})               // 0
```

---

## Extreme Value Functions

### Min

Returns the minimum value in the slice. Returns `ErrEmptySlice` for empty slices.

```go
func Min[T Number](s []T) (T, error)
```

**Complexity:** O(n) time, O(1) space

```go
min, err := Min([]int{3, 1, 4, 1, 5})
// min = 1, err = nil

_, err := Min([]int{})
// err = ErrEmptySlice
```

---

### Max

Returns the maximum value in the slice. Returns `ErrEmptySlice` for empty slices.

```go
func Max[T Number](s []T) (T, error)
```

**Complexity:** O(n) time, O(1) space

```go
max, err := Max([]int{3, 1, 4, 1, 5})
// max = 5, err = nil
```

---

### MinMax

Returns both minimum and maximum values in a single pass. Returns `ErrEmptySlice` for empty slices.

```go
func MinMax[T Number](s []T) (min, max T, err error)
```

**Complexity:** O(n) time, O(1) space — more efficient than calling `Min` and `Max` separately.

```go
min, max, err := MinMax([]int{3, 1, 4, 1, 5, 9, 2, 6})
// min = 1, max = 9, err = nil
```

---

### MinIndex

Returns the index of the minimum value. Returns `-1` for empty slices.

```go
func MinIndex[T Number](s []T) int
```

**Complexity:** O(n) time, O(1) space

```go
MinIndex([]int{3, 1, 4, 1, 5})   // 1
MinIndex([]int{5, 4, 3, 2, 1})   // 4
MinIndex([]int{})                 // -1
```

---

### MaxIndex

Returns the index of the maximum value. Returns `-1` for empty slices.

```go
func MaxIndex[T Number](s []T) int
```

**Complexity:** O(n) time, O(1) space

```go
MaxIndex([]int{3, 1, 4, 1, 5})   // 4
MaxIndex([]int{5, 4, 3, 2, 1})   // 0
MaxIndex([]int{})                 // -1
```

---

### Range

Returns the difference between max and min values (`max - min`). Returns `ErrEmptySlice` for empty slices.

```go
func Range[T Number](s []T) (T, error)
```

**Complexity:** O(n) time, O(1) space

```go
r, err := Range([]int{3, 1, 4, 1, 5, 9, 2, 6})
// r = 8 (9 - 1), err = nil
```

---

## Statistical Functions

### Median

Returns the middle value of a sorted slice. For even-length slices, returns the average of the two middle values. Returns `0` for empty slices.

```go
func Median[T Number](s []T) float64
```

**Complexity:** O(n log n) time (due to sorting), O(n) space

```go
Median([]int{1, 2, 3, 4, 5})     // 3.0
Median([]int{1, 2, 3, 4})        // 2.5
Median([]int{5, 1, 3, 2, 4})     // 3.0 (unsorted input is handled)
```

---

### Mode

Returns the most frequently occurring value(s). Returns `nil` for empty slices. May return multiple values if there are ties.

```go
func Mode[T Number](s []T) []T
```

**Complexity:** O(n) time, O(k) space where k is the number of unique values

```go
Mode([]int{1, 2, 2, 3, 3, 3})    // [3]
Mode([]int{1, 1, 2, 2})          // [1, 2] (tie)
Mode([]int{1, 2, 3})             // [1, 2, 3] (all same frequency)
```

---

### Variance

Returns the population variance (denominator = N). Returns `0` for empty slices.

```go
func Variance[T Number](s []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
Variance([]float64{1, 2, 3, 4, 5})  // 2.0
// Mean = 3, squared differences: 4+1+0+1+4 = 10, variance = 10/5 = 2
```

---

### SampleVariance

Returns the sample variance (denominator = N-1). Returns `0` for slices with fewer than 2 elements.

```go
func SampleVariance[T Number](s []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
SampleVariance([]float64{1, 2, 3, 4, 5})  // 2.5
// Same sum of squares = 10, sample variance = 10/4 = 2.5
```

---

### StdDev

Returns the population standard deviation (square root of population variance). Returns `0` for empty slices.

```go
func StdDev[T Number](s []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
StdDev([]float64{1, 2, 3, 4, 5})  // ~1.414 (sqrt(2))
```

---

### SampleStdDev

Returns the sample standard deviation (square root of sample variance). Returns `0` for slices with fewer than 2 elements.

```go
func SampleStdDev[T Number](s []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
SampleStdDev([]float64{1, 2, 3, 4, 5})  // ~1.581 (sqrt(2.5))
```

---

### Percentile

Returns the value at the given percentile (0-100). Uses linear interpolation between data points. Returns `0` for empty slices or invalid percentiles.

```go
func Percentile[T Number](s []T, p float64) float64
```

**Complexity:** O(n log n) time (due to sorting), O(n) space

```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
Percentile(data, 0)     // 1.0
Percentile(data, 50)    // 5.5 (median)
Percentile(data, 100)   // 10.0
```

---

### Quartiles

Returns Q1 (25th percentile), Q2 (median/50th percentile), and Q3 (75th percentile).

```go
func Quartiles[T Number](s []T) (q1, q2, q3 float64)
```

**Complexity:** O(n log n) time, O(n) space

```go
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
q1, q2, q3 := Quartiles(data)
// q1 = 3.25, q2 = 5.5, q3 = 7.75
```

---

### IQR

Returns the interquartile range (Q3 - Q1).

```go
func IQR[T Number](s []T) float64
```

**Complexity:** O(n log n) time, O(n) space

```go
IQR([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})  // 4.5 (7.75 - 3.25)
```

---

## Number Theory Functions

### IsPrime

Reports whether `n` is a prime number. Returns `false` for `n < 2`.

```go
func IsPrime[T Integer](n T) bool
```

**Complexity:** O(√n) time, O(1) space

```go
IsPrime(17)   // true
IsPrime(15)   // false
IsPrime(1)    // false
IsPrime(2)    // true
```

---

### Primes

Returns all prime numbers up to and including `n` using the Sieve of Eratosthenes.

```go
func Primes(n int) []int
```

**Complexity:** O(n log log n) time, O(n) space

```go
Primes(20)   // [2, 3, 5, 7, 11, 13, 17, 19]
Primes(1)    // nil
```

---

### GCD

Returns the greatest common divisor of `a` and `b` using the Euclidean algorithm.

```go
func GCD[T Integer](a, b T) T
```

**Complexity:** O(log(min(a, b))) time, O(1) space

```go
GCD(12, 18)    // 6
GCD(17, 13)    // 1 (coprime)
GCD(-12, 18)   // 6 (handles negatives)
GCD(0, 5)      // 5
```

---

### LCM

Returns the least common multiple of `a` and `b`.

```go
func LCM[T Integer](a, b T) T
```

**Complexity:** O(log(min(a, b))) time (dominated by GCD), O(1) space

```go
LCM(4, 6)     // 12
LCM(3, 5)     // 15
LCM(0, 5)     // 0
```

---

### Factorial

Returns `n!` (n factorial). Returns `1` for `n <= 1`. Returns `0` if the result would overflow `int64` (n > 20).

```go
func Factorial(n int) int64
```

**Complexity:** O(n) time, O(1) space

```go
Factorial(5)   // 120
Factorial(10)  // 3628800
Factorial(20)  // 2432902008176640000
Factorial(21)  // 0 (overflow)
```

---

### Fibonacci

Returns the nth Fibonacci number (0-indexed): F(0) = 0, F(1) = 1, F(n) = F(n-1) + F(n-2).

```go
func Fibonacci(n int) int64
```

**Complexity:** O(n) time, O(1) space — uses iterative approach

```go
Fibonacci(0)   // 0
Fibonacci(1)   // 1
Fibonacci(10)  // 55
Fibonacci(20)  // 6765
Fibonacci(-1)  // 0
```

---

## Basic Math Functions

### Abs

Returns the absolute value of `x`.

```go
func Abs[T Number](x T) T
```

**Complexity:** O(1) time, O(1) space

```go
Abs(5)    // 5
Abs(-5)   // 5
Abs(0)    // 0
```

---

### Sign

Returns `-1` for negative, `0` for zero, `1` for positive.

```go
func Sign[T Number](x T) int
```

**Complexity:** O(1) time, O(1) space

```go
Sign(5)    // 1
Sign(-5)   // -1
Sign(0)    // 0
```

---

### Clamp

Restricts `x` to the range `[min, max]`.

```go
func Clamp[T Number](x, min, max T) T
```

**Complexity:** O(1) time, O(1) space

```go
Clamp(15, 0, 10)   // 10
Clamp(-5, 0, 10)   // 0
Clamp(5, 0, 10)    // 5
```

---

### Sqrt

Returns the square root using Newton's method. Returns `0` for negative numbers.

```go
func Sqrt(x float64) float64
```

**Complexity:** O(1) time (fixed 100 iterations max), O(1) space

```go
Sqrt(4)     // 2.0
Sqrt(9)     // 3.0
Sqrt(2)     // ~1.414
Sqrt(-1)    // 0
```

---

### Pow

Returns `x` raised to the power of `n`. For negative exponents with integer types, returns `0`.

```go
func Pow[T Number](x T, n int) T
```

**Complexity:** O(n) time, O(1) space

```go
Pow(2, 3)    // 8
Pow(2, 0)    // 1
Pow(2, -1)   // 0 (integer types can't represent fractions)
```

---

### PowFloat

Returns `x` raised to the power of `n` for float types. Supports negative exponents.

```go
func PowFloat(x float64, n int) float64
```

**Complexity:** O(n) time, O(1) space

```go
PowFloat(2, 3)     // 8.0
PowFloat(2, -1)    // 0.5
PowFloat(2, 0)     // 1.0
```

---

## Vector Operations

### DotProduct

Computes the dot product of two vectors. Returns `0` if vectors have different lengths.

```go
func DotProduct[T Number](a, b []T) T
```

**Complexity:** O(n) time, O(1) space

```go
DotProduct([]int{1, 2, 3}, []int{4, 5, 6})   // 32
// 1*4 + 2*5 + 3*6 = 32

DotProduct([]int{1, 2}, []int{1})            // 0 (different lengths)
```

---

### Magnitude

Returns the Euclidean length (L2 norm) of a vector.

```go
func Magnitude[T Number](v []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
Magnitude([]float64{3, 4})   // 5.0
// sqrt(3² + 4²) = sqrt(9 + 16) = sqrt(25) = 5
```

---

### Normalize

Returns a unit vector in the same direction. Returns `nil` if the input is a zero vector.

```go
func Normalize[T Float](v []T) []T
```

**Complexity:** O(n) time, O(n) space

```go
Normalize([]float64{3, 4})   // [0.6, 0.8]
Normalize([]float64{0, 0})   // nil
```

---

### CosineSimilarity

Returns the cosine similarity between two vectors (range: -1 to 1). Returns `0` if vectors have different lengths or either is a zero vector.

```go
func CosineSimilarity[T Number](a, b []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
CosineSimilarity([]float64{1, 0}, []float64{0, 1})        // 0 (orthogonal)
CosineSimilarity([]float64{1, 2, 3}, []float64{2, 4, 6})  // 1 (parallel)
```

---

### EuclideanDistance

Returns the Euclidean (L2) distance between two points.

```go
func EuclideanDistance[T Number](a, b []T) float64
```

**Complexity:** O(n) time, O(1) space

```go
EuclideanDistance([]float64{0, 0}, []float64{3, 4})   // 5.0
```

---

### ManhattanDistance

Returns the Manhattan (L1) distance between two points.

```go
func ManhattanDistance[T Number](a, b []T) T
```

**Complexity:** O(n) time, O(1) space

```go
ManhattanDistance([]int{0, 0}, []int{3, 4})   // 7
// |3-0| + |4-0| = 7
```

---

## Matrix Operations

### Matrix Type

```go
type Matrix[T Number] [][]T
```

---

### MatrixMultiply

Multiplies two matrices. Returns `ErrInvalidDimensions` if dimensions don't match.

For matrices A (m×n) and B (n×p), the result is C (m×p).

```go
func MatrixMultiply[T Number](a, b Matrix[T]) (Matrix[T], error)
```

**Complexity:** O(m × n × p) time, O(m × p) space

```go
a := mathutil.Matrix[int]{ {1, 2}, {3, 4} }
b := mathutil.Matrix[int]{ {5, 6}, {7, 8} }
result, err := MatrixMultiply(a, b)
// result = [[19, 22], [43, 50]]
```

---

### MatrixTranspose

Returns the transpose of a matrix (rows become columns).

```go
func MatrixTranspose[T Number](m Matrix[T]) Matrix[T]
```

**Complexity:** O(m × n) time, O(m × n) space

```go
m := mathutil.Matrix[int]{ {1, 2, 3}, {4, 5, 6} }
result := MatrixTranspose(m)
// result = [[1, 4], [2, 5], [3, 6]]
```

---

### MatrixAdd

Adds two matrices element-wise. Returns `ErrInvalidDimensions` if dimensions don't match.

```go
func MatrixAdd[T Number](a, b Matrix[T]) (Matrix[T], error)
```

**Complexity:** O(m × n) time, O(m × n) space

```go
a := mathutil.Matrix[int]{ {1, 2}, {3, 4} }
b := mathutil.Matrix[int]{ {5, 6}, {7, 8} }
result, err := MatrixAdd(a, b)
// result = [[6, 8], [10, 12]]
```

---

### MatrixScalar

Multiplies a matrix by a scalar.

```go
func MatrixScalar[T Number](m Matrix[T], scalar T) Matrix[T]
```

**Complexity:** O(m × n) time, O(m × n) space

```go
m := mathutil.Matrix[int]{ {1, 2}, {3, 4} }
result := MatrixScalar(m, 2)
// result = [[2, 4], [6, 8]]
```

---

## Utility Functions

### Cumsum

Returns the cumulative sum of the slice.

```go
func Cumsum[T Number](s []T) []T
```

**Complexity:** O(n) time, O(n) space

```go
Cumsum([]int{1, 2, 3, 4})   // [1, 3, 6, 10]
```

---

### Diff

Returns the differences between consecutive elements.

```go
func Diff[T Number](s []T) []T
```

**Complexity:** O(n) time, O(n) space

```go
Diff([]int{1, 3, 6, 10})   // [2, 3, 4]
Diff([]int{1})              // nil
```

---

### Scale

Multiplies all elements by a scalar and returns a new slice.

```go
func Scale[T Number](s []T, scalar T) []T
```

**Complexity:** O(n) time, O(n) space

```go
Scale([]int{1, 2, 3}, 2)   // [2, 4, 6]
```

---

### Add

Performs element-wise addition of two slices. The shorter slice determines the result length.

```go
func Add[T Number](a, b []T) []T
```

**Complexity:** O(min(len(a), len(b))) time, O(min(len(a), len(b))) space

```go
Add([]int{1, 2, 3}, []int{4, 5, 6})   // [5, 7, 9]
```

---

### Subtract

Performs element-wise subtraction (`a - b`). The shorter slice determines the result length.

```go
func Subtract[T Number](a, b []T) []T
```

**Complexity:** O(min(len(a), len(b))) time, O(min(len(a), len(b))) space

```go
Subtract([]int{5, 7, 9}, []int{1, 2, 3})   // [4, 5, 6]
```

---

### Multiply

Performs element-wise multiplication. The shorter slice determines the result length.

```go
func Multiply[T Number](a, b []T) []T
```

**Complexity:** O(min(len(a), len(b))) time, O(min(len(a), len(b))) space

```go
Multiply([]int{1, 2, 3}, []int{4, 5, 6})   // [4, 10, 18]
```

---

### LinSpace

Returns `n` evenly spaced values from `start` to `end` (inclusive).

```go
func LinSpace(start, end float64, n int) []float64
```

**Complexity:** O(n) time, O(n) space

```go
LinSpace(0, 10, 5)    // [0, 2.5, 5, 7.5, 10]
LinSpace(0, 10, 1)    // [0]
```

---

### Arange

Returns values from `start` to `stop` (exclusive) with given `step`.

```go
func Arange[T Number](start, stop, step T) []T
```

**Complexity:** O(n) time, O(n) space where n = (stop-start)/step

```go
Arange(0, 10, 2)      // [0, 2, 4, 6, 8]
Arange(0, 10, 0)      // nil (zero step)
Arange(10, 0, -2)     // [10, 8, 6, 4, 2]
```

---

### MinMaxNormalize

Scales a slice to the range `[0.0, 1.0]` using min-max normalization:

```
result[i] = (s[i] - min) / (max - min)
```

Returns `nil` for nil or empty input. Returns a zero-filled slice when all elements are equal.

```go
func MinMaxNormalize[T Number](s []T) []float64
```

**Complexity:** O(n) time, O(n) space

```go
MinMaxNormalize([]int{0, 5, 10})    // [0.0, 0.5, 1.0]
MinMaxNormalize([]int{3, 3, 3})     // [0.0, 0.0, 0.0]
```

---

### Histogram

Returns the frequency count of values in bins. `bins` specifies the bin edges (n+1 edges for n bins).

```go
func Histogram[T Number](data []T, bins []T) []int
```

**Complexity:** O(m × n) time where m = len(data), n = len(bins)-1; O(n) space

```go
data := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
bins := []int{1, 2, 3, 4, 5}
Histogram(data, bins)   // [1, 2, 3, 4]
// Bin [1,2): 1 value
// Bin [2,3): 2 values
// Bin [3,4): 3 values
// Bin [4,5]: 4 values
```

---

## Complete Example

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/mathutil"
)

func main() {
    numbers := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Aggregations
    fmt.Println("Sum:", mathutil.Sum(numbers))           // 55
    fmt.Println("Average:", mathutil.Average(numbers))   // 5.5

    // Extremes
    min, max, _ := mathutil.MinMax(numbers)
    fmt.Println("Min:", min, "Max:", max)               // 1, 10

    // Statistics
    fmt.Println("Median:", mathutil.Median(numbers))     // 5.5
    fmt.Println("StdDev:", mathutil.StdDev(numbers))     // ~2.87

    // Number theory
    fmt.Println("GCD(48, 18):", mathutil.GCD(48, 18))    // 6
    fmt.Println("IsPrime(17):", mathutil.IsPrime(17))    // true

    // Vector operations
    a := []float64{1, 2, 3}
    b := []float64{4, 5, 6}
    fmt.Println("DotProduct:", mathutil.DotProduct(a, b)) // 32

    // Normalization
    normalized := mathutil.MinMaxNormalize([]int{0, 5, 10})
    fmt.Println("Normalized:", normalized)               // [0, 0.5, 1]
}
```

---

## Performance Notes

| Operation             | Time Complexity  | Space Complexity | Notes                               |
| --------------------- | ---------------- | ---------------- | ----------------------------------- |
| Sum, Product, Average | O(n)             | O(1)             | Single pass                         |
| Min, Max, MinMax      | O(n)             | O(1)             | MinMax does one pass instead of two |
| Median, Percentile    | O(n log n)       | O(n)             | Requires sorting                    |
| Mode                  | O(n)             | O(k)             | k = unique values                   |
| Variance, StdDev      | O(n)             | O(1)             | Two passes (mean + sum of squares)  |
| IsPrime               | O(√n)            | O(1)             | Checks odd divisors only            |
| Primes                | O(n log log n)   | O(n)             | Sieve of Eratosthenes               |
| GCD                   | O(log(min(a,b))) | O(1)             | Euclidean algorithm                 |
| Matrix Multiply       | O(m×n×p)         | O(m×p)           | Standard matrix multiplication      |
| Matrix Transpose      | O(m×n)           | O(m×n)           | Creates new matrix                  |
| DotProduct            | O(n)             | O(1)             | Single pass                         |
| CosineSimilarity      | O(n)             | O(1)             | Computes dot product + magnitudes   |
| Histogram             | O(m×n)           | O(n)             | m = data length, n = bins           |
