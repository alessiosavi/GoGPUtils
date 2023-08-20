package dynamodbutils

import (
	"context"
	"fmt"
	arrayutils "github.com/alessiosavi/GoGPUtils/array"
	awsutils "github.com/alessiosavi/GoGPUtils/aws"
	"github.com/alessiosavi/GoGPUtils/helper"
	stringutils "github.com/alessiosavi/GoGPUtils/string"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var dynamoClient *dynamodb.Client = nil
var once sync.Once

type UpdateType string

type Update struct {
	Type  UpdateType `json:"type,omitempty"`
	Key   string     `json:"key,omitempty"`
	Value string     `json:"value,omitempty"`
}

var RETRY_ATTEMPT int64 = 1

func init() {
	once.Do(func() {
		cfg, err := awsutils.New()
		if err != nil {
			panic(err)
		}
		dynamoClient = dynamodb.New(dynamodb.Options{Credentials: cfg.Credentials, Region: cfg.Region})
		retryTmp := os.Getenv("DYNAMO_RETRY")
		if !stringutils.IsBlank(retryTmp) {
			RETRY_ATTEMPT, err = strconv.ParseInt(retryTmp, 10, 64)
			if err != nil || RETRY_ATTEMPT > 10 {
				log.Println("WARNING! Error setting DYNAMO_RETRY: ", err, RETRY_ATTEMPT)
				RETRY_ATTEMPT = 1
			}
		}
	})
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
	if _, err := dynamoClient.CreateTable(context.Background(), definition); err != nil {
		return err
	}
	return waitForTable(context.Background(), dynamoClient, *definition.TableName)
}

func WriteItem(tableName string, item interface{}) error {
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

func WriteBatchItem(tableName string, items []interface{}) error {
	var writeReqs []types.WriteRequest
	for i := range items {
		item, err := attributevalue.MarshalMap(items[i])
		if err != nil {
			panic(err)
		}
		writeReqs = append(writeReqs, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}
	requestItems := map[string][]types.WriteRequest{tableName: writeReqs}
	result, err := dynamoClient.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{tableName: writeReqs},
	})
	if err != nil {
		for i := 0; err != nil; i++ {
			if i > int(RETRY_ATTEMPT) {
				return err
			}
			time.Sleep(time.Second * time.Duration(i+1))
			if result != nil && len(result.UnprocessedItems) != 0 {
				requestItems = result.UnprocessedItems
			}

			result, err = dynamoClient.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
				RequestItems: requestItems,
			})
		}
	}
	return err
}
func DeleteTable(tableName string) error {
	_, err := dynamoClient.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})

	return err
}

func GetDocument(tableName string, documentQuery map[string]types.AttributeValue, res interface{}) (*dynamodb.GetItemOutput, error) {
	doc, err := dynamoClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		Key:       documentQuery,
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, err
	}
	//var item map[string]string = make(map[string]string)
	//for attrName, attrValue := range doc.Item {
	//	if av, ok := attrValue.(*types.AttributeValueMemberS); ok {
	//		item[attrName] = av.Value
	//	}
	//}
	//log.Println(item)
	if err = attributevalue.UnmarshalMap(doc.Item, &res); err != nil {
		return nil, err
	}
	return doc, err
}

func ScanAsync(tableName, projectExpression string, ch chan<- []map[string]types.AttributeValue, bar *progressbar.ProgressBar) {
	defer close(ch)
	scan, err := dynamoClient.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:            aws.String(tableName),
		ProjectionExpression: aws.String(projectExpression),
		Limit:                aws.Int32(1000),
	})
	if err != nil {
		panic(err)
	}
	ch <- scan.Items

	for len(scan.LastEvaluatedKey) != 0 {
		scan, err = dynamoClient.Scan(context.Background(), &dynamodb.ScanInput{
			TableName:            aws.String(tableName),
			ProjectionExpression: aws.String(projectExpression),
			ExclusiveStartKey:    scan.LastEvaluatedKey,
			Limit:                aws.Int32(1000),
		})
		if err != nil {
			panic(err)
		}
		bar.ChangeMax(bar.GetMax() + len(scan.Items))
		ch <- scan.Items
	}
	return
}

func DeleteAllItems(tableName, projectExpression string) (*string, error) {
	bar := progressbar.Default(1)
	n := 10
	buffer := make(chan []map[string]types.AttributeValue, n)
	done := make(chan bool, n)
	go ScanAsync(tableName, projectExpression, buffer, bar)
	for i := 0; i < n; i++ {
		go DeleteItems(tableName, buffer, bar, done)
	}
	for i := 0; i < n; i++ {
		<-done
	}

	return nil, nil
}

func DeleteItems(tableName string, datas <-chan []map[string]types.AttributeValue, bar *progressbar.ProgressBar, done chan bool) {
	for data := range datas {
		data0, dataN := arrayutils.SplitEqual(data, 25)
		var writeReqs []types.WriteRequest = nil
		for i := range data0 {
			writeReqs = make([]types.WriteRequest, 0, len(data0[i]))
			for k := range data0[i] {
				writeReqs = append(writeReqs, types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: data0[i][k]}})
			}
			requestItems := map[string][]types.WriteRequest{tableName: writeReqs}
			if _, err := dynamoClient.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{RequestItems: requestItems}); err != nil {
				panic(err)
			}
			bar.Add(25)
		}

		writeReqs = make([]types.WriteRequest, 0, len(dataN))
		for i := range dataN {
			writeReqs = append(writeReqs, types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: dataN[i]}})
		}
		requestItems := map[string][]types.WriteRequest{tableName: writeReqs}
		if _, err := dynamoClient.BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{RequestItems: requestItems}); err != nil {
			panic(err)
		}
		bar.Add(len(dataN))
	}
	done <- true
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
