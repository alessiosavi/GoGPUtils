// Package lambda provides helpers for AWS Lambda operations.
//
// # Features
//
//   - Invoke functions (sync and async)
//   - List and describe functions
//   - Deploy code from S3 or ZIP
//   - Manage function configuration
//
// # Client Creation
//
// Create a client using AWS configuration:
//
//	cfg, err := aws.LoadConfig(ctx, aws.WithRegion("us-west-2"))
//	if err != nil {
//	    return err
//	}
//
//	client := lambda.NewClient(cfg)
//
// # Basic Operations
//
//	// Invoke a function synchronously
//	response, err := client.Invoke(ctx, "my-function", []byte(`{"key": "value"}`))
//
//	// Invoke asynchronously (fire and forget)
//	err := client.InvokeAsync(ctx, "my-function", payload)
//
//	// List all functions
//	functions, err := client.ListFunctions(ctx)
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockLambdaAPI{}
//	client := lambda.NewClientWithAPI(mock)
package lambda
