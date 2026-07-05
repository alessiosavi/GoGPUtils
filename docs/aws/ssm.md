---
layout: default
title: SSM Parameter Store
parent: AWS Services
nav_order: 4
---

# SSM Parameter Store

The `aws/ssm` package provides helpers for AWS Systems Manager Parameter Store operations, including getting, listing, and managing parameters.

## Features

- Get parameter values (string, string list, secure string)
- Get multiple parameters in one call
- List parameters with filtering
- Automatic decryption of secure strings

## Client Creation

```go
import "github.com/alessiosavi/GoGPUtils/aws/ssm"

cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
if err != nil {
    return err
}

client, err := ssm.NewClient(cfg)
if err != nil {
    return err
}
```

## Types

### ParameterInfo

Contains metadata about a parameter.

```go
type ParameterInfo struct {
    Name             string
    Type             string
    Value            string
    Version          int64
    LastModifiedDate time.Time
    ARN              string
    DataType         string
}
```

## Get Operations

### GetParameter

Retrieves a parameter value. SecureString parameters are automatically decrypted.

```go
func (c *Client) GetParameter(ctx context.Context, name string, opts ...GetParameterOption) (string, error)
```

**GetParameter Options:**

| Option                    | Description                                     |
| ------------------------- | ----------------------------------------------- |
| `WithDecryption(decrypt)` | Decrypt SecureString parameters (default: true) |

**Example:**

```go
value, err := client.GetParameter(ctx, "/app/config/database_url")
if errors.Is(err, ssm.ErrParameterNotFound) {
    // Handle not found
}

// Without decryption
value, err := client.GetParameter(ctx, "/app/config/value", ssm.WithDecryption(false))
```

### GetParameterInfo

Retrieves a parameter with its metadata.

```go
func (c *Client) GetParameterInfo(ctx context.Context, name string, opts ...GetParameterOption) (*ParameterInfo, error)
```

**Example:**

```go
info, err := client.GetParameterInfo(ctx, "/app/config/database_url")
fmt.Printf("Version: %d, Last Modified: %v\n", info.Version, info.LastModifiedDate)
```

### GetParameters

Retrieves multiple parameters in one call. Returns a map of parameter name to value. Invalid parameter names are returned separately.

```go
func (c *Client) GetParameters(ctx context.Context, names []string, opts ...GetParameterOption) (map[string]string, []string, error)
```

**Example:**

```go
names := []string{"/app/config/a", "/app/config/b", "/app/config/c"}
values, invalid, err := client.GetParameters(ctx, names)

for name, value := range values {
    fmt.Printf("%s = %s\n", name, value)
}

for _, name := range invalid {
    fmt.Printf("Invalid parameter: %s\n", name)
}
```

## List Operations

### ListParametersByPath

Returns all parameters under a path.

```go
func (c *Client) ListParametersByPath(ctx context.Context, path string, opts ...ListParametersByPathOption) ([]ParameterInfo, error)
```

**ListParametersByPath Options:**

| Option                        | Description                         |
| ----------------------------- | ----------------------------------- |
| `WithRecursive(recursive)`    | Enable recursive listing            |
| `WithPathDecryption(decrypt)` | Decrypt SecureString parameters     |
| `WithPathMaxResults(n)`       | Limit number of parameters returned |

**Example:**

```go
// List parameters under a path
params, err := client.ListParametersByPath(ctx, "/app/config/")
for _, p := range params {
    fmt.Printf("%s = %s\n", p.Name, p.Value)
}

// Recursive listing
params, err := client.ListParametersByPath(ctx, "/app/", ssm.WithRecursive(true))
```

### ListParameters

Returns metadata for all parameters (names only).

```go
func (c *Client) ListParameters(ctx context.Context) ([]string, error)
```

**Example:**

```go
names, err := client.ListParameters(ctx)
for _, name := range names {
    fmt.Println(name)
}
```

## Put Operations

### PutParameter

Creates or updates a parameter.

```go
func (c *Client) PutParameter(ctx context.Context, name, value string, opts ...PutParameterOption) error
```

**PutParameter Options:**

| Option                     | Description                             |
| -------------------------- | --------------------------------------- |
| `WithParameterType(t)`     | Set parameter type (default: String)    |
| `WithOverwrite(overwrite)` | Allow overwriting existing parameter    |
| `WithDescription(desc)`    | Set parameter description               |
| `WithKMSKeyID(keyID)`      | Set KMS key for SecureString encryption |

**Example:**

```go
// Simple string parameter
err := client.PutParameter(ctx, "/app/config/database_url", "postgres://...")

// Secure string
err := client.PutParameter(ctx, "/app/secret/api_key", "secret-value",
    ssm.WithParameterType(types.ParameterTypeSecureString),
)

// With overwrite
err := client.PutParameter(ctx, "/app/config/value", "new-value",
    ssm.WithOverwrite(true),
    ssm.WithDescription("Application configuration value"),
)

// With custom KMS key
err := client.PutParameter(ctx, "/app/secret", "value",
    ssm.WithParameterType(types.ParameterTypeSecureString),
    ssm.WithKMSKeyID("alias/my-key"),
)
```

## Delete Operations

### DeleteParameter

Deletes a parameter.

```go
func (c *Client) DeleteParameter(ctx context.Context, name string) error
```

**Example:**

```go
err := client.DeleteParameter(ctx, "/app/config/old-setting")
if errors.Is(err, ssm.ErrParameterNotFound) {
    // Parameter already deleted
}
```

## Error Handling

### Sentinel Errors

| Error                         | Description                      |
| ----------------------------- | -------------------------------- |
| `ErrParameterNotFound`        | Parameter does not exist         |
| `ErrParameterVersionNotFound` | Parameter version does not exist |
| `ErrAccessDenied`             | Access to parameter is denied    |

### Error Handling Example

```go
value, err := client.GetParameter(ctx, "/app/config/key")
if err != nil {
    if errors.Is(err, ssm.ErrParameterNotFound) {
        // Use default value
        value = "default"
    } else if errors.Is(err, ssm.ErrAccessDenied) {
        // Log and alert
        log.Printf("Access denied to parameter")
        return err
    } else {
        return err
    }
}
```

## Testing

```go
type mockSSMAPI struct {
    ssm.API
}

func (m *mockSSMAPI) GetParameter(ctx context.Context, params *ssmsdk.GetParameterInput, optFns ...func(*ssmsdk.Options)) (*ssmsdk.GetParameterOutput, error) {
    return &ssmsdk.GetParameterOutput{
        Parameter: &types.Parameter{
            Name:  aws.String("/app/config/key"),
            Value: aws.String("test-value"),
            Type:  types.ParameterTypeString,
        },
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockSSMAPI{}
    client := ssm.NewClientWithAPI(mock)
    // Use client in tests...
}
```
