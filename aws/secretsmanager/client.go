package secretsmanager

import (
	"context"
	"errors"
	"strings"

	"github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

const serviceName = "secretsmanager"

// API defines the Secrets Manager operations used by this package.
// This interface enables testing with mocks.
type API interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
	ListSecrets(ctx context.Context, params *secretsmanager.ListSecretsInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error)
	DescribeSecret(ctx context.Context, params *secretsmanager.DescribeSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DescribeSecretOutput, error)
	CreateSecret(ctx context.Context, params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.CreateSecretOutput, error)
	UpdateSecret(ctx context.Context, params *secretsmanager.UpdateSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.UpdateSecretOutput, error)
	DeleteSecret(ctx context.Context, params *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.DeleteSecretOutput, error)
}

// Client provides Secrets Manager operations.
type Client struct {
	api API
	cfg *aws.Config
}

// NewClient creates a new Secrets Manager client with the given configuration.
//
// Example:
//
//	cfg, err := aws.LoadConfig(ctx)
//	if err != nil {
//	    return err
//	}
//	client := secretsmanager.NewClient(cfg)
func NewClient(cfg *aws.Config) (*Client, error) {
	if cfg == nil {
		return nil, aws.ErrNilConfig
	}

	smClient := secretsmanager.NewFromConfig(cfg.AWS())

	return &Client{
		api: smClient,
		cfg: cfg,
	}, nil
}

// NewClientWithAPI creates a client with a custom API implementation.
// Useful for testing with mocks.
//
// Example:
//
//	mock := &MockSecretsManagerAPI{}
//	client := secretsmanager.NewClientWithAPI(mock)
func NewClientWithAPI(api API) *Client {
	return &Client{
		api: api,
	}
}

// API returns the underlying Secrets Manager API for direct SDK access.
func (c *Client) API() API {
	return c.api
}

// Common Secrets Manager errors.
var (
	// ErrSecretNotFound is returned when a secret does not exist.
	ErrSecretNotFound = errors.New("secretsmanager: secret not found")

	// ErrSecretDeleted is returned when accessing a deleted secret.
	ErrSecretDeleted = errors.New("secretsmanager: secret is deleted")

	// ErrAccessDenied is returned when access to a secret is denied.
	ErrAccessDenied = errors.New("secretsmanager: access denied")
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

	return strings.Contains(err.Error(), "ResourceNotFoundException")
}

// isInvalidRequest checks if the error indicates an invalid request.
func isInvalidRequest(err error) bool {
	if err == nil {
		return false
	}

	var ir *types.InvalidRequestException
	if errors.As(err, &ir) {
		return true
	}

	return strings.Contains(err.Error(), "InvalidRequestException")
}
