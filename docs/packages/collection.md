---
title: collection
parent: Packages
nav_order: 8
---

# collection

Generic data structures for Go.
{: .fs-6 .fw-300 }

The `collection` package provides well-tested, generic implementations of fundamental data structures. All structures use Go generics to work with any type that satisfies the required constraints.

## Import

```go
import "github.com/alessiosavi/GoGPUtils/collection"
```

## Design Principles

- **Generic**: All structures use type parameters — no `interface{}` or code generation
- **Zero values work**: `NewXxx()` returns ready-to-use structures
- **Optional returns**: `Pop`/`Dequeue` return `(value, bool)` instead of panicking on empty structures
- **Immutable iterators**: `Values()` returns copies, not references
- **Not thread-safe**: For concurrent access, wrap operations with `sync` primitives

---

## Stack

A **LIFO** (Last-In-First-Out) data structure backed by a slice.

### Type Signature

```go
type Stack[T any] struct{ ... }

func NewStack[T any]() *Stack[T]
func NewStackWithCapacity[T any](capacity int) *Stack[T]
```

### Methods

| Method    | Signature                                | Description                                    |
| --------- | ---------------------------------------- | ---------------------------------------------- |
| `Push`    | `func (s *Stack[T]) Push(item T)`        | Add element to top                             |
| `PushAll` | `func (s *Stack[T]) PushAll(items ...T)` | Add multiple elements                          |
| `Pop`     | `func (s *Stack[T]) Pop() (T, bool)`     | Remove and return top element                  |
| `PopN`    | `func (s *Stack[T]) PopN(n int) []T`     | Remove and return up to `n` elements           |
| `PopAll`  | `func (s *Stack[T]) PopAll() []T`        | Remove and return all elements (top to bottom) |
| `Peek`    | `func (s *Stack[T]) Peek() (T, bool)`    | Return top element without removing            |
| `Len`     | `func (s *Stack[T]) Len() int`           | Number of elements                             |
| `IsEmpty` | `func (s *Stack[T]) IsEmpty() bool`      | True if no elements                            |
| `Clear`   | `func (s *Stack[T]) Clear()`             | Remove all elements                            |
| `Values`  | `func (s *Stack[T]) Values() []T`        | Copy of all elements (bottom to top)           |

### Time Complexity

| Operation         | Time           | Space |
| ----------------- | -------------- | ----- |
| `Push`            | O(1) amortized | O(1)  |
| `Pop`             | O(1)           | O(1)  |
| `Peek`            | O(1)           | O(1)  |
| `Len` / `IsEmpty` | O(1)           | O(1)  |
| `Values`          | O(n)           | O(n)  |
| `PopAll`          | O(n)           | O(n)  |
| `Clear`           | O(1)           | O(1)  |

### Usage Examples

```go
// Basic push and pop
stack := collection.NewStack[int]()
stack.Push(1)
stack.Push(2)
stack.Push(3)

val, ok := stack.Pop()   // val=3, ok=true
val, ok = stack.Peek()   // val=2, ok=true (stack unchanged)

// Bulk operations
stack.PushAll(4, 5, 6)
all := stack.PopAll()    // [6, 5, 4, 2, 1] — LIFO order

// Safe empty handling
empty := collection.NewStack[string]()
_, ok = empty.Pop()      // ok=false, no panic
```

---

## Queue

A **FIFO** (First-In-First-Out) data structure backed by a slice.

### Type Signature

```go
type Queue[T any] struct{ ... }

func NewQueue[T any]() *Queue[T]
func NewQueueWithCapacity[T any](capacity int) *Queue[T]
```

### Methods

| Method       | Signature                                   | Description                           |
| ------------ | ------------------------------------------- | ------------------------------------- |
| `Enqueue`    | `func (q *Queue[T]) Enqueue(item T)`        | Add element to back                   |
| `EnqueueAll` | `func (q *Queue[T]) EnqueueAll(items ...T)` | Add multiple elements                 |
| `Dequeue`    | `func (q *Queue[T]) Dequeue() (T, bool)`    | Remove and return front element       |
| `DequeueN`   | `func (q *Queue[T]) DequeueN(n int) []T`    | Remove and return up to `n` elements  |
| `Peek`       | `func (q *Queue[T]) Peek() (T, bool)`       | Return front element without removing |
| `Len`        | `func (q *Queue[T]) Len() int`              | Number of elements                    |
| `IsEmpty`    | `func (q *Queue[T]) IsEmpty() bool`         | True if no elements                   |
| `Clear`      | `func (q *Queue[T]) Clear()`                | Remove all elements                   |
| `Values`     | `func (q *Queue[T]) Values() []T`           | Copy of all elements in FIFO order    |

### Time Complexity

| Operation         | Time           | Space |
| ----------------- | -------------- | ----- |
| `Enqueue`         | O(1) amortized | O(1)  |
| `Dequeue`         | O(n)*          | O(1)  |
| `Peek`            | O(1)           | O(1)  |
| `Len` / `IsEmpty` | O(1)           | O(1)  |
| `Values`          | O(n)           | O(n)  |
| `Clear`           | O(1)           | O(1)  |

\* `Dequeue` is O(n) due to slice re-slicing (`items = items[1:]`). For high-throughput scenarios, consider a ring buffer or linked-list queue.

### Usage Examples

```go
// Basic enqueue and dequeue
queue := collection.NewQueue[string]()
queue.Enqueue("first")
queue.Enqueue("second")
queue.Enqueue("third")

val, ok := queue.Dequeue()  // val="first", ok=true
val, ok = queue.Peek()      // val="second", ok=true (queue unchanged)

// Bulk operations
queue.EnqueueAll("fourth", "fifth")
firstThree := queue.DequeueN(3)  // ["second", "third", "fourth"]

// Safe empty handling
empty := collection.NewQueue[int]()
_, ok = empty.Dequeue()     // ok=false, no panic
```

---

## Set

An **unordered collection of unique elements** backed by a Go map.

### Type Signature

```go
type Set[T comparable] struct{ ... }

func NewSet[T comparable]() *Set[T]
func NewSetWithCapacity[T comparable](capacity int) *Set[T]
func NewSetFrom[T comparable](items []T) *Set[T]
```

### Methods

| Method                | Signature                                                     | Description                         |
| --------------------- | ------------------------------------------------------------- | ----------------------------------- |
| `Add`                 | `func (s *Set[T]) Add(items ...T)`                            | Add one or more elements            |
| `Remove`              | `func (s *Set[T]) Remove(items ...T)`                         | Remove one or more elements         |
| `Contains`            | `func (s *Set[T]) Contains(item T) bool`                      | Check membership                    |
| `ContainsAll`         | `func (s *Set[T]) ContainsAll(items ...T) bool`               | Check if all items are present      |
| `ContainsAny`         | `func (s *Set[T]) ContainsAny(items ...T) bool`               | Check if any item is present        |
| `Len`                 | `func (s *Set[T]) Len() int`                                  | Number of elements                  |
| `IsEmpty`             | `func (s *Set[T]) IsEmpty() bool`                             | True if no elements                 |
| `Clear`               | `func (s *Set[T]) Clear()`                                    | Remove all elements                 |
| `Values`              | `func (s *Set[T]) Values() []T`                               | All elements (order not guaranteed) |
| `Clone`               | `func (s *Set[T]) Clone() *Set[T]`                            | Deep copy of the set                |
| `ForEach`             | `func (s *Set[T]) ForEach(fn func(T))`                        | Iterate over elements               |
| `Filter`              | `func (s *Set[T]) Filter(predicate func(T) bool) *Set[T]`     | New set with matching elements      |
| `Union`               | `func (s *Set[T]) Union(other *Set[T]) *Set[T]`               | Elements in either set              |
| `Intersection`        | `func (s *Set[T]) Intersection(other *Set[T]) *Set[T]`        | Elements in both sets               |
| `Difference`          | `func (s *Set[T]) Difference(other *Set[T]) *Set[T]`          | Elements in `s` but not `other`     |
| `SymmetricDifference` | `func (s *Set[T]) SymmetricDifference(other *Set[T]) *Set[T]` | Elements in exactly one set         |
| `IsSubset`            | `func (s *Set[T]) IsSubset(other *Set[T]) bool`               | All elements of `s` in `other`      |
| `IsSuperset`          | `func (s *Set[T]) IsSuperset(other *Set[T]) bool`             | `s` contains all of `other`         |
| `Equal`               | `func (s *Set[T]) Equal(other *Set[T]) bool`                  | Same elements                       |

### Package-Level Functions

| Function        | Signature                                          | Description              |
| --------------- | -------------------------------------------------- | ------------------------ |
| `ToSliceSorted` | `func ToSliceSorted[T cmp.Ordered](s *Set[T]) []T` | Sorted slice of elements |

### Time Complexity

| Operation                 | Time                 | Space                |
| ------------------------- | -------------------- | -------------------- |
| `Add`                     | O(1) amortized       | O(1)                 |
| `Remove`                  | O(1)                 | O(1)                 |
| `Contains`                | O(1)                 | O(1)                 |
| `Len` / `IsEmpty`         | O(1)                 | O(1)                 |
| `Values`                  | O(n)                 | O(n)                 |
| `Union`                   | O(\|a\| + \|b\|)     | O(\|a\| + \|b\|)     |
| `Intersection`            | O(min(\|a\|, \|b\|)) | O(min(\|a\|, \|b\|)) |
| `Difference`              | O(\|a\|)             | O(\|a\|)             |
| `SymmetricDifference`     | O(\|a\| + \|b\|)     | O(\|a\| + \|b\|)     |
| `IsSubset` / `IsSuperset` | O(\|a\|)             | O(1)                 |
| `Equal`                   | O(min(\|a\|, \|b\|)) | O(1)                 |
| `Clone`                   | O(n)                 | O(n)                 |
| `Filter`                  | O(n)                 | O(k)                 |
| `ToSliceSorted`           | O(n log n)           | O(n)                 |

### Usage Examples

```go
// Basic operations
set := collection.NewSet[int]()
set.Add(1, 2, 3)
set.Add(2)               // Duplicate, ignored

set.Contains(2)          // true
set.Contains(5)          // false
set.Len()                // 3

set.Remove(2)
set.Contains(2)          // false

// Create from slice
fruits := collection.NewSetFrom([]string{"apple", "banana", "apple"})
// fruits = {"apple", "banana"}
```

### Set Operations

```go
a := collection.NewSetFrom([]int{1, 2, 3})
b := collection.NewSetFrom([]int{2, 3, 4})

// Union: elements in either set
union := a.Union(b)              // {1, 2, 3, 4}

// Intersection: elements in both sets
inter := a.Intersection(b)       // {2, 3}

// Difference: elements in a but not in b
diff := a.Difference(b)          // {1}

// Symmetric difference: elements in exactly one set
symDiff := a.SymmetricDifference(b)  // {1, 4}

// Subset / Superset
small := collection.NewSetFrom([]int{1, 2})
large := collection.NewSetFrom([]int{1, 2, 3, 4})

small.IsSubset(large)     // true
large.IsSuperset(small)   // true
small.IsSuperset(large)   // false

// Equality
x := collection.NewSetFrom([]int{1, 2, 3})
y := collection.NewSetFrom([]int{3, 1, 2})
x.Equal(y)                // true (order doesn't matter)
```

### Filtering and Iteration

```go
set := collection.NewSetFrom([]int{1, 2, 3, 4, 5, 6})

// Filter to even numbers
evens := set.Filter(func(n int) bool {
    return n%2 == 0
})  // {2, 4, 6}

// Iterate
set.ForEach(func(n int) {
    fmt.Println(n)
})

// Sorted output (deterministic order)
sorted := collection.ToSliceSorted(set)  // [1, 2, 3, 4, 5, 6]
```

---

## BST (Binary Search Tree)

An **ordered tree structure** for efficient search, insert, and delete operations. Elements must satisfy the `cmp.Ordered` constraint.

### Type Signature

```go
type BST[T cmp.Ordered] struct{ ... }

func NewBST[T cmp.Ordered]() *BST[T]
func NewBSTFrom[T cmp.Ordered](items []T) *BST[T]
```

### Methods

| Method        | Signature                                      | Description                              |
| ------------- | ---------------------------------------------- | ---------------------------------------- |
| `Insert`      | `func (t *BST[T]) Insert(values ...T)`         | Add one or more values                   |
| `Contains`    | `func (t *BST[T]) Contains(value T) bool`      | Check membership                         |
| `Remove`      | `func (t *BST[T]) Remove(value T) bool`        | Delete a value (returns true if found)   |
| `Min`         | `func (t *BST[T]) Min() (T, bool)`             | Smallest value                           |
| `Max`         | `func (t *BST[T]) Max() (T, bool)`             | Largest value                            |
| `Len`         | `func (t *BST[T]) Len() int`                   | Number of elements                       |
| `IsEmpty`     | `func (t *BST[T]) IsEmpty() bool`              | True if no elements                      |
| `Clear`       | `func (t *BST[T]) Clear()`                     | Remove all elements                      |
| `Height`      | `func (t *BST[T]) Height() int`                | Tree height (0 = empty, 1 = single node) |
| `InOrder`     | `func (t *BST[T]) InOrder() []T`               | Sorted values (left, root, right)        |
| `PreOrder`    | `func (t *BST[T]) PreOrder() []T`              | Pre-order traversal (root, left, right)  |
| `PostOrder`   | `func (t *BST[T]) PostOrder() []T`             | Post-order traversal (left, right, root) |
| `LevelOrder`  | `func (t *BST[T]) LevelOrder() []T`            | Breadth-first traversal                  |
| `Values`      | `func (t *BST[T]) Values() []T`                | Alias for `InOrder()` — sorted values    |
| `ForEach`     | `func (t *BST[T]) ForEach(fn func(T))`         | Iterate in sorted order                  |
| `RangeSearch` | `func (t *BST[T]) RangeSearch(min, max T) []T` | Values in inclusive range [min, max]     |

### Time Complexity

| Operation                | Average      | Worst | Space     |
| ------------------------ | ------------ | ----- | --------- |
| `Insert`                 | O(log n)     | O(n)* | O(1)      |
| `Contains`               | O(log n)     | O(n)* | O(1)      |
| `Remove`                 | O(log n)     | O(n)* | O(1)      |
| `Min` / `Max`            | O(log n)     | O(n)* | O(1)      |
| `Height`                 | O(n)         | O(n)  | O(log n)† |
| `InOrder` / `Values`     | O(n)         | O(n)  | O(n)      |
| `PreOrder` / `PostOrder` | O(n)         | O(n)  | O(n)      |
| `LevelOrder`             | O(n)         | O(n)  | O(w)‡     |
| `RangeSearch`            | O(k + log n) | O(n)  | O(k)      |
| `Len` / `IsEmpty`        | O(1)         | O(1)  | O(1)      |

\* Worst case occurs when the tree becomes skewed (e.g., inserting sorted data). For guaranteed O(log n), use a self-balancing tree (AVL, Red-Black).

† `Height` uses recursion; stack space is O(log n) balanced, O(n) skewed.

‡ `LevelOrder` uses a queue; space is O(w) where w = maximum tree width.

### Usage Examples

```go
// Basic operations
tree := collection.NewBST[int]()
tree.Insert(5, 3, 7, 1, 4, 6, 8)

tree.Contains(3)   // true
tree.Contains(9)   // false
tree.Len()         // 7

// Min and Max
min, _ := tree.Min()  // 1
max, _ := tree.Max()  // 8

// Traversals
tree.InOrder()     // [1, 3, 4, 5, 6, 7, 8]  — sorted
tree.PreOrder()    // [5, 3, 1, 4, 7, 6, 8]  — root first
tree.PostOrder()   // [1, 4, 3, 6, 8, 7, 5]  — root last
tree.LevelOrder()  // [5, 3, 7, 1, 4, 6, 8]  — breadth-first

// Range search
tree.RangeSearch(3, 7)  // [3, 4, 5, 6, 7]

// Remove
tree.Remove(3)     // true (found and removed)
tree.Remove(99)    // false (not found)

// Create from slice
words := collection.NewBSTFrom([]string{"cherry", "apple", "banana"})
words.InOrder()    // ["apple", "banana", "cherry"]
```

### Iteration

```go
tree := collection.NewBSTFrom([]int{3, 1, 4, 1, 5})

// Iterate in sorted order
tree.ForEach(func(n int) {
    fmt.Println(n)
})
// Output: 1, 3, 4, 5

// Values is an alias for InOrder
sorted := tree.Values()  // [1, 3, 4, 5]
```

---

## Thread Safety

The data structures in this package are **not thread-safe** by default. For concurrent access, use `sync.Mutex` or `sync.RWMutex`:

```go
type SafeStack[T any] struct {
    mu    sync.Mutex
    stack *collection.Stack[T]
}

func (s *SafeStack[T]) Push(item T) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.stack.Push(item)
}
```

---

## Benchmarks

Run benchmarks with:

```bash
go test -bench=. ./collection/
```

| Benchmark                       | Description                        |
| ------------------------------- | ---------------------------------- |
| `BenchmarkStack_PushPop`        | Push/pop cycle                     |
| `BenchmarkQueue_EnqueueDequeue` | Enqueue/dequeue cycle              |
| `BenchmarkSet_Add`              | Add elements (with duplicates)     |
| `BenchmarkSet_Contains`         | Membership test on 10k elements    |
| `BenchmarkBST_Insert`           | Sequential insert                  |
| `BenchmarkBST_Contains`         | Search in 10k element tree         |
| `BenchmarkBST_InOrder`          | In-order traversal of 10k elements |
