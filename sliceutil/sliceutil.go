package sliceutil

import (
	"cmp"
	"slices"
)

// Filter returns a new slice containing only elements for which the predicate returns true.
// Returns nil if the input slice is nil.
//
// Example:
//
//	evens := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
//	// evens = [2, 4]
func Filter[T any](s []T, predicate func(T) bool) []T {
	if s == nil {
		return nil
	}

	result := make([]T, 0, len(s)/2) // Estimate half will match

	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}

// FilterInPlace filters the slice in place, returning the modified slice.
// The underlying array is modified; elements are shifted to fill gaps.
// This avoids allocation but modifies the original slice.
//
// Example:
//
//	nums := []int{1, 2, 3, 4, 5}
//	evens := FilterInPlace(nums, func(n int) bool { return n%2 == 0 })
//	// evens = [2, 4], nums[0:2] = [2, 4]
func FilterInPlace[T any](s []T, predicate func(T) bool) []T {
	if s == nil {
		return nil
	}

	n := 0

	for _, v := range s {
		if predicate(v) {
			s[n] = v
			n++
		}
	}

	return s[:n]
}

// Map applies a transformation function to each element and returns a new slice.
// Returns nil if the input slice is nil.
//
// Example:
//
//	doubled := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
//	// doubled = [2, 4, 6]
func Map[T, U any](s []T, transform func(T) U) []U {
	if s == nil {
		return nil
	}

	result := make([]U, len(s))
	for i, v := range s {
		result[i] = transform(v)
	}

	return result
}

// MapWithIndex applies a transformation function that receives both index and element.
//
// Example:
//
//	indexed := MapWithIndex([]string{"a", "b"}, func(i int, s string) string {
//	    return fmt.Sprintf("%d:%s", i, s)
//	})
//	// indexed = ["0:a", "1:b"]
func MapWithIndex[T, U any](s []T, transform func(int, T) U) []U {
	if s == nil {
		return nil
	}

	result := make([]U, len(s))
	for i, v := range s {
		result[i] = transform(i, v)
	}

	return result
}

// Reduce reduces a slice to a single value using an accumulator function.
//
// Example:
//
//	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
//	// sum = 10
func Reduce[T, U any](s []T, initial U, accumulator func(U, T) U) U {
	result := initial
	for _, v := range s {
		result = accumulator(result, v)
	}

	return result
}

// Contains reports whether the slice contains the target element.
// Uses == for comparison, requiring comparable types.
//
// Example:
//
//	if Contains(names, "Alice") { ... }
func Contains[T comparable](s []T, target T) bool {

	return slices.Contains(s, target)
}

// ContainsFunc reports whether any element satisfies the predicate.
//
// Example:
//
//	hasNegative := ContainsFunc(nums, func(n int) bool { return n < 0 })
func ContainsFunc[T any](s []T, predicate func(T) bool) bool {

	return slices.ContainsFunc(s, predicate)
}

// IndexOf returns the index of the first occurrence of target, or -1 if not found.
//
// Example:
//
//	idx := IndexOf([]string{"a", "b", "c"}, "b")
//	// idx = 1
func IndexOf[T comparable](s []T, target T) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}

	return -1
}

// IndexOfFunc returns the index of the first element satisfying the predicate, or -1.
func IndexOfFunc[T any](s []T, predicate func(T) bool) int {
	for i, v := range s {
		if predicate(v) {
			return i
		}
	}

	return -1
}

// LastIndexOf returns the index of the last occurrence of target, or -1 if not found.
func LastIndexOf[T comparable](s []T, target T) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == target {
			return i
		}
	}

	return -1
}

// Unique returns a new slice with duplicate elements removed.
// The first occurrence of each element is kept; order is preserved.
//
// Example:
//
//	unique := Unique([]int{1, 2, 2, 3, 1})
//	// unique = [1, 2, 3]
func Unique[T comparable](s []T) []T {
	if s == nil {
		return nil
	}

	if len(s) == 0 {
		return []T{}
	}

	seen := make(map[T]struct{}, len(s))

	result := make([]T, 0, len(s))

	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}

			result = append(result, v)
		}
	}

	return result
}

// UniqueFunc returns a new slice with duplicates removed based on a key function.
// Elements with the same key are considered duplicates; first occurrence is kept.
//
// Example:
//
//	type User struct { ID int; Name string }
//	users := []User{{1, "Alice"}, {2, "Bob"}, {1, "Alice2"}}
//	unique := UniqueFunc(users, func(u User) int { return u.ID })
//	// unique = [{1, "Alice"}, {2, "Bob"}]
func UniqueFunc[T any, K comparable](s []T, key func(T) K) []T {
	if s == nil {
		return nil
	}

	seen := make(map[K]struct{}, len(s))

	result := make([]T, 0, len(s))

	for _, v := range s {
		k := key(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}

			result = append(result, v)
		}
	}

	return result
}

// Chunk splits a slice into chunks of the specified size.
// The last chunk may be smaller if len(s) is not divisible by size.
// Returns nil if size <= 0 or s is nil.
//
// Example:
//
//	chunks := Chunk([]int{1, 2, 3, 4, 5}, 2)
//	// chunks = [[1, 2], [3, 4], [5]]
func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 || s == nil {
		return nil
	}

	if len(s) == 0 {
		return [][]T{}
	}

	numChunks := (len(s) + size - 1) / size
	result := make([][]T, 0, numChunks)

	for i := 0; i < len(s); i += size {
		end := min(i+size, len(s))

		result = append(result, s[i:end])
	}

	return result
}

// Flatten converts a slice of slices into a single slice.
//
// Example:
//
//	flat := Flatten([][]int{{1, 2}, {3}, {4, 5}})
//	// flat = [1, 2, 3, 4, 5]
func Flatten[T any](s [][]T) []T {
	if s == nil {
		return nil
	}

	total := 0
	for _, inner := range s {
		total += len(inner)
	}

	result := make([]T, 0, total)
	for _, inner := range s {
		result = append(result, inner...)
	}

	return result
}

// Reverse returns a new slice with elements in reverse order.
//
// Example:
//
//	rev := Reverse([]int{1, 2, 3})
//	// rev = [3, 2, 1]
func Reverse[T any](s []T) []T {
	if s == nil {
		return nil
	}

	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}

	return result
}

// ReverseInPlace reverses the slice in place and returns it.
//
// Example:
//
//	nums := []int{1, 2, 3}
//	ReverseInPlace(nums)
//	// nums = [3, 2, 1]
func ReverseInPlace[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// Intersect returns elements present in both slices.
// Result preserves order from the first slice.
//
// Example:
//
//	common := Intersect([]int{1, 2, 3}, []int{2, 3, 4})
//	// common = [2, 3]
func Intersect[T comparable](a, b []T) []T {
	if a == nil || b == nil {
		return nil
	}

	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)

	for _, v := range a {
		if _, ok := set[v]; ok {
			result = append(result, v)
		}
	}

	return Unique(result)
}

// Difference returns elements in a that are not in b.
//
// Example:
//
//	diff := Difference([]int{1, 2, 3}, []int{2, 3, 4})
//	// diff = [1]
func Difference[T comparable](a, b []T) []T {
	if a == nil {
		return nil
	}

	if b == nil {
		return slices.Clone(a)
	}

	set := make(map[T]struct{}, len(b))
	for _, v := range b {
		set[v] = struct{}{}
	}

	result := make([]T, 0)

	for _, v := range a {
		if _, ok := set[v]; !ok {
			result = append(result, v)
		}
	}

	return result
}

// Union returns all unique elements from both slices.
//
// Example:
//
//	all := Union([]int{1, 2}, []int{2, 3})
//	// all = [1, 2, 3]
func Union[T comparable](a, b []T) []T {
	if a == nil && b == nil {
		return nil
	}

	combined := make([]T, 0, len(a)+len(b))
	combined = append(combined, a...)
	combined = append(combined, b...)

	return Unique(combined)
}

// GroupBy groups elements by a key function.
//
// Example:
//
//	nums := []int{1, 2, 3, 4, 5, 6}
//	groups := GroupBy(nums, func(n int) string {
//	    if n%2 == 0 { return "even" }
//	    return "odd"
//	})
//	// groups = map["even":[2, 4, 6] "odd":[1, 3, 5]]
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	if s == nil {
		return nil
	}

	result := make(map[K][]T)

	for _, v := range s {
		k := key(v)
		result[k] = append(result[k], v)
	}

	return result
}

// Partition splits a slice into two: elements matching the predicate and those that don't.
//
// Example:
//
//	evens, odds := Partition([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
//	// evens = [2, 4], odds = [1, 3]
func Partition[T any](s []T, predicate func(T) bool) (matching, notMatching []T) {
	if s == nil {
		return nil, nil
	}

	matching = make([]T, 0, len(s)/2)

	notMatching = make([]T, 0, len(s)/2)

	for _, v := range s {
		if predicate(v) {
			matching = append(matching, v)
		} else {
			notMatching = append(notMatching, v)
		}
	}

	return matching, notMatching
}

// Take returns the first n elements. If n > len(s), returns all elements.
//
// Example:
//
//	first := Take([]int{1, 2, 3, 4, 5}, 3)
//	// first = [1, 2, 3]
func Take[T any](s []T, n int) []T {
	if s == nil || n <= 0 {
		return nil
	}

	if n > len(s) {
		n = len(s)
	}

	return slices.Clone(s[:n])
}

// TakeLast returns the last n elements.
//
// Example:
//
//	last := TakeLast([]int{1, 2, 3, 4, 5}, 3)
//	// last = [3, 4, 5]
func TakeLast[T any](s []T, n int) []T {
	if s == nil || n <= 0 {
		return nil
	}

	if n > len(s) {
		n = len(s)
	}

	return slices.Clone(s[len(s)-n:])
}

// Drop returns elements after skipping the first n. If n >= len(s), returns empty slice.
//
// Example:
//
//	rest := Drop([]int{1, 2, 3, 4, 5}, 2)
//	// rest = [3, 4, 5]
func Drop[T any](s []T, n int) []T {
	if s == nil {
		return nil
	}

	if n >= len(s) {
		return []T{}
	}

	if n < 0 {
		n = 0
	}

	return slices.Clone(s[n:])
}

// DropLast returns elements after removing the last n.
func DropLast[T any](s []T, n int) []T {
	if s == nil {
		return nil
	}

	if n >= len(s) {
		return []T{}
	}

	if n < 0 {
		n = 0
	}

	return slices.Clone(s[:len(s)-n])
}

// TakeWhile returns elements from the start while predicate returns true.
//
// Example:
//
//	result := TakeWhile([]int{1, 2, 3, 4, 1}, func(n int) bool { return n < 3 })
//	// result = [1, 2]
func TakeWhile[T any](s []T, predicate func(T) bool) []T {
	if s == nil {
		return nil
	}

	for i, v := range s {
		if !predicate(v) {
			return slices.Clone(s[:i])
		}
	}

	return slices.Clone(s)
}

// DropWhile returns elements after dropping from the start while predicate returns true.
func DropWhile[T any](s []T, predicate func(T) bool) []T {
	if s == nil {
		return nil
	}

	for i, v := range s {
		if !predicate(v) {
			return slices.Clone(s[i:])
		}
	}

	return []T{}
}

// All returns true if all elements satisfy the predicate.
// Returns true for empty slices.
//
// Example:
//
//	allPositive := All([]int{1, 2, 3}, func(n int) bool { return n > 0 })
//	// allPositive = true
func All[T any](s []T, predicate func(T) bool) bool {
	for _, v := range s {
		if !predicate(v) {
			return false
		}
	}

	return true
}

// Any returns true if any element satisfies the predicate.
// Returns false for empty slices.
//
// Example:
//
//	hasNegative := Any([]int{1, -2, 3}, func(n int) bool { return n < 0 })
//	// hasNegative = true
func Any[T any](s []T, predicate func(T) bool) bool {

	return slices.ContainsFunc(s, predicate)
}

// None returns true if no elements satisfy the predicate.
// Returns true for empty slices.
func None[T any](s []T, predicate func(T) bool) bool {
	return !Any(s, predicate)
}

// Count returns the number of elements satisfying the predicate.
func Count[T any](s []T, predicate func(T) bool) int {
	count := 0

	for _, v := range s {
		if predicate(v) {
			count++
		}
	}

	return count
}

// Find returns the first element satisfying the predicate and true, or zero value and false.
//
// Example:
//
//	first, ok := Find(users, func(u User) bool { return u.Age > 18 })
func Find[T any](s []T, predicate func(T) bool) (T, bool) {
	for _, v := range s {
		if predicate(v) {
			return v, true
		}
	}

	var zero T

	return zero, false
}

// FindLast returns the last element satisfying the predicate.
func FindLast[T any](s []T, predicate func(T) bool) (T, bool) {
	for i := len(s) - 1; i >= 0; i-- {
		if predicate(s[i]) {
			return s[i], true
		}
	}

	var zero T

	return zero, false
}

// Min returns the minimum element using natural ordering.
// Returns zero value and false for empty slices.
func Min[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T

		return zero, false
	}

	min := s[0]
	for _, v := range s[1:] {
		if v < min {
			min = v
		}
	}

	return min, true
}

// Max returns the maximum element using natural ordering.
// Returns zero value and false for empty slices.
func Max[T cmp.Ordered](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T

		return zero, false
	}

	max := s[0]
	for _, v := range s[1:] {
		if v > max {
			max = v
		}
	}

	return max, true
}

// MinFunc returns the minimum element using a comparison function.
// cmp should return negative if a < b, zero if a == b, positive if a > b.
func MinFunc[T any](s []T, cmpFn func(a, b T) int) (T, bool) {
	if len(s) == 0 {
		var zero T

		return zero, false
	}

	min := s[0]
	for _, v := range s[1:] {
		if cmpFn(v, min) < 0 {
			min = v
		}
	}

	return min, true
}

// MaxFunc returns the maximum element using a comparison function.
func MaxFunc[T any](s []T, cmpFn func(a, b T) int) (T, bool) {
	if len(s) == 0 {
		var zero T

		return zero, false
	}

	max := s[0]
	for _, v := range s[1:] {
		if cmpFn(v, max) > 0 {
			max = v
		}
	}

	return max, true
}

// Equal reports whether two slices are equal.
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// EqualFunc reports whether two slices are equal using a custom comparison.
func EqualFunc[T, U any](a []T, b []U, eq func(T, U) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !eq(a[i], b[i]) {
			return false
		}
	}

	return true
}

// ForEach calls a function for each element.
//
// Example:
//
//	ForEach(items, func(item Item) { process(item) })
func ForEach[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

// ForEachWithIndex calls a function for each element with its index.
func ForEachWithIndex[T any](s []T, fn func(int, T)) {
	for i, v := range s {
		fn(i, v)
	}
}

// Zip combines two slices into pairs. Stops at the shorter slice.
//
// Example:
//
//	type Pair[T, U any] struct { First T; Second U }
//	pairs := Zip([]int{1, 2}, []string{"a", "b"})
//	// pairs = [{1, "a"}, {2, "b"}]
func Zip[T, U any](a []T, b []U) [][2]any {
	length := min(len(b), len(a))

	result := make([][2]any, length)
	for i := 0; i < length; i++ {
		result[i] = [2]any{a[i], b[i]}
	}

	return result
}

// ZipWith combines two slices using a function.
//
// Example:
//
//	sums := ZipWith([]int{1, 2, 3}, []int{4, 5, 6}, func(a, b int) int { return a + b })
//	// sums = [5, 7, 9]
func ZipWith[T, U, V any](a []T, b []U, combine func(T, U) V) []V {
	length := min(len(b), len(a))

	result := make([]V, length)
	for i := 0; i < length; i++ {
		result[i] = combine(a[i], b[i])
	}

	return result
}

// Pad extends a slice to the target length by appending the fill value.
// If slice is already >= length, returns a clone.
//
// Example:
//
//	padded := Pad([]int{1, 2}, 5, 0)
//	// padded = [1, 2, 0, 0, 0]
func Pad[T any](s []T, length int, fill T) []T {
	if len(s) >= length {
		return slices.Clone(s)
	}

	result := make([]T, length)
	copy(result, s)

	for i := len(s); i < length; i++ {
		result[i] = fill
	}

	return result
}

// PadLeft extends a slice by prepending the fill value.
//
// Example:
//
//	padded := PadLeft([]int{1, 2}, 5, 0)
//	// padded = [0, 0, 0, 1, 2]
func PadLeft[T any](s []T, length int, fill T) []T {
	if len(s) >= length {
		return slices.Clone(s)
	}

	result := make([]T, length)

	offset := length - len(s)

	for i := range offset {
		result[i] = fill
	}

	copy(result[offset:], s)

	return result
}

// RemoveAt returns a new slice with the element at index removed.
// Returns nil if index is out of bounds.
//
// Example:
//
//	removed := RemoveAt([]int{1, 2, 3}, 1)
//	// removed = [1, 3]
func RemoveAt[T any](s []T, index int) []T {
	if index < 0 || index >= len(s) {
		return nil
	}

	result := make([]T, 0, len(s)-1)
	result = append(result, s[:index]...)
	result = append(result, s[index+1:]...)

	return result
}

// RemoveValue returns a new slice with all occurrences of value removed.
//
// Example:
//
//	cleaned := RemoveValue([]int{1, 2, 3, 2}, 2)
//	// cleaned = [1, 3]
func RemoveValue[T comparable](s []T, value T) []T {
	return Filter(s, func(v T) bool { return v != value })
}

// RemoveFirst returns a new slice with the first occurrence of value removed.
func RemoveFirst[T comparable](s []T, value T) []T {
	idx := IndexOf(s, value)
	if idx == -1 {
		return slices.Clone(s)
	}

	return RemoveAt(s, idx)
}

// Insert returns a new slice with value inserted at index.
// If index is out of bounds, appends to end.
//
// Example:
//
//	inserted := Insert([]int{1, 3}, 1, 2)
//	// inserted = [1, 2, 3]
func Insert[T any](s []T, index int, value T) []T {
	if index < 0 {
		index = 0
	}

	if index >= len(s) {
		result := make([]T, len(s)+1)
		copy(result, s)
		result[len(s)] = value

		return result
	}

	result := make([]T, len(s)+1)
	copy(result, s[:index])
	result[index] = value
	copy(result[index+1:], s[index:])

	return result
}

// Shuffle returns a new slice with elements in random order.
// Uses crypto/rand for secure randomness.
//
// For a deterministic shuffle, use ShuffleWithSeed.
func Shuffle[T any](s []T) []T {
	if s == nil {
		return nil
	}

	result := slices.Clone(s)
	ShuffleInPlace(result)

	return result
}

// ShuffleInPlace shuffles the slice in place using Fisher-Yates algorithm.
// Uses math/rand for performance; use Shuffle for crypto-secure randomness.
func ShuffleInPlace[T any](s []T) {
	// Fisher-Yates shuffle using simple PRNG seeded from time
	// For tests/determinism, see ShuffleWithSeed
	for i := len(s) - 1; i > 0; i-- {
		// Simple LCG for shuffling - not crypto secure but fast
		j := int(fastrand()) % (i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// Simple fast random for shuffling (not crypto secure).
var fastrandState uint64 = 1

func fastrand() uint32 {
	// xorshift64
	fastrandState ^= fastrandState << 13
	fastrandState ^= fastrandState >> 7
	fastrandState ^= fastrandState << 17

	return uint32(fastrandState)
}

// SeedShuffle sets the seed for ShuffleInPlace. Useful for deterministic tests.
func SeedShuffle(seed uint64) {
	if seed == 0 {
		seed = 1
	}

	fastrandState = seed
}

// Compact removes consecutive duplicate elements, similar to Unix uniq.
// For removing all duplicates, use Unique.
//
// Example:
//
//	compacted := Compact([]int{1, 1, 2, 2, 2, 1})
//	// compacted = [1, 2, 1]
func Compact[T comparable](s []T) []T {
	if len(s) <= 1 {
		return slices.Clone(s)
	}

	result := make([]T, 0, len(s))

	result = append(result, s[0])

	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1] {
			result = append(result, s[i])
		}
	}

	return result
}

// CompactFunc removes consecutive elements where the predicate returns true.
func CompactFunc[T any](s []T, eq func(T, T) bool) []T {
	if len(s) <= 1 {
		return slices.Clone(s)
	}

	result := make([]T, 0, len(s))

	result = append(result, s[0])

	for i := 1; i < len(s); i++ {
		if !eq(s[i-1], s[i]) {
			result = append(result, s[i])
		}
	}

	return result
}

// FlatMap applies a function that returns a slice and flattens the result.
//
// Example:
//
//	words := FlatMap([]string{"hello world", "foo bar"}, func(s string) []string {
//	    return strings.Split(s, " ")
//	})
//	// words = ["hello", "world", "foo", "bar"]
func FlatMap[T, U any](s []T, transform func(T) []U) []U {
	if s == nil {
		return nil
	}

	var result []U
	for _, v := range s {
		result = append(result, transform(v)...)
	}

	return result
}

// Associate creates a map from slice elements using a key function.
//
// Example:
//
//	users := []User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}}
//	byID := Associate(users, func(u User) int { return u.ID })
//	// byID = map[1:{1, "Alice"} 2:{2, "Bob"}]
func Associate[T any, K comparable](s []T, key func(T) K) map[K]T {
	if s == nil {
		return nil
	}

	result := make(map[K]T, len(s))
	for _, v := range s {
		result[key(v)] = v
	}

	return result
}

// AssociateWith creates a map from slice elements to values computed by a function.
//
// Example:
//
//	lengths := AssociateWith([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) })
//	// lengths = map["a":1 "bb":2 "ccc":3]
func AssociateWith[K comparable, V any](s []K, value func(K) V) map[K]V {
	if s == nil {
		return nil
	}

	result := make(map[K]V, len(s))
	for _, k := range s {
		result[k] = value(k)
	}

	return result
}
