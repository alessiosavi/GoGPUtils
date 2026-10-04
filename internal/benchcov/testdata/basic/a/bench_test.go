package a

import "testing"

func BenchmarkFoo(b *testing.B) {
	for b.Loop() {
		Foo()
	}
}

func BenchmarkT_M(b *testing.B) {
	var t T
	for b.Loop() {
		t.M()
	}
}

func BenchmarkG_Get(b *testing.B) {
	var g G[string, int]
	for b.Loop() {
		g.Get()
	}
}

// BenchmarkT_MP is a composite and covers neither T.M nor T.P.
func BenchmarkT_MP(b *testing.B) {
	var t T
	for b.Loop() {
		t.M()
		t.P()
	}
}

// BenchmarkBaseline is an extra benchmark, which is allowed.
func BenchmarkBaseline(b *testing.B) {
	for b.Loop() {
		_ = helper()
	}
}

// Benchmarkbar is not a benchmark: a lowercase letter follows the prefix.
func Benchmarkbar(b *testing.B) {}

// BenchmarkBar has the wrong signature, so Bar stays uncovered.
func BenchmarkBar(t *testing.T) {}

// BenchmarkS_Len is generic, so it is not a benchmark.
func BenchmarkS_Len[E any](b *testing.B) {}
