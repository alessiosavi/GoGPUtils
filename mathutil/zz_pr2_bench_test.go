package mathutil

import (
	"strconv"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkPR2Histogram fills in the crossover bins between BenchmarkHistogram's
// endpoints, with identical data, edge construction, and measured operations.
func BenchmarkPR2Histogram(b *testing.B) {
	benchkit.Each(b,
		benchkit.TypeCase{Name: "int", Fn: func(b *testing.B) {
			benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
				s := benchkit.Ints(n)
				for i := range s {
					s[i] %= 65536
				}
				for _, bins := range []int{16, 32, 64, 128} {
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
			benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
				s := benchkit.Floats(n)
				for _, bins := range []int{16, 32, 64, 128} {
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
