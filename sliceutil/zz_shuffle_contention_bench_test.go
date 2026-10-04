package sliceutil

import (
	"sync/atomic"
	"testing"
)

type shuffleStruct64 struct{ words [8]uint64 }

// Mean elapsed time for a copy and small shuffle while a large shuffler runs.
// Background completions are sampled strictly inside the timed window, before
// signaling stop; warm-up and teardown completions do not enter the metric.
func benchmarkShuffleContended[T any](b *testing.B, value func(int) T) {
	large, src, work := make([]T, 65536), make([]T, 16), make([]T, 16)
	for i := range large {
		large[i] = value(i)
	}
	for i := range src {
		src[i] = value(i)
	}
	var completed atomic.Uint64
	var stop atomic.Bool
	ready, done := make(chan struct{}), make(chan struct{})
	SeedShuffle(42)
	go func() {
		defer close(done)
		ShuffleInPlace(large)
		close(ready)
		for !stop.Load() {
			ShuffleInPlace(large)
			completed.Add(1)
		}
	}()
	defer func() { stop.Store(true); <-done }()
	<-ready
	b.ResetTimer()
	before := completed.Load()
	for range b.N {
		copy(work, src)
		ShuffleInPlace(work)
	}
	after := completed.Load()
	b.StopTimer()
	// Zero is printed for the caller to reject, rather than fabricated by
	// waiting/yielding in the measured small-operation loop.
	b.ReportMetric(float64(after-before)/float64(b.N), "bg-shuffles/op")
}

func BenchmarkShuffleContended(b *testing.B) {
	b.Run("elem=int", func(b *testing.B) {
		b.Run("n=16", func(b *testing.B) {
			benchmarkShuffleContended(b, func(i int) int { return i })
		})
	})
	b.Run("elem=struct64", func(b *testing.B) {
		b.Run("n=16", func(b *testing.B) {
			benchmarkShuffleContended(b, func(i int) shuffleStruct64 {
				var v shuffleStruct64
				for j := range v.words {
					v.words[j] = uint64(i*8 + j)
				}
				return v
			})
		})
	})
}
