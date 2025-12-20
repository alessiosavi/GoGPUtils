package collection

// Stack is a generic LIFO (Last-In-First-Out) data structure.
// The zero value is not usable; use NewStack to create a Stack.
type Stack[T any] struct {
	items []T
}

// NewStack creates a new empty Stack.
//
// Example:
//
//	stack := NewStack[int]()
//	stack.Push(1)
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

// NewStackWithCapacity creates a Stack with pre-allocated capacity.
// Use this when you know the approximate size to reduce allocations.
func NewStackWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, capacity),
	}
}

// Push adds an element to the top of the stack.
//
// Example:
//
//	stack.Push(42)
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top element.
// Returns false if the stack is empty.
//
// Example:
//
//	val, ok := stack.Pop()
//	if !ok {
//	    // stack was empty
//	}
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	idx := len(s.items) - 1
	item := s.items[idx]
	s.items = s.items[:idx]
	return item, true
}

// Peek returns the top element without removing it.
// Returns false if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of elements in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the stack has no elements.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	s.items = s.items[:0]
}

// Values returns a copy of all elements from bottom to top.
func (s *Stack[T]) Values() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)
	return result
}

// PushAll adds multiple elements to the stack.
// Elements are pushed in order, so the last element becomes the top.
func (s *Stack[T]) PushAll(items ...T) {
	s.items = append(s.items, items...)
}

// PopAll removes and returns all elements from top to bottom.
// The stack will be empty after this operation.
func (s *Stack[T]) PopAll() []T {
	result := make([]T, len(s.items))
	for i := len(s.items) - 1; i >= 0; i-- {
		result[len(s.items)-1-i] = s.items[i]
	}
	s.items = s.items[:0]
	return result
}
