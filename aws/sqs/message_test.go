package sqs

import (
	"context"
	"reflect"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	sqssdk "github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type receiveAPI struct {
	API
	receive func(*sqssdk.ReceiveMessageInput) (*sqssdk.ReceiveMessageOutput, error)
}

func (m *receiveAPI) ReceiveMessage(_ context.Context, in *sqssdk.ReceiveMessageInput, _ ...func(*sqssdk.Options)) (*sqssdk.ReceiveMessageOutput, error) {
	return m.receive(in)
}

func TestReceiveMessages_PreservesBodyAndAttributes(t *testing.T) {
	body := `{"text":"&quot; &amp; &#34; 日本"}`
	attrs := map[string]types.MessageAttributeValue{
		"text":   {DataType: awssdk.String("String"), StringValue: awssdk.String("&amp;")},
		"binary": {DataType: awssdk.String("Binary"), BinaryValue: []byte{0, 1, 255}},
		"number": {DataType: awssdk.String("Number.int"), StringValue: awssdk.String("9007199254740993")},
	}
	client := NewClientWithAPI(&receiveAPI{receive: func(in *sqssdk.ReceiveMessageInput) (*sqssdk.ReceiveMessageOutput, error) {
		if !reflect.DeepEqual(in.MessageSystemAttributeNames, []types.MessageSystemAttributeName{types.MessageSystemAttributeNameAll}) {
			t.Errorf("system attributes = %v", in.MessageSystemAttributeNames)
		}
		if !reflect.DeepEqual(in.MessageAttributeNames, []string{"All"}) {
			t.Error("custom attributes not requested")
		}
		return &sqssdk.ReceiveMessageOutput{Messages: []types.Message{{Body: &body, MessageId: awssdk.String("id"), ReceiptHandle: awssdk.String("receipt"), MD5OfBody: awssdk.String("md5"), Attributes: map[string]string{"SentTimestamp": "123"}, MessageAttributes: attrs}}}, nil
	}})
	messages, err := client.ReceiveMessages(context.Background(), "queue", WithAttributeNames(types.QueueAttributeNameAll), WithMessageAttributeNames("All"))
	if err != nil || len(messages) != 1 {
		t.Fatalf("messages = %v, err = %v", messages, err)
	}
	m := messages[0]
	if m.Body != body {
		t.Errorf("body = %q, want %q", m.Body, body)
	}
	if m.ID != "id" || m.ReceiptHandle != "receipt" || m.MD5OfBody != "md5" || m.Attributes["SentTimestamp"] != "123" {
		t.Errorf("message fields = %+v", m)
	}
	field := reflect.ValueOf(m).FieldByName("MessageAttributes")
	if !field.IsValid() || !reflect.DeepEqual(field.Interface(), attrs) {
		t.Error("custom message attributes were lost")
	}
}
