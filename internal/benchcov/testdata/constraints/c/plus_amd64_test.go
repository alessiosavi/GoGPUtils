package c

import "testing"

func BenchmarkPlus(b *testing.B) {
	for b.Loop() {
		Plus()
	}
}
