package collection

import (
	"cmp"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"slices"
	"testing"
)

// baseQueue is a generic FIFO (First-In-First-Out) data structure.
// The zero value is not usable; use baseNewQueue to create a baseQueue.
type baseQueue[T any] struct {
	items []T
}

// baseNewQueue creates a new empty baseQueue.
//
// Example:
//
//	queue := baseNewQueue[string]()
//	queue.Enqueue("task1")
func baseNewQueue[T any]() *baseQueue[T] {
	return &baseQueue[T]{
		items: make([]T, 0),
	}
}

// baseNewQueueWithCapacity creates a baseQueue with pre-allocated capacity.
func baseNewQueueWithCapacity[T any](capacity int) *baseQueue[T] {
	return &baseQueue[T]{
		items: make([]T, 0, capacity),
	}
}

// Enqueue adds an element to the back of the queue.
//
// Example:
//
//	queue.Enqueue("task")
func (q *baseQueue[T]) Enqueue(item T) {
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
func (q *baseQueue[T]) Dequeue() (T, bool) {
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
func (q *baseQueue[T]) Peek() (T, bool) {
	if len(q.items) == 0 {
		var zero T

		return zero, false
	}

	return q.items[0], true
}

// Len returns the number of elements in the queue.
func (q *baseQueue[T]) Len() int {
	return len(q.items)
}

// IsEmpty returns true if the queue has no elements.
func (q *baseQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Clear removes all elements from the queue.
func (q *baseQueue[T]) Clear() {
	q.items = q.items[:0]
}

// Values returns a copy of all elements in FIFO order.
func (q *baseQueue[T]) Values() []T {
	result := make([]T, len(q.items))
	copy(result, q.items)

	return result
}

// EnqueueAll adds multiple elements to the queue.
// Elements are added in order.
func (q *baseQueue[T]) EnqueueAll(items ...T) {
	q.items = append(q.items, items...)
}

// DequeueN removes and returns up to n elements from the front of the queue in FIFO order.
// If n > Len(), all elements are returned. If n <= 0 or the queue is empty, returns nil.
//
// Example:
//
//	q := baseNewQueue[int]()
//	q.EnqueueAll(1, 2, 3, 4, 5)
//	q.DequeueN(3) // [1, 2, 3]; queue now contains [4, 5]
func (q *baseQueue[T]) DequeueN(n int) []T {
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

// baseStack is a generic LIFO (Last-In-First-Out) data structure.
// The zero value is not usable; use baseNewStack to create a baseStack.
type baseStack[T any] struct {
	items []T
}

// baseNewStack creates a new empty baseStack.
//
// Example:
//
//	stack := baseNewStack[int]()
//	stack.Push(1)
func baseNewStack[T any]() *baseStack[T] {
	return &baseStack[T]{
		items: make([]T, 0),
	}
}

// baseNewStackWithCapacity creates a baseStack with pre-allocated capacity.
// Use this when you know the approximate size to reduce allocations.
func baseNewStackWithCapacity[T any](capacity int) *baseStack[T] {
	return &baseStack[T]{
		items: make([]T, 0, capacity),
	}
}

// Push adds an element to the top of the stack.
//
// Example:
//
//	stack.Push(42)
func (s *baseStack[T]) Push(item T) {
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
func (s *baseStack[T]) Pop() (T, bool) {
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
func (s *baseStack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T

		return zero, false
	}

	return s.items[len(s.items)-1], true
}

// Len returns the number of elements in the stack.
func (s *baseStack[T]) Len() int {
	return len(s.items)
}

// IsEmpty returns true if the stack has no elements.
func (s *baseStack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Clear removes all elements from the stack.
func (s *baseStack[T]) Clear() {
	s.items = s.items[:0]
}

// Values returns a copy of all elements from bottom to top.
func (s *baseStack[T]) Values() []T {
	result := make([]T, len(s.items))
	copy(result, s.items)

	return result
}

// PushAll adds multiple elements to the stack.
// Elements are pushed in order, so the last element becomes the top.
func (s *baseStack[T]) PushAll(items ...T) {
	s.items = append(s.items, items...)
}

// PopAll removes and returns all elements from top to bottom.
// The stack will be empty after this operation.
//
//nolint:modernize // Preserve the verbatim BASE oracle.
func (s *baseStack[T]) PopAll() []T {
	result := make([]T, len(s.items))
	for i := len(s.items) - 1; i >= 0; i-- {
		result[len(s.items)-1-i] = s.items[i]
	}

	s.items = s.items[:0]

	return result
}

// PopN removes and returns up to n elements from the top of the stack in LIFO order.
// If n > Len(), all elements are returned. If n <= 0 or the stack is empty, returns nil.
//
// Example:
//
//	s := baseNewStack[int]()
//	s.PushAll(1, 2, 3, 4, 5) // top is 5
//	s.PopN(3)                 // [5, 4, 3]; stack now contains [1, 2]
func (s *baseStack[T]) PopN(n int) []T {
	if n <= 0 || len(s.items) == 0 {
		return nil
	}

	if n > len(s.items) {
		n = len(s.items)
	}

	result := make([]T, n)
	for i := range n {
		result[i] = s.items[len(s.items)-1-i]
	}

	s.items = s.items[:len(s.items)-n]

	return result
}

// baseBST is a generic Binary Search Tree.
// Elements must be ordered (satisfy cmp.Ordered constraint).
type baseBST[T cmp.Ordered] struct {
	root *baseBSTNode[T]
	size int
}

type baseBSTNode[T cmp.Ordered] struct {
	value T
	left  *baseBSTNode[T]
	right *baseBSTNode[T]
}

// baseNewBST creates a new empty Binary Search Tree.
//
// Example:
//
//	tree := baseNewBST[int]()
//	tree.Insert(5, 3, 7)
func baseNewBST[T cmp.Ordered]() *baseBST[T] {
	return &baseBST[T]{}
}

// baseNewBSTFrom creates a baseBST from a slice.
func baseNewBSTFrom[T cmp.Ordered](items []T) *baseBST[T] {
	tree := baseNewBST[T]()
	tree.Insert(items...)

	return tree
}

// Insert adds one or more values to the tree.
// Duplicate values are ignored.
//
// Example:
//
//	tree.Insert(5, 3, 7, 1, 4)
func (t *baseBST[T]) Insert(values ...T) {
	for _, v := range values {
		t.insert(v)
	}
}

//nolint:gocritic // Preserve the verbatim BASE oracle.
func (t *baseBST[T]) insert(value T) {
	if t.root == nil {
		t.root = &baseBSTNode[T]{value: value}
		t.size++

		return
	}

	current := t.root

	for {
		if value < current.value {
			if current.left == nil {
				current.left = &baseBSTNode[T]{value: value}
				t.size++

				return
			}

			current = current.left
		} else if value > current.value {
			if current.right == nil {
				current.right = &baseBSTNode[T]{value: value}
				t.size++

				return
			}

			current = current.right
		} else {
			// Duplicate, ignore
			return
		}
	}
}

// Contains returns true if the value is in the tree.
func (t *baseBST[T]) Contains(value T) bool {
	return t.find(value) != nil
}

//nolint:gocritic // Preserve the verbatim BASE oracle.
func (t *baseBST[T]) find(value T) *baseBSTNode[T] {
	current := t.root
	for current != nil {
		if value < current.value {
			current = current.left
		} else if value > current.value {
			current = current.right
		} else {
			return current
		}
	}

	return nil
}

// Remove deletes a value from the tree.
// Returns true if the value was found and removed.
//
//nolint:gocritic // Preserve the verbatim BASE oracle.
func (t *baseBST[T]) Remove(value T) bool {
	var parent *baseBSTNode[T]

	current := t.root
	isLeft := false

	// Find the node to remove
	for current != nil && current.value != value {
		parent = current

		if value < current.value {
			current = current.left
			isLeft = true
		} else {
			current = current.right
			isLeft = false
		}
	}

	if current == nil {
		return false
	}

	// Case 1: No children
	if current.left == nil && current.right == nil {
		if parent == nil {
			t.root = nil
		} else if isLeft {
			parent.left = nil
		} else {
			parent.right = nil
		}
	} else if current.left == nil {
		// Case 2: Only right child
		if parent == nil {
			t.root = current.right
		} else if isLeft {
			parent.left = current.right
		} else {
			parent.right = current.right
		}
	} else if current.right == nil {
		// Case 3: Only left child
		if parent == nil {
			t.root = current.left
		} else if isLeft {
			parent.left = current.left
		} else {
			parent.right = current.left
		}
	} else {
		// Case 4: Two children - find in-order successor
		successor := current.right
		successorParent := current

		for successor.left != nil {
			successorParent = successor
			successor = successor.left
		}

		current.value = successor.value

		if successorParent == current {
			successorParent.right = successor.right
		} else {
			successorParent.left = successor.right
		}
	}

	t.size--

	return true
}

// Len returns the number of elements in the tree.
func (t *baseBST[T]) Len() int {
	return t.size
}

// IsEmpty returns true if the tree has no elements.
func (t *baseBST[T]) IsEmpty() bool {
	return t.size == 0
}

// Clear removes all elements from the tree.
func (t *baseBST[T]) Clear() {
	t.root = nil
	t.size = 0
}

// Min returns the minimum value in the tree.
// Returns zero value and false if tree is empty.
func (t *baseBST[T]) Min() (T, bool) {
	if t.root == nil {
		var zero T

		return zero, false
	}

	current := t.root
	for current.left != nil {
		current = current.left
	}

	return current.value, true
}

// Max returns the maximum value in the tree.
// Returns zero value and false if tree is empty.
func (t *baseBST[T]) Max() (T, bool) {
	if t.root == nil {
		var zero T

		return zero, false
	}

	current := t.root
	for current.right != nil {
		current = current.right
	}

	return current.value, true
}

// InOrder returns values in sorted order (left, root, right).
//
// Example:
//
//	tree.Insert(5, 3, 7, 1, 4)
//	tree.InOrder() // [1, 3, 4, 5, 7]
func (t *baseBST[T]) InOrder() []T {
	var result []T

	t.inOrderTraverse(t.root, &result)

	return result
}

func (t *baseBST[T]) inOrderTraverse(node *baseBSTNode[T], result *[]T) {
	if node == nil {
		return
	}

	t.inOrderTraverse(node.left, result)
	*result = append(*result, node.value)
	t.inOrderTraverse(node.right, result)
}

// PreOrder returns values in pre-order (root, left, right).
func (t *baseBST[T]) PreOrder() []T {
	var result []T

	t.preOrderTraverse(t.root, &result)

	return result
}

func (t *baseBST[T]) preOrderTraverse(node *baseBSTNode[T], result *[]T) {
	if node == nil {
		return
	}

	*result = append(*result, node.value)
	t.preOrderTraverse(node.left, result)
	t.preOrderTraverse(node.right, result)
}

// PostOrder returns values in post-order (left, right, root).
func (t *baseBST[T]) PostOrder() []T {
	var result []T

	t.postOrderTraverse(t.root, &result)

	return result
}

func (t *baseBST[T]) postOrderTraverse(node *baseBSTNode[T], result *[]T) {
	if node == nil {
		return
	}

	t.postOrderTraverse(node.left, result)
	t.postOrderTraverse(node.right, result)
	*result = append(*result, node.value)
}

// LevelOrder returns values level by level (breadth-first).
func (t *baseBST[T]) LevelOrder() []T {
	if t.root == nil {
		return nil
	}

	var result []T

	queue := []*baseBSTNode[T]{t.root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node.value)

		if node.left != nil {
			queue = append(queue, node.left)
		}

		if node.right != nil {
			queue = append(queue, node.right)
		}
	}

	return result
}

// Height returns the height of the tree.
// An empty tree has height 0; a single node has height 1.
func (t *baseBST[T]) Height() int {
	return t.height(t.root)
}

func (t *baseBST[T]) height(node *baseBSTNode[T]) int {
	if node == nil {
		return 0
	}

	leftHeight := t.height(node.left)

	rightHeight := t.height(node.right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}

	return rightHeight + 1
}

// ForEach calls a function for each value in sorted order.
func (t *baseBST[T]) ForEach(fn func(T)) {
	t.forEach(t.root, fn)
}

func (t *baseBST[T]) forEach(node *baseBSTNode[T], fn func(T)) {
	if node == nil {
		return
	}

	t.forEach(node.left, fn)
	fn(node.value)
	t.forEach(node.right, fn)
}

// Values returns all values in sorted order.
// Alias for InOrder().
func (t *baseBST[T]) Values() []T {
	return t.InOrder()
}

// RangeSearch returns all values in the inclusive range [min, max] in sorted order.
// Uses the baseBST structure to prune branches outside the range efficiently
// (O(k + log n) for k results rather than O(n)).
// Returns nil if the tree is empty, min > max, or no values fall in the range.
//
// Example:
//
//	tree := baseNewBSTFrom([]int{1, 3, 5, 7, 9})
//	tree.RangeSearch(3, 7) // [3, 5, 7]
//
//nolint:revive // Preserve the verbatim BASE oracle.
func (t *baseBST[T]) RangeSearch(min, max T) []T {
	if t.root == nil || min > max {
		return nil
	}

	var result []T

	t.rangeSearch(t.root, min, max, &result)

	if len(result) == 0 {
		return nil
	}

	return result
}

//nolint:revive // Preserve the verbatim BASE oracle.
func (t *baseBST[T]) rangeSearch(node *baseBSTNode[T], min, max T, result *[]T) {
	if node == nil {
		return
	}

	// Only recurse left if current node value is greater than min
	if node.value > min {
		t.rangeSearch(node.left, min, max, result)
	}

	if node.value >= min && node.value <= max {
		*result = append(*result, node.value)
	}

	// Only recurse right if current node value is less than max
	if node.value < max {
		t.rangeSearch(node.right, min, max, result)
	}
}

func same[T any](t *testing.T, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
func modelContainers[T comparable](t *testing.T, val func(int) T) {
	for _, seed := range []int64{1, 7, 42, 99} {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			r := rand.New(rand.NewSource(seed))
			var q Queue[T]
			var bq baseQueue[T]
			var s Stack[T]
			var bs baseStack[T]
			check := func() {
				t.Helper()
				sameItems(t, q.Values(), bq.Values())
				sameItems(t, s.Values(), bs.Values())
				same(t, q.Len(), bq.Len())
				same(t, s.Len(), bs.Len())
				same(t, q.IsEmpty(), bq.IsEmpty())
				same(t, s.IsEmpty(), bs.IsEmpty())
			}
			add := func(v []T) { q.EnqueueAll(v...); bq.EnqueueAll(v...); s.PushAll(v...); bs.PushAll(v...) }
			pop := func() {
				g, ok := q.Dequeue()
				w, wo := bq.Dequeue()
				sameElement(t, g, w)
				same(t, ok, wo)
				g, ok = s.Pop()
				w, wo = bs.Pop()
				sameElement(t, g, w)
				same(t, ok, wo)
			}
			bulk := func(n int) { sameItems(t, q.DequeueN(n), bq.DequeueN(n)); sameItems(t, s.PopN(n), bs.PopN(n)) }
			// Directed sequences ensure wrapping, growth, retained empty capacity and Clear/refill.
			for _, n := range []int{0, 1, 7, 8, 9, 63, 64, 65, 100, 1024} {
				v := make([]T, n)
				for i := range v {
					v[i] = val(i)
				}
				add(v)
				check()
				for range n / 2 {
					pop()
				}
				add(v)
				check()
				bulk(n + 3)
				check()
				q.Clear()
				bq.Clear()
				s.Clear()
				bs.Clear()
				add(v)
				check()
				bulk(100000)
				check()
			}
			for i := range 10000 {
				switch r.Intn(10) {
				case 0, 1:
					v := val(r.Intn(1000))
					q.Enqueue(v)
					bq.Enqueue(v)
					s.Push(v)
					bs.Push(v)
				case 2:
					pop()
				case 3:
					v := make([]T, r.Intn(40))
					for j := range v {
						v[j] = val(r.Intn(1000))
					}
					add(v)
				case 4:
					ns := []int{-1, 0, 1, r.Intn(40), q.Len() + s.Len() + 1}
					bulk(ns[r.Intn(len(ns))])
				case 5:
					q.Clear()
					bq.Clear()
					s.Clear()
					bs.Clear()
				case 6:
					g, ok := q.Peek()
					w, wo := bq.Peek()
					sameElement(t, g, w)
					same(t, ok, wo)
					g, ok = s.Peek()
					w, wo = bs.Peek()
					sameElement(t, g, w)
					same(t, ok, wo)
				case 7:
					sameItems(t, s.PopAll(), bs.PopAll())
				case 8: // Returned slices must be independent of the container.
					v := q.Values()
					if len(v) > 0 {
						v[0] = val(12345)
					}
					v = s.Values()
					if len(v) > 0 {
						v[0] = val(12345)
					}
				case 9: // Inputs must be copied too.
					v := []T{val(i)}
					add(v)
					v[0] = val(i + 1)
				}
				check()
			}
		})
	}
}
func TestOracleQueueStack(t *testing.T) {
	t.Run("int", func(t *testing.T) { modelContainers(t, func(i int) int { return i }) })
	pool := make([]int, 20000)
	t.Run("ptr", func(t *testing.T) {
		modelContainers(t, func(i int) *int {
			if i%11 == 0 {
				return nil
			}
			return &pool[i%len(pool)]
		})
	})
}
func panicOutcome(f func()) (msg string) {
	defer func() {
		if p := recover(); p != nil {
			msg = fmt.Sprintf("%T:%v", p, p)
		}
	}()
	f()
	return
}
func TestOracleConstructorPanic(t *testing.T) {
	for _, n := range []int{-1, -100, -int(^uint(0)>>1) - 1} {
		same(t, panicOutcome(func() { NewQueueWithCapacity[int](n) }), panicOutcome(func() { baseNewQueueWithCapacity[int](n) }))
	}
}
func countNodes[T cmp.Ordered](n *bstNode[T]) int {
	if n == nil {
		return 0
	}
	return 1 + countNodes(n.left) + countNodes(n.right)
}
func TestOracleBST(t *testing.T) {
	for _, seed := range []int64{1, 7, 42, 99} {
		r := rand.New(rand.NewSource(seed))
		var b BST[int]
		var o baseBST[int]
		for range 10000 {
			v := r.Intn(200) - 100
			switch r.Intn(25) {
			case 0:
				b.Clear()
				o.Clear()
			case 1, 2, 3, 4, 5, 6, 7, 8:
				same(t, b.Remove(v), o.Remove(v))
			default:
				b.Insert(v)
				o.Insert(v)
			}
			same(t, b.InOrder(), o.InOrder())
			same(t, b.PreOrder(), o.PreOrder())
			same(t, b.PostOrder(), o.PostOrder())
			same(t, b.LevelOrder(), o.LevelOrder())
			same(t, b.Values(), o.Values())
			same(t, b.size, countNodes(b.root))
			same(t, b.size, len(o.InOrder()))
			same(t, o.size, len(o.InOrder()))
		}
	}
}
func TestOracleBSTFloats(t *testing.T) {
	vals := []float64{math.NaN(), math.Float64frombits(0x7ff8000000001234), 0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), 1, -1}
	for _, first := range vals {
		var b BST[float64]
		var o baseBST[float64]
		b.Insert(first)
		o.Insert(first)
		check := func() {
			for _, pair := range [][2][]float64{{b.InOrder(), o.InOrder()}, {b.PreOrder(), o.PreOrder()}, {b.PostOrder(), o.PostOrder()}, {b.LevelOrder(), o.LevelOrder()}, {b.Values(), o.Values()}} {
				same(t, pair[0] == nil, pair[1] == nil)
				same(t, len(pair[0]), len(pair[1]))
				for i := range pair[0] {
					same(t, math.Float64bits(pair[0][i]), math.Float64bits(pair[1][i]))
				}
			}
			same(t, b.size, countNodes(b.root))
		}
		for _, v := range vals {
			b.Insert(v)
			o.Insert(v)
			check()
		}
		for _, v := range vals {
			same(t, b.Remove(v), o.Remove(v))
			check()
		}
		b.Clear()
		o.Clear()
		check()
	}
}

func TestOracleBSTValueCopy(t *testing.T) {
	a := NewBSTFrom([]int{0})
	o := baseNewBSTFrom([]int{0})
	b, bo := *a, *o
	a.Insert(1, 2)
	o.Insert(1, 2)
	check := func() {
		same(t, b.InOrder(), bo.InOrder())
		same(t, b.PreOrder(), bo.PreOrder())
		same(t, b.PostOrder(), bo.PostOrder())
		same(t, b.LevelOrder(), bo.LevelOrder())
		same(t, b.Values(), bo.Values())
	}
	same(t, b.Remove(1), bo.Remove(1))
	check()
	same(t, b.Remove(2), bo.Remove(2))
	check()
}

func TestQueueZeroSizedWrap(t *testing.T) {
	c := int(^uint(0) >> 1)
	q := NewQueueWithCapacity[struct{}](c)
	b := baseNewQueueWithCapacity[struct{}](c)
	xs := make([]struct{}, c)
	q.EnqueueAll(xs...)
	b.EnqueueAll(xs...)
	q.DequeueN(2)
	b.DequeueN(2)
	q.EnqueueAll(struct{}{}, struct{}{})
	b.EnqueueAll(struct{}{}, struct{}{})
	want := len(b.DequeueN(c - 1))
	same(t, len(q.DequeueN(c-1)), want)
	same(t, q.Len(), b.Len())
	same(t, len(q.Values()), len(b.Values()))
}
func TestQueueOverflowPanic(t *testing.T) {
	c := int(^uint(0) >> 1)
	q := NewQueueWithCapacity[struct{}](c)
	q.EnqueueAll(make([]struct{}, c)...)
	b := baseNewQueueWithCapacity[struct{}](c)
	b.EnqueueAll(make([]struct{}, c)...)
	same(t, panicOutcome(func() { q.Enqueue(struct{}{}) }), panicOutcome(func() { b.Enqueue(struct{}{}) }))
}

func TestQueueZeroSizedGrowth(t *testing.T) {
	c := int(^uint(0)>>1)/2 + 1
	for _, bulk := range []bool{false, true} {
		q := NewQueueWithCapacity[struct{}](c)
		b := baseNewQueueWithCapacity[struct{}](c)
		xs := make([]struct{}, c)
		q.EnqueueAll(xs...)
		b.EnqueueAll(xs...)
		op := func() {
			if bulk {
				q.EnqueueAll(struct{}{})
			} else {
				q.Enqueue(struct{}{})
			}
		}
		baseOp := func() {
			if bulk {
				b.EnqueueAll(struct{}{})
			} else {
				b.Enqueue(struct{}{})
			}
		}
		same(t, panicOutcome(op), panicOutcome(baseOp))
		same(t, q.Len(), b.Len())
	}
}

func TestRemovedResultOwnership(t *testing.T) {
	q := NewQueue[int]()
	q.EnqueueAll(1, 2, 3)
	v := q.DequeueN(2)
	q.EnqueueAll(4, 5, 6)
	same(t, v, []int{1, 2})
	v[0] = 90
	same(t, q.Values(), []int{3, 4, 5, 6})
	s := NewStack[int]()
	s.PushAll(1, 2, 3)
	v = s.PopN(2)
	s.PushAll(4, 5, 6)
	same(t, v, []int{3, 2})
	v[0] = 90
	same(t, s.Values(), []int{1, 4, 5, 6})
}

func TestOracleConstructors(t *testing.T) {
	same(t, NewQueue[int]().Values(), baseNewQueue[int]().Values())
	same(t, NewStack[int]().Values(), baseNewStack[int]().Values())
	for _, n := range []int{0, 1, 64} {
		same(t, NewQueueWithCapacity[int](n).Values(), baseNewQueueWithCapacity[int](n).Values())
		same(t, NewStackWithCapacity[int](n).Values(), baseNewStackWithCapacity[int](n).Values())
	}
}

func oracleWrappedRetirement[T comparable](t *testing.T, value T) {
	t.Helper()
	for _, clearAll := range []bool{false, true} {
		q := NewQueueWithCapacity[T](8)
		b := baseNewQueueWithCapacity[T](8)
		for range 6 {
			q.Enqueue(value)
			b.Enqueue(value)
		}
		sameItems(t, q.DequeueN(5), b.DequeueN(5))
		for range 6 {
			q.Enqueue(value)
			b.Enqueue(value)
		}
		// The ring's live positions are 5,6,7,0,1,2,3: retire both segments.
		if clearAll {
			q.Clear()
			b.Clear()
		} else {
			sameItems(t, q.DequeueN(6), b.DequeueN(6))
		}
		sameItems(t, q.Values(), b.Values())
		q.Enqueue(value)
		b.Enqueue(value)
		sameItems(t, q.Values(), b.Values())
	}
}

func TestOracleWrappedRetirement(t *testing.T) {
	oracleWrappedRetirement(t, 7)
	value := 7
	oracleWrappedRetirement(t, &value)
}

// Comparable equality preserves pointer identity; DeepEqual would only compare
// the pointed-to values, which can be equal for distinct model payloads.
func sameElement[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func sameItems[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	same(t, got == nil, want == nil)
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
