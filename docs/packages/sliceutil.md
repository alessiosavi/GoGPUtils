---
title: sliceutil
parent: Packages
nav_order: 1
---

# sliceutil

Generic slice operations for Go.
{: .fs-6 .fw-300 }

## Overview

The `sliceutil` package provides generic slice operations that work with any `comparable` or `ordered` type through Go's generics. It is designed with the following principles:

- **Generic**: Works with any type that satisfies the constraints
- **Non-mutating by default**: Functions return new slices unless suffixed with `InPlace`
- **Nil-safe**: All functions handle nil slices gracefully
- **Zero allocations**: Where possible, operations avoid heap allocations

## Installation

```go
import "github.com/alessiosavi/GoGPUtils/sliceutil"
```

## Function Reference

### Filtering

#### Filter

Returns a new slice containing only elements for which the predicate returns true. Returns nil if the input slice is nil.

```go
func Filter[T any](s []T, predicate func(T) bool) []T
```

**Example:**

```go
evens := sliceutil.Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
// evens = [2, 4]
```

#### FilterInPlace

Filters the slice in place, returning the modified slice. The underlying array is modified; elements are shifted to fill gaps. This avoids allocation but modifies the original slice.

```go
func FilterInPlace[T any](s []T, predicate func(T) bool) []T
```

**Example:**

```go
nums := []int{1, 2, 3, 4, 5}
evens := sliceutil.FilterInPlace(nums, func(n int) bool { return n%2 == 0 })
// evens = [2, 4], nums[0:2] = [2, 4]
```

---

### Transformation

#### Map

Applies a transformation function to each element and returns a new slice. Returns nil if the input slice is nil.

```go
func Map[T, U any](s []T, transform func(T) U) []U
```

**Example:**

```go
doubled := sliceutil.Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
// doubled = [2, 4, 6]
```

#### MapWithIndex

Applies a transformation function that receives both index and element.

```go
func MapWithIndex[T, U any](s []T, transform func(int, T) U) []U
```

**Example:**

```go
indexed := sliceutil.MapWithIndex([]string{"a", "b"}, func(i int, s string) string {
    return fmt.Sprintf("%d:%s", i, s)
})
// indexed = ["0:a", "1:b"]
```

#### MapErr

Applies fn to each element of s and returns the transformed slice. Returns the first error encountered (with a nil result slice); returns nil result and nil error for a nil input slice.

```go
func MapErr[T, U any](s []T, fn func(T) (U, error)) ([]U, error)
```

**Example:**

```go
nums, err := sliceutil.MapErr([]string{"1", "2", "three"}, strconv.Atoi)
// err != nil on the third element; nums is nil
```

#### FlatMap

Applies a function that returns a slice and flattens the result.

```go
func FlatMap[T, U any](s []T, transform func(T) []U) []U
```

**Example:**

```go
words := sliceutil.FlatMap([]string{"hello world", "foo bar"}, func(s string) []string {
    return strings.Split(s, " ")
})
// words = ["hello", "world", "foo", "bar"]
```

---

### Reduction

#### Reduce

Reduces a slice to a single value using an accumulator function.

```go
func Reduce[T, U any](s []T, initial U, accumulator func(U, T) U) U
```

**Example:**

```go
sum := sliceutil.Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
// sum = 10
```

---

### Search and Index

#### Contains

Reports whether the slice contains the target element. Uses `==` for comparison, requiring comparable types.

```go
func Contains[T comparable](s []T, target T) bool
```

**Example:**

```go
if sliceutil.Contains(names, "Alice") { ... }
```

#### ContainsFunc

Reports whether any element satisfies the predicate.

```go
func ContainsFunc[T any](s []T, predicate func(T) bool) bool
```

**Example:**

```go
hasNegative := sliceutil.ContainsFunc(nums, func(n int) bool { return n < 0 })
```

#### IndexOf

Returns the index of the first occurrence of target, or -1 if not found.

```go
func IndexOf[T comparable](s []T, target T) int
```

**Example:**

```go
idx := sliceutil.IndexOf([]string{"a", "b", "c"}, "b")
// idx = 1
```

#### IndexOfFunc

Returns the index of the first element satisfying the predicate, or -1.

```go
func IndexOfFunc[T any](s []T, predicate func(T) bool) int
```

#### LastIndexOf

Returns the index of the last occurrence of target, or -1 if not found.

```go
func LastIndexOf[T comparable](s []T, target T) int
```

#### Find

Returns the first element satisfying the predicate and true, or zero value and false.

```go
func Find[T any](s []T, predicate func(T) bool) (T, bool)
```

**Example:**

```go
first, ok := sliceutil.Find(users, func(u User) bool { return u.Age > 18 })
```

#### FindLast

Returns the last element satisfying the predicate.

```go
func FindLast[T any](s []T, predicate func(T) bool) (T, bool)
```

---

### Uniqueness

#### Unique

Returns a new slice with duplicate elements removed. The first occurrence of each element is kept; order is preserved.

```go
func Unique[T comparable](s []T) []T
```

**Example:**

```go
unique := sliceutil.Unique([]int{1, 2, 2, 3, 1})
// unique = [1, 2, 3]
```

#### UniqueFunc

Returns a new slice with duplicates removed based on a key function. Elements with the same key are considered duplicates; first occurrence is kept.

```go
func UniqueFunc[T any, K comparable](s []T, key func(T) K) []T
```

**Example:**

```go
type User struct { ID int; Name string }
users := []User{ {1, "Alice"}, {2, "Bob"}, {1, "Alice2"} }
unique := sliceutil.UniqueFunc(users, func(u User) int { return u.ID })
// unique = [{1, "Alice"}, {2, "Bob"}]
```

#### Compact

Removes consecutive duplicate elements, similar to Unix `uniq`. For removing all duplicates, use `Unique`.

```go
func Compact[T comparable](s []T) []T
```

**Example:**

```go
compacted := sliceutil.Compact([]int{1, 1, 2, 2, 2, 1})
// compacted = [1, 2, 1]
```

#### CompactFunc

Removes consecutive elements where the predicate returns true.

```go
func CompactFunc[T any](s []T, eq func(T, T) bool) []T
```

---

### Chunking and Grouping

#### Chunk

Splits a slice into chunks of the specified size. The last chunk may be smaller if `len(s)` is not divisible by size. Returns nil if size <= 0 or s is nil.

```go
func Chunk[T any](s []T, size int) [][]T
```

**Example:**

```go
chunks := sliceutil.Chunk([]int{1, 2, 3, 4, 5}, 2)
// chunks = [[1, 2], [3, 4], [5]]
```

#### Flatten

Converts a slice of slices into a single slice.

```go
func Flatten[T any](s [][]T) []T
```

**Example:**

```go
flat := sliceutil.Flatten([][]int{ {1, 2}, {3}, {4, 5} })
// flat = [1, 2, 3, 4, 5]
```

#### GroupBy

Groups elements by a key function.

```go
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T
```

**Example:**

```go
nums := []int{1, 2, 3, 4, 5, 6}
groups := sliceutil.GroupBy(nums, func(n int) string {
    if n%2 == 0 { return "even" }
    return "odd"
})
// groups = map["even":[2, 4, 6] "odd":[1, 3, 5]]
```

#### Partition

Splits a slice into two: elements matching the predicate and those that don't.

```go
func Partition[T any](s []T, predicate func(T) bool) (matching, notMatching []T)
```

**Example:**

```go
evens, odds := sliceutil.Partition([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
// evens = [2, 4], odds = [1, 3]
```

---

### Reordering

#### Reverse

Returns a new slice with elements in reverse order.

```go
func Reverse[T any](s []T) []T
```

**Example:**

```go
rev := sliceutil.Reverse([]int{1, 2, 3})
// rev = [3, 2, 1]
```

#### ReverseInPlace

Reverses the slice in place and returns it.

```go
func ReverseInPlace[T any](s []T) []T
```

**Example:**

```go
nums := []int{1, 2, 3}
sliceutil.ReverseInPlace(nums)
// nums = [3, 2, 1]
```

#### Shuffle

Returns a new slice with elements in random order. Uses `math/rand` (NOT `crypto/rand`) for performance; not suitable for security-sensitive use.

```go
func Shuffle[T any](s []T) []T
```

#### ShuffleInPlace

Shuffles the slice in place using Fisher-Yates algorithm. Uses `math/rand` for performance; use `randutil.Shuffle` for crypto-secure randomness.

```go
func ShuffleInPlace[T any](s []T)
```

#### SeedShuffle

Sets the seed for `ShuffleInPlace`. Useful for deterministic tests.

```go
func SeedShuffle(seed uint64)
```

---

### Set Operations

#### Intersect

Returns elements present in both slices. Result preserves order from the first slice.

```go
func Intersect[T comparable](a, b []T) []T
```

**Example:**

```go
common := sliceutil.Intersect([]int{1, 2, 3}, []int{2, 3, 4})
// common = [2, 3]
```

#### Difference

Returns elements in `a` that are not in `b`.

```go
func Difference[T comparable](a, b []T) []T
```

**Example:**

```go
diff := sliceutil.Difference([]int{1, 2, 3}, []int{2, 3, 4})
// diff = [1]
```

#### Union

Returns all unique elements from both slices.

```go
func Union[T comparable](a, b []T) []T
```

**Example:**

```go
all := sliceutil.Union([]int{1, 2}, []int{2, 3})
// all = [1, 2, 3]
```

---

### Slicing

#### Take

Returns the first n elements. If n > len(s), returns all elements.

```go
func Take[T any](s []T, n int) []T
```

**Example:**

```go
first := sliceutil.Take([]int{1, 2, 3, 4, 5}, 3)
// first = [1, 2, 3]
```

#### TakeLast

Returns the last n elements.

```go
func TakeLast[T any](s []T, n int) []T
```

**Example:**

```go
last := sliceutil.TakeLast([]int{1, 2, 3, 4, 5}, 3)
// last = [3, 4, 5]
```

#### Drop

Returns elements after skipping the first n. If n >= len(s), returns empty slice.

```go
func Drop[T any](s []T, n int) []T
```

**Example:**

```go
rest := sliceutil.Drop([]int{1, 2, 3, 4, 5}, 2)
// rest = [3, 4, 5]
```

#### DropLast

Returns elements after removing the last n.

```go
func DropLast[T any](s []T, n int) []T
```

#### TakeWhile

Returns elements from the start while predicate returns true.

```go
func TakeWhile[T any](s []T, predicate func(T) bool) []T
```

**Example:**

```go
result := sliceutil.TakeWhile([]int{1, 2, 3, 4, 1}, func(n int) bool { return n < 3 })
// result = [1, 2]
```

#### DropWhile

Returns elements after dropping from the start while predicate returns true.

```go
func DropWhile[T any](s []T, predicate func(T) bool) []T
```

---

### Predicates

#### All

Returns true if all elements satisfy the predicate. Returns true for empty slices.

```go
func All[T any](s []T, predicate func(T) bool) bool
```

**Example:**

```go
allPositive := sliceutil.All([]int{1, 2, 3}, func(n int) bool { return n > 0 })
// allPositive = true
```

#### Any

Returns true if any element satisfies the predicate. Returns false for empty slices.

```go
func Any[T any](s []T, predicate func(T) bool) bool
```

**Example:**

```go
hasNegative := sliceutil.Any([]int{1, -2, 3}, func(n int) bool { return n < 0 })
// hasNegative = true
```

#### None

Returns true if no elements satisfy the predicate. Returns true for empty slices.

```go
func None[T any](s []T, predicate func(T) bool) bool
```

#### Count

Returns the number of elements satisfying the predicate.

```go
func Count[T any](s []T, predicate func(T) bool) int
```

---

### Min / Max

#### Min

Returns the minimum element using natural ordering. Returns zero value and false for empty slices.

```go
func Min[T cmp.Ordered](s []T) (T, bool)
```

#### Max

Returns the maximum element using natural ordering. Returns zero value and false for empty slices.

```go
func Max[T cmp.Ordered](s []T) (T, bool)
```

#### MinFunc

Returns the minimum element using a comparison function. `cmp` should return negative if a < b, zero if a == b, positive if a > b.

```go
func MinFunc[T any](s []T, cmpFn func(a, b T) int) (T, bool)
```

#### MaxFunc

Returns the maximum element using a comparison function.

```go
func MaxFunc[T any](s []T, cmpFn func(a, b T) int) (T, bool)
```

---

### Equality

#### Equal

Reports whether two slices are equal.

```go
func Equal[T comparable](a, b []T) bool
```

#### EqualFunc

Reports whether two slices are equal using a custom comparison.

```go
func EqualFunc[T, U any](a []T, b []U, eq func(T, U) bool) bool
```

---

### Iteration

#### ForEach

Calls a function for each element.

```go
func ForEach[T any](s []T, fn func(T))
```

**Example:**

```go
sliceutil.ForEach(items, func(item Item) { process(item) })
```

#### ForEachWithIndex

Calls a function for each element with its index.

```go
func ForEachWithIndex[T any](s []T, fn func(int, T))
```

---

### Zipping

#### Zip

Combines two slices into pairs. Stops at the shorter slice.

```go
func Zip[T, U any](a []T, b []U) [][2]any
```

**Example:**

```go
type Pair[T, U any] struct { First T; Second U }
pairs := sliceutil.Zip([]int{1, 2}, []string{"a", "b"})
// pairs = [{1, "a"}, {2, "b"}]
```

#### ZipWith

Combines two slices using a function.

```go
func ZipWith[T, U, V any](a []T, b []U, combine func(T, U) V) []V
```

**Example:**

```go
sums := sliceutil.ZipWith([]int{1, 2, 3}, []int{4, 5, 6}, func(a, b int) int { return a + b })
// sums = [5, 7, 9]
```

---

### Padding

#### Pad

Extends a slice to the target length by appending the fill value. If slice is already >= length, returns a clone.

```go
func Pad[T any](s []T, length int, fill T) []T
```

**Example:**

```go
padded := sliceutil.Pad([]int{1, 2}, 5, 0)
// padded = [1, 2, 0, 0, 0]
```

#### PadLeft

Extends a slice by prepending the fill value.

```go
func PadLeft[T any](s []T, length int, fill T) []T
```

**Example:**

```go
padded := sliceutil.PadLeft([]int{1, 2}, 5, 0)
// padded = [0, 0, 0, 1, 2]
```

---

### Removal and Insertion

#### RemoveAt

Returns a new slice with the element at index removed. Returns nil if index is out of bounds.

```go
func RemoveAt[T any](s []T, index int) []T
```

**Example:**

```go
removed := sliceutil.RemoveAt([]int{1, 2, 3}, 1)
// removed = [1, 3]
```

#### RemoveValue

Returns a new slice with all occurrences of value removed.

```go
func RemoveValue[T comparable](s []T, value T) []T
```

**Example:**

```go
cleaned := sliceutil.RemoveValue([]int{1, 2, 3, 2}, 2)
// cleaned = [1, 3]
```

#### RemoveFirst

Returns a new slice with the first occurrence of value removed.

```go
func RemoveFirst[T comparable](s []T, value T) []T
```

#### Insert

Returns a new slice with value inserted at index. If index is out of bounds, appends to end.

```go
func Insert[T any](s []T, index int, value T) []T
```

**Example:**

```go
inserted := sliceutil.Insert([]int{1, 3}, 1, 2)
// inserted = [1, 2, 3]
```

---

### Association

#### Associate

Creates a map from slice elements using a key function.

```go
func Associate[T any, K comparable](s []T, key func(T) K) map[K]T
```

**Example:**

```go
users := []User{ {ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"} }
byID := sliceutil.Associate(users, func(u User) int { return u.ID })
// byID = map[1:{1, "Alice"} 2:{2, "Bob"}]
```

#### AssociateWith

Creates a map from slice elements to values computed by a function.

```go
func AssociateWith[K comparable, V any](s []K, value func(K) V) map[K]V
```

**Example:**

```go
lengths := sliceutil.AssociateWith([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) })
// lengths = map["a":1 "bb":2 "ccc":3]
```

---

## Go Built-ins vs sliceutil

| Operation        | Go Built-in / Standard Library                       | sliceutil Adds                                                   |
| ---------------- | ---------------------------------------------------- | ---------------------------------------------------------------- |
| Length           | `len(slice)`                                         | —                                                                |
| Append           | `append(slice, items...)`                            | —                                                                |
| Copy             | `copy(dst, src)`                                     | —                                                                |
| Slice expression | `slice[start:end]`                                   | —                                                                |
| Sort             | `slices.Sort`, `sort.Slice`                          | —                                                                |
| Binary search    | `slices.BinarySearch`                                | —                                                                |
| Contains         | `slices.Contains` (Go 1.21+)                         | `ContainsFunc`, `IndexOf`, `LastIndexOf`, `IndexOfFunc`          |
| Filter           | Manual loop                                          | `Filter`, `FilterInPlace`                                        |
| Map              | Manual loop                                          | `Map`, `MapWithIndex`, `MapErr`                                  |
| Reduce           | Manual loop                                          | `Reduce`                                                         |
| Unique           | Manual loop + map                                    | `Unique`, `UniqueFunc`, `Compact`, `CompactFunc`                 |
| Chunk            | Manual loop                                          | `Chunk`                                                          |
| Flatten          | Manual loop                                          | `Flatten`, `FlatMap`                                             |
| GroupBy          | Manual loop + map                                    | `GroupBy`                                                        |
| Partition        | Manual loop                                          | `Partition`                                                      |
| Reverse          | Manual loop                                          | `Reverse`, `ReverseInPlace`                                      |
| Shuffle          | `rand.Shuffle`                                       | `Shuffle`, `ShuffleInPlace`, `SeedShuffle`                       |
| Set ops          | Manual loop + map                                    | `Intersect`, `Difference`, `Union`                               |
| Take/Drop        | Slice expressions                                    | `Take`, `TakeLast`, `Drop`, `DropLast`, `TakeWhile`, `DropWhile` |
| Predicates       | Manual loop                                          | `All`, `Any`, `None`, `Count`                                    |
| Find             | Manual loop                                          | `Find`, `FindLast`                                               |
| Min/Max          | `slices.Min`, `slices.Max` (Go 1.21+)                | `MinFunc`, `MaxFunc`                                             |
| Equal            | `slices.Equal` (Go 1.21+)                            | `EqualFunc`                                                      |
| Zip              | Manual loop                                          | `Zip`, `ZipWith`                                                 |
| Pad              | Manual loop                                          | `Pad`, `PadLeft`                                                 |
| Remove           | `append(slice[:i], slice[i+1:]...)`                  | `RemoveAt`, `RemoveValue`, `RemoveFirst`                         |
| Insert           | `append(slice[:i], append([]T{v}, slice[i:]...)...)` | `Insert`                                                         |
| Associate        | Manual loop + map                                    | `Associate`, `AssociateWith`                                     |
| ForEach          | `for _, v := range slice`                            | `ForEach`, `ForEachWithIndex`                                    |

---

## Complete Example

```go
package main

import (
    "fmt"
    "github.com/alessiosavi/GoGPUtils/sliceutil"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Filter even numbers
    evens := sliceutil.Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println("Evens:", evens) // [2 4 6 8 10]

    // Double each number
    doubled := sliceutil.Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println("Doubled:", doubled) // [2 4 6 8 10 12 14 16 18 20]

    // Sum all numbers
    sum := sliceutil.Reduce(numbers, 0, func(acc, n int) int {
        return acc + n
    })
    fmt.Println("Sum:", sum) // 55

    // Split into chunks
    chunks := sliceutil.Chunk(numbers, 3)
    fmt.Println("Chunks:", chunks) // [[1 2 3] [4 5 6] [7 8 9] [10]]

    // Group by even/odd
    grouped := sliceutil.GroupBy(numbers, func(n int) string {
        if n%2 == 0 {
            return "even"
        }
        return "odd"
    })
    fmt.Println("Grouped:", grouped)
    // map[even:[2 4 6 8 10] odd:[1 3 5 7 9]]

    // Set operations
    a := []int{1, 2, 3, 4}
    b := []int{3, 4, 5, 6}
    fmt.Println("Union:", sliceutil.Union(a, b))           // [1 2 3 4 5 6]
    fmt.Println("Intersect:", sliceutil.Intersect(a, b))   // [3 4]
    fmt.Println("Difference:", sliceutil.Difference(a, b)) // [1 2]
}
```
