package collection

import (
	"slices"
	"testing"
)

func TestQueue_ZeroValue(t *testing.T) {
	t.Run("int", func(t *testing.T) { checkZeroValueQueue(t, 1, 2) })
	t.Run("pointer", func(t *testing.T) {
		a, b := 1, 2
		checkZeroValueQueue(t, &a, &b)
	})
}

func checkZeroValueQueue[T comparable](t *testing.T, a, b T) {
	t.Helper()
	var q Queue[T]
	var zero T
	if q.Len() != 0 || !q.IsEmpty() || len(q.Values()) != 0 {
		t.Fatal("zero-value queue should be empty")
	}
	if got, ok := q.Dequeue(); got != zero || ok {
		t.Fatalf("empty Dequeue() = %v, %v; want %v, false", got, ok, zero)
	}
	if got, ok := q.Peek(); got != zero || ok {
		t.Fatalf("empty Peek() = %v, %v; want %v, false", got, ok, zero)
	}
	for _, n := range []int{-1, 0, 1} {
		if got := q.DequeueN(n); got != nil {
			t.Fatalf("empty DequeueN(%d) = %v, want nil", n, got)
		}
	}
	q.Clear()
	q.EnqueueAll()
	q.Enqueue(a)
	if got, ok := q.Peek(); got != a || !ok {
		t.Fatalf("Peek() = %v, %v; want %v, true", got, ok, a)
	}
	if q.Len() != 1 || q.IsEmpty() || !slices.Equal(q.Values(), []T{a}) {
		t.Fatal("Enqueue/Peek should leave one element")
	}
	if got, ok := q.Dequeue(); got != a || !ok {
		t.Fatalf("Dequeue() = %v, %v; want %v, true", got, ok, a)
	}
	if q.Len() != 0 || !q.IsEmpty() {
		t.Fatal("Dequeue should empty the queue")
	}

	var bulk Queue[T]
	bulk.EnqueueAll(a, b)
	if got := bulk.Values(); !slices.Equal(got, []T{a, b}) {
		t.Fatalf("Values() = %v, want %v", got, []T{a, b})
	}
	if got := bulk.DequeueN(1); !slices.Equal(got, []T{a}) {
		t.Fatalf("DequeueN(1) = %v, want %v", got, []T{a})
	}
	if got := bulk.Values(); !slices.Equal(got, []T{b}) {
		t.Fatalf("Values() after DequeueN = %v, want %v", got, []T{b})
	}
	bulk.Clear()
	if bulk.Len() != 0 || !bulk.IsEmpty() || len(bulk.Values()) != 0 {
		t.Fatal("Clear should empty the queue")
	}
}

func TestStack_ZeroValue(t *testing.T) {
	t.Run("int", func(t *testing.T) { checkZeroValueStack(t, 1, 2) })
	t.Run("pointer", func(t *testing.T) {
		a, b := 1, 2
		checkZeroValueStack(t, &a, &b)
	})
}

func checkZeroValueStack[T comparable](t *testing.T, a, b T) {
	t.Helper()
	var s Stack[T]
	var zero T
	if s.Len() != 0 || !s.IsEmpty() || len(s.Values()) != 0 {
		t.Fatal("zero-value stack should be empty")
	}
	if got, ok := s.Pop(); got != zero || ok {
		t.Fatalf("empty Pop() = %v, %v; want %v, false", got, ok, zero)
	}
	if got, ok := s.Peek(); got != zero || ok {
		t.Fatalf("empty Peek() = %v, %v; want %v, false", got, ok, zero)
	}
	for _, n := range []int{-1, 0, 1} {
		if got := s.PopN(n); got != nil {
			t.Fatalf("empty PopN(%d) = %v, want nil", n, got)
		}
	}
	if got := s.PopAll(); len(got) != 0 {
		t.Fatalf("empty PopAll() = %v, want no elements", got)
	}
	s.Clear()
	s.PushAll()
	s.Push(a)
	if got, ok := s.Peek(); got != a || !ok {
		t.Fatalf("Peek() = %v, %v; want %v, true", got, ok, a)
	}
	if s.Len() != 1 || s.IsEmpty() || !slices.Equal(s.Values(), []T{a}) {
		t.Fatal("Push/Peek should leave one element")
	}
	if got, ok := s.Pop(); got != a || !ok {
		t.Fatalf("Pop() = %v, %v; want %v, true", got, ok, a)
	}
	if s.Len() != 0 || !s.IsEmpty() {
		t.Fatal("Pop should empty the stack")
	}

	var bulk Stack[T]
	bulk.PushAll(a, b)
	if got := bulk.Values(); !slices.Equal(got, []T{a, b}) {
		t.Fatalf("Values() = %v, want %v", got, []T{a, b})
	}
	if got := bulk.PopN(1); !slices.Equal(got, []T{b}) {
		t.Fatalf("PopN(1) = %v, want %v", got, []T{b})
	}
	if got := bulk.PopAll(); !slices.Equal(got, []T{a}) {
		t.Fatalf("PopAll() = %v, want %v", got, []T{a})
	}
	if bulk.Len() != 0 || !bulk.IsEmpty() {
		t.Fatal("PopAll should empty the stack")
	}
	bulk.PushAll(a, b)
	bulk.Clear()
	if bulk.Len() != 0 || !bulk.IsEmpty() || len(bulk.Values()) != 0 {
		t.Fatal("Clear should empty the stack")
	}
}
