package mathutil

import (
	"strconv"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkSum measures summing nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkSum(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Sum(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Sum(s)
				}
			})
		}},
	)
}

// BenchmarkProduct measures multiplying n elements close to one. Float inputs
// are 1+f/n for f from benchkit.Floats; conversion to int gives all ones.
// Both products stay finite and nonzero. Inputs are reused without restoration.
func BenchmarkProduct(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				values := benchkit.Floats(n)
				s := make([]int, n)
				for i, f := range values {
					s[i] = int(1 + f/float64(n))
				}
				b.ReportAllocs()
				for b.Loop() {
					Product(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				for i, f := range s {
					s[i] = 1 + f/float64(n)
				}
				b.ReportAllocs()
				for b.Loop() {
					Product(s)
				}
			})
		}},
	)
}

// BenchmarkAverage measures averaging nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkAverage(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Average(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Average(s)
				}
			})
		}},
	)
}

// BenchmarkMin measures finding the minimum of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMin(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Min(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Min(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMax measures finding the maximum of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMax(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Max(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Max(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMinMax measures finding both extrema of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMinMax(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, _, err := MinMax(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, _, err := MinMax(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMinIndex measures finding the minimum index of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMinIndex(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					MinIndex(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					MinIndex(s)
				}
			})
		}},
	)
}

// BenchmarkMaxIndex measures finding the maximum index of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMaxIndex(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					MaxIndex(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					MaxIndex(s)
				}
			})
		}},
	)
}

// BenchmarkRange measures finding the range of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkRange(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Range(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Range(s); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMedian measures copying, sorting, and finding the median of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMedian(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Median(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Median(s)
				}
			})
		}},
	)
}

// BenchmarkMode measures frequency counting and sorting the modes of n elements.
// Inputs are quantized to eight repeated values: integers [0,8), or floats
// 0/8 through 7/8. Read-only inputs are reused; map/result allocation is measured.
func BenchmarkMode(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 8
				}
				b.ReportAllocs()
				for b.Loop() {
					Mode(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				for i, f := range s {
					s[i] = float64(int(f*8)) / 8
				}
				b.ReportAllocs()
				for b.Loop() {
					Mode(s)
				}
			})
		}},
	)
}

// BenchmarkVariance measures population variance of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkVariance(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Variance(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Variance(s)
				}
			})
		}},
	)
}

// BenchmarkSampleVariance measures sample variance of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkSampleVariance(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					SampleVariance(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					SampleVariance(s)
				}
			})
		}},
	)
}

// BenchmarkStdDev measures population standard deviation of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkStdDev(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					StdDev(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					StdDev(s)
				}
			})
		}},
	)
}

// BenchmarkSampleStdDev measures sample standard deviation of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkSampleStdDev(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					SampleStdDev(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					SampleStdDev(s)
				}
			})
		}},
	)
}

// BenchmarkPercentile measures copying, sorting, and finding the 95th percentile of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkPercentile(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Percentile(s, 95)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Percentile(s, 95)
				}
			})
		}},
	)
}

// BenchmarkQuartiles measures three copy/sort percentile passes over nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkQuartiles(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Quartiles(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Quartiles(s)
				}
			})
		}},
	)
}

// BenchmarkIQR measures the interquartile range using three copy/sort passes over nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkIQR(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					IQR(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					IQR(s)
				}
			})
		}},
	)
}

// BenchmarkIsPrime measures primality for scalar domain cases: 17, 2147483647,
// and 2147483646. int64 keeps divisor squaring representable; no setup or
// restoration occurs in the measured loop.
func BenchmarkIsPrime(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int64", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value int64
			}{
				{"small-prime", 17},
				{"large-prime", 2147483647},
				{"composite", 2147483646},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						IsPrime(tc.value)
					}
				})
			}
		}},
	)
}

// BenchmarkPrimes measures sieving up to the inclusive limit n, with n from
// benchkit.Sizes. Each call allocates its own sieve and output; no restoration
// is needed.
func BenchmarkPrimes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {

		b.ReportAllocs()
		for b.Loop() {
			Primes(n)
		}
	})
}

// BenchmarkGCD measures GCD of two positive 32-bit values, 1000000007 and
// 1000000009. The int64 instantiation keeps their product representable.
// Arguments are reused; no restoration is needed.
func BenchmarkGCD(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int64", Fn: func(b *testing.B) {
			b.Run("case=32-bit-values", func(b *testing.B) {
				a, v := int64(1000000007), int64(1000000009)
				b.ReportAllocs()
				for b.Loop() {
					GCD(a, v)
				}
			})
		}},
	)
}

// BenchmarkLCM measures LCM of two positive 32-bit values, 1000000007 and
// 1000000009. The int64 instantiation keeps their product representable.
// Arguments are reused; no restoration is needed.
func BenchmarkLCM(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int64", Fn: func(b *testing.B) {
			b.Run("case=32-bit-values", func(b *testing.B) {
				a, v := int64(1000000007), int64(1000000009)
				b.ReportAllocs()
				for b.Loop() {
					LCM(a, v)
				}
			})
		}},
	)
}

// BenchmarkFactorial measures Factorial for representable int64 results; n is the
// factorial argument. Arguments are reused without restoration.
func BenchmarkFactorial(b *testing.B) {
	benchkit.Run(b, []int{5, 10, 20}, func(b *testing.B, n int) {

		b.ReportAllocs()
		for b.Loop() {
			Factorial(n)
		}
	})
}

// BenchmarkFibonacci measures Fibonacci for representable int64 results; n is the
// sequence index. Arguments are reused without restoration.
func BenchmarkFibonacci(b *testing.B) {
	benchkit.Run(b, []int{10, 50, 90}, func(b *testing.B, n int) {

		b.ReportAllocs()
		for b.Loop() {
			Fibonacci(n)
		}
	})
}

// BenchmarkAbs measures scalar negative, zero, and positive inputs (-42, 0, 42).
// Each case reuses its argument; no restoration is needed.
func BenchmarkAbs(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value int
			}{
				{"negative", -42},
				{"zero", 0},
				{"positive", 42},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Abs(tc.value)
					}
				})
			}
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value float64
			}{
				{"negative", -42},
				{"zero", 0},
				{"positive", 42},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Abs(tc.value)
					}
				})
			}
		}},
	)
}

// BenchmarkSign measures scalar negative, zero, and positive inputs (-42, 0, 42).
// Each case reuses its argument; no restoration is needed.
func BenchmarkSign(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value int
			}{
				{"negative", -42},
				{"zero", 0},
				{"positive", 42},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Sign(tc.value)
					}
				})
			}
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value float64
			}{
				{"negative", -42},
				{"zero", 0},
				{"positive", 42},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Sign(tc.value)
					}
				})
			}
		}},
	)
}

// BenchmarkClamp measures scalar inputs below, inside, and above [0,10] (-5, 5, 15).
// Each case reuses its argument; no restoration is needed.
func BenchmarkClamp(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value int
			}{
				{"below", -5},
				{"inside", 5},
				{"above", 15},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Clamp(tc.value, 0, 10)
					}
				})
			}
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			for _, tc := range []struct {
				name  string
				value float64
			}{
				{"below", -5},
				{"inside", 5},
				{"above", 15},
			} {
				b.Run("case="+tc.name, func(b *testing.B) {
					b.ReportAllocs()
					for b.Loop() {
						Clamp(tc.value, 0, 10)
					}
				})
			}
		}},
	)
}

// BenchmarkSqrt measures square roots of finite scalar inputs: 0.25, 2, and
// 1e100. Each case reuses its argument; no restoration is needed.
func BenchmarkSqrt(b *testing.B) {
	for _, tc := range []struct {
		name  string
		value float64
	}{
		{"fraction", 0.25},
		{"regular", 2},
		{"large", 1e100},
	} {
		b.Run("case="+tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				Sqrt(tc.value)
			}
		})
	}
}

// BenchmarkPow measures repeated multiplication with base two; n is the
// exponent (2, 10, or 30). Results remain representable for int and float64.
// Arguments are reused without restoration.
func BenchmarkPow(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, []int{2, 10, 30}, func(b *testing.B, n int) {
				base := int(2)
				b.ReportAllocs()
				for b.Loop() {
					Pow(base, n)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, []int{2, 10, 30}, func(b *testing.B, n int) {
				base := float64(2)
				b.ReportAllocs()
				for b.Loop() {
					Pow(base, n)
				}
			})
		}},
	)
}

// BenchmarkPowFloat measures repeated multiplication with base 1.01; n is the
// exponent (-10, 2, 10, or 30). Negative n includes the reciprocal operation.
// Arguments are reused without restoration.
func BenchmarkPowFloat(b *testing.B) {
	benchkit.Run(b, []int{-10, 2, 10, 30}, func(b *testing.B, n int) {
		base := 1.01
		b.ReportAllocs()
		for b.Loop() {
			PowFloat(base, n)
		}
	})
}

// BenchmarkDotProduct measures DotProduct on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkDotProduct(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					DotProduct(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					DotProduct(a, v)
				}
			})
		}},
	)
}

// BenchmarkMagnitude measures the Euclidean magnitude of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMagnitude(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Magnitude(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Magnitude(s)
				}
			})
		}},
	)
}

// BenchmarkNormalize measures allocating a unit vector from n float64 elements
// in [0,1). The nonzero input is reused read-only; no restoration is needed.
func BenchmarkNormalize(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Normalize(s)
				}
			})
		}},
	)
}

// BenchmarkCosineSimilarity measures CosineSimilarity on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkCosineSimilarity(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					CosineSimilarity(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					CosineSimilarity(a, v)
				}
			})
		}},
	)
}

// BenchmarkEuclideanDistance measures EuclideanDistance on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkEuclideanDistance(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					EuclideanDistance(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					EuclideanDistance(a, v)
				}
			})
		}},
	)
}

// BenchmarkManhattanDistance measures ManhattanDistance on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkManhattanDistance(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					ManhattanDistance(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					ManhattanDistance(a, v)
				}
			})
		}},
	)
}

// BenchmarkMatrixMultiply measures MatrixMultiply on square float64 matrices; n is the
// dimension (n rows by n columns), from benchkit.SmallSizes. Compatible
// read-only inputs contain [0,1) values. Result allocation is measured;
// no restoration is needed.
func BenchmarkMatrixMultiply(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.SmallSizes, func(b *testing.B, n int) {
				a := Matrix[float64](benchkit.Matrix(n, n))
				v := Matrix[float64](benchkit.Matrix(n, n))
				b.ReportAllocs()
				for b.Loop() {
					if _, err := MatrixMultiply(a, v); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMatrixTranspose measures MatrixTranspose on square float64 matrices; n is the
// dimension (n rows by n columns), from benchkit.SmallSizes. Compatible
// read-only inputs contain [0,1) values. Result allocation is measured;
// no restoration is needed.
func BenchmarkMatrixTranspose(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.SmallSizes, func(b *testing.B, n int) {
				a := Matrix[float64](benchkit.Matrix(n, n))
				b.ReportAllocs()
				for b.Loop() {
					MatrixTranspose(a)
				}
			})
		}},
	)
}

// BenchmarkMatrixAdd measures MatrixAdd on square float64 matrices; n is the
// dimension (n rows by n columns), from benchkit.SmallSizes. Compatible
// read-only inputs contain [0,1) values. Result allocation is measured;
// no restoration is needed.
func BenchmarkMatrixAdd(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.SmallSizes, func(b *testing.B, n int) {
				a := Matrix[float64](benchkit.Matrix(n, n))
				v := Matrix[float64](benchkit.Matrix(n, n))
				b.ReportAllocs()
				for b.Loop() {
					if _, err := MatrixAdd(a, v); err != nil {
						b.Fatal(err)
					}
				}
			})
		}},
	)
}

// BenchmarkMatrixScalar measures MatrixScalar on square float64 matrices; n is the
// dimension (n rows by n columns), from benchkit.SmallSizes. Compatible
// read-only inputs contain [0,1) values and the scalar is two. Result allocation is measured;
// no restoration is needed.
func BenchmarkMatrixScalar(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.SmallSizes, func(b *testing.B, n int) {
				a := Matrix[float64](benchkit.Matrix(n, n))
				b.ReportAllocs()
				for b.Loop() {
					MatrixScalar(a, 2)
				}
			})
		}},
	)
}

// BenchmarkCumsum measures allocating cumulative sums of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkCumsum(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Cumsum(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Cumsum(s)
				}
			})
		}},
	)
}

// BenchmarkDiff measures allocating consecutive differences of nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkDiff(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Diff(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Diff(s)
				}
			})
		}},
	)
}

// BenchmarkScale measures allocating values scaled by two from nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkScale(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					Scale(s, 2)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					Scale(s, 2)
				}
			})
		}},
	)
}

// BenchmarkAdd measures Add on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkAdd(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Add(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Add(a, v)
				}
			})
		}},
	)
}

// BenchmarkSubtract measures Subtract on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkSubtract(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Subtract(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Subtract(a, v)
				}
			})
		}},
	)
}

// BenchmarkMultiply measures Multiply on two equal-length read-only inputs; n
// is the element count per input. Integer a is in [0,128), with v=127-a;
// float a is in [0,1), with v=1-a. No restoration is needed.
func BenchmarkMultiply(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Ints(n)
				for i := range a {
					a[i] %= 128
				}
				v := make([]int, n)
				for i := range v {
					v[i] = 127 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Multiply(a, v)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				a := benchkit.Floats(n)
				v := make([]float64, n)
				for i := range v {
					v[i] = 1 - a[i]
				}
				b.ReportAllocs()
				for b.Loop() {
					Multiply(a, v)
				}
			})
		}},
	)
}

// BenchmarkLinSpace measures allocating n equally spaced output values over
// [0,1], including both endpoints. Arguments are reused without restoration.
func BenchmarkLinSpace(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {

		b.ReportAllocs()
		for b.Loop() {
			LinSpace(0, 1, n)
		}
	})
}

// BenchmarkArange measures allocating n output values over [0,n) with step one.
// The float64 arguments and successive values are exact integers. Arguments
// are reused without restoration.
func BenchmarkArange(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				stop := n
				b.ReportAllocs()
				for b.Loop() {
					Arange(int(0), stop, int(1))
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				stop := float64(n)
				b.ReportAllocs()
				for b.Loop() {
					Arange(float64(0), stop, float64(1))
				}
			})
		}},
	)
}

// BenchmarkMinMaxNormalize measures allocating min-max normalization of nonconstant nonempty slices;
// n is the number of input elements. Read-only inputs are reused; ints
// are in [0,128), floats in [0,1). No restoration is measured.
func BenchmarkMinMaxNormalize(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 128
				}
				b.ReportAllocs()
				for b.Loop() {
					MinMaxNormalize(s)
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				b.ReportAllocs()
				for b.Loop() {
					MinMaxNormalize(s)
				}
			})
		}},
	)
}

// BenchmarkHistogram measures binning n input elements into bins=8 or bins=256,
// using 9 or 257 sorted edges spanning [0,65536] for int and [0,1] for
// float64. Inputs are uniform within that range and reused read-only. The
// API's edge copy/sort and count allocation are measured; no restoration is needed.
func BenchmarkHistogram(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 65536
				}
				for _, bins := range []int{8, 256} {
					b.Run("bins="+strconv.Itoa(bins), func(b *testing.B) {
						edges := make([]int, bins+1)
						for i := range edges {
							edges[i] = i * (65536 / bins)
						}
						b.ReportAllocs()
						for b.Loop() {
							Histogram(s, edges)
						}
					})
				}
			})
		}},
		benchkit.TypeCase{Name: "float64", Fn: func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				for _, bins := range []int{8, 256} {
					b.Run("bins="+strconv.Itoa(bins), func(b *testing.B) {
						edges := make([]float64, bins+1)
						for i := range edges {
							edges[i] = float64(i) / float64(bins)
						}
						b.ReportAllocs()
						for b.Loop() {
							Histogram(s, edges)
						}
					})
				}
			})
		}},
	)
}
