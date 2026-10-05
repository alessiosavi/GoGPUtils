package sqs

import (
	"strconv"
	"testing"
)

var pr14BatchIDSink string

func BenchmarkPR14GenerateBatchID(b *testing.B) {
	for _, index := range []int{0, 9, 10, 99, 12345} {
		b.Run("index="+strconv.Itoa(index), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				pr14BatchIDSink = generateBatchID(index)
			}
		})
	}
}
