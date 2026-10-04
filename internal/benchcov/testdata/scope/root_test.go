package scope

import "testing"

func BenchmarkRootFn(b *testing.B) {
	for b.Loop() {
		RootFn()
	}
}
