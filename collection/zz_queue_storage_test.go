package collection

import (
	"strconv"
	"testing"
)

func boundQueue[T any](t *testing.T, q *Queue[T]) {
	t.Helper()
	if cap(q.buf) > max(64, 4*q.n) {
		t.Fatalf("capacity=%d, live=%d", cap(q.buf), q.n)
	}
}

func TestQueueFootprint(t *testing.T) {
	for _, size := range []int{1, 8, 63, 64, 65, 100, 1024, 65536} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			var q Queue[int]
			for i := range size {
				old := cap(q.buf)
				q.Enqueue(i)
				if cap(q.buf) > old && cap(q.buf) > 2*q.n+64 {
					t.Fatal("growth bound")
				}
				if old > 0 && q.n > old && cap(q.buf) != 2*old {
					t.Fatal("not explicit doubling")
				}
			}
			for !q.IsEmpty() {
				q.Dequeue()
				boundQueue(t, &q)
			}
			for i := range 16 {
				q.Enqueue(i)
			}
			boundQueue(t, &q)
			q.DequeueN(16)
			boundQueue(t, &q)
			q = *NewQueueWithCapacity[int](size)
			q.EnqueueAll(make([]int, size)...)
			q.DequeueN(size)
			if size > 64 {
				same(t, cap(q.buf), 0)
			} else {
				same(t, cap(q.buf), size)
			}
			q.EnqueueAll(make([]int, size)...)
			q.DequeueN(max(0, size-3))
			boundQueue(t, &q)
			q.EnqueueAll(make([]int, size)...)
			old := cap(q.buf)
			q.Clear()
			same(t, cap(q.buf), old)
			same(t, q.head, 0)
			same(t, q.n, 0)
		})
	}
	q := NewQueueWithCapacity[int](64)
	q.EnqueueAll(make([]int, 64)...)
	q.DequeueN(64)
	same(t, cap(q.buf), 64)
	q = NewQueueWithCapacity[int](1024)
	q.EnqueueAll(make([]int, 1024)...)
	q.DequeueN(769)
	same(t, cap(q.buf), 512)
	q.EnqueueAll(make([]int, 257)...)
	same(t, cap(q.buf), 512)
	q.Enqueue(9)
	same(t, cap(q.buf), 1024)
}

func TestQueueRetiredSegments(t *testing.T) {
	x := 7
	for _, method := range []string{"DequeueN", "Clear"} {
		q := NewQueueWithCapacity[*int](8)
		for range 6 {
			q.Enqueue(&x)
		}
		q.DequeueN(5)
		for range 6 {
			q.Enqueue(&x)
		} // live positions 5,6,7,0,1,2,3
		if method == "Clear" {
			q.Clear()
		} else {
			q.DequeueN(6)
		}
		for i, v := range q.buf {
			live := method == "DequeueN" && i == 3
			if !live && v != nil {
				t.Fatalf("%s: retired slot %d still live", method, i)
			}
		}
	}
}

func TestQueueHysteresis(t *testing.T) {
	q := NewQueueWithCapacity[int](128)
	q.EnqueueAll(make([]int, 65)...)
	for range 1000 {
		q.Enqueue(1)
		q.Dequeue()
		if cap(q.buf) != 128 {
			t.Fatal("steady-state thrash", cap(q.buf))
		}
	}
	q = NewQueueWithCapacity[int](1024)
	q.EnqueueAll(make([]int, 1024)...)
	q.DequeueN(769)
	same(t, cap(q.buf), 512)
	for range 1000 {
		q.Dequeue()
		q.Enqueue(1)
		same(t, cap(q.buf), 512)
	}
}

func TestQueueNoShrinkThrash(t *testing.T) {
	q := NewQueueWithCapacity[int](512)
	q.EnqueueAll(make([]int, 255)...)
	a := testing.AllocsPerRun(1000, func() { q.Dequeue(); q.Enqueue(1) })
	t.Logf("steady depth=255 capacity=512 AllocsPerRun=%g", a)
	if a != 0 {
		t.Fatalf("steady depth reallocates: %g allocs/op", a)
	}
}

func TestQueueWrappedRetirement(t *testing.T) {
	for _, method := range []string{"DequeueN", "Clear"} {
		t.Run(method, func(t *testing.T) {
			x := 7
			q := NewQueueWithCapacity[*int](8)
			q.head, q.n = 5, 7
			for _, i := range []int{5, 6, 7, 0, 1, 2, 3} {
				q.buf[i] = &x
			}
			if q.head != 5 || q.n != 7 {
				t.Fatal("fixture did not wrap")
			}
			if method == "Clear" {
				q.Clear()
			} else {
				q.DequeueN(6)
			}
			for i, v := range q.buf {
				if method == "DequeueN" && i == 3 {
					if v != &x {
						t.Fatal("live item cleared")
					}
					continue
				}
				if v != nil {
					t.Fatalf("%s retired slot %d still live", method, i)
				}
			}
		})
	}
}
