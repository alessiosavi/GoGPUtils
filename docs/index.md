---
title: Home
nav_order: 1
---

# GoGPUtils

A collection of well-tested, idiomatic Go utilities for common programming tasks. Most utility packages use only the standard library; `stringutil` and `textnorm` also use `golang.org/x/text` for Unicode text processing, and `aws` packages use the AWS SDK v2.

[![CI](https://github.com/alessiosavi/GoGPUtils/actions/workflows/ci.yml/badge.svg)](https://github.com/alessiosavi/GoGPUtils/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/alessiosavi/GoGPUtils.svg)](https://pkg.go.dev/github.com/alessiosavi/GoGPUtils)
[![Go Version](https://img.shields.io/github/go-mod/go-version/alessiosavi/GoGPUtils)](https://github.com/alessiosavi/GoGPUtils/blob/master/go.mod)
[![Latest Tag](https://img.shields.io/github/v/tag/alessiosavi/GoGPUtils?sort=semver)](https://github.com/alessiosavi/GoGPUtils/tags)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/alessiosavi/GoGPUtils)

> **v1 - Experimental**: This library is in its initial release phase. APIs may change in future versions.

## Installation

```bash
go get github.com/alessiosavi/GoGPUtils
```

Requires Go 1.27.1 or later.

## Design Philosophy

- **Errors over panics**: All functions return errors instead of panicking
- **Zero global state**: No singletons; all state is explicit
- **Generic when useful**: Uses generics to reduce duplication without over-abstraction
- **Minimal dependencies**: Utilities use the standard library, with `golang.org/x/text` for Unicode text processing; AWS helpers use the AWS SDK v2
- **Context-aware**: Blocking operations accept `context.Context`

## Package Architecture

```mermaid
graph TB
    subgraph Core[Core Utilities]
        S[sliceutil]
        ST[stringutil]
        M[mathutil]
        F[fileutil]
        R[randutil]
        C[cryptoutil]
        COL[collection]
        CA[cache]
    end

    subgraph Text[Text Processing]
        TN[textnorm]
        SW[textnorm/stopwords]
    end

    subgraph AWS[AWS SDK v2]
        AW[aws]
        S3[aws/s3]
        DDB[aws/dynamodb]
        SQS[aws/sqs]
        SSM[aws/ssm]
        SM[aws/secretsmanager]
        L[aws/lambda]
    end

    S3 --> AW
    DDB --> AW
    SQS --> AW
    SSM --> AW
    SM --> AW
    L --> AW
```

## Package Dependency Graph

The core and text packages import no other GoGPUtils package: `stringutil` and `textnorm` use `golang.org/x/text`, and the rest use only the standard library. The AWS service packages import the shared `aws` package. The diagram shows dependencies between local packages, excluding tests:

```mermaid
graph LR
    subgraph Base["Base Packages"]
        AC[aws]
    end

    subgraph Services["Service Packages"]
        AS3[aws/s3]
        ADD[aws/dynamodb]
        ASQ[aws/sqs]
        ASM[aws/ssm]
        ASE[aws/secretsmanager]
        AL[aws/lambda]
    end

    AS3 --> AC
    ADD --> AC
    ASQ --> AC
    ASM --> AC
    ASE --> AC
    AL --> AC
```

## Packages

| Package                                | Description                                                          | Dependencies              |
| -------------------------------------- | -------------------------------------------------------------------- | ------------------------- |
| [`sliceutil`](packages/sliceutil.md)   | Generic slice operations (filter, map, reduce, chunk, etc.)          | None                      |
| [`stringutil`](packages/stringutil.md) | String manipulation and similarity algorithms                        | `golang.org/x/text`       |
| [`textnorm`](packages/textnorm.md)     | Deterministic text normalization pipelines                           | `golang.org/x/text`       |
| [`mathutil`](packages/mathutil.md)     | Mathematical and statistical operations                              | None                      |
| [`fileutil`](packages/fileutil.md)     | File system operations with proper error handling                    | None                      |
| [`cryptoutil`](packages/cryptoutil.md) | Secure AES-GCM encryption                                            | None                      |
| [`randutil`](packages/randutil.md)     | Cryptographically secure random generation                           | None                      |
| [`collection`](packages/collection.md) | Generic data structures (Stack, Queue, Set, BST)                     | None                      |
| [`cache`](packages/cache.md)           | Generic in-memory cache (SIEVE eviction, TTL, deduplicated loading)  | None                      |
| [`aws`](aws/index.md)                  | AWS SDK v2 helpers (S3, DynamoDB, SQS, SSM, Secrets Manager, Lambda) | `aws-sdk-go-v2`           |

## Quick Start

### Slice Operations

```go
import "github.com/alessiosavi/GoGPUtils/sliceutil"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Filter even numbers
evens := sliceutil.Filter(numbers, func(n int) bool {
    return n%2 == 0
})

// Map transformation
doubled := sliceutil.Map(numbers, func(n int) int {
    return n * 2
})

// Reduce to sum
sum := sliceutil.Reduce(numbers, 0, func(acc, n int) int {
    return acc + n
})
```

### String Similarity

```go
import "github.com/alessiosavi/GoGPUtils/stringutil"

// Levenshtein distance
dist := stringutil.LevenshteinDistance("kitten", "sitting") // 3

// Jaro-Winkler similarity
sim := stringutil.JaroWinklerSimilarity("hello", "hallo", 0.1) // ~0.88
```

### Text Normalization

```go
import "github.com/alessiosavi/GoGPUtils/textnorm"

// Search-optimized normalization
normalized := textnorm.SearchPreset().Run("  Café, go!  ")

// Canonical normalization
clean := textnorm.CanonicalPreset().Run("  Hello, World!  ")
```

## Testing

Run all tests:

```bash
go test ./...
```

Run with race detector:

```bash
go test -race ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please ensure:

1. All new code has tests
2. Tests pass with race detector enabled
3. Code follows Go conventions (`gofmt`, `golint`)
4. Public APIs are documented

## License

MIT License - see [LICENSE](LICENSE) for details.
