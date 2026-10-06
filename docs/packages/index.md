---
title: Packages
nav_order: 2
has_children: true
---

# Packages

GoGPUtils provides focused utility packages for common programming tasks. Each package is designed to be independent with minimal cross-package dependencies.

## Core Utilities

Most of these packages use only the standard library; `stringutil` also uses `golang.org/x/text` for Unicode text processing:

| Package                     | Description                        | Coverage                                                              |
| --------------------------- | ---------------------------------- | --------------------------------------------------------------------- |
| [sliceutil](sliceutil.md)   | Generic slice operations           | Filter, Map, Reduce, Chunk, GroupBy, Partition, FlatMap, and more     |
| [stringutil](stringutil.md) | String manipulation and similarity | Case conversion, padding, truncation, Levenshtein, Jaro-Winkler, Dice |
| [mathutil](mathutil.md)     | Mathematical operations            | Sum, Average, Median, StdDev, Matrix ops, Vector ops, Number theory   |
| [fileutil](fileutil.md)     | File system operations             | Read/Write, Copy, Move, List, Find, Touch, Temp files                 |
| [cryptoutil](cryptoutil.md) | AES-GCM encryption                 | Encrypt/Decrypt, Key derivation, Hashing                              |
| [randutil](randutil.md)     | Secure random generation           | Secure strings, IDs, choices, sequences                               |
| [collection](collection.md) | Data structures                    | Stack, Queue, Set, Binary Search Tree                                 |
| [cache](cache.md)           | In-memory cache                    | SIEVE eviction, TTL, GetOrLoad, OnEvict, Stats, iteration             |

## Text Processing

| Package                                                     | Description                  | Dependencies        |
| ----------------------------------------------------------- | ---------------------------- | ------------------- |
| [textnorm](textnorm.md)                                     | Text normalization pipelines | `golang.org/x/text` |
| [textnorm/stopwords](textnorm.md#package-textnormstopwords) | Embedded stopword lists      | None                |

## AWS SDK Helpers

These packages wrap AWS SDK v2 for easier use:

| Package                                     | Description                        | AWS Service     |
| ------------------------------------------- | ---------------------------------- | --------------- |
| [aws](aws/index.md)                         | AWS configuration and shared types | -               |
| [aws/s3](aws/s3.md)                         | S3 object and bucket operations    | S3              |
| [aws/dynamodb](aws/dynamodb.md)             | DynamoDB item and query operations | DynamoDB        |
| [aws/sqs](aws/sqs.md)                       | SQS message operations             | SQS             |
| [aws/ssm](aws/ssm.md)                       | Parameter Store operations         | SSM             |
| [aws/secretsmanager](aws/secretsmanager.md) | Secret management                  | Secrets Manager |
| [aws/lambda](aws/lambda.md)                 | Lambda function operations         | Lambda          |

## Internal Packages

| Package                   | Description                                                             |
| ------------------------- | ----------------------------------------------------------------------- |
| `internal/constraints`    | Generic numeric and ordered type constraints (currently no importers)   |
| `aws/internal/pagination` | Pagination helpers for AWS service clients (currently no importers)     |
| `aws/internal/testutil`   | LocalStack test utilities for AWS integration tests                     |

## Package Interaction Diagram

Only the AWS service packages import another GoGPUtils package (the shared `aws` package); the diagram excludes tests:

```mermaid
graph TB
    subgraph Core["Core Utilities"]
        SL[sliceutil]
        STR[stringutil]
        M[mathutil]
        F[fileutil]
        R[randutil]
        CR[cryptoutil]
        COL[collection]
        CA[cache]
    end

    subgraph Text["Text Processing"]
        TN[textnorm]
        SW[textnorm/stopwords]
    end

    subgraph AWS2["AWS SDK v2"]
        AC[aws]
        S3[aws/s3]
        DDB[aws/dynamodb]
        SQS[aws/sqs]
        SSM[aws/ssm]
        SM[aws/secretsmanager]
        L[aws/lambda]
    end

    S3 --> AC
    DDB --> AC
    SQS --> AC
    SSM --> AC
    SM --> AC
    L --> AC
```
