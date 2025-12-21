// Package dynamodb provides helpers for Amazon DynamoDB operations.
//
// # Features
//
//   - Item operations: Get, Put, Delete, Update with automatic marshaling
//   - Batch operations: BatchGet, BatchWrite with automatic chunking
//   - Scan and Query with pagination support
//   - Table operations: Create, Delete, WaitForActive
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
//	client := dynamodb.NewClient(cfg)
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
//	err := client.PutItem(ctx, "users", user)
//
//	// Get an item
//	var result User
//	err := client.GetItem(ctx, "users", dynamodb.Key{"pk": "user-123"}, &result)
//
//	// Delete an item
//	err := client.DeleteItem(ctx, "users", dynamodb.Key{"pk": "user-123"})
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockDynamoDBAPI{}
//	client := dynamodb.NewClientWithAPI(mock)
package dynamodb
