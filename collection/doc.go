// Package collection provides generic data structures for Go.
//
// All data structures in this package use Go generics to work with
// any type that satisfies the required constraints.
//
// # Stack
//
// A LIFO (Last-In-First-Out) data structure:
//
//	stack := collection.NewStack[int]()
//	stack.Push(1)
//	stack.Push(2)
//	val, ok := stack.Pop()  // val=2, ok=true
//
// # Queue
//
// A FIFO (First-In-First-Out) data structure:
//
//	queue := collection.NewQueue[string]()
//	queue.Enqueue("first")
//	queue.Enqueue("second")
//	val, ok := queue.Dequeue()  // val="first", ok=true
//
// # Set
//
// An unordered collection of unique elements:
//
//	set := collection.NewSet[int]()
//	set.Add(1, 2, 3)
//	set.Contains(2)  // true
//	set.Remove(2)
//
// # BST (Binary Search Tree)
//
// An ordered tree structure for efficient search, insert, and delete:
//
//	tree := collection.NewBST[int]()
//	tree.Insert(5, 3, 7, 1, 4)
//	tree.Contains(3)  // true
//	values := tree.InOrder()  // [1, 3, 4, 5, 7]
//
// # Thread Safety
//
// The data structures in this package are NOT thread-safe by default.
// For concurrent access, use sync primitives around Stack, Queue, Set, and BST.
//
// # Design Decisions
//
// - Generic: All structures use type parameters
// - Zero values work: NewXxx() returns ready-to-use structures
// - Optional returns: Pop/Dequeue return (value, bool) instead of panicking
// - Immutable iterators: Values() returns copies, not references
package collection
