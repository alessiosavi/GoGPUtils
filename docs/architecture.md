---
title: Architecture
nav_order: 4
---

# Architecture

This page describes the structure, design patterns, and relationships between GoGPUtils packages.

## Package Overview

GoGPUtils is organized into three layers: **Core utilities** (no external dependencies), **Text processing** (`stringutil` and `textnorm`, using `golang.org/x/text`), and **AWS SDK wrappers** (many external dependencies).

```mermaid
graph TB
    subgraph Core["Core Utilities (stdlib only)"]
        SL[sliceutil]
        M[mathutil]
        F[fileutil]
        RD[randutil]
        CR[cryptoutil]
        COL[collection]
        CA[cache]
    end

    subgraph Text["Text Processing (+1 dep)"]
        STR[stringutil]
        TN[textnorm]
        SW[textnorm/stopwords]
    end

    subgraph AWS["AWS SDK v2 (+15 deps)"]
        S3[aws/s3]
        DDB[aws/dynamodb]
        SQS[aws/sqs]
        SSM[aws/ssm]
        SM[aws/secretsmanager]
        L[aws/lambda]
    end

    style Core fill:#e8f5e9
    style Text fill:#fff3e0
    style AWS fill:#fce4ec
```

## Package Sizes

| Package              | Lines of Code | Functions | Tests | Test Ratio |
| -------------------- | ------------- | --------- | ----- | ---------- |
| `stringutil`         | 3,609         | 76        | 88    | 1.16x      |
| `sliceutil`          | 2,367         | 117       | 58    | 0.50x      |
| `mathutil`           | 1,991         | 108       | 57    | 0.53x      |
| `fileutil`           | 1,735         | 87        | 42    | 0.48x      |
| `collection`         | 1,683         | 112       | 40    | 0.36x      |
| `textnorm`           | 1,057         | 59        | 38    | 0.64x      |
| `randutil`           | 1,010         | 62        | 27    | 0.44x      |
| `cryptoutil`         | 665           | 34        | 19    | 0.56x      |
| `aws/s3`             | 1,850         | ~45       | ~25   | 0.56x      |
| `aws/dynamodb`       | 1,714         | ~55       | ~20   | 0.36x      |
| `aws/secretsmanager` | 1,353         | ~20       | ~15   | 0.75x      |
| `aws/sqs`            | 1,182         | ~30       | ~15   | 0.50x      |
| `aws/ssm`            | 1,067         | ~22       | ~12   | 0.55x      |
| `aws/lambda`         | 588           | ~15       | ~8    | 0.53x      |

**Observations:**

- `stringutil` is the largest package (similarity algorithms + cleaning + validation)
- `collection` has the most functions (BST, Stack, Queue, Set) but fewest tests per function
- AWS packages are roughly evenly structured: ~50% test ratio

## Import Dependencies

### Core Layer (stdlib only)

The `collection`, `cache`, `sliceutil`, `mathutil`, `fileutil`, `cryptoutil`, and `randutil`
packages import only the Go standard library (excluding tests).

`internal/constraints` defines numeric and ordered type constraints but currently has no importers.
`mathutil` declares its own `Number`, `Integer`, and `Float` constraints;
`sliceutil` and `collection` use `cmp.Ordered` for ordered types.

### Text Layer

Non-standard-library imports are summarized below:

```
stringutil → golang.org/x/text (Unicode normalization and accent removal)
textnorm   → golang.org/x/text (Unicode normalization, case folding, and width folding)
```

`stringutil` and `textnorm` share one external dependency (`golang.org/x/text`).
`textnorm/stopwords` uses only the standard library.

### AWS Layer

All AWS service packages import the shared `aws` package, which contains
`config.go` and `errors.go`, plus the corresponding AWS SDK v2 service packages.
The diagram shows dependencies between local packages, excluding tests:

```mermaid
graph LR
    subgraph Base["Base Packages"]
        AWS[aws]
    end

    subgraph Services["Service Packages"]
        S3[aws/s3]
        DDB[aws/dynamodb]
        SQS[aws/sqs]
        SSM[aws/ssm]
        SM[aws/secretsmanager]
        L[aws/lambda]
    end

    S3 --> AWS
    DDB --> AWS
    SQS --> AWS
    SSM --> AWS
    SM --> AWS
    L --> AWS
```

The service packages do not import one another. `aws/internal/pagination` currently
has no importers. `aws/internal/testutil` imports `aws` and is used by the S3,
DynamoDB, SQS, SSM, and Secrets Manager integration tests.

All service constructors accept a `*aws.Config`, loaded separately with
`aws.LoadConfig`; the S3 constructor also accepts client options.

## Design Patterns

### 1. Non-Mutating by Default

Functions return new values. The original is never modified unless explicitly named:

```go
// Returns NEW slice
filtered := sliceutil.Filter(original, predicate)

// Modifies IN PLACE (explicitly named)
sliceutil.FilterInPlace(&original, predicate)
```

### 2. Error Wrapping

Instead of returning raw errors, AWS packages wrap them:

```go
// aws/errors.go provides consistent wrapping
func WrapError(err error, msg string) error {
    if err == nil { return nil }
    return fmt.Errorf("%s: %w", msg, err)
}

// Usage in every AWS service
if err != nil {
    return errors.WrapError(err, "failed to put object")
}
```

This gives you error chains like: `failed to put object: operation error S3: PutObject, ...`

### 3. Option Pattern (AWS)

AWS service functions use variadic option functions:

```go
// S3 put with options
err := client.PutObject(ctx, bucket, key, data,
    s3.WithMetadata(map[string]string{"version": "1.0"}),
    s3.WithContentType("application/json"),
)
```

### 4. Pipeline Pattern (textnorm)

Text normalization uses immutable pipelines:

```go
pipeline := textnorm.New().
    Then(textnorm.RemoveAccents).
    Then(textnorm.FoldCase).
    Then(textnorm.CollapseWhitespace)

result := pipeline.Run(input)  // input is unchanged
```

## Testing Patterns

### Table-Driven Tests

Every package uses Go's idiomatic table-driven pattern:

```go
func TestSum(t *testing.T) {
    tests := []struct{
        name     string
        input    []int
        expected int
    }{
        {"empty", []int{}, 0},
        {"single", []int{5}, 5},
        {"multiple", []int{1, 2, 3}, 6},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Sum(tt.input)
            if got != tt.expected {
                t.Errorf("Sum() = %v, want %v", got, tt.expected)
            }
        })
    }
}
```

### AWS Integration Tests

AWS packages test against LocalStack (a local AWS emulator):

```go
func TestS3PutAndGet(t *testing.T) {
    ctx := context.Background()
    client := setupTest(t)  // Creates LocalStack container

    bucket := "test-bucket"
    err := client.CreateBucket(ctx, bucket)
    require.NoError(t, err)

    err = client.PutObject(ctx, bucket, "key", []byte("hello"))
    require.NoError(t, err)

    data, err := client.GetObject(ctx, bucket, "key")
    require.NoError(t, err)
    assert.Equal(t, "hello", string(data))
}
```

## Complexity Hotspots

Based on lines of code per function, these are the most complex areas:

| Package                    | Complexity Driver                                               |
| -------------------------- | --------------------------------------------------------------- |
| `stringutil/similarity.go` | 7 similarity algorithms (Levenshtein, Jaro-Winkler, Dice, etc.) |
| `sliceutil/sliceutil.go`   | 59 generic functions with type constraints                      |
| `collection/bst.go`        | BST implementation with recursive operations                    |
| `aws/dynamodb`             | Query/scan builder with many option combinations                |
| `aws/s3`                   | Object operations with multipart upload support                 |

## External Dependencies

| Package  | External Dependencies   | Why                   |
| -------- | ----------------------- | --------------------- |
| Core (6) | 0                       | Standard library only |
| stringutil, textnorm | 1 shared (`golang.org/x/text`) | Unicode text processing |
| textnorm/stopwords | 0              | Embedded word lists   |
| aws/*    | 15 (`aws-sdk-go-v2/*`)  | AWS service clients   |
