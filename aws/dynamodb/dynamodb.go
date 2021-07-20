package dynamodb

import (
	"context"
	"fmt"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"time"
)

func New() (*dynamodb.Client, error) {
	cfg, err := awsutils.New()
	if err != nil {
		return nil, err
	}
	dynamoClient := dynamodb.New(dynamodb.Options{Credentials: cfg.Credentials, Region: cfg.Region})
	return dynamoClient, nil
}

func waitForTable(ctx context.Context, db *dynamodb.Client, tableName string) error {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return errors.Wrap(err, "timed out while waiting for table to become active")
	}
	return err
}
func CreateTable(definition *dynamodb.CreateTableInput) error {
	dynamoClient, err := New()
	if err != nil {
		return err
	}
	_, err = dynamoClient.CreateTable(context.Background(), definition)
	return waitForTable(context.Background(), dynamoClient, *definition.TableName)
}

func WriteItem(tableName string, item interface{}) error {
	dynamoClient, err := New()
	if err != nil {
		return err
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	fmt.Println(helper.MarshalIndent(av))
	_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})

	return err
}

func DeleteTable(tableName string) error {
	dynamoClient, err := New()
	if err != nil {
		return err
	}
	_, err = dynamoClient.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	return err
}

type UpdateType string

const (
	UpdateTypeSet    UpdateType = "SET"
	UpdateTypeRemove UpdateType = "REMOVE"
	UpdateTypeAdd    UpdateType = "ADD"
)

type Update struct {
	Type  UpdateType `json:"type,omitempty"`
	Key   string     `json:"key,omitempty"`
	Value string     `json:"value,omitempty"`
}

//func UpdateItem(tableName string, key map[string]types.AttributeValue, set map[string]string) error {
//	dynamoClient, err := New()
//	if err != nil {
//		return err
//	}
//
//	dynamoClient.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
//		Key:                         key,
//		TableName:                   aws.String(tableName),
//		AttributeUpdates:            nil,
//		ConditionExpression:         nil,
//		ConditionalOperator:         "",
//		Expected:                    nil,
//		ExpressionAttributeNames:    nil,
//		ExpressionAttributeValues:   nil,
//		ReturnConsumedCapacity:      "",
//		ReturnItemCollectionMetrics: "",
//		ReturnValues:                "",
//		UpdateExpression:            nil,
//	})
//	return nil
//}
