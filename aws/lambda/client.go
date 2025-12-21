package lambda

import (
	"context"
	"errors"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

const serviceName = "lambda"

// API defines the Lambda operations used by this package.
// This interface enables testing with mocks.
type API interface {
	Invoke(ctx context.Context, params *lambda.InvokeInput, optFns ...func(*lambda.Options)) (*lambda.InvokeOutput, error)
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
	GetFunction(ctx context.Context, params *lambda.GetFunctionInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionOutput, error)
	GetFunctionConfiguration(ctx context.Context, params *lambda.GetFunctionConfigurationInput, optFns ...func(*lambda.Options)) (*lambda.GetFunctionConfigurationOutput, error)
	UpdateFunctionCode(ctx context.Context, params *lambda.UpdateFunctionCodeInput, optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionCodeOutput, error)
	UpdateFunctionConfiguration(ctx context.Context, params *lambda.UpdateFunctionConfigurationInput, optFns ...func(*lambda.Options)) (*lambda.UpdateFunctionConfigurationOutput, error)
	DeleteFunction(ctx context.Context, params *lambda.DeleteFunctionInput, optFns ...func(*lambda.Options)) (*lambda.DeleteFunctionOutput, error)
	ListTags(ctx context.Context, params *lambda.ListTagsInput, optFns ...func(*lambda.Options)) (*lambda.ListTagsOutput, error)
}

// Client provides Lambda operations.
type Client struct {
	api API
	cfg *aws.Config
}

// NewClient creates a new Lambda client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := lambda.NewClient(cfg)
func NewClient(cfg *aws.Config) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	lambdaClient := lambda.NewFromConfig(cfg.AWS())

	return &Client{
		api: lambdaClient,
		cfg: cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockLambdaAPI{}
//	client := lambda.NewClientWithAPI(mock)
func NewClientWithAPI(api API) *Client {
	return &Client{
		api: api,
	}
}

// API returns the underlying Lambda API for direct SDK access.
func (c *Client) API() API {
	return c.api
}

// Common Lambda errors.
var (
	// ErrFunctionNotFound is returned when a function does not exist.
	ErrFunctionNotFound = errors.New("lambda: function not found")

	// ErrInvocationFailed is returned when a function invocation fails.
	ErrInvocationFailed = errors.New("lambda: invocation failed")

	// ErrFunctionError is returned when a function returns an error.
	ErrFunctionError = errors.New("lambda: function returned error")
)

// isResourceNotFound checks if the error indicates a not found condition.
func isResourceNotFound(err error) bool {
	if err == nil {
		return false
	}

	var rnf *types.ResourceNotFoundException
	if errors.As(err, &rnf) {
		return true
	}

	return strings.Contains(err.Error(), "ResourceNotFoundException") || strings.Contains(err.Error(), "Function not found")
}
