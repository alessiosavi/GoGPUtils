package p

import "testing"

func BenchmarkOne(b *testing.B) {
	for b.Loop() {
		One()
	}
}
