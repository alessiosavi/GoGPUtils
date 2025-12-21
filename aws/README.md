# AWS Utilities for Go

Idiomatic Go helpers for AWS SDK v2, designed for production use.

## Features

- **No global state**: All clients are explicitly created and passed
- **Context-aware**: All operations accept `context.Context` for cancellation and timeouts
- **Testable**: Interfaces enable mocking without real AWS calls
- **Explicit configuration**: No magic; all settings are visible
- **Safe defaults**: Sensible retry policies out of the box
- **Minimal abstraction**: Helpers augment the SDK, not replace it

## Installation

```bash
go get github.com/alessiosavi/GoGPUtils/aws
```

## Supported Services

| Package | Service | Description |
|---------|---------|-------------|
| `aws` | Core | Configuration, errors, common utilities |
| `aws/s3` | S3 | Object storage operations |
| `aws/dynamodb` | DynamoDB | NoSQL database operations |
| `aws/sqs` | SQS | Message queue operations |
| `aws/secretsmanager` | Secrets Manager | Secret storage and retrieval |
| `aws/ssm` | SSM Parameter Store | Configuration parameter storage |
| `aws/lambda` | Lambda | Serverless function operations |

## Quick Start

### Configuration

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

## Service Examples

### S3

```go
// Upload an object
err := client.PutObject(ctx, "bucket", "key", []byte("content"),
    s3.WithContentType("text/plain"),
    s3.WithStorageClass(types.StorageClassStandardIa),
)

// Download an object
data, err := client.GetObject(ctx, "bucket", "key")

// List objects with prefix
objects, err := client.ListObjects(ctx, "bucket",
    s3.WithPrefix("logs/"),
    s3.WithMaxKeys(100),
)

// Check if object exists
exists, err := client.ObjectExists(ctx, "bucket", "key")

// Copy/Move objects
err := client.CopyObject(ctx, "src-bucket", "src-key", "dst-bucket", "dst-key")
err := client.MoveObject(ctx, "bucket", "old-key", "bucket", "new-key")

// Parse S3 URIs
bucket, key, err := s3.ParseS3Path("s3://my-bucket/path/to/file.txt")
```

### DynamoDB

```go
// Define your item struct
type User struct {
    ID    string `dynamodbav:"pk"`
    Email string `dynamodbav:"email"`
    Name  string `dynamodbav:"name"`
}

// Put an item
user := User{ID: "user-123", Email: "alice@example.com", Name: "Alice"}
err := client.PutItem(ctx, "users", user)

// Get an item
var result User
err := client.GetItem(ctx, "users", dynamodb.Key{"pk": "user-123"}, &result)

// Query items
keyExpr := expression.Key("pk").Equal(expression.Value("user-123"))
result, err := client.Query(ctx, "users", keyExpr)

var users []User
err = result.Unmarshal(&users)

// Batch operations
items := []any{user1, user2, user3}
unprocessed, err := client.BatchWriteItems(ctx, "users", items)
```

### SQS

```go
// Send a message
msgID, err := client.SendMessage(ctx, queueURL, "Hello World",
    sqs.WithDelaySeconds(60),
)

// Receive messages with long polling
messages, err := client.ReceiveMessages(ctx, queueURL,
    sqs.WithMaxMessages(10),
    sqs.WithWaitTimeSeconds(20),
)

// Process and delete
for _, msg := range messages {
    process(msg.Body)
    client.DeleteMessage(ctx, queueURL, msg.ReceiptHandle)
}

// Resolve queue name to URL
queueURL, err := client.GetQueueURL(ctx, "my-queue")
```

### Secrets Manager

```go
// Get a secret as string
apiKey, err := client.GetSecretString(ctx, "api-key")

// Get and unmarshal JSON secret
type DBConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

var config DBConfig
err := client.GetSecretJSON(ctx, "db-credentials", &config)

// List all secrets
secrets, err := client.ListSecrets(ctx,
    secretsmanager.WithNameFilter("prod/"),
)
```

### SSM Parameter Store

```go
// Get a parameter (auto-decrypts SecureString)
value, err := client.GetParameter(ctx, "/app/config/database_url")

// Get multiple parameters
values, invalid, err := client.GetParameters(ctx, []string{
    "/app/config/a",
    "/app/config/b",
})

// List parameters by path
params, err := client.ListParametersByPath(ctx, "/app/config/",
    ssm.WithRecursive(true),
)

// Put a parameter
err := client.PutParameter(ctx, "/app/config/key", "value",
    ssm.WithParameterType(types.ParameterTypeSecureString),
    ssm.WithOverwrite(true),
)
```

### Lambda

```go
// Invoke synchronously
result, err := client.Invoke(ctx, "my-function", payload)
if result.HasError() {
    log.Printf("Function error: %s", result.FunctionError)
}

var response MyResponse
err = result.Unmarshal(&response)

// Invoke asynchronously
err := client.InvokeAsync(ctx, "my-function", payload)

// List functions
functions, err := client.ListFunctions(ctx)

// Deploy from S3
err := client.DeployFromS3(ctx, "my-function", "bucket", "code.zip")

// Deploy from local ZIP
err := client.DeployFromZip(ctx, "my-function", "function.zip")
```

## Testing

All clients accept interfaces that can be mocked for testing:

```go
package mypackage_test

import (
    "context"
    "testing"

    "github.com/alessiosavi/GoGPUtils/aws/s3"
)

type mockS3API struct {
    s3.API
    getObjectData []byte
}

func (m *mockS3API) GetObject(ctx context.Context, params *s3sdk.GetObjectInput, optFns ...func(*s3sdk.Options)) (*s3sdk.GetObjectOutput, error) {
    return &s3sdk.GetObjectOutput{
        Body: io.NopCloser(bytes.NewReader(m.getObjectData)),
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockS3API{getObjectData: []byte("test content")}
    client := s3.NewClientWithAPI(mock, nil, nil)

    data, err := client.GetObject(context.Background(), "bucket", "key")
    if err != nil {
        t.Fatal(err)
    }

    if string(data) != "test content" {
        t.Errorf("unexpected data: %s", data)
    }
}
```

## Configuration Options

### AWS Config

```go
cfg, err := aws.LoadConfig(ctx,
    aws.WithRegion("us-west-2"),           // Set region
    aws.WithProfile("production"),         // Use named profile
    aws.WithRetryMaxAttempts(5),           // Set retry attempts
    aws.WithRetryMode(aws.RetryModeAdaptive), // Set retry mode
    aws.WithEndpoint("http://localhost:4566"), // Custom endpoint (LocalStack)
    aws.WithCredentials(credsProvider),    // Custom credentials
)
```

### Client Options

Each service client supports additional options:

```go
// S3 with custom part size for multipart uploads
s3Client, err := s3.NewClient(cfg,
    s3.WithUploaderConcurrency(10),
    s3.WithDownloaderConcurrency(10),
    s3.WithPartSize(10 * 1024 * 1024), // 10 MB parts
)
```

## Error Handling

All errors are wrapped with context:

```go
data, err := client.GetObject(ctx, "bucket", "key")
if err != nil {
    // Check for specific errors
    if errors.Is(err, s3.ErrObjectNotFound) {
        // Handle not found
    }

    // Or check the operation error type
    var opErr *aws.OperationError
    if errors.As(err, &opErr) {
        log.Printf("Service: %s, Operation: %s", opErr.Service, opErr.Operation)
    }
}
```

## Migration from Old Implementation

If migrating from the old AWS helpers (commit b7a0843), note these changes:

1. **No global clients**: Create clients explicitly with `NewClient(cfg)`
2. **No panics**: All errors are returned, never panic
3. **Context required**: All operations require `context.Context`
4. **Explicit imports**: Import each service package separately
5. **Interface-based**: Use `NewClientWithAPI()` for testing

### Before (old)

```go
import "github.com/alessiosavi/GoGPUtils/aws/S3"

// Global client initialized in init()
data, err := S3utils.GetObject("bucket", "key")
```

### After (new)

```go
import (
    "github.com/alessiosavi/GoGPUtils/aws"
    "github.com/alessiosavi/GoGPUtils/aws/s3"
)

cfg, _ := aws.LoadConfig(ctx)
client, _ := s3.NewClient(cfg)
data, err := client.GetObject(ctx, "bucket", "key")
```

## Design Principles

1. **Explicit over implicit**: No hidden state or magic configuration
2. **Composition over inheritance**: Small, focused packages
3. **Errors over panics**: All failures return errors
4. **Interfaces for testing**: Every client can be mocked
5. **Context everywhere**: Cancellation and timeouts are first-class
6. **Safe defaults**: Sensible retry and timeout configurations
7. **Minimal surface**: Only expose what's necessary

## Contributing

Contributions are welcome! Please ensure:

1. All code follows Go conventions (run `gofmt`)
2. Tests cover new functionality
3. Documentation is updated
4. No breaking changes to public API

## License

MIT License - see LICENSE file for details.
