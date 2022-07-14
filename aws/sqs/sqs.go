package sqs

import (
	"context"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
		sqsClient = sqs.New(sqs.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	})
}

func GetMessage(queueName string) ([]types.Message, error) {
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
		break
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
