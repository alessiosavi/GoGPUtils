package a_test

import tst "testing"

func BenchmarkT_P(b *tst.B) {
	for b.Loop() {
	}
}
