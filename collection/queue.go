package collection

// Queue is a generic FIFO (First-In-First-Out) data structure.
// The zero value is not usable; use NewQueue to create a Queue.
type Queue[T any] struct {
	items []T
}

// NewQueue creates a new empty Queue.
//
// Example:
//
//	queue := NewQueue[string]()
//	queue.Enqueue("task1")
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
	}
}

// NewQueueWithCapacity creates a Queue with pre-allocated capacity.
func NewQueueWithCapacity[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0, capacity),
	}
}

// Enqueue adds an element to the back of the queue.
//
// Example:
//
//	queue.Enqueue("task")
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front element.
// Returns false if the queue is empty.
//
// Example:
//
//	val, ok := queue.Dequeue()
//	if !ok {
//	    // queue was empty
//	}
func (q *Queue[T]) Dequeue() (T, bool) {
	if len(q.items) == 0 {
		var zero T

		return zero, false
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// Peek returns the front element without removing it.
// Returns false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T

		return zero, false
	}

	return q.items[0], true
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// IsEmpty returns true if the queue has no elements.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Clear removes all elements from the queue.
func (q *Queue[T]) Clear() {
	q.items = q.items[:0]
}

// Values returns a copy of all elements in FIFO order.
func (q *Queue[T]) Values() []T {
	result := make([]T, len(q.items))
	copy(result, q.items)

	return result
}

// EnqueueAll adds multiple elements to the queue.
// Elements are added in order.
func (q *Queue[T]) EnqueueAll(items ...T) {
	q.items = append(q.items, items...)
}

// DequeueN removes and returns up to n elements from the front of the queue in FIFO order.
// If n > Len(), all elements are returned. If n <= 0 or the queue is empty, returns nil.
//
// Example:
//
//	q := NewQueue[int]()
//	q.EnqueueAll(1, 2, 3, 4, 5)
//	q.DequeueN(3) // [1, 2, 3]; queue now contains [4, 5]
func (q *Queue[T]) DequeueN(n int) []T {
	if n <= 0 || len(q.items) == 0 {
		return nil
	}

	if n > len(q.items) {
		n = len(q.items)
	}

	result := make([]T, n)
	copy(result, q.items[:n])
	q.items = q.items[n:]

	return result
}
