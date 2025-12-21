// Package sqs provides helpers for Amazon SQS operations.
//
// # Features
//
//   - Send and receive messages
//   - Batch operations with automatic chunking
//   - Message attribute handling
//   - Queue URL resolution from name
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
//	client := sqs.NewClient(cfg)
//
// # Basic Operations
//
//	// Send a message
//	msgID, err := client.SendMessage(ctx, queueURL, "Hello World")
//
//	// Receive messages
//	messages, err := client.ReceiveMessages(ctx, queueURL, sqs.WithMaxMessages(10))
//
//	// Delete a message
//	err := client.DeleteMessage(ctx, queueURL, receiptHandle)
//
// # Testing
//
// For testing, use the interface-based client:
//
//	mock := &MockSQSAPI{}
//	client := sqs.NewClientWithAPI(mock)
package sqs
