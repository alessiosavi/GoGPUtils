package dynamodb

import (
	"context"
	"errors"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const serviceName = "dynamodb"

// API defines the DynamoDB operations used by this package.
// This interface enables testing with mocks.
type API interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	DeleteTable(ctx context.Context, params *dynamodb.DeleteTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteTableOutput, error)
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
}

// Client provides DynamoDB operations.
type Client struct {
	api API
	cfg *aws.Config
}

// NewClient creates a new DynamoDB client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := dynamodb.NewClient(cfg)
func NewClient(cfg *aws.Config) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	dynamoClient := dynamodb.NewFromConfig(cfg.AWS())

	return &Client{
		api: dynamoClient,
		cfg: cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockDynamoDBAPI{}
//	client := dynamodb.NewClientWithAPI(mock)
func NewClientWithAPI(api API) *Client {
	return &Client{
		api: api,
	}
}

// API returns the underlying DynamoDB API for direct SDK access.
func (c *Client) API() API {
	return c.api
}

// Common DynamoDB errors.
var (
	// ErrItemNotFound is returned when an item does not exist.
	ErrItemNotFound = errors.New("dynamodb: item not found")

	// ErrTableNotFound is returned when a table does not exist.
	ErrTableNotFound = errors.New("dynamodb: table not found")

	// ErrConditionalCheckFailed is returned when a condition check fails.
	ErrConditionalCheckFailed = errors.New("dynamodb: conditional check failed")

	// ErrProvisionedThroughputExceeded is returned when throughput is exceeded.
	ErrProvisionedThroughputExceeded = errors.New("dynamodb: provisioned throughput exceeded")
)

// isConditionalCheckFailed checks if the error is a conditional check failure.
func isConditionalCheckFailed(err error) bool {
	if err == nil {
		return false
	}

	var ccf *types.ConditionalCheckFailedException
	if errors.As(err, &ccf) {
		return true
	}

	return strings.Contains(err.Error(), "ConditionalCheckFailed")
}

// isResourceNotFound checks if the error indicates a not found condition.
func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}

	var rnf *types.ResourceNotFoundException
	if errors.As(err, &rnf) {
		return true
	}

	return strings.Contains(err.Error(), "ResourceNotFoundException")
}
