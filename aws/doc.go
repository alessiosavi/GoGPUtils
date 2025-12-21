// Package aws provides idiomatic Go helpers for AWS SDK v2.
//
// # Design Philosophy
//
// This package follows several key principles:
//
//   - No global state: All clients are explicitly created and passed
//   - Context-aware: All operations accept context.Context for cancellation
//   - Testable: Interfaces enable mocking without real AWS calls
//   - Explicit configuration: No magic; all settings are visible
//   - Safe defaults: Sensible retry policies and timeouts out of the box
//   - Minimal abstraction: Helpers augment the SDK, not replace it
//
// # Package Structure
//
// The aws package is organized into focused sub-packages:
//
//   - aws: Core configuration and common utilities
//   - aws/s3: S3 object and bucket operations
//   - aws/dynamodb: DynamoDB item and table operations
//   - aws/sqs: SQS message operations
//   - aws/secretsmanager: Secrets Manager operations
//   - aws/ssm: SSM Parameter Store operations
//   - aws/lambda: Lambda function operations
//
// # Configuration
//
// All clients are created with explicit configuration:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
//	if err != nil {
//	    return err
//	}
//
//	s3Client := s3.NewClient(cfg)
//
// # Error Handling
//
// All operations return wrapped errors with context:
//
//	obj, err := client.GetObject(ctx, bucket, key)
//	if errors.Is(err, s3.ErrObjectNotFound) {
//	    // Handle not found
//	}
//
// # Testing
//
// Each client accepts interfaces that can be mocked:
//
//	mock := &MockS3API{}
//	client := s3.NewClientWithAPI(mock)
//	// Use client in tests
package aws
