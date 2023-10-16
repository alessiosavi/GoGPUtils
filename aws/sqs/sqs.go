package sqs

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	guuid "github.com/google/uuid"
	"html"
	"sync"
)

var sqsClient *sqs.Client = nil
var once sync.Once

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		sqsClient = sqs.New(sqs.Options{Credentials: cfg.Credentials, Region: cfg.Region, RetryMaxAttempts: 5, RetryMode: aws.RetryModeAdaptive})
	})
}

func GetMessage(queueName string) ([]sqsTypes.Message, error) {
	url, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return nil, err
	}
	messages, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl: url.QueueUrl,
	})
	if err != nil {
		return nil, err
	}

	for i, message := range messages.Messages {
		*messages.Messages[i].Body = html.UnescapeString(*message.Body)
	}
	return messages.Messages, nil
}

func DeleteMessage(queueName, receiptHandle string) error {
	url, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return err
	}
	_, err = sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      url.QueueUrl,
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func GetQueueURL(queueName string) (string, error) {
	url, err := sqsClient.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *url.QueueUrl, nil
}

func WriteMessage(queueURL, message string) (*sqs.SendMessageOutput, error) {
	return sqsClient.SendMessage(context.Background(), &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &queueURL,
	})
}

func WriteMessages(queueURL string, messages []string) (*sqs.SendMessageBatchOutput, error) {
	var msgs []sqsTypes.SendMessageBatchRequestEntry

	for _, message := range messages {
		guid := guuid.New().String()
		msgs = append(msgs, sqsTypes.SendMessageBatchRequestEntry{
			Id:          &guid,
			MessageBody: &message,
		})
	}

	return sqsClient.SendMessageBatch(context.Background(), &sqs.SendMessageBatchInput{
		Entries:  msgs,
		QueueUrl: &queueURL,
	})
}
