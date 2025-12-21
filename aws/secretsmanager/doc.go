// Package secretsmanager provides helpers for AWS Secrets Manager operations.
//
// # Features
//
//   - Get secret values (string or binary)
//   - Automatic JSON unmarshaling
//   - List secrets with filtering
//   - Secret rotation support
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
//	client := secretsmanager.NewClient(cfg)
//
// # Basic Operations
//
//	// Get a secret as string
//	value, err := client.GetSecretString(ctx, "my-secret")
//
//	// Get and unmarshal a JSON secret
//	var config DBConfig
//	err := client.GetSecretJSON(ctx, "db-credentials", &config)
//
//	// List all secrets
//	secrets, err := client.ListSecrets(ctx)
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockSecretsManagerAPI{}
//	client := secretsmanager.NewClientWithAPI(mock)
package secretsmanager
