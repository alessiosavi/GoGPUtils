// Package testutil provides utilities for testing AWS helpers with LocalStack.
package testutil

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// TestConfig holds configuration for integration tests.
type TestConfig struct {
	AWSConfig *aws.Config
	Endpoint  string
	Region    string
}

const (
	// DefaultLocalStackEndpoint is the default LocalStack endpoint.
	DefaultLocalStackEndpoint = "http://localhost:4566"

	// DefaultRegion is the default AWS region for tests.
	DefaultRegion = "us-east-1"

	// TestAccessKeyID is the test AWS access key ID.
	TestAccessKeyID = "test"

	// TestSecretAccessKey is the test AWS secret access key.
	TestSecretAccessKey = "test"
)

// LocalStackEndpoint returns the LocalStack endpoint from environment or default.
func LocalStackEndpoint() string {
	if endpoint := os.Getenv("LOCALSTACK_ENDPOINT"); endpoint != "" {
		return endpoint
	}

	return DefaultLocalStackEndpoint
}

// SkipIfNoLocalStack skips the test if LocalStack is not available.
func SkipIfNoLocalStack(t *testing.T) {
	t.Helper()

	endpoint := LocalStackEndpoint()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/_localstack/health", nil)
	if err != nil {
		t.Skipf("LocalStack not available: %v", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("LocalStack not available at %s: %v", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Skipf("LocalStack health check failed: status %d", resp.StatusCode)
	}
}

// MustLoadConfig loads AWS configuration for LocalStack tests.
// It fails the test if configuration cannot be loaded.
func MustLoadConfig(t *testing.T) *TestConfig {
	t.Helper()

	cfg, err := LoadConfig(context.Background())
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	return &TestConfig{
		AWSConfig: cfg,
		Endpoint:  LocalStackEndpoint(),
		Region:    DefaultRegion,
	}
}

// LoadConfig loads AWS configuration for LocalStack.
func LoadConfig(ctx context.Context) (*aws.Config, error) {
	endpoint := LocalStackEndpoint()

	// Use custom endpoint resolver for LocalStack
	customResolver := awssdk.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (awssdk.Endpoint, error) {
			return awssdk.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
				PartitionID:       "aws",
				SigningRegion:     DefaultRegion,
			}, nil
		},
	)

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(DefaultRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			TestAccessKeyID,
			TestSecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &aws.Config{Config: cfg}, nil
}

// UniqueID generates a unique identifier for test resources.
func UniqueID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

// UniqueBucketName generates a unique S3 bucket name.
func UniqueBucketName() string {
	return UniqueID("test-bucket")
}

// UniqueTableName generates a unique DynamoDB table name.
func UniqueTableName() string {
	return UniqueID("test-table")
}

// UniqueQueueName generates a unique SQS queue name.
func UniqueQueueName() string {
	return UniqueID("test-queue")
}

// UniqueSecretName generates a unique secret name.
func UniqueSecretName() string {
	return UniqueID("test-secret")
}

// UniqueParameterName generates a unique SSM parameter name.
func UniqueParameterName() string {
	return "/" + UniqueID("test/param")
}
