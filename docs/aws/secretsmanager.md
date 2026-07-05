---
layout: default
title: Secrets Manager
parent: AWS Services
nav_order: 5
---

# Secrets Manager

The `aws/secretsmanager` package provides helpers for AWS Secrets Manager operations, including getting secret values, listing secrets, and managing secrets.

## Features

- Get secret values (string or binary)
- Automatic JSON unmarshaling
- List secrets with filtering
- Create, update, and delete secrets

## Client Creation

```go
import "github.com/alessiosavi/GoGPUtils/aws/secretsmanager"

cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
if err != nil {
    return err
}

client, err := secretsmanager.NewClient(cfg)
if err != nil {
    return err
}
```

## Types

### SecretInfo

Contains metadata about a secret.

```go
type SecretInfo struct {
    Name             string
    ARN              string
    Description      string
    CreatedDate      time.Time
    LastChangedDate  time.Time
    LastAccessedDate time.Time
    Tags             map[string]string
}
```

## Get Operations

### GetSecretString

Retrieves a secret value as a string.

```go
func (c *Client) GetSecretString(ctx context.Context, secretName string, opts ...GetSecretOption) (string, error)
```

**GetSecret Options:**

| Option                    | Description                                                           |
| ------------------------- | --------------------------------------------------------------------- |
| `WithVersionID(id)`       | Retrieve specific version                                             |
| `WithVersionStage(stage)` | Retrieve specific version stage (AWSCURRENT, AWSPREVIOUS, AWSPENDING) |

**Example:**

```go
apiKey, err := client.GetSecretString(ctx, "api-key")
if errors.Is(err, secretsmanager.ErrSecretNotFound) {
    // Handle not found
}

// Get previous version
prevKey, err := client.GetSecretString(ctx, "api-key", secretsmanager.WithVersionStage("AWSPREVIOUS"))
```

### GetSecretBinary

Retrieves a secret value as binary data.

```go
func (c *Client) GetSecretBinary(ctx context.Context, secretName string, opts ...GetSecretOption) ([]byte, error)
```

**Example:**

```go
certData, err := client.GetSecretBinary(ctx, "tls-certificate")
```

### GetSecretJSON

Retrieves a secret and unmarshals it into the provided destination. The secret value must be valid JSON.

```go
func (c *Client) GetSecretJSON(ctx context.Context, secretName string, dest any, opts ...GetSecretOption) error
```

**Example:**

```go
type DBConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
}

var config DBConfig
err := client.GetSecretJSON(ctx, "db-credentials", &config)
// config.Host, config.Port, etc.
```

## List Operations

### ListSecrets

Returns all secrets in the account.

```go
func (c *Client) ListSecrets(ctx context.Context, opts ...ListSecretsOption) ([]SecretInfo, error)
```

**ListSecrets Options:**

| Option                      | Description                      |
| --------------------------- | -------------------------------- |
| `WithMaxResults(n)`         | Limit number of secrets returned |
| `WithNameFilter(prefix)`    | Filter by name prefix            |
| `WithTagFilter(key, value)` | Filter by tag                    |

**Example:**

```go
// List all secrets
secrets, err := client.ListSecrets(ctx)
for _, secret := range secrets {
    fmt.Println(secret.Name)
}

// Filter by name prefix
secrets, err := client.ListSecrets(ctx, secretsmanager.WithNameFilter("prod/"))

// Filter by tag
secrets, err := client.ListSecrets(ctx, secretsmanager.WithTagFilter("env", "production"))
```

### ListSecretNames

Returns just the names of all secrets.

```go
func (c *Client) ListSecretNames(ctx context.Context, opts ...ListSecretsOption) ([]string, error)
```

**Example:**

```go
names, err := client.ListSecretNames(ctx)
```

### DescribeSecret

Returns metadata about a secret without retrieving its value.

```go
func (c *Client) DescribeSecret(ctx context.Context, secretName string) (*SecretInfo, error)
```

**Example:**

```go
info, err := client.DescribeSecret(ctx, "my-secret")
fmt.Printf("Last changed: %v\n", info.LastChangedDate)
```

## Create Operations

### CreateSecretString

Creates a new secret with a string value.

```go
func (c *Client) CreateSecretString(ctx context.Context, name, value string, opts ...CreateSecretOption) error
```

**CreateSecret Options:**

| Option                        | Description                |
| ----------------------------- | -------------------------- |
| `WithCreateDescription(desc)` | Set description            |
| `WithCreateTags(tags)`        | Set tags                   |
| `WithCreateKMSKeyID(keyID)`   | Set KMS key for encryption |

**Example:**

```go
// Simple creation
err := client.CreateSecretString(ctx, "my-secret", "secret-value")

// With metadata
err := client.CreateSecretString(ctx, "my-secret", "secret-value",
    secretsmanager.WithCreateDescription("API key"),
    secretsmanager.WithCreateTags(map[string]string{"env": "prod"}),
    secretsmanager.WithCreateKMSKeyID("alias/my-key"),
)
```

### CreateSecretBinary

Creates a new secret with binary data.

```go
func (c *Client) CreateSecretBinary(ctx context.Context, name string, value []byte, opts ...CreateSecretOption) error
```

**Example:**

```go
err := client.CreateSecretBinary(ctx, "my-cert", certData)
```

## Update Operations

### UpdateSecretString

Updates an existing secret with a new string value.

```go
func (c *Client) UpdateSecretString(ctx context.Context, name, value string) error
```

**Example:**

```go
err := client.UpdateSecretString(ctx, "my-secret", "new-value")
if errors.Is(err, secretsmanager.ErrSecretNotFound) {
    // Secret doesn't exist
}
```

### UpdateSecretBinary

Updates an existing secret with new binary data.

```go
func (c *Client) UpdateSecretBinary(ctx context.Context, name string, value []byte) error
```

**Example:**

```go
err := client.UpdateSecretBinary(ctx, "my-cert", newCertData)
```

## Delete Operations

### DeleteSecret

Deletes a secret. By default, secrets have a recovery window of 30 days.

```go
func (c *Client) DeleteSecret(ctx context.Context, name string, opts ...DeleteSecretOption) error
```

**DeleteSecret Options:**

| Option                     | Description                                      |
| -------------------------- | ------------------------------------------------ |
| `WithForceDelete()`        | Force immediate deletion without recovery window |
| `WithRecoveryWindow(days)` | Set recovery window in days (minimum 7)          |

**Example:**

```go
// Default deletion (30 day recovery window)
err := client.DeleteSecret(ctx, "my-secret")

// Force immediate deletion
err := client.DeleteSecret(ctx, "my-secret", secretsmanager.WithForceDelete())

// Custom recovery window
err := client.DeleteSecret(ctx, "my-secret", secretsmanager.WithRecoveryWindow(7))
```

## Error Handling

### Sentinel Errors

| Error               | Description                |
| ------------------- | -------------------------- |
| `ErrSecretNotFound` | Secret does not exist      |
| `ErrSecretDeleted`  | Secret is deleted          |
| `ErrAccessDenied`   | Access to secret is denied |

### Error Handling Example

```go
apiKey, err := client.GetSecretString(ctx, "api-key")
if err != nil {
    if errors.Is(err, secretsmanager.ErrSecretNotFound) {
        // Secret doesn't exist - may need to create it
        return err
    }
    if errors.Is(err, secretsmanager.ErrSecretDeleted) {
        // Secret was deleted - may need to restore
        return err
    }
    return err
}
```

## Testing

```go
type mockSecretsManagerAPI struct {
    secretsmanager.API
}

func (m *mockSecretsManagerAPI) GetSecretValue(ctx context.Context, params *smsdk.GetSecretValueInput, optFns ...func(*smsdk.Options)) (*smsdk.GetSecretValueOutput, error) {
    return &smsdk.GetSecretValueOutput{
        SecretString: aws.String(`{"host":"localhost","port":5432}`),
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockSecretsManagerAPI{}
    client := secretsmanager.NewClientWithAPI(mock)
    // Use client in tests...
}
```
