package collection

import (
	"runtime"
	"slices"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

func pr10Balanced(lo, hi int, out *[]int) {
	if lo >= hi {
		return
	}
	m := (lo + hi) / 2
	*out = append(*out, m)
	pr10Balanced(lo, m, out)
	pr10Balanced(m+1, hi, out)
}
func BenchmarkBST_LevelOrderShape(b *testing.B) {
	b.Run("shape=balanced", func(b *testing.B) {
		benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
			var vals []int
			pr10Balanced(0, n, &vals)
			tree := NewBSTFrom(vals)
			b.ReportAllocs()
			for b.Loop() {
				tree.LevelOrder()
			}
		})
	})
	b.Run("shape=skewed", func(b *testing.B) {
		benchkit.Run(b, benchkit.SmallSizes, func(b *testing.B, n int) {
			vals := make([]int, n)
			for i := range vals {
				vals[i] = i
			}
			tree := NewBSTFrom(vals)
			b.ReportAllocs()
			for b.Loop() {
				tree.LevelOrder()
			}
		})
	})
}
func BenchmarkQueue_BurstDrainRefill(b *testing.B) {
	benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
		items := benchkit.Ints(n)
		q := NewQueue[int]()
		b.ReportAllocs()
		for b.Loop() {
			q.EnqueueAll(items...)
			for range n {
				q.Dequeue()
			}
			for i := range 16 {
				q.Enqueue(i)
			}
			for range 16 {
				q.Dequeue()
			}
		}
	})
}
func pr10Pointers(n int) []*int {
	v := make([]int, n)
	p := make([]*int, n)
	for i := range v {
		v[i] = i
		p[i] = &v[i]
	}
	return p
}
func BenchmarkStack_ClearPtr(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := pr10Pointers(n)
		s := NewStackWithCapacity[*int](n)
		b.ReportAllocs()
		for b.Loop() {
			s.PushAll(items...)
			if s.Len() != n {
				b.Fatal("bad refill")
			}
			s.Clear()
		}
	})
}
func BenchmarkQueue_ClearPtr(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		items := pr10Pointers(n)
		q := NewQueueWithCapacity[*int](n)
		b.ReportAllocs()
		for b.Loop() {
			q.EnqueueAll(items...)
			if q.Len() != n {
				b.Fatal("bad refill")
			}
			q.Clear()
		}
	})
}

// The heap harness is self-contained in the additive file, and uses only the
// exported collection API. Separate frames keep temporary payload roots out of
// the measurement frame. The container stays live through the second heap read.
//
//go:noinline
func pr10HeapRetiredStack(method string) any {
	s := NewStackWithCapacity[*[1 << 20]byte](8)
	p := new([1 << 20]byte)
	p[0] = 1
	s.Push(nil)
	s.Push(p)
	switch method {
	case "Pop":
		s.Pop()
	case "PopN":
		s.PopN(1)
	case "PopAll":
		s.PopAll()
	case "Clear":
		s.Clear()
	}
	return s
}

//go:noinline
func pr10HeapRetiredQueue(method string) any {
	q := NewQueueWithCapacity[*[1 << 20]byte](8)
	p := new([1 << 20]byte)
	p[0] = 1
	q.Enqueue(p)
	q.Enqueue(nil)
	switch method {
	case "Dequeue":
		q.Dequeue()
	case "DequeueN":
		q.DequeueN(1)
	case "Clear":
		q.Clear()
	}
	return q
}

//go:noinline
func pr10HeapBurstDrainRefill() any {
	q := NewQueue[int]()
	q.EnqueueAll(make([]int, 65536)...)
	for range 65536 {
		q.Dequeue()
	}
	for i := range 16 {
		q.Enqueue(i)
	}
	return q
}

func TestPR10RetainedHeap(t *testing.T) {
	if !testing.Verbose() {
		t.Skip("informational retained-heap harness; run with -v; TestRetention holds the retention assertion")
	}
	cases := []struct {
		name     string
		makeLive func() any
	}{}
	for _, method := range []string{"Pop", "PopN", "PopAll", "Clear"} {
		cases = append(cases, struct {
			name     string
			makeLive func() any
		}{
			"Stack/" + method, func() any { return pr10HeapRetiredStack(method) },
		})
	}
	for _, method := range []string{"Dequeue", "DequeueN", "Clear"} {
		cases = append(cases, struct {
			name     string
			makeLive func() any
		}{
			"Queue/" + method, func() any { return pr10HeapRetiredQueue(method) },
		})
	}
	cases = append(cases, struct {
		name     string
		makeLive func() any
	}{
		"Queue/BurstDrainRefill", pr10HeapBurstDrainRefill,
	})
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var deltas []int64
			for range 5 {
				runtime.GC()
				runtime.GC()
				var before, after runtime.MemStats
				runtime.ReadMemStats(&before)
				live := c.makeLive()
				runtime.GC()
				runtime.GC()
				runtime.ReadMemStats(&after)
				runtime.KeepAlive(live)
				deltas = append(deltas, int64(after.HeapAlloc)-int64(before.HeapAlloc))
			}
			slices.Sort(deltas)
			t.Logf("HEAP %s samples=%v median=%d spread=[%d,%d]\n",
				c.name, deltas, deltas[2], deltas[0], deltas[4])
		})
	}
}
