// Package gputils provides a collection of well-tested, idiomatic Go utilities
// for common programming tasks.
//
// # Design Philosophy
//
// This library follows several key principles:
//
//   - Errors over panics: All functions return errors instead of panicking
//   - Zero global state: No singletons; all state is explicit
//   - Generic when useful: Uses generics to reduce duplication without over-abstraction
//   - Minimal dependencies: Utilities use the standard library, with golang.org/x/text
//     for Unicode text processing in stringutil and textnorm; AWS helpers use the AWS SDK v2
//   - Context-aware: Blocking operations accept context.Context
//
// # Packages
//
// The library is organized into focused packages:
//
//   - sliceutil: Generic slice operations (filter, map, chunk, etc.)
//   - stringutil: String manipulation and similarity algorithms
//   - mathutil: Mathematical and statistical operations
//   - fileutil: File system operations with proper error handling
//   - cryptoutil: Secure AES-GCM encryption
//   - randutil: Cryptographically secure random generation
//   - collection: Generic data structures (Stack, Queue, Set, BST)
//   - textnorm: Deterministic text normalization pipelines
//   - cache: Generic in-memory cache (SIEVE eviction, TTL, deduplicated loading)
//   - aws: AWS SDK v2 helpers (S3, DynamoDB, SQS, SSM, Secrets Manager, Lambda)
//
// # Example Usage
//
//	import "github.com/alessiosavi/GoGPUtils/sliceutil"
//
//	numbers := []int{1, 2, 3, 4, 5}
//	evens := sliceutil.Filter(numbers, func(n int) bool { return n%2 == 0 })
//	// evens = [2, 4]
package gputils
