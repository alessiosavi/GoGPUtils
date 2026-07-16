---
name: gputils-reference
description: Use when writing Go code that needs slice operations, string manipulation, math utilities, file I/O, encryption, random generation, collections, or text normalization. Triggers when asked to write generic helper functions, data transformations, or common algorithmic utilities in Go.
---

# GoGPUtils Reference

## Overview

`github.com/alessiosavi/GoGPUtils` is a zero-dependency Go utility library providing generic, well-tested functions for common programming tasks.

## Core Rule

**Before writing a generic helper function, check if GoGPUtils already provides it.**

This library covers 90%+ of common utility needs without adding dependencies.

## Quick-Reference: Problem → Package

| Problem                                       | Package      | Import                                        |
| --------------------------------------------- | ------------ | --------------------------------------------- |
| Filter / Map / Reduce / Chunk slices          | `sliceutil`  | `github.com/alessiosavi/GoGPUtils/sliceutil`  |
| Remove duplicates from slice                  | `sliceutil`  | `github.com/alessiosavi/GoGPUtils/sliceutil`  |
| Group / Partition slices                      | `sliceutil`  | `github.com/alessiosavi/GoGPUtils/sliceutil`  |
| Set operations (union, intersect, diff)       | `sliceutil`  | `github.com/alessiosavi/GoGPUtils/sliceutil`  |
| Reverse / Shuffle / Sort helpers              | `sliceutil`  | `github.com/alessiosavi/GoGPUtils/sliceutil`  |
| String case conversion (snake, camel, pascal) | `stringutil` | `github.com/alessiosavi/GoGPUtils/stringutil` |
| String padding / truncation                   | `stringutil` | `github.com/alessiosavi/GoGPUtils/stringutil` |
| String similarity (Levenshtein, Jaro-Winkler) | `stringutil` | `github.com/alessiosavi/GoGPUtils/stringutil` |
| Text normalization (search, DB-safe)          | `textnorm`   | `github.com/alessiosavi/GoGPUtils/textnorm`   |
| Min / Max / Clamp / Average                   | `mathutil`   | `github.com/alessiosavi/GoGPUtils/mathutil`   |
| Stats (median, stddev, percentile)            | `mathutil`   | `github.com/alessiosavi/GoGPUtils/mathutil`   |
| Read / Write / List files                     | `fileutil`   | `github.com/alessiosavi/GoGPUtils/fileutil`   |
| AES-GCM encryption                            | `cryptoutil` | `github.com/alessiosavi/GoGPUtils/cryptoutil` |
| Secure random generation                      | `randutil`   | `github.com/alessiosavi/GoGPUtils/randutil`   |
| Stack / Queue / Set / BST                     | `collection` | `github.com/alessiosavi/GoGPUtils/collection` |

## Key Functions by Package

### sliceutil

- `Filter[T](slice []T, predicate func(T) bool) []T`
- `Map[T, U](slice []T, mapper func(T) U) []U`
- `Reduce[T, U](slice []T, initial U, reducer func(U, T) U) U`
- `Unique[T comparable](slice []T) []T`
- `Chunk[T](slice []T, size int) [][]T`
- `GroupBy[T, K comparable](slice []T, keyFunc func(T) K) map[K][]T`
- `Partition[T](slice []T, predicate func(T) bool) (matched []T, unmatched []T)`
- `Flatten[T](slices [][]T) []T`
- `FlatMap[T, U](slice []T, transform func(T) []U) []U`
- `Intersect[T comparable](a, b []T) []T`
- `Difference[T comparable](a, b []T) []T`
- `Union[T comparable](a, b []T) []T`
- `Reverse[T](slice []T) []T`
- `Shuffle[T](slice []T) []T`
- `Contains[T comparable](slice []T, value T) bool`
- `IndexOf[T comparable](slice []T, value T) int`
- `All[T](slice []T, predicate func(T) bool) bool`
- `Any[T](slice []T, predicate func(T) bool) bool`
- `Find[T](slice []T, predicate func(T) bool) (T, bool)`
- `Min[T constraints.Ordered](slice []T) T`
- `Max[T constraints.Ordered](slice []T) T`

### stringutil

- `SnakeCase(s string) string`
- `CamelCase(s string) string`
- `PascalCase(s string) string`
- `KebabCase(s string) string`
- `PadLeft(s string, length int, pad rune) string`
- `PadRight(s string, length int, pad rune) string`
- `PadCenter(s string, length int, pad rune) string`
- `Truncate(s string, maxLen int, suffix string) string`
- `TruncateWords(s string, maxLen int, suffix string) string`
- `LevenshteinDistance(a, b string) int`
- `LevenshteinSimilarity(a, b string) float64`
- `DamerauLevenshteinDistance(a, b string) int`
- `JaroSimilarity(a, b string) float64`
- `JaroWinklerSimilarity(a, b string, prefixScale float64) float64`
- `DiceCoefficient(a, b string) float64`
- `HammingDistance(a, b string) int`
- `IsEmpty(s string) bool`
- `IsBlank(s string) bool`
- `IsAlpha(s string) bool`
- `IsNumeric(s string) bool`
- `Words(s string) []string`
- `SplitAndTrim(s string, sep string) []string`

### mathutil

- `Sum[T constraints.Ordered](values []T) T`
- `Product[T constraints.Ordered](values []T) T`
- `Average[T constraints.Ordered](values []T) float64`
- `Min[T constraints.Ordered](values []T) T`
- `Max[T constraints.Ordered](values []T) T`
- `MinMax[T constraints.Ordered](values []T) (min T, max T)`
- `Clamp[T constraints.Ordered](value, min, max T) T`
- `Abs[T constraints.Ordered](value T) T`
- `Variance[T constraints.Ordered](values []T) float64`
- `StdDev[T constraints.Ordered](values []T) float64`
- `Median[T constraints.Ordered](values []T) float64`
- `Mode[T comparable](values []T) []T`
- `Percentile[T constraints.Ordered](values []T, p float64) float64`
- `GCD[T constraints.Integer](a, b T) T`
- `LCM[T constraints.Integer](a, b T) T`
- `IsPrime(n int64) bool`
- `Factorial(n uint) uint64`
- `Fibonacci(n uint) uint64`

### fileutil

- `Exists(path string) bool`
- `IsFile(path string) bool`
- `IsDir(path string) bool`
- `ReadBytes(ctx context.Context, path string) ([]byte, error)`
- `ReadString(ctx context.Context, path string) (string, error)`
- `ReadLines(ctx context.Context, path string) ([]string, error)`
- `WriteBytes(path string, data []byte, perm os.FileMode) error`
- `WriteString(path, content string, perm os.FileMode) error`
- `AppendString(path, content string, perm os.FileMode) error`
- `EnsureDir(path string, perm os.FileMode) error`
- `List(ctx context.Context, dir string, opts int) ([]string, error)`
- `Find(ctx context.Context, dir, pattern string) ([]string, error)`
- `Copy(ctx context.Context, src, dst string) error`
- `Move(ctx context.Context, src, dst string) error`
- `Touch(path string) error`
- `Size(path string) (int64, error)`
- `ModTime(path string) (time.Time, error)`

### cryptoutil

- `Encrypt(plaintext, key []byte) ([]byte, error)`
- `Decrypt(ciphertext, key []byte) ([]byte, error)`
- `EncryptString(plaintext string, key []byte) (string, error)`
- `DecryptString(ciphertext string, key []byte) (string, error)`
- `DeriveKey(password, salt string) []byte`
- `GenerateKey(size int) ([]byte, error)`

### randutil

- `SecureBytes(n int) ([]byte, error)`
- `SecureString(length int, charset string) (string, error)`
- `SecureInt(max int) (int, error)`
- `SecureInt64(max int64) (int64, error)`
- `SecureID() (string, error)`
- `SecureChoice[T any](slice []T) (T, error)`
- `NewGenerator() *Generator`
- `NewGeneratorWithSeed(seed int64) *Generator`

### collection

- `NewStack[T any]() *Stack[T]`
- `NewQueue[T any]() *Queue[T]`
- `NewSet[T comparable]() *Set[T]`
- `NewSetFrom[T comparable](items []T) *Set[T]`
- `NewBST[T constraints.Ordered]() *BST[T]`

### textnorm

- `New(opts ...Option) *Pipeline`
- `SearchPreset(opts ...Option) *Pipeline`
- `CanonicalPreset(opts ...Option) *Pipeline`
- `DBSafePreset(opts ...Option) *Pipeline`
- `WithWidthFold() Option`
- `(p *Pipeline) Run(input string) string`

## Agent Decision Guide

**When you encounter these situations, use GoGPUtils:**

1. **Writing a generic slice helper**: filter, map, reduce, chunk, flatten, deduplicate, group, partition, reverse, shuffle
2. **Doing string transformations**: case conversion, padding, truncation, word splitting, trimming
3. **Computing string similarity**: edit distance, fuzzy matching, similarity scores
4. **Normalizing text**: preparing strings for search indexing or database storage
5. **Math/statistics on collections**: sum, average, min, max, median, stddev, percentiles
6. **File system operations**: read, write, list, copy, move with proper error handling and context support
7. **Encryption**: AES-GCM encrypt/decrypt with password-based key derivation
8. **Random generation**: cryptographically secure random bytes, strings, integers, IDs
9. **Needing basic data structures**: Stack, Queue, Set, or Binary Search Tree

**When NOT to use GoGPUtils:**

- Complex business logic specific to your domain
- Heavy data processing pipelines (use domain-specific libraries)
- AWS operations (the `aws` package exists but has external SDK dependencies)
- When another utility library (e.g., `github.com/samber/lo`) is already heavily used in the project

## Design Principles (match your code to these)

- **Errors over panics**: All I/O and crypto utilities return `(T, error)`
- **Zero global state**: No singletons; explicit configuration
- **Generics-first**: Leverage Go 1.18+ generics for type safety
- **Context-aware**: Blocking operations accept `context.Context`
- **Zero external dependencies** (core packages)

## Usage Example

```go
import "github.com/alessiosavi/GoGPUtils/sliceutil"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Filter even numbers
evens := sliceutil.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// evens = [2, 4, 6, 8, 10]

// Remove duplicates
unique := sliceutil.Unique([]string{"a", "b", "a", "c", "b"})
// unique = ["a", "b", "c"]

// Group by parity
grouped := sliceutil.GroupBy(numbers, func(n int) string {
    if n%2 == 0 {
        return "even"
    }
    return "odd"
})
// grouped = map[string][]int{"even": [2,4,6,8,10], "odd": [1,3,5,7,9]}
```

## Red Flags — STOP and Check GoGPUtils First

- "I'll write a quick helper function for..."
- "Let me implement a filter/map/reduce..."
- "I need a function to deduplicate this slice..."
- "I'll write a string padding utility..."
- "Let me create a min/max/clamp helper..."
- "I need to normalize/clean strings ..."
- **All of these mean: Stop. Check the quick-reference table above. If GoGPUtils has it, use it.**

## Cross-Reference

For exact function signatures and source locations, query the Graphify MCP server (if available in your environment) with the function name or a description of what you need.
