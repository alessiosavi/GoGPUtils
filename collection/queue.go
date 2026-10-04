package collection

// Queue is a generic FIFO (First-In-First-Out) data structure.
// The zero value is not usable; use NewQueue to create a Queue.
// A Queue must not be copied by value.
type Queue[T any] struct {
	buf     []T
	head, n int
}

// NewQueue creates a new empty Queue.
//
// Example:
//
//	queue := NewQueue[string]()
//	queue.Enqueue("task1")
func NewQueue[T any]() *Queue[T] { return &Queue[T]{} }

// NewQueueWithCapacity creates a Queue with pre-allocated capacity.
func NewQueueWithCapacity[T any](capacity int) *Queue[T] {
	b := make([]T, 0, capacity)
	return &Queue[T]{buf: b[:capacity]}
}

// overflowQueueGrowth produces append's actual runtime panic without allocating a huge backing.
func overflowQueueGrowth() { _ = append(make([]struct{}, int(^uint(0)>>1)), struct{}{}) }
func (q *Queue[T]) index(k int) int {
	distance := len(q.buf) - q.head
	if k >= distance {
		return k - distance
	}
	return q.head + k
}
func (q *Queue[T]) copyTo(dst []T) {
	first := min(len(dst), len(q.buf)-q.head)
	copy(dst, q.buf[q.head:q.head+first])
	copy(dst[first:], q.buf[:len(dst)-first])
}
func (q *Queue[T]) resize(c int) {
	b := make([]T, c)
	q.copyTo(b[:q.n])
	q.buf = b
	q.head = 0
}
func (q *Queue[T]) shrink() {
	c := len(q.buf)
	if q.n == 0 {
		if c > 64 {
			q.buf = nil
			q.head = 0
		}
		return
	}
	// c > 4*n, expressed without overflow.
	if c <= 64 || q.n > (c-1)/4 {
		return
	}
	target := 64
	for target/2 < q.n {
		if target > int(^uint(0)>>1)/2 {
			overflowQueueGrowth()
		}
		target *= 2
	}
	q.resize(target)
}

// Enqueue adds an element to the back of the queue.
//
// Example:
//
//	queue.Enqueue("task")
func (q *Queue[T]) Enqueue(item T) {
	if q.n == len(q.buf) {
		c := len(q.buf)
		if c == int(^uint(0)>>1) {
			overflowQueueGrowth()
		}
		if c > int(^uint(0)>>1)/2 {
			c = int(^uint(0) >> 1)
		} else {
			c *= 2
		}
		if c == 0 {
			c = 8
		}
		q.resize(c)
	}
	q.buf[q.index(q.n)] = item
	q.n++
}

// EnqueueAll adds multiple elements to the queue.
// Elements are added in order.
func (q *Queue[T]) EnqueueAll(items ...T) {
	if len(items) > int(^uint(0)>>1)-q.n {
		overflowQueueGrowth()
	}
	need := q.n + len(items)
	if need > len(q.buf) {
		c := len(q.buf)
		if c > int(^uint(0)>>1)/2 {
			c = int(^uint(0) >> 1)
		} else {
			c *= 2
		}
		q.resize(max(c, need))
	}
	if len(items) == 0 {
		return
	}
	i := q.index(q.n)
	first := min(len(items), len(q.buf)-i)
	copy(q.buf[i:i+first], items[:first])
	copy(q.buf[:len(items)-first], items[first:])
	q.n = need
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
	if q.n == 0 {
		var z T
		return z, false
	}
	v := q.buf[q.head]
	if mayHavePointers[T]() {
		var z T
		q.buf[q.head] = z
	}
	q.head++
	if q.head == len(q.buf) {
		q.head = 0
	}
	q.n--
	q.shrink()
	return v, true
}

// Peek returns the front element without removing it.
// Returns false if the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if q.n == 0 {
		var z T
		return z, false
	}
	return q.buf[q.head], true
}

// Len returns the number of elements in the queue.
func (q *Queue[T]) Len() int { return q.n }

// IsEmpty returns true if the queue has no elements.
func (q *Queue[T]) IsEmpty() bool { return q.n == 0 }
func (q *Queue[T]) clearN(n int) {
	if mayHavePointers[T]() {
		first := min(n, len(q.buf)-q.head)
		clear(q.buf[q.head : q.head+first])
		clear(q.buf[:n-first])
	}
}

// Clear removes all elements from the queue.
func (q *Queue[T]) Clear() {
	q.clearN(q.n)
	q.head = 0
	q.n = 0
}

// Values returns a copy of all elements in FIFO order.
func (q *Queue[T]) Values() []T {
	result := make([]T, q.n)
	q.copyTo(result)
	return result
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
	if n <= 0 || q.n == 0 {
		return nil
	}
	n = min(n, q.n)
	result := make([]T, n)
	q.copyTo(result)
	q.clearN(n)
	q.head = q.index(n)
	q.n -= n
	q.shrink()
	return result
}
