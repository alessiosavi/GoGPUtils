// Package sliceutil provides generic slice operations for Go.
//
// All functions in this package work with any comparable or ordered type
// through Go's generics. The package avoids allocations where possible
// and provides both in-place and copying variants of mutating operations.
//
// # Design Principles
//
//   - Generic: Works with any type that satisfies the constraints
//   - Non-mutating by default: Functions return new slices unless suffixed with "InPlace"
//   - Nil-safe: All functions handle nil slices gracefully
//   - Zero allocations: Where possible, operations avoid heap allocations
//
// # Basic Operations
//
//	// Filter elements
//	evens := sliceutil.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
//
//	// Transform elements
//	doubled := sliceutil.Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
//
//	// Check containment
//	if sliceutil.Contains(names, "Alice") { ... }
//
// # Chunking and Grouping
//
//	// Split into chunks of size 3
//	chunks := sliceutil.Chunk(items, 3)
//
//	// Group items by key
//	groups := sliceutil.GroupBy([]string{"a", "bb"}, func(item string) int { return len(item) })
//
// # Set Operations
//
//	unique := sliceutil.Unique(items)
//	intersection := sliceutil.Intersect(a, b)
//	difference := sliceutil.Difference(a, b)
package sliceutil
