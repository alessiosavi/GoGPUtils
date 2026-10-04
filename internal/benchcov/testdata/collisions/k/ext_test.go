package k_test

import "testing"

func BenchmarkFoo(b *testing.B) {
	for b.Loop() {
	}
}
