package cryptoutil

import (
	"strconv"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

var zzHashStringResult string

func TestZZHashStringAllocations(t *testing.T) {
	for _, n := range []int{16, 31, 32, 33, 511, 512, 513, 1023, 1024, 1025, 65536} {
		t.Run("n="+strconv.Itoa(n), func(t *testing.T) {
			text := benchkit.Text(n, benchkit.Mixed)
			baseAllocs, baseBytes := zzMeasureAllocations(func() { zzHashStringResult = baseHashString(text) })
			gotAllocs, gotBytes := zzMeasureAllocations(func() { zzHashStringResult = HashString(text) })
			t.Logf("allocs %.0f -> %.0f; bytes %d -> %d", baseAllocs, gotAllocs, baseBytes, gotBytes)
			if zzHashStringResult != baseHashString(text) {
				t.Fatal("hash differs from BASE")
			}
			if n <= 32 {
				if gotAllocs > baseAllocs || gotBytes > baseBytes {
					t.Fatal("short-input allocations and bytes must not exceed BASE")
				}
			} else if gotAllocs >= baseAllocs || gotBytes >= baseBytes {
				t.Fatal("streamed input must reduce allocations and bytes versus BASE")
			}
		})
	}
}
