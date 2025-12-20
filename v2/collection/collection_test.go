package collection

import (
	"slices"
	"testing"
)

// ============================================================================
// Stack Tests
// ============================================================================

func TestStack_PushPop(t *testing.T) {
	stack := NewStack[int]()

	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Len() != 3 {
		t.Errorf("Len() = %d, want 3", stack.Len())
	}

	val, ok := stack.Pop()
	if !ok || val != 3 {
		t.Errorf("Pop() = %d, %v; want 3, true", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 2 {
		t.Errorf("Pop() = %d, %v; want 2, true", val, ok)
	}

	val, ok = stack.Pop()
	if !ok || val != 1 {
		t.Errorf("Pop() = %d, %v; want 1, true", val, ok)
	}

	_, ok = stack.Pop()
	if ok {
		t.Error("Pop() on empty stack should return false")
	}
}

func TestStack_Peek(t *testing.T) {
	stack := NewStack[string]()

	_, ok := stack.Peek()
	if ok {
		t.Error("Peek() on empty stack should return false")
	}

	stack.Push("hello")
	val, ok := stack.Peek()
	if !ok || val != "hello" {
		t.Errorf("Peek() = %q, %v; want 'hello', true", val, ok)
	}

	// Stack should be unchanged
	if stack.Len() != 1 {
		t.Error("Peek() should not modify stack")
	}
}

func TestStack_IsEmpty(t *testing.T) {
	stack := NewStack[int]()

	if !stack.IsEmpty() {
		t.Error("IsEmpty() should return true for new stack")
	}

	stack.Push(1)
	if stack.IsEmpty() {
		t.Error("IsEmpty() should return false after push")
	}

	stack.Pop()
	if !stack.IsEmpty() {
		t.Error("IsEmpty() should return true after popping all elements")
	}
}

func TestStack_Clear(t *testing.T) {
	stack := NewStack[int]()
	stack.PushAll(1, 2, 3, 4, 5)

	stack.Clear()

	if !stack.IsEmpty() {
		t.Error("Clear() should empty the stack")
	}
}

func TestStack_Values(t *testing.T) {
	stack := NewStack[int]()
	stack.PushAll(1, 2, 3)

	values := stack.Values()
	want := []int{1, 2, 3}

	if !slices.Equal(values, want) {
		t.Errorf("Values() = %v, want %v", values, want)
	}

	// Values should be a copy
	values[0] = 999
	if stack.Values()[0] == 999 {
		t.Error("Values() should return a copy")
	}
}

func TestStack_PopAll(t *testing.T) {
	stack := NewStack[int]()
	stack.PushAll(1, 2, 3)

	values := stack.PopAll()
	want := []int{3, 2, 1} // Reversed order

	if !slices.Equal(values, want) {
		t.Errorf("PopAll() = %v, want %v", values, want)
	}

	if !stack.IsEmpty() {
		t.Error("Stack should be empty after PopAll()")
	}
}

// ============================================================================
// Queue Tests
// ============================================================================

func TestQueue_EnqueueDequeue(t *testing.T) {
	queue := NewQueue[int]()

	queue.Enqueue(1)
	queue.Enqueue(2)
	queue.Enqueue(3)

	if queue.Len() != 3 {
		t.Errorf("Len() = %d, want 3", queue.Len())
	}

	val, ok := queue.Dequeue()
	if !ok || val != 1 {
		t.Errorf("Dequeue() = %d, %v; want 1, true", val, ok)
	}

	val, ok = queue.Dequeue()
	if !ok || val != 2 {
		t.Errorf("Dequeue() = %d, %v; want 2, true", val, ok)
	}

	val, ok = queue.Dequeue()
	if !ok || val != 3 {
		t.Errorf("Dequeue() = %d, %v; want 3, true", val, ok)
	}

	_, ok = queue.Dequeue()
	if ok {
		t.Error("Dequeue() on empty queue should return false")
	}
}

func TestQueue_Peek(t *testing.T) {
	queue := NewQueue[string]()

	_, ok := queue.Peek()
	if ok {
		t.Error("Peek() on empty queue should return false")
	}

	queue.Enqueue("first")
	queue.Enqueue("second")

	val, ok := queue.Peek()
	if !ok || val != "first" {
		t.Errorf("Peek() = %q, %v; want 'first', true", val, ok)
	}

	// Queue should be unchanged
	if queue.Len() != 2 {
		t.Error("Peek() should not modify queue")
	}
}

func TestQueue_Values(t *testing.T) {
	queue := NewQueue[int]()
	queue.EnqueueAll(1, 2, 3)

	values := queue.Values()
	want := []int{1, 2, 3}

	if !slices.Equal(values, want) {
		t.Errorf("Values() = %v, want %v", values, want)
	}
}

// ============================================================================
// Set Tests
// ============================================================================

func TestSet_AddContains(t *testing.T) {
	set := NewSet[int]()

	set.Add(1, 2, 3)

	if !set.Contains(1) || !set.Contains(2) || !set.Contains(3) {
		t.Error("Set should contain added elements")
	}

	if set.Contains(4) {
		t.Error("Set should not contain unadded elements")
	}

	if set.Len() != 3 {
		t.Errorf("Len() = %d, want 3", set.Len())
	}

	// Adding duplicate should not increase size
	set.Add(1)
	if set.Len() != 3 {
		t.Error("Adding duplicate should not increase size")
	}
}

func TestSet_Remove(t *testing.T) {
	set := NewSetFrom([]int{1, 2, 3, 4, 5})

	set.Remove(3)

	if set.Contains(3) {
		t.Error("Set should not contain removed element")
	}
	if set.Len() != 4 {
		t.Errorf("Len() = %d, want 4", set.Len())
	}

	// Remove non-existent element
	set.Remove(999)
	if set.Len() != 4 {
		t.Error("Removing non-existent element should not change size")
	}
}

func TestSet_ContainsAll(t *testing.T) {
	set := NewSetFrom([]int{1, 2, 3, 4, 5})

	if !set.ContainsAll(1, 3, 5) {
		t.Error("ContainsAll should return true for subset")
	}

	if set.ContainsAll(1, 6) {
		t.Error("ContainsAll should return false when any element missing")
	}
}

func TestSet_ContainsAny(t *testing.T) {
	set := NewSetFrom([]int{1, 2, 3})

	if !set.ContainsAny(3, 4, 5) {
		t.Error("ContainsAny should return true when any element present")
	}

	if set.ContainsAny(4, 5, 6) {
		t.Error("ContainsAny should return false when no elements present")
	}
}

func TestSet_Union(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3})
	b := NewSetFrom([]int{3, 4, 5})

	union := a.Union(b)

	if union.Len() != 5 {
		t.Errorf("Union Len() = %d, want 5", union.Len())
	}

	for _, v := range []int{1, 2, 3, 4, 5} {
		if !union.Contains(v) {
			t.Errorf("Union should contain %d", v)
		}
	}
}

func TestSet_Intersection(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3, 4})
	b := NewSetFrom([]int{3, 4, 5, 6})

	intersection := a.Intersection(b)

	if intersection.Len() != 2 {
		t.Errorf("Intersection Len() = %d, want 2", intersection.Len())
	}

	if !intersection.ContainsAll(3, 4) {
		t.Error("Intersection should contain 3 and 4")
	}
}

func TestSet_Difference(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3, 4})
	b := NewSetFrom([]int{3, 4, 5})

	diff := a.Difference(b)

	if diff.Len() != 2 {
		t.Errorf("Difference Len() = %d, want 2", diff.Len())
	}

	if !diff.ContainsAll(1, 2) {
		t.Error("Difference should contain 1 and 2")
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3})
	b := NewSetFrom([]int{3, 4, 5})

	symDiff := a.SymmetricDifference(b)

	if symDiff.Len() != 4 {
		t.Errorf("SymmetricDifference Len() = %d, want 4", symDiff.Len())
	}

	if !symDiff.ContainsAll(1, 2, 4, 5) {
		t.Error("SymmetricDifference should contain 1, 2, 4, 5")
	}

	if symDiff.Contains(3) {
		t.Error("SymmetricDifference should not contain 3")
	}
}

func TestSet_IsSubset(t *testing.T) {
	a := NewSetFrom([]int{1, 2})
	b := NewSetFrom([]int{1, 2, 3, 4})

	if !a.IsSubset(b) {
		t.Error("a should be subset of b")
	}

	if b.IsSubset(a) {
		t.Error("b should not be subset of a")
	}
}

func TestSet_IsSuperset(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3, 4})
	b := NewSetFrom([]int{1, 2})

	if !a.IsSuperset(b) {
		t.Error("a should be superset of b")
	}
}

func TestSet_Equal(t *testing.T) {
	a := NewSetFrom([]int{1, 2, 3})
	b := NewSetFrom([]int{3, 1, 2})
	c := NewSetFrom([]int{1, 2, 4})

	if !a.Equal(b) {
		t.Error("a and b should be equal")
	}

	if a.Equal(c) {
		t.Error("a and c should not be equal")
	}
}

func TestSet_Clone(t *testing.T) {
	original := NewSetFrom([]int{1, 2, 3})
	clone := original.Clone()

	if !original.Equal(clone) {
		t.Error("Clone should be equal to original")
	}

	// Modifying clone should not affect original
	clone.Add(4)
	if original.Contains(4) {
		t.Error("Modifying clone should not affect original")
	}
}

func TestSet_Filter(t *testing.T) {
	set := NewSetFrom([]int{1, 2, 3, 4, 5, 6})
	evens := set.Filter(func(n int) bool { return n%2 == 0 })

	if evens.Len() != 3 {
		t.Errorf("Filter Len() = %d, want 3", evens.Len())
	}

	if !evens.ContainsAll(2, 4, 6) {
		t.Error("Filter should contain even numbers")
	}
}

// ============================================================================
// BST Tests
// ============================================================================

func TestBST_InsertContains(t *testing.T) {
	tree := NewBST[int]()

	tree.Insert(5, 3, 7, 1, 4, 6, 8)

	for _, v := range []int{1, 3, 4, 5, 6, 7, 8} {
		if !tree.Contains(v) {
			t.Errorf("Tree should contain %d", v)
		}
	}

	if tree.Contains(2) {
		t.Error("Tree should not contain 2")
	}

	if tree.Len() != 7 {
		t.Errorf("Len() = %d, want 7", tree.Len())
	}
}

func TestBST_DuplicateInsert(t *testing.T) {
	tree := NewBST[int]()
	tree.Insert(5, 3, 5, 3, 7, 5)

	if tree.Len() != 3 {
		t.Errorf("Len() = %d, want 3 (duplicates should be ignored)", tree.Len())
	}
}

func TestBST_Remove(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4, 6, 8})

	// Remove leaf
	if !tree.Remove(1) {
		t.Error("Remove(1) should return true")
	}
	if tree.Contains(1) {
		t.Error("Tree should not contain 1 after removal")
	}

	// Remove node with one child
	tree.Remove(6)
	if tree.Contains(6) {
		t.Error("Tree should not contain 6 after removal")
	}

	// Remove node with two children
	tree.Remove(3)
	if tree.Contains(3) {
		t.Error("Tree should not contain 3 after removal")
	}

	// Tree structure should still be valid
	inOrder := tree.InOrder()
	prev := inOrder[0]
	for _, v := range inOrder[1:] {
		if v < prev {
			t.Error("InOrder should be sorted after removals")
		}
		prev = v
	}

	// Remove non-existent
	if tree.Remove(999) {
		t.Error("Remove non-existent should return false")
	}
}

func TestBST_MinMax(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4, 6, 8})

	min, ok := tree.Min()
	if !ok || min != 1 {
		t.Errorf("Min() = %d, %v; want 1, true", min, ok)
	}

	max, ok := tree.Max()
	if !ok || max != 8 {
		t.Errorf("Max() = %d, %v; want 8, true", max, ok)
	}

	// Empty tree
	emptyTree := NewBST[int]()
	_, ok = emptyTree.Min()
	if ok {
		t.Error("Min() on empty tree should return false")
	}
}

func TestBST_InOrder(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4, 6, 8})

	inOrder := tree.InOrder()
	want := []int{1, 3, 4, 5, 6, 7, 8}

	if !slices.Equal(inOrder, want) {
		t.Errorf("InOrder() = %v, want %v", inOrder, want)
	}
}

func TestBST_PreOrder(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4})

	preOrder := tree.PreOrder()
	want := []int{5, 3, 1, 4, 7}

	if !slices.Equal(preOrder, want) {
		t.Errorf("PreOrder() = %v, want %v", preOrder, want)
	}
}

func TestBST_PostOrder(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4})

	postOrder := tree.PostOrder()
	want := []int{1, 4, 3, 7, 5}

	if !slices.Equal(postOrder, want) {
		t.Errorf("PostOrder() = %v, want %v", postOrder, want)
	}
}

func TestBST_LevelOrder(t *testing.T) {
	tree := NewBSTFrom([]int{5, 3, 7, 1, 4, 6, 8})

	levelOrder := tree.LevelOrder()
	want := []int{5, 3, 7, 1, 4, 6, 8}

	if !slices.Equal(levelOrder, want) {
		t.Errorf("LevelOrder() = %v, want %v", levelOrder, want)
	}
}

func TestBST_Height(t *testing.T) {
	tree := NewBST[int]()

	if tree.Height() != 0 {
		t.Error("Empty tree should have height 0")
	}

	tree.Insert(5)
	if tree.Height() != 1 {
		t.Error("Single node tree should have height 1")
	}

	tree.Insert(3, 7)
	if tree.Height() != 2 {
		t.Error("Balanced 3-node tree should have height 2")
	}

	// Create skewed tree
	skewed := NewBSTFrom([]int{1, 2, 3, 4, 5})
	if skewed.Height() != 5 {
		t.Errorf("Skewed tree Height() = %d, want 5", skewed.Height())
	}
}

func TestBST_Clear(t *testing.T) {
	tree := NewBSTFrom([]int{1, 2, 3, 4, 5})

	tree.Clear()

	if !tree.IsEmpty() {
		t.Error("Tree should be empty after Clear()")
	}

	if tree.Len() != 0 {
		t.Error("Len() should be 0 after Clear()")
	}
}

func TestBST_String(t *testing.T) {
	tree := NewBSTFrom([]string{"cherry", "apple", "banana"})

	if !tree.Contains("apple") || !tree.Contains("banana") || !tree.Contains("cherry") {
		t.Error("BST should work with strings")
	}

	inOrder := tree.InOrder()
	want := []string{"apple", "banana", "cherry"}

	if !slices.Equal(inOrder, want) {
		t.Errorf("InOrder() = %v, want %v", inOrder, want)
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkStack_PushPop(b *testing.B) {
	stack := NewStack[int]()

	for i := 0; i < b.N; i++ {
		stack.Push(i)
		stack.Pop()
	}
}

func BenchmarkQueue_EnqueueDequeue(b *testing.B) {
	queue := NewQueue[int]()

	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
		queue.Dequeue()
	}
}

func BenchmarkSet_Add(b *testing.B) {
	set := NewSet[int]()

	for i := 0; i < b.N; i++ {
		set.Add(i % 1000)
	}
}

func BenchmarkSet_Contains(b *testing.B) {
	set := NewSetWithCapacity[int](10000)
	for i := 0; i < 10000; i++ {
		set.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contains(i % 10000)
	}
}

func BenchmarkBST_Insert(b *testing.B) {
	tree := NewBST[int]()

	for i := 0; i < b.N; i++ {
		tree.Insert(i)
	}
}

func BenchmarkBST_Contains(b *testing.B) {
	tree := NewBST[int]()
	for i := 0; i < 10000; i++ {
		tree.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Contains(i % 10000)
	}
}

func BenchmarkBST_InOrder(b *testing.B) {
	tree := NewBST[int]()
	for i := 0; i < 10000; i++ {
		tree.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.InOrder()
	}
}
