// Package dynamodb provides helpers for Amazon DynamoDB operations.
//
// # Features
//
//   - Item operations: Get, Put, Delete with automatic marshaling
//   - Batch writes and deletes with automatic chunking and lossless raw APIs
//   - Scan and Query with pagination support
//   - Additional SDK operations through API(), including BatchGet, Update, and table operations
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
//	client, err := dynamodb.NewClient(cfg)
//	if err != nil {
//	    return err
//	}
//
// # Basic Operations
//
//	// Define a struct for your items
//	type User struct {
//	    ID    string `dynamodbav:"pk"`
//	    Email string `dynamodbav:"email"`
//	    Name  string `dynamodbav:"name"`
//	}
//
//	// Put an item
//	user := User{ID: "user-123", Email: "alice@example.com", Name: "Alice"}
//	err = client.PutItem(ctx, "users", user)
//
//	// Get an item
//	var result User
//	err = client.GetItem(ctx, "users", dynamodb.Key{"pk": "user-123"}, &result)
//
//	// Delete an item
//	err = client.DeleteItem(ctx, "users", dynamodb.Key{"pk": "user-123"})
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockDynamoDBAPI{}
//	client := dynamodb.NewClientWithAPI(mock)
package dynamodb
