package collection

import (
	"runtime"
	"testing"
	"weak"
)

type retentionPayload [1 << 20]byte

// A separate non-inlined frame removes temporary strong pointers and discarded results.
// Capacity 8 and a nil sentinel make the negative control retain a non-zero backing slice.
//
//go:noinline
func retiredStack(method string) (*Stack[*retentionPayload], weak.Pointer[retentionPayload]) {
	s := NewStackWithCapacity[*retentionPayload](8)
	p := new(retentionPayload)
	p[0] = 1
	w := weak.Make(p)
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
	return s, w
}

//go:noinline
func retiredQueue(method string) (*Queue[*retentionPayload], weak.Pointer[retentionPayload]) {
	q := NewQueueWithCapacity[*retentionPayload](8)
	p := new(retentionPayload)
	p[0] = 1
	w := weak.Make(p)
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
	return q, w
}
func TestRetention(t *testing.T) {
	for _, method := range []string{"Pop", "PopN", "PopAll", "Clear"} {
		t.Run("Stack/"+method, func(t *testing.T) {
			s, w := retiredStack(method)
			runtime.GC()
			runtime.GC()
			alive := w.Value() != nil
			runtime.KeepAlive(s)
			if alive {
				t.Fatal("retired payload still reachable")
			}
		})
	}
	for _, method := range []string{"Dequeue", "DequeueN", "Clear"} {
		t.Run("Queue/"+method, func(t *testing.T) {
			q, w := retiredQueue(method)
			runtime.GC()
			runtime.GC()
			alive := w.Value() != nil
			runtime.KeepAlive(q)
			if alive {
				t.Fatal("retired payload still reachable")
			}
		})
	}
}
