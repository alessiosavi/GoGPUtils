package k

import "testing"

func BenchmarkFoo(b *testing.B) {
	for b.Loop() {
		Foo()
	}
}

func BenchmarkT_M(b *testing.B) {
	for b.Loop() {
		T_M()
	}
}
