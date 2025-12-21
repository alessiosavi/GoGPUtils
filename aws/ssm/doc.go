// Package ssm provides helpers for AWS Systems Manager Parameter Store operations.
//
// # Features
//
//   - Get parameter values (string, string list, secure string)
//   - Get multiple parameters in one call
//   - List parameters with filtering
//   - Automatic decryption of secure strings
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
//	client := ssm.NewClient(cfg)
//
// # Basic Operations
//
//	// Get a parameter value
//	value, err := client.GetParameter(ctx, "/app/config/database_url")
//
//	// Get multiple parameters
//	values, err := client.GetParameters(ctx, []string{"/app/config/a", "/app/config/b"})
//
//	// List parameters by path
//	params, err := client.ListParametersByPath(ctx, "/app/config/")
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockSSMAPI{}
//	client := ssm.NewClientWithAPI(mock)
package ssm
