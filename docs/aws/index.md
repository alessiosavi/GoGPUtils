---
layout: default
title: AWS Services
nav_order: 3
has_children: true
---

# AWS Services

The `aws` package provides idiomatic Go helpers for AWS SDK v2, designed for production use. It follows a consistent design philosophy across all service clients.

## Design Philosophy

- **No global state**: All clients are explicitly created and passed
- **Context-aware**: All operations accept `context.Context` for cancellation and timeouts
- **Testable**: Interfaces enable mocking without real AWS calls
- **Explicit configuration**: No magic; all settings are visible
- **Safe defaults**: Sensible retry policies out of the box
- **Minimal abstraction**: Helpers augment the SDK, not replace it

## Package Structure

| Package              | Service             | Description                             |
| -------------------- | ------------------- | --------------------------------------- |
| `aws`                | Core                | Configuration, errors, common utilities |
| `aws/s3`             | S3                  | Object storage operations               |
| `aws/dynamodb`       | DynamoDB            | NoSQL database operations               |
| `aws/sqs`            | SQS                 | Message queue operations                |
| `aws/secretsmanager` | Secrets Manager     | Secret storage and retrieval            |
| `aws/ssm`            | SSM Parameter Store | Configuration parameter storage         |
| `aws/lambda`         | Lambda              | Serverless function operations          |

## Quick Start

### Configuration

All service clients are created from a shared AWS configuration:

```go
package main

import (
    "context"
    "log"

    "github.com/alessiosavi/GoGPUtils/aws"
    "github.com/alessiosavi/GoGPUtils/aws/s3"
)

func main() {
    ctx := context.Background()

    // Load AWS configuration
    cfg, err := aws.LoadConfig(ctx,
        aws.WithRegion("us-west-2"),
        aws.WithRetryMaxAttempts(5),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create S3 client
    s3Client, err := s3.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Use the client
    data, err := s3Client.GetObject(ctx, "my-bucket", "my-key")
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Downloaded %d bytes", len(data))
}
```

## Core Configuration

### LoadConfig

```go
func LoadConfig(ctx context.Context, opts ...ConfigOption) (*Config, error)
```

Loads AWS configuration with the specified options. If no options are provided, it uses default configuration from environment variables, shared credentials, and IAM roles.

### Config Options

| Option                      | Description                               |
| --------------------------- | ----------------------------------------- |
| `WithRegion(region)`        | Set the AWS region                        |
| `WithProfile(profile)`      | Use named profile from shared credentials |
| `WithRetryMaxAttempts(n)`   | Set maximum retry attempts (default: 3)   |
| `WithRetryMode(mode)`       | Set retry mode (standard or adaptive)     |
| `WithEndpoint(url)`         | Custom endpoint for LocalStack/testing    |
| `WithCredentials(provider)` | Explicit credentials provider             |
| `WithLoadOption(opt)`       | Raw AWS SDK config load option            |

### Example

```go
cfg, err := aws.LoadConfig(ctx,
    aws.WithRegion("us-west-2"),
    aws.WithProfile("production"),
    aws.WithRetryMaxAttempts(5),
    aws.WithRetryMode(aws.RetryModeAdaptive),
    aws.WithEndpoint("http://localhost:4566"), // LocalStack
)
```

## Common Errors

All operations return wrapped errors with context. The following sentinel errors are defined in the core `aws` package:

| Error               | Description             |
| ------------------- | ----------------------- |
| `ErrNilConfig`      | Nil config provided     |
| `ErrNilClient`      | Nil client provided     |
| `ErrEmptyBucket`    | Empty bucket name       |
| `ErrEmptyKey`       | Empty key/path          |
| `ErrEmptyTable`     | Empty table name        |
| `ErrEmptyQueue`     | Empty queue name/URL    |
| `ErrEmptySecret`    | Empty secret name       |
| `ErrEmptyParameter` | Empty parameter name    |
| `ErrEmptyFunction`  | Empty function name     |
| `ErrInvalidInput`   | Input validation failed |

### Error Types

```go
// ConfigError represents an error loading AWS configuration
type ConfigError struct {
    Err error
}

// OperationError represents an error during an AWS operation
type OperationError struct {
    Service   string
    Operation string
    Err       error
}

// NotFoundError represents a resource not found error
type NotFoundError struct {
    Resource string
    ID       string
}

// ValidationError represents an input validation error
type ValidationError struct {
    Field   string
    Message string
}
```

### Error Handling Patterns

```go
data, err := client.GetObject(ctx, "bucket", "key")
if err != nil {
    // Check for specific sentinel errors
    if errors.Is(err, s3.ErrObjectNotFound) {
        // Handle not found
    }

    // Check operation error type
    var opErr *aws.OperationError
    if errors.As(err, &opErr) {
        log.Printf("Service: %s, Operation: %s", opErr.Service, opErr.Operation)
    }
}
```

## Testing

All clients accept interfaces that can be mocked for testing:

```go
type mockS3API struct {
    s3.API
}

func (m *mockS3API) GetObject(ctx context.Context, params *s3sdk.GetObjectInput, optFns ...func(*s3sdk.Options)) (*s3sdk.GetObjectOutput, error) {
    return &s3sdk.GetObjectOutput{
        Body: io.NopCloser(bytes.NewReader([]byte("test"))),
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockS3API{}
    client := s3.NewClientWithAPI(mock, nil, nil)
    // Use client in tests...
}
```

## Service Documentation

- [S3](s3.md) - Object storage operations
- [DynamoDB](dynamodb.md) - NoSQL database operations
- [SQS](sqs.md) - Message queue operations
- [SSM Parameter Store](ssm.md) - Configuration parameter storage
- [Secrets Manager](secretsmanager.md) - Secret storage and retrieval
- [Lambda](lambda.md) - Serverless function operations
