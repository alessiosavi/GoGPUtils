---
title: SQS
parent: AWS Services
nav_order: 3
---

# SQS

The `aws/sqs` package provides helpers for Amazon SQS operations, including sending and receiving messages, batch operations, and queue management.

## Features

- Send and receive messages
- Batch operations with automatic chunking
- Message attribute handling
- Queue URL resolution from name

## Client Creation

```go
import "github.com/alessiosavi/GoGPUtils/aws/sqs"

cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
if err != nil {
    return err
}

client, err := sqs.NewClient(cfg)
if err != nil {
    return err
}
```

## Types

### Message

Represents an SQS message.

```go
type Message struct {
    ID            string
    Body          string
    ReceiptHandle string
    Attributes    map[string]string
    MD5OfBody     string
}
```

### BatchError

Represents an error in a batch operation.

```go
type BatchError struct {
    ID      string
    Code    string
    Message string
}
```

## Send Operations

### SendMessage

Sends a message to an SQS queue. Returns the message ID on success.

```go
func (c *Client) SendMessage(ctx context.Context, queueURL, body string, opts ...SendOption) (string, error)
```

**Send Options:**

| Option                         | Description                          |
| ------------------------------ | ------------------------------------ |
| `WithDelaySeconds(seconds)`    | Set message delay in seconds         |
| `WithMessageAttributes(attrs)` | Set message attributes               |
| `WithMessageGroupID(id)`       | Set message group ID for FIFO queues |
| `WithDeduplicationID(id)`      | Set deduplication ID for FIFO queues |

**Example:**

```go
// Simple send
msgID, err := client.SendMessage(ctx, queueURL, "Hello World")

// With delay
msgID, err := client.SendMessage(ctx, queueURL, "Hello World",
    sqs.WithDelaySeconds(60),
)

// With attributes
attrs := map[string]types.MessageAttributeValue{
    "Author": {DataType: aws.String("String"), StringValue: aws.String("John")},
}
msgID, err := client.SendMessage(ctx, queueURL, "Hello World",
    sqs.WithMessageAttributes(attrs),
)

// FIFO queue
msgID, err := client.SendMessage(ctx, queueURL, "Hello World",
    sqs.WithMessageGroupID("group-1"),
    sqs.WithDeduplicationID("unique-id"),
)
```

### SendMessageBatch

Sends multiple messages to an SQS queue. Automatically handles the 10-message limit per batch.

```go
func (c *Client) SendMessageBatch(ctx context.Context, queueURL string, bodies []string) ([]string, []BatchError, error)
```

**Example:**

```go
messages := []string{"msg1", "msg2", "msg3"}
successful, failed, err := client.SendMessageBatch(ctx, queueURL, messages)

for _, f := range failed {
    log.Printf("Failed to send: %s - %s", f.ID, f.Message)
}
```

## Receive Operations

### ReceiveMessages

Receives messages from an SQS queue. Messages are automatically HTML unescaped.

```go
func (c *Client) ReceiveMessages(ctx context.Context, queueURL string, opts ...ReceiveOption) ([]Message, error)
```

**Receive Options:**

| Option                                | Description                        |
| ------------------------------------- | ---------------------------------- |
| `WithMaxMessages(n)`                  | Maximum messages to receive (1-10) |
| `WithVisibilityTimeout(seconds)`      | Visibility timeout in seconds      |
| `WithWaitTimeSeconds(seconds)`        | Enable long polling                |
| `WithAttributeNames(names...)`        | System attributes to retrieve      |
| `WithMessageAttributeNames(names...)` | Message attributes to retrieve     |

**Example:**

```go
// Receive single message
messages, err := client.ReceiveMessages(ctx, queueURL)

// Receive with long polling
messages, err := client.ReceiveMessages(ctx, queueURL,
    sqs.WithMaxMessages(10),
    sqs.WithWaitTimeSeconds(20),
)

// Process messages
for _, msg := range messages {
    fmt.Println(msg.Body)
    // Process message...
    err := client.DeleteMessage(ctx, queueURL, msg.ReceiptHandle)
    if err != nil {
        log.Printf("Failed to delete message: %v", err)
    }
}
```

## Delete Operations

### DeleteMessage

Deletes a message from an SQS queue.

```go
func (c *Client) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error
```

**Example:**

```go
err := client.DeleteMessage(ctx, queueURL, message.ReceiptHandle)
```

### DeleteMessageBatch

Deletes multiple messages from an SQS queue. Automatically handles the 10-message limit per batch.

```go
func (c *Client) DeleteMessageBatch(ctx context.Context, queueURL string, receiptHandles []string) ([]string, []BatchError, error)
```

**Example:**

```go
handles := []string{msg1.ReceiptHandle, msg2.ReceiptHandle}
successful, failed, err := client.DeleteMessageBatch(ctx, queueURL, handles)
```

## Queue Operations

### GetQueueURL

Resolves a queue name to a queue URL.

```go
func (c *Client) GetQueueURL(ctx context.Context, queueName string) (string, error)
```

**Example:**

```go
queueURL, err := client.GetQueueURL(ctx, "my-queue")
if errors.Is(err, sqs.ErrQueueNotFound) {
    // Queue doesn't exist
}
```

### ChangeMessageVisibility

Changes the visibility timeout of a message.

```go
func (c *Client) ChangeMessageVisibility(ctx context.Context, queueURL, receiptHandle string, timeoutSeconds int32) error
```

**Example:**

```go
// Extend processing time by 5 minutes
err := client.ChangeMessageVisibility(ctx, queueURL, receiptHandle, 300)
```

### PurgeQueue

Deletes all messages from a queue. This operation can only be performed once every 60 seconds.

```go
func (c *Client) PurgeQueue(ctx context.Context, queueURL string) error
```

**Example:**

```go
err := client.PurgeQueue(ctx, queueURL)
```

## Error Handling

### Sentinel Errors

| Error                     | Description            |
| ------------------------- | ---------------------- |
| `ErrQueueNotFound`        | Queue does not exist   |
| `ErrMessageNotFound`      | Message does not exist |
| `ErrInvalidReceiptHandle` | Invalid receipt handle |

### Error Handling Example

```go
messages, err := client.ReceiveMessages(ctx, queueURL)
if err != nil {
    if errors.Is(err, sqs.ErrQueueNotFound) {
        // Queue was deleted
        return err
    }
    return err
}

for _, msg := range messages {
    err := processMessage(msg)
    if err != nil {
        // Message will return to queue after visibility timeout
        continue
    }

    err = client.DeleteMessage(ctx, queueURL, msg.ReceiptHandle)
    if err != nil {
        log.Printf("Failed to delete: %v", err)
    }
}
```

## Testing

```go
type mockSQSAPI struct {
    sqs.API
}

func (m *mockSQSAPI) ReceiveMessage(ctx context.Context, params *sqssdk.ReceiveMessageInput, optFns ...func(*sqssdk.Options)) (*sqssdk.ReceiveMessageOutput, error) {
    return &sqssdk.ReceiveMessageOutput{
        Messages: []types.Message{
            {
                MessageId:     aws.String("msg-1"),
                Body:          aws.String("Hello World"),
                ReceiptHandle: aws.String("receipt-1"),
            },
        },
    }, nil
}

func TestMyFunction(t *testing.T) {
    mock := &mockSQSAPI{}
    client := sqs.NewClientWithAPI(mock)
    // Use client in tests...
}
```
