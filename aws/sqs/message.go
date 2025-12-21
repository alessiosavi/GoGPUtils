package sqs

import (
	"context"
	"html"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Message represents an SQS message.
type Message struct {
	ID            string
	Body          string
	ReceiptHandle string
	Attributes    map[string]string
	MD5OfBody     string
}

// SendOption configures SendMessage operations.
type SendOption func(*sendOptions)

type sendOptions struct {
	delaySeconds      int32
	messageAttributes map[string]types.MessageAttributeValue
	groupID           string
	deduplicationID   string
}

// WithDelaySeconds sets the message delay in seconds.
//
// Example:
//
//	msgID, err := client.SendMessage(ctx, queueURL, body, sqs.WithDelaySeconds(60))
func WithDelaySeconds(seconds int32) SendOption {
	return func(o *sendOptions) {
		o.delaySeconds = seconds
	}
}

// WithMessageAttributes sets message attributes.
//
// Example:
//
//	attrs := map[string]types.MessageAttributeValue{
//	    "Author": {DataType: aws.String("String"), StringValue: aws.String("John")},
//	}
//	msgID, err := client.SendMessage(ctx, queueURL, body, sqs.WithMessageAttributes(attrs))
func WithMessageAttributes(attrs map[string]types.MessageAttributeValue) SendOption {
	return func(o *sendOptions) {
		o.messageAttributes = attrs
	}
}

// WithMessageGroupID sets the message group ID for FIFO queues.
//
// Example:
//
//	msgID, err := client.SendMessage(ctx, queueURL, body, sqs.WithMessageGroupID("group-1"))
func WithMessageGroupID(groupID string) SendOption {
	return func(o *sendOptions) {
		o.groupID = groupID
	}
}

// WithDeduplicationID sets the deduplication ID for FIFO queues.
//
// Example:
//
//	msgID, err := client.SendMessage(ctx, queueURL, body, sqs.WithDeduplicationID("unique-id"))
func WithDeduplicationID(id string) SendOption {
	return func(o *sendOptions) {
		o.deduplicationID = id
	}
}

// SendMessage sends a message to an SQS queue.
// Returns the message ID on success.
//
// Example:
//
//	msgID, err := client.SendMessage(ctx, queueURL, "Hello World")
func (c *Client) SendMessage(ctx context.Context, queueURL, body string, opts ...SendOption) (string, error) {
	if queueURL == "" {
		return "", aws.ErrEmptyQueue
	}

	options := &sendOptions{}
	for _, opt := range opts {
		opt(options)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    awssdk.String(queueURL),
		MessageBody: awssdk.String(body),
	}

	if options.delaySeconds > 0 {
		input.DelaySeconds = options.delaySeconds
	}

	if len(options.messageAttributes) > 0 {
		input.MessageAttributes = options.messageAttributes
	}

	if options.groupID != "" {
		input.MessageGroupId = awssdk.String(options.groupID)
	}

	if options.deduplicationID != "" {
		input.MessageDeduplicationId = awssdk.String(options.deduplicationID)
	}

	output, err := c.api.SendMessage(ctx, input)
	if err != nil {
		if isQueueNotFound(err) {
			return "", ErrQueueNotFound
		}

		return "", aws.WrapError(serviceName, "SendMessage", err)
	}

	return awssdk.ToString(output.MessageId), nil
}

// SendMessageBatch sends multiple messages to an SQS queue.
// Automatically handles the 10-message limit per batch.
// Returns the IDs of successfully sent messages and any errors.
//
// Example:
//
//	messages := []string{"msg1", "msg2", "msg3"}
//	successful, failed, err := client.SendMessageBatch(ctx, queueURL, messages)
func (c *Client) SendMessageBatch(ctx context.Context, queueURL string, bodies []string) ([]string, []BatchError, error) {
	if queueURL == "" {
		return nil, nil, aws.ErrEmptyQueue
	}

	if len(bodies) == 0 {
		return nil, nil, nil
	}

	const maxBatchSize = 10

	var (
		successful []string
		failed     []BatchError
	)

	for i := 0; i < len(bodies); i += maxBatchSize {
		end := min(i+maxBatchSize, len(bodies))

		batch := bodies[i:end]
		entries := make([]types.SendMessageBatchRequestEntry, len(batch))

		for j, body := range batch {
			id := generateBatchID(i + j)
			entries[j] = types.SendMessageBatchRequestEntry{
				Id:          awssdk.String(id),
				MessageBody: awssdk.String(body),
			}
		}

		output, err := c.api.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			QueueUrl: awssdk.String(queueURL),
			Entries:  entries,
		})
		if err != nil {
			return successful, failed, aws.WrapError(serviceName, "SendMessageBatch", err)
		}

		for _, s := range output.Successful {
			successful = append(successful, awssdk.ToString(s.MessageId))
		}

		for _, f := range output.Failed {
			failed = append(failed, BatchError{
				ID:      awssdk.ToString(f.Id),
				Code:    awssdk.ToString(f.Code),
				Message: awssdk.ToString(f.Message),
			})
		}
	}

	return successful, failed, nil
}

// BatchError represents an error in a batch operation.
type BatchError struct {
	ID      string
	Code    string
	Message string
}

// ReceiveOption configures ReceiveMessages operations.
type ReceiveOption func(*receiveOptions)

type receiveOptions struct {
	maxMessages             int32
	visibilityTimeout       int32
	waitTimeSeconds         int32
	attributeNames          []types.QueueAttributeName
	messageAttrNames        []string
	receiveRequestAttemptID string
}

// WithMaxMessages sets the maximum number of messages to receive.
// Default is 1, maximum is 10.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL, sqs.WithMaxMessages(10))
func WithMaxMessages(n int32) ReceiveOption {
	return func(o *receiveOptions) {
		o.maxMessages = n
	}
}

// WithVisibilityTimeout sets the visibility timeout in seconds.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL, sqs.WithVisibilityTimeout(60))
func WithVisibilityTimeout(seconds int32) ReceiveOption {
	return func(o *receiveOptions) {
		o.visibilityTimeout = seconds
	}
}

// WithWaitTimeSeconds enables long polling.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL, sqs.WithWaitTimeSeconds(20))
func WithWaitTimeSeconds(seconds int32) ReceiveOption {
	return func(o *receiveOptions) {
		o.waitTimeSeconds = seconds
	}
}

// WithAttributeNames specifies which system attributes to retrieve.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL,
//	    sqs.WithAttributeNames(types.QueueAttributeNameAll))
func WithAttributeNames(names ...types.QueueAttributeName) ReceiveOption {
	return func(o *receiveOptions) {
		o.attributeNames = names
	}
}

// WithMessageAttributeNames specifies which message attributes to retrieve.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL,
//	    sqs.WithMessageAttributeNames("Author", "Timestamp"))
func WithMessageAttributeNames(names ...string) ReceiveOption {
	return func(o *receiveOptions) {
		o.messageAttrNames = names
	}
}

// ReceiveMessages receives messages from an SQS queue.
// Messages are automatically HTML unescaped.
//
// Example:
//
//	messages, err := client.ReceiveMessages(ctx, queueURL,
//	    sqs.WithMaxMessages(10),
//	    sqs.WithWaitTimeSeconds(20),
//	)
//	for _, msg := range messages {
//	    fmt.Println(msg.Body)
//	    // Process message...
//	    client.DeleteMessage(ctx, queueURL, msg.ReceiptHandle)
//	}
func (c *Client) ReceiveMessages(ctx context.Context, queueURL string, opts ...ReceiveOption) ([]Message, error) {
	if queueURL == "" {
		return nil, aws.ErrEmptyQueue
	}

	options := &receiveOptions{
		maxMessages: 1,
	}
	for _, opt := range opts {
		opt(options)
	}

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            awssdk.String(queueURL),
		MaxNumberOfMessages: options.maxMessages,
	}

	if options.visibilityTimeout > 0 {
		input.VisibilityTimeout = options.visibilityTimeout
	}

	if options.waitTimeSeconds > 0 {
		input.WaitTimeSeconds = options.waitTimeSeconds
	}

	if len(options.attributeNames) > 0 {
		input.AttributeNames = options.attributeNames
	}

	if len(options.messageAttrNames) > 0 {
		input.MessageAttributeNames = options.messageAttrNames
	}

	output, err := c.api.ReceiveMessage(ctx, input)
	if err != nil {
		if isQueueNotFound(err) {
			return nil, ErrQueueNotFound
		}

		return nil, aws.WrapError(serviceName, "ReceiveMessage", err)
	}

	messages := make([]Message, len(output.Messages))
	for i, msg := range output.Messages {
		// Unescape HTML entities in the body
		body := awssdk.ToString(msg.Body)
		body = html.UnescapeString(body)

		attrs := make(map[string]string)
		for k, v := range msg.Attributes {
			attrs[string(k)] = v
		}

		messages[i] = Message{
			ID:            awssdk.ToString(msg.MessageId),
			Body:          body,
			ReceiptHandle: awssdk.ToString(msg.ReceiptHandle),
			Attributes:    attrs,
			MD5OfBody:     awssdk.ToString(msg.MD5OfBody),
		}
	}

	return messages, nil
}

// DeleteMessage deletes a message from an SQS queue.
//
// Example:
//
//	err := client.DeleteMessage(ctx, queueURL, message.ReceiptHandle)
func (c *Client) DeleteMessage(ctx context.Context, queueURL, receiptHandle string) error {
	if queueURL == "" {
		return aws.ErrEmptyQueue
	}

	if receiptHandle == "" {
		return ErrInvalidReceiptHandle
	}

	_, err := c.api.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      awssdk.String(queueURL),
		ReceiptHandle: awssdk.String(receiptHandle),
	})
	if err != nil {
		if isQueueNotFound(err) {
			return ErrQueueNotFound
		}

		return aws.WrapError(serviceName, "DeleteMessage", err)
	}

	return nil
}

// DeleteMessageBatch deletes multiple messages from an SQS queue.
// Automatically handles the 10-message limit per batch.
//
// Example:
//
//	handles := []string{msg1.ReceiptHandle, msg2.ReceiptHandle}
//	successful, failed, err := client.DeleteMessageBatch(ctx, queueURL, handles)
func (c *Client) DeleteMessageBatch(ctx context.Context, queueURL string, receiptHandles []string) ([]string, []BatchError, error) {
	if queueURL == "" {
		return nil, nil, aws.ErrEmptyQueue
	}

	if len(receiptHandles) == 0 {
		return nil, nil, nil
	}

	const maxBatchSize = 10

	var (
		successful []string
		failed     []BatchError
	)

	for i := 0; i < len(receiptHandles); i += maxBatchSize {
		end := min(i+maxBatchSize, len(receiptHandles))

		batch := receiptHandles[i:end]
		entries := make([]types.DeleteMessageBatchRequestEntry, len(batch))

		for j, handle := range batch {
			id := generateBatchID(i + j)
			entries[j] = types.DeleteMessageBatchRequestEntry{
				Id:            awssdk.String(id),
				ReceiptHandle: awssdk.String(handle),
			}
		}

		output, err := c.api.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
			QueueUrl: awssdk.String(queueURL),
			Entries:  entries,
		})
		if err != nil {
			return successful, failed, aws.WrapError(serviceName, "DeleteMessageBatch", err)
		}

		for _, s := range output.Successful {
			successful = append(successful, awssdk.ToString(s.Id))
		}

		for _, f := range output.Failed {
			failed = append(failed, BatchError{
				ID:      awssdk.ToString(f.Id),
				Code:    awssdk.ToString(f.Code),
				Message: awssdk.ToString(f.Message),
			})
		}
	}

	return successful, failed, nil
}

// GetQueueURL resolves a queue name to a queue URL.
//
// Example:
//
//	queueURL, err := client.GetQueueURL(ctx, "my-queue")
func (c *Client) GetQueueURL(ctx context.Context, queueName string) (string, error) {
	if queueName == "" {
		return "", aws.ErrEmptyQueue
	}

	output, err := c.api.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: awssdk.String(queueName),
	})
	if err != nil {
		if isQueueNotFound(err) {
			return "", ErrQueueNotFound
		}

		return "", aws.WrapError(serviceName, "GetQueueUrl", err)
	}

	return awssdk.ToString(output.QueueUrl), nil
}

// ChangeMessageVisibility changes the visibility timeout of a message.
//
// Example:
//
//	err := client.ChangeMessageVisibility(ctx, queueURL, receiptHandle, 300)
func (c *Client) ChangeMessageVisibility(ctx context.Context, queueURL, receiptHandle string, timeoutSeconds int32) error {
	if queueURL == "" {
		return aws.ErrEmptyQueue
	}

	if receiptHandle == "" {
		return ErrInvalidReceiptHandle
	}

	_, err := c.api.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          awssdk.String(queueURL),
		ReceiptHandle:     awssdk.String(receiptHandle),
		VisibilityTimeout: timeoutSeconds,
	})
	if err != nil {
		return aws.WrapError(serviceName, "ChangeMessageVisibility", err)
	}

	return nil
}

// PurgeQueue deletes all messages from a queue.
// This operation can only be performed once every 60 seconds.
//
// Example:
//
//	err := client.PurgeQueue(ctx, queueURL)
func (c *Client) PurgeQueue(ctx context.Context, queueURL string) error {
	if queueURL == "" {
		return aws.ErrEmptyQueue
	}

	_, err := c.api.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: awssdk.String(queueURL),
	})
	if err != nil {
		if isQueueNotFound(err) {
			return ErrQueueNotFound
		}

		return aws.WrapError(serviceName, "PurgeQueue", err)
	}

	return nil
}

// generateBatchID generates a unique ID for batch operations.
func generateBatchID(index int) string {
	// Simple numeric ID - AWS allows up to 80 characters
	const digits = "0123456789"
	if index < 10 {
		return string(digits[index])
	}

	result := make([]byte, 0, 10)
	for index > 0 {
		result = append([]byte{digits[index%10]}, result...)
		index /= 10
	}

	return string(result)
}
