---
name: gputils-reference
description: Use when writing Go code that needs slice operations, string manipulation, math utilities, file I/O, encryption, random generation, collections, in-memory caching, or text normalization. Triggers when asked to write generic helper functions, data transformations, or common algorithmic utilities in Go.
---

# GoGPUtils Reference

## Overview

`github.com/alessiosavi/GoGPUtils` is a Go utility library providing generic, well-tested functions for common programming tasks. Most packages use only the standard library; `stringutil` and `textnorm` also use `golang.org/x/text`, and the `aws` packages use the AWS SDK v2. Requires Go 1.27.1 or later.

## Core Rule

**Before writing a generic helper function, check if GoGPUtils already provides it.**

This library covers many common utility needs with minimal dependencies.

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
| In-memory cache (TTL, SIEVE, GetOrLoad)       | `cache`      | `github.com/alessiosavi/GoGPUtils/cache`      |

## Key Functions by Package

### sliceutil

- `Filter[T any](s []T, predicate func(T) bool) []T`
- `Map[T, U any](s []T, transform func(T) U) []U`
- `Reduce[T, U any](s []T, initial U, accumulator func(U, T) U) U`
- `Unique[T comparable](s []T) []T`
- `Chunk[T any](s []T, size int) [][]T`
- `GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T`
- `Partition[T any](s []T, predicate func(T) bool) (matching, notMatching []T)`
- `Flatten[T any](s [][]T) []T`
- `FlatMap[T, U any](s []T, transform func(T) []U) []U`
- `Intersect[T comparable](a, b []T) []T`
- `Difference[T comparable](a, b []T) []T`
- `Union[T comparable](a, b []T) []T`
- `Reverse[T any](s []T) []T`
- `Shuffle[T any](s []T) []T`
- `Contains[T comparable](s []T, target T) bool`
- `IndexOf[T comparable](s []T, target T) int`
- `All[T any](s []T, predicate func(T) bool) bool`
- `Any[T any](s []T, predicate func(T) bool) bool`
- `Find[T any](s []T, predicate func(T) bool) (T, bool)`
- `Min[T cmp.Ordered](s []T) (T, bool)`
- `Max[T cmp.Ordered](s []T) (T, bool)`

### stringutil

- `SnakeCase(s string) string`
- `CamelCase(s string) string`
- `PascalCase(s string) string`
- `KebabCase(s string) string`
- `PadLeft(s string, length int, padChar rune) string`
- `PadRight(s string, length int, padChar rune) string`
- `PadCenter(s string, length int, padChar rune) string`
- `Truncate(s string, maxLen int, suffix string) string`
- `TruncateWords(s string, maxLen int, suffix string) string`
- `LevenshteinDistance(s1, s2 string) int`
- `LevenshteinSimilarity(s1, s2 string) float64`
- `DamerauLevenshteinDistance(s1, s2 string) int`
- `JaroSimilarity(s1, s2 string) float64`
- `JaroWinklerSimilarity(s1, s2 string, prefixScale float64) float64`
- `DiceCoefficient(s1, s2 string) float64`
- `HammingDistance(s1, s2 string) int`
- `IsEmpty(s string) bool`
- `IsBlank(s string) bool`
- `IsAlpha(s string) bool`
- `IsNumeric(s string) bool`
- `Words(s string) []string`
- `SplitAndTrim(s, sep string) []string`

### mathutil

- `Sum[T Number](s []T) T`
- `Product[T Number](s []T) T`
- `Average[T Number](s []T) float64`
- `Min[T Number](s []T) (T, error)`
- `Max[T Number](s []T) (T, error)`
- `MinMax[T Number](s []T) (min, max T, err error)`
- `Clamp[T Number](x, min, max T) T`
- `Abs[T Number](x T) T`
- `Variance[T Number](s []T) float64`
- `StdDev[T Number](s []T) float64`
- `Median[T Number](s []T) float64`
- `Mode[T Number](s []T) []T`
- `Percentile[T Number](s []T, p float64) float64`
- `GCD[T Integer](a, b T) T`
- `LCM[T Integer](a, b T) T`
- `IsPrime[T Integer](n T) bool`
- `Factorial(n int) int64`
- `Fibonacci(n int) int64`

### fileutil

- `Exists(path string) bool`
- `IsFile(path string) bool`
- `IsDir(path string) bool`
- `ReadBytes(ctx context.Context, path string) ([]byte, error)`
- `ReadString(ctx context.Context, path string) (string, error)`
- `ReadLines(ctx context.Context, path string) ([]string, error)`
- `WriteBytes(path string, data []byte, perm fs.FileMode) error`
- `WriteString(path, content string, perm fs.FileMode) error`
- `AppendString(path, content string, perm fs.FileMode) error`
- `EnsureDir(path string, perm fs.FileMode) error`
- `List(ctx context.Context, dir string, opts ListOption) ([]string, error)`
- `Find(ctx context.Context, dir, pattern string) ([]string, error)`
- `Copy(ctx context.Context, src, dst string) error`
- `Move(ctx context.Context, src, dst string) error`
- `Touch(path string) error`
- `Size(path string) (int64, error)`
- `ModTime(path string) (time.Time, error)`

### cryptoutil

- `Encrypt(plaintext, key []byte) (string, error)`
- `Decrypt(ciphertextB64 string, key []byte) ([]byte, error)`
- `EncryptString(plaintext string, key []byte) (string, error)`
- `DecryptString(ciphertextB64 string, key []byte) (string, error)`
- `DeriveKey(password, salt string) []byte`
- `GenerateKey(size int) ([]byte, error)`

### randutil

- `SecureBytes(n int) ([]byte, error)`
- `SecureString(length int, charset string) (string, error)`
- `SecureInt(max int) (int, error)`
- `SecureInt64(max int64) (int64, error)`
- `SecureID() (string, error)`
- `SecureChoice[T any](s []T) (T, error)`
- `NewGenerator() *Generator`
- `NewGeneratorWithSeed(seed uint64) *Generator`

### collection

- `NewStack[T any]() *Stack[T]`
- `NewQueue[T any]() *Queue[T]`
- `NewSet[T comparable]() *Set[T]`
- `NewSetFrom[T comparable](items []T) *Set[T]`
- `NewBST[T cmp.Ordered]() *BST[T]`

### textnorm

- `New() Pipeline`
- `SearchPreset(opts ...PresetOption) Pipeline`
- `CanonicalPreset(opts ...PresetOption) Pipeline`
- `MeaningPreset(opts ...PresetOption) Pipeline`
- `DBSafePreset(opts ...PresetOption) Pipeline`
- `HygienePreset(opts ...PresetOption) Pipeline`
- `WithWidthFold() PresetOption`
- `(p Pipeline) Run(input string) (string, error)`

### cache

- `New[K comparable, V any](cfg Config[K, V]) (*Cache[K, V], error)`
- `(c *Cache[K, V]) Get(key K) (V, bool)`
- `(c *Cache[K, V]) Set(key K, value V)`
- `(c *Cache[K, V]) SetWithTTL(key K, value V, ttl time.Duration)`
- `(c *Cache[K, V]) GetOrLoad(ctx context.Context, key K, load func(context.Context) (V, error)) (V, error)`
- `(c *Cache[K, V]) Delete(key K) bool`
- `(c *Cache[K, V]) Stats() Stats`
- `(c *Cache[K, V]) Close() error`

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
10. **Caching values in memory**: bounded cache with TTL, SIEVE eviction, and deduplicated loading

**When NOT to use GoGPUtils:**

- Complex business logic specific to your domain
- Heavy data processing pipelines (use domain-specific libraries)
- AWS operations (the `aws` package exists but has external SDK dependencies)
- When another utility library (e.g., `github.com/samber/lo`) is already heavily used in the project

## Design Principles (match your code to these)

- **Errors over panics**: All I/O and crypto utilities return `(T, error)`
- **Zero global state**: No singletons; explicit configuration
- **Generics-first**: Leverage generics for type safety (the module requires Go 1.27.1+)
- **Context-aware**: Blocking operations accept `context.Context`
- **Minimal dependencies**: core packages use only the standard library; `stringutil` and `textnorm` add `golang.org/x/text`

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
