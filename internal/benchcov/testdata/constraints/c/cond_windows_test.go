package c

import "testing"

func benchCond(b *testing.B) {
	for b.Loop() {
		Cond()
	}
}
