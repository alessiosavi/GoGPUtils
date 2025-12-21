package sqs

import (
	"context"
	"errors"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const serviceName = "sqs"

// API defines the SQS operations used by this package.
// This interface enables testing with mocks.
type API interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	DeleteMessageBatch(ctx context.Context, params *sqs.DeleteMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageBatchOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
	GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)
	ChangeMessageVisibility(ctx context.Context, params *sqs.ChangeMessageVisibilityInput, optFns ...func(*sqs.Options)) (*sqs.ChangeMessageVisibilityOutput, error)
	PurgeQueue(ctx context.Context, params *sqs.PurgeQueueInput, optFns ...func(*sqs.Options)) (*sqs.PurgeQueueOutput, error)
}

// Client provides SQS operations.
type Client struct {
	api API
	cfg *aws.Config
}

// NewClient creates a new SQS client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := sqs.NewClient(cfg)
func NewClient(cfg *aws.Config) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	sqsClient := sqs.NewFromConfig(cfg.AWS())

	return &Client{
		api: sqsClient,
		cfg: cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockSQSAPI{}
//	client := sqs.NewClientWithAPI(mock)
func NewClientWithAPI(api API) *Client {
	return &Client{
		api: api,
	}
}

// API returns the underlying SQS API for direct SDK access.
func (c *Client) API() API {
	return c.api
}

// Common SQS errors.
var (
	// ErrQueueNotFound is returned when a queue does not exist.
	ErrQueueNotFound = errors.New("sqs: queue not found")

	// ErrMessageNotFound is returned when a message does not exist.
	ErrMessageNotFound = errors.New("sqs: message not found")

	// ErrInvalidReceiptHandle is returned when a receipt handle is invalid.
	ErrInvalidReceiptHandle = errors.New("sqs: invalid receipt handle")
)

// isQueueNotFound checks if the error indicates a queue not found condition.
func isQueueNotFound(err error) bool {
	if err == nil {
		return false
	}

	var qne *types.QueueDoesNotExist
	if errors.As(err, &qne) {
		return true
	}

	return strings.Contains(err.Error(), "QueueDoesNotExist") || strings.Contains(err.Error(), "does not exist")
}
