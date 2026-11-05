package minhash

import (
	"fmt"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

func benchShingles() []uint64 {
	xs := make([]uint64, 300)
	for i := range xs {
		xs[i] = uint64(i) * 7919
	}
	return xs
}

func BenchmarkShingles(b *testing.B) {
	for _, n := range benchkit.Sizes {
		b.Run(fmt.Sprintf("n=%d/k=5", n), func(b *testing.B) {
			tokens := benchkit.Words(n)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := Shingles(tokens, 5); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for _, n := range []int{64, 128, 1024} {
		b.Run(fmt.Sprintf("perm=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := New(n, 42); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkHasher_NumPerm(b *testing.B) {
	h := mustNew(b, 64, 42)
	b.ReportAllocs()
	for b.Loop() {
		h.NumPerm()
	}
}

func BenchmarkHasher_Sign(b *testing.B) {
	for _, n := range []int{64, 128} {
		b.Run(fmt.Sprintf("perm=%d/n=300", n), func(b *testing.B) {
			h, xs := mustNew(b, n, 42), benchShingles()
			b.ReportAllocs()
			for b.Loop() {
				if _, err := h.Sign(xs); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkHasher_SignInto(b *testing.B) {
	for _, n := range []int{64, 128} {
		b.Run(fmt.Sprintf("perm=%d/n=300", n), func(b *testing.B) {
			h, xs := mustNew(b, n, 42), benchShingles()
			dst := make(Signature, n)
			b.ReportAllocs()
			for b.Loop() {
				if err := h.SignInto(dst, xs); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSimilarity(b *testing.B) {
	h := mustNew(b, 128, 42)
	a, c := mustSign(b, h, benchShingles()), mustSign(b, h, benchShingles()[100:])
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Similarity(a, c); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJaccard(b *testing.B) {
	a := benchShingles()
	c := a[100:]
	b.ReportAllocs()
	for b.Loop() {
		Jaccard(a, c)
	}
}

func BenchmarkBandKeys(b *testing.B) {
	for _, n := range []int{64, 128} {
		b.Run(fmt.Sprintf("perm=%d/bands=%d", n, n/4), func(b *testing.B) {
			sig := mustSign(b, mustNew(b, n, 42), benchShingles())
			b.ReportAllocs()
			for b.Loop() {
				if _, err := BandKeys(sig, n/4); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkNewIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := NewIndex(16, 4); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkIndex_Add measures one successful Add to an empty index. A new
// index is constructed outside timing each iteration so state stays bounded.
func BenchmarkIndex_Add(b *testing.B) {
	sig := mustSign(b, mustNew(b, 64, 42), benchShingles())
	b.ReportAllocs()
	for b.Loop() {
		b.StopTimer()
		idx, err := NewIndex(16, 4)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		if err := idx.Add("document", sig); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIndex_Candidates(b *testing.B) {
	idx, err := NewIndex(16, 4)
	if err != nil {
		b.Fatal(err)
	}
	sig := mustSign(b, mustNew(b, 64, 42), benchShingles())
	for i := range 100 {
		if err := idx.Add(fmt.Sprint(i), sig); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := idx.Candidates(sig); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBands(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, _, err := Bands(64, 0.85); err != nil {
			b.Fatal(err)
		}
	}
}
