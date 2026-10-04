package c

import "testing"

func BenchmarkLinuxOnly(b *testing.B) {
	for b.Loop() {
		LinuxOnly()
	}
}
