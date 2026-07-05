---
layout: default
title: DynamoDB
parent: AWS Services
nav_order: 2
---

# DynamoDB

The `aws/dynamodb` package provides helpers for Amazon DynamoDB operations, including item operations, batch operations, and query/scan with pagination support.

## Features

- Item operations: Get, Put, Delete, Update with automatic marshaling
- Batch operations: BatchGet, BatchWrite with automatic chunking
- Scan and Query with pagination support
- Table operations: Create, Delete

## Client Creation

```go
import "github.com/alessiosavi/GoGPUtils/aws/dynamodb"

cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
if err != nil {
    return err
}

client, err := dynamodb.NewClient(cfg)
if err != nil {
    return err
}
```

## Types

### Key

Represents a DynamoDB item key. Keys should use string, number, or binary attribute values.

```go
type Key map[string]any
```

**Example:**

```go
// Simple partition key
key := dynamodb.Key{"pk": "user-123"}

// Composite key (partition + sort)
key := dynamodb.Key{"pk": "user-123", "sk": "profile"}
```

### QueryResult

Contains the results of a Query operation.

```go
type QueryResult struct {
    Items            []map[string]types.AttributeValue
    Count            int32
    ScannedCount     int32
    LastEvaluatedKey map[string]types.AttributeValue
}
```

Methods:

- `Unmarshal(dest any) error` - Unmarshals results into a slice
- `HasMorePages() bool` - Returns true if more results exist

### ScanResult

Contains the results of a Scan operation.

```go
type ScanResult struct {
    Items            []map[string]types.AttributeValue
    Count            int32
    ScannedCount     int32
    LastEvaluatedKey map[string]types.AttributeValue
}
```

Methods:

- `Unmarshal(dest any) error` - Unmarshals results into a slice
- `HasMorePages() bool` - Returns true if more results exist

## Item Operations

### GetItem

Retrieves an item from a DynamoDB table. The result is unmarshaled into the provided destination.

```go
func (c *Client) GetItem(ctx context.Context, tableName string, key Key, dest any) error
```

**Example:**

```go
type User struct {
    ID    string `dynamodbav:"pk"`
    Email string `dynamodbav:"email"`
    Name  string `dynamodbav:"name"`
}

var user User
err := client.GetItem(ctx, "users", dynamodb.Key{"pk": "user-123"}, &user)
if errors.Is(err, dynamodb.ErrItemNotFound) {
    // Handle not found
}
```

### GetItemRaw

Retrieves an item and returns the raw attribute value map.

```go
func (c *Client) GetItemRaw(ctx context.Context, tableName string, key Key) (map[string]types.AttributeValue, error)
```

**Example:**

```go
item, err := client.GetItemRaw(ctx, "users", dynamodb.Key{"pk": "user-123"})
```

### PutItem

Writes an item to a DynamoDB table. The item is automatically marshaled from the provided struct.

```go
func (c *Client) PutItem(ctx context.Context, tableName string, item any) error
```

**Example:**

```go
user := User{ID: "user-123", Email: "alice@example.com", Name: "Alice"}
err := client.PutItem(ctx, "users", user)
```

### PutItemIfNotExists

Writes an item only if it doesn't already exist. Returns `ErrConditionalCheckFailed` if the item exists.

```go
func (c *Client) PutItemIfNotExists(ctx context.Context, tableName string, item any, pkAttribute string) error
```

**Example:**

```go
user := User{ID: "user-123", Email: "alice@example.com"}
err := client.PutItemIfNotExists(ctx, "users", user, "pk")
if errors.Is(err, dynamodb.ErrConditionalCheckFailed) {
    // Item already exists
}
```

### DeleteItem

Deletes an item from a DynamoDB table.

```go
func (c *Client) DeleteItem(ctx context.Context, tableName string, key Key) error
```

**Example:**

```go
err := client.DeleteItem(ctx, "users", dynamodb.Key{"pk": "user-123"})
```

### DeleteItemIfExists

Deletes an item only if it exists. Returns `ErrConditionalCheckFailed` if the item doesn't exist.

```go
func (c *Client) DeleteItemIfExists(ctx context.Context, tableName string, key Key, pkAttribute string) error
```

**Example:**

```go
err := client.DeleteItemIfExists(ctx, "users", dynamodb.Key{"pk": "user-123"}, "pk")
```

## Batch Operations

### BatchWriteItems

Writes multiple items to one or more tables. Automatically handles the 25-item limit per batch.

```go
func (c *Client) BatchWriteItems(ctx context.Context, tableName string, items []any) ([]any, error)
```

**Example:**

```go
items := []any{user1, user2, user3}
unprocessed, err := client.BatchWriteItems(ctx, "users", items)
if len(unprocessed) > 0 {
    log.Printf("%d items were not processed", len(unprocessed))
}
```

### BatchDeleteItems

Deletes multiple items from a table. Automatically handles the 25-item limit per batch.

```go
func (c *Client) BatchDeleteItems(ctx context.Context, tableName string, keys []Key) ([]Key, error)
```

**Example:**

```go
keys := []dynamodb.Key{ {"pk": "user-1"}, {"pk": "user-2"} }
unprocessed, err := client.BatchDeleteItems(ctx, "users", keys)
```

## Query Operations

### Query

Performs a query operation on a DynamoDB table.

```go
func (c *Client) Query(ctx context.Context, tableName string, keyCondition expression.KeyConditionBuilder, opts ...QueryOption) (*QueryResult, error)
```

**Query Options:**

| Option                     | Description                            |
| -------------------------- | -------------------------------------- |
| `WithIndex(indexName)`     | Query a secondary index                |
| `WithLimit(limit)`         | Limit number of items returned         |
| `WithScanForward(forward)` | Set scan direction (true=ascending)    |
| `WithFilter(filter)`       | Add filter expression                  |
| `WithProjection(attrs...)` | Specify attributes to return           |
| `WithConsistentRead()`     | Enable consistent reads                |
| `WithStartKey(key)`        | Set exclusive start key for pagination |

**Example:**

```go
// Query by partition key
keyExpr := expression.Key("pk").Equal(expression.Value("user-123"))
result, err := client.Query(ctx, "users", keyExpr)

// Query with sort key condition
keyExpr := expression.KeyAnd(
    expression.Key("pk").Equal(expression.Value("user-123")),
    expression.Key("sk").BeginsWith("order-"),
)
result, err := client.Query(ctx, "orders", keyExpr)

// Query with options
result, err := client.Query(ctx, "users", keyExpr,
    dynamodb.WithIndex("email-index"),
    dynamodb.WithLimit(10),
    dynamodb.WithScanForward(false),
)

var users []User
err = result.Unmarshal(&users)
```

### QueryAll

Performs a query and automatically paginates through all results. Use with caution on large datasets.

```go
func (c *Client) QueryAll(ctx context.Context, tableName string, keyCondition expression.KeyConditionBuilder, dest any, opts ...QueryOption) error
```

**Example:**

```go
keyExpr := expression.Key("pk").Equal(expression.Value("user-123"))
var users []User
err := client.QueryAll(ctx, "users", keyExpr, &users)
```

## Scan Operations

### Scan

Performs a scan operation on a DynamoDB table. Scans read every item in the table - use Query when possible.

```go
func (c *Client) Scan(ctx context.Context, tableName string, opts ...ScanOption) (*ScanResult, error)
```

**Scan Options:**

| Option                                     | Description                   |
| ------------------------------------------ | ----------------------------- |
| `WithScanLimit(limit)`                     | Limit number of items scanned |
| `WithScanFilter(filter)`                   | Add filter expression         |
| `WithScanProjection(attrs...)`             | Specify attributes to return  |
| `WithScanStartKey(key)`                    | Set exclusive start key       |
| `WithParallelScan(segment, totalSegments)` | Configure parallel scanning   |

**Example:**

```go
result, err := client.Scan(ctx, "users")
var users []User
err = result.Unmarshal(&users)

// Scan with filter
filter := expression.Name("status").Equal(expression.Value("active"))
result, err := client.Scan(ctx, "users", dynamodb.WithScanFilter(filter))

// Parallel scan
result, err := client.Scan(ctx, "users", dynamodb.WithParallelScan(0, 4))
```

### ScanAll

Performs a scan and automatically paginates through all results. Use with extreme caution on large tables.

```go
func (c *Client) ScanAll(ctx context.Context, tableName string, dest any, opts ...ScanOption) error
```

**Example:**

```go
var users []User
err := client.ScanAll(ctx, "users", &users)
```

### ScanCallback

Iterates through all items in a table and calls a callback for each batch.

```go
func (c *Client) ScanCallback(ctx context.Context, tableName string, callback func([]map[string]types.AttributeValue) error, opts ...ScanOption) error
```

**Example:**

```go
err := client.ScanCallback(ctx, "users", func(items []map[string]types.AttributeValue) error {
    var users []User
    attributevalue.UnmarshalListOfMaps(items, &users)
    for _, user := range users {
        process(user)
    }
    return nil
})
```

## Error Handling

### Sentinel Errors

| Error                              | Description            |
| ---------------------------------- | ---------------------- |
| `ErrItemNotFound`                  | Item does not exist    |
| `ErrTableNotFound`                 | Table does not exist   |
| `ErrConditionalCheckFailed`        | Condition check failed |
| `ErrProvisionedThroughputExceeded` | Throughput exceeded    |

### Error Handling Example

```go
var user User
err := client.GetItem(ctx, "users", key, &user)
if err != nil {
    if errors.Is(err, dynamodb.ErrItemNotFound) {
        // Create new user
    } else if errors.Is(err, dynamodb.ErrTableNotFound) {
        // Table doesn't exist - fatal error
        return err
    } else if errors.Is(err, dynamodb.ErrProvisionedThroughputExceeded) {
        // Backoff and retry
        time.Sleep(time.Second)
        return retry()
    }
    return err
}
```

## Testing

```go
type mockDynamoDBAPI struct {
    dynamodb.API
}

func (m *mockDynamoDBAPI) GetItem(ctx context.Context, params *dynamodbsdk.GetItemInput, optFns ...func(*dynamodbsdk.Options)) (*dynamodbsdk.GetItemOutput, error) {
    return &dynamodbsdk.GetItemOutput{
        Item: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "user-123"},
        },
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockDynamoDBAPI{}
    client := dynamodb.NewClientWithAPI(mock)
    // Use client in tests...
}
```
