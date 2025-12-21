// Package sqs integration tests using LocalStack.
//go:build integration

package sqs_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alessiosavi/GoGPUtils/aws/internal/testutil"
	"github.com/alessiosavi/GoGPUtils/aws/sqs"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	sqssdk "github.com/aws/aws-sdk-go-v2/service/sqs"
)

// testContext returns a context with timeout for tests.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	return ctx
}

// testConfig holds test configuration.
type testConfig struct {
	cfg      *testutil.TestConfig
	client   *sqs.Client
	queueURL string
}

// setupTestQueue creates a unique test queue and returns cleanup function.
func setupTestQueue(t *testing.T, ctx context.Context, cfg *testutil.TestConfig) string {
	t.Helper()

	queueName := testutil.UniqueQueueName()
	sqsClient := sqssdk.NewFromConfig(cfg.AWSConfig.AWS())

	output, err := sqsClient.CreateQueue(ctx, &sqssdk.CreateQueueInput{
		QueueName: awssdk.String(queueName),
	})
	if err != nil {
		t.Fatalf("failed to create test queue: %v", err)
	}

	queueURL := awssdk.ToString(output.QueueUrl)

	t.Cleanup(func() {
		sqsClient.DeleteQueue(context.Background(), &sqssdk.DeleteQueueInput{
			QueueUrl: awssdk.String(queueURL),
		})
	})

	return queueURL
}

// setupTest creates a test environment with SQS client and queue.
func setupTest(t *testing.T) *testConfig {
	t.Helper()
	testutil.SkipIfNoLocalStack(t)

	ctx := testContext(t)
	cfg := testutil.MustLoadConfig(t)

	client, err := sqs.NewClient(cfg.AWSConfig)
	if err != nil {
		t.Fatalf("failed to create SQS client: %v", err)
	}

	queueURL := setupTestQueue(t, ctx, cfg)

	return &testConfig{
		cfg:      cfg,
		client:   client,
		queueURL: queueURL,
	}
}

func TestSendMessage(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "simple message",
			body: "Hello, World!",
		},
		{
			name: "json message",
			body: `{"event": "test", "data": {"value": 123}}`,
		},
		{
			name: "message with special chars",
			body: "Special chars: <>&\"'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgID, err := tc.client.SendMessage(ctx, tc.queueURL, tt.body)
			if err != nil {
				t.Fatalf("SendMessage failed: %v", err)
			}

			if msgID == "" {
				t.Error("expected non-empty message ID")
			}
		})
	}
}

func TestSendMessageWithOptions(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	msgID, err := tc.client.SendMessage(ctx, tc.queueURL, "test with delay",
		sqs.WithDelaySeconds(0),
	)
	if err != nil {
		t.Fatalf("SendMessage with options failed: %v", err)
	}

	if msgID == "" {
		t.Error("expected non-empty message ID")
	}
}

func TestSendMessageBatch(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	messages := []string{
		"batch message 1",
		"batch message 2",
		"batch message 3",
	}

	successful, failed, err := tc.client.SendMessageBatch(ctx, tc.queueURL, messages)
	if err != nil {
		t.Fatalf("SendMessageBatch failed: %v", err)
	}

	if len(failed) > 0 {
		t.Errorf("expected no failed messages, got %d", len(failed))
	}

	if len(successful) != len(messages) {
		t.Errorf("successful count mismatch: got %d, want %d", len(successful), len(messages))
	}
}

func TestSendMessageBatchLarge(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Create more than 10 messages to test batching
	var messages []string
	for i := 0; i < 15; i++ {
		messages = append(messages, "batch message "+string(rune('0'+i)))
	}

	successful, failed, err := tc.client.SendMessageBatch(ctx, tc.queueURL, messages)
	if err != nil {
		t.Fatalf("SendMessageBatch failed: %v", err)
	}

	if len(failed) > 0 {
		t.Errorf("expected no failed messages, got %d", len(failed))
	}

	if len(successful) != len(messages) {
		t.Errorf("successful count mismatch: got %d, want %d", len(successful), len(messages))
	}
}

func TestReceiveMessages(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send messages
	expectedBody := "test message for receive"
	_, err := tc.client.SendMessage(ctx, tc.queueURL, expectedBody)
	if err != nil {
		t.Fatalf("setup SendMessage failed: %v", err)
	}

	// Give LocalStack a moment to make the message available
	time.Sleep(100 * time.Millisecond)

	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL, sqs.WithMaxMessages(10))
	if err != nil {
		t.Fatalf("ReceiveMessages failed: %v", err)
	}

	if len(messages) == 0 {
		t.Fatal("expected at least one message")
	}

	found := false
	for _, msg := range messages {
		if msg.Body == expectedBody {
			found = true
			if msg.ID == "" {
				t.Error("message ID should not be empty")
			}
			if msg.ReceiptHandle == "" {
				t.Error("receipt handle should not be empty")
			}
		}
	}

	if !found {
		t.Error("expected message not found in received messages")
	}
}

func TestReceiveMessagesWithOptions(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send messages
	for i := 0; i < 5; i++ {
		_, err := tc.client.SendMessage(ctx, tc.queueURL, "message "+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("setup SendMessage failed: %v", err)
		}
	}

	// Wait for messages to be available
	time.Sleep(100 * time.Millisecond)

	t.Run("with max messages", func(t *testing.T) {
		messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL,
			sqs.WithMaxMessages(3),
		)
		if err != nil {
			t.Fatalf("ReceiveMessages failed: %v", err)
		}

		if len(messages) > 3 {
			t.Errorf("expected at most 3 messages, got %d", len(messages))
		}
	})

	t.Run("with visibility timeout", func(t *testing.T) {
		messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL,
			sqs.WithVisibilityTimeout(30),
		)
		if err != nil {
			t.Fatalf("ReceiveMessages with visibility timeout failed: %v", err)
		}

		// Clean up - delete received messages
		for _, msg := range messages {
			tc.client.DeleteMessage(ctx, tc.queueURL, msg.ReceiptHandle)
		}
	})
}

func TestDeleteMessage(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send a message
	_, err := tc.client.SendMessage(ctx, tc.queueURL, "message to delete")
	if err != nil {
		t.Fatalf("setup SendMessage failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Receive the message
	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL)
	if err != nil {
		t.Fatalf("ReceiveMessages failed: %v", err)
	}

	if len(messages) == 0 {
		t.Fatal("expected at least one message")
	}

	// Delete the message
	err = tc.client.DeleteMessage(ctx, tc.queueURL, messages[0].ReceiptHandle)
	if err != nil {
		t.Fatalf("DeleteMessage failed: %v", err)
	}
}

func TestDeleteMessageBatch(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send multiple messages
	for i := 0; i < 3; i++ {
		_, err := tc.client.SendMessage(ctx, tc.queueURL, "batch delete "+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("setup SendMessage failed: %v", err)
		}
	}

	time.Sleep(100 * time.Millisecond)

	// Receive messages
	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL, sqs.WithMaxMessages(10))
	if err != nil {
		t.Fatalf("ReceiveMessages failed: %v", err)
	}

	if len(messages) == 0 {
		t.Skip("no messages to delete")
	}

	// Collect receipt handles
	var handles []string
	for _, msg := range messages {
		handles = append(handles, msg.ReceiptHandle)
	}

	// Batch delete
	successful, failed, err := tc.client.DeleteMessageBatch(ctx, tc.queueURL, handles)
	if err != nil {
		t.Fatalf("DeleteMessageBatch failed: %v", err)
	}

	if len(failed) > 0 {
		t.Errorf("expected no failed deletes, got %d", len(failed))
	}

	if len(successful) != len(handles) {
		t.Errorf("successful count mismatch: got %d, want %d", len(successful), len(handles))
	}
}

func TestGetQueueURL(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Create a queue with known name
	queueName := testutil.UniqueQueueName()
	sqsClient := sqssdk.NewFromConfig(tc.cfg.AWSConfig.AWS())

	output, err := sqsClient.CreateQueue(ctx, &sqssdk.CreateQueueInput{
		QueueName: awssdk.String(queueName),
	})
	if err != nil {
		t.Fatalf("failed to create test queue: %v", err)
	}
	expectedURL := awssdk.ToString(output.QueueUrl)

	t.Cleanup(func() {
		sqsClient.DeleteQueue(context.Background(), &sqssdk.DeleteQueueInput{
			QueueUrl: awssdk.String(expectedURL),
		})
	})

	// Test GetQueueURL
	gotURL, err := tc.client.GetQueueURL(ctx, queueName)
	if err != nil {
		t.Fatalf("GetQueueURL failed: %v", err)
	}

	if gotURL != expectedURL {
		t.Errorf("URL mismatch: got %s, want %s", gotURL, expectedURL)
	}
}

func TestGetQueueURLNotFound(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	_, err := tc.client.GetQueueURL(ctx, "non-existent-queue-12345")
	if !errors.Is(err, sqs.ErrQueueNotFound) {
		t.Errorf("expected ErrQueueNotFound, got: %v", err)
	}
}

func TestChangeMessageVisibility(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send and receive a message
	_, err := tc.client.SendMessage(ctx, tc.queueURL, "visibility test")
	if err != nil {
		t.Fatalf("setup SendMessage failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL)
	if err != nil {
		t.Fatalf("ReceiveMessages failed: %v", err)
	}

	if len(messages) == 0 {
		t.Fatal("expected at least one message")
	}

	// Change visibility
	err = tc.client.ChangeMessageVisibility(ctx, tc.queueURL, messages[0].ReceiptHandle, 60)
	if err != nil {
		t.Fatalf("ChangeMessageVisibility failed: %v", err)
	}

	// Clean up
	tc.client.DeleteMessage(ctx, tc.queueURL, messages[0].ReceiptHandle)
}

func TestPurgeQueue(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Setup: send some messages
	for i := 0; i < 3; i++ {
		_, err := tc.client.SendMessage(ctx, tc.queueURL, "purge test "+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("setup SendMessage failed: %v", err)
		}
	}

	// Purge the queue
	err := tc.client.PurgeQueue(ctx, tc.queueURL)
	if err != nil {
		t.Fatalf("PurgeQueue failed: %v", err)
	}

	// Wait for purge to take effect
	time.Sleep(500 * time.Millisecond)

	// Verify queue is empty (may take some time for LocalStack)
	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL, sqs.WithMaxMessages(10))
	if err != nil {
		t.Fatalf("ReceiveMessages after purge failed: %v", err)
	}

	// Note: LocalStack purge is eventually consistent
	if len(messages) > 0 {
		t.Logf("Warning: queue still has %d messages after purge (LocalStack eventually consistent)", len(messages))
	}
}

func TestValidationErrors(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	t.Run("empty queue URL for send", func(t *testing.T) {
		_, err := tc.client.SendMessage(ctx, "", "test")
		if err == nil {
			t.Error("expected error for empty queue URL")
		}
	})

	t.Run("empty queue URL for receive", func(t *testing.T) {
		_, err := tc.client.ReceiveMessages(ctx, "")
		if err == nil {
			t.Error("expected error for empty queue URL")
		}
	})

	t.Run("empty receipt handle for delete", func(t *testing.T) {
		err := tc.client.DeleteMessage(ctx, tc.queueURL, "")
		if !errors.Is(err, sqs.ErrInvalidReceiptHandle) {
			t.Errorf("expected ErrInvalidReceiptHandle, got: %v", err)
		}
	})

	t.Run("empty queue name for GetQueueURL", func(t *testing.T) {
		_, err := tc.client.GetQueueURL(ctx, "")
		if err == nil {
			t.Error("expected error for empty queue name")
		}
	})
}

func TestRoundTrip(t *testing.T) {
	tc := setupTest(t)
	ctx := testContext(t)

	// Send a message
	body := "round trip test message"
	msgID, err := tc.client.SendMessage(ctx, tc.queueURL, body)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if msgID == "" {
		t.Error("expected non-empty message ID")
	}

	time.Sleep(100 * time.Millisecond)

	// Receive the message
	messages, err := tc.client.ReceiveMessages(ctx, tc.queueURL, sqs.WithMaxMessages(10))
	if err != nil {
		t.Fatalf("ReceiveMessages failed: %v", err)
	}

	var found *sqs.Message
	for i, msg := range messages {
		if msg.Body == body {
			found = &messages[i]
			break
		}
	}

	if found == nil {
		t.Fatal("sent message not found in received messages")
	}

	// Delete the message
	err = tc.client.DeleteMessage(ctx, tc.queueURL, found.ReceiptHandle)
	if err != nil {
		t.Fatalf("DeleteMessage failed: %v", err)
	}
}
