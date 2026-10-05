package sqs

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	sqssdk "github.com/aws/aws-sdk-go-v2/service/sqs"
)

// baseGenerateBatchID is copied verbatim from BASE 62dc6b5, except for its name.
func baseGenerateBatchID(index int) string {
	// Simple numeric ID - AWS allows up to 80 characters
	const digits = "0123456789"
	if index < 10 {
		return string(digits[index])
	}

	result := make([]byte, 0, 10)
	for index > 0 {
		result = append([]byte{digits[index%10]}, result...)
		index /= 10
	}

	return string(result)
}

func TestPR14GenerateBatchIDBaseOracle(t *testing.T) {
	check := func(index int) {
		t.Helper()
		if got, want := generateBatchID(index), baseGenerateBatchID(index); got != want {
			t.Fatalf("index=%d: got %q, want %q", index, got, want)
		}
	}
	for index := range 100001 {
		check(index)
	}
	for boundary := 10; ; boundary *= 10 {
		check(boundary - 1)
		check(boundary)
		check(boundary + 1)
		if boundary > math.MaxInt/10 {
			break
		}
	}
	check(math.MaxInt - 1)
	check(math.MaxInt)
	rng := rand.New(rand.NewPCG(14, 2))
	for range 1000 {
		check(rng.Int())
	}
}

func pr14BatchIDPanic(fn func(int) string, index int) (panicType reflect.Type, message string) {
	defer func() {
		if p := recover(); p != nil {
			panicType, message = reflect.TypeOf(p), fmt.Sprint(p)
		}
	}()
	fn(index)
	return
}

func TestPR14GenerateBatchIDNegativePanic(t *testing.T) {
	for _, index := range []int{-1, math.MinInt} {
		t.Run(fmt.Sprintf("index=%d", index), func(t *testing.T) {
			wantType, wantMessage := pr14BatchIDPanic(baseGenerateBatchID, index)
			gotType, gotMessage := pr14BatchIDPanic(generateBatchID, index)
			if wantType == nil || gotType != wantType || gotMessage != wantMessage {
				t.Fatalf("panic = (%v, %q), BASE = (%v, %q)", gotType, gotMessage, wantType, wantMessage)
			}
		})
	}
}

type pr14BatchAPI struct {
	API
	sends   []*sqssdk.SendMessageBatchInput
	deletes []*sqssdk.DeleteMessageBatchInput
}

func (a *pr14BatchAPI) SendMessageBatch(_ context.Context, in *sqssdk.SendMessageBatchInput, _ ...func(*sqssdk.Options)) (*sqssdk.SendMessageBatchOutput, error) {
	a.sends = append(a.sends, in)
	return &sqssdk.SendMessageBatchOutput{}, nil
}

func (a *pr14BatchAPI) DeleteMessageBatch(_ context.Context, in *sqssdk.DeleteMessageBatchInput, _ ...func(*sqssdk.Options)) (*sqssdk.DeleteMessageBatchOutput, error) {
	a.deletes = append(a.deletes, in)
	return &sqssdk.DeleteMessageBatchOutput{}, nil
}

func TestPR14BatchWrappersBaseIDs(t *testing.T) {
	const queueURL = "https://sqs.example.test/queue"
	inputs := make([]string, 25)
	for i := range inputs {
		inputs[i] = fmt.Sprintf("payload-%d", i)
	}
	wantSizes := []int{10, 10, 5}
	t.Run("SendMessageBatch", func(t *testing.T) {
		api := &pr14BatchAPI{}
		client := NewClientWithAPI(api)
		if _, _, err := client.SendMessageBatch(t.Context(), queueURL, inputs); err != nil {
			t.Fatal(err)
		}
		if len(api.sends) != len(wantSizes) {
			t.Fatalf("got %d batches, want %d", len(api.sends), len(wantSizes))
		}
		index := 0
		for batch, request := range api.sends {
			if awssdk.ToString(request.QueueUrl) != queueURL || len(request.Entries) != wantSizes[batch] {
				t.Fatalf("batch %d: queue=%v entries=%d", batch, request.QueueUrl, len(request.Entries))
			}
			for _, entry := range request.Entries {
				if got, want := awssdk.ToString(entry.Id), baseGenerateBatchID(index); got != want {
					t.Errorf("index=%d: ID=%q, BASE=%q", index, got, want)
				}
				if got := awssdk.ToString(entry.MessageBody); got != inputs[index] {
					t.Errorf("index=%d: body=%q, want %q", index, got, inputs[index])
				}
				index++
			}
		}
	})
	t.Run("DeleteMessageBatch", func(t *testing.T) {
		api := &pr14BatchAPI{}
		client := NewClientWithAPI(api)
		if _, _, err := client.DeleteMessageBatch(t.Context(), queueURL, inputs); err != nil {
			t.Fatal(err)
		}
		if len(api.deletes) != len(wantSizes) {
			t.Fatalf("got %d batches, want %d", len(api.deletes), len(wantSizes))
		}
		index := 0
		for batch, request := range api.deletes {
			if awssdk.ToString(request.QueueUrl) != queueURL || len(request.Entries) != wantSizes[batch] {
				t.Fatalf("batch %d: queue=%v entries=%d", batch, request.QueueUrl, len(request.Entries))
			}
			for _, entry := range request.Entries {
				if got, want := awssdk.ToString(entry.Id), baseGenerateBatchID(index); got != want {
					t.Errorf("index=%d: ID=%q, BASE=%q", index, got, want)
				}
				if got := awssdk.ToString(entry.ReceiptHandle); got != inputs[index] {
					t.Errorf("index=%d: receipt=%q, want %q", index, got, inputs[index])
				}
				index++
			}
		}
	})
}
