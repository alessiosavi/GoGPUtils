---
layout: default
title: Lambda
parent: AWS Services
nav_order: 6
---

# Lambda

The `aws/lambda` package provides helpers for AWS Lambda operations, including invoking functions, listing functions, and deploying code.

## Features

- Invoke functions (sync and async)
- List and describe functions
- Deploy code from S3 or ZIP
- Manage function configuration

## Client Creation

```go
import "github.com/alessiosavi/GoGPUtils/aws/lambda"

cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
if err != nil {
    return err
}

client, err := lambda.NewClient(cfg)
if err != nil {
    return err
}
```

## Types

### InvokeResult

Contains the result of a Lambda invocation.

```go
type InvokeResult struct {
    StatusCode      int32
    Payload         []byte
    FunctionError   string
    LogResult       string
    ExecutedVersion string
}
```

Methods:

- `Unmarshal(dest any) error` - Unmarshals payload into destination
- `HasError() bool` - Returns true if function returned an error

### FunctionInfo

Contains information about a Lambda function.

```go
type FunctionInfo struct {
    Name          string
    ARN           string
    Description   string
    Runtime       string
    Handler       string
    MemorySize    int32
    Timeout       int32
    LastModified  string
    CodeSize      int64
    State         string
    StateReason   string
    Architectures []string
}
```

## Invoke Operations

### Invoke

Calls a Lambda function synchronously and returns the result.

```go
func (c *Client) Invoke(ctx context.Context, functionName string, payload []byte, opts ...InvokeOption) (*InvokeResult, error)
```

**Invoke Options:**

| Option                     | Description                                          |
| -------------------------- | ---------------------------------------------------- |
| `WithInvocationType(t)`    | Set invocation type (RequestResponse, Event, DryRun) |
| `WithLogType(t)`           | Include logs in response (LogTypeTail)               |
| `WithQualifier(qualifier)` | Function version or alias                            |

**Example:**

```go
payload := []byte(`{"name": "John"}`)
result, err := client.Invoke(ctx, "my-function", payload)
if err != nil {
    return err
}

if result.HasError() {
    return fmt.Errorf("function error: %s", result.FunctionError)
}

fmt.Println(string(result.Payload))

// With logs
result, err := client.Invoke(ctx, "my-function", payload, lambda.WithLogType(types.LogTypeTail))
fmt.Println("Logs:", result.LogResult)

// Invoke specific version
result, err := client.Invoke(ctx, "my-function", payload, lambda.WithQualifier("prod"))
```

### InvokeJSON

Calls a Lambda function with a JSON payload. The input is automatically marshaled to JSON.

```go
func (c *Client) InvokeJSON(ctx context.Context, functionName string, input any, opts ...InvokeOption) (*InvokeResult, error)
```

**Example:**

```go
type MyInput struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

input := MyInput{Name: "John", Age: 30}
result, err := client.InvokeJSON(ctx, "my-function", input)
```

### InvokeAsync

Calls a Lambda function asynchronously (fire and forget).

```go
func (c *Client) InvokeAsync(ctx context.Context, functionName string, payload []byte) error
```

**Example:**

```go
err := client.InvokeAsync(ctx, "my-function", payload)
```

## List Operations

### ListFunctions

Returns all Lambda functions in the account.

```go
func (c *Client) ListFunctions(ctx context.Context) ([]FunctionInfo, error)
```

**Example:**

```go
functions, err := client.ListFunctions(ctx)
for _, fn := range functions {
    fmt.Printf("%s (%s, %dMB)\n", fn.Name, fn.Runtime, fn.MemorySize)
}
```

### ListFunctionNames

Returns just the names of all Lambda functions.

```go
func (c *Client) ListFunctionNames(ctx context.Context) ([]string, error)
```

**Example:**

```go
names, err := client.ListFunctionNames(ctx)
```

### GetFunction

Retrieves information about a Lambda function.

```go
func (c *Client) GetFunction(ctx context.Context, functionName string) (*FunctionInfo, error)
```

**Example:**

```go
info, err := client.GetFunction(ctx, "my-function")
fmt.Printf("Runtime: %s, Memory: %dMB\n", info.Runtime, info.MemorySize)
```

## Deploy Operations

### DeployFromS3

Updates a function's code from an S3 object.

```go
func (c *Client) DeployFromS3(ctx context.Context, functionName, bucket, key string) error
```

**Example:**

```go
err := client.DeployFromS3(ctx, "my-function", "my-bucket", "deployments/code.zip")
```

### DeployFromZip

Updates a function's code from a local ZIP file.

```go
func (c *Client) DeployFromZip(ctx context.Context, functionName, zipPath string) error
```

**Example:**

```go
err := client.DeployFromZip(ctx, "my-function", "function.zip")
```

### DeployFromBytes

Updates a function's code from a ZIP byte slice.

```go
func (c *Client) DeployFromBytes(ctx context.Context, functionName string, zipData []byte) error
```

**Example:**

```go
zipData, err := os.ReadFile("function.zip")
if err != nil {
    return err
}
err = client.DeployFromBytes(ctx, "my-function", zipData)
```

## Delete Operations

### DeleteFunction

Deletes a Lambda function.

```go
func (c *Client) DeleteFunction(ctx context.Context, functionName string) error
```

**Example:**

```go
err := client.DeleteFunction(ctx, "my-function")
if errors.Is(err, lambda.ErrFunctionNotFound) {
    // Function already deleted
}
```

## Tag Operations

### GetTags

Retrieves the tags for a Lambda function.

```go
func (c *Client) GetTags(ctx context.Context, functionARN string) (map[string]string, error)
```

**Example:**

```go
tags, err := client.GetTags(ctx, "arn:aws:lambda:us-west-2:123456789:function:my-function")
for k, v := range tags {
    fmt.Printf("%s = %s\n", k, v)
}
```

## Error Handling

### Sentinel Errors

| Error                 | Description                |
| --------------------- | -------------------------- |
| `ErrFunctionNotFound` | Function does not exist    |
| `ErrInvocationFailed` | Function invocation failed |
| `ErrFunctionError`    | Function returned an error |

### Error Handling Example

```go
result, err := client.Invoke(ctx, "my-function", payload)
if err != nil {
    if errors.Is(err, lambda.ErrFunctionNotFound) {
        // Function doesn't exist
        return err
    }
    return err
}

if result.HasError() {
    // Function executed but returned an error
    log.Printf("Function error: %s", result.FunctionError)
    return lambda.ErrFunctionError
}

var response MyResponse
err = result.Unmarshal(&response)
if err != nil {
    return err
}
```

## Testing

```go
type mockLambdaAPI struct {
    lambda.API
}

func (m *mockLambdaAPI) Invoke(ctx context.Context, params *lambdasdk.InvokeInput, optFns ...func(*lambdasdk.Options)) (*lambdasdk.InvokeOutput, error) {
    return &lambdasdk.InvokeOutput{
        StatusCode: 200,
        Payload:    []byte(`{"status":"ok"}`),
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockLambdaAPI{}
    client := lambda.NewClientWithAPI(mock)
    // Use client in tests...
}
```
