package c

import "testing"

func BenchmarkPlain(b *testing.B) {
	for b.Loop() {
		Plain()
	}
}

func BenchmarkTagged(b *testing.B) {
	for b.Loop() {
	}
}

// BenchmarkCond delegates to a helper defined per build condition.
func BenchmarkCond(b *testing.B) { benchCond(b) }
