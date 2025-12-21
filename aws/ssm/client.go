package ssm

import (
	"context"
	"errors"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

const serviceName = "ssm"

// API defines the SSM operations used by this package.
// This interface enables testing with mocks.
type API interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
}

// Client provides SSM Parameter Store operations.
type Client struct {
	api API
	cfg *aws.Config
}

// NewClient creates a new SSM client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := ssm.NewClient(cfg)
func NewClient(cfg *aws.Config) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	ssmClient := ssm.NewFromConfig(cfg.AWS())

	return &Client{
		api: ssmClient,
		cfg: cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockSSMAPI{}
//	client := ssm.NewClientWithAPI(mock)
func NewClientWithAPI(api API) *Client {
	return &Client{
		api: api,
	}
}

// API returns the underlying SSM API for direct SDK access.
func (c *Client) API() API {
	return c.api
}

// Common SSM errors.
var (
	// ErrParameterNotFound is returned when a parameter does not exist.
	ErrParameterNotFound = errors.New("ssm: parameter not found")

	// ErrParameterVersionNotFound is returned when a parameter version does not exist.
	ErrParameterVersionNotFound = errors.New("ssm: parameter version not found")

	// ErrAccessDenied is returned when access to a parameter is denied.
	ErrAccessDenied = errors.New("ssm: access denied")
)

// isParameterNotFound checks if the error indicates a not found condition.
func isParameterNotFound(err error) bool {
	if err == nil {
		return false
	}

	var pnf *types.ParameterNotFound
	if errors.As(err, &pnf) {
		return true
	}

	return strings.Contains(err.Error(), "ParameterNotFound")
}
