---
title: Home
nav_order: 1
---

# GoGPUtils

A collection of well-tested, idiomatic Go utilities for common programming tasks. Zero external dependencies for core utilities.

[![Go](https://github.com/alessiosavi/GoGPUtils/actions/workflows/go.yml/badge.svg)](https://github.com/alessiosavi/GoGPUtils/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/alessiosavi/GoGPUtils)](https://goreportcard.com/report/github.com/alessiosavi/GoGPUtils)
[![GoDoc](https://godoc.org/github.com/alessiosavi/GoGPUtils?status.svg)](https://godoc.org/github.com/alessiosavi/GoGPUtils)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **v1 - Experimental**: This library is in its initial release phase. APIs may change in future versions.

## Installation

```bash
go get github.com/alessiosavi/GoGPUtils
```

## Design Philosophy

- **Errors over panics**: All functions return errors instead of panicking
- **Zero global state**: No singletons; all state is explicit
- **Generic when useful**: Uses generics to reduce duplication without over-abstraction
- **Minimal dependencies**: Core library has zero external dependencies
- **Context-aware**: Blocking operations accept `context.Context`

## Package Architecture

```mermaid
graph TB
    subgraph Core[Core Utilities - Zero Dependencies]
        S[sliceutil]
        ST[stringutil]
        M[mathutil]
        F[fileutil]
        R[randutil]
        C[cryptoutil]
        COL[collection]
    end

    subgraph Text[Text Processing]
        TN[textnorm]
        SW[textnorm/stopwords]
    end

    subgraph AWS[AWS SDK v2]
        AW[aws/config]
        S3[aws/s3]
        DDB[aws/dynamodb]
        SQS[aws/sqs]
        SSM[aws/ssm]
        SM[aws/secretsmanager]
        L[aws/lambda]
    end

    TN --> SW
    S3 --> AW
    DDB --> AW
    SQS --> AW
    SSM --> AW
    SM --> AW
    L --> AW
```

## Package Dependency Graph

```mermaid
graph LR
    subgraph Internal[Internal]
        IC[internal/constraints]
    end

    subgraph Utils[Utilities]
        SL[sliceutil]
        STR[stringutil]
        MT[mathutil]
        FL[fileutil]
        RD[randutil]
        CR[cryptoutil]
        CL[collection]
    end

    subgraph TextProc[Text]
        TX[textnorm]
        SW2[stopwords]
    end

    subgraph AWS2[AWS]
        AC[aws/config]
        AS3[aws/s3]
        ADD[aws/dynamodb]
        ASQ[aws/sqs]
        ASM[aws/ssm]
        ASE[aws/secretsmanager]
        AL[aws/lambda]
        AIP[aws/internal/pagination]
        AIT[aws/internal/testutil]
    end

    MT --> IC
    SL --> IC
    STR --> IC
    TX --> STR
    TX --> SW2
    ADD --> AC
    AS3 --> AC
    ASQ --> AC
    ASM --> AC
    ASE --> AC
    AL --> AC
    AS3 --> AIP
    ADD --> AIP
    ASQ --> AIP
    ADD --> AIT
    AS3 --> AIT
    ASE --> AIT
    ASM --> AIT
    ASQ --> AIT
    AL --> AIT
```

## Packages

| Package                                | Description                                                          | Dependencies              |
| -------------------------------------- | -------------------------------------------------------------------- | ------------------------- |
| [`sliceutil`](packages/sliceutil.md)   | Generic slice operations (filter, map, reduce, chunk, etc.)          | `internal/constraints`    |
| [`stringutil`](packages/stringutil.md) | String manipulation and similarity algorithms                        | `internal/constraints`    |
| [`textnorm`](packages/textnorm.md)     | Deterministic text normalization pipelines                           | `stringutil`, `stopwords` |
| [`mathutil`](packages/mathutil.md)     | Mathematical and statistical operations                              | `internal/constraints`    |
| [`fileutil`](packages/fileutil.md)     | File system operations with proper error handling                    | None                      |
| [`cryptoutil`](packages/cryptoutil.md) | Secure AES-GCM encryption                                            | None                      |
| [`randutil`](packages/randutil.md)     | Cryptographically secure random generation                           | None                      |
| [`collection`](packages/collection.md) | Generic data structures (Stack, Queue, Set, BST)                     | None                      |
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
