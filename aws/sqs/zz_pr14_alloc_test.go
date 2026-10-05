package sqs

import (
	"fmt"
	"testing"
)

var pr14BatchIDAllocSink string

func TestPR14GenerateBatchIDAllocations(t *testing.T) {
	for _, tc := range []struct {
		index int
		want  float64
	}{
		{0, 0},
		{9, 0},
		{10, 0},
		{99, 0},
		{12345, 1},
	} {
		t.Run(fmt.Sprintf("index=%d", tc.index), func(t *testing.T) {
			base := testing.AllocsPerRun(1000, func() {
				pr14BatchIDAllocSink = baseGenerateBatchID(tc.index)
			})
			got := testing.AllocsPerRun(1000, func() {
				pr14BatchIDAllocSink = generateBatchID(tc.index)
			})
			t.Logf("allocations: BASE=%g candidate=%g", base, got)
			if got > tc.want || got > base {
				t.Errorf("allocations = %g, want <= %g and <= BASE %g", got, tc.want, base)
			}
		})
	}
}
