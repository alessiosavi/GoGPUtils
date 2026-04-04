// Package stringutil provides string manipulation utilities.
//
// The package is organized into several categories:
//
// # Search and Indexing
//
// Functions for finding substrings and patterns:
//
//	indices := stringutil.AllIndexes("banana", "an")  // [1, 3]
//	ok := stringutil.HasAnyPrefix(s, "http://", "https://")
//	ok := stringutil.ContainsAll(s, "foo", "bar")
//
// # Transformation
//
// Functions for transforming strings:
//
//	reversed := stringutil.Reverse("hello")  // "olleh"
//	truncated := stringutil.Truncate(s, 100, "...")
//	padded := stringutil.PadLeft("42", 5, '0')  // "00042"
//
// # Validation
//
// Functions for checking string properties:
//
//	if stringutil.IsNumeric(s) { ... }
//	if stringutil.IsAlpha(s) { ... }
//	if stringutil.IsPalindrome(s) { ... }
//
// # Similarity (see similarity.go)
//
// Algorithms for measuring string similarity:
//
//	distance := stringutil.LevenshteinDistance("kitten", "sitting")  // 3
//	score := stringutil.JaroWinklerSimilarity("martha", "marhta", 0.1) // ~0.96
//	coefficient := stringutil.DiceCoefficient("night", "nacht")
//
// All functions are designed to be nil-safe and handle edge cases gracefully.
package stringutil
